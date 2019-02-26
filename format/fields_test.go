////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format_test

import (
	"testing"
	"gitlab.com/elixxir/primitives/format"
	"math/rand"
	"errors"
	"fmt"
	"bytes"
)

// Ensures that you can read to and write from a field, and that the contents
// are represented in the serialization
// Not currently tested: any overlap of fields
func testField(get func() []byte, set func([]byte) int, ser func() []byte,
	length int) error {
	// Ensure expected field length
	testLen := len(get())
	if testLen != length {
		return errors.New(fmt.Sprintf("Test len was %v, " +
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
		return errors.New(fmt.Sprintf("Test len was %v, " +
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

// Test each field of the payload
func TestPayload(t *testing.T) {
	payload := format.NewPayload()
	var err error
	err = testField(payload.GetSenderID, payload.SetSenderID,
		payload.SerializePayload, format.MP_SID_LEN)
	if err != nil {
		t.Errorf("Sender ID failed: %v", err.Error())
	}
	err = testField(payload.GetSender().Bytes, payload.SetSenderID,
		payload.SerializePayload, format.MP_SID_LEN)
	if err != nil {
		t.Errorf("Sender ID by id.User failed: %v", err.Error())
	}
	err = testField(payload.GetPayload, payload.SetPayload,
		payload.SerializePayload, format.MP_PAYLOAD_LEN)
	if err != nil {
		t.Errorf("Payload failed: %v", err.Error())
	}
}

func TestAssociatedData(t *testing.T) {
	data := format.NewAssociatedData()
	var err error
	err = testField(data.GetRecipientID,
		data.SetRecipientID,
		data.SerializeAssociatedData,
		format.AD_RID_LEN)
	if err != nil {
		t.Errorf("Recipient ID failed: %v", err.Error())
	}
	err = testField(data.GetKeyFingerprint,
		data.SetKeyFingerprint,
		data.SerializeAssociatedData,
		format.AD_KEYFP_LEN)
	if err != nil {
		t.Errorf("Recipient ID failed: %v", err.Error())
	}
	err = testField(data.GetMAC,
		data.SetMAC,
		data.SerializeAssociatedData,
		format.AD_MAC_LEN)
	if err != nil {
		t.Errorf("Recipient ID failed: %v", err.Error())
	}
	err = testField(data.GetRecipient().Bytes,
		data.SetRecipientID,
		data.SerializeAssociatedData,
		format.AD_RID_LEN)
	if err != nil {
		t.Errorf("Recipient ID by id.User failed: %v", err.Error())
	}
}
