package ephemeral

import (
	"crypto"
	"encoding/binary"
	"github.com/pkg/errors"
	"gitlab.com/xx_network/crypto/large"
	"gitlab.com/xx_network/primitives/id"
	"math"
	"time"
)

var period = int64(time.Hour * 24)
var numOffsets int64 = 1 << 16

// Ephemeral ID type alias
type Id uint64

// GetId returns ephemeral ID based on passed in ID
func GetId(id *id.ID, size uint) (Id, error) {
	iid, err := GetIntermediaryId(id)
	if err != nil {
		return math.MaxUint64, err
	}
	return GetIdFromIntermediary(iid, size)
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
func GetIdFromIntermediary(iid []byte, size uint) (Id, error) {
	b2b := crypto.BLAKE2b_256.New()
	if size > 64 {
		return math.MaxUint64, errors.New("Cannot generate ID with size > 64")
	}
	salt := getRotationSalt(iid)

	_, err := b2b.Write(iid)
	if err != nil {
		return math.MaxUint64, err
	}
	_, err = b2b.Write(salt)
	if err != nil {
		return math.MaxUint64, err
	}
	eid := b2b.Sum(nil)
	return Id(large.NewIntFromBytes(eid).Uint64()), nil
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
