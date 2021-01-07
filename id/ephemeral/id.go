package ephemeral

import (
	"crypto"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/crypto/large"
	"gitlab.com/xx_network/primitives/id"
	"io"
	"math"
	"time"
)

var period = int64(time.Hour * 24)
var numOffsets int64 = 1 << 16

// Ephemeral ID type alias
type Id [8]byte

// Clear an ID down to the correct size
func (eid *Id) clear(size uint) {
	var mask uint64 = math.MaxUint64 >> (64 - size)
	maskedId := binary.BigEndian.Uint64(eid[:]) & mask
	binary.BigEndian.PutUint64(eid[:], maskedId)
}

// Fill cleared bits of an ID with random data from passed in rng
func (eid *Id) fill(size uint, rng io.Reader) error {
	rand := Id{}
	_, err := rng.Read(rand[:])
	if err != nil {
		return err
	}
	var mask uint64 = math.MaxUint64 << size
	maskedRand := mask & binary.BigEndian.Uint64(rand[:])
	binary.BigEndian.PutUint64(eid[:], maskedRand|binary.BigEndian.Uint64(eid[:]))
	return nil
}

// GetId returns ephemeral ID based on passed in ID
func GetId(id *id.ID, size uint, rng csprng.Source) (Id, error) {
	iid, err := GetIntermediaryId(id)
	if err != nil {
		return Id{}, err
	}
	return GetIdFromIntermediary(iid, size, rng)
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
func GetIdFromIntermediary(iid []byte, size uint, rng csprng.Source) (Id, error) {
	b2b := crypto.BLAKE2b_256.New()
	if size > 64 {
		return Id{}, errors.New("Cannot generate ID with size > 64")
	}
	salt := getRotationSalt(iid)

	_, err := b2b.Write(iid)
	if err != nil {
		return Id{}, err
	}
	_, err = b2b.Write(salt)
	if err != nil {
		return Id{}, err
	}
	eid := Id{}
	copy(eid[:], b2b.Sum(nil))
	fmt.Printf("unmodified: %+v\n", eid)

	eid.clear(size)
	fmt.Printf("cleared: %+v\n", eid)
	err = eid.fill(size, rng)
	if err != nil {
		return Id{}, err
	}
	fmt.Printf("filled: %+v\n", eid)

	return eid, nil
}

// getRotationSalt returns rotation salt based on ID hash and timestamp
func getRotationSalt(idHash []byte) []byte {
	hashNum := large.NewIntFromBytes(idHash)
	offset := large.NewInt(1).Mod(hashNum, large.NewInt(numOffsets)).Int64()
	ts := time.Now().UnixNano()
	timestampPhase := ts % period
	var saltNum uint64
	if timestampPhase < offset {
		saltNum = uint64((ts - period) / period)
	} else {
		saltNum = uint64(ts / period)
	}
	salt := make([]byte, 8)
	binary.BigEndian.PutUint64(salt, saltNum)
	return salt
}
