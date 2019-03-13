////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package switchboard

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/id"
	"reflect"
	"strconv"
	"sync"
)

type Item interface {
	// To reviewer: Is this the correct name for this method? It's always the
	// sender ID in the client, but that might not be the case on the nodes
	GetSender() *id.User
	GetOuterType() int32
	GetInnerType() int32
}

// This is an interface so you can receive callbacks through the Gomobile boundary
type Listener interface {
	Hear(item Item, isHeardElsewhere bool)
}

type listenerRecord struct {
	l  Listener
	id string
}

type Switchboard struct {
	// By matching with the keys for each level of the map,
	// you can find the listeners that meet each criterion
	listeners map[id.User]map[int32]map[int32][]*listenerRecord
	lastID    int
	mux       sync.RWMutex
}

var Listeners = NewSwitchboard()

func NewSwitchboard() *Switchboard {
	return &Switchboard{
		listeners: make(map[id.User]map[int32]map[int32][]*listenerRecord),
		lastID: 0,
	}
}

// Add a new listener to the map
// Returns ID of the new listener. Keep this around if you want to be able to
// delete the listener later.
//
// user: 0 for all,
// or any user ID to listen for messages from a particular user.
// messageType: 0 for all, or any message type to listen for messages of that
// type.
// newListener: something implementing the Listener callback interface.
// Don't pass nil to this.
//
// If a message matches multiple listeners, all of them will hear the message.
func (lm *Switchboard) Register(user *id.User,
	outerType int32, innerType int32,
	newListener Listener) string {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	lm.lastID++
	if lm.listeners[*user] == nil {
		lm.listeners[*user] = make(map[int32]map[int32][]*listenerRecord)
	}

	if lm.listeners[*user][outerType] == nil {
		lm.listeners[*user][outerType] = make(map[int32][]*listenerRecord)
	}

	newListenerRecord := &listenerRecord{
		l:  newListener,
		id: strconv.Itoa(lm.lastID),
	}
	lm.listeners[*user][outerType][innerType] = append(
		lm.listeners[*user][outerType][innerType],
		newListenerRecord)

	return newListenerRecord.id
}

func (lm *Switchboard) Unregister(listenerID string) {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	// Iterate over all listeners in the map.
	for u, perUser := range lm.listeners {
		for outerType, perOuterType := range perUser {
			for innerType, perInnerType := range perOuterType {
				for i, listener := range perInnerType {
					if listener.id == listenerID {
						// this matches, so remove listener from data structure
						lm.listeners[u][outerType][innerType] = append(
							perInnerType[:i], perInnerType[i+1:]...)
						// since the id is unique per listener,
						// we can terminate the loop early.
						return
					}
				}
			}
		}
	}
}

func (lm *Switchboard) matchListeners(item Item) []*listenerRecord {

	matches := make([]*listenerRecord, 0)

	// 8 cases total, for matching both specific and general listeners
	// This seems inefficient
	for _, listener := range lm.listeners[*item.GetSender()][item.
		GetOuterType()][item.GetInnerType()] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][item.
		GetOuterType()][item.GetInnerType()] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*item.GetSender()][item.GetOuterType()][0] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][item.GetOuterType()][0] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*item.GetSender()][0][0] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][0][0] {
		matches = append(matches, listener)
	}
	// Match all, but with generic outer type
	for _, listener := range lm.listeners[*item.GetSender()][0][item.GetInnerType()] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][0][item.GetInnerType()] {
		matches = append(matches, listener)
	}

	return matches
}

// Broadcast a message to the appropriate listeners
func (lm *Switchboard) Speak(item Item) {
	lm.mux.RLock()
	defer lm.mux.RUnlock()

	// Matching listeners include those that match all criteria perfectly,
	// as well as those that don't care about certain criteria.
	matches := lm.matchListeners(item)

	if len(matches) > 0 {
		// notify all normal listeners
		for _, listener := range matches {
			jww.INFO.Printf("Hearing on listener %v of type %v",
				listener.id, reflect.TypeOf(listener.l))
			// If you want to be able to hear things on the switchboard on
			// multiple goroutines, you should call Speak() on the switchboard
			// from multiple goroutines
			listener.l.Hear(item, len(matches) > 1)
		}
	} else {
		jww.ERROR.Printf(
			"Message of type %v, %v from user %q didn't match any listeners in"+
				" the map", item.GetOuterType(), item.GetInnerType(),
			item.GetSender())
		// dump representation of the map
		for u, perUser := range lm.listeners {
			for outerType, perOuterType := range perUser {
				for innerType, perInnerType := range perOuterType {
					for i, listener := range perInnerType {

						jww.ERROR.Printf("Listener %v: %v, user %v, "+
							"outertype %v, type %v, ",
							i, listener.id, u, outerType, innerType)
					}
				}
			}
		}
	}
}
