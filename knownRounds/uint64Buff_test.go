///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package knownRounds

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

// Happy path of get().
func Test_uint64Buff_get(t *testing.T) {
	// Generate test positions and expected value
	testData := []struct {
		pos   int
		value bool
	}{
		{0, false},
		{64, true},
		{319, false},
		{320, false},
	}
	u64b := uint64Buff{0, ones, 0, ones, 0}

	for i, data := range testData {
		value := u64b.get(data.pos)
		if value != data.value {
			t.Errorf("get() returned incorrect value for bit at position %d (round %d)."+
				"\n\texpected: %v\n\treceived: %v", data.pos, i, data.value, value)
		}
	}
}

// Happy path of set().
func Test_uint64Buff_set(t *testing.T) {
	// Generate test positions and expected buffers
	testData := []struct {
		pos  int
		buff uint64Buff
	}{
		{0, uint64Buff{0x8000000000000000, ones, 0, ones, 0}},
		{64, uint64Buff{0, ones, 0, ones, 0}},
		{320, uint64Buff{0x8000000000000000, ones, 0, ones, 0}},
		{15, uint64Buff{0x1000000000000, ones, 0, ones, 0}},
	}

	for i, data := range testData {
		u64b := uint64Buff{0, ones, 0, ones, 0}
		u64b.set(data.pos)
		if !reflect.DeepEqual(u64b, data.buff) {
			t.Errorf("Resulting buffer after setting bit at position %d (round %d)."+
				"\n\texpected: %064b\n\treceived: %064b"+
				"\n\033[38;5;59m               0123456789012345678901234567890123456789012345678901234567890123 4567890123456789012345678901234567890123456789012345678901234567 8901234567890123456789012345678901234567890123456789012345678901 2345678901234567890123456789012345678901234567890123456789012345 6789012345678901234567890123456789012345678901234567890123456789 0123456789012345678901234567890123456789012345678901234567890123"+
				"\n\u001B[38;5;59m               0         1         2         3         4         5         6          7         8         9         0         1         2          3         4         5         6         7         8         9          0         1         2         3         4         5          6         7         8         9         0         1          2         3         4         5         6         7         8"+
				"\n\u001B[38;5;59m               0         0         0         0         0         0         0          0         0         0         1         1         1          1         1         1         1         1         1         1          2         2         2         2         2         2          2         2         2         2         3         3          3         3         3         3         3         3         3",
				data.pos, i, data.buff, u64b)
		}
	}
}

// Tests that clearRange() clears the correct bits.
func Test_uint64Buff_clearRange(t *testing.T) {
	// Generate test ranges and expected buffer
	testData := []struct {
		start, end int
		buff       uint64Buff
	}{
		{0, 63, uint64Buff{1, ones, ones, ones, ones}},
		{0, 64, uint64Buff{0, ones, ones, ones, ones}},
		{0, 65, uint64Buff{0, 0x7FFFFFFFFFFFFFFF, ones, ones, ones}},
		{0, 319, uint64Buff{0, 0, 0, 0, 1}},
		{0, 320, uint64Buff{0, 0, 0, 0, 0}},
		{1, 318, uint64Buff{0x8000000000000000, 0, 0, 0, 3}},
		{1, 330, uint64Buff{0, 0, 0, 0, 0}},
		{0, 1200, uint64Buff{0, 0, 0, 0, 0}},
		{0, 400, uint64Buff{0, 0, 0, 0, 0}},
		{36, 354, uint64Buff{0x30000000, 0, 0, 0, 0}},
		{0, 0, uint64Buff{ones, ones, ones, ones, ones}},
		{0, 1, uint64Buff{0x7FFFFFFFFFFFFFFF, ones, ones, ones, ones}},
		{5, 27, uint64Buff{0xF800001FFFFFFFFF, ones, ones, ones, ones}},
		{5, 110, uint64Buff{0xF800000000000000, 0x3FFFF, ones, ones, ones}},
		{310, 5, uint64Buff{0x7FFFFFFFFFFFFFF, ones, ones, ones, 0xFFFFFFFFFFFFFC00}},
	}

	for i, data := range testData {
		u64b := uint64Buff{ones, ones, ones, ones, ones}
		u64b.clearRange(data.start, data.end)
		if !reflect.DeepEqual(u64b, data.buff) {
			t.Errorf("Resulting buffer after clearing range %d to %d is incorrect (round %d)."+
				"\n\texpected: %064b\n\treceived: %064b"+
				"\n\033[38;5;59m               0123456789012345678901234567890123456789012345678901234567890123 4567890123456789012345678901234567890123456789012345678901234567 8901234567890123456789012345678901234567890123456789012345678901 2345678901234567890123456789012345678901234567890123456789012345 6789012345678901234567890123456789012345678901234567890123456789 0123456789012345678901234567890123456789012345678901234567890123"+
				"\n\u001B[38;5;59m               0         1         2         3         4         5         6          7         8         9         0         1         2          3         4         5         6         7         8         9          0         1         2         3         4         5          6         7         8         9         0         1          2         3         4         5         6         7         8"+
				"\n\u001B[38;5;59m               0         0         0         0         0         0         0          0         0         0         1         1         1          1         1         1         1         1         1         1          2         2         2         2         2         2          2         2         2         2         3         3          3         3         3         3         3         3         3",
				data.start, data.end, i, data.buff, u64b)
		}
	}
}

// Tests that copy() copies the correct bits.
func Test_uint64Buff_copy(t *testing.T) {
	// Generate test ranges and expected copied value

	const numTests = 100
	const maxBuffSize = 100

	prng := rand.New(rand.NewSource(42))

	for i := 0; i < numTests; i++ {
		lenBuf := 0
		for lenBuf == 0 {
			lenBuf = int(prng.Uint32() % maxBuffSize)
		}

		buf := make(uint64Buff, lenBuf)
		for j := 0; j < lenBuf; j++ {
			buf[j] = prng.Uint64()
		}
		subsampleStart, subsampleEnd := 0, 0
		for subsampleEnd-subsampleStart == 0 {
			subsampleStart = int(prng.Uint64() % uint64(lenBuf*64))
			subsampleDelta := int(prng.Uint64() % (uint64(lenBuf*64 - subsampleStart)))
			subsampleEnd = subsampleStart + subsampleDelta
		}

		// delta := subsampleEnd-subsampleStart

		copied := buf.copy(subsampleStart, subsampleEnd)

		// check edge regions
		for j := 0; j < subsampleStart%64; j++ {
			if !copied.get(j) {
				t.Errorf("Round %v position %v < substampeStart %v(%v) is "+
					"false when should be true", i, j, subsampleStart, subsampleStart%64)
			}
		}
		// dont test the edge case where the last element is the last in the
		// last block because nothing will have been filled in to test
		if (subsampleEnd/64 - subsampleStart/64) != len(copied) {
			for j := subsampleEnd % 64; j < 64; j++ {
				if copied.get(((len(copied) - 1) * 64) + j) {
					t.Errorf("Round %v position %v (%v) > substampeEnd %v(%v) is "+
						"true when should be false", i, ((len(copied)-1)*64)+j, j,
						subsampleEnd, subsampleEnd%64)
				}
			}
		}
		// check all in between bits are correct
		for j := subsampleStart % 64; j < subsampleEnd-subsampleStart; j++ {
			if copied.get(j) != buf.get(j+(subsampleStart/64)*64) {
				t.Errorf("Round %v copy position %v not the same as original"+
					" position %v (%v + %v)", i, j%64, (j+subsampleStart)%64,
					subsampleStart, j)
			}
		}
	}
}

// Happy path of convertLoc().
func Test_uint64Buff_convertLoc(t *testing.T) {
	// Generate test position and expected block index and offset
	testData := []struct {
		pos         int
		bin, offset int
	}{
		{0, 0, 0},
		{5, 0, 5},
		{63, 0, 63},
		{64, 1, 0},
		{127, 1, 63},
		{128, 2, 0},
		{319, 4, 63},
		{320, 0, 0},
	}

	u64b := uint64Buff{0, 0, 0, 0, 0}

	for i, data := range testData {
		bin, offset := u64b.convertLoc(data.pos)
		if bin != data.bin || offset != data.offset {
			t.Errorf("convert() returned incorrect values for position %d "+
				"(round %d).\n\texpected: bin: %3d  offset: %3d"+
				"\n\treceived: bin: %3d  offset: %3d",
				data.pos, i, data.bin, data.offset, bin, offset)
		}
	}
}

// Happy path of convertEnd().
func Test_uint64Buff_convertEnd(t *testing.T) {
	// Generate test position and expected block index and offset
	testData := []struct {
		pos         int
		bin, offset int
	}{
		{0, 0, 0},
		{5, 0, 5},
		{63, 0, 63},
		{64, 0, 64},
		{65, 1, 1},
		{127, 1, 63},
		{128, 1, 64},
		{319, 4, 63},
		{320, 4, 64},
	}

	u64b := uint64Buff{0, 0, 0, 0, 0}

	for i, data := range testData {
		bin, offset := u64b.convertEnd(data.pos)
		if bin != data.bin || offset != data.offset {
			t.Errorf("convert() returned incorrect values for position %d "+
				"(round %d).\n\texpected: bin: %3d  offset: %3d"+
				"\n\treceived: bin: %3d  offset: %3d",
				data.pos, i, data.bin, data.offset, bin, offset)
		}
	}
}

// Tests happy path of getBin().
func Test_uint64Buff_getBin(t *testing.T) {
	// Generate test block indexes and the expected index in the buffer
	testData := []struct {
		block       int
		expectedBin int
	}{
		{0, 0},
		{4, 4},
		{5, 0},
		{15, 0},
		{82, 2},
	}

	u64b := uint64Buff{0, 0, 0, 0, 0}
	for i, data := range testData {
		bin := u64b.getBin(data.block)
		if bin != data.expectedBin {
			t.Errorf("getBin() returned incorrect block index for index %d "+
				"(round %d).\n\texpected: %d\n\treceived: %d",
				data.block, i, data.expectedBin, bin)
		}
	}
}

// Tests that delta() returns the correct delta for the given range.
func Test_uint64Buff_delta(t *testing.T) {
	// Generate test ranges and the expected delta
	testData := []struct {
		start, end    int
		expectedDelta int
	}{
		{0, 0, 1},
		{5, 5, 1},
		{170, 170, 1},
		{670, 670, 1},
		{63, 64, 1},
		{0, 63, 1},
		{0, 64, 1},
		{0, 65, 2},
		{5, 35, 1},
		{5, 75, 2},
		{0, 75, 2},
		{0, 319, 5},
		{0, 400, 7},
		{35, 354, 6},
		{63, 354, 6},
		{45, 5, 6},
		{130, 5, 4},
		{230, 65, 4},
		{310, 5, 2},
		{310, 64, 2},
		{310, 65, 3},
	}

	u64b := uint64Buff{ones, ones, ones, ones, ones}

	for i, data := range testData {
		delta := u64b.delta(data.start, data.end)
		if delta != data.expectedDelta {
			t.Errorf("delta() returned incorrect value for range %d to %d (round %d)."+
				"\n\texpected: %d\n\treceived: %d",
				data.start, data.end, i, data.expectedDelta, delta)
		}
	}
}

// Tests that bitMaskRange() produces the correct bit mask for the range.
func Test_bitMaskRange(t *testing.T) {
	// Generate test ranges and the expected mask
	testData := []struct {
		start, end   int
		expectedMask uint64
	}{
		{0, 0, 0b1111111111111111111111111111111111111111111111111111111111111111},
		{63, 63, 0b1111111111111111111111111111111111111111111111111111111111111111},
		{0, 65, 0b0000000000000000000000000000000000000000000000000000000000000000},
		{5, 25, 0b1111100000000000000000000111111111111111111111111111111111111111},
		{15, 15, 0b1111111111111111111111111111111111111111111111111111111111111111},
		{32, 62, 0b1111111111111111111111111111111100000000000000000000000000000011},
		{62, 32, 0b0000000000000000000000000000000011111111111111111111111111111100},
		{62, 32, 0b0000000000000000000000000000000011111111111111111111111111111100},
		{5, 65, 0b1111100000000000000000000000000000000000000000000000000000000000},
		{75, 85, 0b0000000000000000000000000000000000000000000000000000000000000000},
		{65, 65, 0b0000000000000000000000000000000000000000000000000000000000000000},
	}

	for i, data := range testData {
		testMask := bitMaskRange(data.start, data.end)
		if testMask != data.expectedMask {
			t.Errorf("Generated mask for range %d to %d is incorrect (round %d)."+
				"\n\texpected: %064b\n\treceived: %064b"+
				"\n              0123456789012345678901234567890123456789012345678901234567890123"+
				"\n              0         1         2         3         4         5         6",
				data.start, data.end, i, data.expectedMask, testMask)
		}
	}
}

// Happy path.
func TestUint64Buff_marshal_unmarshal(t *testing.T) {
	testData := []struct {
		buff uint64Buff
	}{
		{uint64Buff{1}},
		{uint64Buff{0x7FFFFFFFFFFFFFFF}},
		{uint64Buff{1, ones, ones, ones, ones}},
		{uint64Buff{0, ones, ones, ones, ones}},
		{uint64Buff{0, 0x7FFFFFFFFFFFFFFF, ones, ones, ones}},
		{uint64Buff{0, 0x7FFFFFFFFFFFFFFF, 0x7FFFFFFFFFFFFFFF, ones, ones}},
		{uint64Buff{0, 0, 0, 0, 1}},
		{uint64Buff{0, 0, 0, 0, 0}},
		{uint64Buff{0x8000000000000000, 0, 0, 0, 3}},
		{uint64Buff{0, 0, 0, 0, 0}},
		{uint64Buff{0, 0, 0, 0, 0}},
		{uint64Buff{0, 0, 0, 0, 0}},
		{uint64Buff{0x30000000, 0, 0, 0, 0}},
		{uint64Buff{ones, ones, ones, ones, ones}},
		{uint64Buff{0x7FFFFFFFFFFFFFFF, ones, ones, ones, ones, 0x13374AFB434FF, 0, 0, 0, 0x5}},
		{uint64Buff{0xF800001FFFFFFFFF, 0xF800001FFFFFFFFF, ones, ones, ones, ones}},
		{uint64Buff{0xF800001FFFFFFFFF, ones, ones, ones, ones}},
		{uint64Buff{0xF800000000000000, 0x3FFFF, ones, ones, ones}},
		{uint64Buff{0x7FFFFFFFFFFFFFF, ones, ones, ones, 0xFFFFFFFFFFFFFC00}},
	}

	for i, data := range testData {
		buff := data.buff.marshal()
		u64b := unmarshal(buff)
		if !reflect.DeepEqual(data.buff, u64b) {
			t.Errorf("Failed to marshal and unmarshal buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, data.buff, u64b)
		}
	}
}

func TestKnownRounds_Marshal2(t *testing.T) {
	data := `{"BitStream":"////v6jgxPkU///v////////Af/3/////////wH/3//t//////8B//////////3/Af/////r8////wH/////v//f//8B/////////f3/Af///////////wEB//v/////////Ae/////f////+wH///////3///8B//v////3v///Af7/9//7//7//wH///////2///8B////////////AQG///2////+//8B7/////3/////Afv////////7/wG///////////8B//ef////////Af//3/////3//wH/////9/+//v8B////////////AQH////f////+/cB/f//////3///Af//////////3wH////////3//8B///q/f//////Af///e///////wH///////////8BA///+///7////QH///////7///8B//////3/////Af////3//////wH/7/////7///8B//////3/////Ab///////////wH///////////8BAf/3/////////wH/5/////////8B+/////////+3Af/f/////////wH//9///+////8B/////7////f/Af//v/f//////wH/+/////////8B/7//////////Af/f9////////wH///////////8BAf//v////7///wHf//////////8B/////////7//Af//////////7wH////////f//8B////////////AQH////////3/v8B//+///////f/Af/////+/////wH/7/////f+//8B//+////7//7/Af/////////9/wH/////////9/8Bv////7//////Af//3////////gH////3//+//v8B/7v//////+//Af///////////wEB//7/9/f/////Af////////7//wH//9//9/+///8B/////////+//Af///////////wEB///////19v//Af/97////7///wH//////f///f8B//v///7//9//Af///////////wEB///////////7Af/////////P/wH//////9///7sB////7/v///9e//f////////zAev///+//////wH////+//////8B///v////////Af///////////wEB//////+7////Af//////+///vwH///////////8BAr/////3//v//wH///+///////8B/f////v/////Af////7//9//2gH/////3/////8B////v/v/////Af/////3//f//wH/7/////////8B////+//+7///Af//3////////wH////////7//8B/9/3//////9///3/////////Af/+/////////wH7//f///////8B///7/9//////Af///f//+////wH//7//////7/8B///r///1////Af///////////wEB///f////////Ab///97//////wH///////////8BAf/////7/f///wH//f///9/7/78B//7////7/7//Af///////////wEB/////f+/////Af///////////wEB///+v//////7Af/////+3////wH///////////8BAv+///////fv/wH///7/7////+8B////////////AQH//v////////8B///////en//3Af///////////wEB/////9//v///Af3/////3////wH///////////8BAf//97//u////wH/////+/////8B/7//////////Af///////////wEB//v//////f/7Af//////96///wH/////9/////8B1/v3//f/////Af/////////7/wH///f9//////8B///v/////7//Af//////3/3//wH///////////8BAff////////97wH//9////+//+8B////////////AQH9//////////8B////////7///Af///f/23//9zwH///////////8BAf/3////v////wH///v//f////8B//////7/////Af///////////wEB//////7/////Af///+///////wH//7/7///7//8B/////////f//Af+//f///////wH///////////cB/9//////////Af/7///7////vwH/////9////f8B/+//////////Af//////9///9QH//7////////8B////////////AQH93//3/////P8B9/v/9//////fAf/f////t///9wH///////f///8B//////////v/Af//////v////wH///////////8BAf////////+//gHf/////////f8B/+//////////Af+/////v////wH////+//////8B///9///////vAf///v///////wH3//////////8B//v/////////Af/9/////////wH7/7//3/////8B/////////v//Ad/////v/v2//wH////////3//8B///////////PAf//9///9////wH/////+////v8B//3/////////Ae//97/////3/wH////////9//8B/////////+//Af///////+7/vwHf////9/////8B/////////f//Af/7////////7gH///////f///8B///+//3/////Af//////////7wH///////////8BAf///e///////wH///////f///0B////3///////Af///6+/////9wH///7/////v/sB///v///3////Aff//////f///wH//fv//9////8B////////////AQH/////+v/9//4B////////////AQL+/7////////8B////3///7///Af///f//3////wH/9/////3/9/8B/////////9//Af///f///////wH/9/////+///8B//////3/4///Af//+///9////wG9+////////+8B///////////3Af//9/+//////wH////3//////8B/////7/3////Af/f///X///f7wH///////////8BAf////+//v//9wH//7//////7/8B//7//////f/uAf/////f/////wH//f/////+//8B///////37/v/Af///f///////QH//////////78Bv///////////Af/nv///////vwH//////f/3//8B//3/////////Af///////7///wH7//////////cB3//////v////Af//3////////wG/v//////r//8B////////v///Af///////9///wH/+/7//v/v//8B////////////AQH/3u////////8B////////////AQH//+///f////8B/////////9//Ad///////////wH/3//3//////8B////////////AQH////9///3//8B9///////////Af/////////3/wH///////////8BAfv/+////////wH//////9//v/8B+v//////////Aff//////+//3wH///3f+//9//4B///f////////Ab///f/////z/wH///////////8BAf/////////3/wH///////+///8B////////////AQH//////////78B3////f3///v/Af//////////+wH89//e//////8B////////s///Af///////f//+wH7////v/////8B3/////2/////Af//9////////wH/v///////9/8B//+//7////3/Af//3/7/v//u/wH///////////8BAf/////777///wH7/////v///98B////////////AQH/////7/7//fsB//f///////f/Af/////////f/wG+/////9////8B/v///v7///9//9/////+////Af+///v//////wH///////+///8B////////////AQH/+/////////8B////////3///Af//////3////wH/////+/////8Bv//////////vAf////////3//wHu/9////////8B///////v//9/3/////////7/Af///////////wEB////////7///Af///////////wEB///3////////Af/////////7/wH//////////98B////////////AQH////3//3///8B////////////AQH//////v////8B/b/3////////Af///////////wEB/////////9//Af/+////////f+////v/////+wH///////////8BAv///9/////7/wH////d//////8B/f//////////Af///////////wEC//////////v/Af//9/////v95wH/////////3/8B//////v////fAf/u7v///9//+wH///////////8BAd//////9////wH///////////4B//////3//f27Af///////////wEB///9/+v//f/+Af/////73////wH///////////8BAf/9/////////wH//////v////8B/////v//////Af//+////////wH/v/////73//8B/7///+//////Af//7/v//9//7wH///+///////8B////////////AQH/3//7//////8B////7//+//+/Af//8//7v//7/wH///37//////8B//////////f/Af//////9//v+wHf///+///3//8B////////7//fAf/v/////////wH/////9/////8B////////////AQL////v//3///0B////////////AQH/////+//v//8B/////////+//Af/f//f//////AH+/f3///v/3/4B////////3///Af3///f/5////wH////9//////8B//+////////fAdf//////////wH///////////8BA///////////+wH7/v////////8B/////////v//Af////2//////wH/////3////zf//f/3/////98B/P////////+/Af////f//////wH//////////98B///f////////Af//+////////gHf///////+//8B////////7/P/Af//////+////wH///////////8BAf//v////////wH///////7///8B///v/f/////9Af///////9///QH//////////vsB/9//v///////Af///+/9///9/wH///////////cB/t7/3v//////Af/////////v/wH///////7///0B/////f//////Af//////+////wH///+7////3/8B/7/3+//93///Ab/v////v7///wH//v////////4B///f/9//////Af7v/7///f///wH7/////////v8B///7////////Af//9////////wH///////+///8B////////////AQH///+///3///8B////////////AQH/7////9f///8B///99////f/9Af///////+///wH///////////8BAf///9////77/wH///////////8BAf/f/////////wH/vb////3//78B////m////v3/Af//7/////7//wH///f///////8B9///7///////Af///////////wEB//////7/////Af/7///////v/QH/+//7//////8B//+9///9+ve/Af////f9/////wH///////////8BAf/f//+//9///gH+/9//3/////8B//+///v////3Af///v///f//nwH/+/////////4B/+/////////vAf///////////wEB9/+f/////+//Af/////////7+gH///////////8BAf///d/7///v/wH///////////8BAf/e/////9///QH///////////8BAv///+//v/7//wH+///7//////8B////////9v//Af///////////wEB///f///7///vAf////////+//wH/////v////+8B//////////f/Ab//7/7///+//wG///////7/+/8B///9////////Af/f//f////p/wGf/9//v/////8B/7//////////Af/f/v//+////wH7////////9/8B////////////AQH///////7///8B3e///f//////Af/7/////////wH///////////8BAf//n/7//////wH/vf////7///8B/////////f//Af///////////wEC/7/9/v//////Af//3//fvv/e9wH+////+/////8B//////7/////Af///f///////wG/////v//77/8B//7/////v///Af////+9///7/wH/////////+/8B///8////7/9////7/+//9///Af///////////wEB///////+/9//Af///////f///wH///f//+////cB////////////AQH9///////+//8Bv//7//////r8Af/+///7/////wH/3/////b///8B//////+/////Af/f/////////wH////9//7/3f8B/////9///9d///////////9/////3/////z9Af//7//3///7/wH///////////0Bv///////////Af///////////wEB3////////7//Af3///+//9/9X//3/////////wH//e///////f8B////v//9//+/Aff////d///9/wH//////v37//8B7///////v///Af//////////f/////f////+/wHv///f///9//8B//////7/////Af///v/v/////wH///////////8BAv/////+/////wH///////////4B////////////AQL/+/v/////v/8Bv//x9//77///Af7//7///9/7/wG//t///////3/+///////9//8B///v//fv////Af///////////wEC///////9////Af/+/8f//////wH///////////8BAf3//v/+//3/7wH///////////8BAf///////7/99wH////f//////8B///9/////v//Af//v/3//v///wH/////v/////8B///////+////Afz/////v////wH/////+/7///8B/////////7//Af///////v///gH//+/f//3///8B/9/////9///6Af///////////wEB2///////////Af//v7///////wH//8eZ//3///8B//3/+v//////Af//+////////wH///v7/////v8B/v/7//////9///v//7//////Af/////////3/wH//f////////8B/////+/9////Af//v///3////wH//9///7////8B////////////AQH/////9/////8B//3//////+//Af///////////wEB//39////////Af//////////7wH///////////8BAf////v//////wH+////////z/8B////////////AQH////+/v////8B///9///7////Af/Y/v///+///wH////////f//8B////////////AQH/3//3v//+//8B3//////+///vAf/////7///7/wH///////////8BAf/7///d/////wG///////////8B/////f//////Af///////////wEB//+///vf////Af///v/+///f/wH//9////f///8B///s///93///Ab/////////3/wH///////+///8Bv///////////Af///////////wEB/9///f//////Af/////////+/wH//////r////8B////////v///Af/f/////////wH//////+f/3/8B////7//7//+/Af//////////+gH////////3//8B////v//f////Af///////////wEB///////7//+7Af///////////wEB3//9////////Af/////f//+//wH///////////sB9/v/////////Af/////+/////wHv/f////v///8B////////////AQH////7/v////sB/f//////////Af+/////////7wH/+/////////8B///+//////f/Af7X/////////wH//////v/u//8B/////////7//Af////f//////wH//////////98B////////vf//Af//////7//3/wH3//////////cB/9///d77////Af////+///v//wH////////f//8B////v//f////Af3/9//d/////wHr///3/9/v//8B/////9//////Af///////////wEB9//9///9////Ae///9//////3wH///////////8BAf/////////+/wH3//7//////v8B///9////77//Af/+/////////wH3///777//3/8B///////+7/9///////////f/Af//37f///+//wH///////////8BAf//u//////v/wH/////29//+78B//f//////v//Af/////+/////wH/9////////3////////7//3///////v////8B+/33////////Ab///+/7/v///wH///77//////8Bvv////3f///3Af//9////////wH/v/////////cB3///////////Af//+/////39/wH7//////////cB/////////9v/Af//9////P///wH////9v///+/8B///z/////v//Af/////////v/wH////fv/////4B//n////////7Af//vvv////3/wH//////P//9/8B///////v///+Af+////////v/wH///////////8BAf7v/////////wH9/////v//v78B/////+//////Af/v//z+///39wH///////////8BAf+////////v7wH/////9/////8B/////7//3//9Af/////////v7wH///////////8BAf/33///9////wH//////7///74B/////7f/////Af/3///////+/wH///////////8BAf/////9/v///wH////y//////8B///+///n////Af//9v////f+/wH////////v3/8B/////////7//Af/++7///f///wH///////////8BAvf3/////////wH////v//////8B9///5///v//fAf/////f/f///wH///////////8BAf//+///v///7wH///ff9////3//////////3/4B////v/+/////Af/f/////+/3/wH///////////sB/////////9//Af///v//+v///wH3//////////8B///X//v///f/Aff/////////f/////26/////wH////////+//8B/f//////v///Af////f//+///wH2/////9////0B//7//////9++Af///////t///wH///////2///8B///e///f/v//Af///////////wEB/7//////////Af////////+//wHf//////////8B/v//v+//////Af//////+////wH7//////vv/78B////////////AQH////vv+//13/////P3/////8B+///////+///Af///////////wEB///f//+/////Afn///////v+/wH////X/f//v78B/9////7/////Af6///////f/9wH//////////33//u///7////8B//vz//u/////Af//////7////wH///////7/+/8B//////////7/Af//v//////9/wH///////v///8B3///7///3//5Af+//////////wH////9/7////8B/////////v//Af///f3/////vwH7//////////8B////v/3///79Af///////////wEB////x//////zAf//////v////wH7/////7////8B////////////AQH///7v9//fv98B/////////v9/7///////v///Af7//f///9/+/wH/nffZ////9/oBlu+Yvvvl6eo76///v8CogJB0","FirstUnchecked":94919909,"LastChecked":94969701}`
	kr := &KnownRounds{}

	err := kr.Unmarshal([]byte(data))
	if err != nil {
		t.Errorf("Unmarshal() returned an error: %+v", err)
	}

	t.Log(kr)
	t.Logf("%064b", kr.bitStream)

	t.Log(kr.Checked(94969696))
}

// printBuff prints the buffer and mask in binary with their start and end point
// labeled.
func printBuff(buff, mask uint64Buff, buffStart, buffEnd, maskStart, maskEnd int) {
	fmt.Printf("\n\u001B[38;5;59m         0         0         0         0         0         0         0          0         0         0         1         1         1          1         1         1         1         1         1         1          2         2         2         2         2         2          2         2         2         2         3         3          3         3         3         3         3         3         3" +
		"\n\u001B[38;5;59m         0         1         2         3         4         5         6          7         8         9         0         1         2          3         4         5         6         7         8         9          0         1         2         3         4         5          6         7         8         9         0         1          2         3         4         5         6         7         8" +
		"\n\033[38;5;59m         0123456789012345678901234567890123456789012345678901234567890123 4567890123456789012345678901234567890123456789012345678901234567 8901234567890123456789012345678901234567890123456789012345678901 2345678901234567890123456789012345678901234567890123456789012345 6789012345678901234567890123456789012345678901234567890123456789 01234567890123456789012345678901234567890123456789012345678901\n")
	fmt.Printf("        %*s%*s\n", buffStart+2+buffStart/64, "S", (buffEnd+2+buffEnd/64)-(buffStart+2+buffStart/64), "E")
	fmt.Printf("buff:   %064b\n", buff)
	fmt.Printf("mask:   %064b\n", mask)
	fmt.Printf("        %*s%*s\n", maskStart+2+maskStart/64, "S", (maskEnd+2+maskEnd/64)-(maskStart+2+maskStart/64), "E")
}
