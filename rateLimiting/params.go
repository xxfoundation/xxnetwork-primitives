////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package rateLimiting

import "time"

// MapParams holds the values used for new BucketMaps.
type MapParams struct {
	// Capacity for newly created buckets
	Capacity uint32

	// The leak rate is calculated by LeakedTokens / LeakDuration
	LeakedTokens uint32
	LeakDuration time.Duration

	// How often to look for and discard stale buckets
	PollDuration time.Duration

	// Age of stale buckets when discarded
	BucketMaxAge time.Duration
}

// BucketParams structure holds all the values to save and restore a Bucket.
type BucketParams struct {
	Key        string  // Unique bucket key
	Capacity   uint32  // Maximum number of tokens the bucket can hold
	Remaining  uint32  // Current number of tokens in the bucket
	LeakRate   float64 // Rate that the bucket leaks tokens at [tokens/ns]
	LastUpdate int64   // Time that the bucket was most recently updated
	Locked     bool    // Prevents auto deletion when stale
	Whitelist  bool    // No limit for adding tokens to bucket
}
