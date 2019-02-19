////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"encoding/binary"
	"encoding/hex"
	"gitlab.com/elixxir/crypto/hash"
	"math"
	"reflect"
	"testing"
)

const REG_CODE string = "AB9DBDC17945CA1DECC31D8B8366337967E682252C71B956B7BAD49ABA69BF2A"

func TestDisassembleRegistrationCode(t *testing.T) {

	regcode, err := hex.DecodeString(REG_CODE)
	if err != nil {
		t.Error(err.Error())
	}

	// Make a right-sized slice for bigendian encoding
	regcode32Bits := make([]byte, 4)
	copy(regcode32Bits[1:], regcode[REGPIN_START:REGPIN_END])
	regpinExpected := binary.BigEndian.Uint32(regcode32Bits)

	regkeyExpected := regcode[REGKEY_START:REGKEY_END]

	regkey, regpin := DisassembleRegistrationCode(regcode)

	if regpin != regpinExpected {
		t.Errorf("Test of DisassembleRegistrationCode failed: Regestration"+
			" Pin codes did not match (%x) differed from expected (%x)", regpin,
			regpinExpected)
	}

	if !reflect.DeepEqual(regkey, regkeyExpected) {
		t.Errorf("Test of DisassembleRegistrationCode failed: Regestration"+
			" key did not match (%v) differed from expected (%v)",
			hex.EncodeToString(regkey),
			hex.EncodeToString(regkeyExpected))
	}
}

func TestRegistrationHash(t *testing.T) {

	regcode := make([]byte, REGCODE_LEN)
	regBytes, err := hex.DecodeString(REG_CODE)
	if err != nil {
		t.Error(err.Error())
	}
	copy(regcode[len(regcode)-len(regBytes):], regBytes)

	hasher, _ := hash.NewCMixHash()

	hasher.Write(regcode)

	expectedHash := hasher.Sum(nil)

	regkey, regpin := DisassembleRegistrationCode(regcode)

	h, _ := RegistrationHash(regkey, regpin)

	if !reflect.DeepEqual(h, expectedHash) {
		t.Errorf("Test of RegistrationHash did not match: Regestration hash"+
			" did not match (%v) differed from expected (%v)",
			hex.EncodeToString(h),
			hex.EncodeToString(expectedHash))
	}

	_, err = RegistrationHash(regkey, math.MaxUint32)

	if err == nil {
		t.Errorf("Test of RegistrationHash did not match: Out of Range pin" +
			" accepted")
	}
}
