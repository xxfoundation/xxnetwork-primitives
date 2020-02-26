package dataStructures

import (
	"strings"
	"testing"
)

type Tester struct {
	Id int
}

func id(val interface{}) int {
	return val.(*Tester).Id
}

func setup() *RingBuff {
	rb := NewRingBuff(5, id)
	for i := 1; i <= 5; i++ {
		rb.Push(&Tester{
			Id: i,
		})
	}
	return rb
}

func TestRingBuff_Get(t *testing.T) {
	rb := setup()
	val := rb.Get().(*Tester)
	if val.Id != 5 {
		t.Errorf("Did not get most recent ID")
	}
}

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

func TestRingBuff_Len(t *testing.T) {
	rb := setup()
	if rb.Len() != 5 {
		t.Errorf("Got wrong count")
	}
}
