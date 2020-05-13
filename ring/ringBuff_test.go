////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package ring

import (
	"strings"
	"testing"
)

// Basic interface to use for testing
type Tester struct {
	Id int
}

// ID func for tester object
func id(val interface{}) int {
	return val.(*Tester).Id
}
func comp(old, new interface{}) bool {
	if old != nil {
		return false
	}
	return true
}

func getTester(id int) *Tester {
	return &Tester{Id: id}
}

// Setup func for tests
func setup() *Buff {
	rb := NewBuff(5, id)
	for i := 1; i <= 5; i++ {
		rb.Push(&Tester{
			Id: i,
		})
	}
	return rb
}

func TestBuff_GetNewestId(t *testing.T) {
	rb := NewBuff(5, id)
	id := rb.GetNewestId()
	if id != -1 {
		t.Error("Should have returned -1 for empty buff")
	}

	rb = setup()
	_ = rb.UpsertById(getTester(7), comp)
	id = rb.GetNewestId()
	if id != 7 {
		t.Error("Should have returned last pushed id")
	}
}

func TestBuff_GetOldestId(t *testing.T) {
	rb := NewBuff(5, id)
	id := rb.GetOldestId()
	if id != -1 {
		t.Error("Should have returned -1 for empty buff")
	}

	rb = setup()
	rb.Push(getTester(6))
	id = rb.GetOldestId()
	if id != 2 {
		t.Errorf("Should have returned 2, instead got: %d", id)
	}

	_ = rb.UpsertById(getTester(22), comp)
	id = rb.GetOldestId()
	if id != 22 {
		t.Errorf("Should have returned 22, instead got %d", id)
	}
}

func TestBuff_GetOldestIndex(t *testing.T) {
	rb := NewBuff(5, id)
	i := rb.GetOldestIndex()
	if i != -1 {
		t.Error("Should have returned -1 for empty buff")
	}

	rb = setup()
	rb.Push(getTester(6))
	i = rb.GetOldestIndex()
	if i != 1 {
		t.Errorf("Should have returned 1, instead got: %d", i)
	}
}

// Test the Get function on ringbuff
func TestBuff_Get(t *testing.T) {
	rb := setup()
	val := rb.Get().(*Tester)
	if val.Id != 5 {
		t.Errorf("Did not get most recent ID")
	}
}

// Test the GetById function of ringbuff
func TestBuff_GetById(t *testing.T) {
	rb := setup()
	val, err := rb.GetById(3)
	if err != nil {
		t.Error("Failed to get id 3")
	}
	val = val.(*Tester).Id
	if val != 3 {
		t.Error("Got the wrong id")
	}
}

// Test the basic push function on RingBuff
func TestBuff_Push(t *testing.T) {
	rb := setup()
	oldFirst := rb.old
	rb.Push(&Tester{
		Id: 6,
	})
	if rb.old != oldFirst+1 {
		t.Error("Didn't increment old properly")
	}
	val := rb.Get().(*Tester)
	if val.Id != 6 {
		t.Error("Did not get newest id")
	}
}

// Test ID upsert on ringbuff (bulk of cases)
func TestBuff_UpsertById(t *testing.T) {
	rb := setup()
	err := rb.UpsertById(&Tester{
		Id: 8,
	}, comp)
	if err != nil {
		t.Errorf("Error on initial upsert: %+v", err)
	}
	if rb.Get().(*Tester).Id != 8 {
		t.Error("Failed to get correct ID")
	}

	val, _ := rb.GetById(7)
	if val != nil {
		t.Errorf("Should have gotten nil value for id 7")
	}

	err = rb.UpsertById(&Tester{
		Id: 7,
	}, comp)
	if err != nil {
		t.Errorf("Failed to upsert old ID: %+v", err)
	}

	err = rb.UpsertById(&Tester{
		Id: 7,
	}, comp)
	if err == nil {
		t.Errorf("Should have received error for failed comp function")
	}

	val, _ = rb.GetById(7)
	if val.(*Tester).Id != 7 {
		t.Errorf("Should have gotten id 7")
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
	comp := func(old, new interface{}) bool {
		if old != nil {
			return false
		}
		return true
	}
	rb := setup()
	err := rb.UpsertById(&Tester{
		Id: -5,
	}, comp)
	if err == nil {
		t.Error("This should have errored: id was too low")
	}
	err = rb.UpsertById(&Tester{
		Id: 6,
	}, comp)
	if err != nil {
		t.Errorf("Should have inserted using first case: %+v", err)
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
	if val.(*Tester).Id != 1 {
		t.Error("Didn't get correct ID")
	}

	rb.Push(&Tester{
		Id: 6,
	})
	val, _ = rb.GetByIndex(0)
	if val.(*Tester).Id != 2 {
		t.Error("Didn't get correct ID after pushing")
	}

	_, err := rb.GetByIndex(25)
	if err == nil {
		t.Errorf("Should have received index out of bounds err")
	}
}

// Test the GetById function of ringbuff
func TestBuff_GetNewerById(t *testing.T) {
	rb := setup()

	list, err := rb.GetNewerById(3)
	if err != nil {
		t.Error("Failed to get newer than id 3")
	}

	if len(list) != 2 {
		t.Errorf("list has wrong number of entrees: %s", list)
	}

	if list[0].(*Tester).Id != 4 {
		t.Error("list has wrong number first element")
	}
	if list[1].(*Tester).Id != 5 {
		t.Error("list has wrong number second element")
	}

	//test you get all when the id is less than the oldest id
	list, err = rb.GetNewerById(-1)
	if len(list) != 5 {
		t.Errorf("list has wrong number of entrees: %s", list)
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
