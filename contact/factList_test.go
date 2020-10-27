package contact

import (
	"reflect"
	"testing"
)

// Test that the num function returns the correct number of Facts
func TestFactList_Num(t *testing.T) {
	e1 := Fact{
		Fact: "testing",
		T:    Phone,
	}
	e2 := Fact{
		Fact: "othertest",
		T:    Email,
	}

	Fs := new([]Fact)
	C := Contact{Facts: *Fs}
	FL := FactList{source: &C}

	FL.source.Facts = append(FL.source.Facts, e1)
	FL.source.Facts = append(FL.source.Facts, e2)

	if FL.Num() != 2 {
		t.Error("Num returned incorrect number of Facts in FactList")
	}
}

// Test the get function gets the right Fact at the index
func TestFactList_Get(t *testing.T) {
	e := Fact{
		Fact: "testing",
		T:    Phone,
	}

	Fs := new([]Fact)
	C := Contact{Facts: *Fs}
	FL := FactList{source: &C}

	FL.source.Facts = append(FL.source.Facts, e)

	if !reflect.DeepEqual(e, FL.Get(0)) {
		t.Error("Expected fact does not match got fact")
	}
}

// Test the add function adds a Fact to the source correctly
func TestFactList_Add(t *testing.T) {
	// Expected fact
	e := Fact{
		Fact: "testing",
		T:    Phone,
	}

	Fs := new([]Fact)
	C := Contact{Facts: *Fs}
	FL := FactList{source: &C}

	err := FL.Add("testing", int(Phone))
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(e, FL.source.Facts[0]) {
		t.Error("Expected fact does not match added fact")
	}
}

func TestFactList_AddInvalidType(t *testing.T) {
	Fs := new([]Fact)
	C := Contact{Facts: *Fs}
	FL := FactList{source: &C}

	err := FL.Add("testing", 200)

	if err == nil {
		t.Error("Adding a Fact to FactList with type 200 did not return an error!")
	}
}
