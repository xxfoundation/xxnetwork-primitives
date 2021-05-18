///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package knownRounds

import (
	"bytes"
	"fmt"
	"gitlab.com/xx_network/primitives/id"
	"math"
	"reflect"
	"testing"
)

// Tests happy path of NewKnownRound().
func TestNewKnownRound(t *testing.T) {
	expectedKR := &KnownRounds{
		bitStream:      uint64Buff{0, 0, 0, 0, 0},
		firstUnchecked: 0,
		lastChecked:    0,
		fuPos:          0,
	}

	testKR := NewKnownRound(310)

	if !reflect.DeepEqual(testKR, expectedKR) {
		t.Errorf("NewKnownRound() did not produce the expected KnownRounds."+
			"\n\texpected: %v\n\treceived: %v",
			expectedKR, testKR)
	}
}

// Tests happy path of KnownRounds.Marshal().
func TestKnownRounds_Marshal_Unmarshal(t *testing.T) {
	testKR := &KnownRounds{
		bitStream:      uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0},
		firstUnchecked: 55,
		lastChecked:    270,
		fuPos:          55,
	}

	data := testKR.Marshal()

	newKR := &KnownRounds{}
	err := newKR.Unmarshal(data)
	if err != nil {
		t.Errorf("Unmarshal() produced an error: %+v", err)
	}

	if !reflect.DeepEqual(testKR, newKR) {
		t.Errorf("Original KnownRounds does not match Unmarshalled."+
			"\nexpected: %+v\nreceived: %+v", testKR, newKR)
	}
}

// Tests happy path of KnownRounds.Marshal().
func TestKnownRounds_Marshal(t *testing.T) {
	testKR := &KnownRounds{
		bitStream:      uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0},
		firstUnchecked: 75,
		lastChecked:    150,
		fuPos:          75,
	}

	expectedData := []byte{75, 0, 0, 0, 0, 0, 0, 0, 150, 0, 0, 0, 0, 0, 0, 0,
		255, 255, 255, 255, 255, 255, 255, 255, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}

	data := testKR.Marshal()

	if !bytes.Equal([]byte(expectedData), data) {
		t.Errorf("Marshal() produced incorrect data."+
			"\n\texpected: %+v\n\treceived: %+v", expectedData, data)
	}

}

// Tests happy path of KnownRounds.Unmarshal().
func TestKnownRounds_Unmarshal(t *testing.T) {
	testKR := &KnownRounds{
		bitStream:      uint64Buff{0, math.MaxUint64, 0, 0, 0},
		firstUnchecked: 75,
		lastChecked:    150,
		fuPos:          11,
	}

	data := testKR.Marshal()

	newKR := NewKnownRound(310)
	err := newKR.Unmarshal(data)
	if err != nil {
		t.Errorf("Unmarshal() produced an unexpected error."+
			"\n\texpected: %+v\n\treceived: %+v", nil, err)
	}

	if !reflect.DeepEqual(newKR, testKR) {
		t.Errorf("Unmarshal() produced an incorrect KnownRounds from the data."+
			"\n\texpected: %v\n\treceived: %v", testKR, newKR)
	}
}

// Tests that KnownRounds.Unmarshal() errors when the new bit stream is too
// small.
func TestKnownRounds_Unmarshal_SizeError(t *testing.T) {
	testKR := &KnownRounds{
		bitStream:      uint64Buff{0, math.MaxUint64, 0, 0, 0},
		firstUnchecked: 75,
		lastChecked:    150,
		fuPos:          11,
	}

	data := testKR.Marshal()

	newKR := NewKnownRound(1)
	err := newKR.Unmarshal(data)
	if err == nil {
		t.Error("Unmarshal() did not produce an error when the size of new " +
			"KnownRound bit stream is too small.")
	}
}

// Tests that KnownRounds.Unmarshal() errors when given invalid JSON data.
func TestKnownRounds_Unmarshal_JsonError(t *testing.T) {
	newKR := NewKnownRound(1)
	err := newKR.Unmarshal([]byte("hello"))
	if err == nil {
		t.Error("Unmarshal() did not produce an error on invalid JSON data.")
	}
}

// Tests happy path of KnownRounds.Check().
func TestKnownRounds_Check(t *testing.T) {
	// Generate test round IDs and expected buffers
	testData := []struct {
		rid, expectedLastChecked id.Round
		buff                     uint64Buff
	}{
		{0, 200, uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0}},
		{75, 200, uint64Buff{4503599627370496, math.MaxUint64, 0, math.MaxUint64, 0}},
		{95, 200, uint64Buff{4294967296, math.MaxUint64, 0, math.MaxUint64, 0}},
		{150, 200, uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0}},
		{320, 320, uint64Buff{0, math.MaxUint64, 0, 0, 0x8000000000000000}},
		{519, 519, uint64Buff{0, 0, 0x100000000000000, 0, 0}},
	}

	for i, data := range testData {
		kr := KnownRounds{
			bitStream:      uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0},
			firstUnchecked: 75,
			lastChecked:    200,
			fuPos:          11,
		}

		kr.Check(data.rid)
		if !reflect.DeepEqual(kr.bitStream, data.buff) {
			t.Errorf("Incorrect resulting buffer after checking round ID %d (round %d)."+
				"\n\texpected: %064b\n\treceived: %064b"+
				"\n\033[38;5;59m               0123456789012345678901234567890123456789012345678901234567890123 4567890123456789012345678901234567890123456789012345678901234567 8901234567890123456789012345678901234567890123456789012345678901 2345678901234567890123456789012345678901234567890123456789012345 6789012345678901234567890123456789012345678901234567890123456789 0123456789012345678901234567890123456789012345678901234567890123"+
				"\n\u001B[38;5;59m               0         1         2         3         4         5         6          7         8         9         0         1         2          3         4         5         6         7         8         9          0         1         2         3         4         5          6         7         8         9         0         1          2         3         4         5         6         7         8"+
				"\n\u001B[38;5;59m               0         0         0         0         0         0         0          0         0         0         1         1         1          1         1         1         1         1         1         1          2         2         2         2         2         2          2         2         2         2         3         3          3         3         3         3         3         3         3",
				data.rid, i, data.buff, kr.bitStream)
		}

		if kr.lastChecked != data.expectedLastChecked {
			t.Errorf("Check() did not modify the the lastChecked round correctly for round ID %d (round %d)."+
				"\n\texpected: %d\n\treceived: %d", data.rid, i, data.expectedLastChecked, kr.lastChecked)
		}
	}
}

// Tests happy path of KnownRounds.Check() with a new KnownRounds.
func TestKnownRounds_Check_NewKR(t *testing.T) {
	// Generate test round IDs and expected buffers
	testData := []struct {
		rid, expectedLastChecked id.Round
		buff                     uint64Buff
	}{
		{1, 1, uint64Buff{0x4000000000000000, 0, 0, 0, 0}},
		{0, 1, uint64Buff{0x8000000000000000, 0, 0, 0, 0}},
		{75, 75, uint64Buff{0, 0x10000000000000, 0, 0, 0}},
		{319, 319, uint64Buff{0, 0, 0, 0, 1}},
	}

	for i, data := range testData {
		kr := NewKnownRound(310)
		kr.Check(data.rid)
		if !reflect.DeepEqual(kr.bitStream, data.buff) {
			t.Errorf("Resulting buffer after checking round ID %d (round %d)."+
				"\n\texpected: %064b\n\treceived: %064b",
				data.rid, i, data.buff, kr.bitStream)
		}

		if kr.lastChecked != data.expectedLastChecked {
			t.Errorf("Check() did not modify the the lastChecked round correctly for round ID %d (round %d)."+
				"\n\texpected: %d\n\treceived: %d", data.rid, i, data.expectedLastChecked, kr.lastChecked)
		}
	}
}

// Happy path of KnownRounds.Checked().
func TestKnownRounds_Checked(t *testing.T) {
	// Generate test positions and expected value
	testData := []struct {
		rid   id.Round
		value bool
	}{
		{75, false},
		{76, false},
		{123, false},
		{124, false},
		{74, true},
		{60, true},
		{0, true},
		{319, false},
		{320, false},
	}
	kr := KnownRounds{
		bitStream:      uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0},
		firstUnchecked: 75,
		lastChecked:    200,
		fuPos:          11,
	}

	for i, data := range testData {
		value := kr.Checked(data.rid)
		if value != data.value {
			t.Errorf("Checked() returned incorrect value for round ID %d (round %d)."+
				"\n\texpected: %v\n\treceived: %v", data.rid, i, data.value, value)
		}
	}
}

// Happy path of KnownRounds.Checked() with a new KnownRounds.
func TestKnownRounds_Checked_NewKR(t *testing.T) {
	// Generate test positions and expected value
	testData := []struct {
		rid   id.Round
		value bool
	}{
		{0, false},
		{1, false},
		{2, false},
		{320, false},
	}

	for i, data := range testData {
		kr := NewKnownRound(5)
		value := kr.Checked(data.rid)
		if value != data.value {
			t.Errorf("Checked() returned incorrect value for round ID %d (round %d)."+
				"\n\texpected: %v\n\treceived: %v", data.rid, i, data.value, value)
		}
	}
}

// Tests happy path of KnownRounds.Forward().
func TestKnownRounds_Forward(t *testing.T) {
	// Generate test round IDs and expected buffers
	testData := []struct {
		rid, expectedFirstChecked, expectedLastChecked id.Round
		expectedFusPos                                 int
	}{
		{75, 75, 200, 11},
		{76, 76, 200, 12},
		{192, 192, 200, 128},
		{150, 192, 200, 128},
		{200, 200, 200, 136},
		{210, 210, 210, 210 % 64},
	}
	kr := KnownRounds{
		bitStream:      uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0},
		firstUnchecked: 75,
		lastChecked:    200,
		fuPos:          11,
	}

	for i, data := range testData {
		kr.bitStream = uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0}
		kr.Forward(data.rid)
		if kr.firstUnchecked != data.expectedFirstChecked {
			t.Errorf("Forward() did not modify the the firstUnchecked round correctly for round ID %d (round %d)."+
				"\n\texpected: %d\n\treceived: %d", data.rid, i, data.expectedFirstChecked, kr.firstUnchecked)
		}
		if kr.lastChecked != data.expectedLastChecked {
			t.Errorf("Forward() did not modify the the lastChecked round correctly for round ID %d (round %d)."+
				"\n\texpected: %d\n\treceived: %d", data.rid, i, data.expectedLastChecked, kr.lastChecked)
		}
		if kr.fuPos != data.expectedFusPos {
			t.Errorf("Forward() did not modify the the fuPos round correctly for round ID %d (round %d)."+
				"\n\texpected: %d\n\treceived: %d", data.rid, i, data.expectedFusPos, kr.fuPos)
		}
	}
}

// Tests happy path of KnownRounds.Forward() with a new KnownRounds.
func TestKnownRounds_Forward_NewKR(t *testing.T) {
	// Generate test round IDs and expected buffers
	testData := []struct {
		rid, expectedFirstUnchecked, expectedLastChecked id.Round
		expectedFusPos                                   int
	}{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{2, 2, 2, 2},
		{320, 320, 320, 0},
	}

	for i, data := range testData {
		kr := NewKnownRound(5)
		kr.Forward(data.rid)
		if kr.firstUnchecked != data.expectedFirstUnchecked {
			t.Errorf("Forward() did not modify the the firstUnchecked round correctly for round ID %d (round %d)."+
				"\n\texpected: %d\n\treceived: %d", data.rid, i, data.expectedFirstUnchecked, kr.firstUnchecked)
		}
		if kr.lastChecked != data.expectedLastChecked {
			t.Errorf("Forward() did not modify the the lastChecked round correctly for round ID %d (round %d)."+
				"\n\texpected: %d\n\treceived: %d", data.rid, i, data.expectedLastChecked, kr.lastChecked)
		}
		if kr.fuPos != data.expectedFusPos {
			t.Errorf("Forward() did not modify the the fuPos round correctly for round ID %d (round %d)."+
				"\n\texpected: %d\n\treceived: %d", data.rid, i, data.expectedFusPos, kr.fuPos)
		}
	}
}

// Test happy path of KnownRounds.RangeUnchecked().
func TestKnownRounds_RangeUnchecked(t *testing.T) {
	// Generate test round IDs and expected buffers
	testData := []struct {
		oldestUnknown, expected id.Round
		has, unknown            []id.Round
	}{
		{55, 141, makeRange(55, 127), makeRange(128, 140)},
		{65, 141, makeRange(65, 127), makeRange(128, 140)},
		{75, 141, makeRange(75, 127), makeRange(128, 140)},
		{85, 141, makeRange(85, 127), makeRange(128, 140)},
		{191, 191, nil, nil},
		{192, 192, nil, nil},
		{292, 292, nil, nil},
	}
	roundCheck := func(id id.Round) bool {
		return true
	}

	// Bitstream = 0x00000000, 0xFFFFFFFF, 0x00000000, 0xFFFFFFFF, 0x00000000
	//                            ^ firstUnchecked, fuPos (position 75)
	//                                                   ^ lastChecked (position 191)
	//                                       ^-^ unknown (position 128 - 140)
	//                           ^-----^ has (position 64 - 127)
	//                     xx also has (position 55 - 63), outside defined buffer
	for i, data := range testData {
		kr := KnownRounds{
			bitStream:      uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0},
			firstUnchecked: 75,
			lastChecked:    191,
			fuPos:          75,
		}

		earliestRound, has, unknown := kr.RangeUnchecked(data.oldestUnknown, 50, roundCheck)

		if earliestRound != data.expected {
			t.Errorf("RangeUnchecked() did not return the correct round (%d)."+
				"\n\texpected: %d\n\treceived: %d",
				i, data.expected, earliestRound)
		}

		if !reflect.DeepEqual(data.has, has) {
			t.Errorf("RangeUnchecked() did not return the correct has list (%d)."+
				"\n\texpected: %v\n\treceived: %v",
				i, data.has, has)
		}

		if !reflect.DeepEqual(data.unknown, unknown) {
			t.Errorf("RangeUnchecked() did not return the correct unknown list (%d)."+
				"\n\texpected: %v\n\treceived: %v",
				i, data.unknown, unknown)
		}
	}
}

// Test happy path of KnownRounds.RangeUnchecked() with a new KnownRounds.
func TestKnownRounds_RangeUnchecked_NewKR(t *testing.T) {
	// Generate test round IDs and expected buffers
	testData := []struct {
		oldestUnknown, expected id.Round
		has, unknown            []id.Round
	}{
		{55, 55, nil, nil},
		{65, 65, nil, nil},
		{75, 75, nil, nil},
		{85, 85, nil, nil},
		{191, 191, nil, nil},
		{192, 192, nil, nil},
		{292, 292, nil, nil},
	}
	roundCheck := func(id id.Round) bool {
		return id%2 == 1
	}

	for i, data := range testData {
		kr := NewKnownRound(310)

		earliestRound, has, unknown := kr.RangeUnchecked(data.oldestUnknown, 50, roundCheck)

		if earliestRound != data.expected {
			t.Errorf("RangeUnchecked() did not return the correct round (%d)."+
				"\n\texpected: %d\n\treceived: %d",
				i, data.expected, earliestRound)
		}

		if !reflect.DeepEqual(data.has, has) {
			t.Errorf("RangeUnchecked() did not return the correct has list (%d)."+
				"\n\texpected: %v\n\treceived: %v",
				i, data.has, has)
		}

		if !reflect.DeepEqual(data.unknown, unknown) {
			t.Errorf("RangeUnchecked() did not return the correct unknown list (%d)."+
				"\n\texpected: %v\n\treceived: %v",
				i, data.unknown, unknown)
		}
	}
}

// Test happy path of KnownRounds.RangeUncheckedMasked().
func TestKnownRounds_RangeUncheckedMasked(t *testing.T) {
	expectedKR := KnownRounds{
		bitStream:      uint64Buff{42949672960, 18446744073709551615, 0, 18446744073709551615, 0},
		firstUnchecked: 15,
		lastChecked:    191,
		fuPos:          0,
	}
	kr := KnownRounds{
		bitStream:      uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0},
		firstUnchecked: 15,
		lastChecked:    191,
		fuPos:          0,
	}
	kr2 := &KnownRounds{
		bitStream:      uint64Buff{math.MaxUint64},
		firstUnchecked: 20,
		lastChecked:    47,
		fuPos:          0,
	}

	roundCheck := func(id id.Round) bool {
		return id%2 == 1
	}

	kr.RangeUncheckedMasked(kr2, roundCheck, 5)
	if !reflect.DeepEqual(expectedKR, kr) {
		t.Errorf("RangeUncheckedMasked() incorrectl modified KnownRounds."+
			"\n\texpected: %+v\n\treceived: %+v", expectedKR, kr)
	}
	fmt.Printf("kr.bitStream: %+v\n", kr.bitStream)
}

// Happy path of getBitStreamPos().
func TestKnownRounds_getBitStreamPos(t *testing.T) {
	// Generate test round IDs and their expected positions
	testData := []struct {
		rid id.Round
		pos int
	}{
		{75, 11},
		{76, 12},
		{123, 59},
		{124, 60},
		{74, 10},
		{60, -4},
		{0, -64},
		{319, 255},
		{320, 256},
	}
	kr := KnownRounds{
		bitStream:      uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0},
		firstUnchecked: 75,
		lastChecked:    85,
		fuPos:          11,
	}
	for i, data := range testData {
		pos := kr.getBitStreamPos(data.rid)
		if pos != data.pos {
			t.Errorf("getBitStreamPos() returned incorrect position for round ID %d (round %d)."+
				"\n\texpected: %v\n\treceived: %v", data.rid, i, data.pos, pos)
		}
	}
}

/*
// Test happy path of KnownRounds.RangeUncheckedMasked().
func TestKnownRounds_RangeUncheckedMasked_2(t *testing.T) {
	expectedKR := KnownRounds{
		bitStream:      make(uint64Buff, 245),
		firstUnchecked: 30,
		lastChecked:    57,
		fuPos:          30,
	}
	expectedKR.bitStream[0] = 0xFFFFFFFD40000040
	kr := KnownRounds{
		bitStream:      make(uint64Buff, 245),
		firstUnchecked: 30,
		lastChecked:    39,
		fuPos:          30,
	}
	kr.bitStream[0] = 0xFFFFFFFC00000000

	mask := &KnownRounds{
		bitStream:      uint64Buff{0xFEFFFFFBFFFFFFC0},
		firstUnchecked: 7,
		lastChecked:    57,
		fuPos:          7,
	}

	roundCheck := func(id id.Round) bool {
		return id%2 == 1
	}

	kr.RangeUncheckedMasked(mask, roundCheck, 5)
	if !reflect.DeepEqual(expectedKR, kr) {
		t.Errorf("RangeUncheckedMasked() incorrect modified KnownRounds."+
			"\n\texpected: %+v\n\treceived: %+v", expectedKR, kr)
	}
	fmt.Printf("kr.bitStream: %064b\n", kr.bitStream)
}*/

// // Test happy path of KnownRounds.RangeUncheckedMasked().
// func TestKnownRounds_RangeUncheckedMasked_3(t *testing.T) {
// 	expectedKR := KnownRounds{
// 		bitStream:      make(uint64Buff, 245),
// 		firstUnchecked: 30,
// 		lastChecked:    57,
// 		fuPos:          30,
// 	}
// 	expectedKR.bitStream[0] = 0b1111111111010000000000000000000000000000000000000000000000000000
// 	kr := KnownRounds{
// 		bitStream:      make(uint64Buff, 245),
// 		firstUnchecked: 9,
// 		lastChecked:    9,
// 		fuPos:          9,
// 	}
// 	kr.bitStream[0] = 0b1111111111000000000000000000000000000000000000000000000000000000
//
// 	mask := &KnownRounds{
// 		bitStream:      uint64Buff{0b1111111101111111111111111111111111111111111111111111111111111111, 0b1111100000000000000000000000000000000000000000000000000000000000},
// 		firstUnchecked: 8,
// 		lastChecked:    68,
// 		fuPos:          8,
// 	}
//
// 	roundCheck := func(id id.Round) bool {
// 		return id%2 == 1
// 	}
//
// 	kr.RangeUncheckedMasked(mask, roundCheck, 5)
// 	if !reflect.DeepEqual(expectedKR, kr) {
// 		t.Errorf("RangeUncheckedMasked() incorrect modified KnownRounds."+
// 			"\n\texpected: %064b\n\treceived: %064b", expectedKR, kr)
// 	}
// 	fmt.Printf("kr.bitStream: %064b\n", kr.bitStream)
// }

//
// // Tests that KnownRounds.subSample() produces the correct buffer for a new
// // KnownRounds.
// func TestKnownRounds_subSample(t *testing.T) {
// 	kr := NewKnownRound(1)
// 	expectedU64b := make(uint64Buff, 3)
//
// 	fmt.Printf("kr: %+v\n", kr)
//
// 	u64b, length := kr.subSample(5, 189)
// 	if !reflect.DeepEqual(expectedU64b, u64b) {
// 		t.Errorf("subSample() returned incorrect buffer." +
// 			"\n\texpected: %064b\n\treceived: %064b", expectedU64b, u64b)
// 	}
//
// 	if len(expectedU64b) != length {
// 		t.Errorf("subSample() returned incorrect buffer length." +
// 			"\n\texpected: %d\n\treceived: %d", len(expectedU64b), length)
// 	}
// }
//
// // Tests that KnownRounds.subSample() produces the correct buffer for a new
// // KnownRounds.
// func TestKnownRounds_subSample2(t *testing.T) {
// 	kr := &KnownRounds{
// 		bitStream:      make(uint64Buff, 15626),
// 		firstUnchecked: 23,
// 		lastChecked:    22,
// 		fuPos:          23,
// 	}
// 	mask := &KnownRounds{
// 		bitStream:      make(uint64Buff, 1),
// 		firstUnchecked: 0,
// 		lastChecked:    1,
// 		fuPos:          0,
// 	}
// 	fmt.Printf("mask: %+v\n", mask)
// 	mask.Forward(kr.firstUnchecked)
// 	fmt.Printf("mask: %+v\n", mask)
// 	expectedU64b := make(uint64Buff, 3)
//
//
// 	u64b, length := kr.subSample(mask.firstUnchecked, mask.lastChecked)
// 	if !reflect.DeepEqual(expectedU64b, u64b) {
// 		t.Errorf("subSample() returned incorrect buffer." +
// 			"\n\texpected: %064b\n\treceived: %064b", expectedU64b, u64b)
// 	}
//
// 	if len(expectedU64b) != length {
// 		t.Errorf("subSample() returned incorrect buffer length." +
// 			"\n\texpected: %d\n\treceived: %d", len(expectedU64b), length)
// 	}
// }
//
// func TestKnownRounds_RangeUncheckedMasked2(t *testing.T) {
// 	kr := &KnownRounds{
// 		bitStream:      make(uint64Buff, 15626),
// 		firstUnchecked: 23,
// 		lastChecked:    22,
// 		fuPos:          23,
// 	}
// 	mask := &KnownRounds{
// 		bitStream:      make(uint64Buff, 1),
// 		firstUnchecked: 0,
// 		lastChecked:    1,
// 		fuPos:          0,
// 	}
//
// 	roundCheck := func(id id.Round) bool {
// 		return id%2 == 1
// 	}
//
// 	kr.RangeUncheckedMasked(mask, roundCheck, 500)
// }

func makeRange(min, max int) []id.Round {
	a := make([]id.Round, max-min+1)
	for i := range a {
		a[i] = id.Round(min + i)
	}
	return a
}
