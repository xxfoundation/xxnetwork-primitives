package format

import (
	"bytes"
	"math/rand"
	"testing"
)

// Happy path.
func TestNewFingerprint(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	fpBytes := make([]byte, KeyFPLen)
	prng.Read(fpBytes)

	fp := NewFingerprint(fpBytes)
	if !bytes.Equal(fpBytes, fp[:]) {
		t.Errorf("NewFingerprint() failed to copy the correct bytes into the "+
			"Fingerprint.\nexpected: %+v\nreceived: %+v", fpBytes, fp)
	}

	// Ensure that the data is copied
	fpBytes[2] = 'x'

	if fp[2] == 'x' {
		t.Errorf("NewFingerprint() failed to create a copy of the data.")
	}
}
