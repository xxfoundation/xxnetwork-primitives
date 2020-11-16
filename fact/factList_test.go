package fact

import (
	"reflect"
	"testing"
)

func TestFactList_StringifyUnstringify(t *testing.T) {
	expected := FactList{}
	expected = append(expected, Fact{
		Fact: "vivian@elixxir.io",
		T:    Email,
	})
	expected = append(expected, Fact{
		Fact: "(270) 301-5797",
		T:    Phone,
	})

	FlString := expected.Stringify()
	// Manually check and verify that the string version is as expected
	t.Log(FlString)

	actual, _, err := UnstringifyFactList(FlString)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Error("fact lists weren't equal")
	}
}
