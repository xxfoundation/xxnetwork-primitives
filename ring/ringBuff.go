////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

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
	"sync"
)

// A circular buffer with the ability to use IDs as position and locks built in
type Buff struct {
	buff                  []interface{}
	count, oldest, newest int
	mux                   sync.RWMutex
}

// Initialize a new ring buffer with length n
func NewBuff(n int) *Buff {
	rb := &Buff{
		buff:   make([]interface{}, n),
		count:  n,
		oldest: 0,
		newest: -1,
	}
	return rb
}

// Get the ID of the newest item in the buffer
func (rb *Buff) GetNewestId() int {
	rb.mux.RLock()
	defer rb.mux.RUnlock()

	return rb.newest
}

// Get the IDof the oldest item in the buffer
func (rb *Buff) GetOldestId() int {
	rb.mux.RLock()
	defer rb.mux.RUnlock()

	return rb.oldest
}

// Push a round to the buffer
func (rb *Buff) Push(val interface{}) {
	rb.mux.Lock()
	defer rb.mux.Unlock()

	rb.push(val)
}

// push a round to a relative index in the buffer
func (rb *Buff) UpsertById(newId int, val interface{}) error {
	rb.mux.Lock()
	defer rb.mux.Unlock()

	// Make sure the id isn't too old
	if rb.oldest > newId {
		return errors.Errorf("Did not upsert value %+v; id is older than first tracked", val)
	}

	// Get most recent ID so we can figure out where to put this
	firstEmptyID := rb.newest + 1

	//fill the buffer up until the newID
	for i := firstEmptyID; i <= newId; i++ {
		rb.push(nil)
	}

	//add the data at the correct location
	index := newId % rb.count
	rb.buff[index] = val

	return nil
}

// Retreive the most recent entry
func (rb *Buff) Get() interface{} {
	rb.mux.RLock()
	defer rb.mux.RUnlock()

	mostRecentIndex := rb.newest % rb.count
	return rb.buff[mostRecentIndex]
}

// Retrieve an entry with the given ID
func (rb *Buff) GetById(id int) (interface{}, error) {
	rb.mux.RLock()
	defer rb.mux.RUnlock()

	return rb.getById(id)
}

// Retrieve an entry at the given index
func (rb *Buff) GetByIndex(i int) (interface{}, error) {
	rb.mux.RLock()
	defer rb.mux.RUnlock()

	if i < 0 || i >= rb.count {
		return nil, errors.Errorf("Could not get item at index %d: index out of bounds", i)
	}

	return rb.buff[i], nil
}

//retrieve all entries newer than the passed one
func (rb *Buff) GetNewerById(id int) ([]interface{}, error) {
	rb.mux.RLock()
	defer rb.mux.RUnlock()

	if id < rb.oldest {
		id = rb.oldest - 1
	}

	if id > rb.newest {
		return nil, errors.Errorf("requested ID %d is higher than the"+
			" newest id %d", id, rb.newest)
	}

	list := make([]interface{}, rb.newest-id)

	for i := id + 1; i <= rb.newest; i++ {
		//error is suppressed because it only occurs when out of bounds,
		//but bounds are already assured in this function
		list[i-id-1], _ = rb.getById(i)
	}

	return list, nil
}

// Return length of the structure
func (rb *Buff) Len() int {
	rb.mux.RLock()
	defer rb.mux.RUnlock()

	return rb.count
}

// next is a helper function for ringbuff
// it handles incrementing the old & new markers
func (rb *Buff) next() {
	rb.newest++
	if rb.newest >= rb.count {
		rb.oldest++
	}
}

// Push a round to the buffer
func (rb *Buff) push(val interface{}) {
	rb.next()
	rb.buff[rb.newest%rb.count] = val
}

// Retrieve an entry with the given ID for internal use without getting the read
// lock
func (rb *Buff) getById(id int) (interface{}, error) {

	// Check it's not before our first known id
	if id < rb.oldest {
		return nil, errors.Errorf("requested ID %d is lower than oldest id %d", id, rb.newest)
	}

	// Check it's not after our last known id
	if id > rb.newest {
		return nil, errors.Errorf("requested id %d is higher than most recent id %d", id, rb.oldest)
	}

	return rb.buff[id%rb.count], nil
}
