package switchboard

import (
	"gitlab.com/elixxir/primitives/id"
	"strings"
	"testing"
	"time"
)

//testing message structure
type Message struct {
	Contents    []byte
	Sender      *id.ID
	MessageType int32
}

func (m *Message) GetSender() *id.ID {
	return m.Sender
}

func (m *Message) GetMessageType() int32 {
	return m.MessageType
}

func TestNew(t *testing.T) {
	sw := New()

	if sw.id == nil {
		t.Errorf("did not create an id map")
	}

	if sw.messageType == nil {
		t.Errorf("did not create a messageType map")
	}
}

func TestSwitchboard_RegisterListener_Errors(t *testing.T) {
	sw := New()
	_, err := sw.RegisterListener(nil, 0, &funcListener{})

	if err == nil {
		t.Errorf("A nil userID should have caused an error")
	}

	if err != nil && !strings.Contains(err.Error(), "cannot register listener to nil user") {
		t.Errorf("A nil userID caused the wrong error")
	}

	_, err = sw.RegisterListener(id.NewIdFromUInt(42, id.User, t), 0, nil)

	if err == nil {
		t.Errorf("A nil listener should have caused an error")
	}

	if err != nil && !strings.Contains(err.Error(), "cannot register nil listener") {
		t.Errorf("A nil listener caused the wrong error")
	}
}

func TestSwitchboard_RegisterListener(t *testing.T) {
	sw := New()

	l := &funcListener{}

	uid := id.NewIdFromUInt(42, id.User, t)

	mt := int32(69)

	lid, err := sw.RegisterListener(uid, mt, l)

	//check the returns
	if err != nil {
		t.Errorf("Register Listener should not have errored: %s", err)
	}

	if lid.messageType != mt {
		t.Errorf("ListenerID message type is wrong")
	}

	if !lid.userID.Cmp(uid) {
		t.Errorf("ListenerID userID is wrong")
	}

	if lid.listener != l {
		t.Errorf("ListenerID listener is wrong")
	}

	//check that the listener is registered in the appropriate location
	setID := sw.id.Get(uid)

	if !setID.Has(l) {
		t.Errorf("Listener is not registered by ID")
	}

	setType := sw.messageType.Get(mt)

	if !setType.Has(l) {
		t.Errorf("Listener is not registered by Message Type")
	}

}

func TestSwitchboard_RegisterFunc_Errors(t *testing.T) {
	sw := New()
	_, err := sw.RegisterFunc("test", nil, 0, func(Item) {})

	if err == nil {
		t.Errorf("A nil userID should have caused an error")
	}

	if err != nil && !strings.Contains(err.Error(), "cannot register listener to nil user") {
		t.Errorf("A nil userID caused the wrong error")
	}

	_, err = sw.RegisterFunc("test", id.NewIdFromUInt(42, id.User, t), 0, nil)

	if err == nil {
		t.Errorf("A nil listener func should have caused an error")
	}

	if err != nil && !strings.Contains(err.Error(), "cannot register listener with nil func") {
		t.Errorf("A nil listener caused func the wrong error")
	}
}

func TestSwitchboard_RegisterFunc(t *testing.T) {
	sw := New()

	heard := false

	l := func(Item) { heard = true }

	uid := id.NewIdFromUInt(42, id.User, t)

	mt := int32(69)

	lid, err := sw.RegisterFunc("test", uid, mt, l)

	//check the returns
	if err != nil {
		t.Errorf("Register Listener should not have errored: %s", err)
	}

	if lid.messageType != mt {
		t.Errorf("ListenerID message type is wrong")
	}

	if !lid.userID.Cmp(uid) {
		t.Errorf("ListenerID userID is wrong")
	}

	//check that the listener is registered in the appropriate location
	setID := sw.id.Get(uid)

	if !setID.Has(lid.listener) {
		t.Errorf("Listener is not registered by ID")
	}

	setType := sw.messageType.Get(mt)

	if !setType.Has(lid.listener) {
		t.Errorf("Listener is not registered by Message Type")
	}

	lid.listener.Hear(nil)
	if !heard {
		t.Errorf("Func listener not registered correctly")
	}
}

func TestSwitchboard_RegisterChan_Errors(t *testing.T) {
	sw := New()
	_, err := sw.RegisterChannel("test", nil, 0, make(chan Item))

	if err == nil {
		t.Errorf("A nil userID should have caused an error")
	}

	if err != nil && !strings.Contains(err.Error(), "cannot register listener to nil user") {
		t.Errorf("A nil userID caused the wrong error")
	}

	_, err = sw.RegisterChannel("test", id.NewIdFromUInt(42, id.User, t), 0, nil)

	if err == nil {
		t.Errorf("A nil channel func should have caused an error")
	}

	if err != nil && !strings.Contains(err.Error(), "cannot register listener with nil channel") {
		t.Errorf("A nil listener caused func the wrong error")
	}
}

func TestSwitchboard_RegisterChan(t *testing.T) {
	sw := New()

	ch := make(chan Item, 1)

	uid := id.NewIdFromUInt(42, id.User, t)

	mt := int32(69)

	lid, err := sw.RegisterChannel("test", uid, mt, ch)

	//check the returns
	if err != nil {
		t.Errorf("Register Listener should not have errored: %s", err)
	}

	if lid.messageType != mt {
		t.Errorf("ListenerID message type is wrong")
	}

	if !lid.userID.Cmp(uid) {
		t.Errorf("ListenerID userID is wrong")
	}

	//check that the listener is registered in the appropriate location
	setID := sw.id.Get(uid)

	if !setID.Has(lid.listener) {
		t.Errorf("Listener is not registered by ID")
	}

	setType := sw.messageType.Get(mt)

	if !setType.Has(lid.listener) {
		t.Errorf("Listener is not registered by Message Type")
	}

	lid.listener.Hear(nil)
	select {
	case <-ch:
	case <-time.After(5 * time.Millisecond):
		t.Errorf("Chan listener not registered correctly")
	}
}

//tests all combinations of hits and misses for speak
func TestSwitchboard_Speak(t *testing.T) {

	uids := []*id.ID{&id.ZeroID, id.NewIdFromUInt(42, id.User, t), id.NewIdFromUInt(69, id.User, t)}
	mts := []int32{AnyType, 42, 69}

	for _, uidReg := range uids {
		for _, mtReg := range mts {

			//create the registrations
			sw := New()
			ch1 := make(chan Item, 1)
			ch2 := make(chan Item, 1)

			_, err := sw.RegisterChannel("test", uidReg, mtReg, ch1)

			if err != nil {
				t.Errorf("Register Listener should not have errored: %s", err)
			}

			_, err = sw.RegisterChannel("test", uidReg, mtReg, ch2)

			if err != nil {
				t.Errorf("Register Listener should not have errored: %s", err)
			}

			//send every possible message
			for _, uid := range uids {
				for _, mt := range mts {
					if uid.Cmp(&id.ZeroID) || mt == AnyType {
						continue
					}

					m := &Message{
						Contents:    []byte{0, 1, 2, 3},
						Sender:      uid,
						MessageType: mt,
					}

					sw.Speak(m)

					shouldHear := (m.Sender.Cmp(uidReg) || uidReg.Cmp(&id.ZeroID)) && (m.MessageType == mtReg || mtReg == AnyType)

					var heard1 bool

					select {
					case <-ch1:
						heard1 = true
					case <-time.After(5 * time.Millisecond):
						heard1 = false
					}

					if shouldHear != heard1 {
						t.Errorf("Correct operation not recorded "+
							"for listener 1: Expected: %v, Occured: %v",
							shouldHear, heard1)
					}

					var heard2 bool

					select {
					case <-ch2:
						heard2 = true
					case <-time.After(5 * time.Millisecond):
						heard2 = false
					}

					if shouldHear != heard2 {
						t.Errorf("Correct operation not recorded "+
							"for listener 2: Expected: %v, Occured: %v",
							shouldHear, heard2)
					}
				}
			}
		}
	}
}

/*
import (
	"bytes"
	"gitlab.com/elixxir/primitives/id"
	"sync"
	"testing"
	"time"
)

type MockListener struct {
	NumHeard        int
	IsFallback      bool
	LastMessage     []byte
	LastMessageType int32
	mux             sync.Mutex
}

type Message struct {
	Contents    []byte
	Sender      *id.ID
	MessageType int32
}

func (m *Message) GetSender() *id.ID {
	return m.Sender
}

func (m *Message) GetMessageType() int32 {
	return m.MessageType
}

func (ml *MockListener) Hear(item Item, isHeardElsewhere bool, i ...interface{}) {
	ml.mux.Lock()
	defer ml.mux.Unlock()

	msg := item.(*Message)

	if !isHeardElsewhere || !ml.IsFallback {
		ml.NumHeard++
		ml.LastMessage = msg.Contents
		ml.LastMessageType = msg.GetMessageType()
	}

	if len(i) > 0 {
		hearChan := i[0].(chan struct{})
		hearChan <- struct{}{}
	}
}

var specificUser *id.ID
var specificMessageType int32 = 5
var delay = 10 * time.Millisecond

func OneListenerSetup(t *testing.T) (*Switchboard, *MockListener) {
	var listeners *Switchboard
	listeners = NewSwitchboard()
	specificUser = id.NewIdFromUInts([4]uint64{0, 0, 0, 5}, id.User, t)
	// add one listener to the map
	fullyMatchedListener := &MockListener{}
	listeners.Register(specificUser, specificMessageType,
		fullyMatchedListener)
	return listeners, fullyMatchedListener
}

func TestListenerMap_SpeakOne(t *testing.T) {
	// set up
	listeners, fullyMatchedListener := OneListenerSetup(t)

	// speak
	listeners.Speak(&Message{
		Contents:    []byte("hmmmm"),
		Sender:      specificUser,
		MessageType: specificMessageType,
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
	listeners, fullyMatchedListener := OneListenerSetup(t)

	// speak
	for i := 0; i < 20; i++ {
		go listeners.Speak(&Message{
			Contents:    make([]byte, 0),
			Sender:      specificUser,
			MessageType: specificMessageType,
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
	listeners, fullyMatchedListener := OneListenerSetup(t)
	nonzeroUser := id.NewIdFromUInts([4]uint64{0, 0, 0, 786}, id.User, t)

	// speak
	listeners.Speak(&Message{
		MessageType: specificMessageType,
		Contents:    make([]byte, 0),
		Sender:      nonzeroUser,
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
	listeners, fullyMatchedListener := OneListenerSetup(t)

	// speak
	listeners.Speak(&Message{
		MessageType: specificMessageType + 1,
		Contents:    make([]byte, 0),
		Sender:      specificUser,
	})

	// determine whether the listener heard the message
	time.Sleep(delay)
	expected := 0
	if fullyMatchedListener.NumHeard != expected {
		t.Errorf("The listener heard %v messages instead of %v",
			fullyMatchedListener.NumHeard, expected)
	}
}

var zeroUser = id.ZeroUser
var zeroMessageType int32

func WildcardListenerSetup() (*Switchboard, *MockListener) {
	var listeners *Switchboard
	listeners = NewSwitchboard()
	// add one listener to the map
	wildcardListener := &MockListener{}
	listeners.Register(&zeroUser, zeroMessageType, wildcardListener)
	return listeners, wildcardListener
}

func TestListenerMap_SpeakWildcard(t *testing.T) {
	// set up
	listeners, wildcardListener := WildcardListenerSetup()
	specificUser = id.NewIdFromUInts([4]uint64{0, 0, 0, 5}, id.User, t)

	// speak
	listeners.Speak(&Message{
		Contents:    make([]byte, 0),
		Sender:      specificUser,
		MessageType: specificMessageType + 1,
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
	specificUser = id.NewIdFromUInts([4]uint64{0, 0, 0, 5}, id.User, t)

	individualListeners := make([]*MockListener, 0)

	// one user, many types
	for messageType := int32(1); messageType <= int32(20); messageType++ {
		newListener := MockListener{}
		listeners.Register(specificUser, messageType,
			&newListener)
		individualListeners = append(individualListeners, &newListener)
	}
	// wildcard listener for the user
	userListener := &MockListener{}
	listeners.Register(specificUser, zeroMessageType, userListener)
	// wildcard listener for all messages
	wildcardListener := &MockListener{}
	listeners.Register(&zeroUser, zeroMessageType, wildcardListener)

	// send to all types for our user
	for messageType := int32(1); messageType <= int32(20); messageType++ {
		go listeners.Speak(&Message{
			MessageType: messageType,
			Contents:    make([]byte, 0),
			Sender:      specificUser,
		})
	}
	// send to all types for a different user
	otherUser := id.NewIdFromUInt(98, id.User, t)
	for messageType := int32(1); messageType <= int32(20); messageType++ {
		go listeners.Speak(&Message{
			MessageType: messageType,
			Contents:    make([]byte, 0),
			Sender:      otherUser,
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
	specificUser = id.NewIdFromUInts([4]uint64{0, 0, 0, 5}, id.User, t)
	// add one normal and one fallback listener to the map
	fallbackListener := &MockListener{}
	fallbackListener.IsFallback = true
	listeners.Register(&zeroUser, zeroMessageType, fallbackListener)
	specificListener := &MockListener{}
	listeners.Register(specificUser, specificMessageType,
		specificListener)

	// send exactly one message to each of them
	listeners.Speak(&Message{
		MessageType: specificMessageType,
		Contents:    make([]byte, 0),
		Sender:      specificUser,
	})
	listeners.Speak(&Message{
		MessageType: specificMessageType + 1,
		Contents:    make([]byte, 0),
		Sender:      specificUser,
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
	listeners, listener := OneListenerSetup(t)
	expected := []byte{0x01, 0x02, 0x03, 0x04}
	listeners.Speak(&Message{
		MessageType: specificMessageType,
		Contents:    expected,
		Sender:      specificUser,
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
	specificUser = id.NewIdFromUInts([4]uint64{0, 0, 0, 5}, id.User, t)
	listenerID := listeners.Register(specificUser, specificMessageType,
		&MockListener{})
	listeners.Unregister(listenerID)
	if len(listeners.listeners[*specificUser][specificMessageType]) != 0 {
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
	hearChan := make(chan struct{}, 5)
	listeners.Register(&id.ZeroUser, 3, l, hearChan)
	// Should match
	listeners.Speak(&Message{
		Contents:    []byte("Test 0"),
		Sender:      id.NewIdFromUInt(8, id.User, t),
		MessageType: 3,
	})
	tmr := time.NewTimer(time.Second)

	select {
	case <-hearChan:
		if l.NumHeard != 1 {
			t.Error("ListenerFunc heard but didn't Listener")
		}
	case <-tmr.C:
		t.Error("ListenerFunc did not hear")
	}

	l.NumHeard = 0
	// Should not match
	listeners.Speak(&Message{
		Contents:    []byte("Test 2"),
		Sender:      id.NewIdFromUInt(8, id.User, t),
		MessageType: 0,
	})

	select {
	case <-hearChan:
		t.Error("ListenerFunc heard but should not have")
	case <-tmr.C:
		if l.NumHeard != 0 {
			t.Error("ListenerFunc should not have heard")
		}
	}
}*/
