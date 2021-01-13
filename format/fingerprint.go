///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package format

type Fingerprint [KeyFPLen]byte

// NewFingerprint generates a new Fingerprint with the provided bytes.
func NewFingerprint(b []byte) Fingerprint {
	var fp Fingerprint
	copy(fp[:], b[:])
	return fp
}
