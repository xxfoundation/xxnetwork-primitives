////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

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
