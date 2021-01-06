package ephemeral

import (
	"crypto"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/crypto/large"
	"gitlab.com/xx_network/primitives/id"
	"math"
	"time"
)

var period = int64(time.Hour * 24)
var numOffsets int64 = 1 << 16

// Ephemeral ID type alias
type Id [8]byte

// GetId returns ephemeral ID based on passed in ID
func GetId(id *id.ID, size uint) (Id, error) {
	iid, err := GetIntermediaryId(id)
	if err != nil {
		return Id{}, err
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
	eid := b2b.Sum(nil)
	var mask float64 = 0
	for i := uint(0); i < size; i++ {
		mask += math.Pow(2, float64(i))
	}
	bitMask := large.NewInt(0).LeftShift(large.NewInt(int64(mask)), 64-size)
	eidBits := large.NewInt(0).And(bitMask, large.NewIntFromBytes(eid))
	fmt.Printf("EID BITS: %+v\n", eidBits.Bytes())

	rand := Id{}
	rng := csprng.NewSystemRNG()
	_, err = rng.Read(rand[:])
	if err != nil {
		return Id{}, err
	}
	var inverseMask float64 = 0
	for i := uint(0); i < 64-size; i++ {
		inverseMask += math.Pow(2, float64(i))
	}
	inverseBitMask := large.NewInt(int64(inverseMask))
	randBits := large.NewInt(0).And(inverseBitMask, large.NewIntFromBytes(rand[:]))
	fmt.Printf("RAND BITS: %+v\n", randBits.Bytes())

	finalEid := large.NewInt(0).Or(randBits, eidBits)
	fmt.Printf("Final bits: %+v\n", finalEid.Bytes())

	ret := Id{}
	copy(ret[:], finalEid.Bytes())
	return ret, nil
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
