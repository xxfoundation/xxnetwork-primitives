////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import "testing"

// Guarantees that all fields of a NewMessage on init are not nil
func TestNewMessage_Init(t *testing.T) {
	msg := NewMessage()
	if msg.Payload == nil || msg.AssociatedData == nil {
		t.Error("An embedded struct was nil")
	}
}
