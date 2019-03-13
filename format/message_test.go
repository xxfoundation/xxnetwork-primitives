////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"gitlab.com/elixxir/primitives/id"
	"math/rand"
	"reflect"
	"testing"
)

// Guarantees that all fields of a NewMessage on init are not nil
func TestNewMessage_Init(t *testing.T) {
	msg := NewMessage()
	if msg.Payload == nil || msg.AssociatedData == nil {
		t.Error("An embedded struct was nil")
	}
}

func TestDeepCopy(t *testing.T) {
	// Generate test data for each structure
	r := rand.New(rand.NewSource(0))
	testBytes := make([]byte, TOTAL_LEN)
	_, err := r.Read(testBytes)
	if err != nil {
		t.Error(err.Error())
	}

	// Make a deep copy of each structure, and make sure that changing one
	// doesn't change the other
	data := DeserializeAssociatedData(testBytes)
	dataCopy := data.DeepCopy()
	payload := DeserializePayload(testBytes)
	payloadCopy := payload.DeepCopy()

	if !reflect.DeepEqual(data, dataCopy) {
		t.Error("Datas should have been equal before mutation, but weren't")
	}
	if !reflect.DeepEqual(payload, payloadCopy) {
		t.Error("Payloads should have been equal before mutation, but weren't")
	}
	// Mutate each copy. The originals and copies should now be different
	dataCopy.SetRecipient(id.NewUserFromUint(5, t))
	payloadCopy.SetSender(id.NewUserFromUint(5, t))
	if reflect.DeepEqual(data, dataCopy) {
		t.Error("Datas should have been different after mutation, but weren't")
	}
	if reflect.DeepEqual(payload, payloadCopy) {
		t.Error("Payloads should have been different after mutation, " +
			"but weren't")
	}
}
