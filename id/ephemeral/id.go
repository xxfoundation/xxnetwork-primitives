package ephemeral

import (
	"bytes"
	"crypto"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/xx_network/primitives/id"
	"hash"
	"io"
	"math"
	"time"
)

const Period = int64(time.Hour * 24)
const NumOffsets int64 = 1 << 16
const NsPerOffset = Period / NumOffsets

// Ephemeral Ids reserved for specific actions:
// All zero's denote a dummy ID
// All one's denote a payment
var ReservedIDs = []Id{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{1, 1, 1, 1, 1, 1, 1, 1},
}

// Ephemeral ID type alias
type Id [8]byte

// Ephemeral Id object which contains the ID
// and the start and end time for the salt window
type ProtoIdentity struct {
	Id    Id
	Start time.Time
	End   time.Time
}

// Return ephemeral ID as a uint64
func (eid *Id) UInt64() uint64 {
	return binary.BigEndian.Uint64(eid[:])
}

// Return ephemeral ID as an int64
func (eid *Id) Int64() int64 {
	ux := binary.BigEndian.Uint64(eid[:])
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x
}

// Clear an ID down to the correct size
func (eid Id) Clear(size uint) Id {
	newId := Id{}
	var mask uint64 = math.MaxUint64 >> (64 - size)
	maskedId := binary.BigEndian.Uint64(eid[:]) & mask
	binary.BigEndian.PutUint64(newId[:], maskedId)
	return newId
}

// Fill cleared bits of an ID with random data from passed in rng
// Accepts the size of the ID in bits & an RNG reader
func (eid Id) Fill(size uint, rng io.Reader) (Id, error) {
	newId := Id{}
	rand := Id{}
	_, err := rng.Read(rand[:])
	if err != nil {
		return Id{}, err
	}
	var mask uint64 = math.MaxUint64 << size
	maskedRand := mask & binary.BigEndian.Uint64(rand[:])
	maskedEid := maskedRand | binary.BigEndian.Uint64(eid[:])
	binary.BigEndian.PutUint64(newId[:], maskedEid)
	return newId, nil
}

// Load an ephemeral ID from raw bytes
func Marshal(data []byte) (Id, error) {
	if len(data) > len(Id{}) || len(data) < len(Id{}) || data == nil {
		return Id{}, errors.New(fmt.Sprintf("Ephemeral ID must be of size %d", len(Id{})))
	}
	eid := Id{}
	copy(eid[:], data)
	return eid, nil
}

// GetIdsByRange returns ephemeral IDs based on passed in ID and a time range
// Accepts an ID, ID size in bits, timestamp in nanoseconds and a time range
// returns a list of ephemeral IDs
func GetIdsByRange(id *id.ID, size uint, timestamp time.Time,
	timeRange time.Duration) ([]ProtoIdentity, error) {

	if size > 64 {
		return []ProtoIdentity{}, errors.New("Cannot generate ID with size > 64")
	}

	iid, err := GetIntermediaryId(id)
	if err != nil {
		return []ProtoIdentity{}, err
	}

	var idList []ProtoIdentity
	timeStop := timestamp.Add(timeRange)

	for timeStop.After(timestamp) {

		newId, start, end, err := GetIdFromIntermediary(iid, size, timestamp.UnixNano())
		if err != nil {
			return []ProtoIdentity{}, err
		}

		idList = append(idList, ProtoIdentity{
			Id:    newId,
			Start: start,
			End:   end,
		})

		//make the timestamp into the next Period
		timestamp = end.Add(time.Nanosecond)
	}
	return idList, nil
}

// GetId returns ephemeral ID based on passed in ID
// Accepts an ID, ID size in bits, and timestamp in nanoseconds
// returns ephemeral ID, start & end timestamps for salt window
func GetId(id *id.ID, size uint, timestamp int64) (Id, time.Time, time.Time, error) {
	iid, err := GetIntermediaryId(id)
	if err != nil {
		return Id{}, time.Time{}, time.Time{}, err
	}
	return GetIdFromIntermediary(iid, size, timestamp)
}

// GetIntermediaryId returns an intermediary ID for ephemeral ID creation (ID hash)
func GetIntermediaryId(id *id.ID) ([]byte, error) {
	b2b := crypto.BLAKE2b_256.New()
	_, err := b2b.Write(id.Marshal())
	if err != nil {
		return nil, err
	}
	idHash := b2b.Sum(nil)
	return idHash, nil
}

// GetIdFromIntermediary returns the ephemeral ID from intermediary (id hash)
// Accepts an intermediary ephemeral ID, ID size in bits, and timestamp in nanoseconds
// returns ephemeral ID, start & end timestamps for salt window
func GetIdFromIntermediary(iid []byte, size uint, timestamp int64) (Id, time.Time, time.Time, error) {
	b2b := crypto.BLAKE2b_256.New()
	if size > 64 {
		return Id{}, time.Time{}, time.Time{}, errors.New("Cannot generate ID with size > 64")
	}
	salt, start, end := getRotationSalt(iid, timestamp)

	// Continually generate an ephemeral Id until we land on
	// an id not within the reserved list of Ids
	eid := Id{}
	var err error
	for reserved := true; reserved; reserved = IsReserved(eid) {
		eid, err = getIdFromIntermediaryHelper(b2b, iid, salt, size)
		if err != nil {
			return Id{}, start, end, err
		}
	}
	return eid, start, end, nil
}

// Helper function which generates a single ephemeral Id
func getIdFromIntermediaryHelper(b2b hash.Hash, iid, salt []byte, size uint) (Id, error) {
	eid := Id{}

	_, err := b2b.Write(iid)
	if err != nil {
		return Id{}, err
	}
	_, err = b2b.Write(salt)
	if err != nil {
		return Id{}, err
	}

	copy(eid[:], b2b.Sum(nil))

	cleared := eid.Clear(size)
	copy(eid[:], cleared[:])

	return eid, err
}

// Checks if the Id passed in is  among
// the reserved global reserved ID list.
// Returns true if reserved, false if non-reserved
func IsReserved(eid Id) bool {
	for _, r := range ReservedIDs {
		if bytes.Equal(eid[:], r[:]) {
			return true
		}
	}
	return false
}

// getRotationSalt returns rotation salt based on ID hash and timestamp
func getRotationSalt(idHash []byte, timestamp int64) ([]byte, time.Time, time.Time) {
	offset := GetOffset(idHash)
	start, end, saltNum := GetOffsetBounds(offset, timestamp)
	salt := make([]byte, 8)
	binary.BigEndian.PutUint64(salt, saltNum)
	return salt, start, end
}

func GetOffset(intermediaryId []byte) int64 {
	hashNum := binary.BigEndian.Uint64(intermediaryId)
	offset := int64((hashNum % uint64(NumOffsets)) * uint64(NsPerOffset))
	return offset
}

func GetOffsetNum(offset int64) int64 {
	return offset / NsPerOffset
}

func GetOffsetBounds(offset, timestamp int64) (time.Time, time.Time, uint64) {
	timestampPhase := timestamp % Period
	var start, end int64
	timestampNum := timestamp / Period
	var saltNum uint64
	if timestampPhase < offset {
		start = (timestampNum-1)*Period + offset
		end = start + Period
		saltNum = uint64((timestamp - Period) / Period)
	} else {
		start = timestampNum*Period + offset
		end = start + Period
		saltNum = uint64(timestamp / Period)
	}
	return time.Unix(0, start), time.Unix(0, end), saltNum
}

func HandleQuantization(start time.Time) (int64, int32) {
	currentOffset := (start.UnixNano() / NsPerOffset) % NumOffsets
	epoch := start.UnixNano() / NsPerOffset
	return currentOffset, int32(epoch)
}
