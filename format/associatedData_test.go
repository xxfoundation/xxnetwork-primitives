////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"bytes"
	"errors"
	"fmt"
	"gitlab.com/elixxir/primitives/id"
	"math/rand"
	"testing"
)

// Ensures that you can read to and write from a field, and that the contents
// are represented in the serialization
// Not currently tested: any overlap of fields
func testField(get func() []byte, set func([]byte) int, ser func() []byte,
	length int) error {
	// Ensure expected field length
	testLen := len(get())
	if testLen != length {
		return errors.New(fmt.Sprintf("Test len was %v, "+
			"but expected len was %v", testLen, length))
	}

	// Populate field with data
	r := rand.New(rand.NewSource(0))
	testBytes := make([]byte, testLen)
	_, err := r.Read(testBytes)
	if err != nil {
		return err
	}
	setLen := set(testBytes)
	if setLen != testLen {
		return errors.New(fmt.Sprintf("Test len was %v, "+
			"but %v bytes were set", testLen, setLen))
	}

	// Ensure the field was populated and serialized
	if !bytes.Equal(get(), testBytes) {
		return errors.New(fmt.Sprintf("Got %v from field, but expected %v",
			get(), testBytes))
	}
	if !bytes.Contains(ser(), testBytes) {
		return errors.New("Test data wasn't included in the serialization")
	}
	return nil
}

func testVariableField(get func() []byte, set func([]byte) int, ser func() []byte,
	length int) error {

	// Populate field with data
	r := rand.New(rand.NewSource(0))
	testBytes := make([]byte, length)
	_, err := r.Read(testBytes)
	if err != nil {
		return err
	}
	setLen := set(testBytes)
	if setLen != length {
		return errors.New(fmt.Sprintf("Test len was %v, "+
			"but %v bytes were set", length, setLen))
	}

	// Ensure the field was populated and serialized
	if !bytes.Equal(get(), testBytes) {
		return errors.New(fmt.Sprintf("Got %v from field, but expected %v",
			get(), testBytes))
	}
	if !bytes.Contains(ser(), testBytes) {
		return errors.New("Test data wasn't included in the serialization")
	}
	return nil
}

// Make sure that SetRecipient and SetSender set the field correctly with id.User
func TestSetUser(t *testing.T) {
	u := new(id.User).SetUints(&[4]uint64{3298561, 1083657, 2836259, 187653})
	payload := NewPayload()
	payload.SetSender(u)
	if !id.Equal(u, payload.GetSender()) {
		t.Errorf("Sender not set correctly. Got: %x, expected %x",
			payload.GetSender(), u)
	}
	data := NewAssociatedData()
	data.SetRecipient(u)
	if !id.Equal(u, data.GetRecipient()) {
		t.Errorf("Recipient not set correctly. Got: %x, expected %x",
			data.GetRecipient(), u)
	}
}

// Guarantees that all fields of a NewAssociatedData are the valid length
func TestNewAssociatedData(t *testing.T) {
	associatedData := NewAssociatedData()

	if len(associatedData.recipientID) != AD_RID_LEN {
		t.Error("Recipient length wasn't right")
	}
	if len(associatedData.keyFingerprint) != AD_KEYFP_LEN {
		t.Error("Keyfp length wasn't right")
	}
	if len(associatedData.timestamp) != AD_TIMESTAMP_LEN {
		t.Error("Timestamp length wasn't right")
	}
	if len(associatedData.mac) != AD_MAC_LEN {
		t.Error("MAC length wasn't right")
	}
	if len(associatedData.rmic) != AD_RMIC_LEN {
		t.Error("MAC length wasn't right")
	}
}

func TestAssociatedData_RecipientID(t *testing.T) {
	data := NewAssociatedData()
	var err error
	err = testField(data.GetRecipientID,
		data.SetRecipientID,
		data.SerializeAssociatedData,
		AD_RID_LEN)
	if err != nil {
		t.Errorf("Recipient ID failed: %v", err.Error())
	}

	err = testField(data.GetRecipient().Bytes,
		data.SetRecipientID,
		data.SerializeAssociatedData,
		AD_RID_LEN)
	if err != nil {
		t.Errorf("Recipient ID by id.User failed: %v", err.Error())
	}
}

func TestAssociatedData_KeyFingerprint(t *testing.T) {
	data := NewAssociatedData()

	fp := Fingerprint{}

	if data.GetKeyFingerprint() != fp {
		t.Errorf("Finger not initialized properly")
	}

	fp = Fingerprint{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31}

	data.SetKeyFingerprint(fp)
	if data.GetKeyFingerprint() != fp {
		t.Errorf("Fingerprint failed")
	}

	fp = Fingerprint{1,0,3,2,5,4,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31}

	data.SetKeyFingerprint(fp)

	if data.GetKeyFingerprint() != fp {
		t.Errorf("Fingerprint failed")
	}

	serAssociatedData := data.SerializeAssociatedData()
	fpFromSer := serAssociatedData[AD_KEYFP_START:AD_KEYFP_END]

	if bytes.Compare(fp[:], fpFromSer) != 0 {
		t.Errorf("Fingerprint get doesn't match serijalized associated data")
	}
}

func TestAssociatedData_Mac(t *testing.T) {
	data := NewAssociatedData()
	var err error

	err = testField(data.GetMAC,
		data.SetMAC,
		data.SerializeAssociatedData,
		AD_MAC_LEN)
	if err != nil {
		t.Errorf("MAC failed: %v", err.Error())
	}
}

func TestAssociatedData_RecipientMIC(t *testing.T) {
	data := NewAssociatedData()
	var err error

	err = testField(data.GetRecipientMIC,
		data.SetRecipientMIC,
		data.SerializeAssociatedData,
		AD_RMIC_LEN)
	if err != nil {
		t.Errorf("Recipient MIC failed: %v", err.Error())
	}
}

func TestAssociatedData_Timestamp(t *testing.T) {
	data := NewAssociatedData()
	var err error

	err = testField(data.GetTimestamp, data.SetTimestamp,
		data.SerializeAssociatedData, AD_TIMESTAMP_LEN)
	if err != nil {
		t.Errorf("Timestamp failed: %v", err.Error())
	}
}
