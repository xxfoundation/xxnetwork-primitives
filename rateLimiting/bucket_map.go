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

// BucketMap structure is a collection of buckets in a map, each with the same
// initial capacity and leak rate.
type BucketMap struct {
	buckets      map[string]*Bucket // Map of the buckets
	newCapacity  uint               // The capacity of newly created buckets
	newLeakRate  float64            // The leak rate of newly created buckets
	cleanPeriod  time.Duration      // Duration between stale bucket removals
	maxDuration  time.Duration      // Max time of inactivity before removal
	sync.RWMutex                    // Only allows one writer at a time
}

// CreateBucketMap creates and returns a new BucketMap structure with the
// specified initial fields. It also starts the state bucket worker.
func CreateBucketMap(newCapacity uint, newLeakRate float64, cleanPeriod,
	maxDuration time.Duration) *BucketMap {
	newBucketMap := &BucketMap{
		buckets:     make(map[string]*Bucket),
		newCapacity: newCapacity,
		newLeakRate: newLeakRate,
		cleanPeriod: cleanPeriod,
		maxDuration: maxDuration,
	}

	// Start the process to remove state buckets (older than maxDuration)
	go newBucketMap.StaleBucketWorker()

	return newBucketMap
}

// CreateBucketMapFromParams creates a new BucketMap from a parameter object.
func CreateBucketMapFromParams(params Params) *BucketMap {
	return CreateBucketMap(params.Capacity, params.LeakRate, params.CleanPeriod,
		params.MaxDuration)
}

// LookupBucket returns the bucket with the specified key. If no bucket exists,
// then a new one is created, inserted into the map, and returned.
func (bm *BucketMap) LookupBucket(key string) *Bucket {
	// Get the bucket and a boolean determining if it exists in the map
	bm.RLock()
	foundBucket, exists := bm.buckets[key]
	bm.RUnlock()

	// Check if the bucket exists
	if exists {
		// If the bucket exists, then return it
		return foundBucket
	} else {
		// If the bucket does not exist, lock the thread and check check again
		// to ensure no other changes are made

		bm.Lock()
		defer bm.Unlock()

		foundBucket, exists = bm.buckets[key]

		// If the bucket does not exist, then create a new one
		if !exists {
			// NOTE: I was unable to test that the key corresponds to the
			// correct bucket. If you end up putting actual values in here, you
			// may want to test for that.
			bm.buckets[key] = Create(bm.newCapacity, bm.newLeakRate)
			foundBucket = bm.buckets[key]
		}

		return foundBucket
	}
}

// StaleBucketWorker clears stale buckets from the map every cleanPeriod. This
// functions is meant to be run in a separate thread.
func (bm *BucketMap) StaleBucketWorker() {
	// Create a new ticker channel that will send every cleanPeriod
	c := time.Tick(bm.cleanPeriod)

	// Run the ticker
	for range c {
		bm.clearStaleBuckets()
	}
}

// clearStaleBuckets loops through the bucket map and removes stale buckets. A
// stale bucket is a bucket that has not been updated in maxDuration.
func (bm *BucketMap) clearStaleBuckets() {
	bm.RLock()
	defer bm.RUnlock()

	// Loop through each bucket in the map
	for key, bucket := range bm.buckets {
		// Calculate time since the bucket's last update
		duration := time.Since(bucket.lastUpdate)

		// Check if the bucket is stale
		if duration >= bm.maxDuration {
			// Remove the stale bucket from the map
			bm.RUnlock()
			bm.Lock()
			delete(bm.buckets, key)
			bm.Unlock()
			bm.RLock()
		}
	}
}
