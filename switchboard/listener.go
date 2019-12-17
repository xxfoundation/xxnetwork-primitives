////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package switchboard

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/id"
	"reflect"
	"strconv"
)

type Item interface {
	// To reviewer: Is this the correct name for this method? It's always the
	// sender ID in the client, but that might not be the case on the nodes
	GetSender() *id.User
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
	// listenersMap is a structure holding a syncMap
	// of our listenersMap designed to fetch and store them asynchronously
	listenersMap listenersMap
	lastID       int
}

var Listeners = NewSwitchboard()

func NewSwitchboard() *Switchboard {
	return &Switchboard{
		listenersMap: listenersMap{},
		lastID:       0,
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
// If a message matches multiple listenersMap, all of them will hear the message.
func (lm *Switchboard) Register(user *id.User,
	messageType int32, newListener Listener) string {

	lm.lastID++

	newListenerRecord := &listenerRecord{
		l:  newListener,
		id: strconv.Itoa(lm.lastID),
	}

	lm.listenersMap.StoreListener(user, messageType, newListenerRecord)

	return newListenerRecord.id
}

func (lm *Switchboard) Unregister(listenerID string) {
	lm.listenersMap.RemoveListener(listenerID)
}


func (lm *Switchboard) matchListeners(item Item) []*listenerRecord {
	matches := make([]*listenerRecord, 0)

	// 4 cases total, for matching both specific and general listenersMap
	// This seems inefficient
	matches = lm.listenersMap.GetMatches(matches, item.GetSender(), 0)
	matches = lm.listenersMap.GetMatches(matches, id.ZeroID, 0)
	// Match all, but with generic outer type
	matches = lm.listenersMap.GetMatches(matches, item.GetSender(), item.GetMessageType())
	matches = lm.listenersMap.GetMatches(matches, id.ZeroID, item.GetMessageType())

	return matches
}



// Broadcast a message to the appropriate listenersMap
func (lm *Switchboard) Speak(item Item) {
	// Matching listenersMap include those that match all criteria perfectly,
	// as well as those that don't care about certain criteria.
	matches := lm.matchListeners(item)

	if len(matches) > 0 {
		// notify all normal listenersMap
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
			"Message of type %v from user %q didn't match any listenersMap in"+
				" the map", item.GetMessageType(), item.GetSender())
		// dump representation of the map
		lm.listenersMap.String()
	}
}
