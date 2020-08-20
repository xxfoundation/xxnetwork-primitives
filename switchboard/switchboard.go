package switchboard

import (
	"errors"
	"github.com/golang-collections/collections/set"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/id"
	"sync"
)

type Switchboard struct {
	id          *byId
	messageType *byType

	mux sync.RWMutex
}

// New generates and returns a new switchboard object.
func New() *Switchboard {
	return &Switchboard{
		id:          newById(),
		messageType: newByType(),
	}
}

// Registers a new listener. Returns the ID of the new listener.
// Keep this around if you want to be able to delete the listener later.
//
// name is used for debug printing and not checked for uniqueness
//
// user: 0 for all, or any user ID to listen for messages from a particular
// user. 0 can be id.ZeroUser or id.ZeroID
// messageType: 0 for all, or any message type to listen for messages of that
// type. 0 can be switchboard.AnyType
// newListener: something implementing the Listener interface. Do not
// pass nil to this.
//
// If a message matches multiple listeners, all of them will hear the message.
func (sw *Switchboard) RegisterListener(user *id.ID, messageType int32,
	newListener Listener) (ListenerID, error) {

	if user == nil {
		return ListenerID{}, errors.New("cannot register listener to nil user")
	}

	if newListener == nil {
		return ListenerID{}, errors.New("cannot register nil listener")
	}

	sw.mux.Lock()

	sw.id.Add(user, newListener)
	sw.messageType.Add(messageType, newListener)

	sw.mux.Unlock()

	return ListenerID{
		userID:      user,
		messageType: messageType,
		listener:    newListener,
	}, nil
}

// Registers a new listener built around the passed function.
// Returns the ID of the new listener.
// Keep this around if you want to be able to delete the listener later.
//
// name is used for debug printing and not checked for uniqueness
//
// user: 0 for all, or any user ID to listen for messages from a particular
// user. 0 can be id.ZeroUser or id.ZeroID
// messageType: 0 for all, or any message type to listen for messages of that
// type. 0 can be switchboard.AnyType
// newListener: a function implementing the ListenerFunc function type.
// Do not pass nil to this.
//
// If a message matches multiple listeners, all of them will hear the message.
func (sw *Switchboard) RegisterFunc(name string, user *id.ID, messageType int32,
	newListener ListenerFunc) (ListenerID, error) {
	if newListener == nil {
		return ListenerID{},
			errors.New("cannot register listener with nil func")
	}
	fl := newFuncListener(newListener, name)
	return sw.RegisterListener(user, messageType, fl)
}

// Registers a new listener built around the passed channel.
// Returns the ID of the new listener.
// Keep this around if you want to be able to delete the listener later.
//
// name is used for debug printing and not checked for uniqueness
//
// user: 0 for all, or any user ID to listen for messages from a particular
// user. 0 can be id.ZeroUser or id.ZeroID
// messageType: 0 for all, or any message type to listen for messages of that
// type. 0 can be switchboard.AnyType
// newListener: an item channel.
// Do not pass nil to this.
//
// If a message matches multiple listeners, all of them will hear the message.
func (sw *Switchboard) RegisterChannel(name string, user *id.ID,
	messageType int32, newListener chan Item) (ListenerID, error) {
	if newListener == nil {
		return ListenerID{},
			errors.New("cannot register listener with nil channel")
	}
	cl := newChanListener(newListener, name)
	return sw.RegisterListener(user, messageType, cl)
}

// Speak broadcasts a message to the appropriate listeners.
func (sw *Switchboard) Speak(item Item) {
	sw.mux.RLock()
	defer sw.mux.RUnlock()

	// Matching listeners: include those that match all criteria perfectly, as
	// well as those that do not care about certain criteria
	matches := sw.matchListeners(item)

	if matches.Len() > 0 {
		matches.Do(func(i interface{}) {
			r := i.(Listener)
			go r.Hear(item)
		})
	} else {
		jww.ERROR.Printf(
			"Message of type %v from user %q didn't match any listeners in"+
				" the map", item.GetMessageType(), item.GetSender())
	}
}

// Unregister removes the listener with the specified ID so it will no longer
// get called
func (sw *Switchboard) Unregister(listenerID ListenerID) {
	sw.mux.Lock()

	idSet := sw.id.Get(listenerID.userID)
	idSet.Remove(listenerID.listener)

	typeSet := sw.messageType.Get(listenerID.messageType)
	typeSet.Remove(listenerID.listener)

	sw.mux.Unlock()
}

func (sw *Switchboard) matchListeners(item Item) *set.Set {
	idSet := sw.id.Get(item.GetSender())
	typeSet := sw.messageType.Get(item.GetMessageType())
	return idSet.Intersection(typeSet)
}
