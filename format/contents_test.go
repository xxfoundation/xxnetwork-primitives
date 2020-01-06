////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"bytes"
	"encoding/base64"
	"math/rand"
	"testing"
)

// Tests that NewContents() properly sets Content's serial and position.
func TestNewContents(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, ContentsLen)
	rand.Read(randSlice)

	// Create new Contents
	c := NewContents(randSlice)

	// Test fields
	if !bytes.Equal(c.serial, randSlice) {
		t.Errorf("NewContents() did not properly set Content's serial"+
			"\n\treceived: %v\n\texpected: %v",
			c.serial, randSlice)
	}

	if c.position != invalidPosition {
		t.Errorf("NewContents() did not properly set Content's position"+
			"\n\treceived: %v\n\texpected: %v",
			c.position, invalidPosition)
	}

	// Check serial's length
	if len(c.serial) != ContentsLen {
		t.Errorf("NewContents() did not create a serial with the correct length"+
			"\n\treceived: %v\n\texpected: %v",
			len(c.serial), ContentsLen)
	}
}

// Tests that NewContents() panics when the new serial is not the same length as
// serial.
func TestNewContents_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, ContentsLen-5)
	rand.Read(randSlice)

	// Defer to an error when NewContents() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewContents() did not panic when expected")
		}
	}()

	// Create new Contents with wrong size
	NewContents(randSlice)
}

// Tests that Contents is always constructed the same way.
func TestContents_Consistency(t *testing.T) {
	// Define the serial to check against (64-bit encoded)
	serial, _ := base64.StdEncoding.DecodeString("U4x/lrFkvxuXu59LtHLon1s" +
		"UhPJSCcnZND6SugndnVLf15tNdkKbYXoMn58NO6VbDMDWFEyIhTWEGsvgcJsHWAg/Yd" +
		"N1vAK0HfT5GSnhj9qeb4LlTnSOgeeeS71v40zcuoQ+6NY+jE/+HOvqVG2PrBPdGqwEz" +
		"i6ih3xVec+ix44bC6+uiBuCp1EQikLtPJA8qkNGWnhiBhaXiu0M48bE8657w+BJW1cS" +
		"/v2+DBAoh+EA2s0tiF9pLLYH2gChHBxwceeWotwtwlpbdLLhKXBeJz8FySMmgo4rBW4" +
		"4F2WOEGFJiUf980RBDtTBFgI/qONXa2/tJ/+JdLrAyv2a0FaSsTYZ5ziWTf3Hno1TQ3" +
		"NmHP1m10/sHhuJSRq3I25LdSFikM8r60LDyicyhWDxqsBnzqbov0bUqytGgEAsX7KCD" +
		"ohdMmDx3peCg9Sgmjb5bCCUF0bj7U2mRqmui0+ntPw6ILr6GnXtMnqGuLDDmvHP0rO1" +
		"EhnqeVM6v0SNLEedMmB1M5BZFMjMHPCdo54Okp0C")

	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, ContentsLen)
	rand.Read(randSlice)

	// Create new Contents
	c := NewContents(randSlice)

	if !bytes.Equal(c.Get(), serial) {
		t.Errorf("Contents's serial does not match the hardcoded serial"+
			"\n\treceived: %v\n\texpected: %v",
			c.Get(), serial)
	}
}

// Tests that Get() returns the correct bytes set to Content's serial.
func TestContents_Get(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, ContentsLen)
	rand.Read(randSlice)

	// Create new Contents
	c := NewContents(randSlice)

	if !bytes.Equal(c.Get(), randSlice) {
		t.Errorf("Get() did not return the correct data from "+
			"Content's serial"+
			"\n\treceived: %v\n\texpected: %v",
			c.Get(), randSlice)
	}
}

// Tests that Set() sets the correct bytes to Content's serial and copies the
// correct number of bytes.
func TestContents_Set(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, ContentsLen)
	rand.Read(randSlice)

	// Create new Contents
	c := NewContents(make([]byte, ContentsLen))

	// Set Content's serial
	c.Set(randSlice)

	if !bytes.Equal(c.Get(), randSlice) {
		t.Errorf("Set() did not properly set Content's serial"+
			"\n\treceived: %v\n\texpected: %v",
			c.Get(), randSlice)
	}
}

// Tests that Set() panics when the new serial is not the same length as serial.
func TestContents_Set_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	serialRand := make([]byte, ContentsLen-5)
	rand.Read(serialRand)

	// Create new Contents
	c := NewContents(make([]byte, ContentsLen))

	// Defer to an error when Set() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Set() did not panic when expected")
		}
	}()

	// Set Content's serial
	c.Set(serialRand)
}

// Tests that GetRightAligned() returns the correct right-aligned bytes set to
// Content's serial.
func TestContents_GetRightAligned(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, ContentsLen)
	rand.Read(randSlice)

	// Create new Contents and set position
	c := NewContents(randSlice)
	c.position = 15

	if !bytes.Equal(c.GetRightAligned(), randSlice[c.position:]) {
		t.Errorf("GetRightAligned() did not return the correct "+
			"right-aligned data from Content's serial"+
			"\n\treceived: %v\n\texpected: %v",
			c.GetRightAligned(), randSlice[c.position:])
	}
}

// Tests that GetRightAligned() panics when position is invalid.
func TestContents_GetRightAligned_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, ContentsLen)
	rand.Read(randSlice)

	// Defer to an error when GetRightAligned() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetRightAligned() did not panic when expected")
		}
	}()

	// Create new Contents and get right aligned data
	c := NewContents(randSlice)
	c.GetRightAligned()
}

// Tests that SetRightAligned() sets the correct bytes right-aligned to Content's
// serial.
func TestContents_SetRightAligned(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, ContentsLen-PadMinLen)
	rand.Read(randSlice)

	// Create new Contents
	c := NewContents(make([]byte, ContentsLen))

	// Set Content's serial right-aligned
	size := c.SetRightAligned(randSlice)

	if !bytes.Equal(c.GetRightAligned(), randSlice) {
		t.Errorf("SetRightAligned() did not properly set "+
			"Content's serial"+
			"\n\treceived: %v\n\texpected: %v",
			c.GetRightAligned(), randSlice)
	}

	if size != ContentsLen-PadMinLen {
		t.Errorf("SetRightAligned() did not copy the correct number of"+
			"bytes into Content's serial"+
			"\n\treceived: %v\n\texpected: %v",
			size, ContentsLen-PadMinLen)
	}
}

// Tests that SetRightAligned() panics when the new serial is not smaller than
// contents minus the minimum padding length.
func TestContents_SetRightAligned_Panic(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, ContentsLen)
	rand.Read(randSlice)

	// Create new Contents
	c := NewContents(make([]byte, ContentsLen))

	// Defer to an error when SetRightAligned() does not panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetRightAligned() did not panic when expected")
		}
	}()

	// Set Content's serial right-aligned
	c.SetRightAligned(randSlice)
}

// Tests that Content's position is set correctly after calling
// SetRightAligned().
func TestContents_GetPosition(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, ContentsLen-PadMinLen)
	rand.Read(randSlice)

	// Create new Contents
	c := NewContents(make([]byte, ContentsLen))

	// Set Content's serial right-aligned
	c.SetRightAligned(randSlice)

	if c.GetPosition() != PadMinLen {
		t.Errorf("GetPosition() did not return the correct content "+
			"starting position"+
			"\n\treceived: %v\n\texpected: %v",
			c.GetPosition(), PadMinLen)
	}
}

// Tests that changes made to a copy by DeepCopy() does not reflect to the
// original contents.
func TestContents_DeepCopy(t *testing.T) {
	// Generate random byte slice
	rand.Seed(42)
	randSlice := make([]byte, ContentsLen)
	randSlice2 := make([]byte, ContentsLen)
	rand.Read(randSlice)

	// Create new Contents
	c := NewContents(randSlice)

	// Create copy and change the serial
	contentsCopy := c.DeepCopy()
	rand.Read(randSlice2)
	contentsCopy.Set(randSlice2)

	if bytes.Equal(c.serial, contentsCopy.serial) {
		t.Errorf("DeepCopy() did not properly create a new copy of Contents"+
			"\n\treceived: %v\n\texpected: %v",
			c.serial, contentsCopy.serial)
	}
}
