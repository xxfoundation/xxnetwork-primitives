////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package rateLimiting

// The bucket map contains a list of leaky buckets that each track and limit
// the rate of usage. The map has an optional database backend where buckets are
// backed up for retrieval on restart.

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
)

// BucketMap structure is a collection of buckets in a map, each with the same
// capacity and leak rate. The structure contains an optional database backend.
type BucketMap struct {
	buckets      map[string]*Bucket // Map of the buckets
	capacity     uint32             // The capacity of buckets in the map
	leakRate     float64            // The leak rate of buckets in the map
	pollDuration time.Duration      // Duration between polls for stale buckets
	bucketMaxAge time.Duration      // Max time of inactivity before removal
	sync.RWMutex

	// Database to back up/restore map from. If no database is being used, then
	// this value should remain nil.
	db Storage
}

// CreateBucketMap creates a new BucketMap structure and starts the stale bucket
// removal thread. The leak rate is calculated by dividing leaked by
// leakDuration. If a database is being used, then all the stored buckets are
// reloaded from storage on creation of the map.
//
// NOTE: If db is nil, then the database will not be used. If the quit channel
// is not provided, then the stale bucket removal service will not start.
func CreateBucketMap(capacity, leaked uint32, leakDuration, pollDuration,
	bucketMaxAge time.Duration, db Storage, quit chan struct{}) *BucketMap {

	// Calculate the leak rate [tokens/nanosecond]
	leakRate := float64(leaked) / float64(leakDuration.Nanoseconds())

	bm := &BucketMap{
		buckets:      make(map[string]*Bucket),
		capacity:     capacity,
		leakRate:     leakRate,
		pollDuration: pollDuration,
		bucketMaxAge: bucketMaxAge,
		db:           db,
	}

	// If the database is enabled, load all the buckets into memory
	if bm.db != nil {
		bm.addAllBuckets(bm.db.RetrieveAllBuckets())
	}

	// Start the process to poll for stale buckets
	if quit != nil {
		go bm.staleBucketWorker(quit)
	}

	return bm
}

// CreateBucketMapFromParams creates a new BucketMap from the buckets MapParams
// structure.
//
// NOTE: If db is nil, then the database will not be used. If the quit channel
// is not provided, then the stale bucket removal service will not start.
func CreateBucketMapFromParams(params *MapParams, db Storage,
	quit chan struct{}) *BucketMap {
	return CreateBucketMap(params.Capacity, params.LeakedTokens,
		params.LeakDuration, params.PollDuration, params.BucketMaxAge, db, quit)
}

// LookupBucket returns the bucket in the map with the specified key. If no
// bucket exists, then a new one is added to the map and returned.
func (bm *BucketMap) LookupBucket(key string) *Bucket {
	bm.RLock()
	// Check if the bucket exists
	foundBucket, exists := bm.buckets[key]
	bm.RUnlock()

	if !exists {
		// If the database is being used, then generate a function to be used to
		// update a database buckets tokens. Otherwise, use nil to signify no
		// database is being used.
		var addToDb func(uint32, int64)
		if bm.db != nil {
			addToDb = bm.createAddToDbFunc(key)
		}

		bm.Lock()
		foundBucket, exists = bm.buckets[key]
		if !exists {
			foundBucket = CreateBucketFromLeakRatio(bm.capacity, bm.leakRate, addToDb)
			bm.buckets[key] = foundBucket
		}
		bm.Unlock()

		// Insert into Storage if enabled
		if bm.db != nil {
			bm.db.UpsertBucket(&BucketParams{
				Key:        key,
				Capacity:   foundBucket.capacity,
				Remaining:  foundBucket.remaining,
				LeakRate:   foundBucket.leakRate,
				LastUpdate: foundBucket.lastUpdate,
				Locked:     foundBucket.locked,
				Whitelist:  foundBucket.whitelist,
			})
		}
	}

	return foundBucket
}

// AddBucket adds a new bucket to the map. The leak rate is calculated by
// dividing leaked by leakDuration.
func (bm *BucketMap) AddBucket(key string, capacity, leaked uint32,
	leakDuration time.Duration) *Bucket {

	// Calculate the leak rate [tokens/nanosecond]
	leakRate := float64(leaked) / float64(leakDuration.Nanoseconds())

	return bm.addBucketFromLeakRatio(key, capacity, leakRate)
}

// addBucketFromLeakRatio adds a new bucket to the map regardless if one already
// exists for the key. The added bucket is locked so that it can only be removed
// manually. If a bucket already exists in the map, the new bucket inherits a
// portion of the tokens of existing bucket. Otherwise, the new bucket is empty.
func (bm *BucketMap) addBucketFromLeakRatio(key string, capacity uint32,
	leakRate float64) *Bucket {

	var addToDb func(uint32, int64)
	if bm.db != nil {
		addToDb = bm.createAddToDbFunc(key)
	}

	// Create new locked bucket
	newBucket := CreateBucketFromLeakRatio(capacity, leakRate, addToDb)
	newBucket.locked = true

	bm.Lock()

	// If a bucket already exists, then the new bucket will use a portion of the
	// current bucket's remaining value.
	bucket, exist := bm.buckets[key]
	if exist {
		newBucket.remaining = bucket.remaining % capacity
	}

	// Insert the new bucket
	bm.buckets[key] = newBucket

	bm.Unlock()

	// Insert into Storage if enabled
	if bm.db != nil {
		bm.db.UpsertBucket(&BucketParams{
			Key:        key,
			Capacity:   newBucket.capacity,
			Remaining:  newBucket.remaining,
			LeakRate:   newBucket.leakRate,
			LastUpdate: newBucket.lastUpdate,
			Locked:     newBucket.locked,
			Whitelist:  newBucket.whitelist,
		})
	}

	return newBucket
}

// addAllBuckets creates a new bucket for each BucketParam and inserts it into
// the map. If a bucket already exists, then it is overwritten.
func (bm *BucketMap) addAllBuckets(params []*BucketParams) {
	for _, bp := range params {
		var addToDb func(uint32, int64)
		if bm.db != nil {
			addToDb = bm.createAddToDbFunc(bp.Key)
		}
		bm.buckets[bp.Key] = CreateBucketFromParams(bp, addToDb)
	}
}

// AddToWhitelist adds the list of entries to the bucket map and set them as
// whitelisted.
func (bm *BucketMap) AddToWhitelist(entries []string) {
	bm.Lock()
	for _, key := range entries {
		bucket, exists := bm.buckets[key]
		if exists {
			bucket.locked = true
			bucket.whitelist = true
		} else {
			var addToDb func(uint32, int64)
			if bm.db != nil {
				addToDb = bm.createAddToDbFunc(key)
			}
			newBucket := CreateBucketFromLeakRatio(bm.capacity, bm.leakRate, addToDb)
			newBucket.locked = true
			newBucket.whitelist = true
			bm.buckets[key] = newBucket
		}
	}
	bm.Unlock()

	if bm.db != nil {
		bm.RLock()
		for _, key := range entries {
			bm.db.UpsertBucket(&BucketParams{
				Key:        key,
				Capacity:   bm.buckets[key].capacity,
				Remaining:  bm.buckets[key].remaining,
				LeakRate:   bm.buckets[key].leakRate,
				LastUpdate: bm.buckets[key].lastUpdate,
				Locked:     bm.buckets[key].locked,
				Whitelist:  bm.buckets[key].whitelist,
			})
		}
		bm.RUnlock()
	}
}

// DeleteBucket removes the bucket with the specified key from the map. If the
// bucket does not exist, then an error is returned.
func (bm *BucketMap) DeleteBucket(key string) error {
	bm.Lock()
	defer bm.Unlock()

	// Return an error if the bucket does not exist
	_, exists := bm.buckets[key]
	if !exists {
		return errors.Errorf("Could not delete bucket with key %s. No bucket "+
			"exists in map.", key)
	}

	// Delete the bucket from the map
	delete(bm.buckets, key)

	// Delete the bucket from the database, if enabled
	if bm.db != nil {
		return bm.db.DeleteBucket(key)
	}

	return nil
}

// staleBucketWorker periodically clears stale buckets from the map every
// pollDuration. The quit channel stops the ticker. This function is meant to be
// run in its own thread.
func (bm *BucketMap) staleBucketWorker(quit chan struct{}) {
	// Create a new ticker that will poll every pollDuration
	ticker := time.NewTicker(bm.pollDuration)

	jww.DEBUG.Printf("Starting StaleBucketWorker in separate thread polling "+
		"every %s.", bm.pollDuration)

	for {
		select {
		case <-ticker.C:
			bm.clearStaleBuckets()
		case <-quit:
			jww.DEBUG.Printf("Stopping StaleBucketWorker thread.")
			ticker.Stop()
			return
		}
	}
}

// clearStaleBuckets loops through the bucket map and removes stale buckets. A
// bucket is stale when the elapsed time since its last update is greater than
// bucketMaxAge, and it has no tokens remaining. Stale buckets that are locked
// are not deleted.
func (bm *BucketMap) clearStaleBuckets() {

	// Get current time for calculating bucket ages
	now := time.Now().UnixNano()

	// Copy the bucket map
	bmCopy := map[string]*Bucket{}
	bm.RLock()
	for key, bucket := range bm.buckets {
		bmCopy[key] = bucket
	}
	bm.RUnlock()

	// Find stale buckets in the map and add keys to a list
	var staleBuckets []string
	for key, b := range bmCopy {
		if !b.locked {
			// Calculate the age of the bucket
			bucketAge := now - b.lastUpdate

			// Add bucket to list if it is stale
			if bucketAge >= bm.bucketMaxAge.Nanoseconds() && b.Remaining() == 0 {
				staleBuckets = append(staleBuckets, key)
			}
		}
	}

	// Delete the stale buckets from the list
	if len(staleBuckets) > 0 {
		bm.Lock()
		for _, key := range staleBuckets {
			delete(bm.buckets, key)
		}
		bm.Unlock()

		// Delete the stale buckets from the database, if enabled
		if bm.db != nil {
			for _, key := range staleBuckets {
				err := bm.db.DeleteBucket(key)
				jww.WARN.Printf("Could not delete stale bucket with key %s: %v",
					key, err)
			}
		}
	}
}

// createAddToDbFunc generates the anonymous function that is passed to a new
// bucket so that it can update the remaining tokens in the database.
func (bm *BucketMap) createAddToDbFunc(key string) func(uint32, int64) {
	return func(remaining uint32, lastUpdate int64) {
		err := bm.db.AddToBucket(key, remaining, lastUpdate)
		if err != nil {
			jww.FATAL.Panicf("Could not add tokens to bucket %s in "+
				"database: %v", key, err)
		}
	}
}
