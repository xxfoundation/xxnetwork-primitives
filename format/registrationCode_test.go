////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/crypto/hash"
	"math"
	"reflect"
	"testing"
)

const REG_CODE string = "AB9DBDC17945CA1DECC31D8B8366337967E682252C71B956B7BAD49ABA69BF2A"

func TestDisassembleRegistrationCode(t *testing.T) {

	regcode := cyclic.NewIntFromString(REG_CODE, 16).LeftpadBytes(REGCODE_LEN)

	regpinExpected := uint32(cyclic.NewIntFromBytes(
		regcode[REGPIN_START:REGPIN_END]).Uint64())

	regkeyExpected := regcode[REGKEY_START:REGKEY_END]

	regkey, regpin := DisassembleRegistrationCode(regcode)

	if regpin != regpinExpected {
		t.Errorf("Test of DisassembleRegistrationCode failed: Regestration"+
			" Pin codes did not match (%v) differed from expected (%v)", regpin,
			regpinExpected)
	}

	if !reflect.DeepEqual(regkey, regkeyExpected) {
		t.Errorf("Test of DisassembleRegistrationCode failed: Regestration"+
			" key did not match (%v) differed from expected (%v)",
			cyclic.NewIntFromBytes(regkey).Text(32),
			cyclic.NewIntFromBytes(regkeyExpected).Text(32))
	}
}

func TestRegistrationHash(t *testing.T) {

	regcode := cyclic.NewIntFromString(REG_CODE, 16).LeftpadBytes(REGCODE_LEN)

	hasher, _ := hash.NewCMixHash()

	hasher.Write(regcode)

	expectedHash := hasher.Sum(nil)

	regkey, regpin := DisassembleRegistrationCode(regcode)

	h, _ := RegistrationHash(regkey, regpin)

	if !reflect.DeepEqual(h, expectedHash) {
		t.Errorf("Test of RegistrationHash did not match: Regestration hash"+
			" did not match (%v) differed from expected (%v)",
			cyclic.NewIntFromBytes(h).Text(32),
			cyclic.NewIntFromBytes(expectedHash).Text(32))
	}

	_, err := RegistrationHash(regkey, math.MaxUint32)

	if err == nil {
		t.Errorf("Test of RegistrationHash did not match: Out of Range pin" +
			" accepted")
	}
}
