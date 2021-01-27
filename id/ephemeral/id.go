package ephemeral

import (
	"crypto"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/xx_network/primitives/id"
	"io"
	"math"
	"time"
)

var period = uint64(time.Hour * 24)
var numOffsets uint64 = 1 << 16
var nsPerOffset = period / numOffsets

// Ephemeral ID type alias
type Id [8]byte

// Ephemeral Id object which contains the ID
// and the start and end time for the salt window
type EphemeralId struct {
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
func Marshal(data []byte) (*Id, error) {
	if len(data) > len(Id{}) || len(data) < len(Id{}) || data == nil {
		return nil, errors.New(fmt.Sprintf("Ephemeral ID must be of size %d", len(Id{})))
	}
	eid := &Id{}
	copy(eid[:], data)
	return eid, nil
}

// GetIdByRange returns ephemeral IDs based on passed in ID and a time range
// Accepts an ID, ID size in bits, timestamp in nanoseconds and a time range
// returns a list of ephemeral IDs
func GetIdByRange(id *id.ID, size uint, timestamp int64,
	timeRange time.Duration) ([]EphemeralId, error) {

	if size > 64 {
		return []EphemeralId{}, errors.New("Cannot generate ID with size > 64")
	}

	iid, err := GetIntermediaryId(id)
	if err != nil {
		return []EphemeralId{}, err
	}

	idList := make([]EphemeralId, 0)

	idsToGenerate := int64(timeRange) / int64(period)

	for i := int64(0); i < idsToGenerate; i++ {
		nextTimestamp := timestamp + i*int64(period)
		newId, start, end, err := GetIdFromIntermediary(iid, size, nextTimestamp)
		if err != nil {
			return []EphemeralId{}, err
		}

		ephId := EphemeralId{
			Id:    newId,
			Start: start,
			End:   end,
		}

		idList = append(idList, ephId)
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
	salt, start, end := getRotationSalt(iid, uint64(timestamp))

	_, err := b2b.Write(iid)
	if err != nil {
		return Id{}, start, end, err
	}
	_, err = b2b.Write(salt)
	if err != nil {
		return Id{}, start, end, err
	}
	eid := Id{}
	copy(eid[:], b2b.Sum(nil))

	cleared := eid.Clear(size)
	copy(eid[:], cleared[:])

	return eid, start, end, nil
}

// getRotationSalt returns rotation salt based on ID hash and timestamp
func getRotationSalt(idHash []byte, timestamp uint64) ([]byte, time.Time, time.Time) {
	hashNum := binary.BigEndian.Uint64(idHash)
	offset := (hashNum % numOffsets) * nsPerOffset
	timestampPhase := timestamp % period
	var start, end uint64
	timestampNum := timestamp / period
	var saltNum uint64
	if timestampPhase < offset {
		start = (timestampNum-1)*period + offset
		end = start + period
		saltNum = (timestamp - period) / period
	} else {
		start = timestampNum*period + offset
		end = start + period
		saltNum = timestamp / period
	}
	salt := make([]byte, 8)
	binary.BigEndian.PutUint64(salt, saltNum)
	return salt, time.Unix(0, int64(start)), time.Unix(0, int64(end))
}
