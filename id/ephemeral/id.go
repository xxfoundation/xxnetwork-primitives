////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package ephemeral

import (
	"crypto"
	"crypto/hmac"
	"encoding/binary"
	"encoding/json"
	"hash"
	"io"
	"math"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/xx_network/primitives/id"
)

const (
	// IdLen is the length of an ephemeral [Id].
	IdLen = 8

	Period            = int64(time.Hour * 24)
	NumOffsets  int64 = 1 << 16
	NsPerOffset       = Period / NumOffsets

	// Minimum and maximum size for new Id.
	maxSize = 64
	minSize = 1
)

// ReservedIDs are ephemeral IDs reserved for specific actions:
//   - All zeros denote a dummy ID
//   - All ones denote a payment
var ReservedIDs = []Id{
	{0, 0, 0, 0, 0, 0, 0, 0}, // Dummy ID
	{1, 1, 1, 1, 1, 1, 1, 1}, // Payment
}

// Id is the ephemeral ID type.
type Id [IdLen]byte

// ProtoIdentity contains the ID and the start and end time for the salt window.
type ProtoIdentity struct {
	Id    Id
	Start time.Time
	End   time.Time
}

// UInt64 returns the ephemeral ID as a uint64.
func (eid *Id) UInt64() uint64 {
	return binary.BigEndian.Uint64(eid[:])
}

// Int64 returns the ephemeral ID as an int64.
func (eid *Id) Int64() int64 {
	ux := binary.BigEndian.Uint64(eid[:])
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x
}

// Clear clears all the bits in the ID outside of the given size.
func (eid Id) Clear(size uint) Id {
	var newId Id
	var mask uint64 = math.MaxUint64 >> (64 - size)
	maskedId := binary.BigEndian.Uint64(eid[:]) & mask
	binary.BigEndian.PutUint64(newId[:], maskedId)
	return newId
}

// Fill sets the bits of an ID with random data from passed in rng. The size of
// the ID is in bits.
func (eid Id) Fill(size uint, rng io.Reader) (Id, error) {
	var newId, rand Id
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

// Marshal loads an ephemeral ID from raw bytes.
func Marshal(data []byte) (Id, error) {
	if data == nil || len(data) != IdLen {
		return Id{}, errors.Errorf("Ephemeral ID must be of size %d", IdLen)
	}
	var eid Id
	copy(eid[:], data)
	return eid, nil
}

// GetIdsByRange returns ephemeral IDs based on passed in ID and a time range.
// Accepts an ID, ID size in bits, timestamp, and a time range. Returns a list
// of ephemeral IDs.
func GetIdsByRange(id *id.ID, size uint, timestamp time.Time,
	timeRange time.Duration) ([]ProtoIdentity, error) {

	if size > maxSize {
		return []ProtoIdentity{},
			errors.Errorf("Cannot generate ID with size > %d", maxSize)
	}

	iid, err := GetIntermediaryId(id)
	if err != nil {
		return []ProtoIdentity{}, err
	}

	var idList []ProtoIdentity
	timeStop := timestamp.Add(timeRange)

	for timeStop.After(timestamp) {

		newId, start, end, err :=
			GetIdFromIntermediary(iid, size, timestamp.UnixNano())
		if err != nil {
			return []ProtoIdentity{}, err
		}

		idList = append(idList, ProtoIdentity{
			Id:    newId,
			Start: start,
			End:   end,
		})

		// Make the timestamp into the next Period
		timestamp = end.Add(time.Nanosecond)
	}
	return idList, nil
}

// GetId returns ephemeral ID based on the passed in ID.
// Accepts an ID, ID size in bits, and timestamp in nanoseconds.
// Returns an ephemeral ID and the start and end timestamps for the salt window.
func GetId(id *id.ID, size uint, timestamp int64) (Id, time.Time, time.Time, error) {
	iid, err := GetIntermediaryId(id)
	if err != nil {
		return Id{}, time.Time{}, time.Time{}, err
	}
	return GetIdFromIntermediary(iid, size, timestamp)
}

// GetIntermediaryId returns an intermediary ID for the ephemeral ID creation
// (ID hash).
func GetIntermediaryId(id *id.ID) ([]byte, error) {
	b2b := crypto.BLAKE2b_256.New()
	_, err := b2b.Write(id.Marshal())
	if err != nil {
		return nil, err
	}
	idHash := b2b.Sum(nil)
	return idHash, nil
}

// GetIdFromIntermediary returns the ephemeral ID from intermediary (ID hash).
// Accepts an intermediary ephemeral ID, ID size in bits, and timestamp in
// nanoseconds.
// Returns an ephemeral ID and the start and end timestamps for salt window.
func GetIdFromIntermediary(iid []byte, size uint, timestamp int64) (
	Id, time.Time, time.Time, error) {
	b2b := crypto.BLAKE2b_256.New()
	if size > maxSize || size < minSize {
		return Id{}, time.Time{}, time.Time{}, errors.Errorf("Cannot generate "+
			"ID, size must be between %d and %d", minSize, maxSize)
	}
	salt, start, end := getRotationSalt(iid, timestamp)

	// Continually generate an ephemeral ID until we land on an ID not within
	// the reserved list of IDs
	var eid Id
	var err error
	for reserved := true; reserved; reserved = IsReserved(eid) {
		eid, err = getIdFromIntermediary(b2b, iid, salt, size)
		if err != nil {
			return Id{}, start, end, err
		}
	}
	return eid, start, end, nil
}

// getIdFromIntermediary generates an ephemeral Id from an intermediary ID and
// salt using the provided hash.
func getIdFromIntermediary(
	b2b hash.Hash, iid, salt []byte, size uint) (Id, error) {
	var eid Id

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

// IsReserved checks if the Id is among the reserved global reserved ID list.
// Returns true if reserved, false if non-reserved.
func IsReserved(eid Id) bool {
	for _, r := range ReservedIDs {
		if hmac.Equal(eid[:], r[:]) {
			return true
		}
	}
	return false
}

// getRotationSalt returns rotation salt based on ID hash and timestamp.
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

// MarshalJSON marshals the ephemeral [Id] into valid JSON. This function
// adheres to the [json.Marshaler] interface.
func (eid Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(eid[:])
}

// UnmarshalJSON unmarshalls the JSON into the ephemeral [Id]. This function
// adheres to the [json.Unmarshaler] interface.
func (eid *Id) UnmarshalJSON(data []byte) error {
	var idBytes []byte
	err := json.Unmarshal(data, &idBytes)
	if err != nil {
		return err
	}

	newEid, err := Marshal(idBytes)
	if err != nil {
		return err
	}

	*eid = newEid

	return nil
}
