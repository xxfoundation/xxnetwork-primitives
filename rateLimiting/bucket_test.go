package rateLimiting

import (
	"math/rand"
	"testing"
	"time"
)

// Test that Create() generates a new bucket with all the expected fields.
func TestCreate(t *testing.T) {
	// Setup expected values
	expectedCapacity := uint(rand.Uint32())
	expectedLeakRate := rand.Float64()

	// Max time (in nanoseconds) between update and test
	testUpdateWait := int64(1000)

	// Create new Bucket
	b := Create(expectedCapacity, expectedLeakRate)

	// Test fields for expected results
	if b.capacity != expectedCapacity {
		t.Errorf("Create() generated Bucket with incorrect capacity."+
			"\n\texpected: %v\n\treceived: %v", expectedCapacity, b.capacity)
	}

	if b.remaining != 0 {
		t.Errorf("Create() generated Bucket with incorrect remaining."+
			"\n\texpected: %v\n\treceived: %v", 0, b.remaining)
	}

	if b.leakRate != expectedLeakRate {
		t.Errorf("Create() generated Bucket with incorrect rate."+
			"\n\texpected: %v\n\treceived: %v", expectedLeakRate, b.leakRate)
	}

	// Check that the lastUpdate occurred recently
	if time.Since(b.lastUpdate).Nanoseconds() > testUpdateWait {
		t.Errorf("Create() generated Bucket with incorrect lastUpdate."+
			"The time between creation and testing was greater than %v ns"+
			"\n\texpected: ~ %v\n\treceived:   %v",
			testUpdateWait, time.Now(), b.lastUpdate)
	}
}

// Tests that Capacity() returns the correct value for a new Bucket.
func TestCapacity(t *testing.T) {
	// Setup expected values
	expectedCapacity := uint(rand.Uint32())

	// Create new Bucket
	b := Create(expectedCapacity, rand.Float64())

	if b.Capacity() != expectedCapacity {
		t.Errorf("Capacity() returned incorrect capacity for bucket."+
			"\n\texpected: %v\n\treceived: %v", expectedCapacity, b.Capacity())
	}
}

// Tests that Remaining() returns the correct value for a new Bucket.
func TestRemaining(t *testing.T) {
	// Create new Bucket
	b := Create(uint(rand.Uint32()), rand.Float64())

	if b.Remaining() != 0 {
		t.Errorf("Remaining() returned incorrect remaining for the bucket."+
			"\n\texpected: %v\n\treceived: %v", 0, b.Remaining())
	}
}

func TestAdd_UnderLeakRate(t *testing.T) {
	/*	// Setup expected values
		expectedCapacity := uint(rand.Uint32())
		expectedLeakRate := rand.Float64()
		b := Create(expectedCapacity, expectedLeakRate)*/
	b := Create(10, 0.000000003) // 3 per second

	addReturnVal := b.Add(9)

	if addReturnVal != true {
		t.Errorf("Add() failed to add when adding under the leak rate\n\treceived: %v\n\texpected: %v", addReturnVal, true)
	}

	if b.remaining != 9 {
		t.Errorf("Add() returned incorrect remaining when adding under the leak rate\n\treceived: %v\n\texpected: %v", b.remaining, 9)
	}

	time.Sleep(2 * time.Second)
	addReturnVal = b.Add(7)

	if addReturnVal != true {
		t.Errorf("Add() failed to add when adding under the leak rate\n\treceived: %v\n\texpected: %v", addReturnVal, true)
	}

	if b.remaining != 10 {
		t.Errorf("Add() returned incorrect remaining when adding under the leak rate\n\treceived: %v\n\texpected: %v", b.remaining, 10)
	}

	time.Sleep(4 * time.Second)
	addReturnVal = b.Add(2)

	if addReturnVal != true {
		t.Errorf("Add() failed to add when adding under the leak rate\n\treceived: %v\n\texpected: %v", addReturnVal, true)
	}

	if b.remaining != 2 {
		t.Errorf("Add() returned incorrect remaining when adding under the leak rate\n\treceived: %v\n\texpected: %v", b.remaining, 2)
	}

	addReturnVal = b.Add(6)

	if addReturnVal != true {
		t.Errorf("Add() failed to add when adding under the leak rate\n\treceived: %v\n\texpected: %v", addReturnVal, true)
	}

	if b.remaining != 8 {
		t.Errorf("Add() returned incorrect remaining when adding under the leak rate\n\treceived: %v\n\texpected: %v", b.remaining, 8)
	}
}

func TestAdd_OverLeakRate(t *testing.T) {
	b := Create(10, 0.000000003) // 3 per second

	addReturnVal := b.Add(17)

	if addReturnVal != false {
		t.Errorf("Add() incorrectly added when adding over the leak rate\n\treceived: %v\n\texpected: %v", addReturnVal, false)
	}

	if b.remaining != 17 {
		t.Errorf("Add() returned incorrect remaining when adding over the leak rate\n\treceived: %v\n\texpected: %v", b.remaining, 17)
	}

	time.Sleep(6 * time.Second)
	addReturnVal = b.Add(712)

	if addReturnVal != false {
		t.Errorf("Add() incorrectly added when adding over the leak rate\n\treceived: %v\n\texpected: %v", addReturnVal, false)
	}

	if b.remaining != 712 {
		t.Errorf("Add() returned incorrect remaining when adding over the leak rate\n\treceived: %v\n\texpected: %v", b.remaining, 712)
	}

	addReturnVal = b.Add(85)

	if addReturnVal != false {
		t.Errorf("Add() incorrectly added when adding over the leak rate\n\treceived: %v\n\texpected: %v", addReturnVal, false)
	}

	if b.remaining != 797 {
		t.Errorf("Add() returned incorrect remaining when adding over the leak rate\n\treceived: %v\n\texpected: %v", b.remaining, 797)
	}
}

func TestAdd_ReturnToNormal(t *testing.T) {
	b := Create(10, 0.000000003) // 3 per second

	addReturnVal := b.Add(12)

	if addReturnVal != false {
		t.Errorf("Add() incorrectly added\n\treceived: %v\n\texpected: %v", addReturnVal, false)
	}

	if b.remaining != 12 {
		t.Errorf("Add() returned incorrect remaining\n\treceived: %v\n\texpected: %v", b.remaining, 12)
	}

	time.Sleep(2 * time.Second)
	addReturnVal = b.Add(2)

	if addReturnVal != true {
		t.Errorf("Add() failed to add\n\treceived: %v\n\texpected: %v", addReturnVal, true)
	}

	if b.remaining != 8 {
		t.Errorf("Add() returned incorrect remaining\n\treceived: %v\n\texpected: %v", b.remaining, 8)
	}

	time.Sleep(1 * time.Second)
	addReturnVal = b.Add(5)

	if addReturnVal != true {
		t.Errorf("Add() failed to add\n\treceived: %v\n\texpected: %v", addReturnVal, true)
	}

	if b.remaining != 10 {
		t.Errorf("Add() returned incorrect remaining\n\treceived: %v\n\texpected: %v", b.remaining, 10)
	}

	time.Sleep(3 * time.Second)

	addReturnVal = b.Add(1)
	addReturnVal = b.Add(1)
	addReturnVal = b.Add(1)
	addReturnVal = b.Add(1)
	addReturnVal = b.Add(1)
	addReturnVal = b.Add(1)
	addReturnVal = b.Add(1)
	addReturnVal = b.Add(1)

	if addReturnVal != true {
		t.Errorf("Add() failed to add\n\treceived: %v\n\texpected: %v", addReturnVal, true)
	}

	if b.remaining != 9 {
		t.Errorf("Add() returned incorrect remaining\n\treceived: %v\n\texpected: %v", b.remaining, 9)
	}

	addReturnVal = b.Add(7)

	if addReturnVal != false {
		t.Errorf("Add() failed to add\n\treceived: %v\n\texpected: %v", addReturnVal, false)
	}

	if b.remaining != 16 {
		t.Errorf("Add() returned incorrect remaining\n\treceived: %v\n\texpected: %v", b.remaining, 16)
	}

	time.Sleep(3 * time.Second)
	addReturnVal = b.Add(7)

	if addReturnVal != false {
		t.Errorf("Add() incorrectly added\n\treceived: %v\n\texpected: %v", addReturnVal, false)
	}

	if b.remaining != 14 {
		t.Errorf("Add() returned incorrect remaining\n\treceived: %v\n\texpected: %v", b.remaining, 14)
	}
}

func TestAddLock(t *testing.T) {
	b := Create(10, 0.000000001)

	result := make(chan bool)

	b.Lock()

	go func() {
		b.Add(15)
		result <- true
	}()

	select {
	case <-result:
		t.Errorf("Add() did not correctly lock the thread")
	case <-time.After(5 * time.Second):
		return
	}
}
