package contact

import (
	"reflect"
	"testing"
)

// Test NewFact() function returns a correctly formatted Fact
func TestNewFact(t *testing.T) {
	// Expected result
	e := Fact{
		Fact: "testing",
		T:    1,
	}

	g, err := NewFact(Email, "testing")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(e, g) {
		t.Errorf("The returned Fact did not match the expected Fact")
	}
}

// Test Get() function correctly gets the fact string
func TestFact_Get(t *testing.T) {
	f := Fact{
		Fact: "testing",
		T:    1,
	}

	if f.Get() != f.Fact {
		t.Errorf("f.Get() did not return the same value as f.Fact")
	}
}

// Test Type() function correctly gets the type number
func TestFact_Type(t *testing.T) {
	f := Fact{
		Fact: "testing",
		T:    1,
	}

	if f.Type() != int(f.T) {
		t.Errorf("f.Type() did not return the same value as int(f.T)")
	}
}

// Test Stringify() creates a string of the Fact
// The output is verified to work in the test below
func TestFact_Stringify(t *testing.T) {
	f := Fact{
		Fact: "testing",
		T:    1,
	}

	expected := "Etesting"
	got := f.Stringify()
	t.Log(got)

	if got != expected {
		t.Errorf("Marshalled object from Got did not match Expected.\n\tGot: %v\n\tExpected: %v", got, expected)
	}
}

// Test the UnstringifyFact function creates a Fact from a string
// NOTE: this test does not pass, with error "Unknown Fact FactType: Etesting"
func TestFact_UnstringifyFact(t *testing.T) {
	// Expected fact from above test
	e := Fact{
		Fact: "testing",
		T:    Email,
	}

	// Stringify-ed Fact from above test
	m := "Etesting"
	f, err := UnstringifyFact(m)
	if err != nil {
		t.Error(err)
	}

	t.Log(f.Fact)
	t.Log(f.T)

	if !reflect.DeepEqual(e, f) {
		t.Errorf("The returned Fact did not match the expected Fact")
	}
}
