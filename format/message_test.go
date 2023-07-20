////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestMessage_VersionDetection(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	// Generate message parts
	fp := NewFingerprint(makeAndFillSlice(KeyFPLen, 'c'))
	mac := makeAndFillSlice(MacLen, 'd')
	ephemeralRID := makeAndFillSlice(EphemeralRIDLen, 'e')
	identityFP := makeAndFillSlice(SIHLen, 'f')
	contents := makeAndFillSlice(MinimumPrimeSize*2-AssociatedDataSize-1, 'g')

	// Set message parts
	msg.SetKeyFP(fp)
	msg.SetMac(mac)
	msg.SetEphemeralRID(ephemeralRID)
	msg.SetSIH(identityFP)
	msg.SetContents(contents)

	copy(msg.version, []byte{123})

	msgBytes := msg.Marshal()
	_, err := Unmarshal(msgBytes)
	if err == nil {
		// t.Error("version detection fail")
		t.Logf("Version detection disabled")
	}
}

func TestMessage_Smoke(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	// Generate message parts
	fp := NewFingerprint(makeAndFillSlice(KeyFPLen, 'c'))
	mac := makeAndFillSlice(MacLen, 'd')
	ephemeralRID := makeAndFillSlice(EphemeralRIDLen, 'e')
	identityFP := makeAndFillSlice(SIHLen, 'f')
	contents := makeAndFillSlice(MinimumPrimeSize*2-AssociatedDataSize-1, 'g')

	// Set message parts
	msg.SetKeyFP(fp)
	msg.SetMac(mac)
	msg.SetEphemeralRID(ephemeralRID)
	msg.SetSIH(identityFP)
	msg.SetContents(contents)

	if !bytes.Equal(fp.Bytes(), msg.keyFP) {
		t.Errorf("keyFp data was corrupted.\nexpected: %+v\nreceived: %+v",
			fp.Bytes(), msg.keyFP)
	}

	if !bytes.Equal(mac, msg.mac) {
		t.Errorf("MAC data was corrupted.\nexpected: %+v\nreceived: %+v",
			mac, msg.mac)
	}

	if !bytes.Equal(ephemeralRID, msg.ephemeralRID) {
		t.Errorf("ephemeralRID data was corrupted.\nexpected: %+v\nreceived: %+v",
			ephemeralRID, msg.ephemeralRID)
	}

	if !bytes.Equal(identityFP, msg.sih) {
		t.Errorf("sih data was corrupted.\nexpected: %+v\nreceived: %+v",
			identityFP, msg.sih)
	}

	if !bytes.Equal(contents, append(msg.contents1, msg.contents2...)) {
		t.Errorf("contents data was corrupted.\nexpected: %+v\nreceived: %+v",
			contents, append(msg.contents1, msg.contents2...))
	}
}

// Happy path.
func TestNewMessage(t *testing.T) {
	numPrimeBytes := MinimumPrimeSize
	expectedMsg := Message{
		data:         make([]byte, 2*numPrimeBytes),
		payloadA:     make([]byte, numPrimeBytes),
		payloadB:     make([]byte, numPrimeBytes),
		keyFP:        make([]byte, KeyFPLen),
		version:      make([]byte, 1),
		contents1:    make([]byte, numPrimeBytes-KeyFPLen-1),
		mac:          make([]byte, MacLen),
		contents2:    make([]byte, numPrimeBytes-MacLen-RecipientIDLen),
		ephemeralRID: make([]byte, EphemeralRIDLen),
		sih:          make([]byte, SIHLen),
		rawContents:  make([]byte, 2*numPrimeBytes-RecipientIDLen),
	}

	msg := NewMessage(MinimumPrimeSize)

	if !reflect.DeepEqual(expectedMsg, msg) {
		t.Errorf("NewMessage did not return the expected Message."+
			"\nexpected: %+v\nreceived: %+v", expectedMsg, msg)
	}
}

// Error path: panics if provided prime size is too small.
func TestNewMessage_NumPrimeBytesError(t *testing.T) {
	// Defer to an error when NewMessage does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewMessage did not panic when the minimum prime size " +
				"is too small.")
		}
	}()

	_ = NewMessage(MinimumPrimeSize - 1)
}

// Happy path.
func TestMessage_Marshal_Unmarshal(t *testing.T) {
	m := NewMessage(256)
	prng := rand.New(rand.NewSource(time.Now().UnixNano()))
	payload := make([]byte, 256)
	prng.Read(payload)
	m.SetPayloadA(payload)
	prng.Read(payload)
	m.SetPayloadB(payload)
	copy(m.version, []byte{messagePayloadVersion})

	messageData := m.Marshal()
	newMsg, err := Unmarshal(messageData)

	if err != nil {
		t.Errorf("Unmarshal failure: %#v", err)
	}

	if !reflect.DeepEqual(m, newMsg) {
		t.Errorf("Failed to Marshal and Unmarshal message."+
			"\nexpected: %#v\nreceived: %#v", m, newMsg)
	}
}

func TestMessage_Marshal_UnmarshalImmutable(t *testing.T) {
	m := NewMessage(256)
	prng := rand.New(rand.NewSource(time.Now().UnixNano()))
	payload := make([]byte, 256)
	prng.Read(payload)
	m.SetPayloadA(payload)
	prng.Read(payload)
	m.SetPayloadB(payload)
	copy(m.version, []byte{messagePayloadVersion})

	m.ephemeralRID[0] = 42
	m.sih[0] = 42

	messageData := m.MarshalImmutable()
	newMsg, err := Unmarshal(messageData)

	if err != nil {
		t.Errorf("Unmarshal failure: %#v", err)
	}

	if newMsg.ephemeralRID[0] != 0 {
		t.Errorf("MarshalImmutable did not clear EphemeralRID: %v != 0",
			newMsg.ephemeralRID)
	}
	if newMsg.sih[0] != 0 {
		t.Errorf("MarshalImmutable did not clear SIH: %v != 0",
			newMsg.sih)
	}
}

// Happy path.
func TestMessage_Version(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	if msg.Version() != 0 {
		t.Errorf("Unexpected version for new Message."+
			"\nexpected: %d\nreceived: %d", 0, msg.Version())
	}

	copy(msg.version, []byte{123})

	if msg.Version() != 123 {
		t.Errorf("Unexpected version."+
			"\nexpected: %d\nreceived: %d", 123, msg.Version())
	}
}

// Happy path.
func TestMessage_Copy(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	msgCopy := msg.Copy()

	contents := make([]byte, MinimumPrimeSize*2-AssociatedDataSize-1)
	copy(contents, "test")
	msgCopy.SetContents(contents)

	if bytes.Equal(msg.GetContents(), contents) {
		t.Errorf("Copy failed to make a copy of the message; modifications " +
			"to copy reflected in original.")
	}
}

// Happy path.
func TestMessage_GetPrimeByteLen(t *testing.T) {
	primeSize := 250
	m := NewMessage(primeSize)

	if m.GetPrimeByteLen() != primeSize {
		t.Errorf("GetPrimeByteLen returned incorrect prime size."+
			"\nexpected: %d\nreceived: %d", primeSize, m.GetPrimeByteLen())
	}
}

// Happy path.
func TestMessage_GetPayloadA(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	testData := []byte("test")
	copy(msg.payloadA, testData)
	payload := msg.GetPayloadA()
	if !bytes.Equal(testData, payload[:len(testData)]) {
		t.Errorf("GetPayloadA did not properly retrieve payload A."+
			"\nexpected: %s\nreceived: %s", testData, payload[:len(testData)])
	}

	// Ensure that the data is copied
	payload[14] = 'x'
	if msg.payloadA[14] == payload[14] {
		t.Error("GetPayloadA did not make a copy; modifications to copy " +
			"reflected in original.")
	}
}

// Happy path.
func TestMessage_SetPayloadA(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	payload := make([]byte, len(msg.payloadA))
	copy(payload, "test")
	msg.SetPayloadA(payload)

	if !bytes.Equal(payload, msg.payloadA) {
		t.Errorf("SetPayloadA failed to set payload A correctly."+
			"\nexpected: %s\nreceived: %s", payload, msg.payloadA)
	}
}

// Error path: length of provided payload incorrect.
func TestMessage_SetPayloadA_LengthError(t *testing.T) {
	payload := make([]byte, MinimumPrimeSize/4)
	msg := NewMessage(MinimumPrimeSize)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetPayloadA failed to panic when the length of the "+
				"provided payload (%d) is not the same as the message payload "+
				"length (%d).", len(payload), len(msg.GetPayloadA()))
		}
	}()

	msg.SetPayloadA(payload)
}

// Happy path.
func TestMessage_GetPayloadB(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	testData := []byte("test")
	copy(msg.payloadB, testData)
	payload := msg.GetPayloadB()
	if !bytes.Equal(testData, payload[:len(testData)]) {
		t.Errorf("GetPayloadB did not properly retrieve payload B."+
			"\nexpected: %s\nreceived: %s", testData, payload[:len(testData)])
	}

	// Ensure that the data is copied
	payload[14] = 'x'
	if msg.payloadB[14] == payload[14] {
		t.Error("GetPayloadB did not make a copy; modifications to copy " +
			"reflected in original.")
	}
}

// Happy path.
func TestMessage_SetPayloadB(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	payload := make([]byte, len(msg.payloadB))
	copy(payload, "test")
	msg.SetPayloadB(payload)

	if !bytes.Equal(payload, msg.payloadB) {
		t.Errorf("SetPayloadB failed to set payload B correctly."+
			"\nexpected: %s\nreceived: %s", payload, msg.payloadB)
	}
}

// Error path: length of provided payload incorrect.
func TestMessage_SetPayloadB_LengthError(t *testing.T) {
	payload := make([]byte, MinimumPrimeSize/4)
	msg := NewMessage(MinimumPrimeSize)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetPayloadB failed to panic when the length of the "+
				"provided payload (%d) is not the same as the message payload "+
				"length (%d).", len(payload), len(msg.GetPayloadB()))
		}
	}()

	msg.SetPayloadB(payload)
}

// Happy path.
func TestMessage_ContentsSize(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	if msg.ContentsSize() != MinimumPrimeSize*2-AssociatedDataSize-1 {
		t.Errorf("ContentsSize returned the wrong content size."+
			"\nexpected: %d\nreceived: %d",
			MinimumPrimeSize*2-AssociatedDataSize-1, msg.ContentsSize())
	}
}

// Happy path.
func TestMessage_GetContents(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	contents := makeAndFillSlice(MinimumPrimeSize*2-AssociatedDataSize-1, 'a')

	copy(msg.contents1, contents[:len(msg.contents1)])
	copy(msg.contents2, contents[len(msg.contents1):])

	retrieved := msg.GetContents()

	if !bytes.Equal(retrieved, contents) {
		t.Errorf("GetContents did not return the expected contents."+
			"\nexpected: %s\nreceived: %s", contents, retrieved)
	}
}

// Happy path: set contents that is large enough to fit in both contents.
func TestMessage_SetContents(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	contents := makeAndFillSlice(MinimumPrimeSize*2-AssociatedDataSize-1, 'a')

	msg.SetContents(contents)

	if !bytes.Equal(msg.contents1, contents[:len(msg.contents1)]) {
		t.Errorf("SetContents did not set contents1 correctly."+
			"\nexpected: %s\nreceived: %s",
			contents[:len(msg.contents1)], msg.contents1)
	}

	if !bytes.Equal(msg.contents2, contents[len(msg.contents1):]) {
		t.Errorf("SetContents did not set contents2 correctly."+
			"\nexpected: %s\nreceived: %s",
			contents[len(msg.contents1):], msg.contents2)
	}
}

// Happy path: set contents that is small enough to fit in the first contents.
func TestMessage_SetContents_ShortContents(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	contents := makeAndFillSlice(MinimumPrimeSize-KeyFPLen-1, 'a')

	msg.SetContents(contents)

	if !bytes.Equal(msg.contents1, contents[:len(msg.contents1)]) {
		t.Errorf("SetContents did not set contents1 correctly."+
			"\nexpected: %s\nreceived: %s",
			contents[:len(msg.contents1)], msg.contents1)
	}

	expectedContents2 :=
		make([]byte, MinimumPrimeSize-MacLen-EphemeralRIDLen-SIHLen)
	if !bytes.Equal(msg.contents2, expectedContents2) {
		t.Errorf("SetContents did not set contents2 correctly."+
			"\nexpected: %+v\nreceived: %+v", expectedContents2, msg.contents2)
	}
}

// Error path: content size too large.
func TestMessage_SetContents_ContentsTooLargeError(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	contents := makeAndFillSlice(MinimumPrimeSize*2, 'a')

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetContents failed to panic when the length of the "+
				"provided contents (%d) is larger than the max content length "+
				"(%d).", len(contents), len(msg.contents1)+len(msg.contents2))
		}
	}()

	msg.SetContents(contents)
}

// Happy path.
func TestMessage_GetRawContentsSize(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	expectedLen := (2 * MinimumPrimeSize) - RecipientIDLen

	if msg.GetRawContentsSize() != expectedLen {
		t.Errorf("GetRawContentsSize did not return the expected size."+
			"\nexpected: %d\nreceived: %d", expectedLen, msg.GetRawContentsSize())
	}
}

// Happy path.
func TestMessage_GetRawContents(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	// Created expected data
	var expectedRawContents []byte
	keyFP := makeAndFillSlice(KeyFPLen, 'a')
	mac := makeAndFillSlice(MacLen, 'b')
	contents1 := makeAndFillSlice(MinimumPrimeSize-KeyFPLen-1, 'c')
	contents2 := makeAndFillSlice(MinimumPrimeSize-MacLen-RecipientIDLen, 'd')
	expectedRawContents = append(expectedRawContents, keyFP...)
	expectedRawContents = append(expectedRawContents, byte(messagePayloadVersion))
	expectedRawContents = append(expectedRawContents, contents1...)
	expectedRawContents = append(expectedRawContents, mac...)
	expectedRawContents = append(expectedRawContents, contents2...)

	// Copy contents into message
	copy(msg.keyFP, keyFP)
	copy(msg.mac, mac)
	copy(msg.version, []byte{messagePayloadVersion})
	copy(msg.contents1, contents1)
	copy(msg.contents2, contents2)

	// Make sure the 1st and middle+1 bits are 1
	msg.payloadA[0] |= 0b10000000
	msg.payloadB[0] |= 0b10000000

	rawContents := msg.GetRawContents()
	if !bytes.Equal(expectedRawContents, rawContents) {
		t.Errorf("GetRawContents did not return the expected raw contents."+
			"\nexpected: %s\nreceived: %s", expectedRawContents, rawContents)
	}

	if rawContents[0]&0b10000000 != 0 {
		t.Errorf("First bit not set to zero")
	}

	fmt.Println(rawContents[msg.GetPrimeByteLen()])

	if rawContents[msg.GetPrimeByteLen()]&0b10000000 != 0 {
		t.Errorf("middle plus one bit not set to zero")
	}

}

// Happy path.
func TestMessage_SetRawContents(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	spLen := (2 * MinimumPrimeSize) - RecipientIDLen
	sp := make([]byte, spLen)

	fp := makeAndFillSlice(len(msg.keyFP), 'f')
	mac := makeAndFillSlice(len(msg.mac), 'm')
	c1 := makeAndFillSlice(len(msg.contents1), 'a')
	c2 := makeAndFillSlice(len(msg.contents2), 'b')

	copy(sp[:KeyFPLen], fp)
	copy(sp[MinimumPrimeSize:MinimumPrimeSize+MacLen], mac)

	copy(sp[KeyFPLen:MinimumPrimeSize], c1)
	copy(sp[MinimumPrimeSize+MacLen:2*MinimumPrimeSize-RecipientIDLen], c2)

	msg.SetRawContents(sp)

	if bytes.Contains(msg.keyFP, []byte("a")) ||
		bytes.Contains(msg.keyFP, []byte("b")) ||
		bytes.Contains(msg.keyFP, []byte("m")) ||
		!bytes.Contains(msg.keyFP, []byte("f")) {
		t.Errorf("Setting raw payload failed, key fingerprint contains "+
			"wrong data: %s", msg.keyFP)
	}

	if bytes.Contains(msg.mac, []byte("a")) ||
		bytes.Contains(msg.mac, []byte("b")) ||
		!bytes.Contains(msg.mac, []byte("m")) ||
		bytes.Contains(msg.mac, []byte("f")) {
		t.Errorf(
			"Setting raw payload failed, mac contains wrong data: %s", msg.mac)
	}

	if !bytes.Contains(msg.contents1, []byte("a")) ||
		bytes.Contains(msg.contents1, []byte("b")) ||
		bytes.Contains(msg.contents1, []byte("m")) ||
		bytes.Contains(msg.contents1, []byte("f")) {
		t.Errorf("Setting raw payload failed, contents1 contains wrong data: %s",
			msg.contents1)
	}

	if bytes.Contains(msg.contents2, []byte("a")) ||
		!bytes.Contains(msg.contents2, []byte("b")) ||
		bytes.Contains(msg.contents2, []byte("m")) ||
		bytes.Contains(msg.contents2, []byte("f")) {
		t.Errorf("Setting raw payload failed, contents2 contains wrong data: %s",
			msg.contents2)
	}

}

// Error path: length of provided raw contents incorrect.
func TestMessage_SetRawContents_LengthError(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	defer func() {
		if r := recover(); r == nil {
			t.Error("SetRawContents failed to panic when length of the " +
				"provided data is incorrect.")
		}
	}()

	msg.SetRawContents(make([]byte, MinimumPrimeSize))
}

// Happy path.
func TestMessage_GetKeyFP(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	keyFP := NewFingerprint(makeAndFillSlice(SIHLen, 'e'))
	msg.keyFP[0] |= 0b10000000
	copy(msg.keyFP, keyFP.Bytes())

	if keyFP != msg.GetKeyFP() {
		t.Errorf("GetKeyFP failed to get the correct keyFP."+
			"\nexpected: %+v\nreceived: %+v", keyFP, msg.GetKeyFP())
	}

	// Ensure that the data is copied
	keyFP[2] = 'x'
	if msg.sih[2] == 'x' {
		t.Error("GetKeyFP failed to make a copy of keyFP.")
	}

	if msg.GetKeyFP()[0]&0b10000000 != 0 {
		t.Errorf("First bit not set to zero")
	}
}

// Happy path.
func TestMessage_SetKeyFP(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	fp := NewFingerprint(makeAndFillSlice(SIHLen, 'e'))

	msg.SetKeyFP(fp)

	if !bytes.Equal(fp.Bytes(), msg.keyFP) {
		t.Errorf("SetKeyFP failed to set keyFP."+
			"\nexpected: %+v\nreceived: %+v", fp, msg.keyFP)
	}
}

// Error path: first bit of provided data is not 0.
func TestMessage_SetKeyFP_FirstBitError(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	fp := NewFingerprint([]byte{0b11111111})

	defer func() {
		if r := recover(); r == nil {
			t.Error("SetKeyFP failed to panic when the first bit of the " +
				"provided data is not 0.")
		}
	}()

	msg.SetKeyFP(fp)
}

// Happy path.
func TestMessage_GetMac(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	mac := makeAndFillSlice(MacLen, 'm')
	copy(msg.mac, mac)
	msg.mac[0] |= 0b10000000

	if !bytes.Equal(mac, msg.GetMac()) {
		t.Errorf("GetMac failed to get the correct MAC."+
			"\nexpected: %+v\nreceived: %+v", mac, msg.GetMac())
	}

	// Ensure that the data is copied
	mac[2] = 'x'
	if msg.mac[2] == 'x' {
		t.Error("GetMac failed to make a copy of mac.")
	}

	if msg.GetMac()[0]&0b10000000 != 0 {
		t.Errorf("First bit not set to zero")
	}
}

// Happy path.
func TestMessage_SetMac(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	mac := makeAndFillSlice(MacLen, 'm')

	msg.SetMac(mac)

	if !bytes.Equal(mac, msg.mac) {
		t.Errorf("SetMac failed to set the MAC."+
			"\nexpected: %+v\nreceived: %+v", mac, msg.mac)
	}
}

// Error path: first bit of provided data is not 0.
func TestMessage_SetMac_FirstBitError(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	mac := make([]byte, MacLen)
	mac[0] = 0b11111111

	defer func() {
		if r := recover(); r == nil {
			t.Error("SetMac failed to panic when the first bit of the " +
				"provided data is not 0.")
		}
	}()

	msg.SetMac(mac)
}

// Error path: the length of the provided data is incorrect.
func TestMessage_SetMac_LenError(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	mac := makeAndFillSlice(MacLen+1, 'm')

	defer func() {
		if r := recover(); r == nil {
			t.Error("SetMac failed to panic when the length of the provided " +
				"MAC is wrong.")
		}
	}()

	msg.SetMac(mac)
}

// Happy path.
func TestMessage_GetEphemeralRID(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	ephemeralRID := makeAndFillSlice(EphemeralRIDLen, 'e')
	copy(msg.ephemeralRID, ephemeralRID)

	if !bytes.Equal(ephemeralRID, msg.GetEphemeralRID()) {
		t.Errorf("GetEphemeralRID failed to get the correct ephemeralRID."+
			"\nexpected: %+v\nreceived: %+v", ephemeralRID, msg.GetEphemeralRID())
	}

	// Ensure that the data is copied
	ephemeralRID[2] = 'x'
	if msg.ephemeralRID[2] == 'x' {
		t.Error("GetEphemeralRID failed to make a copy of ephemeralRID.")
	}
}

// Happy path.
func TestMessage_SetEphemeralRID(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	ephemeralRID := makeAndFillSlice(EphemeralRIDLen, 'e')

	// Ensure that the data is copied
	msg.SetEphemeralRID(ephemeralRID)
	if !bytes.Equal(ephemeralRID, msg.ephemeralRID) {
		t.Errorf("SetEphemeralRID failed to set the ephemeralRID."+
			"\nexpected: %+v\nreceived: %+v", ephemeralRID, msg.ephemeralRID)
	}
}

// Error path: provided ephemeral recipient ID data too short.
func TestMessage_SetEphemeralRID_LengthError(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetEphemeralRID failed to panic when the length of " +
				"the provided data is incorrect.")
		}
	}()

	msg.SetEphemeralRID(make([]byte, EphemeralRIDLen*2))
}

// Happy path.
func TestMessage_GetIdentityFP(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	identityFP := makeAndFillSlice(SIHLen, 'e')
	copy(msg.sih, identityFP)

	if !bytes.Equal(identityFP, msg.GetSIH()) {
		t.Errorf("GetSIH failed to get the correct sih."+
			"\nexpected: %+v\nreceived: %+v", identityFP, msg.GetSIH())
	}

	// Ensure that the data is copied
	identityFP[2] = 'x'
	if msg.sih[2] == 'x' {
		t.Error("GetSIH failed to make a copy of sih.")
	}
}

// Happy path.
func TestMessage_SetIdentityFP(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	identityFP := makeAndFillSlice(SIHLen, 'e')

	msg.SetSIH(identityFP)
	if !bytes.Equal(identityFP, msg.sih) {
		t.Errorf("SetSIH failed to set the sih."+
			"\nexpected: %+v\nreceived: %+v", identityFP, msg.sih)
	}
}

// Error path: size of provided data is incorrect.
func TestMessage_SetIdentityFP_LengthError(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetSIH failed to panic when the length of " +
				"the provided data is incorrect.")
		}
	}()

	msg.SetSIH(make([]byte, SIHLen*2))
}

// Tests that digests come out correctly and are different.
func TestMessage_Digest(t *testing.T) {

	expectedA := "/9SqCYEP3uUixw1ua1D7"
	expectedB := "v+183UhPfK61KCNeSClT"

	msgA := NewMessage(MinimumPrimeSize)

	contentsA := makeAndFillSlice(MinimumPrimeSize*2-AssociatedDataSize-1, 'a')

	msgA.SetContents(contentsA)

	digestA := msgA.Digest()

	if digestA != expectedA {
		t.Errorf("Digest A does not match expected: "+
			"DigestA: %s, Expected: %s", digestA, expectedA)
	}

	msgB := NewMessage(MinimumPrimeSize)

	contentsB := makeAndFillSlice(MinimumPrimeSize*2-AssociatedDataSize-1, 'b')

	msgB.SetContents(contentsB)

	digestB := msgB.Digest()

	if digestB != expectedB {
		t.Errorf("Digest B does not match expected: "+
			"DigestB: %s, Expected: %s", digestB, expectedB)
	}

	if digestA == digestB {
		t.Errorf("Digest A and Digest B are the same, they "+
			"should be diferent; A: %s, B: %s", digestA, digestB)
	}
}

// Unit test of Message.GoString.
func TestMessage_GoString(t *testing.T) {
	// Create message
	msg := NewMessage(MinimumPrimeSize)
	msg.SetKeyFP(NewFingerprint(makeAndFillSlice(KeyFPLen, 'c')))
	msg.SetMac(makeAndFillSlice(MacLen, 'd'))
	msg.SetEphemeralRID(makeAndFillSlice(EphemeralRIDLen, 'e'))
	msg.SetSIH(makeAndFillSlice(SIHLen, 'f'))
	msg.SetContents(makeAndFillSlice(MinimumPrimeSize*2-AssociatedDataSize-1, 'g'))

	expected :=
		"format.Message{keyFP:Y2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2M=, " +
			"MAC:ZGRkZGRkZGRkZGRkZGRkZGRkZGRkZGRkZGRkZGRkZGQ=, " +
			"ephemeralRID:7306357456645743973, " +
			"sih:ZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZg==, " +
			"contents:\"gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg\"}"

	if expected != msg.GoString() {
		t.Errorf("GoString returned incorrect string."+
			"\nexpected: %s\nreceived: %s", expected, msg.GoString())
	}
}

// Unit test of Message.GoString with an empty Message.
func TestMessage_GoString_EmptyMessage(t *testing.T) {
	var msg Message

	expected := "format.Message{keyFP:<nil>, MAC:<nil>, " +
		"ephemeralRID:<nil>, sih:<nil>, contents:\"\"}"

	if expected != msg.GoString() {
		t.Errorf("GoString returned incorrect string."+
			"\nexpected: %s\nreceived: %s", expected, msg.GoString())
	}
}

func TestMessage_SetGroupBits(t *testing.T) {

	var msgsToTest []Message

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			msg := generateMsg()
			if i == 1 {
				msg.payloadA[0] |= 0b10000000
			}
			if j == 1 {
				msg.payloadB[0] |= 0b10000000
			}
			msgsToTest = append(msgsToTest, msg)
		}
	}

	count := 0
	for _, msg := range msgsToTest {
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				msg.SetGroupBits(i == 1, j == 1)
				if int(msg.payloadA[0])>>7 != i {
					t.Errorf("first group bit not set. Expected %t, "+
						"got: %t on test %d-A", i == 1, int(msg.payloadA[0])>>8 == 1, count)
				}
				if int(msg.payloadB[0])>>7 != j {
					t.Errorf("second group bit not set. Expected %t, "+
						"got: %t on test %d-B", j == 1, int(msg.payloadB[0])>>8 == 1, count)
				}
				count++
			}
		}
	}
}

func TestSetFirstBit(t *testing.T) {
	b := []byte{0, 0, 0}
	setFirstBit(b, true)
	if b[0] != 0b10000000 {
		t.Errorf("first bit did not set")
	}

	b = []byte{255, 0, 0}
	setFirstBit(b, false)
	if b[0] != 0b01111111 {
		t.Errorf("first bit did not get unset set")
	}
}

func generateMsg() Message {
	msg := NewMessage(MinimumPrimeSize)

	// Created expected data
	var expectedRawContents []byte
	keyFP := makeAndFillSlice(KeyFPLen, 'a')
	mac := makeAndFillSlice(MacLen, 'b')
	contents1 := makeAndFillSlice(MinimumPrimeSize-KeyFPLen, 'c')
	contents2 := makeAndFillSlice(MinimumPrimeSize-MacLen-RecipientIDLen, 'd')
	expectedRawContents = append(expectedRawContents, keyFP...)
	expectedRawContents = append(expectedRawContents, contents1...)
	expectedRawContents = append(expectedRawContents, mac...)
	expectedRawContents = append(expectedRawContents, contents2...)

	// Copy contents into message
	copy(msg.keyFP, keyFP)
	copy(msg.mac, mac)
	copy(msg.contents1, contents1)
	copy(msg.contents2, contents2)

	return msg
}

// makeAndFillSlice creates a slice of the specified size filled with the
// specified rune.
func makeAndFillSlice(size int, r rune) []byte {
	buff := make([]byte, size)
	buff = bytes.Map(func(r2 rune) rune { return r }, buff)
	return buff
}
