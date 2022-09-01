////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"encoding/base64"
)

type Fingerprint [KeyFPLen]byte

// NewFingerprint generates a new Fingerprint with the provided bytes.
func NewFingerprint(b []byte) Fingerprint {
	var fp Fingerprint
	copy(fp[:], b[:])
	return fp
}

// Bytes returns the fingerprint as a byte slice.
func (fp Fingerprint) Bytes() []byte {
	return fp[:]
}

// String returns the fingerprint as a base 64 encoded string. This functions
// satisfies the fmt.Stringer interface.
func (fp Fingerprint) String() string {
	return base64.StdEncoding.EncodeToString(fp.Bytes())
}
