////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package ring

import (
	"fmt"
	"strings"
	"testing"
)

// Setup func for tests
func setup() *Buff {
	rb := NewBuff(5)
	for i := 0; i < 5; i++ {
		v := i
		rb.Push(&v)
	}
	return rb
}

func TestBuff_GetNewestId(t *testing.T) {
	rb := NewBuff(5)
	id := rb.GetNewestId()
	if id != -1 {
		t.Error("Should have returned -1 for empty buff")
	}

	rb = setup()
	_ = rb.UpsertById(7, &struct{}{})
	id = rb.GetNewestId()
	if id != 7 {
		t.Error("Should have returned last pushed id")
	}
}

func TestBuff_GetOldestId(t *testing.T) {
	rb := NewBuff(5)
	id := rb.GetOldestId()
	if id != 0 {
		t.Error("Should have returned 0 for empty buff")
	}

	rb = setup()
	//rb.Push(&struct{}{})
	id = rb.GetOldestId()
	fmt.Println("new", rb.GetNewestId())
	if id != 0 {
		t.Errorf("Should have returned 0, instead got: %d", id)
	}

	_ = rb.UpsertById(22, &struct{}{})
	id = rb.GetOldestId()
	if id != 18 {
		t.Errorf("Should have returned 18, instead got %d", id)
	}
}

// Test the Get function on ringbuff
func TestBuff_Get(t *testing.T) {
	rb := setup()
	val := rb.Get().(*int)
	if *val != 4 {
		t.Errorf("Did not get most recent ID")
	}
}

// Test the GetById function of ringbuff
func TestBuff_GetById(t *testing.T) {
	rb := setup()
	val, err := rb.GetById(3)
	if err != nil {
		t.Errorf("Failed to get id 3: %s", err.Error())
	}
	v := *val.(*int)
	if v != 3 {
		t.Errorf("Got the wrong id: expected: %v, Recieved: %v", 3, v)
	}
}

// Test the basic push function on RingBuff
func TestBuff_Push(t *testing.T) {
	rb := setup()
	oldFirst := rb.oldest
	v := 6
	rb.Push(&v)
	if rb.oldest != oldFirst+1 {
		t.Error("Didn't increment old properly")
	}
	val := rb.Get().(*int)
	if *val != v {
		t.Error("Did not get newest id")
	}
}

// Test ID upsert on ringbuff (bulk of cases)
func TestBuff_UpsertById(t *testing.T) {
	rb := setup()
	v := 15
	err := rb.UpsertById(v, &v)
	if err != nil {
		t.Errorf("Error on initial upsert: %+v", err)
	}
	if *rb.Get().(*int) != v {
		t.Error("Failed to get correct ID")
	}

	val, _ := rb.GetById(7)
	if val != nil {
		t.Errorf("Should have gotten nil value for id 7")
	}
	v = 14
	err = rb.UpsertById(v, &v)
	if err != nil {
		t.Errorf("Failed to upsert old ID: %+v", err)
	}

	val, _ = rb.GetById(v)
	if *val.(*int) != v {
		t.Errorf("Should have gotten id %v, recieved %v", v, *val.(*int))
	}

	_, err = rb.GetById(20)
	if err != nil && !strings.Contains(err.Error(), "is higher than most recent id") {
		t.Error("Did not get proper error for id too high")
	}

	_, err = rb.GetById(0)
	if err != nil && !strings.Contains(err.Error(), "is lower than oldest id") {
		t.Error("Did not get proper error for id too high")
	}
}

// Test upserting by id on ringbuff
func TestBuff_UpsertById2(t *testing.T) {
	rb := setup()
	v := -5
	err := rb.UpsertById(v, &v)
	if err == nil {
		t.Error("This should have errored: id was too low")
	}
	v = 6
	err = rb.UpsertById(v, &v)
	if err != nil {
		t.Errorf("Should have inserted, insert valid: %+v", err)
	}

}

// test the length function of ring buff
func TestBuff_Len(t *testing.T) {
	rb := setup()
	if rb.Len() != 5 {
		t.Errorf("Got wrong count")
	}
}

// Test GetByIndex in ringbuff
func TestBuff_GetByIndex(t *testing.T) {
	rb := setup()
	val, _ := rb.GetByIndex(0)
	if *val.(*int) != 0 {
		t.Error("Didn't get correct ID")
	}
	v := 5
	rb.Push(&v)
	val, err := rb.GetByIndex(0)
	if err != nil {
		t.Errorf("Get by index should not error: %s", err)
	}
	if *val.(*int) != v {
		t.Error("Didn't get correct ID after pushing")
	}

	_, err = rb.GetByIndex(25)
	if err == nil {
		t.Errorf("Should have received index out of bounds err")
	}
}

// Test the GetById function of ringbuff
func TestBuff_GetNewerById(t *testing.T) {
	rb := setup()

	list, err := rb.GetNewerById(2)
	if err != nil {
		t.Error("Failed to get newer than id 2")
	}

	if len(list) != 2 {
		t.Errorf("list has wrong number of entrees: %s", list)
	}

	if *list[0].(*int) != 3 {
		t.Error("list has wrong number first element")
	}
	if *list[1].(*int) != 4 {
		t.Error("list has wrong number second element")
	}

	//test you get all when the id is less than the oldest id
	list, err = rb.GetNewerById(-1)
	if len(list) != 5 {
		t.Errorf("list has wrong number of entrees: %v", len(list))
	}

	//test you get an error when you ask for something newer than the newest
	list, err = rb.GetNewerById(42)
	if list != nil {
		t.Errorf("list should be nil")
	}

	if err == nil {
		t.Errorf("error should have been returned")
	} else if !strings.Contains(err.Error(), "is higher than the newest id") {
		t.Errorf("wrong error returned: %s", err.Error())
	}

}

// test that when the ring buffer is filled and one more object is added that it
// overwrites the first object
func TestBuff_Overflow(t *testing.T) {
	capacity := 111
	rb := NewBuff(capacity)
	expected := make([]int, 111)
	for i := 0; i < capacity; i++ {
		v := i
		rb.Push(&v)
		expected[i] = i
	}

	for i := 0; i < capacity; i++ {
		val, _ := rb.GetByIndex(i)
		v := *val.(*int)
		if v != expected[i] {
			t.Errorf("Element %v not as expected. Expected: %v, Recieved: %v", i, expected[i], v)
		}
	}

	if rb.Len() != capacity {
		t.Errorf("size of the buffer is wrong")
	}

	vNew := 111
	rb.push(&vNew)

	expected[0] = vNew
	for i := 0; i < capacity; i++ {
		val, _ := rb.GetByIndex(i)
		v := *val.(*int)
		if v != expected[i] {
			t.Errorf("Element %v not as expected. Expected: %v, Recieved: %v", i, expected[i], v)
		}
	}

	if rb.Len() != capacity {
		t.Errorf("size of the buffer is wrong")
	}
}

// test that when the ring buffer is filled and one more object is added that it
// overwrites the first object
func TestBuff_MajorOverflow(t *testing.T) {
	capacity := 111
	rb := NewBuff(capacity)
	expected := make([]int, 111)
	for i := 0; i < capacity*10+42; i++ {
		v := i
		rb.Push(&v)
		expected[i%capacity] = i
	}

	if rb.Len() != capacity {
		t.Errorf("size of the buffer is wrong")
	}

	for i := 0; i < capacity; i++ {
		val, _ := rb.GetByIndex(i)
		v := *val.(*int)
		if v != expected[i] {
			t.Errorf("Element %v not as expected. Expected: %v, Recieved: %v", i, expected[i], v)
		}
	}

	vNew := 123456789
	rb.push(&vNew)

	expected[(capacity*10+42)%capacity] = vNew
	for i := 0; i < capacity; i++ {
		val, _ := rb.GetByIndex(i)
		v := *val.(*int)
		if v != expected[i] {
			t.Errorf("Element %v not as expected. Expected: %v, Recieved: %v", i, expected[i], v)
		}
	}

	if rb.Len() != capacity {
		t.Errorf("size of the buffer is wrong")
	}
}
