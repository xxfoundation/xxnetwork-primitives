////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package rateLimiting

// Error messages to return once a user has hit the rate limti
const (
	ClientRateLimitErr = "Too many messages sent from ID %v with IP address %s in a specific time frame by a user"
)
