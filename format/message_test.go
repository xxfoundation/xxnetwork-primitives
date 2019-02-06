////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"fmt"
	"gitlab.com/elixxir/crypto/cyclic"
	"testing"
	"gitlab.com/elixxir/primitives/userid"
	"math/rand"
	"bytes"
)

func TestNewMessage(t *testing.T) {

	tests := uint64(3)

	testStrings := make([][]byte, tests)

	testStrings[0] = testText[0 : DATA_LEN/2]
	testStrings[1] = testText[0:DATA_LEN]

	testStrings[2] = testText[0 : 2*DATA_LEN]

	expectedSlices := make([][][]byte, tests)

	expectedSlices[0] = make([][]byte, 1)

	expectedSlices[0][0] = []byte(testStrings[0])

	expectedSlices[1] = make([][]byte, 2)

	expectedSlices[1][0] = ([]byte(testStrings[1]))[0:DATA_LEN]

	expectedSlices[2] = make([][]byte, 3)

	expectedSlices[2][0] = ([]byte(testStrings[2]))[0:DATA_LEN]
	expectedSlices[2][1] = ([]byte(testStrings[2]))[DATA_LEN : 2*DATA_LEN]
	expectedSlices[2][2] = ([]byte(testStrings[2]))[2*DATA_LEN:]

	for i := uint64(0); i < tests; i++ {
		msglst, _ := NewMessage(id.NewUserIDFromUint(i+1, t),
			id.NewUserIDFromUint(i+1, t),
			testStrings[i])

		for indx, msg := range msglst {

			if uint64(i+1) != msg.senderID.Uint64() {
				t.Errorf("Test of NewMessage failed on test %v:%v, "+
					"sID did not match;\n  Expected: %v, Received: %v", i,
					indx, i, msg.senderID)
			}

			if uint64(i+1) != msg.recipientID.Uint64() {
				t.Errorf("Test of NewMessage failed on test %v:%v, "+
					"rID did not match;\n  Expected: %v, Received: %v", i,
					indx, i, msg.recipientID)
			}

			expct := cyclic.NewIntFromBytes(expectedSlices[i][indx])

			if msg.data.Cmp(expct) != 0 {
				t.Errorf("Test of NewMessage failed on test %v:%v, "+
					"bytes did not match;\n Value Expected: %v, Value Received: %v", i,
					indx, expct.Text(16), msg.data.Text(16))
			}

			serial := msg.SerializeMessage()
			deserial := DeserializeMessage(serial)

			pldSuccess, pldErr := payloadEqual(msg.Payload, deserial.Payload)

			if !pldSuccess {
				t.Errorf("Test of NewMessage failed on test %v:%v, "+
					"postserial Payload did not match: %s", i, indx, pldErr)
			}

			rcpSuccess, rcpErr := recipientEqual(msg.Recipient,
				deserial.Recipient)

			if !rcpSuccess {
				t.Errorf("Test of NewMessage failed on test %v:%v, "+
					"postserial Recipient did not match: %s", i, indx, rcpErr)
			}

		}

	}

}

func payloadEqual(p1 Payload, p2 Payload) (bool, string) {
	if p1.data.Cmp(p2.data) != 0 {
		return false, fmt.Sprintf("data; Expected %v, Recieved: %v",
			p1.data.Text(16), p2.data.Text(16))
	}

	if p1.senderID.Cmp(p2.senderID) != 0 {
		return false, fmt.Sprintf("sender; Expected %v, Recieved: %v",
			p1.senderID.Text(16), p2.senderID.Text(16))
	}

	if p1.payloadMIC.Cmp(p2.payloadMIC) != 0 {
		return false, fmt.Sprintf("payloadMIC; Expected %v, Recieved: %v",
			p1.payloadMIC.Text(16), p2.payloadMIC.Text(16))
	}

	if p1.payloadInitVect.Cmp(p2.payloadInitVect) != 0 {
		return false, fmt.Sprintf("payloadInitVect; Expected %v, Recieved: %v",
			p1.payloadInitVect.Text(16), p2.payloadInitVect.Text(16))
	}

	return true, ""

}

func recipientEqual(r1 Recipient, r2 Recipient) (bool, string) {
	if r1.recipientID.Cmp(r2.recipientID) != 0 {
		return false, fmt.Sprintf("recipientID; Expected %v, Recieved: %v",
			r1.recipientID.Text(16), r2.recipientID.Text(16))
	}

	if r1.recipientEmpty.Cmp(r2.recipientEmpty) != 0 {
		return false, fmt.Sprintf("empty; Expected %v, Recieved: %v",
			r1.recipientEmpty.Text(16), r2.recipientEmpty.Text(16))
	}

	if r1.recipientMIC.Cmp(r2.recipientMIC) != 0 {
		return false, fmt.Sprintf("recipientMIC; Expected %v, Recieved: %v",
			r1.recipientMIC.Text(16), r2.recipientMIC.Text(16))
	}

	if r1.recipientInitVect.Cmp(r2.recipientInitVect) != 0 {
		return false, fmt.Sprintf("payloadInitVect; Expected %v, Recieved: %v",
			r1.recipientInitVect.Text(16), r2.recipientInitVect.Text(16))
	}

	return true, ""

}

//TODO: Test End cases, messages over 2x length, at max length, and others.
var testText = []byte("Lorem ipsum dolor sit amet, " +
	"consectetur adipiscing elit. Sed" +
	" maximus convallis libero in laoreet. Aenean venenatis auctor condimentum." +
	" Suspendisse sed sapien purus. Ut molestie, mauris id porta ultrices, justo" +
	" nisi bibendum diam, quis facilisis metus ipsum nec dui. Nunc turpis felis," +
	" tristique nec viverra non, ultricies at elit. Ut pretium erat non porta" +
	" bibendum. Cras diam nulla, lobortis vel commodo luctus, dapibus nec nunc." +
	" Pellentesque ac commodo orci. Pellentesque nec nisi maximus, varius odio" +
	" eget, suscipit est. In viverra pretium lobortis. Fusce quis efficitur " +
	" libero. Sed eleifend dictum nulla sed tempus. Donec a tristique dolor, " +
	" quis mattis tellus. Nullam massa elit, ullamcorper ac consectetur ut, " +
	" tincidunt vel erat. Vivamus ut mauris eu ligula pretium tristique id in " +
	" justo. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce" +
	" porttitor, massa non iaculis faucibus, magna metus venenatis nisi," +
	" sodales fringilla enim nulla a erat. Vestibulum posuere ligula a mi " +
	" mollis, quis sodales ipsum hendrerit. Duis a iaculis felis, at " +
	" tristique ligula. In vulputate arcu quam, sit amet consequat lorem" +
	" convallis varius. Donec efficitur semper metus, a sodales dolor " +
	" vestibulum eu. Aliquam et laoreet massa. Phasellus cursus ligula ac " +
	" gravida vehicula. Etiam vitae malesuada nunc. Nunc vitae massa ex. " +
	" Mauris ullamcorper, nunc et rutrum lacinia, est nulla consectetur ex," +
	" non faucibus nulla eros imperdiet justo. Aenean ut velit a odio pretium" +
	" dictum ac nec dui. Vestibulum vulputate nulla vel elit ornare maximus." +
	" Sed egestas diam vel arcu venenatis, nec pulvinar ligula placerat. " +
	" Praesent sed interdum magna. Integer in diam lacus. Sed congue enim eros," +
	" ut ultricies erat porttitor sed. Nullam neque risus, bibendum eu risus ut," +
	" fermentum viverra dolor. Cras non iaculis augue, id euismod metus. In hac" +
	" habitasse platea dictumst. Aenean convallis dignissim commodo. Duis ut" +
	" ultricies turpis. Duis mollis finibus mi dignissim efficitur. Maecenas" +
	" eleifend mi porttitor convallis sed.")

// Proves that NewMessage returns an error in the correct cases
func TestNewMessage_Errors(t *testing.T) {
	// Use new id.UserID because using id.
	// ZeroID would result in the pointers being equal as well
	// The test should rely on comparing the underlying data,
	// not the memory address
	// Creating message designated for sending to zero user ID should fail
	_, err := NewMessage(new(id.UserID), new(id.UserID), []byte("some text"))
	if err == nil {
		t.Error("Didn't get an expected error from creating new message to" +
			" zero user")
	}

	// However, message designated for sending from zero user ID should
	// create successfully. Populating the sender ID should be optional for some
	// use cases (untraceable return address, for example.) At the time of
	// writing, the infrastructure required to support communications that don't
	// specify a return ID hasn't been built.
	_, err2 := NewMessage(new(id.UserID), id.NewUserIDFromUint(5,
		t), []byte("some more text"))
	if err2 != nil {
		t.Errorf("Got an unexpected error from creating new message from zero"+
			" user: %v", err2.Error())
	}
}

// Proves that the data you get out of the payload is equal to the data you
// put in
func TestMessage_GetPayload(t *testing.T) {
	rng := rand.New(rand.NewSource(87321))
	data := make([]byte, DATA_LEN)
	_, err := rng.Read(data)
	if err != nil {
		t.Errorf("Got error from data generation: %s", err.Error())
	}
	msg := Message{
		Payload: Payload{
			data: cyclic.NewIntFromBytes(data),
		},
		Recipient: Recipient{},
	}
	if !bytes.Equal(data, msg.GetPayload()) {
		t.Errorf("Message payload was %q, expected %q", msg.GetPayload(), data)
	}
}

func TestMessage_GetRecipient(t *testing.T) {
	recipient := make([]byte, RID_LEN)
	rng := rand.New(rand.NewSource(9319))
	_, err := rng.Read(recipient)
	if err != nil {
		t.Error(err.Error())
	}
	msg := Message{Payload: Payload{},
		Recipient: Recipient{
			recipientID: cyclic.NewIntFromBytes(recipient),
		}}
	if !bytes.Equal(recipient, msg.GetRecipient()[:]) {
		t.Errorf("Message recipient was %q, expected %q",
			*msg.GetRecipient(), recipient)
	}
}

func TestMessage_GetSender(t *testing.T) {
	sender := make([]byte, SID_LEN)
	rng := rand.New(rand.NewSource(9812))
	_, err := rng.Read(sender)
	if err != nil {
		t.Error(err.Error())
	}
	msg := Message{Payload: Payload{senderID: cyclic.NewIntFromBytes(sender)},
		Recipient: Recipient{}}
	if !bytes.Equal(sender, msg.GetSender()[:]) {
		t.Errorf("Message sender was %q, expected %q",
			*msg.GetSender(), sender)
	}
}
