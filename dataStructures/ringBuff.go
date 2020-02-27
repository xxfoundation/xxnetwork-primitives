////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package dataStructures

/*
 * The RingBuffer data structure is used to store information on rounds and updates
 * It functions like a typical Circluar buffer, with some slight modifications
 * First, it is made generic by using interface{} instead of a defined type
 * Second, it requires an id function to be passed in which gets an ID from whatever the underlying object is
 * Finally, it allows for manipulation of data using both normal indeces and ID values as counters
 */

import (
	"github.com/pkg/errors"
	"math"
	"sync"
)

// These function types are passed into ringbuff, allowing us to make it semi-generic
type idFunc func(interface{}) int
type compFunc func(interface{}, interface{}) bool

// A circular buffer with the ability to use IDs as position and locks built in
type RingBuff struct {
	buff              []interface{}
	count, head, tail int
	id                idFunc
	lock              sync.RWMutex
}

// Initialize a new ring buffer with length n
func NewRingBuff(n int, id idFunc) *RingBuff {
	rb := &RingBuff{
		buff:  make([]interface{}, n),
		count: n,
		head:  -1,
		tail:  0,
		id:    id,
	}
	return rb
}

// Push a round to the buffer
func (rb *RingBuff) Push(val interface{}) {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	rb.push(val)
}

// push a round to a relative index in the buffer
func (rb *RingBuff) UpsertById(val interface{}, comp compFunc) error {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	// Make sure the id isn't too old
	newId := rb.id(val)
	if rb.head != -1 && rb.id(rb.buff[rb.head]) > newId {
		return errors.Errorf("Did not upsert value %+v; id is older than first tracked", val)
	}

	// Get most recent ID so we can figure out where to put this
	mostRecentIndex := (rb.tail + rb.count - 1) % rb.count
	lastId := rb.id(rb.buff[mostRecentIndex])
	if lastId+1 == newId {
		// last id is the previous one; we can just push
		rb.push(val)
	} else if (lastId + 1) < newId {
		// there are id's between the last and the current; increment using dummy entries
		for i := lastId + 1; i < newId; i++ {
			rb.push(nil)
		}
		rb.push(val)
	} else if lastId+1 > newId {
		// this is an old ID, check the comp function and insert if true
		i := rb.getIndex(newId - (lastId + 1))
		if comp(rb.buff[i], val) {
			rb.buff[i] = val
		} else {
			return errors.Errorf("Did not upsert value %+v; comp function returned false", val)
		}
	}
	return nil
}

// Retreive the most recent entry
func (rb *RingBuff) Get() interface{} {
	rb.lock.RLock()
	defer rb.lock.RUnlock()

	mostRecentIndex := (rb.tail + rb.count - 1) % rb.count
	return rb.buff[mostRecentIndex]
}

// Retrieve an entry with the given ID
func (rb *RingBuff) GetById(id int) (interface{}, error) {
	rb.lock.RLock()
	defer rb.lock.RUnlock()

	// Check it's not before our first known id
	firstId := rb.id(rb.buff[rb.head])
	if id < firstId {
		return nil, errors.Errorf("requested ID %d is lower than oldest id %d", id, firstId)
	}

	// Check it's not after our last known id
	lastId := rb.id(rb.Get())
	if id > lastId {
		return nil, errors.Errorf("requested id %d is higher than most recent id %d", id, lastId)
	}

	index := rb.getIndex(id - firstId)
	return rb.buff[index], nil
}

// Retrieve an entry at the given index
func (rb *RingBuff) GetByIndex(i int) (interface{}, error) {
	rb.lock.RLock()
	defer rb.lock.RUnlock()

	if math.Abs(float64(i)) > float64(rb.Len()) { // FIXME: this shouldn't be float but we don't have a package where it's not float
		return nil, errors.Errorf("Could not get item at index %d: index out of bounds", i)
	}
	return rb.buff[rb.getIndex(i)], nil
}

// Return length of the structure
func (rb *RingBuff) Len() int {
	rb.lock.RLock()
	defer rb.lock.RUnlock()

	return rb.count
}

// next is a helper function for ringbuff
// it handles incrementing the head & tail markers
func (rb *RingBuff) next() {
	rb.tail = (rb.tail + 1) % rb.count
	if rb.tail-1 == rb.head {
		rb.head = (rb.head + 1) % rb.count
	}
	if rb.head == -1 {
		rb.head = 0
	}
}

// getIndex is a helper function for ringbuff
// it returns an index relative to the head/tail position of the buffer
func (rb *RingBuff) getIndex(i int) int {
	var index int
	if i < 0 {
		index = (rb.tail + rb.count + i) % rb.count
	} else {
		index = (rb.head + i) % rb.count
	}
	return index
}

// Push a round to the buffer
func (rb *RingBuff) push(val interface{}) {
	rb.buff[rb.tail] = val
	rb.next()
}
