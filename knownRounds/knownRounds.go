// Tracks which rounds have been checked and which are unchecked using a bit
// stream.
package knownRounds

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/xx_network/primitives/id"
)

// KnownRounds structure tracks which rounds are known and which are unknown.
// Each bit in bitStream corresponds to a round ID and if it is set, it means
// the round has been checked. All rounds before firstUnchecked are known to be
// checked. All rounds after lastChecked are unknown.
type KnownRounds struct {
	bitStream      uint64Buff // Buffer of check/unchecked rounds
	firstUnchecked id.Round   // ID of the first round that us unchecked
	lastChecked    id.Round   // ID of the last round that is checked
	fuPos          int        // Bit position of firstUnchecked in bitStream
}

// DiskKnownRounds structure is used to as an intermediary to marshal and
// unmarshal KnownRounds.
type DiskKnownRounds struct {
	BitStream                   []uint64
	FirstUnchecked, LastChecked uint64
}

// NewKnownRound creates a new empty KnownRounds in the default state with a
// bit stream of the specified size.
func NewKnownRound(size int) *KnownRounds {
	return &KnownRounds{
		bitStream:      make(uint64Buff, size),
		firstUnchecked: 0,
		lastChecked:    1,
		fuPos:          0,
	}
}

// Marshal returns the JSON encoding of DiskKnownRounds, which contains the
// compressed information from KnownRounds. The bit stream is compressed such
// that the firstUnchecked occurs in the first block of the bit stream.
func (kr *KnownRounds) Marshal() ([]byte, error) {
	// Calculate length of compressed bit stream.
	startPos := kr.getBitStreamPos(kr.firstUnchecked)
	endPos := kr.getBitStreamPos(kr.lastChecked)
	length := kr.bitStream.delta(startPos, endPos)

	// Generate DiskKnownRounds with bit stream of the correct size
	dkr := DiskKnownRounds{
		BitStream:      make([]uint64, length),
		FirstUnchecked: uint64(kr.firstUnchecked),
		LastChecked:    uint64(kr.lastChecked),
	}

	// Copy only the blocks between firstUnchecked and lastChecked to the stream
	startBlock, _ := kr.bitStream.convertLoc(startPos)
	for i := uint(0); i < length; i++ {
		dkr.BitStream[i] = kr.bitStream[(i+startBlock)%uint(len(kr.bitStream))]
	}

	return json.Marshal(dkr)
}

// Unmarshal parses the JSON-encoded data and stores it in the KnownRounds. An
// error is returned if the bit stream data is larger than the KnownRounds bit
// stream.
func (kr *KnownRounds) Unmarshal(data []byte) error {
	// Unmarshal JSON data
	dkr := &DiskKnownRounds{}
	err := json.Unmarshal(data, dkr)
	if err != nil {
		return err
	}

	// Return an error if the passed in bit stream is larger than the new stream
	if len(kr.bitStream) < len(dkr.BitStream) {
		return errors.Errorf("KnownRounds bitStream size of %d is too small "+
			"for passed in bit stream of size %d.",
			len(kr.bitStream), len(dkr.BitStream))
	}

	// Copy values over
	copy(kr.bitStream, dkr.BitStream)
	kr.firstUnchecked = id.Round(dkr.FirstUnchecked)
	kr.lastChecked = id.Round(dkr.LastChecked)
	kr.fuPos = int(dkr.FirstUnchecked % 64)

	return nil
}

// Checked determines if the round has been checked.
func (kr *KnownRounds) Checked(rid id.Round) bool {
	if rid < kr.firstUnchecked {
		return true
	} else if rid > kr.lastChecked {
		return false
	}

	pos := kr.getBitStreamPos(rid)

	return kr.bitStream.get(pos)
}

// Check denotes a round has been checked. If the passed in round occurred after
// the last checked round, then every round between them is set as unchecked and
// the passed in round becomes the last checked round.
func (kr *KnownRounds) Check(rid id.Round) {
	if rid < kr.firstUnchecked {
		return
	}
	pos := kr.getBitStreamPos(rid)

	// Set round as checked
	kr.bitStream.set(pos)

	// If the round ID is newer, then set it as the last checked ID and uncheck
	// all the newly added rounds in the buffer
	if rid > kr.lastChecked {
		kr.bitStream.clearRange(kr.getBitStreamPos(kr.lastChecked+1), pos)
		kr.lastChecked = rid
	}
}

// Forward sets all rounds before the given round ID as checked.
func (kr *KnownRounds) Forward(rid id.Round) {
	if rid > kr.lastChecked {
		kr.firstUnchecked = rid
		kr.lastChecked = rid - 1
		kr.fuPos = int(rid % 64)
	} else if rid >= kr.firstUnchecked {
		for ; kr.bitStream.get(kr.getBitStreamPos(rid)) &&
			rid <= kr.lastChecked; rid++ {
		}
		kr.fuPos = kr.getBitStreamPos(rid)
		kr.firstUnchecked = rid
	}
}

// RangeUnchecked runs the passed function over the range all unchecked round
// IDs up to the passed newestRound to determine if they should be checked.
func (kr *KnownRounds) RangeUnchecked(newestRound id.Round,
	roundCheck func(id id.Round) bool) {

	// If the newest round is in the range of known rounds, then skip checking
	if newestRound < kr.firstUnchecked {
		return
	}
	internalEnd := kr.lastChecked

	// Check all the rounds after the last checked round
	if newestRound >= kr.lastChecked {
		for i := kr.lastChecked; i <= newestRound; i++ {
			if roundCheck(i) {
				kr.Check(i)
			}
		}
	} else {
		internalEnd = newestRound
	}

	// Check all unknown rounds between first unchecked and last checked
	for i := kr.firstUnchecked; i < internalEnd; i++ {
		if !kr.Checked(i) && roundCheck(i) {
			kr.Check(i)
		}
	}
}

// assumption mask << kr and therefore do not have to deal with discontinuity
func (kr *KnownRounds) RangeUncheckedMasked(mask *KnownRounds,
	roundCheck func(id id.Round) bool, maxChecked int) {

	mask.Forward(kr.firstUnchecked)
	subSample, delta := kr.subSample(mask.firstUnchecked, mask.lastChecked)
	fmt.Printf("delta: %d\n", delta)
	fmt.Printf("subSample: %064b\n", subSample)
	result := subSample.implies(mask.bitStream)
	numChecked := 0

	for i := mask.firstUnchecked + id.Round(delta) - 1; i >= mask.firstUnchecked && numChecked < maxChecked; i, numChecked = i-1, numChecked+1 {
		if !result.get(int(i-mask.firstUnchecked)) && roundCheck(i) {
			kr.Check(i)
		}
	}

	for i := kr.firstUnchecked; i < mask.firstUnchecked && numChecked < maxChecked; i, numChecked = i+1, numChecked+1 {
		if !kr.Checked(i) && roundCheck(i) {
			kr.Check(i)
		}
	}
}

func (kr *KnownRounds) subSample(start, end id.Round) (uint64Buff, int) {
	// numBlocks := int(end/64 - start/64)
	numBlocks := kr.bitStream.delta(kr.getBitStreamPos(start), kr.getBitStreamPos(end))
	fmt.Printf("start: %d\n", start)
	fmt.Printf("end: %d\n", end)
	fmt.Printf("numBlocks: %d\n", numBlocks)

	if kr.lastChecked < end {
		end = kr.lastChecked
	}

	buff := kr.bitStream.copy(kr.getBitStreamPos(start), kr.getBitStreamPos(end+1))

	return buff.extend(int(numBlocks)), int(end - start)
}

// Get the position of the bit in the bit stream for the given round ID.
func (kr *KnownRounds) getBitStreamPos(rid id.Round) int {
	var delta int
	if rid < kr.firstUnchecked {
		delta = -int(kr.firstUnchecked - rid)
	} else {
		delta = int(rid - kr.firstUnchecked)
	}

	return kr.fuPos + delta
}
