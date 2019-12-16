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
	"sync"
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

//We use this structure to map these value to a listenerRecord array in the sync.Map
type listenerMapId struct {
	userId      id.User
	messageType int32
}

type Switchboard struct {
	// By matching with the keys for each level of the map,
	// you can find the listeners that meet each criterion
	listeners sync.Map
	//listenerIds maps listener ids to listenerMapIds making the reverse search in unregister o(k).
	listenerIds sync.Map
	lastID      int
}

var Listeners = NewSwitchboard()

func NewSwitchboard() *Switchboard {
	return &Switchboard{
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
	messageType int32, newListener Listener) string {

	mapId := listenerMapId{*user, messageType}

	lm.lastID++

	newListenerRecord := &listenerRecord{
		l:  newListener,
		id: strconv.Itoa(lm.lastID),
	}

	listenerArray, ok := getListenerRecords(mapId, lm.listeners)
	newListenerRecordSlice := []*listenerRecord{}
	if ok {
		//sync map returns an interface, so give it a type then append and save
		newListenerRecordSlice = append(listenerArray, newListenerRecord)
	} else {
		newListenerRecordSlice = append(newListenerRecordSlice, newListenerRecord)
	}

	//store listener id into the map that links them back to map ids for the o(1) search in unregister
	lm.listenerIds.Store(newListenerRecord.id, mapId)
	lm.listeners.Store(mapId, newListenerRecordSlice)

	return newListenerRecord.id
}

func (lm *Switchboard) Unregister(listenerID string) {
	// This method uses a map of listenerIds to listenerMapId objects so we know where the listener object is making the
	// search o(k).
	unregisterId_i, ok := lm.listenerIds.Load(listenerID)

	if ok {
		unregisterMapId := unregisterId_i.(listenerMapId)
		listeners, ok := getListenerRecords(unregisterMapId, lm.listeners)
		if ok {
			for i := range listeners {
				if listenerID == listeners[i].id {
					//In deleting here is it important to maintain order? quicker solution if not
					newListeners := deleteElem(i, listeners)
					lm.listenerIds.Delete(listenerID)
					lm.listeners.Store(unregisterMapId, newListeners)
					return
				}
			}

		} else {
			// Could not be found therefore doesnt exist
			return
		}
	}
}

func deleteElem(loc int, records []*listenerRecord) []*listenerRecord {
	//removes listener and keeps order
	copy(records[loc:], records[loc+1:])        // Shift a[i+1:] left one index.
	records[len(records)-1] = &listenerRecord{} // Erase last element (write zero value).
	records = records[:len(records)-1]          // Truncate slice.

	return records
}

func (lm *Switchboard) matchListeners(item Item) []*listenerRecord {
	matches := make([]*listenerRecord, 0)

	// 8 cases total, for matching both specific and general listeners
	// This seems inefficient
	matches = getMatches(matches, *item.GetSender(), item.GetMessageType(), lm)
	matches = getMatches(matches, *id.ZeroID, item.GetMessageType(), lm)
	matches = getMatches(matches, *item.GetSender(), 0, lm)
	matches = getMatches(matches, *id.ZeroID, 0, lm)
	matches = getMatches(matches, *item.GetSender(), 0, lm)
	matches = getMatches(matches, *id.ZeroID, 0, lm)
	// Match all, but with generic outer type
	matches = getMatches(matches, *item.GetSender(), item.GetMessageType(), lm)
	matches = getMatches(matches, *id.ZeroID, item.GetMessageType(), lm)

	return matches
}

//loops through the listener getting all matches and returning the appended matches object
func getMatches(matches []*listenerRecord, user id.User, messageType int32, lm *Switchboard) []*listenerRecord {

	mapId := listenerMapId{user, messageType}
	listeners, ok := getListenerRecords(mapId, lm.listeners)
	if ok {
		for _, listener := range listeners {
			matches = appendIfUnique(matches, listener)
		}
	}
	return matches
}

func appendIfUnique(matches []*listenerRecord, newListener *listenerRecord) []*listenerRecord {
	// Search for the listener ID
	found := false
	for _, l := range matches {
		found = found || (l.id == newListener.id)
	}
	if !found {
		// If we didn't find it, it's OK to append it to the slice
		return append(matches, newListener)
	} else {
		// We already matched this listener, and shouldn't append it
		return matches
	}
}

// Broadcast a message to the appropriate listeners
func (lm *Switchboard) Speak(item Item) {
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
			"Message of type %v from user %q didn't match any listeners in"+
				" the map", item.GetMessageType(), item.GetSender())
		// dump representation of the map
		printListenersMap(lm.listeners)
	}
}

// Loops through the sync map to print out a representation of the map in the error logs. that can be used for debugging.
func printListenersMap(lm sync.Map) {
	lm.Range(func(key interface{}, value interface{}) bool {
		mapId := key.(listenerMapId)
		listeners := value.([]*listenerRecord)
		for i := range listeners {
			jww.ERROR.Printf("Listener %v: %v, user %v, "+
				" type %v, ", i, listeners[i].id, mapId.userId, mapId.messageType)
			return true
		}
		return false
	})
}

// This function loads the listenerRecord slice from our sync map, and typing the interface
func getListenerRecords(mapId listenerMapId, listenerMap sync.Map) ([]*listenerRecord, bool) {
	listenerArray_i, ok := listenerMap.Load(mapId)
	if ok {
		return listenerArray_i.([]*listenerRecord), ok
	}
	return nil, ok
}
