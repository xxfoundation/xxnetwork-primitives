package switchboard

import (
	"github.com/golang-collections/collections/set"
	"gitlab.com/elixxir/primitives/id"
)

type byId struct {
	list    map[id.ID]*set.Set
	generic *set.Set
}

func newById() *byId {
	bi := &byId{
		list:    make(map[id.ID]*set.Set),
		generic: set.New(),
	}

	//make the zero IDs, which are defined as any, all point to the generic
	bi.list[id.ZeroUser] = bi.generic
	bi.list[id.ZeroID] = bi.generic

	return bi
}

func (bi *byId) Get(uid *id.ID) *set.Set {
	lookup, ok := bi.list[*uid]
	if !ok {
		return bi.generic
	} else {
		return lookup.Union(bi.generic)
	}
}

func (bi *byId) Add(uid *id.ID, r Listener) *set.Set {
	s, ok := bi.list[*uid]
	if !ok {
		s = set.New(r)
		bi.list[*uid] = s
	} else {
		s.Insert(r)
	}

	return s
}
