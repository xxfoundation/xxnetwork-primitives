package switchboard

import (
	"github.com/golang-collections/collections/set"
)

type byType struct {
	list    map[int32]*set.Set
	generic *set.Set
}

// builds a new byType structure
// registers an AnyType as generic
func newByType() *byType {
	bt := &byType{
		list:    make(map[int32]*set.Set),
		generic: set.New(),
	}

	// make the zero messages, which are defined as AnyType,
	// point to the generic
	bt.list[AnyType] = bt.generic

	return bt
}

// returns a set associated with the passed messageType unioned with the
// generic return
func (bt *byType) Get(messageType int32) *set.Set {
	lookup, ok := bt.list[messageType]
	if !ok {
		return bt.generic
	} else {
		return lookup.Union(bt.generic)
	}
}

// adds a listener to a set for the given messageType. Creates a new set to add
// it to if the set does not exist
func (bt *byType) Add(messageType int32, r Listener) *set.Set {
	s, ok := bt.list[messageType]
	if !ok {
		s = set.New(r)
		bt.list[messageType] = s
	} else {
		s.Insert(r)
	}

	return s
}

// Removes the passed listener from the set for messageType and
// deletes the set if it is empty and the type is not AnyType
func (bt *byType) Remove(mt int32, l Listener) {
	s, ok := bt.list[mt]
	if ok {
		s.Remove(l)

		if s.Len() == 0 && mt != AnyType {
			delete(bt.list, mt)
		}
	}
}
