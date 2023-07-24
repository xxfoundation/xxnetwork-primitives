////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package rateLimiting

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// Test that CreateBucketFromLeakRatio generates a new bucket with all the
// expected fields.
func TestCreateBucketWithLeakRate(t *testing.T) {
	// Setup expected values
	expectedCapacity := rand.Uint32()
	expectedLeakRate := rand.Float64()

	// Create new Bucket
	b := CreateBucketFromLeakRatio(expectedCapacity, expectedLeakRate, nil)

	// Test fields for expected results
	if b.capacity != expectedCapacity {
		t.Errorf("CreateBucketFromLeakRatio generated Bucket with incorrect "+
			"capacity.\nexpected: %d\nreceived: %d",
			expectedCapacity, b.capacity)
	}

	if b.remaining != 0 {
		t.Errorf("CreateBucketFromLeakRatio generated Bucket with incorrect "+
			"remaining.\nexpected: %d\nreceived: %d", 0, b.remaining)
	}

	if b.leakRate != expectedLeakRate {
		t.Errorf("CreateBucketFromLeakRatio generated Bucket with incorrect "+
			"leak rate.\nexpected: %f\nreceived: %f",
			expectedLeakRate, b.leakRate)
	}

	// Check that the lastUpdate occurred recently
	if time.Now().UnixNano()-b.lastUpdate > time.Second.Nanoseconds() {
		t.Errorf("CreateBucketFromLeakRatio generated Bucket with old "+
			"lastUpdate.\nreceived: %d", b.lastUpdate)
	}

	if b.locked {
		t.Errorf("CreateBucketFromLeakRatio generated Bucket with incorrect "+
			"lock.\nexpected: %t\nreceived: %t", false, b.locked)
	}
}

// Test that CreateBucket generates a new bucket with all the expected fields.
func TestCreateBucket(t *testing.T) {
	// Setup expected values
	expectedCapacity := rand.Uint32()
	expectedLeakRate := 3.0 / float64(5*time.Millisecond)

	// Create new Bucket
	b := CreateBucket(expectedCapacity, 3.0, 5*time.Millisecond, nil)

	if b.leakRate != expectedLeakRate {
		t.Errorf("New Bucket has incorrect leak rate."+
			"\nexpected: %f\nreceived: %f", expectedLeakRate, b.leakRate)
	}
}

// Test that CreateBucketFromDB produces expected bucket.
func TestCreateBucketFromDB(t *testing.T) {
	// Setup expected values
	expectedCapacity := rand.Uint32()
	expectedLeaked := rand.Uint32()
	expectedLeakDuration := time.Duration(rand.Uint32())
	expectedInBucket := rand.Uint32()
	timestamp := time.Now().UnixNano()
	updateDB := func(uint32, int64) {}

	testBucket := CreateBucketFromDB(expectedCapacity, expectedLeaked,
		expectedLeakDuration, expectedInBucket, timestamp, updateDB)

	expectedBucket := &Bucket{
		capacity:   expectedCapacity,
		remaining:  expectedInBucket,
		leakRate:   float64(expectedLeaked) / float64(expectedLeakDuration.Nanoseconds()),
		lastUpdate: timestamp,
		locked:     false,
		whitelist:  false,
		updateDB:   updateDB,
	}

	expectedJson, err := expectedBucket.MarshalJSON()
	if err != nil {
		t.Fatalf(err.Error())
	}

	testJson, err := testBucket.MarshalJSON()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !bytes.Equal(expectedJson, testJson) {
		t.Errorf("CreateBucketFromDB produced an incorrect bucket."+
			"\nexepcted: %+v\nreceived: %+v", expectedBucket, testBucket)

	}
}

// Tests that CreateBucketFromParams produces the expected bucket.
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
		t.Errorf("CreateBucketFromParams produced an incorrect bucket."+
			"\nexepcted: %+v\nreceived: %+v", expectedBucket, testBucket)
	}
}

// Tests that Bucket.Capacity returns the correct value for a new Bucket.
func TestBucket_Capacity(t *testing.T) {
	// Setup expected values
	expectedCapacity := rand.Uint32()

	// Create new Bucket
	b := CreateBucketFromLeakRatio(expectedCapacity, rand.Float64(), nil)

	if b.Capacity() != expectedCapacity {
		t.Errorf("Capacity returned incorrect capacity."+
			"\nexpected: %d\nreceived: %d", expectedCapacity, b.Capacity())
	}
}

// Tests that Bucket.Remaining returns the correct value for a new Bucket.
func TestBucket_Remaining(t *testing.T) {
	// Create new Bucket
	b := CreateBucketFromLeakRatio(rand.Uint32(), rand.Float64(), nil)

	if b.Remaining() != 0 {
		t.Errorf("Remaining returned incorrect remaining for the bucket."+
			"\nexpected: %d\nreceived: %d", 0, b.Remaining())
	}
}

// Tests that IsLocked returns false for a new Bucket.
func TestBucket_IsLocked(t *testing.T) {
	// Create new Bucket
	b := CreateBucketFromLeakRatio(rand.Uint32(), rand.Float64(), nil)

	if b.IsLocked() {
		t.Errorf("IsLocked returned incorrect locked status."+
			"\nexpected: %t\nreceived: %t", false, b.IsLocked())
	}
}

// Tests that Bucket.IsWhitelisted returns false for a new Bucket.
func TestBucket_IsWhitelisted(t *testing.T) {
	// Create new Bucket
	b := CreateBucketFromLeakRatio(rand.Uint32(), rand.Float64(), nil)

	if b.IsWhitelisted() {
		t.Errorf("IsWhitelisted returned incorrect whitelist status."+
			"\nexpected: %t\nreceived: %t", false, b.IsWhitelisted())
	}
}

// Tests that Bucket.IsFull returns false for a new Bucket and true for a
// filled bucket.
func TestBucket_IsFull(t *testing.T) {
	// Create new Bucket
	b := CreateBucketFromLeakRatio(rand.Uint32()/2, rand.Float64(), nil)

	if b.IsFull() {
		t.Errorf("IsFull returned incorrect value for new bucket."+
			"\nexpected: %t\nreceived: %t", false, b.IsFull())
	}

	b.Add(b.capacity * 2)

	if !b.IsFull() {
		t.Errorf("IsFull returned incorrect value for a filled bucket."+
			"\nexpected: %t\nreceived: %t", true, b.IsFull())
	}
}

// Tests that Bucket.IsEmpty returns true for a new Bucket and false for a
// filled bucket.
func TestBucket_IsEmpty(t *testing.T) {
	// Create new Bucket
	b := CreateBucketFromLeakRatio(rand.Uint32(), rand.Float64(), nil)

	if !b.IsEmpty() {
		t.Errorf("IsEmpty returned incorrect value for new bucket."+
			"\nexpected: %t\nreceived: %t", true, b.IsFull())
	}

	b.Add(b.capacity / 2)

	if b.IsEmpty() {
		t.Errorf("IsEmpty returned incorrect value for a filled bucket."+
			"\nexpected: %t\nreceived: %t", false, b.IsEmpty())
	}
}

// Bucket.Add happy path.
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

		if success, _ := b.Add(r.tokensToAdd); !success {
			t.Errorf("Add(%d) added tokens past bucket capacity (round %d). "+
				"[cap: %d, rem: %d]", r.tokensToAdd, i, b.capacity, b.remaining)
		}

		if b.remaining != r.expectedRem {
			t.Errorf("Incorrect number of tokens remaining in bucket (round %d)."+
				"\nexpected: %d\nreceived: %d", i, r.expectedRem, b.remaining)
		}
	}
}

// Tests that Bucket.Add returns false when adding tokens over capacity.
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

		if success, _ := b.Add(r.tokensToAdd); success != r.expectedAddReturn {
			t.Errorf("Add(%d) added tokens past bucket capacity (round %d). "+
				"[cap: %d, rem: %d]", r.tokensToAdd, i, b.capacity, b.remaining)
		}

		if b.remaining != r.expectedRem {
			t.Errorf("Incorrect number of tokens remaining in bucket (round %d)."+
				"\nexpected: %d\nreceived: %d", i, r.expectedRem, b.remaining)
		}
	}
}

// Tests that Bucket.Add updates the database bucket when it is enabled.
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
	db := &BucketParams{
		Key:        "keyA",
		Capacity:   b.capacity,
		Remaining:  b.remaining,
		LeakRate:   b.leakRate,
		LastUpdate: b.lastUpdate,
		Locked:     b.locked,
		Whitelist:  b.whitelist,
	}

	b.updateDB = func(remaining uint32, lastUpdate int64) {
		db.Remaining = remaining
		db.LastUpdate = lastUpdate
	}

	// Add expected values and test
	for i, r := range testData {
		time.Sleep(r.sleepTime * duration)

		if success, _ := b.Add(r.tokensToAdd); !success {
			t.Errorf("Add(%d) added tokens past bucket capacity (round %d). "+
				"[cap: %d, rem: %d]", r.tokensToAdd, i, b.capacity, b.remaining)
		}

		if b.remaining != r.expectedRem {
			t.Errorf("Incorrect number of tokens remaining in bucket (round %d)."+
				"\nexpected: %d\nreceived: %d", i, r.expectedRem, b.remaining)
		}

		if b.remaining != db.Remaining {
			t.Errorf("Incorrect number of tokens remaining in database bucket "+
				"(round %d).\nexpected: %d\nreceived: %d",
				i, b.remaining, db.Remaining)
		}

		if b.lastUpdate != db.LastUpdate {
			t.Errorf("Incorrect LastUpdate in database bucket (round %d)."+
				"\nexpected: %d\nreceived: %d", i, b.lastUpdate, db.LastUpdate)
		}
	}
}

// Tests that Bucket.Add always returns true for a whitelisted bucket.
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

		if success, _ := b.Add(r.tokensToAdd); !success {
			t.Errorf("Add(%d) failed on a whitelisted bucket (round %d). "+
				"[cap: %d, rem: %d]", r.tokensToAdd, i, b.capacity, b.remaining)
		}

		if b.remaining != r.expectedRem {
			t.Errorf("Incorrect number of tokens remaining in bucket (round %d)."+
				"\nexpected: %d\nreceived: %d", i, r.expectedRem, b.remaining)
		}
	}
}

// Tests that Bucket.Add is thread safe.
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
		t.Errorf("Add did not correctly lock the thread.")
	case <-time.After(50 * time.Millisecond):
		return
	}
}

// Tests that Bucket.MarshalJSON and Bucket.UnmarshalJSON can serialize and
// deserialize buckets
func TestBucket_MarshalUnmarshal(t *testing.T) {
	b := CreateBucketFromLeakRatio(10, 1, nil)

	data, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	var b2 Bucket
	err = json.Unmarshal(data, &b2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(b, &b2) {
		t.Error("buckets should be equal after serialization/deserialization")
	}
}

// Tests that addToDB can be set up after marshalling/unmarshalling
func TestBucket_addToDB(t *testing.T) {
	called := false

	var b2 Bucket
	b2.SetAddToDB(func(u uint32, i int64) {
		called = true
	})
	b2.updateDB(0, 0)
	if !called {
		t.Error("addToDb should have been called")
	}
}
