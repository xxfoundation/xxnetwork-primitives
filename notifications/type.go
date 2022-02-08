///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package notifications

import "fmt"

// Provider indicates the notification provider that will be called at the
// higher level. There are different notifications providers for different
// mobile operating systems. The specific provider is enumerated below.
type Provider uint32

// Enumeration of different mobile OSes.
const (
	UNKNOWN = Provider(iota) // Unknown is the unfilled zero field.
	APNS
	FCM
	HUAWEI
)

// Stringer to get the name of the Provider.
func (t Provider) String() string {
	switch t {
	case UNKNOWN:
		return "Unknown or unspecified notifications provider"
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
