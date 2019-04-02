////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package switchboard

import (
	"bytes"
	"gitlab.com/elixxir/primitives/id"
	"sync"
	"testing"
	"time"
	"gitlab.com/elixxir/primitives/format"
)

type MockListener struct {
	NumHeard        int
	IsFallback      bool
	LastMessage     []byte
	LastMessageType int32
	mux             sync.Mutex
}

type Message struct {
	Contents  []byte
	Sender    *id.User
	CryptoType format.CryptoType
	MessageType int32
}

func (m *Message) GetSender() *id.User {
	return m.Sender
}

func (m *Message) GetMessageType() int32 {
	return m.MessageType
}

func (m *Message) GetCryptoType() format.CryptoType {
	return m.CryptoType
}

func (ml *MockListener) Hear(item Item, isHeardElsewhere bool) {
	ml.mux.Lock()
	defer ml.mux.Unlock()

	msg := item.(*Message)

	if !isHeardElsewhere || !ml.IsFallback {
		ml.NumHeard++
		ml.LastMessage = msg.Contents
		ml.LastMessageType = msg.GetMessageType()
	}
}

var specificUser = new(id.User).SetUints(&[4]uint64{0, 0, 0, 5})
var specificMessageType int32 = 5
var specificCryptoType format.CryptoType = 2
var delay = 10 * time.Millisecond

func OneListenerSetup() (*Switchboard, *MockListener) {
	var listeners *Switchboard
	listeners = NewSwitchboard()
	// add one listener to the map
	fullyMatchedListener := &MockListener{}
	// TODO different type for message types?
	listeners.Register(specificUser, specificCryptoType, specificMessageType,
		fullyMatchedListener)
	return listeners, fullyMatchedListener
}

func TestListenerMap_SpeakOne(t *testing.T) {
	// set up
	listeners, fullyMatchedListener := OneListenerSetup()

	// speak
	listeners.Speak(&Message{
		Contents:  []byte("hmmmm"),
		Sender:    specificUser,
		MessageType: specificMessageType,
		CryptoType: specificCryptoType,
	})

	// determine whether the listener heard the message
	time.Sleep(delay)
	expected := 1
	if fullyMatchedListener.NumHeard != 1 {
		t.Errorf("The listener heard %v messages instead of %v",
			fullyMatchedListener.NumHeard, expected)
	}
}

func TestListenerMap_SpeakManyToOneListener(t *testing.T) {
	// set up
	listeners, fullyMatchedListener := OneListenerSetup()

	// speak
	for i := 0; i < 20; i++ {
		go listeners.Speak(&Message{
			Contents:  make([]byte, 0),
			Sender:    specificUser,
			MessageType: specificMessageType,
			CryptoType: specificCryptoType,
		})
	}

	// determine whether the listener heard the message
	time.Sleep(delay)
	expected := 20
	if fullyMatchedListener.NumHeard != expected {
		t.Errorf("The listener heard %v messages instead of %v",
			fullyMatchedListener.NumHeard, expected)
	}
}

func TestListenerMap_SpeakToAnother(t *testing.T) {
	// set up
	listeners, fullyMatchedListener := OneListenerSetup()

	// speak
	listeners.Speak(&Message{
		MessageType: specificMessageType,
		CryptoType: specificCryptoType,
		Contents:  make([]byte, 0),
		Sender:    nonzeroUser,
	})

	// determine whether the listener heard the message
	time.Sleep(delay)
	expected := 0
	if fullyMatchedListener.NumHeard != expected {
		t.Errorf("The listener heard %v messages instead of %v",
			fullyMatchedListener.NumHeard, expected)
	}
}

func TestListenerMap_SpeakDifferentType(t *testing.T) {
	// set up
	listeners, fullyMatchedListener := OneListenerSetup()

	// speak
	listeners.Speak(&Message{
		MessageType: specificMessageType + 1,
		CryptoType: specificCryptoType + 1,
		Contents:  make([]byte, 0),
		Sender:    specificUser,
	})

	// determine whether the listener heard the message
	time.Sleep(delay)
	expected := 0
	if fullyMatchedListener.NumHeard != expected {
		t.Errorf("The listener heard %v messages instead of %v",
			fullyMatchedListener.NumHeard, expected)
	}
}

var zeroUser = id.ZeroID
var nonzeroUser = new(id.User).SetUints(&[4]uint64{0, 0, 0, 786})
var zeroMessageType int32
var zeroCryptoType format.CryptoType

func WildcardListenerSetup() (*Switchboard, *MockListener) {
	var listeners *Switchboard
	listeners = NewSwitchboard()
	// add one listener to the map
	wildcardListener := &MockListener{}
	// TODO different type for message types?
	listeners.Register(zeroUser, zeroCryptoType, zeroMessageType,
		wildcardListener)
	return listeners, wildcardListener
}

func TestListenerMap_SpeakWildcard(t *testing.T) {
	// set up
	listeners, wildcardListener := WildcardListenerSetup()

	// speak
	listeners.Speak(&Message{
		Contents:  make([]byte, 0),
		Sender:    specificUser,
		MessageType: specificMessageType + 1,
		CryptoType: specificCryptoType + 1,
	})

	// determine whether the listener heard the message
	time.Sleep(delay)
	expected := 1
	if wildcardListener.NumHeard != expected {
		t.Errorf("The listener heard %v messages instead of %v",
			wildcardListener.NumHeard, expected)
	}
}

func TestListenerMap_SpeakManyToMany(t *testing.T) {
	listeners := NewSwitchboard()

	individualListeners := make([]*MockListener, 0)

	// one user, many types
	for messageType := int32(1); messageType <= int32(20); messageType++ {
		newListener := MockListener{}
		listeners.Register(specificUser, specificCryptoType, messageType,
			&newListener)
		individualListeners = append(individualListeners, &newListener)
	}
	// wildcard listener for the user
	userListener := &MockListener{}
	listeners.Register(specificUser, zeroCryptoType, zeroMessageType, userListener)
	// wildcard listener for all messages
	wildcardListener := &MockListener{}
	listeners.Register(zeroUser, zeroCryptoType, zeroMessageType, wildcardListener)

	// send to all types for our user
	for messageType := int32(1); messageType <= int32(20); messageType++ {
		go listeners.Speak(&Message{
			MessageType: messageType,
			CryptoType: specificCryptoType,
			Contents:  make([]byte, 0),
			Sender:    specificUser,
		})
	}
	// send to all types for a different user
	otherUser := id.NewUserFromUint(98, t)
	for messageType := int32(1); messageType <= int32(20); messageType++ {
		go listeners.Speak(&Message{
			MessageType: messageType,
			CryptoType: specificCryptoType,
			Contents:  make([]byte, 0),
			Sender:    otherUser,
		})
	}

	time.Sleep(delay)

	expectedIndividuals := 1
	expectedUserWildcard := 20
	expectedAllWildcard := 40
	for i := 0; i < len(individualListeners); i++ {
		if individualListeners[i].NumHeard != expectedIndividuals {
			t.Errorf("Individual listener got %v messages, "+
				"expected %v messages", individualListeners[i].NumHeard, expectedIndividuals)
		}
	}
	if userListener.NumHeard != expectedUserWildcard {
		t.Errorf("User wildcard got %v messages, expected %v message",
			userListener.NumHeard, expectedUserWildcard)
	}
	if wildcardListener.NumHeard != expectedAllWildcard {
		t.Errorf("User wildcard got %v messages, expected %v message",
			wildcardListener.NumHeard, expectedAllWildcard)
	}
}

func TestListenerMap_SpeakFallback(t *testing.T) {
	var listeners *Switchboard
	listeners = NewSwitchboard()
	// add one normal and one fallback listener to the map
	fallbackListener := &MockListener{}
	fallbackListener.IsFallback = true
	listeners.Register(zeroUser, zeroCryptoType, zeroMessageType, fallbackListener)
	specificListener := &MockListener{}
	listeners.Register(specificUser, specificCryptoType, specificMessageType,
		specificListener)

	// send exactly one message to each of them
	listeners.Speak(&Message{
		MessageType: specificMessageType,
		CryptoType: specificCryptoType,
		Contents:  make([]byte, 0),
		Sender:    specificUser,
	})
	listeners.Speak(&Message{
		MessageType: specificMessageType + 1,
		CryptoType: specificCryptoType + 1,
		Contents:  make([]byte, 0),
		Sender:    specificUser,
	})

	time.Sleep(delay)

	expected := 1

	if specificListener.NumHeard != expected {
		t.Errorf("Specific listener: Expected %v, got %v messages", expected,
			specificListener.NumHeard)
	}
	if fallbackListener.NumHeard != expected {
		t.Errorf("Fallback listener: Expected %v, got %v messages", expected,
			fallbackListener.NumHeard)
	}
}

func TestListenerMap_SpeakBody(t *testing.T) {
	listeners, listener := OneListenerSetup()
	expected := []byte{0x01, 0x02, 0x03, 0x04}
	listeners.Speak(&Message{
		MessageType: specificMessageType,
		CryptoType: specificCryptoType,
		Contents:  expected,
		Sender:    specificUser,
	})
	time.Sleep(delay)
	if !bytes.Equal(listener.LastMessage, expected) {
		t.Errorf("Received message was %v, expected %v",
			listener.LastMessage, expected)
	}
	if listener.LastMessageType != specificMessageType {
		t.Errorf("Received message type was %v, expected %v",
			listener.LastMessageType, specificMessageType)
	}
}

func TestListenerMap_Unregister(t *testing.T) {
	listeners := NewSwitchboard()
	listenerID := listeners.Register(specificUser, specificCryptoType, specificMessageType,
		&MockListener{})
	listeners.Unregister(listenerID)
	if len(listeners.listeners[*specificUser][specificCryptoType][specificMessageType]) != 0 {
		t.Error("The listener was still in the map after we stopped" +
			" listening on it")
	}
}

// The following tests show correct behavior in certain type situations.
// In all cases, the listeners are listening to all users, because these tests
// are about types.
// This test demonstrates correct behavior when the crypto and message types
// are both specified.
func TestListenerMap_SpecificListener(t *testing.T) {
	listeners := NewSwitchboard()
	l := &MockListener{}
	listeners.Register(id.ZeroID, 5, 3, l)
	// Should match
	listeners.Speak(&Message{
		Contents:    []byte("Test 0"),
		Sender:      id.NewUserFromUint(8, t),
		CryptoType:  5,
		MessageType: 3,
	})
	if l.NumHeard != 1 {
		t.Error("Listener should have heard")
	}

	// Reset the count
	l.NumHeard = 0
	// Should not match
	listeners.Speak(&Message{
		Contents:    []byte("Test 1"),
		Sender:      id.NewUserFromUint(8, t),
		CryptoType:  0,
		MessageType: 3,
	})
	if l.NumHeard != 0 {
		t.Error("Listener should not have heard")
	}

	l.NumHeard = 0
	// Should not match
	listeners.Speak(&Message{
		Contents:    []byte("Test 2"),
		Sender:      id.NewUserFromUint(8, t),
		CryptoType:  5,
		MessageType: 0,
	})
	if l.NumHeard != 0 {
		t.Error("Listener should not have heard")
	}

	l.NumHeard = 0
	// Should not match
	listeners.Speak(&Message{
		Contents:    []byte("Test 3"),
		Sender:      id.NewUserFromUint(8, t),
		CryptoType:  0,
		MessageType: 0,
	})
	if l.NumHeard != 0 {
		t.Error("Listener should not have heard")
	}
}

func TestListenerMap_ZeroCryptoType(t *testing.T) {
	listeners := NewSwitchboard()
	l := &MockListener{}
	// This listener should always get exactly one message if the messageType
	// is 3. The crypto type used should not affect this behavior.
	listeners.Register(id.ZeroID, 0, 3, l)

	// Should match once
	listeners.Speak(&Message{
		Contents:    []byte("Test 0"),
		Sender:      id.NewUserFromUint(8, t),
		CryptoType:  5,
		MessageType: 3,
	})
	if l.NumHeard != 1 {
		t.Error("Listener should have heard")
	}

	// Reset the count
	l.NumHeard = 0
	// Should match once
	listeners.Speak(&Message{
		Contents:    []byte("Test 1"),
		Sender:      id.NewUserFromUint(8, t),
		CryptoType:  0,
		MessageType: 3,
	})
	if l.NumHeard != 1 {
		t.Error("Listener should have heard once")
	}

	// Reset the count
	l.NumHeard = 0
	// Should not match
	listeners.Speak(&Message{
		Contents:    []byte("Test 2"),
		Sender:      id.NewUserFromUint(8, t),
		CryptoType:  0,
		MessageType: 0,
	})
	if l.NumHeard != 0 {
		t.Error("Listener should not have heard")
	}

	// Reset the count
	l.NumHeard = 0
	// Should not match
	listeners.Speak(&Message{
		Contents:    []byte("Test 3"),
		Sender:      id.NewUserFromUint(8, t),
		CryptoType:  5,
		MessageType: 0,
	})
	if l.NumHeard != 0 {
		t.Error("Listener should not have heard")
	}
}
