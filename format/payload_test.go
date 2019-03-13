package format

import "testing"


// Guarantees that all fields of a NewPayload are the valid length
func TestNewPayload(t *testing.T) {
	payload := NewPayload()
	if len(payload.payloadData) != MP_PAYLOAD_LEN {
		t.Error("Payload length wasn't right")
	}
	if len(payload.senderID) != MP_SID_LEN {
		t.Error("Sender length wasn't right")
	}
}

// Test each field of the payload
func TestPayload(t *testing.T) {
	payload := NewPayload()
	var err error
	err = testField(payload.GetSenderID, payload.SetSenderID,
		payload.SerializePayload, MP_SID_LEN)
	if err != nil {
		t.Errorf("Sender ID failed: %v", err.Error())
	}
	err = testField(payload.GetSender().Bytes, payload.SetSenderID,
		payload.SerializePayload, MP_SID_LEN)
	if err != nil {
		t.Errorf("Sender ID by id.User failed: %v", err.Error())
	}
	// These functions return variable length according to size of actual data
	// so must be tested with different function
	err = testVariableField(payload.GetPayloadData, payload.SetPayloadData,
		payload.SerializePayload, MP_PAYLOAD_LEN)
	if err != nil {
		t.Errorf("Payload Data failed: %v", err.Error())
	}
	err = testVariableField(payload.GetPayload, payload.SetPayload,
		payload.SerializePayload, TOTAL_LEN)
	if err != nil {
		t.Errorf("Payload failed: %v", err.Error())
	}
	err = testVariableField(payload.GetPayload, payload.SetSplitPayload,
		payload.SerializePayload, MP_PAYLOAD_LEN)
	if err != nil {
		t.Errorf("Payload Split failed: %v", err.Error())
	}
}