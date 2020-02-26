package dataStructures

import (
	"github.com/pkg/errors"
	"sync"
)

type idFunc func(interface{}) int
type compFunc func(interface{}, interface{}) bool

type RingBuff struct {
	buff             []interface{}
	len, first, last int
	id               idFunc
	lock             sync.Mutex
}

// next is a helper function for ringbuff
// it handles incrementing the first & last markers
func (rb *RingBuff) next() {
	rb.last = (rb.last + 1) % rb.len
	if rb.last == rb.first {
		rb.first = (rb.first + 1) % rb.len
	}
	if rb.first == -1 {
		rb.first = 0
	}
}

// getIndex is a helper function for ringbuff
// it returns an index relative to the first/last position of the buffer
func (rb *RingBuff) getIndex(i int) int {
	var index int
	if i < 0 {
		index = (rb.last + rb.len + i) % rb.len
	} else {
		index = (rb.first + i) % rb.len
	}
	return index
}

// Initialize a new ring buffer with length n
func New(n int, id idFunc) *RingBuff {
	rb := &RingBuff{
		buff:  make([]interface{}, n),
		len:   n,
		first: -1,
		last:  0,
		id:    id,
	}
	return rb
}

// Push a round to the buffer
func (rb *RingBuff) push(val interface{}) {
	rb.buff[rb.last] = val
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

	if rb.id(rb.buff[rb.first]) > newId {
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
	mostRecentIndex := (rb.last + rb.len - 1) % rb.len
	return rb.buff[mostRecentIndex]
}

func (rb *RingBuff) GetById(id int) (interface{}, error) {
	firstId := rb.id(rb.buff[rb.first])
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

func (rb *RingBuff) Len() int {
	return rb.len
}
