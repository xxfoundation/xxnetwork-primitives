////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package switchboard

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/format"
	"reflect"
	"strconv"
	"sync"
)

type Item interface {
	// To reviewer: Is this the correct name for this method? It's always the
	// sender ID in the client, but that might not be the case on the nodes
	GetSender() *id.User
	GetCryptoType() format.CryptoType
	GetMessageType() int32
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
	listeners map[id.User]map[format.CryptoType]map[int32][]*listenerRecord
	lastID    int
	mux       sync.RWMutex
}

var Listeners = NewSwitchboard()

func NewSwitchboard() *Switchboard {
	return &Switchboard{
		listeners: make(map[id.User]map[format.CryptoType]map[int32][]*listenerRecord),
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
	cryptoType format.CryptoType, messageType int32,
	newListener Listener) string {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	lm.lastID++
	if lm.listeners[*user] == nil {
		lm.listeners[*user] = make(map[format.CryptoType]map[int32][]*listenerRecord)
	}

	if lm.listeners[*user][cryptoType] == nil {
		lm.listeners[*user][cryptoType] = make(map[int32][]*listenerRecord)
	}

	newListenerRecord := &listenerRecord{
		l:  newListener,
		id: strconv.Itoa(lm.lastID),
	}
	lm.listeners[*user][cryptoType][messageType] = append(
		lm.listeners[*user][cryptoType][messageType],
		newListenerRecord)

	return newListenerRecord.id
}

func (lm *Switchboard) Unregister(listenerID string) {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	// Iterate over all listeners in the map.
	for u, perUser := range lm.listeners {
		for outerType, perCryptoType := range perUser {
			for innerType, perMessageType := range perCryptoType {
				for i, listener := range perMessageType {
					if listener.id == listenerID {
						// this matches, so remove listener from data structure
						lm.listeners[u][outerType][innerType] = append(
							perMessageType[:i], perMessageType[i+1:]...)
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
		GetCryptoType()][item.GetMessageType()] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][item.
		GetCryptoType()][item.GetMessageType()] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*item.GetSender()][item.GetCryptoType()][0] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][item.GetCryptoType()][0] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*item.GetSender()][format.None][0] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][format.None][0] {
		matches = append(matches, listener)
	}
	// Match all, but with generic outer type
	for _, listener := range lm.listeners[*item.GetSender()][format.None][item.GetMessageType()] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][format.None][item.GetMessageType()] {
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
				" the map", item.GetCryptoType(), item.GetMessageType(),
			item.GetSender())
		// dump representation of the map
		for u, perUser := range lm.listeners {
			for outerType, perCryptoType := range perUser {
				for innerType, perMessageType := range perCryptoType {
					for i, listener := range perMessageType {

						jww.ERROR.Printf("Listener %v: %v, user %v, "+
							"outertype %v, type %v, ",
							i, listener.id, u, outerType, innerType)
					}
				}
			}
		}
	}
}
