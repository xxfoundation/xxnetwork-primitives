package id

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
)

// Tests that setting the bytes with NewNodeFromBytes() populates the node ID with all
// the same bytes.
func TestNodeID_SetBytes(t *testing.T) {
	idBytes := make([]byte, NodeIdLen)
	rand.Read(idBytes)
	id := NewNodeFromBytes(idBytes)

	if !bytes.Equal(id[:], idBytes) {
		t.Errorf("NewNodeFromBytes() incorrectly set the NodeID bytes"+
			"\n\treceived: %v\n\texpected: %v", id[:], idBytes)
	}
}

// Tests that providing invalid input (wrong length) to NewNodeFromBytes() returns an
// array of all zeros.
func TestNodeID_SetBytes_Error(t *testing.T) {
	var idBytes []byte
	id := NewNodeFromBytes(idBytes)

	if !bytes.Equal(id[:], make([]byte, NodeIdLen)) {
		t.Errorf("NewNodeFromBytes() on nil data did not set all bytes to zero"+
			"\n\treceived: %v\n\texpected: %v", id[:], make([]byte, NodeIdLen))
	}
}

// Tests that Bytes() correctly converts a node ID to an identical byte slice.
func TestNodeID_Bytes(t *testing.T) {
	idBytes := make([]byte, NodeIdLen)
	rand.Read(idBytes)
	id := NewNodeFromBytes(idBytes)

	if !bytes.Equal(id[:], id.Bytes()) {
		t.Errorf("Bytes() returned incorrect byte slice of NodeID"+
			"\n\treceived: %v\n\texpected: %v", id[:], idBytes)
	}
}

// Tests that Bytes() correctly makes a new copy of the bytes.
func TestNodeID_Bytes_Copy(t *testing.T) {
	idBytes := make([]byte, NodeIdLen)
	rand.Read(idBytes)
	id := NewNodeFromBytes(idBytes)

	nodeBytes := id.Bytes()

	// Modify the original
	for j := 0; j < NodeIdLen; j++ {
		id[j] = ^id[j]
	}

	if !bytes.Equal(nodeBytes, idBytes) {
		t.Errorf("Bytes() returned incorrect byte slice of Node ID"+
			"\n\treceived: %v\n\texpected: %v", nodeBytes, idBytes)
	}
}

// Tests that Cmp() returns true when two node IDs are equal and returns false
// when they are not equal.
func TestNodeID_Cmp(t *testing.T) {
	id1 := Node{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	id2 := Node{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	id3 := Node{1, 2, 3, 4, 5, 6, 7, 8, 9, 11}

	if !id1.Cmp(&id2) {
		t.Errorf("Cmp() incorrectly determined the two IDs are not equal"+
			"\n\treceived: %v\n\texpected: %v", id1, id2)
	}

	if id3.Cmp(&id1) {
		t.Errorf("Cmp() incorrectly determined the two IDs are equal"+
			"\n\treceived: %v\n\texpected: %v", id3, id2)
	}
}

// Test that DeepCopy() returns an exact copy of the node ID and that changing
// the original node ID does not change the newly created one.
func TestNodeID_DeepCopy(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	for i := 0; i < 100; i++ {
		var original Node
		rng.Read(original[:])
		deepCopy := (&original).DeepCopy()

		if !reflect.DeepEqual(original, *deepCopy) {
			t.Errorf("DeepCopy() produced a copy that does not equal the "+
				"original on attempt #%d\n\treceived: %v\n\texpected: %v",
				i, original, *deepCopy)
		}

		// Modify the original
		for j := 0; j < NodeIdLen; j++ {
			original[j] = ^original[j]
		}

		if reflect.DeepEqual(original, *deepCopy) {
			t.Errorf("DeepCopy() produced a copy that is linked to the "+
				"original on attempt #%d\n\treceived: %v\n\texpected: %v",
				i, *deepCopy, original)
		}
	}
}

func TestNodeID_DeepCopy_Error(t *testing.T) {
	var original *Node
	deepCopy := original.DeepCopy()

	if deepCopy != nil {
		t.Errorf("DeepCopy() did not return nil when NodeID is nil"+
			"\n\treceived: %v\n\texpected: %v",
			deepCopy, original)
	}
}

func TestNode_String(t *testing.T) {
	// A node ID should produce the same string each time if the underlying data
	// is the same
	id1 := Node{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	id2 := Node{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// However, if you change it, the produced string should differ
	id3 := Node{1, 2, 3, 4, 5, 6, 7, 8, 9, 11}

	if id1.String() != id2.String() {
		t.Error("id1 and id2 are identical, " +
			"and the strings they produce should be identical, but aren't")
	}
	if id3.String() == id1.String() {
		t.Error("id1 and id3 are not identical, " +
			"and the strings they produce should not be identical, " +
			"but they are")
	}
}
