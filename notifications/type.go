///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package notifications

import "fmt"

// Type indicates the notification code path that will be followed.
// There are different notification logic in the higher levels for different
// mobile operating systems. The mobile OSes are enumerated below.
type Type uint8

// Enumeration of different mobile OSes.
const (
	APNS = Type(iota)
	FCM
	HUAWEI
)

// Stringer to get the name of the Type.
func (t Type) String() string {
	switch t {
	case APNS:
		return "APNS"
	case FCM:
		return "FCM"
	case HUAWEI:
		return "HUAWEI"
	default:
		return fmt.Sprintf("UNKNOWN TYPE: %d", t)
	}
}
