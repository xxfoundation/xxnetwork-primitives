////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package switchboard

import (
	"bytes"
	"gitlab.com/elixxir/primitives/cmixproto"
	"gitlab.com/elixxir/primitives/id"
	"sync"
	"testing"
	"time"
)

type MockListener struct {
	NumHeard        int
	IsFallback      bool
	LastMessage     []byte
	LastMessageType cmixproto.InnerType
	mux             sync.Mutex
}

type Message struct {
	contents  []byte
	sender    *id.User
	innerType cmixproto.InnerType
	outerType cmixproto.OuterType
}

func (m *Message) GetSender() *id.User {
	return m.sender
}

func (m *Message) GetInnerType() cmixproto.InnerType {
	return m.innerType
}

func (m *Message) GetOuterType() cmixproto.OuterType {
	return m.outerType
}

func (ml *MockListener) Hear(item Item, isHeardElsewhere bool) {
	ml.mux.Lock()
	defer ml.mux.Unlock()

	msg := item.(*Message)

	if !isHeardElsewhere || !ml.IsFallback {
		ml.NumHeard++
		ml.LastMessage = msg.contents
		ml.LastMessageType = msg.GetInnerType()
	}
}

var specificUser = new(id.User).SetUints(&[4]uint64{0, 0, 0, 5})
var specificInnerType = cmixproto.InnerType_TEXT_MESSAGE
var specificOuterType = cmixproto.OuterType_E2E
var delay = 10 * time.Millisecond

func OneListenerSetup() (*Switchboard, *MockListener) {
	var listeners *Switchboard
	listeners = NewSwitchboard()
	// add one listener to the map
	fullyMatchedListener := &MockListener{}
	// TODO different type for message types?
	listeners.Register(specificUser, specificOuterType, specificInnerType,
		fullyMatchedListener)
	return listeners, fullyMatchedListener
}

func TestListenerMap_SpeakOne(t *testing.T) {
	// set up
	listeners, fullyMatchedListener := OneListenerSetup()

	// speak
	listeners.Speak(&Message{
		contents:  []byte("hmmmm"),
		sender:    specificUser,
		innerType: specificInnerType,
		outerType: specificOuterType,
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
			contents:  make([]byte, 0),
			sender:    specificUser,
			innerType: specificInnerType,
			outerType: specificOuterType,
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
		innerType: specificInnerType,
		outerType: specificOuterType,
		contents:  make([]byte, 0),
		sender:    nonzeroUser,
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
		innerType: specificInnerType + 1,
		outerType: specificOuterType + 1,
		contents:  make([]byte, 0),
		sender:    specificUser,
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
var zeroInnerType cmixproto.InnerType
var zeroOuterType cmixproto.OuterType

func WildcardListenerSetup() (*Switchboard, *MockListener) {
	var listeners *Switchboard
	listeners = NewSwitchboard()
	// add one listener to the map
	wildcardListener := &MockListener{}
	// TODO different type for message types?
	listeners.Register(zeroUser, zeroOuterType, zeroInnerType,
		wildcardListener)
	return listeners, wildcardListener
}

func TestListenerMap_SpeakWildcard(t *testing.T) {
	// set up
	listeners, wildcardListener := WildcardListenerSetup()

	// speak
	listeners.Speak(&Message{
		contents:  make([]byte, 0),
		sender:    specificUser,
		innerType: specificInnerType + 1,
		outerType: specificOuterType + 1,
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
	for messageType := cmixproto.InnerType(1); messageType <= cmixproto.
		InnerType(20); messageType++ {
		newListener := MockListener{}
		listeners.Register(specificUser, specificOuterType, messageType,
			&newListener)
		individualListeners = append(individualListeners, &newListener)
	}
	// wildcard listener for the user
	userListener := &MockListener{}
	listeners.Register(specificUser, zeroOuterType, zeroInnerType, userListener)
	// wildcard listener for all messages
	wildcardListener := &MockListener{}
	listeners.Register(zeroUser, zeroOuterType, zeroInnerType, wildcardListener)

	// send to all types for our user
	for messageType := cmixproto.InnerType(1); messageType <= cmixproto.
		InnerType(20); messageType++ {
		go listeners.Speak(&Message{
			innerType: messageType,
			outerType: specificOuterType,
			contents:  make([]byte, 0),
			sender:    specificUser,
		})
	}
	// send to all types for a different user
	otherUser := id.NewUserFromUint(98, t)
	for messageType := cmixproto.InnerType(1); messageType <= cmixproto.
		InnerType(20); messageType++ {
		go listeners.Speak(&Message{
			innerType: messageType,
			outerType: specificOuterType,
			contents:  make([]byte, 0),
			sender:    otherUser,
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
	listeners.Register(zeroUser, zeroOuterType, zeroInnerType, fallbackListener)
	specificListener := &MockListener{}
	listeners.Register(specificUser, specificOuterType, specificInnerType,
		specificListener)

	// send exactly one message to each of them
	listeners.Speak(&Message{
		innerType: specificInnerType,
		outerType: specificOuterType,
		contents:  make([]byte, 0),
		sender:    specificUser,
	})
	listeners.Speak(&Message{
		innerType: specificInnerType + 1,
		outerType: specificOuterType + 1,
		contents:  make([]byte, 0),
		sender:    specificUser,
	})

	time.Sleep(delay)

	expected := 1

	if specificListener.NumHeard != expected {
		t.Errorf("Specific listener: Expected %v, got %v messages", expected,
			specificListener.NumHeard)
	}
	if fallbackListener.NumHeard != expected {
		t.Errorf("Fallback listener: Expected %v, got %v messages", expected,
			specificListener.NumHeard)
	}
}

func TestListenerMap_SpeakBody(t *testing.T) {
	listeners, listener := OneListenerSetup()
	expected := []byte{0x01, 0x02, 0x03, 0x04}
	listeners.Speak(&Message{
		innerType: specificInnerType,
		outerType: specificOuterType,
		contents:  expected,
		sender:    specificUser,
	})
	time.Sleep(delay)
	if !bytes.Equal(listener.LastMessage, expected) {
		t.Errorf("Received message was %v, expected %v",
			listener.LastMessage, expected)
	}
	if listener.LastMessageType != specificInnerType {
		t.Errorf("Received message type was %v, expected %v",
			listener.LastMessageType, specificInnerType)
	}
}

func TestListenerMap_Unregister(t *testing.T) {
	listeners := NewSwitchboard()
	listenerID := listeners.Register(specificUser, specificOuterType, specificInnerType,
		&MockListener{})
	listeners.Unregister(listenerID)
	if len(listeners.listeners[*specificUser][specificOuterType][specificInnerType]) != 0 {
		t.Error("The listener was still in the map after we stopped" +
			" listening on it")
	}
}
