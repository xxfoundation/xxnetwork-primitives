////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"errors"
	"fmt"
	// TODO Should we also delete this to completely remove the crypto dependency?
	"gitlab.com/elixxir/crypto/hash"
)

const REGCODE_LEN uint64 = 32

//Registration pin in 24 bits, in hex every character is 4 bits,
// so 24/4 = 6 characters
const REGPIN_LEN uint64 = 3
const REGPIN_START uint64 = 0
const REGPIN_END uint64 = REGPIN_START + REGPIN_LEN

// Max value for 6 digit registration key, (2^24)-1
const REGPIN_MAX uint32 = uint32((1 << (REGPIN_LEN * 8)) - 1)

const REGKEY_LEN uint64 = REGCODE_LEN - REGPIN_LEN
const REGKEY_START uint64 = REGPIN_END
const REGKEY_END uint64 = REGKEY_START + REGKEY_LEN

// Takes a Registration Code and returns the Registration Key and
// Registration Pin
func DisassembleRegistrationCode(regcode []byte) ([]byte, uint32) {
	// convert byteslice to uint32 manually because binary.BigEndian expects
	// a slice with a particular length
	var pin uint32
	for i := uint64(0); i < REGPIN_LEN; i++ {
		pin <<= 8
		pin |= uint32(regcode[REGPIN_START+i])
	}
	return regcode[REGKEY_START:REGKEY_END], pin
}

// Takes a Registration Key and Registration Pin, combines them,
// and returns the Registration Hash
func RegistrationHash(regkey []byte, regpin uint32) ([]byte, error) {

	//Make sure the pin is in range
	if regpin > REGPIN_MAX {
		return nil, errors.New(fmt.Sprintf(
			"Could not make Registration Hash: Regestrtation Pin too long"+
				"; Max: %v, Recieved: %v", REGPIN_MAX, regpin))
	}

	//Rebuild the registration code
	regcode := make([]byte, REGCODE_LEN)

	//Turn the pin into a byte slice and copy it into the registration code
	regpinBytes := make([]byte, REGPIN_LEN)
	for i := int(REGPIN_LEN - 1); i >= 0; i-- {
		regpinBytes[i] = byte(regpin)
		regpin >>= 8
	}
	copy(regcode[REGPIN_START:REGPIN_END], regpinBytes)

	copy(regcode[REGKEY_START:REGKEY_END], regkey)

	//Get the object to hash the code with
	hasher, err := hash.NewCMixHash()

	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Could not make Regestration Hash: Could not get Regestartion"+
				" Hash: %v", err.Error()))
	}

	//Hash the code
	hasher.Write(regcode)

	b := hasher.Sum(nil)

	//Return the correct region
	return b[0:REGCODE_LEN], nil
}
