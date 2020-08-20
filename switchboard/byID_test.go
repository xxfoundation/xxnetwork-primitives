package switchboard

import (
	"github.com/golang-collections/collections/set"
	"gitlab.com/elixxir/primitives/id"
	"testing"
)

func TestById_newById(t *testing.T) {
	nbi := newById()

	if nbi.list == nil {
		t.Errorf("No list created")
	}

	if nbi.generic == nil {
		t.Errorf("No generic created")
	}

	if nbi.generic != nbi.list[id.ZeroUser] {
		t.Errorf("zero user not registered as generic")
	}

	if nbi.generic != nbi.list[id.ZeroID] {
		t.Errorf("zero id not registered as generic")
	}
}

func TestById_Get_Empty(t *testing.T) {
	nbi := newById()

	uid := id.NewIdFromUInt(42, id.User, t)

	s := nbi.Get(uid)

	if s.Len() != 0 {
		t.Errorf("Should not have returned a set")
	}
}

func TestById_Get_Selected(t *testing.T) {
	nbi := newById()

	uid := id.NewIdFromUInt(42, id.User, t)

	set1 := set.New(0)

	nbi.list[*uid] = set1

	s := nbi.Get(uid)

	if s.Len() == 0 {
		t.Errorf("Should have returned a set")
	}

	if !s.SubsetOf(set1) || !set1.SubsetOf(s) {
		t.Errorf("Wrong set returned")
	}
}

func TestById_Get_Generic(t *testing.T) {
	nbi := newById()

	uid := id.NewIdFromUInt(42, id.User, t)

	nbi.generic.Insert(0)

	s := nbi.Get(uid)

	if s.Len() == 0 {
		t.Errorf("Should have returned a set")
	}

	if !s.SubsetOf(nbi.generic) || !nbi.generic.SubsetOf(s) {
		t.Errorf("Wrong set returned")
	}
}

func TestById_Get_GenericSelected(t *testing.T) {
	nbi := newById()

	uid := id.NewIdFromUInt(42, id.User, t)

	set1 := set.New(0)

	nbi.list[*uid] = set1

	nbi.generic.Insert(1)

	s := nbi.Get(uid)

	if s.Len() == 0 {
		t.Errorf("Should have returned a set")
	}

	setUnion := set1.Union(nbi.generic)

	if !s.SubsetOf(setUnion) || !setUnion.SubsetOf(s) {
		t.Errorf("Wrong set returned")
	}
}
