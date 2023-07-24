////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package rateLimiting

// Storage is the generic interface used by the BucketMap for permanent storage.
type Storage interface {
	// UpsertBucket inserts the BucketParams into Storage with the unique
	// BucketParams key. If a bucket already exists with the same key, its
	// values are updated.
	UpsertBucket(bp *BucketParams)

	// AddToBucket updates the remaining and lastUpdate of the bucket with the
	// given key in Storage. If the bucket does not exist, an error is returned.
	AddToBucket(key string, remaining uint32, lastUpdate int64) error

	// RetrieveBucket returns a BucketParams for the bucket with the given key
	// from Storage. If the bucket does not exist, an error is returned.
	RetrieveBucket(key string) (*BucketParams, error)

	// RetrieveAllBuckets returns an array of all the buckets found in Storage.
	RetrieveAllBuckets() []*BucketParams

	// DeleteBucket deletes the bucket with the given key from Storage. If no
	// bucket is found, an error is returned.
	DeleteBucket(key string) error
}
