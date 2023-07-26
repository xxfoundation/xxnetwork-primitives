////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Package knownRounds tracks which rounds have been checked and which are
// unchecked using a bit stream.
package knownRounds

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"

	"gitlab.com/xx_network/primitives/id"
)

type RoundCheckFunc func(id id.Round) bool

// KnownRounds structure tracks which rounds are known and which are unknown.
// Each bit in bitStream corresponds to a round ID and if it is set, it means
// the round has been checked. All rounds before firstUnchecked are known to be
// checked. All rounds after lastChecked are unknown.
type KnownRounds struct {
	bitStream      uint64Buff // Buffer of check/unchecked rounds
	firstUnchecked id.Round   // ID of the first round that us unchecked
	lastChecked    id.Round   // ID of the last round that is checked
	fuPos          int        // The bit position of firstUnchecked in bitStream
}

// DiskKnownRounds structure is used to as an intermediary to marshal and
// unmarshal KnownRounds.
type DiskKnownRounds struct {
	BitStream                   []byte
	FirstUnchecked, LastChecked uint64
}

// NewKnownRound creates a new empty KnownRounds in the default state with a
// bit stream that can hold the given number of rounds.
func NewKnownRound(roundCapacity int) *KnownRounds {
	return &KnownRounds{
		bitStream:      make(uint64Buff, (roundCapacity+63)/64),
		firstUnchecked: 0,
		lastChecked:    0,
		fuPos:          0,
	}
}

// NewFromParts creates a new KnownRounds from the given firstUnchecked,
// lastChecked, fuPos, and uint64 buffer.
func NewFromParts(
	buff []uint64, firstUnchecked, lastChecked id.Round, fuPos int) *KnownRounds {
	return &KnownRounds{
		bitStream:      buff,
		firstUnchecked: firstUnchecked,
		lastChecked:    lastChecked,
		fuPos:          fuPos,
	}
}

// Marshal returns the JSON encoding of DiskKnownRounds, which contains the
// compressed information from KnownRounds. The bit stream is compressed such
// that the firstUnchecked occurs in the first block of the bit stream.
func (kr *KnownRounds) Marshal() []byte {
	// Calculate length of compressed bit stream.
	startPos := kr.getBitStreamPos(kr.firstUnchecked)
	endPos := kr.getBitStreamPos(kr.lastChecked)
	length := kr.bitStream.delta(startPos, endPos)

	// Copy only the blocks between firstUnchecked and lastChecked to the stream
	startBlock, _ := kr.bitStream.convertLoc(startPos)
	bitStream := make(uint64Buff, length)
	for i := 0; i < length; i++ {
		bitStream[i] = kr.bitStream[(i+startBlock)%len(kr.bitStream)]
	}

	// Create new buffer
	buf := bytes.Buffer{}

	// Add firstUnchecked to buffer
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(kr.firstUnchecked))
	buf.Write(b)

	// Add lastChecked to buffer
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(kr.lastChecked))
	buf.Write(b)

	// Add marshaled bitStream to buffer
	buf.Write(bitStream.marshal())

	return buf.Bytes()
}

// Unmarshal parses the JSON-encoded data and stores it in the KnownRounds. An
// error is returned if the bit stream data is larger than the KnownRounds bit
// stream.
func (kr *KnownRounds) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)

	if buf.Len() < 16 {
		return errors.Errorf("KnownRounds Unmarshal: "+
			"size of data %d < %d expected", buf.Len(), 16)
	}

	// Get firstUnchecked and lastChecked and calculate fuPos
	kr.firstUnchecked = id.Round(binary.LittleEndian.Uint64(buf.Next(8)))
	kr.lastChecked = id.Round(binary.LittleEndian.Uint64(buf.Next(8)))
	kr.fuPos = int(kr.firstUnchecked % 64)

	// Unmarshal the bitStream from the rest of the bytes
	bitStream, err := unmarshal(buf.Bytes())
	if err != nil {
		return errors.Errorf("Failed to unmarshal bitstream: %+v", err)
	}

	// Handle the copying in of the bit stream
	if len(kr.bitStream) == 0 {
		// If there is no bitstream, like in the wire representations, then make
		// the size equal to what is coming in
		kr.bitStream = bitStream
	} else if len(kr.bitStream) >= len(bitStream) {
		// If a size already exists and the data fits within it, then copy it
		// into the beginning of the buffer
		copy(kr.bitStream, bitStream)
	} else {
		// If the passed in data is larger than the internal buffer, then return
		// an error
		return errors.Errorf("KnownRounds bitStream size of %d is too small "+
			"for passed in bit stream of size %d.",
			len(kr.bitStream), len(bitStream))
	}

	return nil
}

// KrChanges map contains a list of changes between two KnownRounds bit streams.
// The key is the index of the changed word and the value contains the change.
type KrChanges map[int]uint64

// OutputBuffChanges returns the current KnownRounds' firstUnchecked,
// lastChecked, fuPos, and a list of changes between the given uint64 buffer and
// the current KnownRounds bit stream. An error is returned if the two buffers
// are not of the same length.
func (kr *KnownRounds) OutputBuffChanges(
	old []uint64) (KrChanges, id.Round, id.Round, int, error) {

	// Return an error if they are not the same length
	if len(old) != len(kr.bitStream) {
		return nil, 0, 0, 0, errors.Errorf("length of old buffer %d is "+
			"not the same as length of the current buffer %d",
			len(old), len(kr.bitStream))
	}

	// Create list of changes
	changes := make(KrChanges)
	for i, word := range kr.bitStream {
		if word != old[i] {
			changes[i] = word
		}
	}

	return changes, kr.firstUnchecked, kr.lastChecked, kr.fuPos, nil
}

func (kr KnownRounds) GetFirstUnchecked() id.Round   { return kr.firstUnchecked }
func (kr KnownRounds) GetLastChecked() id.Round      { return kr.lastChecked }
func (kr KnownRounds) GetFuPos() int                 { return kr.fuPos }
func (kr KnownRounds) GetBitStream() []uint64        { return kr.bitStream.deepCopy() }
func (kr KnownRounds) MarshalBitStream1Byte() []byte { return kr.bitStream.marshal1ByteVer2() }
func (kr KnownRounds) MarshalBitStream2Byte() []byte { return kr.bitStream.marshal2BytesVer2() }
func (kr KnownRounds) MarshalBitStream4Byte() []byte { return kr.bitStream.marshal4BytesVer2() }
func (kr KnownRounds) MarshalBitStream8Byte() []byte { return kr.bitStream.marshal8BytesVer2() }

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
// the passed in round becomes the last checked round. Will panic if the buffer
// is not large enough to hold the current data and the new data
func (kr *KnownRounds) Check(rid id.Round) {
	if abs(int(kr.lastChecked-rid))/(len(kr.bitStream)*64) > 0 {
		jww.FATAL.Panicf("Cannot check a round outside the current scope. " +
			"Scope is KnownRounds size more rounds than last checked. A call " +
			"to Forward can be used to fix the scope.")
	}
	kr.check(rid)
}

func (kr *KnownRounds) ForceCheck(rid id.Round) {
	if rid < kr.firstUnchecked {
		return
	} else if kr.lastChecked < rid &&
		int(rid-kr.firstUnchecked) > (len(kr.bitStream)*64) {
		kr.Forward(rid - id.Round(len(kr.bitStream)*64))
	}

	kr.check(rid)
}

// Check denotes a round has been checked. If the passed in round occurred after
// the last checked round, then every round between them is set as unchecked and
// the passed in round becomes the last checked round. Will shift the buffer
// forward, erasing old data, if the buffer is not large enough to hold the new
// checked input
func (kr *KnownRounds) check(rid id.Round) {
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

	if kr.getBitStreamPos(kr.firstUnchecked) == pos {
		if kr.getBitStreamPos(kr.lastChecked) == pos {
			kr.fuPos = kr.getBitStreamPos(rid + 1)
			kr.firstUnchecked = rid + 1
			kr.lastChecked = rid + 1
			kr.bitStream.clear(kr.fuPos)
		} else {
			kr.migrateFirstUnchecked(rid)
		}
	}

	// Handle cases where rid lapse firstUnchecked one or more times.
	if rid > kr.firstUnchecked && (rid-kr.firstUnchecked) >= id.Round(kr.Len()) {
		newFu := rid + 1 - id.Round(kr.Len())
		kr.fuPos = kr.getBitStreamPos(newFu)
		kr.firstUnchecked = rid + 1 - id.Round(kr.Len())
		kr.migrateFirstUnchecked(rid)
	}

	// Set round as checked
	kr.bitStream.set(pos)
}

// abs returns the absolute value of the passed in integer.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// migrateFirstUnchecked moves firstUnchecked to the next unchecked round or
// sets it to lastUnchecked if all rounds are checked.
func (kr *KnownRounds) migrateFirstUnchecked(rid id.Round) {
	for ; kr.bitStream.get(kr.getBitStreamPos(rid)) && rid <= kr.lastChecked; rid++ {
	}
	kr.fuPos = kr.getBitStreamPos(rid)
	kr.firstUnchecked = rid
}

// Forward sets all rounds before the given round ID as checked.
func (kr *KnownRounds) Forward(rid id.Round) {
	if rid > kr.lastChecked {
		kr.firstUnchecked = rid
		kr.lastChecked = rid
		kr.fuPos = int(rid % 64)
	} else if rid > kr.firstUnchecked {
		kr.migrateFirstUnchecked(rid)
	}
}

// RangeUnchecked runs the passed function over all rounds starting with oldest
// unknown and ending with
func (kr *KnownRounds) RangeUnchecked(oldestUnknown id.Round, threshold uint,
	roundCheck func(id id.Round) bool, maxPickups int) (
	earliestRound id.Round, has, unknown []id.Round) {

	newestRound := kr.lastChecked

	// Calculate how far back we should go back to check rounds
	// If the newest round is smaller than the threshold, then our oldest round
	// is zero. But the earliest round ID is 1.
	oldestPossibleEarliestRound := id.Round(1)
	if newestRound > id.Round(threshold) {
		oldestPossibleEarliestRound = newestRound - id.Round(threshold)
	}

	earliestRound = kr.lastChecked + 1
	has = make([]id.Round, 0, maxPickups)

	// If the oldest unknown round is outside the range we are attempting to
	// check, then skip checking
	if oldestUnknown > kr.lastChecked {
		jww.TRACE.Printf(
			"RangeUnchecked: oldestUnknown (%d) > kr.lastChecked (%d)",
			oldestUnknown, kr.lastChecked)
		return oldestUnknown, nil, nil
	}

	// Loop through all rounds from the oldest unknown to the last checked round
	// and check them, if possible
	for i := oldestUnknown; i <= kr.lastChecked; i++ {

		// If the source does not know about the round, set that round as
		// unknown and don't check it
		if !kr.Checked(i) {
			if i < oldestPossibleEarliestRound {
				unknown = append(unknown, i)
			} else if i < earliestRound {
				earliestRound = i
			}
			continue
		}

		// check the round
		hasRound := roundCheck(i)

		// If checking is not complete and the round is earlier than the
		// earliest round, then set it to the earliest round
		if hasRound {
			has = append(has, i)
			// Do not pick up too many messages at once
			if len(has) >= maxPickups {
				nextRound := i + 1
				if (nextRound) < earliestRound {
					earliestRound = nextRound
				}
				break
			}
		}
	}

	// Return the next round
	return earliestRound, has, unknown
}

// RangeUncheckedMasked masks the bit stream with the provided mask.
func (kr *KnownRounds) RangeUncheckedMasked(mask *KnownRounds,
	roundCheck RoundCheckFunc, maxChecked int) {

	kr.RangeUncheckedMaskedRange(mask, roundCheck, 0, math.MaxUint64, maxChecked)
}

// RangeUncheckedMaskedRange masks the bit stream with the provided mask.
func (kr *KnownRounds) RangeUncheckedMaskedRange(mask *KnownRounds,
	roundCheck RoundCheckFunc, start, end id.Round, maxChecked int) {

	numChecked := 0

	if mask.firstUnchecked != mask.lastChecked {
		mask.Forward(kr.firstUnchecked)
		subSample, delta := kr.subSample(mask.firstUnchecked, mask.lastChecked)
		// FIXME: it is inefficient to make a copy of the mask here.
		result := subSample.implies(mask.bitStream)

		for i := mask.firstUnchecked + id.Round(delta) - 1; i >= mask.firstUnchecked && numChecked < maxChecked; i, numChecked = i-1, numChecked+1 {
			if !result.get(int(i-mask.firstUnchecked)) && roundCheck(i) {
				kr.Check(i)
			}
		}
	}

	if start < kr.firstUnchecked {
		start = kr.firstUnchecked
	}

	if end > mask.firstUnchecked {
		end = mask.firstUnchecked
	}

	for i := start; i < end && numChecked < maxChecked; i, numChecked = i+1, numChecked+1 {
		if !kr.Checked(i) && roundCheck(i) {
			kr.Check(i)
		}
	}
}

// subSample returns a sub sample of the KnownRounds buffer from the start to
// end round and its length.
func (kr *KnownRounds) subSample(start, end id.Round) (uint64Buff, int) {
	// Get the number of blocks spanned by the range
	numBlocks := kr.bitStream.delta(kr.getBitStreamPos(start),
		kr.getBitStreamPos(end))

	if start > kr.lastChecked {
		return make(uint64Buff, numBlocks), numBlocks
	}

	copyEnd := end
	if kr.lastChecked < end {
		copyEnd = kr.lastChecked
	}

	// Create a sub sample of the buffer
	buff := kr.bitStream.copy(kr.getBitStreamPos(start),
		kr.getBitStreamPos(copyEnd+1))

	// Return a buffer of the correct size and its length
	return buff.extend(numBlocks), abs(int(end - start))
}

// Truncate returns a subs ample of the KnownRounds buffer from last checked.
func (kr *KnownRounds) Truncate(start id.Round) *KnownRounds {
	if start <= kr.firstUnchecked {
		return kr
	}

	// Return a buffer of the correct size and its length
	newKr := &KnownRounds{
		bitStream:      kr.bitStream.deepCopy(),
		firstUnchecked: kr.firstUnchecked,
		lastChecked:    kr.lastChecked,
		fuPos:          kr.fuPos,
	}

	newKr.migrateFirstUnchecked(start)

	return newKr
}

// Get the position of the bit in the bit stream for the given round ID.
func (kr *KnownRounds) getBitStreamPos(rid id.Round) int {
	var delta int
	if rid < kr.firstUnchecked {
		delta = -int(kr.firstUnchecked - rid)
	} else {
		delta = int(rid - kr.firstUnchecked)
	}

	pos := (kr.fuPos + delta) % kr.Len()
	if pos < 0 {
		return kr.Len() + pos
	}
	return pos

}

// Len returns the max number of round IDs the buffer can hold.
func (kr *KnownRounds) Len() int {
	return len(kr.bitStream) * 64
}
