////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package format

import (
	"bytes"
	id2 "gitlab.com/xx_network/primitives/id"
	"math/rand"
	"testing"
)

func TestContentsSize(t *testing.T) {

	m := NewMessage(256)

	prng := rand.New(rand.NewSource(42))

	for i := 0; i < 1000; i++ {
		size := uint16(prng.Uint64() % (1<<14 - 1))

		m.setContentsSize(size)

		gotSize := m.getContentsSize()

		if size != gotSize {
			t.Errorf("Reconstructed size not correct; "+
				"intial: %v, reconstructed: %v", size, gotSize)
		}

	}

}

func TestMessage_GetPrimeByteLen(t *testing.T) {
	const primeSize = 250
	m := NewMessage(primeSize)

	if m.GetPrimeByteLen() != primeSize {
		t.Errorf("returned prime size is incorrect")
	}
}

func TestMessage_Smoke(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	timestamp := make([]byte, TimestampLen)
	timestamp = bytes.Map(func(r rune) rune {
		return 'a'
	}, timestamp)

	fp := Fingerprint{}
	keyFp := make([]byte, KeyFPLen)
	keyFp = bytes.Map(func(r rune) rune {
		return 'c'
	}, keyFp)
	copy(fp[:], keyFp)

	mac := make([]byte, MacLen)
	mac = bytes.Map(func(r rune) rune {
		return 'd'
	}, mac)

	recipientId := id2.ID{}
	idData := make([]byte, RecipientIDLen)
	idData = bytes.Map(func(r rune) rune {
		return 'e'
	}, idData)
	copy(recipientId[:], idData)

	contents := make([]byte, MinimumPrimeSize*2-AssociatedDataSize)
	contents = bytes.Map(func(r rune) rune {
		return 'f'
	}, contents)

	msg.SetTimestamp(timestamp)

	msg.SetKeyFP(fp)

	msg.SetMac(mac)

	msg.SetRecipientID(&recipientId)

	msg.SetContents(contents)

	if bytes.Compare(idData, msg.recipientID) != 0 {
		t.Errorf("recipient ID was corrupted.  Original: %+v, Current: %+v", idData, msg.recipientID)
	}

	if bytes.Compare(timestamp, msg.timestamp) != 0 {
		t.Errorf("timestamp data was corrupted.  Original: %+v, Current: %+v", timestamp, msg.timestamp)
	}

	if bytes.Compare(mac, msg.mac) != 0 {
		t.Errorf("mac data was corrupted.  Original: %+v, Current: %+v", mac, msg.mac)
	}

	if bytes.Compare(keyFp, msg.keyFP) != 0 {
		t.Errorf("keyFp data was corrupted.  Original: %+v, Current: %+v", keyFp, msg.keyFP)
	}

	if bytes.Compare(append(msg.contents1, msg.contents2...), contents) != 0 {
		t.Errorf("contents data was corrupted.  Original: %+v, Current(pt1): %+v, Current(pt2: %+v",
			contents, msg.contents1, msg.contents2)
	}
}

func TestNewMessage_Panic(t *testing.T) {
	// Defer to an error when NewMessage() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewMessage() did not panic when expected")
		}
	}()
	_ = NewMessage(MinimumPrimeSize - 1)
}

func TestMessage_ContentsSize(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	if msg.ContentsSize() != MinimumPrimeSize*2-AssociatedDataSize {
		t.Errorf("Contents size somehow wrong")
	}
}

func TestMessage_Copy(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	msgCopy := msg.Copy()

	s := []byte("test")
	contents := make([]byte, MinimumPrimeSize*2-AssociatedDataSize)
	copy(contents, s)

	msgCopy.SetContents(contents)

	if bytes.Compare(msg.GetContents(), contents) == 0 {
		t.Errorf("The copy is still pointing at the original data")
	}
}

func TestMessage_GetContents(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	s := []byte("test")
	contents := make([]byte, MinimumPrimeSize*2-AssociatedDataSize)
	copy(contents, s)

	copy(msg.contents1, contents[:len(msg.contents1)])
	copy(msg.contents2, contents[len(msg.contents1):])
	msg.setContentsSize(uint16(len(contents)))

	retrieved := msg.GetContents()

	if bytes.Compare(retrieved, contents) != 0 {
		t.Errorf("Did not properly get contents of message: %+v", retrieved)
	}
}

func TestMessage_SetContents(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	c := make([]byte, MinimumPrimeSize*2-AssociatedDataSize)
	contents := bytes.Map(func(r rune) rune {
		return 'a'
	}, c)

	msg.SetContents(contents)

	if bytes.Compare(msg.contents1, contents[:len(msg.contents1)]) != 0 {
		t.Errorf("contents 1 not as expected")
	}
	if bytes.Compare(msg.contents2, contents[len(msg.contents1):]) != 0 {
		t.Errorf("contents 2 not as expected")
	}
}

func TestMessage_GetData(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	copy(msg.data, "test")
	data := msg.GetData()
	if string(data[:4]) != "test" {
		t.Errorf("Didn't properly retrieve data field")
	}

	data[42] = 'q'

	if msg.data[42] == 'q' {
		t.Errorf("Changes to retrieved data changed the message object")
	}
}

func TestMessage_SetKeyFP(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	fp := Fingerprint{}
	copy(fp[:], "test")
	msg.SetKeyFP(fp)

	setFp := Fingerprint{}
	copy(setFp[:], msg.keyFP)
	if bytes.Compare(fp[:], setFp[:]) != 0 {
		t.Errorf("Set fp %+v does not match original %+v", setFp, fp)
	}
}

func TestMessage_GetKeyFP(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	copy(msg.keyFP, "test")
	fp := msg.GetKeyFP()
	if string(fp[:4]) != "test" {
		t.Errorf("Didn't properly retrieve keyFP")
	}

	fp[14] = 'x'

	if msg.keyFP[14] == 'x' {
		t.Errorf("Change to retrieved fingerprint altered message field")
	}
}

func TestMessage_SetMac_WrongLen(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetMac() did not panic when given wrong length input")
		}
	}()

	msg := NewMessage(MinimumPrimeSize)

	msg.SetMac([]byte("mac"))
}

func TestMessage_SetMac_BadFormat(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetMac() did not panic when given input w/o first byte set to 0")
		}
	}()

	msg := NewMessage(MinimumPrimeSize)

	mac := make([]byte, MacLen)
	mac[0] |= 0x80
	msg.SetMac(mac)

}

func TestMessage_SetMac(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	mac := make([]byte, MacLen)
	copy(mac, "mac")
	mac[0] = 0
	msg.SetMac(mac)

	if bytes.Compare(msg.mac, mac) != 0 {
		t.Errorf("Failed to set mac field")
	}
}

func TestMessage_GetMac(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	copy(msg.mac, "test")
	mac := msg.GetMac()
	if string(mac[:4]) != "test" {
		t.Errorf("Didn't properly retrieve MAC")
	}

	mac[14] = 'x'

	if msg.mac[14] == 'x' {
		t.Errorf("Change to retrieved mac field altered message field")
	}
}

func TestMessage_SetPayloadA_WrongLen(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetPayloadA() did not panic when given input of wrong len")
		}
	}()

	msg := NewMessage(MinimumPrimeSize)

	msg.SetPayloadA([]byte("test"))
}

func TestMessage_SetPayloadA(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	payloadA := make([]byte, len(msg.payloadA))
	copy(payloadA, "test")
	msg.SetPayloadA(payloadA)

	if bytes.Compare(payloadA, msg.payloadA) != 0 {
		t.Errorf("Failed to set the payload a field properly")
	}
}

func TestMessage_GetPayloadA(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	copy(msg.payloadA, "test")
	payloadA := msg.GetPayloadA()
	if string(payloadA[:4]) != "test" {
		t.Errorf("Did not properly retrieve payload A")
	}

	payloadA[14] = 'x'

	if msg.payloadA[14] == 'x' {
		t.Errorf("Change to retreived payloadA field altered message field")
	}
}

func TestMessage_SetPayloadB_WrongLen(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetPayloadB() did not panic when given input of wrong len")
		}
	}()

	msg.SetPayloadB([]byte("test"))
}

func TestMessage_SetPayloadB(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	payloadB := make([]byte, len(msg.payloadB))
	copy(payloadB, "test")
	msg.SetPayloadB(payloadB)

	if bytes.Compare(msg.payloadB, payloadB) != 0 {
		t.Errorf("Did not set payloadB field properly")
	}
}

func TestMessage_GetPayloadB(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	copy(msg.payloadB, "test")
	payloadB := msg.GetPayloadB()
	if string(payloadB[:4]) != "test" {
		t.Errorf("Did not properly retrieve payload B")
	}

	payloadB[14] = 'x'

	if msg.payloadB[14] == 'x' {
		t.Errorf("Change to retreived payloadB field altered message field")
	}
}

func TestMessage_SetRecipientID(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	id := id2.NewIdFromString("testid", id2.Gateway, t)
	msg.SetRecipientID(id)

	if bytes.Compare(id[:], msg.recipientID) != 0 {
		t.Errorf("Did not set recipient id field properly")
	}
}

func TestMessage_GetRecipientID(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	copy(msg.recipientID, "test")
	recipientId := msg.GetRecipientID()
	if string(recipientId[:4]) != "test" {
		t.Errorf("Did not properly retrieve recipient ID")
	}

	recipientId[14] = 'x'

	if msg.recipientID[14] == 'x' {
		t.Errorf("Change to retrieved recipientID altered message field")
	}
}

func TestMessage_SetTimestamp_WrongLen(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetTimestamp() did not panic when given input of wrong len")
		}
	}()

	msg := NewMessage(MinimumPrimeSize)
	msg.SetTimestamp([]byte("test"))
}

func TestMessage_SetTimestamp(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	timestamp := make([]byte, TimestampLen)
	copy(timestamp, "test")
	msg.SetTimestamp(timestamp)

	if bytes.Compare(timestamp, msg.timestamp) != 0 {
		t.Errorf("Failed to set timestamp field properly")
	}
}

func TestMessage_GetTimestamp(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	copy(msg.timestamp, "test")
	timestamp := msg.GetTimestamp()
	if string(timestamp[:4]) != "test" {
		t.Errorf("Did not properly retrieve timestamp")
	}

	timestamp[14] = 'x'

	if msg.timestamp[14] == 'x' {
		t.Errorf("Change to retrieved timestamp altered message field")
	}
}

func TestMessage_GetSecretPayload(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	copy(msg.contents1, "contents1")
	copy(msg.contents2, "contents2")
	copy(msg.size, []byte("xx"))

	secret := msg.GetSecretPayload()
	if !bytes.Contains(secret, []byte("contents1")) {
		t.Errorf("SEcret payload did not include contents 1")
	}
	if !bytes.Contains(secret, []byte("contents2")) {
		t.Errorf("Secret payload did not include contents 2")
	}
	if !bytes.Contains(secret, []byte("xx")) {
		t.Errorf("Secret payload did not include len")
	}
}

func TestMessage_GetSecretPayloadSize(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)

	expectedLen := (2 * MinimumPrimeSize) - AssociatedDataSize + SizeLen

	if msg.GetSecretPayloadSize() != expectedLen {
		t.Errorf("Didn't get expected length")
	}
}

func TestMessage_SetSecretPayload(t *testing.T) {
	msg := NewMessage(MinimumPrimeSize)
	spLen := (2 * MinimumPrimeSize) - AssociatedDataSize + SizeLen
	sp := make([]byte, spLen)
	copy(sp[:SizeLen], "xx")

	c1 := make([]byte, len(msg.contents1))
	c1 = bytes.Map(func(r rune) rune {
		return 'a'
	}, c1)

	c2 := make([]byte, len(msg.contents2))
	c2 = bytes.Map(func(r rune) rune {
		return 'b'
	}, c2)

	copy(sp[SizeLen:SizeLen+len(msg.contents1)], c1)
	copy(sp[SizeLen+len(msg.contents1):SizeLen+len(msg.contents1)+len(msg.contents2)], c2)

	msg.SetSecretPayload(sp)

	if bytes.Contains(msg.size, []byte("a")) || bytes.Contains(msg.size, []byte("b")) ||
		!bytes.Contains(msg.size, []byte("xx")) {
		t.Errorf("Setting secret payload failed")
	}
	if bytes.Contains(msg.contents1, []byte("x")) || bytes.Contains(msg.contents1, []byte("b")) ||
		!bytes.Contains(msg.contents1, []byte("aa")) {
		t.Errorf("Setting secret payload failed")
	}
	if bytes.Contains(msg.contents2, []byte("x")) || bytes.Contains(msg.contents2, []byte("a")) ||
		!bytes.Contains(msg.contents2, []byte("bb")) {
		t.Errorf("Setting secret payload failed")
	}
}
