////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package switchboard

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/cmixproto"
	"gitlab.com/elixxir/primitives/id"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
)

type Item interface {
	// To reviewer: Is this the correct name for this method? It's always the
	// sender ID in the client, but that might not be the case on the nodes
	GetSender() *id.User
	GetOuterType() cmixproto.OuterType
	GetInnerType() cmixproto.InnerType
}

// This is an interface so you can receive callbacks through the Gomobile boundary
type Listener interface {
	Hear(item Item, isHeardElsewhere bool)
}

// Implementation of Listener interface using a channel
type listeningQueue struct {
	bufSize    int32
	buf        []Item
	writeIndex int32
	queue      chan []Item
}

// Set up a listening queue and add it to the switchboard
func (s *Switchboard) ListenChannel(bufSize int32, channelBufferSize int,
	outerType cmixproto.OuterType, innerType cmixproto.InnerType,
	sender *id.User) (id string, messageQueue chan []Item) {
	l := listeningQueue{
		bufSize:    bufSize,
		buf:        make([]Item, bufSize),
		writeIndex: 0,
		queue:      make(chan []Item, channelBufferSize),
	}
	id = s.Register(sender, outerType, innerType, l)
	return id, l.queue
}

// Empty remaining elements in the buffer
func (l *listeningQueue) Flush() {
	result := l.buf
	l.buf = make([]Item, 0, l.bufSize)
	atomic.StoreInt32(&l.writeIndex, 0)
	l.queue <- result
}

// Returns the number of items in the listening buffer.
// May be helpful when determining when to do the final flush.
func (l *listeningQueue) Len() int32 {
	return l.writeIndex
}

// TODO What happens if you use pointer receiver? Should test whether it still works correctly.
// Multiple threads can write to this buffer simultaneously through the
// switchboard using this method
func (l listeningQueue) Hear(item Item, isHeardElsewhere bool) {
	writeIndex := atomic.LoadInt32(&l.writeIndex)
	for {
		// Only try to CAS if the write index is in the correct range
		// If it's outside the correct range, it should return to the correct
		// range soon, because another goroutine is about to flush the buffer
		if writeIndex < l.bufSize {
			// If the CAS succeeds, it's safe to write to the buffer at writeIndex
			if atomic.CompareAndSwapInt32(&l.writeIndex, writeIndex, writeIndex+1) {
				break
			}
		}
		writeIndex = atomic.LoadInt32(&l.writeIndex)
	}
	l.buf[writeIndex] = item
	if writeIndex >= l.bufSize {
		l.Flush()
	}
}

type listenerRecord struct {
	l  Listener
	id string
}

type Switchboard struct {
	// By matching with the keys for each level of the map,
	// you can find the listeners that meet each criterion
	listeners map[id.User]map[cmixproto.OuterType]map[cmixproto.InnerType][]*listenerRecord
	lastID    int
	mux       sync.RWMutex
}

var Listeners = NewSwitchboard()

func NewSwitchboard() *Switchboard {
	return &Switchboard{
		listeners: make(map[id.User]map[cmixproto.OuterType]map[cmixproto.
			InnerType][]*listenerRecord),
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
	outerType cmixproto.OuterType, innerType cmixproto.InnerType,
	newListener Listener) string {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	lm.lastID++
	if lm.listeners[*user] == nil {
		lm.listeners[*user] = make(map[cmixproto.OuterType]map[cmixproto.InnerType][]*listenerRecord)
	}

	if lm.listeners[*user][outerType] == nil {
		lm.listeners[*user][outerType] = make(map[cmixproto.InnerType][]*listenerRecord)
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
	for _, listener := range lm.listeners[*item.GetSender()][item.
		GetOuterType()][item.GetInnerType()] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][item.
		GetOuterType()][item.GetInnerType()] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*item.GetSender()][cmixproto.
		OuterType_NONE][cmixproto.InnerType_NO_TYPE] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][cmixproto.
		OuterType_NONE][cmixproto.InnerType_NO_TYPE] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*item.GetSender()][item.
		GetOuterType()][cmixproto.InnerType_NO_TYPE] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][item.
		GetOuterType()][cmixproto.InnerType_NO_TYPE] {
		matches = append(matches, listener)
	}
	// Match all, but with generic outer type
	for _, listener := range lm.listeners[*item.GetSender()][cmixproto.
		OuterType_NONE][item.GetInnerType()] {
		matches = append(matches, listener)
	}
	for _, listener := range lm.listeners[*id.ZeroID][cmixproto.
		OuterType_NONE][item.GetInnerType()] {
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
				" the map", item.GetOuterType().String(), item.GetInnerType().String(),
				item.GetSender())
		// dump representation of the map
		for u, perUser := range lm.listeners {
			for outerType, perOuterType := range perUser {
				for innerType, perInnerType := range perOuterType {
					for i, listener := range perInnerType {

						jww.ERROR.Printf("Listener %v: %v, user %v, " +
							"outertype %v, type %v, ",
							i, listener.id, u, outerType.String(),
							innerType.String())
					}
				}
			}
		}
	}
}
