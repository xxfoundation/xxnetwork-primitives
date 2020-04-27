////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package ring

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
type Buff struct {
	buff            []interface{}
	count, old, new int
	id              idFunc
	sync.RWMutex
}

// Initialize a new ring buffer with length n
func NewBuff(n int, id idFunc) *Buff {
	rb := &Buff{
		buff:  make([]interface{}, n),
		count: n,
		old:   -1,
		new:   0,
		id:    id,
	}
	return rb
}

// Get the ID of the newest item in the buffer
func (rb *Buff) GetNewestId() int {
	rb.RLock()
	defer rb.RUnlock()

	return rb.getNewestId()
}

// Get the IDof the oldest item in the buffer
func (rb *Buff) GetOldestId() int {
	rb.RLock()
	defer rb.RUnlock()

	return rb.getOldestId()
}

// Get the index of the oldest item in the buffer, not including nil values inserted by UpsertById
func (rb *Buff) GetOldestIndex() int {
	rb.RLock()
	defer rb.RUnlock()

	return rb.getOldestIndex()
}

// Push a round to the buffer
func (rb *Buff) Push(val interface{}) {
	rb.Lock()
	defer rb.Unlock()

	rb.push(val)
}

// push a round to a relative index in the buffer
func (rb *Buff) UpsertById(val interface{}, comp compFunc) error {
	rb.Lock()
	defer rb.Unlock()

	// Make sure the id isn't too old
	newId := rb.id(val)
	if rb.old != -1 && rb.id(rb.buff[rb.old]) > newId {
		return errors.Errorf("Did not upsert value %+v; id is older than first tracked", val)
	}

	// Get most recent ID so we can figure out where to put this
	mostRecentIndex := (rb.new + rb.count - 1) % rb.count
	mostRecentId := rb.id(rb.buff[mostRecentIndex])
	if mostRecentId+1 == newId {
		// last id is the previous one; we can just push
		rb.push(val)
	} else if (mostRecentId + 1) < newId {
		// there are id's between the last and the current; increment using dummy entries
		for i := mostRecentId + 1; i < newId; i++ {
			rb.push(nil)
		}
		rb.push(val)
	} else if mostRecentId+1 > newId {
		// this is an old ID, check the comp function and insert if true
		i := rb.getIndex(newId - (mostRecentId + 1))
		if comp(rb.buff[i], val) {
			rb.buff[i] = val
		} else {
			return errors.Errorf("Did not upsert value %+v; comp function returned false", val)
		}
	}
	return nil
}

// Retreive the most recent entry
func (rb *Buff) Get() interface{} {
	rb.RLock()
	defer rb.RUnlock()

	mostRecentIndex := (rb.new + rb.count - 1) % rb.count
	return rb.buff[mostRecentIndex]
}

// Retrieve an entry with the given ID
func (rb *Buff) GetById(id int) (interface{}, error) {
	rb.RLock()
	defer rb.RUnlock()

	// Check it's not before our first known id
	firstId := rb.id(rb.buff[rb.getOldestIndex()])
	if id < firstId {
		return nil, errors.Errorf("requested ID %d is lower than oldest id %d", id, firstId)
	}

	// Check it's not after our last known id
	lastId := rb.id(rb.Get())
	if id > lastId {
		return nil, errors.Errorf("requested id %d is higher than most recent id %d", id, lastId)
	}

	index := rb.getIndex(id - lastId - 1)
	return rb.buff[index], nil
}

// Retrieve an entry at the given index
func (rb *Buff) GetByIndex(i int) (interface{}, error) {
	rb.RLock()
	defer rb.RUnlock()

	if math.Abs(float64(i)) > float64(rb.Len()) { // FIXME: this shouldn't be float but we don't have a package where it's not float
		return nil, errors.Errorf("Could not get item at index %d: index out of bounds", i)
	}
	return rb.buff[rb.getIndex(i)], nil
}

//retrieve all entries newer than the passed one
func (rb *Buff) GetNewerById(id int) ([]interface{}, error) {
	rb.RLock()
	defer rb.RUnlock()

	new := rb.getNewestId()
	old := rb.getOldestId()

	if id < old {
		id = old - 1
	}

	if id > new {
		return nil, errors.Errorf("requested ID %d is higher than the"+
			" newest id %d", id, new)
	}

	numNewer := new - id
	list := make([]interface{}, numNewer)

	for i := 0; i < numNewer; i++ {
		//error is suppressed because it only occurs when out of bounds,
		//but bounds are already assured in this function
		list[i], _ = rb.GetById(id + i + 1)
	}

	return list, nil
}

// Return length of the structure
func (rb *Buff) Len() int {
	rb.RLock()
	defer rb.RUnlock()

	return rb.count
}

// next is a helper function for ringbuff
// it handles incrementing the old & new markers
func (rb *Buff) next() {
	rb.new = (rb.new + 1) % rb.count
	if rb.new-1 == rb.old {
		rb.old = (rb.old + 1) % rb.count
	}
	if rb.old == -1 {
		rb.old = 0
	}
}

// getIndex is a helper function for ringbuff
// it returns an index relative to the old/new position of the buffer
func (rb *Buff) getIndex(i int) int {
	var index int
	if i < 0 {
		index = (rb.new + rb.count + i) % rb.count
	} else {
		index = (rb.old + i) % rb.count
	}
	return index
}

// Push a round to the buffer
func (rb *Buff) push(val interface{}) {
	rb.buff[rb.new] = val
	rb.next()
}

// Get the ID of the newest item in the buffer
func (rb *Buff) getNewestId() int {
	mostRecentIndex := (rb.new + rb.count - 1) % rb.count
	if rb.buff[mostRecentIndex] == nil {
		return -1
	}
	return rb.id(rb.Get())
}

// Get the IDof the oldest item in the buffer
func (rb *Buff) getOldestId() int {
	if rb.old == -1 {
		return -1
	}
	return rb.id(rb.buff[rb.getOldestIndex()])
}

// Get the index of the oldest item in the buffer, not including nil values inserted by UpsertById
func (rb *Buff) getOldestIndex() int {
	if rb.old == -1 {
		return -1
	}
	var last = rb.old
	for ; rb.buff[last] == nil; last = (last + 1) % rb.count {
	}
	return last
}
