////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package rateLimiting

// Note: Do not remove the error below. It is used in client and gateway.

// ClientRateLimitErr is returned once a user has hit the rate limit.
const ClientRateLimitErr = "Too many messages sent from ID %v with IP address " +
	"%s in a specific time frame by a user"
