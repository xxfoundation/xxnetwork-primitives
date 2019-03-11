////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import "testing"

// Guarantees that all fields of a NewMessage are accessible
func TestNewMessage(t *testing.T) {
	msg := NewMessage()
	if msg.Payload == nil || msg.AssociatedData == nil {
		t.Error("An embedded struct was nil")
	}
	if len(msg.payloadData) != MP_PAYLOAD_LEN {
		t.Error("Payload length wasn't right")
	}
	if len(msg.senderID) != MP_SID_LEN {
		t.Error("Sender length wasn't right")
	}
	if len(msg.recipientID) != AD_RID_LEN {
		t.Error("Recipient length wasn't right")
	}
	if len(msg.keyFingerprint) != AD_KEYFP_LEN {
		t.Error("Keyfp length wasn't right")
	}
	if len(msg.timestamp) != AD_TIMESTAMP_LEN {
		t.Error("Timestamp length wasn't right")
	}
	if len(msg.mac) != AD_MAC_LEN {
		t.Error("MAC length wasn't right")
	}
	if len(msg.rmic) != AD_RMIC_LEN {
		t.Error("MAC length wasn't right")
	}
}
