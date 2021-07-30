package region

import (
	"testing"
)

// Make a Rust array with all the country codes, for the blockchain
func TestMakeRustArray(t *testing.T) {
	for i, s := range countryBins {
		t.Logf("[[%d, %d], %d],", int(i[0]), int(i[1]), s)
	}
}
