////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package knownRounds

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"gitlab.com/xx_network/primitives/id"
)

// Tests happy path of NewKnownRound.
func TestNewKnownRound(t *testing.T) {
	expectedKR := &KnownRounds{
		bitStream:      uint64Buff{0, 0, 0, 0, 0},
		firstUnchecked: 0,
		lastChecked:    0,
		fuPos:          0,
	}

	testKR := NewKnownRound(310)

	if !reflect.DeepEqual(testKR, expectedKR) {
		t.Errorf("NewKnownRound did not produce the expected KnownRounds."+
			"\nexpected: %v\nreceived: %v",
			expectedKR, testKR)
	}
}

// Happy path.
func TestNewFromParts(t *testing.T) {
	expected := &KnownRounds{
		bitStream:      uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0},
		firstUnchecked: 75,
		lastChecked:    150,
		fuPos:          75,
	}

	received := NewFromParts(expected.bitStream, expected.firstUnchecked,
		expected.lastChecked, expected.fuPos)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("NewFromParts did not return the expected KnownRounds."+
			"\nexpected: %v\nreceived: %v", expected, received)
	}
}

// Tests happy path of KnownRounds.Marshal.
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
		t.Errorf("Unmarshal produced an error: %+v", err)
	}

	if !reflect.DeepEqual(testKR, newKR) {
		t.Errorf("Original KnownRounds does not match Unmarshalled."+
			"\nexpected: %+v\nreceived: %+v", testKR, newKR)
	}
}

// Tests happy path of KnownRounds.Marshal.
func TestKnownRounds_Marshal(t *testing.T) {
	testKR := &KnownRounds{
		bitStream:      uint64Buff{0, math.MaxUint64, 0, math.MaxUint64, 0},
		firstUnchecked: 75,
		lastChecked:    150,
		fuPos:          75,
	}

	expectedData := []byte{75, 0, 0, 0, 0, 0, 0, 0, 150, 0, 0, 0, 0, 0, 0, 0, 2,
		1, 255, 8, 0, 8}

	data := testKR.Marshal()

	if !bytes.Equal(expectedData, data) {
		t.Errorf("Marshal produced incorrect data."+
			"\nexpected: %+v\nreceived: %+v", expectedData, data)
	}

}

// Tests happy path of KnownRounds.Unmarshal.
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
		t.Errorf("Unmarshal produced an unexpected error."+
			"\nexpected: %+v\nreceived: %+v", nil, err)
	}

	if !reflect.DeepEqual(newKR, testKR) {
		t.Errorf("Unmarshal produced an incorrect KnownRounds from the data."+
			"\nexpected: %v\nreceived: %v", testKR, newKR)
	}
}

// Tests that KnownRounds.Unmarshal errors when the new bit stream is too
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
		t.Error("Unmarshal did not produce an error when the size of new " +
			"KnownRound bit stream is too small.")
	}
}

// Tests that KnownRounds.Unmarshal errors when given invalid JSON data.
func TestKnownRounds_Unmarshal_JsonError(t *testing.T) {
	newKR := NewKnownRound(1)
	err := newKR.Unmarshal([]byte("hello"))
	if err == nil {
		t.Error("Unmarshal did not produce an error on invalid JSON data.")
	}
}

// Happy path.
func TestKnownRounds_OutputBuffChanges(t *testing.T) {
	// Generate test round IDs and expected buffers
	const max = math.MaxUint64
	testData := []struct {
		current KnownRounds
		old     []uint64
		changes KrChanges
	}{{
		current: KnownRounds{uint64Buff{}, 75, 320, 75},
		old:     []uint64{},
		changes: KrChanges{},
	}, {
		current: KnownRounds{uint64Buff{0, max, 0, max, 0}, 75, 320, 75},
		old:     []uint64{0, max, 0, max, 0},
		changes: KrChanges{},
	}, {
		current: KnownRounds{uint64Buff{0, max, 0, max, 0}, 75, 320, 75},
		old:     []uint64{0, max, 0, max, 0},
		changes: KrChanges{},
	}, {
		current: KnownRounds{uint64Buff{1, max, 0, max, 0}, 75, 320, 75},
		old:     []uint64{0, max, 0, max, 0},
		changes: KrChanges{0: 1},
	}, {
		current: KnownRounds{uint64Buff{0, max, 0, max, 0}, 75, 320, 75},
		old:     []uint64{max, 0, max, 0, max},
		changes: KrChanges{0: 0, 1: max, 2: 0, 3: max, 4: 0},
	}}

	for i, data := range testData {
		changes, firstUnchecked, lastChecked, fuPos, err :=
			data.current.OutputBuffChanges(data.old)
		if err != nil {
			t.Errorf("OutputBuffChanges produced an error (%d): %+v", i, err)
		}

		if data.current.firstUnchecked != firstUnchecked {
			t.Errorf("OutputBuffChanges returned incorrect firstUnchecked (%d)."+
				"\nexpected: %d\nreceived: %d",
				i, data.current.firstUnchecked, firstUnchecked)
		}

		if data.current.lastChecked != lastChecked {
			t.Errorf("OutputBuffChanges returned incorrect lastChecked (%d)."+
				"\nexpected: %d\nreceived: %d",
				i, data.current.lastChecked, lastChecked)
		}

		if data.current.fuPos != fuPos {
			t.Errorf("OutputBuffChanges returned incorrect fuPos (%d)."+
				"\nexpected: %d\nreceived: %d", i, data.current.fuPos, fuPos)
		}

		if !reflect.DeepEqual(data.changes, changes) {
			t.Errorf("OutputBuffChanges returned incorrect changes (%d)."+
				"\nexpected: %v\nreceived: %v", i, data.changes, changes)
		}
	}
}

// Error path: buffers are not the same length.
func TestKnownRounds_OutputBuffChanges_IncorrectLengthError(t *testing.T) {
	// Generate test round IDs and expected buffers
	const max = math.MaxUint64
	testData := []struct {
		current KnownRounds
		old     []uint64
	}{{
		current: KnownRounds{uint64Buff{0, max, 0, max, 0}, 75, 320, 75},
		old:     []uint64{0, max, 0},
	}, {
		current: KnownRounds{uint64Buff{0, max, 0}, 75, 320, 75},
		old:     []uint64{0, max, 0, max, 0},
	}}

	expectedErr := "not the same as length of the current buffer"
	for i, data := range testData {
		_, _, _, _, err := data.current.OutputBuffChanges(data.old)
		if err == nil || !strings.Contains(err.Error(), expectedErr) {
			t.Errorf("OutputBuffChanges did not produce the expected error "+
				"when the buffers are the wrong lengths (%d)."+
				"\nexpected: %s\nreceived: %+v", i, expectedErr, err)
		}
	}
}

// Tests that KnownRounds.GetFirstUnchecked returns the expected value.
func TestKnownRounds_GetFirstUnchecked(t *testing.T) {
	kr := KnownRounds{
		bitStream:      uint64Buff{0, 1, 2, 3, 4, 5, 6, 7},
		firstUnchecked: 65,
		lastChecked:    556,
		fuPos:          1,
	}

	if kr.firstUnchecked != kr.GetFirstUnchecked() {
		t.Errorf("GetFirstUnchecked did not return the expected value."+
			"\nexpected: %d\nreceived: %d", kr.firstUnchecked, kr.GetFirstUnchecked())
	}
}

// Tests that KnownRounds.GetLastChecked returns the expected value.
func TestKnownRounds_GetLastChecked(t *testing.T) {
	kr := KnownRounds{
		bitStream:      uint64Buff{0, 1, 2, 3, 4, 5, 6, 7},
		firstUnchecked: 65,
		lastChecked:    556,
		fuPos:          1,
	}

	if kr.lastChecked != kr.GetLastChecked() {
		t.Errorf("GetLastChecked did not return the expected value."+
			"\nexpected: %d\nreceived: %d", kr.lastChecked, kr.GetLastChecked())
	}
}

// Tests that KnownRounds.GetFuPos returns the expected value.
func TestKnownRounds_GetFuPos(t *testing.T) {
	kr := KnownRounds{
		bitStream:      uint64Buff{0, 1, 2, 3, 4, 5, 6, 7},
		firstUnchecked: 65,
		lastChecked:    556,
		fuPos:          1,
	}

	if kr.fuPos != kr.GetFuPos() {
		t.Errorf("GetFuPos did not return the expected value."+
			"\nexpected: %d\nreceived: %d", kr.fuPos, kr.GetFuPos())
	}
}

// Tests that KnownRounds.GetBitStream returns the expected value.
func TestKnownRounds_GetBitStream(t *testing.T) {
	kr := KnownRounds{
		bitStream:      uint64Buff{0, 1, 2, 3, 4, 5, 6, 7},
		firstUnchecked: 65,
		lastChecked:    556,
		fuPos:          1,
	}

	if !reflect.DeepEqual([]uint64(kr.bitStream), kr.GetBitStream()) {
		t.Errorf("GetFuPos did not return the expected value."+
			"\nexpected: %#v\nreceived: %#v", kr.bitStream, kr.GetBitStream())
	}
}

// Tests happy path of KnownRounds.Check.
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
			t.Errorf("Incorrect resulting buffer after checking round ID %d (%d)."+
				"\nexpected: %064b\nreceived: %064b"+
				"\n\033[38;5;59m               0123456789012345678901234567890123456789012345678901234567890123 4567890123456789012345678901234567890123456789012345678901234567 8901234567890123456789012345678901234567890123456789012345678901 2345678901234567890123456789012345678901234567890123456789012345 6789012345678901234567890123456789012345678901234567890123456789 0123456789012345678901234567890123456789012345678901234567890123"+
				"\n\u001B[38;5;59m               0         1         2         3         4         5         6          7         8         9         0         1         2          3         4         5         6         7         8         9          0         1         2         3         4         5          6         7         8         9         0         1          2         3         4         5         6         7         8"+
				"\n\u001B[38;5;59m               0         0         0         0         0         0         0          0         0         0         1         1         1          1         1         1         1         1         1         1          2         2         2         2         2         2          2         2         2         2         3         3          3         3         3         3         3         3         3",
				data.rid, i, data.buff, kr.bitStream)
		}

		if kr.lastChecked != data.expectedLastChecked {
			t.Errorf("Check did not modify the lastChecked round correctly "+
				"for round ID %d (%d).\nexpected: %d\nreceived: %d",
				data.rid, i, data.expectedLastChecked, kr.lastChecked)
		}
	}
}

// Tests happy path of KnownRounds.Check with a new KnownRounds.
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
			t.Errorf("Resulting buffer after checking round ID %d (%d)."+
				"\nexpected: %064b\nreceived: %064b",
				data.rid, i, data.buff, kr.bitStream)
		}

		if kr.lastChecked != data.expectedLastChecked {
			t.Errorf("Check did not modify the lastChecked round correctly "+
				"for round ID %d (%d).\nexpected: %d\nreceived: %d",
				data.rid, i, data.expectedLastChecked, kr.lastChecked)
		}
	}
}

// Happy path of KnownRounds.Checked.
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
			t.Errorf("Checked returned incorrect value for round ID %d (%d)."+
				"\nexpected: %v\nreceived: %v", data.rid, i, data.value, value)
		}
	}
}

// Happy path of KnownRounds.Checked with a new KnownRounds.
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
			t.Errorf("Checked returned incorrect value for round ID %d (%d)."+
				"\nexpected: %v\nreceived: %v", data.rid, i, data.value, value)
		}
	}
}

// Tests happy path of KnownRounds.Forward.
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
			t.Errorf("Forward did not modify the firstUnchecked round "+
				"correctly for round ID %d (%d).\nexpected: %d\nreceived: %d",
				data.rid, i, data.expectedFirstChecked, kr.firstUnchecked)
		}
		if kr.lastChecked != data.expectedLastChecked {
			t.Errorf("Forward did not modify the lastChecked round correctly "+
				"or round ID %d (%d).\nexpected: %d\nreceived: %d",
				data.rid, i, data.expectedLastChecked, kr.lastChecked)
		}
		if kr.fuPos != data.expectedFusPos {
			t.Errorf("Forward did not modify the fuPos round correctly for "+
				"round ID %d (%d).\nexpected: %d\nreceived: %d",
				data.rid, i, data.expectedFusPos, kr.fuPos)
		}
	}
}

// Tests happy path of KnownRounds.Forward with a new KnownRounds.
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
			t.Errorf("Forward did not modify the firstUnchecked round "+
				"correctly for round ID %d (%d).\nexpected: %d\nreceived: %d",
				data.rid, i, data.expectedFirstUnchecked, kr.firstUnchecked)
		}
		if kr.lastChecked != data.expectedLastChecked {
			t.Errorf("Forward did not modify the lastChecked round correctly "+
				"for round ID %d (%d).\nexpected: %d\nreceived: %d",
				data.rid, i, data.expectedLastChecked, kr.lastChecked)
		}
		if kr.fuPos != data.expectedFusPos {
			t.Errorf("Forward did not modify the fuPos round correctly for "+
				"round ID %d (%d).\nexpected: %d\nreceived: %d",
				data.rid, i, data.expectedFusPos, kr.fuPos)
		}
	}
}

// Test happy path of KnownRounds.RangeUnchecked.
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

		earliestRound, has, unknown :=
			kr.RangeUnchecked(data.oldestUnknown, 50, roundCheck, 1000)

		if earliestRound != data.expected {
			t.Errorf("RangeUnchecked did not return the correct round (%d)."+
				"\nexpected: %d\nreceived: %d",
				i, data.expected, earliestRound)
		}

		if len(data.has) != len(has) {
			t.Errorf("RangeUnchecked did not return the correct has list (%d)."+
				"\nexpected: %v\nreceived: %v",
				i, data.has, has)
		}

		if !reflect.DeepEqual(data.unknown, unknown) {
			t.Errorf("RangeUnchecked did not return the correct unknown list (%d)."+
				"\nexpected: %v\nreceived: %v",
				i, data.unknown, unknown)
		}
	}
}

// Test happy path of KnownRounds.RangeUnchecked with a new KnownRounds.
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

		earliestRound, has, unknown :=
			kr.RangeUnchecked(data.oldestUnknown, 50, roundCheck, 1000)

		if earliestRound != data.expected {
			t.Errorf("RangeUnchecked did not return the correct round (%d)."+
				"\nexpected: %d\nreceived: %d",
				i, data.expected, earliestRound)
		}

		if !reflect.DeepEqual(data.has, has) {
			t.Errorf("RangeUnchecked did not return the correct has list (%d)."+
				"\nexpected: %v\nreceived: %v",
				i, data.has, has)
		}

		if !reflect.DeepEqual(data.unknown, unknown) {
			t.Errorf("RangeUnchecked did not return the correct unknown list (%d)."+
				"\nexpected: %v\nreceived: %v",
				i, data.unknown, unknown)
		}
	}
}

// Test happy path of KnownRounds.RangeUncheckedMasked.
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
		t.Errorf("RangeUncheckedMasked incorrectl modified KnownRounds."+
			"\nexpected: %+v\nreceived: %+v", expectedKR, kr)
	}
	fmt.Printf("kr.bitStream: %+v\n", kr.bitStream)
}

// Happy path of getBitStreamPos.
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
		{60, 316},
		{0, 256},
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
			t.Errorf("getBitStreamPos returned incorrect position for round "+
				"ID %d (%d).\nexpected: %v\nreceived: %v",
				data.rid, i, data.pos, pos)
		}
	}
}

/*
// Test happy path of KnownRounds.RangeUncheckedMasked.
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
		t.Errorf("RangeUncheckedMasked incorrect modified KnownRounds."+
			"\nexpected: %+v\nreceived: %+v", expectedKR, kr)
	}
	fmt.Printf("kr.bitStream: %064b\n", kr.bitStream)
}*/

// // Test happy path of KnownRounds.RangeUncheckedMasked.
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
// 		t.Errorf("RangeUncheckedMasked incorrect modified KnownRounds."+
// 			"\nexpected: %064b\nreceived: %064b", expectedKR, kr)
// 	}
// 	fmt.Printf("kr.bitStream: %064b\n", kr.bitStream)
// }

//
// // Tests that KnownRounds.subSample produces the correct buffer for a new
// // KnownRounds.
// func TestKnownRounds_subSample(t *testing.T) {
// 	kr := NewKnownRound(1)
// 	expectedU64b := make(uint64Buff, 3)
//
// 	fmt.Printf("kr: %+v\n", kr)
//
// 	u64b, length := kr.subSample(5, 189)
// 	if !reflect.DeepEqual(expectedU64b, u64b) {
// 		t.Errorf("subSample returned incorrect buffer." +
// 			"\nexpected: %064b\nreceived: %064b", expectedU64b, u64b)
// 	}
//
// 	if len(expectedU64b) != length {
// 		t.Errorf("subSample returned incorrect buffer length." +
// 			"\nexpected: %d\nreceived: %d", len(expectedU64b), length)
// 	}
// }
//
// // Tests that KnownRounds.subSample produces the correct buffer for a new
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
// 		t.Errorf("subSample returned incorrect buffer." +
// 			"\nexpected: %064b\nreceived: %064b", expectedU64b, u64b)
// 	}
//
// 	if len(expectedU64b) != length {
// 		t.Errorf("subSample returned incorrect buffer length." +
// 			"\nexpected: %d\nreceived: %d", len(expectedU64b), length)
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

func TestKnownRounds_Truncate(t *testing.T) {
	kr := KnownRounds{
		bitStream:      uint64Buff{math.MaxUint64, 0, math.MaxUint64, 0},
		firstUnchecked: 64,
		lastChecked:    130,
		fuPos:          1,
	}

	newKR := kr.Truncate(74)

	if newKR.firstUnchecked != 127 {
		t.Errorf("Failed to truncate. First unchecked not migrated correctly."+
			"\nexpected: %d\nreceived: %d", 127, newKR.firstUnchecked)
	}

	krBytes := kr.Marshal()
	newKrBytes := newKR.Marshal()

	if len(newKrBytes) >= len(krBytes) {
		t.Errorf("Marshalled truncated KR larger than original."+
			"\nexpected: %d\nrecived: %d", len(krBytes), len(newKrBytes))
	}
}

// Same as test above but checking in the case that the circular buffer has
// wrapped around.
func TestKnownRounds_Truncate_Wrap_Around(t *testing.T) {
	kr := KnownRounds{
		bitStream:      uint64Buff{math.MaxUint64, 0, math.MaxUint64, 0},
		firstUnchecked: 320,
		lastChecked:    390,
		fuPos:          1,
	}

	newKR := kr.Truncate(330)

	if newKR.firstUnchecked != 383 {
		t.Errorf("Failed to truncate. First unchecked not migrated correctly."+
			"\nexpected: %d\nreceived: %d", 383, newKR.firstUnchecked)
	}

	krBytes := kr.Marshal()
	newKrBytes := newKR.Marshal()

	if len(newKrBytes) >= len(krBytes) {
		t.Errorf("Marshalled truncated KR larger than original."+
			"\nexpected: %d\nrecived: %d", len(krBytes), len(newKrBytes))
	}
}

// Simulate saving and reading from the database by:
// 1. make random edits to the KnownRounds
// 2. save after each random edit (KnownRounds.OutputBuffChanges)
// 3. reconstructs the KnownRounds from the saved data (NewFromParts)
// 4. compare the original KnownRounds to the reconstructed KnownRounds
func TestKnownRounds_Database_Simulation(t *testing.T) {
	prng := rand.New(rand.NewSource(42))
	n := 255

	kr := &KnownRounds{
		bitStream:      makeRandomUint64Slice(n, prng),
		firstUnchecked: 5,
		lastChecked:    id.Round(n * 64),
		fuPos:          5,
	}

	saved := kr
	var err error
	var changes KrChanges

	for i := 0; i < 100; i++ {
		t.Logf("%d  %v", i, kr)
		// Modify random round
		kr.Check(id.Round(prng.Int63n(int64(kr.lastChecked))))
		t.Logf("%d  %v", i, kr)

		// Save changes
		changes, saved.firstUnchecked, saved.lastChecked, saved.fuPos, err =
			kr.OutputBuffChanges(saved.bitStream)
		if err != nil {
			t.Errorf("Failed to output changed (%d): %+v", i, err)
		}

		// Apply changes to saved KnownRounds
		for j, word := range changes {
			saved.bitStream[j] = word
		}

		// Reconstructs the KnownRounds from the saved data
		newKR := NewFromParts(saved.bitStream,
			saved.firstUnchecked, saved.lastChecked, saved.fuPos)

		// Compare the original KnownRounds to the reconstructed KnownRounds
		if !reflect.DeepEqual(kr, newKR) {
			t.Errorf("Reconstructed KnownRounds does not match original."+
				"\nexpected: %v\nreceived: %v", kr, newKR)
		}
	}
}

func makeRandomUint64Slice(n int, prng *rand.Rand) []uint64 {
	uints := make([]uint64, n)
	for i := range uints {
		uints[i] = prng.Uint64()
	}
	return uints
}

func makeRange(min, max int) []id.Round {
	a := make([]id.Round, max-min+1)
	for i := range a {
		a[i] = id.Round(min + i)
	}
	return a
}

func TestKnownRounds_Len(t *testing.T) {
	kr := NewKnownRound(0)

	decodeString := []byte{
		174, 69, 206, 0, 0, 0, 0, 0, 150, 73, 206, 0, 0, 0, 0, 0, 2, 1, 0, 136}

	err := kr.Unmarshal(decodeString)
	if err != nil {
		t.Errorf("Failed to unmarshal: %+v", err)
	}
}
