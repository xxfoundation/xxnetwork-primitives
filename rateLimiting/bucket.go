////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package rateLimiting

import (
	"sync"
	"time"
)

// This is an implementation of the leaky bucket algorithm:
// https://en.wikipedia.org/wiki/Leaky_bucket

// Bucket structure tracks the capacity and rate at which the remaining capacity
// decreases.
type Bucket struct {
	capacity   uint      // Maximum number of items the bucket can hold
	remaining  uint      // Current number of items in the bucket
	leakRate   float64   // Rate that the bucket leaks at [items/nanosecond]
	lastUpdate time.Time // Time that the bucket was most recently updated
	sync.Mutex
}

// Create generates a empty bucket with the specified capacity and leak rate.
func Create(capacity uint, leakRate float64) *Bucket {
	return &Bucket{
		capacity:   capacity,
		remaining:  0, // Start with an empty bucket
		leakRate:   leakRate,
		lastUpdate: time.Now(),
	}
}

// Capacity returns the max number of items allowed in the bucket.
func (b *Bucket) Capacity() uint {
	return b.capacity
}

// Remaining returns the remaining space in the bucket.
func (b *Bucket) Remaining() uint {
	return b.remaining
}

// Add adds the specified number of items to the bucket and updates the number
// of items remaining. Returns true if the items were added; otherwise, returns
// false if there was insufficient capacity to do so.
func (b *Bucket) Add(items uint) bool {
	b.Lock()
	defer b.Unlock()

	// Update the number of remaining items in the bucket
	b.update()

	// Add the items to the bucket
	b.remaining += items

	// If the items went over capacity, return false
	if b.remaining > b.capacity {
		return false
	} else {
		return true
	}
}

// update calculates and updates the number of remaining items in the bucket
// since the last update. This function is not thread-safe. It must only be used
// in a thread-safe manner.
func (b *Bucket) update() {
	// Calculate the time elapsed since the last update, in nanoseconds
	elapsedTime := time.Since(b.lastUpdate).Nanoseconds()

	// Calculate the number of items that have leaked over the elapsed time
	itemsLeaked := uint(float64(elapsedTime) * b.leakRate)

	// Update the number of remaining items in the bucket by subtract the number
	// of leaked items from the remaining items ensuring that remaining is no
	// less than zero
	if itemsLeaked > b.remaining {
		b.remaining = 0
	} else {
		b.remaining -= itemsLeaked
	}

	// Update timestamp
	b.lastUpdate = time.Now()
}
