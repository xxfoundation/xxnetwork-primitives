package id

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
)

// Tests that setting the bytes with NewGateway() populates the Gateway ID with
// all the same bytes.
func TestGatewayID_NewGateway(t *testing.T) {
	idBytes := make([]byte, GatewayIdLen)
	rand.Read(idBytes)
	nodeId := NewNodeFromBytes(idBytes)
	id := nodeId.NewGateway()

	if !bytes.Equal(id[:], idBytes) {
		t.Errorf("NewGateway() incorrectly set the Gateway ID bytes"+
			"\n\treceived: %v\n\texpected: %v", id[:], idBytes)
	}
}

// Tests that providing invalid input (wrong length) to NewGateway() returns an
// array of all zeros.
func TestGatewayID_NewGateway_Error(t *testing.T) {
	var idBytes []byte
	nodeId := NewNodeFromBytes(idBytes)
	id := nodeId.NewGateway()

	if !bytes.Equal(id[:], make([]byte, GatewayIdLen)) {
		t.Errorf("NewGateway() on nil data did not set all bytes to zero"+
			"\n\treceived: %v\n\texpected: %v", id[:], make([]byte, GatewayIdLen))
	}
}

// Tests that Bytes() correctly converts a Gateway ID to an identical byte
// slice.
func TestGatewayID_Bytes(t *testing.T) {
	idBytes := make([]byte, GatewayIdLen)
	rand.Read(idBytes)
	nodeId := NewNodeFromBytes(idBytes)
	id := nodeId.NewGateway()

	if !bytes.Equal(id.Bytes(), idBytes) {
		t.Errorf("Bytes() returned incorrect byte slice of Gateway ID"+
			"\n\treceived: %v\n\texpected: %v", id.Bytes(), idBytes)
	}
}

// Tests that Bytes() correctly makes a new copy of the bytes.
func TestGatewayID_Bytes_Copy(t *testing.T) {
	idBytes := make([]byte, GatewayIdLen)
	rand.Read(idBytes)
	nodeId := NewNodeFromBytes(idBytes)
	id := nodeId.NewGateway()

	gatewayBytes := id.Bytes()

	// Modify the original
	for j := 0; j < GatewayIdLen; j++ {
		id[j] = ^id[j]
	}

	if !bytes.Equal(gatewayBytes, idBytes) {
		t.Errorf("Bytes() returned incorrect byte slice of Gateway ID"+
			"\n\treceived: %v\n\texpected: %v", gatewayBytes, idBytes)
	}
}

// Tests that Cmp() returns true when two Gateway IDs are equal and returns
// false when they are not equal.
func TestGatewayID_Cmp(t *testing.T) {
	id1 := Gateway{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	id2 := Gateway{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	id3 := Gateway{1, 2, 3, 4, 5, 6, 7, 8, 9, 11}

	if !id1.Cmp(&id2) {
		t.Errorf("Cmp() incorrectly determined the two IDs are not equal"+
			"\n\treceived: %v\n\texpected: %v", id1, id2)
	}

	if id3.Cmp(&id1) {
		t.Errorf("Cmp() incorrectly determined the two IDs are equal"+
			"\n\treceived: %v\n\texpected: %v", id3, id2)
	}
}

// Test that DeepCopy() returns an exact copy of the Gateway ID and that
// changing the original Gateway ID does not change the newly created one.
func TestGatewayID_DeepCopy(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	for i := 0; i < 100; i++ {
		var original Gateway
		rng.Read(original[:])
		deepCopy := (&original).DeepCopy()

		if !reflect.DeepEqual(original, *deepCopy) {
			t.Errorf("DeepCopy() produced a copy that does not equal the "+
				"original on attempt #%d\n\treceived: %v\n\texpected: %v",
				i, original, *deepCopy)
		}

		// Modify the original
		for j := 0; j < GatewayIdLen; j++ {
			original[j] = ^original[j]
		}

		if reflect.DeepEqual(original, *deepCopy) {
			t.Errorf("DeepCopy() produced a copy that is linked to the "+
				"original on attempt #%d\n\treceived: %v\n\texpected: %v",
				i, *deepCopy, original)
		}
	}
}

// Test that DeepCopy() returns an error when the Gateway ID is nil.
func TestGatewayID_DeepCopy_Error(t *testing.T) {
	var original *Gateway
	deepCopy := original.DeepCopy()

	if deepCopy != nil {
		t.Errorf("DeepCopy() did not return nil when GatewayID is nil"+
			"\n\treceived: %v\n\texpected: %v",
			deepCopy, original)
	}
}

// Test that String() produces the same string each time when the underlying
// data is the same.
func TestGateway_String(t *testing.T) {
	id1 := Gateway{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	id2 := Gateway{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	id3 := Gateway{1, 2, 3, 4, 5, 6, 7, 8, 9, 11}

	if id1.String() != id2.String() {
		t.Errorf("id1 and id2 are identical, "+
			"and the strings they produce should be identical, but are not"+
			"\n\treceived: %v\n\texpected: %v",
			id1.String(), id2.String())
	}

	if id3.String() == id1.String() {
		t.Errorf("id1 and id3 are not identical, "+
			"and the strings they produce should not be identical, "+
			"but they are"+
			"\n\treceived: %v\n\texpected: %v",
			id3.String(), id1.String())
	}

	if id3.String() == id1.String() {
		t.Errorf("id1 and id3 are not identical, "+
			"and the strings they produce should not be identical, "+
			"but they are"+
			"\n\treceived: %v\n\texpected: %v",
			id3.String(), id1.String())
	}
}

// Test that String() outputs a string with "-Gateway" appended to it.
func TestGateway_String_Append(t *testing.T) {
	id1 := Gateway{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	if id1.String()[len(id1.String())-8:] != "-Gateway" {
		t.Errorf("Last 8 characters of string incorrect"+
			"\n\treceived: %v\n\texpected: %v",
			id1.String()[len(id1.String())-8:], "-Gateway")
	}
}

func TestNewTmpGateway(t *testing.T) {
	tmp := NewTmpGateway()
	t.Logf(tmp.String())
	expected := "dG1wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=-Gateway"
	if tmp.String() != expected {
		t.Logf("failed creating a new tmp Gateway")
		t.Fail()
	}
}
