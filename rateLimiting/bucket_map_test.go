////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package rateLimiting

import (
	"math"
	"testing"
	"time"
)

// Test functionality of CreateBucketMap() to create empty map, it's capacity,
// and the leak rate.
func TestCreateBucketMap(t *testing.T) {
	bm := CreateBucketMap(5, 0.003, 10*time.Second, 3*time.Second)

	if len(bm.buckets) != 0 {
		t.Errorf("CreateBucketMap() created a non-empty map (length greater than 0)\n\treceived: %v\n\texpected: %v", len(bm.buckets), 0)
	}

	if bm.newCapacity != 5 {
		t.Errorf("CreateBucketMap() created incorrect newCapacity\n\treceived: %v\n\texpected: %v", bm.newCapacity, 5)
	}

	if bm.newLeakRate != 0.003 {
		t.Errorf("CreateBucketMap() created incorrect newLeakRate\n\treceived: %v\n\texpected: %v", bm.newLeakRate, 0.003)
	}
}

// Test functionality of TestCreateBucketMapFromParams() to create empty map,
// it's capacity, and the leak rate.
func TestCreateBucketMapFromParams(t *testing.T) {
	params := Params{
		Capacity:      5,
		LeakRate:      0.003,
		CleanPeriod:   10 * time.Second,
		MaxDuration:   3 * time.Second,
		WhitelistFile: "",
	}
	bm := CreateBucketMapFromParams(params)

	if len(bm.buckets) != 0 {
		t.Errorf("CreateBucketMap() created a non-empty map (length greater than 0)."+
			"\n\treceived: %v\n\texpected: %v", len(bm.buckets), 0)
	}

	if bm.newCapacity != params.Capacity {
		t.Errorf("CreateBucketMap() created incorrect newCapacity."+
			"\n\treceived: %v\n\texpected: %v", bm.newCapacity, params.Capacity)
	}

	if bm.newLeakRate != params.LeakRate {
		t.Errorf("CreateBucketMap() created incorrect newLeakRate."+
			"\n\treceived: %v\n\texpected: %v", bm.newLeakRate, params.LeakRate)
	}
}

// Ensures LookupBucket() correctly creates new buckets that were not previously
// present in the map.
func TestLookupBucket_NewBuckets(t *testing.T) {
	bm := CreateBucketMap(math.MaxUint32, math.MaxFloat64, 10*time.Second, 3*time.Second)
	buckets := make([]*Bucket, 10)

	for i := 0; i < 10; i++ {
		buckets[i] = bm.LookupBucket(string(i))

		if len(bm.buckets) != (i + 1) {
			t.Errorf("LookupBucket() has incorrect length when looking up new buckets\n\treceived: %v\n\texpected: %v", len(bm.buckets), i+1)
		}

		if buckets[i].capacity != math.MaxUint32 {
			t.Errorf("LookupBucket() set the incorrect capacity when creating a new bucket\n\treceived: %v\n\texpected: %v", buckets[i].capacity, math.MaxUint32)
		}

		if buckets[i].leakRate != math.MaxFloat64 {
			t.Errorf("LookupBucket() set the incorrect rate when creating a new bucket\n\treceived: %v\n\texpected: %v", buckets[i].leakRate, math.MaxFloat64)
		}
	}
}

// Ensures LookupBucket() correctly creates recalls buckets that already exist
// in the map
func TestLookupBucket_RecallBuckets(t *testing.T) {
	bm := CreateBucketMap(5, 0.123, 10*time.Second, 3*time.Second)
	buckets := make([]*Bucket, 10)

	for i := 0; i < 10; i++ {
		buckets[i] = bm.LookupBucket(string(i))
	}

	for i := 0; i < 9; i++ {
		bu := bm.LookupBucket(string(i))

		b, ok := bm.buckets[string(i)]

		if &bu == &buckets[i] {
			t.Errorf("LookupBucket() did not return a bucket when looking up existing buckets\n\treceived: %v\n\texpected: %v", &bu, &buckets[i])
		}

		if !ok {
			t.Errorf("LookupBucket() did not find a bucket for the given key %v when looking up existing buckets\n\treceived: %v\n\texpected: %v", i, ok, true)
		}

		if len(bm.buckets) != 10 {
			t.Errorf("LookupBucket() has incorrect length when looking up existing buckets\n\treceived: %v\n\texpected: %v", len(bm.buckets), 10)
		}

		if b.capacity != 5 {
			t.Errorf("LookupBucket() set the incorrect capacity when looking up existing buckets\n\treceived: %v\n\texpected: %v", b.capacity, 5)
		}

		if b.leakRate != 0.123 {
			t.Errorf("LookupBucket() set the incorrect rate when looking up existing buckets\n\treceived: %v\n\texpected: %v", b.leakRate, 0.123)
		}
	}
}

// Tests the thread locking functionality of LookupBucket().
func TestLookupBucket_Lock(t *testing.T) {
	bm := CreateBucketMap(0, 0, 10*time.Second, 3*time.Second)

	result := make(chan bool)

	bm.Lock()

	go func() {
		bm.LookupBucket("0")
		result <- true
	}()

	select {
	case <-result:
		t.Errorf("LookupBucket() did not correctly lock the thread")
	case <-time.After(5 * time.Second):
		return
	}
}

// Ensures clearStaleBuckets() removes stale buckets and not buckets in use.
func TestClearStaleBuckets(t *testing.T) {
	// Create a fill bucket map
	bm := CreateBucketMap(math.MaxUint32, math.MaxFloat64, 10*time.Second, 3*time.Second)
	buckets := make([]*Bucket, 20)
	for i := 0; i < 10; i++ {
		buckets[i] = bm.LookupBucket(string(i))
	}

	time.Sleep(5 * time.Second)

	for i := 10; i < 20; i++ {
		buckets[i] = bm.LookupBucket(string(i))
	}

	time.Sleep(1 * time.Second)

	bm.clearStaleBuckets()

	if len(bm.buckets) != 10 {
		t.Errorf("clearOldBuckets() did not clear out the correct number of buckets\n\treceived: %v\n\texpected: %v", len(bm.buckets), 10)
	}

	for i := 10; i < 20; i++ {
		_, exists := bm.buckets[string(i)]

		if !exists {
			t.Errorf("clearOldBuckets() cleared key (%d) that was not old enough", i)
		}
	}
}

// Ensures StaleBucketWorker() periodically runs and removes stale buckets.
func TestStaleBucketWorker(t *testing.T) {
	// Create a fill bucket map
	bm := CreateBucketMap(math.MaxUint32, math.MaxFloat64, 3*time.Second, 5*time.Second)
	buckets := make([]*Bucket, 20)
	for i := 0; i < 10; i++ {
		buckets[i] = bm.LookupBucket(string(i))
	}

	time.Sleep(5 * time.Second)

	for i := 10; i < 20; i++ {
		buckets[i] = bm.LookupBucket(string(i))
	}

	time.Sleep(4 * time.Second)

	if len(bm.buckets) != 10 {
		t.Errorf("clearOldBuckets() did not clear out the correct number of buckets\n\treceived: %v\n\texpected: %v", len(bm.buckets), 10)
	}

	for i := 10; i < 20; i++ {
		_, exists := bm.buckets[string(i)]

		if !exists {
			t.Errorf("clearOldBuckets() cleared key (%d) that was not old enough", i)
		}
	}
}
