package switchboard

import (
	"github.com/golang-collections/collections/set"
)

type byType struct {
	list    map[int32]*set.Set
	generic *set.Set
}

func newByType() *byType {
	bt := &byType{
		list:    make(map[int32]*set.Set),
		generic: set.New(),
	}

	//make the zero messages, which are defined as any, all point to the generic
	bt.list[AnyType] = bt.generic

	return bt
}

func (bt *byType) Get(messageType int32) *set.Set {
	lookup, ok := bt.list[messageType]
	if !ok {
		return bt.generic
	} else {
		return lookup.Union(bt.generic)
	}
}

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
