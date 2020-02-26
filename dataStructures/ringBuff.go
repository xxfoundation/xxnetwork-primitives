package dataStructures

/*
 * The dataStructures package contains data structures for use in other repos
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

// next is a helper function for ringbuff
// it handles incrementing the head & tail markers
func (rb *RingBuff) next() {
	rb.tail = (rb.tail + 1) % rb.count
	if rb.tail == rb.head {
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
	newId := rb.id(val)

	if rb.id(rb.buff[rb.head]) > newId {
		return errors.Errorf("Did not upsert value %+v; id is older than first tracked", val)
	}

	lastId := rb.id(rb.Get())
	if lastId+1 == newId {
		rb.push(val)
	} else if (lastId + 1) < newId {
		for i := lastId + 1; i < newId; i++ {
			rb.push(nil)
		}
		rb.push(val)
	} else if lastId+1 > newId {
		i := rb.getIndex(newId - (lastId + 1))
		if comp(rb.buff[i], val) {
			rb.buff[i] = val
		} else {
			return errors.Errorf("Did not upsert value %+v; comp function returned false", val)
		}
	}
	return nil
}

func (rb *RingBuff) Get() interface{} {
	rb.lock.RLock()
	defer rb.lock.RUnlock()

	mostRecentIndex := (rb.tail + rb.count - 1) % rb.count
	return rb.buff[mostRecentIndex]
}

func (rb *RingBuff) GetById(id int) (interface{}, error) {
	rb.lock.RLock()
	defer rb.lock.RUnlock()

	firstId := rb.id(rb.buff[rb.head])
	if id < firstId {
		return nil, errors.Errorf("requested ID %d is lower than oldest id %d", id, firstId)
	}

	lastId := rb.id(rb.Get())
	if id > lastId {
		return nil, errors.Errorf("requested id %d is higher than most recent id %d", id, lastId)
	}

	index := rb.getIndex(id - firstId)
	return rb.buff[index], nil
}

func (rb *RingBuff) GetByIndex(i int) (interface{}, error) {
	if math.Abs(float64(i)) > float64(rb.Len()) {
		return nil, errors.Errorf("Could not get item at index %d: index out of bounds", i)
	}
	return rb.buff[rb.getIndex(i)], nil
}

func (rb *RingBuff) Len() int {
	rb.lock.RLock()
	defer rb.lock.RUnlock()

	return rb.count
}
