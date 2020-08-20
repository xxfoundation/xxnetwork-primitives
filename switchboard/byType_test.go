package switchboard

import (
	"github.com/golang-collections/collections/set"
	"testing"
)

func TestByType_newByType(t *testing.T) {
	nbt := newByType()

	if nbt.list == nil {
		t.Errorf("No list created")
	}

	if nbt.generic == nil {
		t.Errorf("No generic created")
	}

	if nbt.generic != nbt.list[0] {
		t.Errorf("zero message type not registered as generic")
	}

}

func TestByType_Get_Empty(t *testing.T) {
	nbt := newByType()

	s := nbt.Get(42)

	if s.Len() != 0 {
		t.Errorf("Should not have returned a set")
	}
}

func TestByType_Get_Selected(t *testing.T) {
	nbt := newByType()

	m := int32(42)

	set1 := set.New(0)

	nbt.list[m] = set1

	s := nbt.Get(m)

	if s.Len() == 0 {
		t.Errorf("Should have returned a set")
	}

	if !s.SubsetOf(set1) || !set1.SubsetOf(s) {
		t.Errorf("Wrong set returned")
	}
}

func TestByType_Get_Generic(t *testing.T) {
	nbt := newByType()

	m := int32(42)

	nbt.generic.Insert(0)

	s := nbt.Get(m)

	if s.Len() == 0 {
		t.Errorf("Should have returned a set")
	}

	if !s.SubsetOf(nbt.generic) || !nbt.generic.SubsetOf(s) {
		t.Errorf("Wrong set returned")
	}
}

func TestByType_Get_GenericSelected(t *testing.T) {
	nbt := newByType()

	m := int32(42)

	nbt.generic.Insert(1)

	set1 := set.New(0)

	nbt.list[m] = set1

	s := nbt.Get(m)

	if s.Len() == 0 {
		t.Errorf("Should have returned a set")
	}

	setUnion := set1.Union(nbt.generic)

	if !s.SubsetOf(setUnion) || !setUnion.SubsetOf(s) {
		t.Errorf("Wrong set returned")
	}
}
