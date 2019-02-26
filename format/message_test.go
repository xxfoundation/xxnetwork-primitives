////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"gitlab.com/elixxir/primitives/id"
	"math/rand"
	"testing"
)

func TestNewMessage(t *testing.T) {

	tests := uint64(3)

	testStrings := [][]byte{
		testText[0 : DATA_LEN/2],
		testText[0:DATA_LEN],
		testText[0 : 2*DATA_LEN],
	}

	expectedSlices := make([][]byte, tests)

	expectedSlices[0] = []byte(testStrings[0])
	expectedSlices[1] = []byte(testStrings[1])[0:DATA_LEN]
	// Since the third slice is too long to fit in one message, the third
	// test case should result in a message including everything that can
	// fit in one message, but also return an error.
	expectedSlices[2] = []byte(testStrings[2])[0:DATA_LEN]

	expectedErrors := []bool{false, false, true}

	for i := uint64(0); i < tests; i++ {
		msg, err := NewMessage(id.NewUserFromUint(i+1, t),
			id.NewUserFromUint(i+1, t),
			testStrings[i])

		// Make sure we get an error on the third string, which is too long
		if (err != nil) != expectedErrors[i] {
			t.Errorf("Didn't get the expected error from NewMessage at index"+
				" %v", i)
		}

		expectedSender := id.NewUserFromUint(i+1, t)
		if !bytes.Equal(msg.GetSender().Bytes(), expectedSender.Bytes()) {
			t.Errorf("Test of NewMessage failed on test %v: "+
				"sID did not match;\n  Expected: %v, Received: %v", i,
				i, msg.senderID)
		}

		expectedRecipient := id.NewUserFromUint(i+1, t)
		if !bytes.Equal(expectedRecipient.Bytes(), msg.GetRecipient().Bytes()) {
			t.Errorf("Test of NewMessage failed on test %v:, "+
				"rID did not match;\n  Expected: %v, Received: %v", i,
				i, msg.recipientID)
		}

		expct := expectedSlices[i]

		if !bytes.Contains(msg.data, expct) {
			t.Errorf("Test of NewMessage failed on test %v:, "+
				"bytes did not match;\n Value Expected: %v, Value Received: %v", i,
				hex.EncodeToString(expct), hex.EncodeToString(msg.data))
		}

		serial := msg.SerializeMessage()
		deserial := DeserializeMessage(serial)

		pldSuccess, pldErr := payloadEqual(msg.Message, deserial.Message)

		if !pldSuccess {
			t.Errorf("Test of NewMessage failed on test %v:, "+
				"postserial Payload did not match: %s", i, pldErr)
		}

		rcpSuccess, rcpErr := recipientEqual(msg.AssociatedData,
			deserial.AssociatedData)

		if !rcpSuccess {
			t.Errorf("Test of NewMessage failed on test %v:, "+
				"postserial Recipient did not match: %s", i, rcpErr)
		}

	}

}

func payloadEqual(p1 *Message, p2 *Message) (bool, string) {
	e := hex.EncodeToString
	// Use Contains instead of Equal here because the byte slice includes
	// trailing zeroes after the end of the string. Package users are
	// responsible for trimming these trailing zeroes currently. Once we migrate
	// to a better padding scheme this will become unnecessary.
	if !bytes.Contains(p2.data, p1.data) {
		return false, fmt.Sprintf("data; Expected %v, Recieved: %v",
			e(p1.data), e(p2.data))
	}

	if !bytes.Equal(p1.senderID, p2.senderID) {
		return false, fmt.Sprintf("sender; Expected %v, Recieved: %v",
			e(p1.senderID), e(p2.senderID))
	}

	if !bytes.Equal(p1.messageMIC, p2.messageMIC) {
		return false, fmt.Sprintf("messageMIC; Expected %v, Recieved: %v",
			e(p1.messageMIC), e(p2.messageMIC))
	}

	if !bytes.Equal(p1.messageInitVect, p2.messageInitVect) {
		return false, fmt.Sprintf("messageInitVect; Expected %v, Recieved: %v",
			e(p1.messageInitVect), e(p2.messageInitVect))
	}

	return true, ""
}

func recipientEqual(r1 *AssociatedData, r2 *AssociatedData) (bool, string) {
	e := hex.EncodeToString
	if !bytes.Equal(r1.recipientID, r2.recipientID) {
		return false, fmt.Sprintf("recipientID; Expected %v, Recieved: %v",
			e(r1.recipientID), e(r2.recipientID))
	}

	if !bytes.Equal(r1.recipientMIC, r2.recipientMIC) {
		return false, fmt.Sprintf("recipientMIC; Expected %v, Recieved: %v",
			e(r1.recipientMIC), e(r2.recipientMIC))
	}

	if !bytes.Equal(r1.recipientInitVect, r2.recipientInitVect) {
		return false, fmt.Sprintf("messageInitVect; Expected %v, Recieved: %v",
			e(r1.recipientInitVect), e(r2.recipientInitVect))
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
	_, err := NewMessage(new(id.User), new(id.User), []byte("some text"))
	if err == nil {
		t.Error("Didn't get an expected error from creating new message to" +
			" zero user")
	}

	// However, message designated for sending from zero user ID should
	// create successfully. Populating the sender ID should be optional for some
	// use cases (untraceable return address, for example.) At the time of
	// writing, the infrastructure required to support communications that don't
	// specify a return ID hasn't been built.
	_, err2 := NewMessage(new(id.User), id.NewUserFromUint(5,
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
	msg := MaryPoppins{
		Message: &Message{
			data: data,
		},
		AssociatedData: &AssociatedData{},
	}
	if !bytes.Equal(data, msg.GetPayload()) {
		t.Errorf("MaryPoppins payload was %q, expected %q", msg.GetPayload(), data)
	}
}

func TestMessage_GetRecipient(t *testing.T) {
	recipient := make([]byte, RID_LEN)
	rng := rand.New(rand.NewSource(9319))
	_, err := rng.Read(recipient)
	if err != nil {
		t.Error(err.Error())
	}
	msg := MaryPoppins{Message: &Message{},
		AssociatedData: &AssociatedData{
			recipientID: recipient,
		}}
	if !bytes.Equal(recipient, msg.GetRecipient()[:]) {
		t.Errorf("MaryPoppins recipient was %q, expected %q",
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
	msg := MaryPoppins{Message: &Message{senderID: sender},
		AssociatedData: &AssociatedData{}}
	if !bytes.Equal(sender, msg.GetSender()[:]) {
		t.Errorf("MaryPoppins sender was %q, expected %q",
			*msg.GetSender(), sender)
	}
}
