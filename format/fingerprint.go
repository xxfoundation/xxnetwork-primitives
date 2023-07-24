////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package format

import (
	"encoding/base64"
	"encoding/json"

	"github.com/pkg/errors"
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

// MarshalJSON adheres to the json.Marshaler interface.
func (fp Fingerprint) MarshalJSON() ([]byte, error) {
	return json.Marshal(fp[:])
}

// UnmarshalJSON adheres to the json.Unmarshaler interface.
func (fp *Fingerprint) UnmarshalJSON(data []byte) error {
	var fpBytes []byte
	err := json.Unmarshal(data, &fpBytes)
	if err != nil {
		return err
	}

	if len(fpBytes) != KeyFPLen {
		return errors.Errorf("length of fingerprint must be %d", KeyFPLen)
	}

	copy(fp[:], fpBytes[:])

	return nil
}
