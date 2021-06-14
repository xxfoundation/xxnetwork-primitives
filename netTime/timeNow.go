///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Package netTime provides a custom time function that should provide the
// current accurate time used by the network from a custom time service.
package netTime

import (
	"time"
)

type NowFunc func() time.Time

// Now returns the current accurate time. The function must be set an accurate
// time service that returns the current time with an accuracy of +/- 300 ms.
var Now NowFunc = time.Now
