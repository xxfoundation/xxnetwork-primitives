////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package rateLimiting

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// Test that CreateBucketFromLeakRatio() generates a new bucket with all the
// expected fields.
func TestCreateBucketWithLeakRate(t *testing.T) {
	// Setup expected values
	expectedCapacity := rand.Uint32()
	expectedLeakRate := rand.Float64()

	// Create new Bucket
	b := CreateBucketFromLeakRatio(expectedCapacity, expectedLeakRate, nil)

	// Test fields for expected results
	if b.capacity != expectedCapacity {
		t.Errorf("CreateBucketFromLeakRatio() generated Bucket with incorrect "+
			"capacity.\n\texpected: %v\n\treceived: %v",
			expectedCapacity, b.capacity)
	}

	if b.remaining != 0 {
		t.Errorf("CreateBucketFromLeakRatio() generated Bucket with incorrect "+
			"remaining.\n\texpected: %v\n\treceived: %v", 0, b.remaining)
	}

	if b.leakRate != expectedLeakRate {
		t.Errorf("CreateBucketFromLeakRatio() generated Bucket with incorrect "+
			"leak rate.\n\texpected: %v\n\treceived: %v",
			expectedLeakRate, b.leakRate)
	}

	// Check that the lastUpdate occurred recently
	if time.Now().UnixNano()-b.lastUpdate > time.Second.Nanoseconds() {
		t.Errorf("CreateBucketFromLeakRatio() generated Bucket with old "+
			"lastUpdate.\n\treceived: %v", b.lastUpdate)
	}

	if b.locked {
		t.Errorf("CreateBucketFromLeakRatio() generated Bucket with incorrect "+
			"lock.\n\texpected: %v\n\treceived: %v",
			expectedLeakRate, b.leakRate)
	}
}

// Test that CreateBucket() generates a new bucket with all the expected fields.
func TestCreateBucket(t *testing.T) {
	// Setup expected values
	expectedCapacity := rand.Uint32()
	expectedLeakRate := 3.0 / float64(5*time.Millisecond)

	// Create new Bucket
	b := CreateBucket(expectedCapacity, 3.0, 5*time.Millisecond, nil)

	if b.leakRate != expectedLeakRate {
		t.Errorf("CreateBucketFromLeakRatio() generated Bucket with incorrect leak rate."+
			"\n\texpected: %v\n\treceived: %v", expectedLeakRate, b.leakRate)
	}
}

// Tests that CreateBucketFromParams() produces the expected bucket.
func TestCreateBucketFromParams(t *testing.T) {
	expectedBucket := &Bucket{
		capacity:   rand.Uint32(),
		remaining:  rand.Uint32(),
		leakRate:   rand.Float64(),
		lastUpdate: time.Now().UnixNano(),
		locked:     true,
		whitelist:  true,
	}

	params := &BucketParams{
		Capacity:   expectedBucket.capacity,
		Remaining:  expectedBucket.remaining,
		LeakRate:   expectedBucket.leakRate,
		LastUpdate: expectedBucket.lastUpdate,
		Locked:     expectedBucket.locked,
		Whitelist:  expectedBucket.whitelist,
	}

	testBucket := CreateBucketFromParams(params, nil)

	if !reflect.DeepEqual(expectedBucket, testBucket) {
		t.Errorf("CreateBucketFromParams() produced an incorrect bucket."+
			"\n\texepcted: %+v\n\treceived: %+v", expectedBucket, testBucket)
	}
}

// Tests that Capacity() returns the correct value for a new Bucket.
func TestBucket_Capacity(t *testing.T) {
	// Setup expected values
	expectedCapacity := rand.Uint32()

	// Create new Bucket
	b := CreateBucketFromLeakRatio(expectedCapacity, rand.Float64(), nil)

	if b.Capacity() != expectedCapacity {
		t.Errorf("Capacity() returned incorrect capacity."+
			"\n\texpected: %v\n\treceived: %v", expectedCapacity, b.Capacity())
	}
}

// Tests that Remaining() returns the correct value for a new Bucket.
func TestBucket_Remaining(t *testing.T) {
	// Create new Bucket
	b := CreateBucketFromLeakRatio(rand.Uint32(), rand.Float64(), nil)

	if b.Remaining() != 0 {
		t.Errorf("Remaining() returned incorrect remaining for the bucket."+
			"\n\texpected: %v\n\treceived: %v", 0, b.Remaining())
	}
}

// Tests that IsLocked() returns false for a new Bucket.
func TestBucket_IsLocked(t *testing.T) {
	// Create new Bucket
	b := CreateBucketFromLeakRatio(rand.Uint32(), rand.Float64(), nil)

	if b.IsLocked() {
		t.Errorf("IsLocked() returned incorrect locked status."+
			"\n\texpected: %v\n\treceived: %v", false, b.IsLocked())
	}
}

// Tests that IsWhitelisted() returns false for a new Bucket.
func TestBucket_IsWhitelisted(t *testing.T) {
	// Create new Bucket
	b := CreateBucketFromLeakRatio(rand.Uint32(), rand.Float64(), nil)

	if b.IsWhitelisted() {
		t.Errorf("IsWhitelisted() returned incorrect whitelist status."+
			"\n\texpected: %v\n\treceived: %v", false, b.IsWhitelisted())
	}
}

// Add() happy path.
func TestBucket_Add(t *testing.T) {
	// Generate test data
	testData := []struct {
		tokensToAdd uint32
		expectedRem uint32
		sleepTime   time.Duration // Multiplied by duration defined below
	}{
		{9, 9, 0},
		{7, 10, 2},
		{10, 10, 5},
		{0, 1, 3},
	}

	// Set up bucket with a leak rate of 3 per millisecond
	duration := 60 * time.Millisecond
	leakRate := 3.0 / float64(duration.Nanoseconds())
	b := CreateBucketFromLeakRatio(10, leakRate, nil)

	// Add expected values and test
	for i, r := range testData {
		time.Sleep(r.sleepTime * duration)

		if !b.Add(r.tokensToAdd) {
			t.Errorf("Add(%d) added tokens past bucket capacity (round %d). "+
				"[cap: %d, rem: %d]", r.tokensToAdd, i, b.capacity, b.remaining)
		}

		if b.remaining != r.expectedRem {
			t.Errorf("Incorrect number of tokens remaining in bucket (round %d)."+
				"\n\texpected: %v\n\treceived: %v",
				i, r.expectedRem, b.remaining)
		}
	}
}

// Tests that Add() returns false when adding tokens over capacity.
func TestBucket_Add_OverCapacity(t *testing.T) {
	// Generate test data
	testData := []struct {
		tokensToAdd       uint32
		expectedRem       uint32
		expectedAddReturn bool
		sleepTime         time.Duration // Multiplied by duration defined below
	}{
		{7, 7, true, 2},
		{9, 11, false, 1},
		{10, 16, false, 1},
		{1, 7, true, 2},
	}

	// Set up bucket with a leak rate of 5 per millisecond
	duration := 30 * time.Millisecond
	leakRate := 5.0 / float64(duration.Nanoseconds())
	b := CreateBucketFromLeakRatio(10, leakRate, nil)

	// Add expected values and test
	for i, r := range testData {
		time.Sleep(r.sleepTime * duration)

		if b.Add(r.tokensToAdd) != r.expectedAddReturn {
			t.Errorf("Add(%d) added tokens past bucket capacity (round %d). "+
				"[cap: %d, rem: %d]", r.tokensToAdd, i, b.capacity, b.remaining)
		}

		if b.remaining != r.expectedRem {
			t.Errorf("Incorrect number of tokens remaining in bucket (round %d)."+
				"\n\texpected: %v\n\treceived: %v", i, r.expectedRem, b.remaining)
		}
	}
}

// Tests that Add() updates the database bucket when it is enabled.
func TestBucket_Add_DB(t *testing.T) {
	// Generate test data
	testData := []struct {
		tokensToAdd uint32
		expectedRem uint32
		sleepTime   time.Duration // Multiplied by duration defined below
	}{
		{9, 9, 0},
		{7, 10, 2},
		{10, 10, 5},
		{0, 1, 3},
	}

	// Set up bucket with a leak rate of 3 per millisecond
	duration := 60 * time.Millisecond
	leakRate := 3.0 / float64(duration.Nanoseconds())
	b := CreateBucketFromLeakRatio(10, leakRate, nil)

	// Set up mock bucket database with addToDb function
	bucketDB := &BucketParams{
		Key:        "keyA",
		Capacity:   b.capacity,
		Remaining:  b.remaining,
		LeakRate:   b.leakRate,
		LastUpdate: b.lastUpdate,
		Locked:     b.locked,
		Whitelist:  b.whitelist,
	}

	b.addToDb = func(remaining uint32, lastUpdate int64) {
		bucketDB.Remaining = remaining
		bucketDB.LastUpdate = lastUpdate
	}

	// Add expected values and test
	for i, r := range testData {
		time.Sleep(r.sleepTime * duration)

		if !b.Add(r.tokensToAdd) {
			t.Errorf("Add(%d) added tokens past bucket capacity (round %d). "+
				"[cap: %d, rem: %d]", r.tokensToAdd, i, b.capacity, b.remaining)
		}

		if b.remaining != r.expectedRem {
			t.Errorf("Incorrect number of tokens remaining in bucket (round %d)."+
				"\n\texpected: %v\n\treceived: %v",
				i, r.expectedRem, b.remaining)
		}

		if b.remaining != bucketDB.Remaining {
			t.Errorf("Incorrect number of tokens remaining in database bucket (round %d)."+
				"\n\texpected: %v\n\treceived: %v",
				i, b.remaining, bucketDB.Remaining)
		}

		if b.lastUpdate != bucketDB.LastUpdate {
			t.Errorf("Incorrect LastUpdate in database bucket (round %d)."+
				"\n\texpected: %v\n\treceived: %v",
				i, b.lastUpdate, bucketDB.LastUpdate)
		}
	}
}

// Tests that Add() always returns true for a whitelisted bucket.
func TestBucket_Add_Whitelist(t *testing.T) {
	// Generate test data
	testData := []struct {
		tokensToAdd uint32
		expectedRem uint32
		sleepTime   time.Duration // Multiplied by duration defined below
	}{
		{9, 9, 0},
		{10, 13, 2},
		{20, 20, 5},
		{0, 11, 3},
	}

	// Set up bucket with a leak rate of 3 per millisecond
	duration := 60 * time.Millisecond
	leakRate := 3.0 / float64(duration.Nanoseconds())
	b := CreateBucketFromLeakRatio(10, leakRate, nil)
	b.whitelist = true

	// Add expected values and test
	for i, r := range testData {
		time.Sleep(r.sleepTime * duration)

		if !b.Add(r.tokensToAdd) {
			t.Errorf("Add(%d) failed on a whitelisted bucket (round %d). "+
				"[cap: %d, rem: %d]", r.tokensToAdd, i, b.capacity, b.remaining)
		}

		if b.remaining != r.expectedRem {
			t.Errorf("Incorrect number of tokens remaining in bucket (round %d)."+
				"\n\texpected: %v\n\treceived: %v",
				i, r.expectedRem, b.remaining)
		}
	}
}

// Tests that Add() is thread safe.
func TestAdd_ThreadSafe(t *testing.T) {
	b := CreateBucketFromLeakRatio(10, 1, nil)
	result := make(chan bool)

	b.Lock()
	defer b.Unlock()

	go func() {
		b.Add(15)
		result <- true
	}()

	select {
	case <-result:
		t.Errorf("Add() did not correctly lock the thread.")
	case <-time.After(50 * time.Millisecond):
		return
	}
}
