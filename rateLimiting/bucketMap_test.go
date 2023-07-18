////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package rateLimiting

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// Test implementation of the database
type bucketDB map[string]*BucketParams

func (db bucketDB) UpsertBucket(bp *BucketParams) {
	db[bp.Key] = bp
}

func (db bucketDB) AddToBucket(key string, remaining uint32, lastUpdate int64) error {
	_, exists := db[key]
	if exists {
		db[key].Remaining = remaining
		db[key].LastUpdate = lastUpdate
		return nil
	}
	return errors.Errorf("Failed to find bucket with key %s.", key)
}

func (db bucketDB) RetrieveBucket(key string) (*BucketParams, error) {
	bp, exists := db[key]
	if exists {
		return bp, nil
	}

	return nil, errors.Errorf("Failed to find bucket with key %s.", key)
}

func (db bucketDB) RetrieveAllBuckets() []*BucketParams {
	var params []*BucketParams
	for _, bp := range db {
		params = append(params, bp)
	}
	return params
}

func (db bucketDB) DeleteBucket(key string) error {
	_, exists := db[key]
	if exists {
		delete(db, key)
		return nil
	}
	return errors.Errorf("Failed to find bucket in map.")
}

// Test that CreateBucketMap produces a new BucketMap with all the correct
// values.
func TestCreateBucketMap(t *testing.T) {
	expectedBM := BucketMap{
		buckets:      make(map[string]*Bucket),
		capacity:     5,
		leakRate:     0.000003,
		pollDuration: 10 * time.Second,
		bucketMaxAge: 3 * time.Second,
	}

	bm := CreateBucketMap(expectedBM.capacity, 3, time.Millisecond,
		expectedBM.pollDuration, expectedBM.bucketMaxAge, nil, nil)

	if !reflect.DeepEqual(&expectedBM, bm) {
		t.Errorf("CreateBucketMap returned an incorrect BucketMap."+
			"\nexpected: %+v\nreceived: %+v", &expectedBM, bm)
	}
}

// Test that CreateBucketMap loads the buckets from the database when it is
// supplied.
func TestCreateBucketMap_DB(t *testing.T) {
	expectedBM := BucketMap{
		buckets:      make(map[string]*Bucket),
		capacity:     5,
		leakRate:     0.000003,
		pollDuration: 10 * time.Second,
		bucketMaxAge: 3 * time.Second,
	}

	// Test buckets
	db := bucketDB{
		"keyA": {"keyA", 32, 24, 50.832, 6337, true, false},
		"keyB": {"keyB", 70, 50, 84.511, 1798, true, true},
		"keyC": {"keyC", 12, 21, 18.631, 9050, false, false},
		"keyD": {"keyD", 37, 31, 84.077, 1468, false, true},
		"keyE": {"keyE", 19, 26, 39.331, 5167, true, true},
		"keyF": {"keyF", 20, 12, 89.203, 4294, true, false},
		"keyG": {"keyG", 56, 24, 10.622, 6494, true, true},
	}

	bm := CreateBucketMap(expectedBM.capacity, 3, time.Millisecond,
		expectedBM.pollDuration, expectedBM.bucketMaxAge, db, nil)

	for key, bp := range db {
		b, exists := bm.buckets[key]
		if !exists {
			t.Errorf("CreateBucketMap did not load bucket %s from database.", key)
		} else if b.capacity != bp.Capacity || b.remaining != bp.Remaining ||
			b.leakRate != bp.LeakRate || b.lastUpdate != bp.LastUpdate ||
			b.locked != bp.Locked || b.whitelist != bp.Whitelist {
			expectedBucket := CreateBucketFromParams(bp, bm.createAddToDbFunc(key))
			t.Errorf("CreateBucketMap did not load the correct values from "+
				"the database for bucket %s.\nexpected: %#v\nreceived: %#v",
				key, expectedBucket, b)
		}
	}
}

// Tests that TestCreateBucketMapFromParams creates a new BucketMap with
// values matching the input MapParams.
func TestCreateBucketMapFromParams(t *testing.T) {
	expectedLeakRate := 0.000003
	testParams := &MapParams{
		Capacity:     5,
		LeakedTokens: 3,
		LeakDuration: time.Millisecond,
		PollDuration: 10 * time.Second,
		BucketMaxAge: 3 * time.Second,
	}
	bm := CreateBucketMapFromParams(testParams, nil, nil)

	if bm.capacity != testParams.Capacity {
		t.Errorf("CreateBucketMap returned incorrect capacity."+
			"\nexpected: %d\nreceived: %d", testParams.Capacity, bm.capacity)
	}

	if bm.leakRate != expectedLeakRate {
		t.Errorf("CreateBucketMap returned incorrect leakRate."+
			"\nexpected: %f\nreceived: %f", expectedLeakRate, bm.leakRate)
	}

	if bm.pollDuration != testParams.PollDuration {
		t.Errorf("CreateBucketMap returned incorrect pollDuration."+
			"\nexpected: %s\nreceived: %s",
			testParams.PollDuration, bm.pollDuration)
	}

	if bm.bucketMaxAge != testParams.BucketMaxAge {
		t.Errorf("CreateBucketMap returned incorrect bucketMaxAge."+
			"\nexpected: %s\nreceived: %s",
			testParams.BucketMaxAge, bm.bucketMaxAge)
	}
}

// Tests that BucketMap.LookupBucket returns the correct bucket or creates it
// if it does not exist. This is done by adding tokens to buckets that are added
// and checking if the number of tokens match when looking up the bucket.
//
// This test is run twice. The first time no database is used and the second
// time a database is used and tests that the bucket is added to the database
// and that the number of tokens is correct.
func TestBucketMap_LookupBucket(t *testing.T) {
	// Generate test buckets
	testData := []struct {
		key       string
		tokens    uint32
		addTokens bool
	}{
		{"keyA", rand.Uint32(), true},
		{"keyB", rand.Uint32(), true},
		{"keyC", 0, false},
		{"keyD", 0, false},
		{"keyE", rand.Uint32(), true},
		{"keyF", 0, false},
	}

	// Test buckets in two rounds: Round 1 will test with no database and round
	// two tests with a database
	for i := 0; i < 2; i++ {
		// Build BucketMap
		bm := &BucketMap{}
		db := bucketDB{}
		if i == 1 {
			t.Logf("Running with database.")
			bm = CreateBucketMap(10, 3, 0, 0, 0, db, nil)
		} else {
			t.Logf("Running without database.")
			bm = CreateBucketMap(10, 3, 0, 0, 0, nil, nil)
		}

		// Add some of the buckets to the map
		for _, b := range testData {
			if b.addTokens {
				bucket := bm.LookupBucket(b.key)
				bucket.Add(b.tokens)
			}
		}

		// Lookup all the buckets and see if the added buckets still exist if the
		// new buckets were added
		for _, b := range testData {
			bucket := bm.LookupBucket(b.key)
			if bucket.Remaining() != b.tokens {
				t.Errorf("LookupBucket returned incorrect bucket %s, it has "+
					"incorrect number of tokens.\nexpected: %d\nreceived: %d",
					b.key, b.tokens, bucket.Remaining())
			}
			if bm.db != nil {
				bp, exists := db[b.key]
				if !exists {
					t.Errorf("LookupBucket failed to add the bucket %s to the "+
						"database.", b.key)
				}
				if bp.Remaining != bucket.Remaining() {
					t.Errorf("Bucket %s in the database has incorrect number "+
						"of tokens.\nexpected: %d\nreceived: %d",
						b.key, bucket.Remaining(), bp.Remaining)
				}
			}
		}
	}
}

// Tests that BucketMap.LookupBucket can look up an existing bucket while a
// read lock is enabled.
func TestBucketMap_LookupBucket_ReadLockWhileReading(t *testing.T) {
	bm := CreateBucketMap(10, 3, 0, 0, 0, nil, nil)
	result := make(chan bool)

	// Add the test bucket
	bm.LookupBucket("test")

	bm.RLock()
	defer bm.RUnlock()

	go func() {
		bm.LookupBucket("test")
		result <- true
	}()

	select {
	case <-result:
		return
	case <-time.After(50 * time.Millisecond):
		t.Errorf("LookupBucket stalled when it should have succesfully " +
			"looked up an existing bucket.")
	}
}

// Tests that BucketMap.LookupBucket cannot look up an existing bucket while
// there is a write lock.
func TestBucketMap_LookupBucket_WriteLockWhileReading(t *testing.T) {
	bm := CreateBucketMap(10, 3, 0, 0, 0, nil, nil)
	result := make(chan bool)

	// Add the test bucket
	bm.LookupBucket("test")

	bm.Lock()
	defer bm.Unlock()

	go func() {
		bm.LookupBucket("test")
		result <- true
	}()

	select {
	case <-result:
		t.Errorf("LookupBucket completed when it should have been waiting " +
			"for the Lock to release.")
	case <-time.After(50 * time.Millisecond):
		return
	}
}

// Tests that BucketMap.LookupBucket cannot add a new bucket to the map while a
// read lock is enabled.
func TestBucketMap_LookupBucket_ReadLockWhileWriting(t *testing.T) {
	bm := CreateBucketMap(10, 3, 0, 0, 0, nil, nil)
	result := make(chan bool)

	bm.RLock()
	defer bm.RUnlock()

	go func() {
		bm.LookupBucket("test")
		result <- true
	}()

	select {
	case <-result:
		t.Errorf("LookupBucket completed when it should have been waiting " +
			"for the RLock to release.")
	case <-time.After(50 * time.Millisecond):
		return
	}
}

// Tests that BucketMap.AddBucket correctly adds new buckets and overwrites
// existing buckets. This function runs twice, once without a database and a
// second time with a database.
func TestBucketMap_AddBucket(t *testing.T) {
	// Generate test buckets
	testData := []struct {
		key          string
		tokensAdded  uint32
		capacity     uint32
		remaining    uint32
		addInitially bool
	}{
		{"keyA", 37, 20, 17, true},
		{"keyB", 0, 25, 0, false},
		{"keyC", 173, 15, 8, true},
		{"keyD", 0, 1, 0, false},
	}

	// Test buckets in two rounds: Round 1 will test with no database and round
	// two tests with a database
	for i := 0; i < 2; i++ {
		// Build BucketMap
		bm := &BucketMap{}
		db := bucketDB{}
		if i == 1 {
			t.Logf("Running with database.")
			bm = CreateBucketMap(10, 3, 0, 0, 0, db, nil)
		} else {
			t.Logf("Running without database.")
			bm = CreateBucketMap(10, 3, 0, 0, 0, nil, nil)
		}

		// Add some of the predefined buckets to the map
		for _, b := range testData {
			if b.addInitially {
				bucket := bm.LookupBucket(b.key)
				bucket.Add(b.tokensAdded)
			}
		}

		// Add all the buckets in the map and ensure they exist and overwrote
		// any other buckets that were in the map.
		for _, b := range testData {
			bm.AddBucket(b.key, b.capacity, 30, time.Second)
			bucket, exists := bm.buckets[b.key]
			if !exists {
				t.Errorf("AddBucket did not add the bucket for key %s.", b.key)
			}
			if bucket.capacity != b.capacity {
				t.Errorf("AddBucket did not overwrite the exisitng bucket."+
					"\ncapacity expected: %d\ncapacity received: %d",
					b.capacity, bucket.capacity)
			}
			if bucket.remaining != b.remaining {
				t.Errorf("AddBucket did not overwrite the exisitng bucket."+
					"\nremaining expected: %d\nremaining received: %d",
					b.remaining, bucket.remaining)
			}
			if bm.db != nil {
				bp, exists := db[b.key]
				if !exists {
					t.Errorf("AddBucket failed to add the bucket %s to the "+
						"database.", b.key)
				}
				if bp.Capacity != bucket.Capacity() {
					t.Errorf("Bucket %s in the database has incorrect number "+
						"of tokens.\nexpected: %d\nreceived: %d",
						b.key, bucket.Capacity(), bp.Capacity)
				}
				if bp.Remaining != bucket.Remaining() {
					t.Errorf("Bucket %s in the database has incorrect number "+
						"of tokens.\nexpected: %d\nreceived: %d",
						b.key, bucket.Remaining(), bp.Remaining)
				}
			}
		}
	}
}

// Tests that BucketMap.AddBucket cannot add a new bucket to the map while a
// read lock is enabled.
func TestBucketMap_AddBucket_ReadLock(t *testing.T) {
	bm := CreateBucketMap(10, 3, 0, 0, 0, nil, nil)
	result := make(chan bool)

	bm.RLock()
	defer bm.RUnlock()

	go func() {
		bm.AddBucket("test", 0, 0, 0)
		result <- true
	}()

	select {
	case <-result:
		t.Errorf("AddBucket completed when it should have been waiting " +
			"for the RLock to release.")
	case <-time.After(50 * time.Millisecond):
		return
	}
}

// Tests that BucketMap.addAllBuckets adds all the buckets with the correct
// values.
func TestBucketMap_AddAllBuckets(t *testing.T) {
	// Generate test buckets
	testBP := []*BucketParams{
		{"keyA", 32, 24, 50.832, 6337, true, false},
		{"keyB", 70, 50, 84.511, 1798, true, true},
		{"keyC", 12, 21, 18.631, 9050, false, false},
		{"keyD", 37, 31, 84.077, 1468, false, true},
		{"keyE", 19, 26, 39.331, 5167, true, true},
		{"keyF", 20, 12, 89.203, 4294, true, false},
		{"keyG", 56, 24, 10.622, 6494, true, true},
		{"", 56, 24, 10.622, 6494, true, true},
	}

	// Create new BucketMap
	bm := CreateBucketMap(5, 3, 0, 0, 0, nil, nil)

	// Add the buckets
	bm.addAllBuckets(testBP)

	// Check that all the buckets exist
	for _, bp := range testBP {
		b, exists := bm.buckets[bp.Key]
		if !exists {
			t.Errorf(
				"addAllBuckets failed to add the bucket with key %s.", bp.Key)
		} else if !reflect.DeepEqual(b, CreateBucketFromParams(bp, nil)) {
			t.Errorf("addAllBuckets created bucket %s with incorrect values."+
				"\nexpected: %+v\nreceived: %+v", bp.Key,
				CreateBucketFromParams(bp, nil), b)
		}
	}
}

// Tests that BucketMap.AddToWhitelist correctly adds the whitelisted items or
// modifies existing buckets to be whitelisted, if they are already in the map.
// This function runs twice, once without a database and a second time with a
// database.
func TestBucketMap_AddToWhitelist(t *testing.T) {
	testKeys := []string{"keyA", "keyB", "keyC", "keyD", ""}

	// Test buckets in two rounds: Round 1 will test with no database and round
	// two tests with a database
	for i := 0; i < 2; i++ {
		// Build BucketMap
		bm := &BucketMap{}
		db := bucketDB{}
		if i == 1 {
			t.Logf("Running with database.")
			bm = CreateBucketMap(10, 3, 0, 0, 0, db, nil)
		} else {
			t.Logf("Running without database.")
			bm = CreateBucketMap(10, 3, 0, 0, 0, nil, nil)
		}

		// Insert buckets into map
		for _, key := range testKeys {
			bm.LookupBucket(key)
		}

		wlKeys := []string{"keyA", "keyC", "keyE", "keyF", ""}
		bm.AddToWhitelist(wlKeys)

		// Check if buckets are on the whitelist
		for _, key := range wlKeys {
			b, exists := bm.buckets[key]
			if !exists {
				t.Errorf(
					"AddToWhitelist did not add the key %s to the map.", key)
			}
			if b.whitelist != true {
				t.Errorf(
					"AddToWhitelist did not mark the key %s as whitelisted.", key)
			}

			if bm.db != nil {
				bp, exists := db[key]
				if !exists {
					t.Errorf("AddToWhitelist failed to add the bucket %s to "+
						"the database.", key)
				}
				if bp.Whitelist != b.whitelist {
					t.Errorf(
						"Bucket %s in the database is not on whitelist.", key)
				}
			}
		}
	}
}

// Tests that BucketMap.AddToWhitelist cannot mark buckets as whitelisted
// while a read lock is enabled.
func TestBucketMap_AddToWhitelist_ReadLock(t *testing.T) {
	bm := CreateBucketMap(10, 3, 0, 0, 0, nil, nil)
	result := make(chan bool)

	bm.RLock()
	defer bm.RUnlock()

	go func() {
		bm.AddToWhitelist([]string{})
		result <- true
	}()

	select {
	case <-result:
		t.Errorf("AddToWhitelist completed when it should have been waiting " +
			"for the RLock to release.")
	case <-time.After(50 * time.Millisecond):
		return
	}
}

// Tests that BucketMap.DeleteBucket deletes the correct buckets. This
// function runs twice, once without a database and a second time with a
// database.
func TestBucketMap_DeleteBucket(t *testing.T) {
	// Generate test buckets
	testBP := []*BucketParams{
		{"keyA", 32, 24, 50.832, 6337, true, false},
		{"keyB", 70, 50, 84.511, 1798, true, true},
		{"keyC", 12, 21, 18.631, 9050, false, false},
		{"keyD", 37, 31, 84.077, 1468, false, true},
		{"keyE", 19, 26, 39.331, 5167, true, true},
		{"keyF", 20, 12, 89.203, 4294, true, false},
		{"keyG", 56, 24, 10.622, 6494, true, true},
	}

	// Add buckets to database
	db := bucketDB{}
	for _, bp := range testBP {
		db[bp.Key] = bp
	}

	// Test buckets in two rounds: Round 1 will test with no database and round
	// two tests with a database
	for i := 0; i < 2; i++ {
		// Build BucketMap
		bm := &BucketMap{}
		if i == 1 {
			t.Logf("Running with database.")
			bm = CreateBucketMap(10, 3, 0, 0, 0, db, nil)
		} else {
			t.Logf("Running without database.")
			bm = CreateBucketMap(10, 3, 0, 0, 0, nil, nil)
		}

		// Insert all the buckets into the map
		for _, bp := range testBP {
			bm.buckets[bp.Key] = CreateBucketFromParams(bp, nil)
		}

		// Delete half the buckets and see if they exist
		var expectedLength int
		for i, bp := range testBP {
			if i%2 == 1 {
				err := bm.DeleteBucket(bp.Key)
				if err != nil {
					t.Errorf("DeleteBucket returned an error when deleting key %s."+
						"\nexpected: %v\nreceived: %v", bp.Key, nil, err)
				}

				bucket, exists := bm.buckets[bp.Key]
				if exists {
					t.Errorf("DeleteBucket did not delete the bucket with key %s."+
						"\nbucket: %+v", bp.Key, bucket)
				}

				if bm.db != nil {
					_, exists := db[bp.Key]
					if exists {
						t.Errorf("DeleteBucket did not delete the bucket "+
							"with key %s from the database.", bp.Key)
					}
				}

			} else {
				_, exists := bm.buckets[bp.Key]
				if !exists {
					t.Errorf("DeleteBucket deleted the bucket with key %s."+
						"\nbucket params: %+v", bp.Key, bp)
				}

				if bm.db != nil {
					_, exists := db[bp.Key]
					if !exists {
						t.Errorf("DeleteBucket deleted bucket with key %s "+
							"from the database.", bp.Key)
					}
				}

				expectedLength++
			}
		}

		if len(bm.buckets) != expectedLength {
			t.Errorf("DeleteBucket did not delete all the correct buckets."+
				"\nexpected length: %d\nreceived length: %d\nbucket map: %+v",
				expectedLength, len(bm.buckets), bm.buckets)
		}
	}
}

// Tests that BucketMap.DeleteBucket returns an error when the bucket does not
// exist.
func TestBucketMap_DeleteBucket_Error(t *testing.T) {
	bm := CreateBucketMap(5, 3, 0, 0, 0, nil, nil)

	err := bm.DeleteBucket("test")
	if err == nil {
		t.Errorf("DeleteBucket did not return an error when the bucket " +
			"does not exist.")
	}
}

// Tests that BucketMap.DeleteBucket cannot delete buckets from the map while
// a read lock is enabled.
func TestBucketMap_DeleteBucket_ReadLock(t *testing.T) {
	bm := CreateBucketMap(5, 3, 0, 0, 0, nil, nil)
	result := make(chan bool)

	bm.RLock()
	defer bm.RUnlock()

	go func() {
		_ = bm.DeleteBucket("test")
		result <- true
	}()

	select {
	case <-result:
		t.Errorf("DeleteBucket completed when it should have been waiting " +
			"for the RLock to release.")
	case <-time.After(50 * time.Millisecond):
		return
	}
}

// Tests that BucketMap.StaleBucketWorker deletes stale buckets in a separate
// thread and that the quit channel stops it.
func TestBucketMap_StaleBucketWorker(t *testing.T) {
	quit := make(chan struct{})
	bm := CreateBucketMap(
		5, 3, time.Millisecond, 30*time.Millisecond, time.Second, nil, quit)
	bm.buckets = map[string]*Bucket{
		"keyA": CreateBucketFromParams( // Stale
			&BucketParams{"keyA", 10, 0, 1,
				time.Now().Add(-3 * time.Second).UnixNano(), false, false}, nil),
		"keyB": CreateBucketFromParams( // Stale
			&BucketParams{"keyB", 10, 0, 1,
				time.Now().AddDate(0, -1, -13).UnixNano(), false, false}, nil),
	}

	time.Sleep(2 * bm.pollDuration)

	if len(bm.buckets) != 0 {
		t.Errorf("staleBucketWorker did not delete the stale buckets."+
			"\nexpected length: %d\nreceived length: %d", 0, len(bm.buckets))
	}

	quit <- struct{}{}

	bm.buckets = map[string]*Bucket{
		"keyA": CreateBucketFromParams(&BucketParams{"keyA", 10, 6, 1, // Stale
			time.Now().Add(-3 * time.Second).UnixNano(), false, false}, nil),
		"keyB": CreateBucketFromParams(&BucketParams{"keyB", 10, 6, 1, // Stale
			time.Now().AddDate(0, -1, -13).UnixNano(), false, false}, nil),
	}

	time.Sleep(2 * bm.pollDuration)

	if len(bm.buckets) != 2 {
		t.Errorf("staleBucketWorker deleted the buckets when it should have " +
			"been stopped.")
	}

}

// Tests that BucketMap.clearStaleBuckets removes all the stale buckets from
// the map. This function runs twice, once without a database and a second time
// with a database.
func TestBucketMap_ClearStaleBuckets(t *testing.T) {
	// Generate test data and bucket parameters
	now := time.Now().UnixNano()
	testData := []struct {
		stale bool
		p     *BucketParams
	}{
		{false, &BucketParams{"keyA", 10, 0, 1, now, false, false}},                                 // Not stale
		{true, &BucketParams{"keyB", 10, 0, 1, now - 3*time.Second.Nanoseconds(), false, false}},    // Stale
		{false, &BucketParams{"keyC", 10, 7, 1, now - 6*time.Second.Nanoseconds(), true, false}},    // Stale but locked
		{false, &BucketParams{"keyD", 10, 7, 1, now, true, false}},                                  // Not stale but locked
		{true, &BucketParams{"keyE", 10, 0, 1, now - 50*time.Hour.Nanoseconds(), false, false}},     // Stale
		{false, &BucketParams{"keyF", 10, 100, 1, now - 3*time.Second.Nanoseconds(), false, false}}, // Not stale
	}

	// Add buckets to database
	db := bucketDB{}
	for _, bp := range testData {
		db[bp.p.Key] = bp.p
	}

	// Test buckets in two rounds: Round 1 will test with no database and round
	// two tests with a database
	for i := 0; i < 2; i++ {
		// Build BucketMap
		bm := &BucketMap{}
		if i == 1 {
			t.Logf("Running with database.")
			bm = CreateBucketMap(5, 3, 0, 0, time.Second, db, nil)
		} else {
			t.Logf("Running without database.")
			bm = CreateBucketMap(5, 3, 0, 0, time.Second, nil, nil)
		}

		// Create new BucketMap and add buckets to it
		for _, bp := range testData {
			bm.buckets[bp.p.Key] = CreateBucketFromParams(bp.p, nil)
		}

		// Clear stale buckets
		bm.clearStaleBuckets()

		// Check if buckets marked stale are removed and non-stale bucket are kept
		for _, bp := range testData {
			b, exists := bm.buckets[bp.p.Key]
			if !bp.stale && !exists {
				t.Errorf("clearStaleBuckets deleted non-stale bucket %s.",
					bp.p.Key)
			} else if bp.stale && exists {
				t.Errorf("clearStaleBuckets did not delete stale bucket %s."+
					"\nbucket: %+v", bp.p.Key, b)
			}

			if bm.db != nil {
				_, exists := db[bp.p.Key]
				if !bp.stale && !exists {
					t.Errorf("clearStaleBuckets deleted non-stale bucket %s "+
						"from the database.", bp.p.Key)
				} else if bp.stale && exists {
					t.Errorf("clearStaleBuckets did not delete stale bucket "+
						"%s from the database.\nbucket: %+v", bp.p.Key, b)
				}
			}
		}
	}
}

// Tests that BucketMap.clearStaleBuckets can search for stale buckets during
// a read lock.
func TestBucketMap_ClearStaleBuckets_ReadLockWhileReading(t *testing.T) {
	// Create new BucketMap and add a non-stale bucket to it
	bm := CreateBucketMap(5, 3, 0, 0, time.Second, nil, nil)
	bm.buckets["keyA"] = CreateBucketFromParams(&BucketParams{"keyA", 10, 0, 1,
		time.Now().UnixNano(), false, false}, nil)

	result := make(chan bool)

	bm.RLock()
	defer bm.RUnlock()

	go func() {
		bm.clearStaleBuckets()
		result <- true
	}()

	select {
	case <-result:
		return
	case <-time.After(50 * time.Millisecond):
		t.Errorf("clearStaleBuckets stalled when it should have succesfully " +
			"finished without finding any stale buckets.")
	}
}

// Tests that BucketMap.clearStaleBuckets cannot search for stale buckets
// during a write lock.
func TestBucketMap_ClearStaleBuckets_WriteLockWhileReading(t *testing.T) {
	// Create new BucketMap and add a non-stale bucket to it
	bm := CreateBucketMap(5, 3, 0, 0, time.Second, nil, nil)
	bm.buckets["keyA"] = CreateBucketFromParams(&BucketParams{"keyA", 10, 0, 1,
		time.Now().UnixNano(), false, false}, nil)

	result := make(chan bool)

	bm.Lock()
	defer bm.Unlock()

	go func() {
		bm.clearStaleBuckets()
		result <- true
	}()

	select {
	case <-result:
		t.Errorf("clearStaleBuckets completed when it should have been waiting " +
			"for the Lock to release.")
	case <-time.After(50 * time.Millisecond):
		return
	}
}

// Tests that BucketMap.clearStaleBuckets cannot remove stale buckets from the
// map while a read lock is enabled.
func TestBucketMap_ClearStaleBuckets_ReadLockWhileWriting(t *testing.T) {
	// Create new BucketMap and add a stale bucket to it
	bm := CreateBucketMap(5, 3, 0, 0, time.Second, nil, nil)
	bm.buckets["keyA"] = CreateBucketFromParams(&BucketParams{"keyA", 10, 0, 1,
		time.Now().Add(-3 * time.Second).UnixNano(), false, false}, nil)

	result := make(chan bool)

	bm.RLock()
	defer bm.RUnlock()

	go func() {
		bm.clearStaleBuckets()
		result <- true
	}()

	select {
	case <-result:
		t.Errorf("clearStaleBuckets completed when it should have been waiting " +
			"for the RLock to release.")
	case <-time.After(50 * time.Millisecond):
		return
	}
}

// Tests that BucketMap.createAddToDbFunc generates an anonymous function that
// panics when it attempts to operate on a bucket that does not exist in the
// database.
func TestBucketMap_CreateAddToDbFunc_Panic(t *testing.T) {
	bm := CreateBucketMap(5, 3, 0, 0, time.Second, bucketDB{}, nil)
	addFunc := bm.createAddToDbFunc("keyA")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("createAddToDbFunc produced a function that does not " +
				"panic when the bucket does not exist in the database.")
		}
	}()

	addFunc(rand.Uint32(), rand.Int63())
}
