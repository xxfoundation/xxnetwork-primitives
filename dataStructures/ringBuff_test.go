package dataStructures

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

// Setup func for tests
func setup() *RingBuff {
	rb := NewRingBuff(5, id)
	for i := 1; i <= 5; i++ {
		rb.Push(&Tester{
			Id: i,
		})
	}
	return rb
}

// Test the Get function on ringbuff
func TestRingBuff_Get(t *testing.T) {
	rb := setup()
	val := rb.Get().(*Tester)
	if val.Id != 5 {
		t.Errorf("Did not get most recent ID")
	}
}

// Test the GetById function of ringbuff
func TestRingBuff_GetById(t *testing.T) {
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
func TestRingBuff_Push(t *testing.T) {
	rb := setup()
	oldFirst := rb.head
	rb.Push(&Tester{
		Id: 6,
	})
	if rb.head != oldFirst+1 {
		t.Error("Didn't increment head properly")
	}
	val := rb.Get().(*Tester)
	if val.Id != 6 {
		t.Error("Did not get newest id")
	}
}

// Test ID upsert on ringbuff (bulk of cases)
func TestRingBuff_UpsertById(t *testing.T) {
	comp := func(old, new interface{}) bool {
		if old != nil {
			return false
		}
		return true
	}
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
func TestRingBuff_UpsertById2(t *testing.T) {
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
func TestRingBuff_Len(t *testing.T) {
	rb := setup()
	if rb.Len() != 5 {
		t.Errorf("Got wrong count")
	}
}

// Test GetByIndex in ringbuff
func TestRingBuff_GetByIndex(t *testing.T) {
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
