///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package knownRounds

import (
	"encoding/base64"
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
	testData := []uint64Buff{
		{1},
		{0x7FFFFFFFFFFFFFFF},
		{1, ones, ones, ones, ones},
		{0, ones, ones, ones, ones},
		{0, 0x7FFFFFFFFFFFFFFF, ones, ones, ones},
		{0, 0x7FFFFFFFFFFFFFFF, 0x7FFFFFFFFFFFFFFF, ones, ones},
		{0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0},
		{0x8000000000000000, 0, 0, 0, 3},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0x30000000, 0, 0, 0, 0},
		{ones, ones, ones, ones, ones},
		{0x7FFFFFFFFFFFFFFF, ones, ones, ones, ones, 0x13374AFB434FF, 0, 0, 0, 0x5},
		{0xF800001FFFFFFFFF, 0xF800001FFFFFFFFF, ones, ones, ones, ones},
		{0xF800001FFFFFFFFF, ones, ones, ones, ones},
		{0xF800000000000000, 0x3FFFF, ones, ones, ones},
		{0x7FFFFFFFFFFFFFF, ones, ones, ones, 0xFFFFFFFFFFFFFC00},
	}

	for i, data := range testData {

		buff := data.marshal()
		u64b, err := unmarshal(buff)
		if err != nil {
			t.Errorf("unmarshal produced an error: %+v", err)
		}
		if !reflect.DeepEqual(data, u64b) {
			t.Errorf("Failed to marshal and unmarshal 1 byte buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, data, u64b)
		}
	}
}

// Happy path.
func TestUint64Buff_marshal_unmarshal_Bytes(t *testing.T) {
	testData := []uint64Buff{
		{1},
		{0x7FFFFFFFFFFFFFFF},
		{1, ones, ones, ones, ones},
		{0, ones, ones, ones, ones},
		{0, 0x7FFFFFFFFFFFFFFF, ones, ones, ones},
		{0, 0x7FFFFFFFFFFFFFFF, 0x7FFFFFFFFFFFFFFF, ones, ones},
		{0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0},
		{0x8000000000000000, 0, 0, 0, 3},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0x30000000, 0, 0, 0, 0},
		{ones, ones, ones, ones, ones},
		{0x7FFFFFFFFFFFFFFF, ones, ones, ones, ones, 0x13374AFB434FF, 0, 0, 0, 0x5},
		{0xF800001FFFFFFFFF, 0xF800001FFFFFFFFF, ones, ones, ones, ones},
		{0xF800001FFFFFFFFF, ones, ones, ones, ones},
		{0xF800000000000000, 0x3FFFF, ones, ones, ones},
		{0x7FFFFFFFFFFFFFF, ones, ones, ones, 0xFFFFFFFFFFFFFC00},
	}

	fmt.Printf("%4s   %4s   %4s   %4s   %4s\n", "orig", "1B", "2B", "4B", "8B")
	fmt.Println("==================================")
	for i, data := range testData {

		buff := data.marshal1Byte()
		u64b := unmarshal1Byte(buff)
		f1bLen := len(buff)
		if !reflect.DeepEqual(data, u64b) {
			t.Errorf("Failed to marshal and unmarshal 1 byte buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, data, u64b)
		}

		buff = data.marshal2Bytes()
		u64b = unmarshal2Bytes(buff)
		f2bLen := len(buff)
		if !reflect.DeepEqual(data, u64b) {
			t.Errorf("Failed to marshal and unmarshal 2 bytes buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, data, u64b)
		}

		buff = data.marshal4Bytes()
		u64b = unmarshal4Bytes(buff)
		f4bLen := len(buff)
		if !reflect.DeepEqual(data, u64b) {
			t.Errorf("Failed to marshal and unmarshal 4 bytes buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, data, u64b)
		}

		buff = data.marshal8Bytes()
		u64b = unmarshal8Bytes(buff)
		f8bLen := len(buff)
		if !reflect.DeepEqual(data, u64b) {
			t.Errorf("Failed to marshal and unmarshal buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, data, u64b)
		}

		origLen := len(data) * 8
		fmt.Printf("%4d   %4d   %4d   %4d   %4d\n", origLen, f1bLen, f2bLen, f4bLen, f8bLen)
		fmt.Printf("       %4.0f%%  %4.0f%%  %4.0f%%  %4.0f%%\n",
			100-float64(f1bLen)/float64(origLen)*100,
			100-float64(f2bLen)/float64(origLen)*100,
			100-float64(f4bLen)/float64(origLen)*100,
			100-float64(f8bLen)/float64(origLen)*100)
		fmt.Println("----------------------------------")
	}
}

// Tests the compression of different marshal word sizes.
func TestUint64Buff_marshal_unmarshal_Size(t *testing.T) {
	testData := []string{
		"c21ABgAAAAC7DUEGAAAAAAF//wXv/wi//wH3/wN//wG//w23/xH3/v8D+v8L+/8L/P8H9/8H/v8Iv/8E9/8Fnt//A7//AfP/Ab//B/f/Af7/Af7/B+//CN//E/f/D+//CP7/GK//Bb//Av7/BX//BH//Af1//wm//wH3/wbv/wX7/wm//w539/8Hf/8Df/8D/f8X/v8Dv/8B/f8K9/8L7/8B9/8G9/8B+/3/Cr//Au//AW//B9//A+//BN//C/7/Ba//Be//Bf3/Cv37/wjv/x5//v8E3/8H9/8B7/8F3/8Iv/8D/f3/Ce//AX9//xb3/wLf/wP+/wPt/wLv/wK//u//Ae//Av7/Dvn/CX//C3//A79//wLf/wX3/wLf/x733/f/Aef/Avf/Bd/v/v8B+/8Q3/8D/v8D7/8j7/8Bv/8I/v8D9/8F3/8Fv/8o/u//Af7/BO/3/f8Cv/8B+3//Br//Gv7f3/8Dv/8L73//Afv/Cfv/CPf/A3/7/wH9/wP7/wK//wT+/wn3/wTv/w5//wK//wu/f/8C8/8M/f8K7/7/Br+/3/8D/f8B9/8F3/8Dv/8Gvv8I+/8I++/3/wn9/wTvv/8D/f8L+/8I9/8H/v8K/v8D9/8G+v8B3/8B9/8N3/3/Bf7/Cf3/BPf/D7/9/we//wXv/wT3/xX7/wO//wL7/wfv/w33/f8Cf/8Lf/8C3/8Gd/8Cf/8lf/8Bf3//D/7/Ab//Gf3f/wG/f/8D+8v/A/3/Ge+//wH3/wPb/wT9/wb3/wI2/wF//wbf/wj9/wK//wHf/wvv/w7+/wTv+/8B9/8G9/8G/f8Cv/8B/v8B+/8L7/8B/v8E98//Ct//Bff/Bb//A/v/Bbv/Cvf/Aff/At8//w/+/wW//wH+/wTf/wP7/wnv/wP3/wl/f/8Hf/8E+9//Dvz+/wX3/wPv/xd//wL7/wf9/wL7/wS//wTf/wN//wJ//wK//wN/+/f/Hff/An/v/wZ//wX3/wf7/ff/GX//Cd//Av7+/wX7/wF//wj1/wXv/v8C9/8G/v8Cz/8Cn/8J7/8E7/8G/f8Bv/8Bf/8F/f8C+/8X7/8Bv/8B97//BPv/B7//Dd//Af3/Avb+/wLv/wXf/wG//wL75/f/Av3/CH//Bn37/wbv/wTz/wXv/xB//xD+/wjO+/8E/v8C9/8E99//A9//Bfv/FL//A/3/DOb/C/v/Env/DOf/Ed//Av7/Af37/wf+/wL8/xX9/wh//we//wHv/xL9/wz9/wN//wS//wLr/wp//wO//wv3/wF//wXv/wH3/wZ//xP5/w+//wb9/y73/w79/wH7+/8D93//Af3/Ar/+9/f/BX//Av7/Avf/Ad//Cfv/Ef73/wv3/xrf/wK//wGfv/8C+/f/Cff/Bvf/BP3/Afd//wPf/xDv/wN//wHv/v8N9f8C/f3/A/7/Bd//Bff/Dft3/wr+/wnf/wXf/w7+/wP7/wLrf/8C7/8B3/8D+/8Hv/8E3/8B8/8Hv/8F+/3/Avv/A/7/Dd/2/wfv/w5//wHv/wfv+/8M/f8GX/8H3/v/Au/6/w/v/wV/v/8E+/8M7/8L+/8B/f8Fv/8MP/8Bv/8K/v8C8/8B9/8L/f8D3/8C7/8S/v8I/P8M99//Avf/Af3/Bvv/Bf77/wX3/wL+/wP5/wHf/wr+/wG//wt//wTv/wJ//wP3/wd//wH9/wPv/wLv/wf+/v8G/X//Ae//E/7/Ab/v/xPt/v8D9/8Bv/8B7/1//wb9/wW//wff/wa/f3//A7//Bt//BPf/Dv3/DX/3/w3f/wPf/w7v/wjv/wb3/wP7/xv9/wL+/wu//wW//wh//wzv/xn+/wnf/wH+/wV//wHv/wL9/wF//v8Qv/8K7/8K3/8G7/8G9/8K/v8O33//Au//AX//CP77/we//wu//wj+f/8Df/8C+/8M9/8L/v8D9/f/Ee//CPfv/wr2/wbf/yfn/wr+/wp7/wOf/wi3/wjf/wPf/wjf+/8F9/8I/f8Ev/8D7/8E7/8Cf/8C3/8G9/8Jv/8D+/8F9/8O3/8D7/8E/v8C3/8Cv/8C79//BN//D/5/+/8Nv+//Ae7/CO//A7//Hfv/CP7/BH//Be//Afv/Af7/Ab//Cd/f/wj9/wTe/wJ//wG/P/7/Af3/Afv/Ce//Av7/Ae/5/wX7/wL9/r//Dd//Bfv+/wP3/wH6/wZ/f/8I9/8Bf/8K/f8cf/8E+/8C1/8D+/8E7/v/F/f+/wm//wf9/wL9/xn7/f8G/f8Ev/7/B/3/Bvv/Au//B7//Bv7/Dv3/C/v9/wX3/wW//wK//wm/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wj3/xDf/wG//wv7/wf+7/8Hv/8F3/8L/f8P/v8G/v8Bv/8T+6//Ad//BJf/B+//Bvv/EPf/Bfvf/w3+/wP+/w79/xDf/wH79/8F7t//Avv/A/f/B9/3/wH+/wh//xH9/wL7/wT3/wTm/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v8g3/8Ff/8F3/8B7/8Qv/8P+/8I9/8E/f8F5/8E7/8Q3f8H9/8N3/8Ev/f/Cr7/A7//Cf3/AX//C9//BPv/A/v/Af7/Ce//A+/9/wL7/wnr/wXf/wd/9/8C3/8B8f8If/7/A+//An/9/v8Bf/8D9/8N9/8G3/8ef/8Bf/8K/v8G9/8Fv/8B9/8C7/8O3/8F/v8D/f8Q3/8BP/8Dr/8G+/8C3+//Bff/Bf3/EL//Bb79/wv3/wf+/wL79/8Q/v8C+/59/wF//wFqV97stB2f",
		"zZtABgAAAAD8DUEGAAAAAAEL+/8B9/8G9/8G/f8Cv/8B/v8B+/8L7/8B/v8E94//Ct//Av7/CL//A/v/Bbv/Cvf/Aff/At8//wf9/wf+/wW//wH+/wTf/wP7/wnv/wP3/wl/f7//Bn//BHvf/w78/v8J7/8Xf/8C+/8H/f8C+/8C+/8Bv/8E3/8Df/8Cf/8Cv/8Df/8B9/8d9/8Cf+//Bn//Bff/B/v99/8P+/8Jf/8B7/8H3/8C/v8G+/8Bf/8E3/8D9f8G/v8C9/8G/v8C3/8Cn/8J7/8E7/8G/f8Bv/8Bf/8C/v8C/f8C+/8W3+//Ab//Aff/Bbv/B7//Dd//Af3/Av7+/wLv/wXf/wG//wLr5/f/Av3/D337/wbv/wF//wL7/wXv/xB//xD+/wjO+/8E/v8C9/8E99//A9//Bfv/FL//A/3/DOb/Hnv/DOf/Ed//BP37/wr8/xX9/wh//wbfv/8B7/8S9f8M/f8Df/8Ev/8C6/8Kf/8B+/8Bv/8Jf/8B9/8Bf/8F7/8B9/8Gf/8C/f8Q+f8W/f8H7/8Z/f8M9/8B/f8M/f8B+/v/A/d//wH9/wK//vf3/wV3/wL+/wL3/wHf/wn7/wv9/wX+9/8L9/8E7/8V3/8Cv/8Bn7//Avv/Cvf/BPv/Aff/Bvd//wPf/wL+/w3v/wHv/wF//wHv/v8N9f8C/f3/A/7/Bd//De//Bft//wr+/wnf/wL3/wLf/xL7/wLrf/8C7/8B3/8D+/8Hv/8G8/8N+/3/Avv/Ed/3/wfv/xDv/wfv/wP7/wn9/wZf/wff+9//Ae/6/w/v/wV/v/8R5/8L+/8B/f8Fv/8Mv/8Bv/8N8/8B9/8L/f8D3/8C7/8S/v8I/P8N3/8C9/8B/f8G+/8F/vv/Bff/Av7/A/n/Dr//C3//BO//An//A/f/B3//Be//Au//B/7+/wb9/wLv/xP+/wLv/xPt/v8D9/8Bv/8B7/8Bf/8G/f8E/f8I3/8Gv39//wO//wv3/w79/w1/9/8N3/8D3/8Df/8K7/8P9/8D+/8D9/8X/f8C/v8G/f8Ev/8Fv/8D7/8Ef/8M7/8B7/8X/v8J3/8B/v8Ff/8B7/8C/f8Bf/7/G+//Ct//Bu//Bvf/Cv7/Dt9//wLv/wF//wj+/wi//wu//wj+f/8Df/8D+/8L5/8L/v8D9/8C+/8K3/8E7/8I9+//Cvb/Bv7/HO//Cuf/Av3/B/7/A3//Bnv/A5//CL//A/v/BN//A9//Bd//A/v/Bff/CP3/CO//BO//Bd//Bvf/Cb//A/v/Bff/Dt//A+//BP7/At//Ar//Au/f/xT+f/8P7/8B7v8L97//Hfv/CP7/BH//Be//Afv/Af7/Ab//Ct//CP3/BN//An//Ab8//v8B/f8B+/8J7/8C/v8C+f8F+/8C/f8Bv/8N3/8F+/7/A/f/Afr/Bn9//wj3/wF//wr9/w6//wt//wF//wT7/wLX/wP7/wTv+/8P9/8H9/7/Cb//Bv39/wK9/xX9/wT9/wb9/wS//wj9/wF//wT7/wLv/we//wX9/v8O/f8L+/3/Bff7/wS//wK/+/8Iv3//Af3/Aff/CPv/BPv/A/7+/wj73/8K9f8C7/8C/f8Gv/8B9/8Q3/8Bv/8Hv/8D+/8I7/8E/v8Cv/8F3/8D/v8H/f8D/f8L/v7/Bf7/Ab//E/uv/wHf/wTX/wfv/wX9+/8W+9//Ef7/Dv3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/B/d//xH9/wL7/wT3/wTm/wr3/wH+/n//Aff9/wNv9f27/wTv/wb+/w2//xC//wHf/wV//wXf/wHv/xC//w/7/wj3/wT9/wXn/wTv/xDd/wf3/w3f/wX3/wK//we+/w39/wF//xD7/wG//wH7/wl//wHv/wPv/f8C+/8J7/8F3/8Hf/f/At//AfH/Cf7/A+//An/9/u//BPf/Dff/Bt//Hn//AX//Cv7/Bvf/Bb//BO//Dt//Bf7/A/3/EN9/P/8D7/8G+/8C3/8Bv/8E9/8I3/8Nv/8Fvv8U/v8C+/f/EP7/BH//BN/e+/8Bvf8BbT86JdgdkA8=",
		"+T1ABgAAAAD8DUEGAAAAAAE6u5Yc0gEPv/8I/f8Bf/8Fv/8E/f8B9/8F+/8H/n//B3//C3//Eb/9/wT93/8Hv/8C9/8Cv/ff/wP9/wPv/w7+/wp//wP9/wT3/wSv/wL7/wS7/wR//wL9/wXv/wL3/wbf/xLv/wX39/8K9/8Bv/8E9/8E/f8H9/8C9f8G/f8C/P8B9/8Q+9//CPvf/wG//wr39/8E+/8G+/8Lb/8C+/8C9/8D/f8C3/8E3/8D3/8K9/8Dv/8Bvv8D/v8K/fv/A/3/A+93/xD+/wjP/wH9f3//A/f/B/v/Bu//Ab/v/wT9/wLf/wHt/wPf/wjv/wL9/wT9/wP3/wG//wf+/wH9/wPv/ws//wT+/wPf/wF//wHs93/+/wL9/wL3/wL7/w57/wi//wNz/wZ+/wT9/wXf/wL3/wJ//wff/wHv/wb3/wh//wH3/wH9/wb+/wHv/f8Gb/8Hv/t//wP9/wj+/wV//wT+/wbv/w73/wL3/f8Bv/8Hv/3/Ae//Am+/ff39/wj+/wX+/wT9/wT+/wbP/wjv/wK//wXv3/8K+/7/AX//EPv/Av7/Eb//GPv/B/3/C3//EX9/7/8Gf/8Cv/8E773/A7//Aff/Cr//Efr/Bb//A/7/CP7/Avv/B/3/Av7+/wn3/v8B7/8Fv/v/Ad7+v/8C7/8G3/8I+f8E7/8L7/8B9/8Cv/8J3/8H+/8H/f8Cf/8Be/8Dv/8G/f8C/f8Dv/8F9/8Dx/n/Ad//Au//Bd//Bv3/Dr//AXfv/wH++/8Bf/8F+X//An/u/wO//wP7/wR//wL+u/8Df/8D+/8E9/vff/8C9/8F/u9//wJ/f/8Bv/8Dz/8Dnv8B+/8Jf/8Cv/8E2/8D/f8B3/8Uv/8D3/8Bv/8H+/8G9/8Cv/8b9/8B9/8H+/8Xf/8C1/8R7/8B7/8Bv/3+/wLv/wbf/wF//wH9/wv3/wPX/wL9/wP9/wfv/wKv3/8Bf/8C/f8I1/8C7/8G7/8G/P8R9/8E7/v/BP7/Afv/A3/+/wLf/wL+/wLv/wHf/wfv/xTv/wr7/wu//wH3/wN//wG//wT9/wi3/wW//wv3/n//Avr/C/v/C/z/B/f/B/7/CL//BPf/BZ7f/wO//wH7/wG//wf3/wH+/wH+/wfv/wL7/wXf/wXv/w33/w/v/wj+/wb9/wu//wWv/wW//wL+/wV//wL9/wF//wH9f/8E/f8Ev/8B9/8Df/8C7/8F+/8Jv/8Jv/8Ef/f/B3//A3//G/7/Bf3/Cvf/Dff/Bvf/Av3/At//B7//Au//AW//Bd//Ad//CN//C/7/Ba//Be//Bf3/Cv37/wjv/wX7/xh//v8E3/8H9/8B7/8FX/8M/f3/C39//xb3/wLf/wP+/wPt/wLv/wK//u//Ae//BN//DPn/D3//BX//A79//wH3/wb3/wLf/x733/f/Aef/Avf/Bd/v/v8B+/8Q3/8D/v8S9/8U7/8K/v8D9/8F3/8E97//KP7v/wH+/wTv9/3/Ar//Aft//wL7/xN//wvf3/8Dv/8L73//Bfv/Bfv/CPf/A3/7/wH9/wP7/f8Bv/8E/v8O7/8Ff/8If/8Cv/8Lv3//AvP/GPb/Br+/3/8D/f8B8/8Jv/8Gvv8R++//D++//wP97/8K+/8I9/8H/v8C/f8H/v8D9/8G/v8B3/8B9/8N3/3/Bf7/Dvf/Cf3/Bb/9/w3v/xr7/wb7/wfv/w33/f8Cf/8Lf/8C3/8Df/8Cd/8Cf/8nf3//D/7/Ab//EPv/CP3f/wG/f/8D+8//A/3/Ge+//wH3/wPb/wT9/wb3/wI2/wF//wbf/wj93/8Bv/8B3/8L7/8E7/8O7/v/Aff/Bvf/Bv3/Ar//Af7/Afv/Ar//CO//Af7/BPfP/wrf/wu//wP7/wW7/wr3/wH3/wLfP/8L/v8D/v8Fv/8B+v8E3/8D+/8J7/8D9/8Jf3//B3//BPvf/w78/v8J7/8Xf/8C+/8H/f8C+/8D37//BN//A3//An//Ar//A3//Aff/Hff/An/v/wZ//wX3/wf7/ff/En//Bn//BN//BN//Av7/Bvv/AX//CPX+/wX+/wL3/wb+/wLf/wKf7/8I7/8E7/8G/f8Bv/8Bf/8B/v8D/f8C+/8Gv/8Q7/8Bv/8B9/8F+/8Hv/8N3/8B/f8B/f7+/wLv/wXe/wG//wL75/f/Av3/D337/wbv/wT7/wXv/xB//w7+/wH+/wjO+/8E/v8C9/8E99//Au/f/wX7/xS//wP9/wzm/x57/wzn/xHf/wT9+/8K/P8V/f8If/8Hv/8B7/8S9f8M/f8Df/8Ev/8C6/8Kf/8Dv/8D/f8H9/8Bf/8B+/8D7/8B9/8Gf/8T6f8Q/f8F/f8u9/8O/f8B+/v/A/d//wH9/wK//vf3/wV//wL+/wL3/wGf/wn7/xH+9/8E+/8G9/8Gf/8T3/8Cv/8Bn7//Avv/Cvf/A3//Avf/Bvd//wPf/xDv/wN//wHv/v8B/f8L9f8C/f3/A/7/BP3f/xP7f/8H/v8C/v8J3/8C3/8C3/8Mv/8F+/8C63//Au//Ad//A/n/B7//BvP/Dfv9/wL7/xHf9/8H7/8Q7/8H7/8E7/8I/f8GX/8H3/v/Au/6/w/v/wV/v/8R7/8L+/8B9f8Fv/8Mvv8Bv/8N8/8B9/8L/f8D3/8C7/8S/v8I/P8F3/8H3/8C9/8B/f8G+/8F/vv/Bff/Av7/A/n/Dr//C3//BO//An//A/f/B3//Be//Au//B/7+/wb9/wLv/xP+/wLv/xPt/v8D9/8Bv/8B7/8Bf/8G/f8H3/8F3/8Gv39//wO//wG//wn3/w79/w1/9/8N3/8D3/8O7/8P9/8D+/8b/f8C/v8F/v8Fv/8Fv/8If/8M7/8Z/v8Dv/8F3/8B/v8Ff/8B7f8C/f8Bf/7/Br//FO//Ct//Bu//Bvf/Cv7/Dt9//wLv/wF//wj+/wX9/wK//wu//wj+f/8Df/8P9/8L/v8D9/8Mv/8F7/8I9+//CH//Afb/Cf3/If3/Auf/Cf3+/wp7/wOf/wO//wS//wjf/wPf/wn7/f8E9/8I/f8I7/8E7/8F3/8G9/8Jv/8D+/8F9/8O3/8D7/8E/v8C3/8Cv/8C79//FP5//w/v/wHu/wy//wS//xj7/wj+/wR//wXv/wH7/wH6/wG//wrP/wj9/wO/3/8Cf/8Bvz/+/wH9/wHz/wnv/wL+/wL5/wX7/wL9/wG//wnv/wPf/wXr/v8D9/8B+v8Gf3//CPf/AX//Cv3/CP3/E3//BPv/Atf/A/v/BO/7/xf3/v8Jv/8H/P8Cvf8a/f8G/f8Ev/8I/f8G+/8C7/8E9/8Cv/8G/v8O/f8L+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wj3/xDf/wG//wv7/wjv/we//wXf/wv9/w/+/wb+/wG//wL7/xD7r/8B3/8E1/8H7/8G+/8Bv/8O9/8F+9//Ef7/Dv3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/Bt//AX//Ef3/Avv9/wP3/wTm/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v8I3/8X3/8Ff/8F3/8B7/8Qv/8P+/8I9/8E/f8F5/8E7/8Q3f8H9/8N3/8F9/8Kvv8N/f8Bf/8Q+/8D+/8L7/8D7/3/Avn/Ce//Bd//B3/3/wLf/wHx/wn+/wPv/wJ//f7/Bff/Dff/Bt//Hn//AX//Cv7/Bvf/Bb//BO//B/v/Bt//Bf7/A/3/BL//C9//AT//A+//Bvv/At//Bvf+/xW//wW+/w37/wb+/wL79/8Q/v8Eb/8E397/Av3/Ae0/vnf4HZAP",
		"yJtABgAAAAAEDkEGAAAAAAELe/8B9/8G9/8G/f8Cv/8B/v8B+/8L7/8B/v8E98//Ct//Afv/Cb//A/v/Bbv/Cvf/Aff/At8//w/+/wW//wH+/wTf/wP7/wnv/wP3/wl/f/8Hf/8E+9//Dvz2/wnv/xd//wL7/wf9/wHf+/8Ev/8E3/8Df/8Cf/8Cv/8Df/8B9/8Rf/8L9/8Cf+//Bn//Bff/B/v99/8Zf/8J3/8C/v8G+/8Bf/8B+/8G9f8G/v8C9/8Gfv8C3/8Cn/8J7/8E7/8G/f8Bv/8Bf/8F/f8C+/8X7/8Bv/8B9/8F+/8Hv/8N3/8B/f8C/v7/Au//Bd//Ab//Avvn99//Af3/B/v/B337/wbv/wT7/wXv/xB//xD+/wjO+/8E/v8C9/8E99//A9//Bfv/FL//A/3/DOb/Hnv/DOf/Ed//BP37/wPv/wb8/xX9/wh//we//wHv/wbv/wv1/wjf/wP9/wN//wS/f/8B6/8Hv/8Cf/8Dv/8L9/8Bf/8F7/8B9/8Gf/8Q7/8C+f8G/f8P/f8f9/8O9/8I3/8F/f8B+/v/A/d//wH9/wK//vf3/wV//wL+/wL3/wHf/wn7/xH+9/8G/f8E9/8F7/8U3/8Cv/8Bn7//Avv/Cvf/Bvf/Bvd//wPf/wX9/wrv/wN//wHv/v8N9f8C/f3/A/7/Bd//E/t//v8J/v8J3/8F3/8S+/8C63//Au//Ad//A/v/B7//BvP/Dfv9/wL7/xHf9/8H7/8Q7/8H7/8N/f8GX/8H3/v/Au/6/w/v/wV/v/8R7/8L+/8B/f8Fv/8Mv/8Bv/8Dv/8J8/8B9/8L/f8D3/8C7/8S/v8I/P8N3/8C9/8B/f8G+/8F7vv/Bff/Av77/wL5/wl//wS//wt//wTv/wJ//wP3/wd//wXv/wLv/wf+/v8G/f8C7/8T/v8C6/8T7f7/A/f/Ab//Ae//AX//Bv3/Dd//Br9/f/8Dv/8Jv/8B9/8H/f8G/f8F/v8Hf/f/Dd//A9//Cn//A+//Cv7/BPf/A/v/G/3/Av7/Bvv/BL//Bb//CH//DO//Gf7/Cd//Af7/BX//Ae//Av3/AX/+/wHv/wj9/xDv/wrf/wP7/wLv/wb3/wr+/w7ff/8C7/8Bf/8I/t//B7//BN//Br//CP5//wN//w/3/wv+/wP3/xLv/wj37/8E9/8F9v8I3/8l5/8C3/8H/v8Ke/8Dn/8Gf/8Bv/8I3/8Bv/8B3/8J+/7/BPf/CP3/CO//Avv/Ae//Bd//Bvf/B9//Ab//A/v/Bff/Dt//A+//BO7/At//Ar//Au/f/wz9/wf+f/8P7/8B7v8Mv/8Uv/8I+/8I/v8Efv8F7/8B+/8B/v8Bv/8H7/8C3/8I/f8E3/8Cf/8Bvz/+/wH9/wH7/wnv/wL+/wL5/wX7/wL9/wG//wr+/wLf/wX7/v8D9/8B+v8Gf3//CPf/AX//CP7/Af3/HH//BPv/Atf/A/v/BO/7/w/+/wf3/v8Jv/8H/f8Cvf8Hf/8Kf/8H/f8G/f8Ev/8I/f8G+/8C7/8Hv/8G/v8O/f8L+/3/Bfe7/wS//wK//wm/f/8B/f8B9/8B9/8G+/8D3/v/A/7+/wj73/8Jf/X/Au//Av3/CPf/EN//Ab//Cf3/Afv/CO//B7//Bd//Cd//Af3/D/7/Bv7/Ab//CL//Cvuv/wHf/wTX/wfv/wb7/xD3/wT++9//Du//Av7/Dv3/Evv3/wXu3/8C+/8B/v8B9/8H3/f/Af7/CH//Ef3/Avv/BPf/A7/m/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v8O/f8R3/8Ff/8F3/8B7/8E/f8Lv/8F/v8J+/8I9/8E/f8F5/8E7/8L3/8E3f8H9/8N3/8Cv/8C9/8Kvv8N/f8Bf/8Q+/8D+/8L7/8D7/3/Avv/Ce//Bd//Bff/AX/3/wLf/wHx/wn+/wPv/wJ//f7/Bff/Cvv/Avf/Bt/v/x1//wF//wr+/wb3/wW//wTv/w7f/wX+/wP9/xDf/wE//wPv/wb7/wLf/wb3/xTf/wG//wW+/xT+/wL79/8Jv/8G/v8Ef/8E39//Av3/Aj++f/gdmAiP/wT7/wI=",
		"zZtABgAAAAAEDkEGAAAAAAEL+/8B9/8G9/8Gvf8Cv/8B/v8B+/8G+/8E7/8B/v8E98//Ct//C7//A/v/AX//A7v/Cvf/Aff/At8//w/+/wW//wH+/wTf/wP7/wnv/wP3/wl/f+//Bn//BPvf/w78/v8J7/8Xf/8C+/8H/f8C+/8Ev/8E3/8Df+//AX//Ar//A3//Aff/Af3/G/f/Af1/7/8Gf/8F9/8H+/33/xl//wX7/wPf/wL+/wb7/wF/3/8H9f8G/v8C9/8G/v8C3/8Cn/8J7/8Cv/8B7/8G/f8Bvv8Bf/8F/f8C+/8X7/8Bv/8B9/8F+/8Hv/8N3/8B/f8C/v7/Au//Bd//Ab//Avvn9/8C/f8Pffv/Bu//BPv/Be//Cv7/BX//EP7/CE77/wT+/wL3/wT33/8D3/8F+/8Uv/8D/f8Lv+b/Hnv/DOf/Ed//BP37/wr8/wHv/xP9/wh//we//wHv/xL1/wz9/wN//wS//v8B6/8Kf/f/Ar//C/f/AX//Be//Aff/Afv/BH//E/n/BP7/EPf9/y73/w79/wH7+/8D93//Af3/Ar/+9/f/BX//Av7/Avf/Ad//Cfv/Ef73/wv3/xrf/wK//wGfv/8B/vv/Cvf/Bvf/Bvd//wPf/wH7/w7P/wN//wHv/v8H/v8F9f8C/f3/A/7/Bd//E/t//wr+/wnf/wXf/wn+/wj7/wLrf/8C7/8B3/8D+/8Hvv8G8/8N+/3/Avv/Ed/3/wfv/xDv/wfv/w39/wZf/wff+/8C7/r/CO//Bu//BN9/v/8I/v8I7/8Ev/8G+/8B/f8Fv/8Mv/8Bv/8F3/8H8/8B9/8B/f8J/f8D3/8C7/8G3/8L/v8I/P8N3/8C9/8B/f8E7/8B+/8F/vv/Bff/Av7/Ae//Afn/De+//wr3f/8E7/8Cf/8C9/f/B3//Be//Au//B/7+/wb9/wLv/xP+/wLv/xPt/v8D9/8Bv/8B7/8Bf/8G/f8N3/8Gv39//wO//wv3/w79/w1/9/8N3/8D3/8F7/8I7/8P9/8D+/8Pv/8L/fv/Af7/C7//Bb//CH//DO//Av3/Fv7/Bvv/At//Af7/BX//Ae//AX/9/wF//v8b7/8K3/8G7/8G9/8K/v8O33//Au//AX//CP7/CL+//wq//wj+f/8Df/8Kv/8E9/8H+/8D/v8D9/8Gf/8L7/8Bv/8G9+//BP3/Bfb/GL//Fef/Cv7/Cnv/A5//CL//CN//A9//Cfv/Bff/CP3/CO//Av3/Ae//Bd//Bvf/Cb//A/v/Bff/CH//Bd//A+//BP7/Af3f/wK//wLv3/8L7/8I/n//A/7/C+//Ae7f/wu//xz9+/8I/v8Ef/7/BO//Afv/Af7/Ab//Ct//B/v9/wTf/wJ//wG/P/7/Af3/Afv/Ce//Avr/Avn/Bfv/Av3/Ab/9/wv33/8F+/7/A/P/Afr/Bn9//wj3/wF//wl//f8N9/8Of/8E+/8C1/8D+7//A+/7/xf3/v8Jv/8H/f8Cvf8B3/8Q+/8C/f8E/f8F/f3/BL//CP3/Au//A/v/Au//B7//Bv7/Dv3/C/v9/wX3+/8Ev/8Cv/8Jv3//Af3/Aff/CPv/BPv/A/7+/wj73/8K9f8C7/8C/f8I9/8Q3/8Bv/8L+/8I7/8Hv/8F3/8L/f8P/v8G/v8Bv/8B9/8N/f8D+6//Ad//BNf/B+//Btv/EPf/Bfvf/xH+/wf3/wb9/xL79/8F7t//Avv/A/f/B9/3/wH+/wh//xH9/wL7/wT3/wK//wHm/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v8Mv/8T3/8Ff/8F3/8B7/8Qv/8P+/8I9/8E/f8F5/8E7/8G/f8J3f8H9/v/DN//A+//Aff/Cr7/Av3/Cv3/AX//CP7/B/v/A/v/C+//A+/9/wL7/wnv/wXf/wH9/wV/9/8C3/8B8f8G9/8C/v8D7/8Cf/3+/wX3/w33/wbf/x5//wF//wr+/wb3/wW9/wTv/wvf/wLf/wX+/wP9/xDf/wE//wPv/wb7/wLf/wb3/wvf/wq//wW+/xT+/wL79/8Q/v8Ef/8E39//Av3/Aj++f/gfmAiP/wT7/wI=",
		"/Z1ABgAAAAAGDkEGAAAAAAFAygABAj//Avv/Bbn/Cvf/Aff/At8//w/++/8Ev/8B/v8E3/8D+/8J7/8D9/3/CH9//wd//wT73/8O/P7/A/3/Be//Ff7/AX//Avv/B/3v/wH7/wS//wTf/wN//wJ//wK//wN/v/f/Gf3/A/f/An/v/wZ//wX3/wf7/ff/DN//DH//Cd//Av7/Bvv/AX//CPX/Bv7/Avf/Bv7/At//Ap//CP3v/wTv/wb9/wG//wF//wX9/wL7/wS//w2//wTv/wG//wH3/wHv/wP7/we//w3f/wH9/wL+/v8B/u//Bd//Ab//Af775/f/Av3/Bu//CH37/wbv/wT7/wXv/xB//xD+/wL3/wXO+/8E/v8C9/8E99//A9//Bfv/Bvf/Db//A/3/DOb/Hnv/DOf/Ed//BP37/wr8/xX9/wh//we//wHv/xL1/wz9/wN//wS//wLr/wp//wO//wT3/wb3/wF//wXv/wH3/wZ//w79/wT5/xLf/wP9/y73/w3v/f8B+/v/A/d//wH9/wK//vf3/wV//wL+/wL3/wHf/wn7/wj+/wj+9/8Bv/8J9/3/Gd//Ar//AZ+//wL7/wr3/wb3/wb3f/8D3/8D7/8M7/8Df/8B7/7/DfX/Av39/wP+/wXf/wn9/wn7f/8K/v8E9/8E3/8F3/8S+/8C63//Au//Ad//A/v/B7//BvP/Dfv9/wL7/wTv/wzf9/8H7/8Q7/8H7/8N/f8GX/8H3/v/Au/6/w/v/wV/v/8K/f8G7/8L+/8B/f8Fv/8G9/8Fv/8Bv/8N8/8B9/8L/f8D3/8C7/8S/v8I/P8Cf/8K3/8C9+/9/wb7/wX++/8F9/8C/v8D+f8Ov/8D7/8Hf/8E7/8Cf/8D9/8Hf/8F7/8C7/8B3/8F/v7/Bv3/Au//E/7/Au//E+3+/wP3/wG//wHv/wF//wb9/w3f/wa/f3//A7//CP3/Avf/Dv3/DX/3/w3f/wPf/w7v/w/3v/8C+/8b/f8C/v8Lv/8Fv/8If/8M7/8Z/v8J3/8B/v7/BH//Ae//Av3/AX/+/xvv/wL3/wff/wbv/wb3/wr+/w7ff/8C7/8Bf/8B/f8G/v8Iv/8Lv3//B/5//wLvf/8P9/8L/v8D9/8K+/8H7/8I9+//Cvb/Luf/AX//CP7/Cnv/A5//CL//CN//A9//B/3/Afv/Bff/Bvv/Af3/CO//BO/+/wTf/wb3/wm//wP7/wX3/w7f/wPv/wT+/wLf/wK//wLv3/8Q9/8D/n//D+//Ae7/A9//CL//GX//A/v/CP7/BH//Be//Afv/Af7/Ab//Ct//Bf7/Av3/BN//An//Ab8//v8B/f8B+/8J7/8C/v8C+f8F+/8C/f8Bv/8N3/8F+/r/A/f/Afr/Bn99/wj3/wF//wG//wj9/wN//xh//wT7/wLX/wP7/wTv+/8X9/7/Cb//B/3/Ar3/Ff3/BP3/BP3/Af3/BL//CP3/Bvv/Au//B7//Bv7/Dv3/C/v9/wX3+/8Ev/8Cv/8Jv3//Af3/Aff/CPv/BPv/A/7+/wj73/8K9f8C7/8C/f8I9/8Q3/8Bv/8L+/8I7/8B/f8Fv/8F3/8J7/8B/f8P/v8G/v8Bv/8T+6//Ad//BNf/B+//Bnv/DH//A/f/Bfvf/xH+/w79/xL79/8F7t//Avv/A/f/B9/3/wH+/wh//wb9/wr9/wH9+/8E9/8E5v8K9/8B/v8Bf/8B9/3/A2/19bv/BO//Bv7/Dr//Ed//BX//Bd//Ae//Bn//Cb/7/w39+/8I9/8E/f8F5/8E7/8Q3f8H9/8Nz/8F9/8G/v8Dvv8N/f8Bf/8Q8/8D+/8L7/v/Au/9/wL7/wnv/wXf/wd/t/8C3/8B8f8J/v8D7/8Cf/3+/wX3/wv7/wH3/wbf/x5//wF//wr+/wb3/wW//wTv/w7f/wX+/wP9/xDf/wE//wPv/wP3/wL7/wLf/wb3/wZ//wrf/wS//wW+7/8T/v8C+/f/EP7/BH//BN/f/wL9/wJ/vn/4H5oMq/8E+/8C",
		"x8BABgAAAAAJDkEGAAAAAAH+/wTf/wTf/wXf7/8R+/8C63//Au//Ad//A/v/At//BL//BvP/Dfv9/wL7/xHf9/8H7/8L/v8E7/8H7/8N/f8GX9//Bt/7/wLv+v8Pz/8Ff7//Ee//C/v/Af3/Bb//DL//Ab//DfP/Aff/C/3/A9//Au//Ev7/CPz/Dd//Avf/Af3/Bvv/Bf77/wX3/wL+/wP5/w6//wt//wH+/wLv/wJ//wP3/wd//wXv/wLv/wf+/v8G/f8C7/8T/v8C7/8T7f7/A/f/Ab//Ae//AX//Bv3/Dd//Br9/f/8Dv/8L9/8O/f8I+/8Ef/f/Dd//A9//Du//D/f/A/v/G/3/Av7/C7//Bb//Bfv/An//Cb//Au//B/f/Ef7/Cd//Af7/BX//Ae//Av3ff/7/CPv/Cvv/B+//Ct+//wXv/wb3/wr+/w7ff/8C7/8Bf/8I/v8Iv/8Lv/8I/n9//wJ//wy//wL3/wv+/wP3/xLv/wj37/8K9v8u5/8K/v8Ke/8Dn/8Iv/8I3/8D3/8J+/8F9/8Bf/8G/f8E/f8D7/8E7/8F2/8G9/8Jv/8D+/8F9/8O3/8D7/8E/v8C3/8Cv/8C79//BX//Dvp//w/v/wHu/wn9/wK//wt//wr3/wb7/wPf/wT+/wR//wXv/wH7/wH+/wG//wrf/wj9/wTf/wJ//wG/P/7/Af3/Afv/Ce//Av7/Avn/Bfv/Av3/Ab//BPv/CN//Bfv+/wP3/wH6/wZ/f/8I939//wr9/wG//xP7/wZ//wT7f/8B1/8D+/8D/u/7/wvf/wv3/v8Jv/8H/f8Cvf8B+/8M9/8G/f8E/f8G/f8Ev/8I/f8G+/8C7/8Hv/8G/v8O/f8H9/8D+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wP+/wT3/xDf/wG//wv7/wS//wPv/w3f/wv9/w/+/wb+/wG//xP7r/8B3/8E1/8H7/8G+/8K3/8F9/8F+9//Ef7/Dv3/Db//BPv3/wXu3/8C+/8D9/8C/v8E3/f/Af7/A/3/BH//Ef33/wH7/wT3/wTm/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v8g3/8Ff9//BN//Ae//EL//CPv/Bvv/Bd//Avf/BP3/Bef/BO//EN3/B/f/Dd//BXf/Cr7/Df3/AXf/Cu//Bfv/A/v/C+//Af3/Ae/9/wL7/wd//wHv/wXf/wd/9/8C3/8B8f8Cv/8G/v8D7/8CX/3+/wX3/w33/wbf/x5//wF//wr+/wXf9/8Fv/8E7/8O3/8F/v8D/f8Q3/8BP/8D6/8G+/8C3/8G9/8J7/8Mv/8Fvv8U/v8C+/f/EP7/BH//BN/f/wL9/wJ/vn/4H5oMqn//A/v/Ag==",
		"zZtABgAAAAARDkEGAAAAAAEL+/8B9/8G9/8G/f8Cv/8B/v8B+/8L7/8B/v8E94//Ct//Av7/CL//A/v/Bbv/Cvf/Aff/At8//wf9/wf+/wW//wH+/wTf/wP7/wnv/wP3/wl/f7//Bn//BHvf/w78/v8J7/8Xf/8C+/8H/f8C+/8C+/8Bv/8E3/8Df/8Cf/8Cv/8Df/8B9/8d9/8Cf+//Bn//Bff/B/v99/8P+/8Jf/8B7/8H3/8C/v8G+/8Bf/8E3/8D9f8G/v8C9/8G/v8C3/8Cn/8J7/8E7/8G/f8Bv/8Bf/8C/v8C/f8C+/8W3+//Ab//Aff/Bbv/B7//Dd//Af3/Av7+/wLv/wXf/wG//wLr5/f/Av3/D337/wbv/wF//wL7/wXv/xB//xD+/wjO+/8E/v8C9/8E99//A9//Bfv/FL//A/3/DOb/Hnv/DOf/Ed//BP37/wr8/xX9/wh//wbfv/8B7/8S9f8M/f8Df/8Ev/8C6/8Kf/8B+/8Bv/8Jf/8B9/8Bf/8F7/8B9/8Gf/8C/f8Q+f8W/f8H7/8Z/f8M9/8B/f8M/f8B+/v/A/d//wH9/wK//vf3/wV3/wL+/wL3/wHf/wn7/wv9/wX+9/8L9/8E7/8V3/8Cv/8Bn7//Avv/Cvf/BPv/Aff/Bvd//wPf/wL+/w3v/wHv/wF//wHv/v8N9f8C/f3/A/7/Bd//De//Bft//wr+/wnf/wL3/wLf/xL7/wLrf/8C7/8B3/8D+/8Hv/8G8/8N+/3/Avv/Ed/3/wfv/xDv/wfv/wP7/wn9/wZf/wff+9//Ae/6/w/v/wV/v/8R5/8L+/8B/f8Fv/8Mv/8Bv/8N8/8B9/8L/f8D3/8C7/8S/v8I/P8N3/8C9/8B/f8G+/8F/vv/Bff/Av7/A/n/Dr//C3//BO//An//A/f/B3//Be//Au//B/7+/wb9/wLv/xP+/wLv/xPt/v8D9/8Bv/8B7/8Bf/8G/f8E/f8I3/8Gv39//wO//wv3/w79/w1/9/8N3/8D3/8Df/8K7/8P9/8D+/8D9/8X/f8C/v8G/f8Ev/8Fv/8D7/8Ef/8M7/8B7/8X/v8J3/8B/v8Ff/8B7/8C/f8Bf/7/G+//Ct//Bu//Bvf/Cv7/Dt9//wLv/wF//wj+/wi//wu//wj+f/8Df/8D+/8L5/8L/v8D9/8C+/8K3/8E7/8I9+//Cvb/Bv7/HO//Cuf/Av3/B/7/A3//Bnv/A5//CL//A/v/BN//A9//Bd//A/v/Bff/CP3/CO//BO//Bd//Bvf/Cb//A/v/Bff/Dt//A+//BP7/At//Ar//Au/f/xT+f/8P7/8B7v8L97//Hfv/CP7/BH//Be//Afv/Af7/Ab//Ct//CP3/BN//An//Ab8//v8B/f8B+/8J7/8C/v8C+f8F+/8C/f8Bv/8N3/8F+/7/A/f/Afr/Bn9//wj3/wF//wr9/w6//wt//wF//wT7/wLX/wP7/wTv+/8P9/8H9/7/Cb//Bv39/wK9/xX9/wT9/wb9/wS//wj9/wF//wT7/wLv/we//wX9/v8O/f8L+/3/Bff7/wS//wK/+/8Iv3//Af3/Aff/CPv/BPv/A/7+/wj73/8K9f8C7/8C/f8Gv/8B9/8Q3/8Bv/8Hv/8D+/8I7/8E/v8Cv/8F3/8D/v8H/f8D/f8L/v7/Bf7/Ab//E/uv/wHf/wTX/wfv/wX9+/8W+9//Ef7/Dv3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/B/d//xH9/wL7/wT3/wTm/wr3/wH+/n//Aff9/wNv9f27/wTv/wb+/w2//xC//wHf/wV//wXf/wHv/xC//w/7/wj3/wT9/wXn/wTv/xDd/wf3/w3f/wX3/wK//we+/w39/wF//xD7/wG//wH7/wl//wHv/wPv/f8C+/8J7/8F3/8Hf/f/At//AfH/Cf7/A+//An/9/u//BPf/Dff/Bt//Hn//AX//Cv7/Bvf/Bb//BO//Dt//Bf7/A/3/EN9/P/8D7/8G+/8C3/8Bv/8E9/8I3/8Nv/8Fvv8U/v8C+/f/EP7/BH//BN/f+/8B/f8Cf75/+B+aTqpAf/8C+/8C",
		"r9lABgAAAAATDkEGAAAAAAH/Bf7/Dt9//wLv/wF//wj+/wif/wu//wPf/wT+f/8Df/8P9/8L/v8D9/8S7/8I9+//Cvb/Luf/Cv7/Cnv/A5//CL//CN//A9//Cfv/Bff/CP3/CO//BO//Bd//Bvf/Cb//A3v/Bff/Dt//A+//BP7/At//Ar//Au/f/xT+f/8Cv/8M7/8B7v8Mv/8d+/8I/v8Ef/8D/f8B7/8B+/8B/v8Bv/8K3/8I/f8E37//AX//Ab8//v8B/f8B+/8F+/8D7/8C/v8C+f8F+/8C/f8Bv/8N3/8F+/7/A/f/Afr/Bn9//wj3/wF//wr9/xh//wN//wT7/wLX/wP7/wTv+/8T+/8D9/7/Cb//B/3/Ar3/Ff3/BP3/BPv/Af3/BL//CP3/Bvv/Au//B7//Bv7/Dv3/C/v9/wX3+/8Ev/8Cv/8Jv3//Af3/Aff/CPv/BPv/A/b+/wj73/8K9f8C7/8C/f8I9/8Q3/8Bv/8B/f8J+/8Gf/8B7/8Hv/8D7/8B3/8L/f8P/v8G/v8Bv/8T+6//Ad//BNf/B+//Bvv/EPf/Bfvf/xH+/w79/xL79/8F7t//Avv/A/f/B9/3/wH+/wh//wn3/wf9/wL7/wT3/wTm/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v8P/v8Q3/8Ff/8F3/8B7/8K7/8Fv/8P+/8Fv/8C9/8E/f8F5+//A+//EN3/B/f/A7//Cd//Bff/Cr7/Df3/AX//EPv/A/v/C+//A+/9/wL7/wnv/wXf/wd/9/8C3/8B8f8J/v8D79//AX/9/v8F9/8E+/8I9/8G3/8ef/8Bf/8K/v8G9/8Fv/8E7/8O3/8F/v8D/f8Q3/8BP/8B9/8B7/8G+/8C3/8G9/8Wv/8Fvv8U/v8C+/f/Cfv/Bv7/BH//BN/f/wL9/wJ//n/6H5peqkBf/wL7/wI=",
		"+nlABgAAAAATDkEGAAAAAAH/AfqyDQABSgAB3/8I3/8Ff/8F/v8Fr/8F7/8F/f8H7/8C/fv/CO//E9//Cn/+/wTf/wf3/wHv/wXf/wz9/f8Lf3//Fvf/At//A/7/A+3/Au//Ar/+7/8B7/8R+f8Vf/8Dv3//CPf7/wHf/x733/f/Aef/Avf/Bd/v/v8B+/8Q3/8D/v8n7/8K/v8D9/8F3/8Fv/8o/u//Af7/BO/3/f8Cv/8B+3//It/f/wO//wvvf/8L+/8I9/8Df/v/Af3/A/v/Arv/BP7/Du//Dn//Ar//C79//wLz/xj+/wa/v9//A/3/Aff/Cb//Br7/Efvv/w/vv/8D/f8L+/8I9/8H/v8D3/8G/v8D9/8G/v8B3/8B9/8N3/3/Bf7/Bn//B/f/D7/9/wf9/wXv/xr7/wb7/wfv/w33/f8Cf/8Lf/8C3/8Gd/8Cf/8L+/8P+/8Lf3//B7//B/7/Ab//Gf3f/wG/f/8D+8//A/3/A/7/En//Au+//wH3/wPb/wT9/wb3/wI2/wF//wbf/wj9/wK//wHf/wvv/xPv+/8B9/8G9/8G/f8Cv/8B/vv7/wvv/wH+/wT3z/8K3/8Lv/8D+/8Fu/8K9/8B9/8C3z//D/7/Bbf/Af7/BN//A/v/Ce//A/f/CX9//wd//wT73/7/Dfz+/f8I7/8Fv/8Rf/8C+/8H/f8C+/8Ev/8E39//An//An//Ar//A3//Aff/DP3/EPf/An/v/wz3/wf7/ff/GX//A7//Bd//Av7/Bvv/AX//CPX/Bv7/Avf/BPf/Af7/At//Ap//Ce//BO//Bv3/Ab//AX//Bf3/Avv/F+//Ab//Aff/Bfv/B7//Dd//Ad3/Av7+/wLv/wXf/wG//wL75/f/Av3/D337/wbv/wT7/wXv/xB//xD+/wjO+/8E/v8C9/8E99//A9//Bfv/FL/v/wL9/wzm/xz+/wF7/wzn/xb9+/8K/P8F7/8P/f8If/8Hv/8B7/8S9f8M/f8Df/8Ev/8C6/8Kf/8Dv/8L9/8Bf/8F7/8B9/8F93//E/n/Fv3/Dff/GO//B/f/Dv3/Afv7/wP3f/8B/f8Cv/739/8Ff/8C/v8C9/8B3/8J+/8J/v8H/vf/C/f/Gt//Ar//AZ+//wL7/wr3/wb3/wb3f/8D3/8Q7/8B3/8Bf/8B7/7/DfX/Av39/wP+/wXf/xP7f/3/Cf7/B/v/Ad//Bd//Evv/Aut//wLv/wHf/wP79/8Gv/8G8/8B+/8L+/3/Avv/Ed/3/wfv/xDv/wfv/w39/wZf/wff+/8C7/r/D+//BX+//xHv/wV//wX7/wH9/wW//wy//wG//wu//wHz/wH3/wv9/wPf/wLv/xL+/wj8/w3f/wL3/wH9/wb7/wX++/8F9/8C/v8D+f8Ov/8F+/8Ff/8E7/8Cf/8D9/8Hf/8F7/8C7/8H/v7/Bv3/Au//Cr//CP7/Au//Be//De3+/wP3/wG//wHv/wF//wb9/wz+3/8Gv39//wL+v/8L9/8F7/8I/f8Nf/f/Dd//Aff/Ad//DP3/Ae//D/f/Ar/7/xHf/wn9/wL+/wu//wW//wh//wzv/xn+/wnf/wH+/wV//wHv/wL9/wF//v8b7/8F9/8E3/8G7/8G9/8K/v8G3/8H33//Au//AX//CP7/CL9//wq//wj+f9//An//D/f/C/7/A/f/Eu//CPfv/wr2/yp//wPn/wl//v8Ke/8Dn/8Iv/8I3/8D3/8I9/v/Bff/CP3/CO//BO//Bd//Bvf/Cb//A/v/Bff/Dt//A+//BP7/At//Ar//Au/f/xT+f/8P7/8B7v8B3/8Kv/8d+/8I/v8Ef/8F7/8B+/8B/v8Bv/8K3/8D+/8E/f8E3/8Cf/8Bvz/+/wH9/wH7/wnv/wL+/wL5/wX7/wL9/wG//w3f/wX7/v8D9/8B+v8Gf3//Bf3/Avf/AX//Cv3/E7//CH//BPv/Atf/A/v/BO/7/xf3/v8D9/8Fv/8H/f8CvP8V/f8E/f8G/f8Ev/8I/f8G+/8C7/8Hv/8G/v8O/f8L+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wX7/wL3/xDf/wG//wv7/wjv/w3f/wv9/w/+/wb+/wG//xH7/wH7r/8B3/8E1/8H7/8G+/8I3/8N+9//Efb/Dn3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/CH//EfX/Avv/BPf/BOb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wP7/wL+/yDf/f8Ef/8F3/8B7/8F/v8Kv/8P+/8Cf/8F9/8E/f8F5/8E7/8E7/8L3f8H9/8N3/8F9/8Kvv8C+/8K/f8Bf/8Df/8M+9//Avv/C+//A+/9/wL7/wnv/wXf/wdv9/8C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Aff/CPv/BPv/A/7+/wj73/8K9f8C7/8C+f8I9/8Q3/8Bv/8L+/8I7/8Hv/8F3/8L/f8P/v8G/v8Bv/8J9/8J+6//Ad//BNf/B+//Bvv/Cd//Bvf/Bfvf/xH+/w79/xL79/8F7t//Avv/A/f/A/3/A9/3/wH+/wh//xH9/wL7/wT3/wTm/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v3/H9//A/f/AX//Bd//Ae//BH//C7//D/v/COf/BP3/Bef/BO//EN3/B/f/DZ//Bff/Bf3/BL7/B/f/Bf3/AX//EPv/A/v/C+//A+/9/wL7/wnv/wXf/wd/9/8C3/8B8f8J/v8D7/8Cf/3+v/8E9/8N9/8G3/8ef/8Bf/8K/v8G9/8Fv/8E7/8O3/8F/v8D/f8B+/8L3/8C3/8BP/8D7/8G+/8C3/8G9/8Wv/8Fvv8I9/8L/v8C+/f/FX//BN/f/wV//wF/+h+eXqpgWH//Afv/Ag==",
		"zZtABgAAAAAZDkEGAAAAAAEL+/8B9/8G9/8Gvf8Cv/8B/v8B+/8G+/8E7/8B/v8E98//Ct//C7//A/v/AX//A7v/Cvf/Aff/At8//w/+/wW//wH+/wTf/wP7/wnv/wP3/wl/f+//Bn//BPvf/w78/v8J7/8Xf/8C+/8H/f8C+/8Ev/8E3/8Df+//AX//Ar//A3//Aff/Af3/G/f/Af1/7/8Gf/8F9/8H+/33/xl//wX7/wPf/wL+/wb7/wF/3/8H9f8G/v8C9/8G/v8C3/8Cn/8J7/8Cv/8B7/8G/f8Bvv8Bf/8F/f8C+/8X7/8Bv/8B9/8F+/8Hv/8N3/8B/f8C/v7/Au//Bd//Ab//Avvn9/8C/f8Pffv/Bu//BPv/Be//Cv7/BX//EP7/CE77/wT+/wL3/wT33/8D3/8F+/8Uv/8D/f8Lv+b/Hnv/DOf/Ed//BP37/wr8/wHv/xP9/wh//we//wHv/xL1/wz9/wN//wS//v8B6/8Kf/f/Ar//C/f/AX//Be//Aff/Afv/BH//E/n/BP7/EPf9/y73/w79/wH7+/8D93//Af3/Ar/+9/f/BX//Av7/Avf/Ad//Cfv/Ef73/wv3/xrf/wK//wGfv/8B/vv/Cvf/Bvf/Bvd//wPf/wH7/w7P/wN//wHv/v8H/v8F9f8C/f3/A/7/Bd//E/t//wr+/wnf/wXf/wn+/wj7/wLrf/8C7/8B3/8D+/8Hvv8G8/8N+/3/Avv/Ed/3/wfv/xDv/wfv/w39/wZf/wff+/8C7/r/CO//Bu//BN9/v/8I/v8I7/8Ev/8G+/8B/f8Fv/8Mv/8Bv/8F3/8H8/8B9/8B/f8J/f8D3/8C7/8G3/8L/v8I/P8N3/8C9/8B/f8E7/8B+/8F/vv/Bff/Av7/Ae//Afn/De+//wr3f/8E7/8Cf/8C9/f/B3//Be//Au//B/7+/wb9/wLv/xP+/wLv/xPt/v8D9/8Bv/8B7/8Bf/8G/f8N3/8Gv39//wO//wv3/w79/w1/9/8N3/8D3/8F7/8I7/8P9/8D+/8Pv/8L/fv/Af7/C7//Bb//CH//DO//Av3/Fv7/Bvv/At//Af7/BX//Ae//AX/9/wF//v8b7/8K3/8G7/8G9/8K/v8O33//Au//AX//CP7/CL+//wq//wj+f/8Df/8Kv/8E9/8H+/8D/v8D9/8Gf/8L7/8Bv/8G9+//BP3/Bfb/GL//Fef/Cv7/Cnv/A5//CL//CN//A9//Cfv/Bff/CP3/CO//Av3/Ae//Bd//Bvf/Cb//A/v/Bff/CH//Bd//A+//BP7/Af3f/wK//wLv3/8L7/8I/n//A/7/C+//Ae7f/wu//xz9+/8I/v8Ef/7/BO//Afv/Af7/Ab//Ct//B/v9/wTf/wJ//wG/P/7/Af3/Afv/Ce//Avr/Avn/Bfv/Av3/Ab/9/wv33/8F+/7/A/P/Afr/Bn9//wj3/wF//wl//f8N9/8Of/8E+/8C1/8D+7//A+/7/xf3/v8Jv/8H/f8Cvf8B3/8Q+/8C/f8E/f8F/f3/BL//CP3/Au//A/v/Au//B7//Bv7/Dv3/C/v9/wX3+/8Ev/8Cv/8Jv3//Af3/Aff/CPv/BPv/A/7+/wj73/8K9f8C7/8C/f8I9/8Q3/8Bv/8L+/8I7/8Hv/8F3/8L/f8P/v8G/v8Bv/8B9/8N/f8D+6//Ad//BNf/B+//Btv/EPf/Bfvf/xH+/wf3/wb9/xL79/8F7t//Avv/A/f/B9/3/wH+/wh//xH9/wL7/wT3/wK//wHm/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v8Mv/8T3/8Ff/8F3/8B7/8Qv/8P+/8I9/8E/f8F5/8E7/8G/f8J3f8H9/v/DN//A+//Aff/Cr7/Av3/Cv3/AX//CP7/B/v/A/v/C+//A+/9/wL7/wnv/wXf/wH9/wV/9/8C3/8B8f8G9/8C/v8D7/8Cf/3+/wX3/w33/wbf/x5//wF//wr+/wb3/wW9/wTv/wvf/wLf/wX+/wP9/xDf/wE//wPv/wb7/wLf/wb3/wvf/wq//wW+/xT+/wL79/8a39//BX//AX/6H55eqmBYf/8B+/8C",
		"+T1ABgAAAAAZDkEGAAAAAAE6u5Yc0gEPv/8I/f8Bf/8Fv/8E/f8B9/8F+/8H/n//B3//C3//Eb/9/wT93/8Hv/8C9/8Cv/ff/wP9/wPv/w7+/wp//wP9/wT3/wSv/wL7/wS7/wR//wL9/wXv/wL3/wbf/xLv/wX39/8K9/8Bv/8E9/8E/f8H9/8C9f8G/f8C/P8B9/8Q+9//CPvf/wG//wr39/8E+/8G+/8Lb/8C+/8C9/8D/f8C3/8E3/8D3/8K9/8Dv/8Bvv8D/v8K/fv/A/3/A+93/xD+/wjP/wH9f3//A/f/B/v/Bu//Ab/v/wT9/wLf/wHt/wPf/wjv/wL9/wT9/wP3/wG//wf+/wH9/wPv/ws//wT+/wPf/wF//wHs93/+/wL9/wL3/wL7/w57/wi//wNz/wZ+/wT9/wXf/wL3/wJ//wff/wHv/wb3/wh//wH3/wH9/wb+/wHv/f8Gb/8Hv/t//wP9/wj+/wV//wT+/wbv/w73/wL3/f8Bv/8Hv/3/Ae//Am+/ff39/wj+/wX+/wT9/wT+/wbP/wjv/wK//wXv3/8K+/7/AX//EPv/Av7/Eb//GPv/B/3/C3//EX9/7/8Gf/8Cv/8E773/A7//Aff/Cr//Efr/Bb//A/7/CP7/Avv/B/3/Av7+/wn3/v8B7/8Fv/v/Ad7+v/8C7/8G3/8I+f8E7/8L7/8B9/8Cv/8J3/8H+/8H/f8Cf/8Be/8Dv/8G/f8C/f8Dv/8F9/8Dx/n/Ad//Au//Bd//Bv3/Dr//AXfv/wH++/8Bf/8F+X//An/u/wO//wP7/wR//wL+u/8Df/8D+/8E9/vff/8C9/8F/u9//wJ/f/8Bv/8Dz/8Dnv8B+/8Jf/8Cv/8E2/8D/f8B3/8Uv/8D3/8Bv/8H+/8G9/8Cv/8b9/8B9/8H+/8Xf/8C1/8R7/8B7/8Bv/3+/wLv/wbf/wF//wH9/wv3/wPX/wL9/wP9/wfv/wKv3/8Bf/8C/f8I1/8C7/8G7/8G/P8R9/8E7/v/BP7/Afv/A3/+/wLf/wL+/wLv/wHf/wfv/xTv/wr7/wu//wH3/wN//wG//wT9/wi3/wW//wv3/n//Avr/C/v/C/z/B/f/B/7/CL//BPf/BZ7f/wO//wH7/wG//wf3/wH+/wH+/wfv/wL7/wXf/wXv/w33/w/v/wj+/wb9/wu//wWv/wW//wL+/wV//wL9/wF//wH9f/8E/f8Ev/8B9/8Df/8C7/8F+/8Jv/8Jv/8Ef/f/B3//A3//G/7/Bf3/Cvf/Dff/Bvf/Av3/At//B7//Au//AW//Bd//Ad//CN//C/7/Ba//Be//Bf3/Cv37/wjv/wX7/xh//v8E3/8H9/8B7/8FX/8M/f3/C39//xb3/wLf/wP+/wPt/wLv/wK//u//Ae//BN//DPn/D3//BX//A79//wH3/wb3/wLf/x733/f/Aef/Avf/Bd/v/v8B+/8Q3/8D/v8S9/8U7/8K/v8D9/8F3/8E97//KP7v/wH+/wTv9/3/Ar//Aft//wL7/xN//wvf3/8Dv/8L73//Bfv/Bfv/CPf/A3/7/wH9/wP7/f8Bv/8E/v8O7/8Ff/8If/8Cv/8Lv3//AvP/GPb/Br+/3/8D/f8B8/8Jv/8Gvv8R++//D++//wP97/8K+/8I9/8H/v8C/f8H/v8D9/8G/v8B3/8B9/8N3/3/Bf7/Dvf/Cf3/Bb/9/w3v/xr7/wb7/wfv/w33/f8Cf/8Lf/8C3/8Df/8Cd/8Cf/8nf3//D/7/Ab//EPv/CP3f/wG/f/8D+8//A/3/Ge+//wH3/wPb/wT9/wb3/wI2/wF//wbf/wj93/8Bv/8B3/8L7/8E7/8O7/v/Aff/Bvf/Bv3/Ar//Af7/Afv/Ar//CO//Af7/BPfP/wrf/wu//wP7/wW7/wr3/wH3/wLfP/8L/v8D/v8Fv/8B+v8E3/8D+/8J7/8D9/8Jf3//B3//BPvf/w78/v8J7/8Xf/8C+/8H/f8C+/8D37//BN//A3//An//Ar//A3//Aff/Hff/An/v/wZ//wX3/wf7/ff/En//Bn//BN//BN//Av7/Bvv/AX//CPX+/wX+/wL3/wb+/wLf/wKf7/8I7/8E7/8G/f8Bv/8Bf/8B/v8D/f8C+/8Gv/8Q7/8Bv/8B9/8F+/8Hv/8N3/8B/f8B/f7+/wLv/wXe/wG//wL75/f/Av3/D337/wbv/wT7/wXv/xB//w7+/wH+/wjO+/8E/v8C9/8E99//Au/f/wX7/xS//wP9/wzm/x57/wzn/xHf/wT9+/8K/P8V/f8If/8Hv/8B7/8S9f8M/f8Df/8Ev/8C6/8Kf/8Dv/8D/f8H9/8Bf/8B+/8D7/8B9/8Gf/8T6f8Q/f8F/f8u9/8O/f8B+/v/A/d//wH9/wK//vf3/wV//wL+/wL3/wGf/wn7/xH+9/8E+/8G9/8Gf/8T3/8Cv/8Bn7//Avv/Cvf/A3//Avf/Bvd//wPf/xDv/wN//wHv/v8B/f8L9f8C/f3/A/7/BP3f/xP7f/8H/v8C/v8J3/8C3/8C3/8Mv/8F+/8C63//Au//Ad//A/n/B7//BvP/Dfv9/wL7/xHf9/8H7/8Q7/8H7/8E7/8I/f8GX/8H3/v/Au/6/w/v/wV/v/8R7/8L+/8B9f8Fv/8Mvv8Bv/8N8/8B9/8L/f8D3/8C7/8S/v8I/P8F3/8H3/8C9/8B/f8G+/8F/vv/Bff/Av7/A/n/Dr//C3//BO//An//A/f/B3//Be//Au//B/7+/wb9/wLv/xP+/wLv/xPt/v8D9/8Bv/8B7/8Bf/8G/f8H3/8F3/8Gv39//wO//wG//wn3/w79/w1/9/8N3/8D3/8O7/8P9/8D+/8b/f8C/v8F/v8Fv/8Fv/8If/8M7/8Z/v8Dv/8F3/8B/v8Ff/8B7f8C/f8Bf/7/Br//FO//Ct//Bu//Bvf/Cv7/Dt9//wLv/wF//wj+/wX9/wK//wu//wj+f/8Df/8P9/8L/v8D9/8Mv/8F7/8I9+//CH//Afb/Cf3/If3/Auf/Cf3+/wp7/wOf/wO//wS//wjf/wPf/wn7/f8E9/8I/f8I7/8E7/8F3/8G9/8Jv/8D+/8F9/8O3/8D7/8E/v8C3/8Cv/8C79//FP5//w/v/wHu/wy//wS//xj7/wj+/wR//wXv/wH7/wH6/wG//wrP/wj9/wO/3/8Cf/8Bvz/+/wH9/wHz/wnv/wL+/wL5/wX7/wL9/wG//wnv/wPf/wXr/v8D9/8B+v8Gf3//CPf/AX//Cv3/CP3/E3//BPv/Atf/A/v/BO/7/xf3/v8Jv/8H/P8Cvf8a/f8G/f8Ev/8I/f8G+/8C7/8E9/8Cv/8G/v8O/f8L+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wj3/xDf/wG//wv7/wjv/we//wXf/wv9/w/+/wb+/wG//wL7/xD7r/8B3/8E1/8H7/8G+/8Bv/8O9/8F+9//Ef7/Dv3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/Bt//AX//Ef3/Avv9/wP3/wTm/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v8I3/8X3/8Ff/8F3/8B7/8Qv/8P+/8I9/8E/f8F5/8E7/8Q3f8H9/8N3/8F9/8Kvv8N/f8Bf/8Q+/8D+/8L7/8D7/3/Avn/Ce//Bd//B3/3/wLf/wHx/wn+/wPv/wJ//f7/Bff/Dff/Bt//Hn//AX//Cv7/Bvf/Bb//BO//B/v/Bt//Bf7/A/3/BL//C9//AT//A+//Bvv/At//Bvf+/xW//wW+/w37/wb+/wL79/8V7/8E39//BX//AX/6H55eqmBYf/8B+9//AQ==",
		"r9lABgAAAAAZDkEGAAAAAAH/Bf7/Dt9//wLv/wF//wj+/wif/wu//wPf/wT+f/8Df/8P9/8L/v8D9/8S7/8I9+//Cvb/Luf/Cv7/Cnv/A5//CL//CN//A9//Cfv/Bff/CP3/CO//BO//Bd//Bvf/Cb//A3v/Bff/Dt//A+//BP7/At//Ar//Au/f/xT+f/8Cv/8M7/8B7v8Mv/8d+/8I/v8Ef/8D/f8B7/8B+/8B/v8Bv/8K3/8I/f8E37//AX//Ab8//v8B/f8B+/8F+/8D7/8C/v8C+f8F+/8C/f8Bv/8N3/8F+/7/A/f/Afr/Bn9//wj3/wF//wr9/xh//wN//wT7/wLX/wP7/wTv+/8T+/8D9/7/Cb//B/3/Ar3/Ff3/BP3/BPv/Af3/BL//CP3/Bvv/Au//B7//Bv7/Dv3/C/v9/wX3+/8Ev/8Cv/8Jv3//Af3/Aff/CPv/BPv/A/b+/wj73/8K9f8C7/8C/f8I9/8Q3/8Bv/8B/f8J+/8Gf/8B7/8Hv/8D7/8B3/8L/f8P/v8G/v8Bv/8T+6//Ad//BNf/B+//Bvv/EPf/Bfvf/xH+/w79/xL79/8F7t//Avv/A/f/B9/3/wH+/wh//wn3/wf9/wL7/wT3/wTm/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v8P/v8Q3/8Ff/8F3/8B7/8K7/8Fv/8P+/8Fv/8C9/8E/f8F5+//A+//EN3/B/f/A7//Cd//Bff/Cr7/Df3/AX//EPv/A/v/C+//A+/9/wL7/wnv/wXf/wd/9/8C3/8B8f8J/v8D79//AX/9/v8F9/8E+/8I9/8G3/8ef/8Bf/8K/v8G9/8Fv/8E7/8O3/8F/v8D/f8Q3/8BP/8B9/8B7/8G+/8C3/8G9/8Wv/8Fvv8U/v8C+/f/Cfv/EN/f/wV//wF/+h+eXqpgWH//Afv/Ag==",
		"wKBABgAAAAAjDkEGAAAAAAF//wT73/8B/v8M/P7/Bu//Au//Eu//BH//Avv/B/3/Avv/BL//BN//Avt//wJ//wK//wN//wH3/xXf/wf3/wJ/7/8Gf/8Ev/f/B/v99/8Zf/8J3/8C/v8G+99//wj1/wb+/wL3/wH9/wT+/wLf/wKf/wnv/wTv/wP3/wL9/wG//wF//wK//wL9/wL7/xfv/wG//wH3/wX7/we//wH7/wvf/wH9/wL+/v8C7/8F3/8Bv/8C++f3/wL9/w99+/8C/f8D7/8E+/8C3/8C7/8H7/8If+//D/7/CN77/wT+/wL2/wT33/8D3/8D/v8B+/8Uv/8D/f8M5v8ee/8J9/8C5/8R3/8E/fv/Cvz/Ff3/CH//B7//Ae//Af7/EPX/B7//BP3/A3//BL//Auv/B/3/An//A7//C/f/AX//Be//Aff/Bn//Af7/DP7/BPn/Fv3/Hvf/D/f/Dv3/Afv79/8C93//Af3/Ar/+9/f/Ae//A3//Av7/Avf/Ad//Cfv/Ef73/wv3/xrf/wK//wEfv/8C+/8K9/8G9/8Gt3//A9//EO3/A3//Ae/+/wR//wj1/wL9/f8D/v8E/t//C/v/B/t//wr+/wnf/f8E3/8S+/8C63//Au//Ad//A/v/B7//BvP/Dfv9/wL7/xHf9/8B3/8F7/8H/v8I7/8H7/8N/f8GX7//Bt/7/wLv+v8P7/8Ff7//Ee/+/wr7/wH9/wW//wy//wG//w3z/wH3/wv9/wH7/wHf/wLv/wa//wv+/wj8/w3f/wL3/wH9/wb7/wX++/8E7/f/Av7/A/n/Dr//C3//BO//An//A/f/B3//Be//Af7v/wf+/v8G/f8C7/8Cf/8Q/v8C7/8T7f7/A/f/Ab//Ae//AX//Bv3/Dd//Br9/f/8Dv/8L9/8O/f8Nf/f/Dd//A1//Du//Cd//Bff/A/v/D/v/C/3/Av7/C7//Bb//CH//DO//Gf7/Cd//Af7/BX//Ae//AvX/AX/+/xvv/wXv/wTf/wbv/wHv/wT3/wr+/w7ff/8C7/8Bf/8I/v8Iv/8Lv/8I/n//A3//D/P/C/7/A/f/Eu//CPfv/wr2/y7n/wr+/wp7/wOf/wR//wO//wjf/wL73/8J+/8D/f8B9/8I/f8Bv/8G7/8C9/8B7/8F3/8G9/8Jv/8D+/8F9/8O3/8D7/8E/v8B79//Ar//Au/f/wP9/xD+f/8P7/8B7v8Mv/8d+/8I/v8Ef/8C9/8C7/8B+/8B/v8Bv/8D/f8G3/8E9/8D/f8E3/8Cf/8BvT/+/wH9/wH7/wnv/wL+/wL5/wX7/wL9/wG//w3f/wX7/v8D9/8B+v8Gf3//CPf/AX//Cv3/DH//D3//BPv/Atf/Ab//Afv/BO/7/wXf/xH3/v8Bf/8Hv/8H/f8Cvf8P/f8F/f8E/f8G/f8Ev/8I/f8G+/8C7/8Hv/8G/v8E/v8J/f8L+/3/Av3/Avf7/wS//wK//wm/f/8B/ff3/wj7/wT7/wP+/v8I+9//CvX/Au//Av3/CPf/EN//AZ//C/v/CO//B5//Bd//C/3/D/7/Bv7/Ab//E/uv/wHf/wTX/wfv/wb7/w37/wL3/wX73/8N7/8D/v8I9/8F/f8S+/f/Be7f/wL7/wP3/wff9/7+/wh//wb3/wr9/wL7/wH9/wL3/wTm/wr3/wH+/wF//wH3/b//Am/1/bv/BO//Bv7/IN//BX//Bd//Ae//Bfv/Cr//Av7/DPv/Ar//Bff/BP3/Bef/BO//B/f/CN3/B/f/Dd//Bff/Afv/CL7/Af3/C/3/AX//EPv/A/v/CPf/Au//A+/9/wG/+/8Jz/8F3/8Ff/8Bf/f/At//AfH/Cf7/A+//An/9/v8F9/8N9/8G3/8ef/8Be/8Kvv8G9/8Fv/8E7/8O3/8F/v8D/f8Gf/8J3/8BP/8D7/8G+/8C3/8G9/8V77//Bb7/FP7/Avv3/wj3/wz3/wTf3/8Ff/8Bf/ofnl6rYFjAH/v/Ag==",
		"x8BABgAAAAAjDkEGAAAAAAH+/wTf/wTf/wXf7/8R+/8C63//Au//Ad//A/v/At//BL//BvP/Dfv9/wL7/xHf9/8H7/8L/v8E7/8H7/8N/f8GX9//Bt/7/wLv+v8Pz/8Ff7//Ee//C/v/Af3/Bb//DL//Ab//DfP/Aff/C/3/A9//Au//Ev7/CPz/Dd//Avf/Af3/Bvv/Bf77/wX3/wL+/wP5/w6//wt//wH+/wLv/wJ//wP3/wd//wXv/wLv/wf+/v8G/f8C7/8T/v8C7/8T7f7/A/f/Ab//Ae//AX//Bv3/Dd//Br9/f/8Dv/8L9/8O/f8I+/8Ef/f/Dd//A9//Du//D/f/A/v/G/3/Av7/C7//Bb//Bfv/An//Cb//Au//B/f/Ef7/Cd//Af7/BX//Ae//Av3ff/7/CPv/Cvv/B+//Ct+//wXv/wb3/wr+/w7ff/8C7/8Bf/8I/v8Iv/8Lv/8I/n9//wJ//wy//wL3/wv+/wP3/xLv/wj37/8K9v8u5/8K/v8Ke/8Dn/8Iv/8I3/8D3/8J+/8F9/8Bf/8G/f8E/f8D7/8E7/8F2/8G9/8Jv/8D+/8F9/8O3/8D7/8E/v8C3/8Cv/8C79//BX//Dvp//w/v/wHu/wn9/wK//wt//wr3/wb7/wPf/wT+/wR//wXv/wH7/wH+/wG//wrf/wj9/wTf/wJ//wG/P/7/Af3/Afv/Ce//Av7/Avn/Bfv/Av3/Ab//BPv/CN//Bfv+/wP3/wH6/wZ/f/8I939//wr9/wG//xP7/wZ//wT7f/8B1/8D+/8D/u/7/wvf/wv3/v8Jv/8H/f8Cvf8B+/8M9/8G/f8E/f8G/f8Ev/8I/f8G+/8C7/8Hv/8G/v8O/f8H9/8D+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wP+/wT3/xDf/wG//wv7/wS//wPv/w3f/wv9/w/+/wb+/wG//xP7r/8B3/8E1/8H7/8G+/8K3/8F9/8F+9//Ef7/Dv3/Db//BPv3/wXu3/8C+/8D9/8C/v8E3/f/Af7/A/3/BH//Ef33/wH7/wT3/wTm/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v8g3/8Ff9//BN//Ae//EL//CPv/Bvv/Bd//Avf/BP3/Bef/BO//EN3/B/f/Dd//BXf/Cr7/Df3/AXf/Cu//Bfv/A/v/C+//Af3/Ae/9/wL7/wd//wHv/wXf/wd/9/8C3/8B8f8Cv/8G/v8D7/8CX/3+/wX3/w33/wbf/x5//wF//wr+/wXf9/8Fv/8E7/8O3/8F/v8D/f8Q3/8BP/8D6/8G+/8C3/8G9/8J7/8Mv/8Fvv8U/v8C+/f/Gt/f/wV//wF/+h+fXqtgWMAf+/8C",
		"zZtABgAAAAAjDkEGAAAAAAEL+/8B9/8G9/8Gvf8Cv/8B/v8B+/8G+/8E7/8B/v8E98//Ct//C7//A/v/AX//A7v/Cvf/Aff/At8//w/+/wW//wH+/wTf/wP7/wnv/wP3/wl/f+//Bn//BPvf/w78/v8J7/8Xf/8C+/8H/f8C+/8Ev/8E3/8Df+//AX//Ar//A3//Aff/Af3/G/f/Af1/7/8Gf/8F9/8H+/33/xl//wX7/wPf/wL+/wb7/wF/3/8H9f8G/v8C9/8G/v8C3/8Cn/8J7/8Cv/8B7/8G/f8Bvv8Bf/8F/f8C+/8X7/8Bv/8B9/8F+/8Hv/8N3/8B/f8C/v7/Au//Bd//Ab//Avvn9/8C/f8Pffv/Bu//BPv/Be//Cv7/BX//EP7/CE77/wT+/wL3/wT33/8D3/8F+/8Uv/8D/f8Lv+b/Hnv/DOf/Ed//BP37/wr8/wHv/xP9/wh//we//wHv/xL1/wz9/wN//wS//v8B6/8Kf/f/Ar//C/f/AX//Be//Aff/Afv/BH//E/n/BP7/EPf9/y73/w79/wH7+/8D93//Af3/Ar/+9/f/BX//Av7/Avf/Ad//Cfv/Ef73/wv3/xrf/wK//wGfv/8B/vv/Cvf/Bvf/Bvd//wPf/wH7/w7P/wN//wHv/v8H/v8F9f8C/f3/A/7/Bd//E/t//wr+/wnf/wXf/wn+/wj7/wLrf/8C7/8B3/8D+/8Hvv8G8/8N+/3/Avv/Ed/3/wfv/xDv/wfv/w39/wZf/wff+/8C7/r/CO//Bu//BN9/v/8I/v8I7/8Ev/8G+/8B/f8Fv/8Mv/8Bv/8F3/8H8/8B9/8B/f8J/f8D3/8C7/8G3/8L/v8I/P8N3/8C9/8B/f8E7/8B+/8F/vv/Bff/Av7/Ae//Afn/De+//wr3f/8E7/8Cf/8C9/f/B3//Be//Au//B/7+/wb9/wLv/xP+/wLv/xPt/v8D9/8Bv/8B7/8Bf/8G/f8N3/8Gv39//wO//wv3/w79/w1/9/8N3/8D3/8F7/8I7/8P9/8D+/8Pv/8L/fv/Af7/C7//Bb//CH//DO//Av3/Fv7/Bvv/At//Af7/BX//Ae//AX/9/wF//v8b7/8K3/8G7/8G9/8K/v8O33//Au//AX//CP7/CL+//wq//wj+f/8Df/8Kv/8E9/8H+/8D/v8D9/8Gf/8L7/8Bv/8G9+//BP3/Bfb/GL//Fef/Cv7/Cnv/A5//CL//CN//A9//Cfv/Bff/CP3/CO//Av3/Ae//Bd//Bvf/Cb//A/v/Bff/CH//Bd//A+//BP7/Af3f/wK//wLv3/8L7/8I/n//A/7/C+//Ae7f/wu//xz9+/8I/v8Ef/7/BO//Afv/Af7/Ab//Ct//B/v9/wTf/wJ//wG/P/7/Af3/Afv/Ce//Avr/Avn/Bfv/Av3/Ab/9/wv33/8F+/7/A/P/Afr/Bn9//wj3/wF//wl//f8N9/8Of/8E+/8C1/8D+7//A+/7/xf3/v8Jv/8H/f8Cvf8B3/8Q+/8C/f8E/f8F/f3/BL//CP3/Au//A/v/Au//B7//Bv7/Dv3/C/v9/wX3+/8Ev/8Cv/8Jv3//Af3/Aff/CPv/BPv/A/7+/wj73/8K9f8C7/8C/f8I9/8Q3/8Bv/8L+/8I7/8Hv/8F3/8L/f8P/v8G/v8Bv/8B9/8N/f8D+6//Ad//BNf/B+//Btv/EPf/Bfvf/xH+/wf3/wb9/xL79/8F7t//Avv/A/f/B9/3/wH+/wh//xH9/wL7/wT3/wK//wHm/wr3/wH+/wF//wH3/f8Db/X9u/8E7/8G/v8Mv/8T3/8Ff/8F3/8B7/8Qv/8P+/8I9/8E/f8F5/8E7/8G/f8J3f8H9/v/DN//A+//Aff/Cr7/Av3/Cv3/AX//CP7/B/v/A/v/C+//A+/9/wL7/wnv/wXf/wH9/wV/9/8C3/8B8f8G9/8C/v8D7/8Cf/3+/wX3/w33/wbf/x5//wF//wr+/wb3/wW9/wTv/wvf/wLf/wX+/wP9/xDf/wE//wPv/wb7/wLf/wb3/wvf/wq//wW+/xT+/wL79/8a39//BX//AX/6H59fq2BYwB/7/wI=",
		"+nlABgAAAAAjDkEGAAAAAAH/AfqyDQABSgAB3/8I3/8Ff/8F/v8Fr/8F7/8F/f8H7/8C/fv/CO//E9//Cn/+/wTf/wf3/wHv/wXf/wz9/f8Lf3//Fvf/At//A/7/A+3/Au//Ar/+7/8B7/8R+f8Vf/8Dv3//CPf7/wHf/x733/f/Aef/Avf/Bd/v/v8B+/8Q3/8D/v8n7/8K/v8D9/8F3/8Fv/8o/u//Af7/BO/3/f8Cv/8B+3//It/f/wO//wvvf/8L+/8I9/8Df/v/Af3/A/v/Arv/BP7/Du//Dn//Ar//C79//wLz/xj+/wa/v9//A/3/Aff/Cb//Br7/Efvv/w/vv/8D/f8L+/8I9/8H/v8D3/8G/v8D9/8G/v8B3/8B9/8N3/3/Bf7/Bn//B/f/D7/9/wf9/wXv/xr7/wb7/wfv/w33/f8Cf/8Lf/8C3/8Gd/8Cf/8L+/8P+/8Lf3//B7//B/7/Ab//Gf3f/wG/f/8D+8//A/3/A/7/En//Au+//wH3/wPb/wT9/wb3/wI2/wF//wbf/wj9/wK//wHf/wvv/xPv+/8B9/8G9/8G/f8Cv/8B/vv7/wvv/wH+/wT3z/8K3/8Lv/8D+/8Fu/8K9/8B9/8C3z//D/7/Bbf/Af7/BN//A/v/Ce//A/f/CX9//wd//wT73/7/Dfz+/f8I7/8Fv/8Rf/8C+/8H/f8C+/8Ev/8E39//An//An//Ar//A3//Aff/DP3/EPf/An/v/wz3/wf7/ff/GX//A7//Bd//Av7/Bvv/AX//CPX/Bv7/Avf/BPf/Af7/At//Ap//Ce//BO//Bv3/Ab//AX//Bf3/Avv/F+//Ab//Aff/Bfv/B7//Dd//Ad3/Av7+/wLv/wXf/wG//wL75/f/Av3/D337/wbv/wT7/wXv/xB//xD+/wjO+/8E/v8C9/8E99//A9//Bfv/FL/v/wL9/wzm/xz+/wF7/wzn/xb9+/8K/P8F7/8P/f8If/8Hv/8B7/8S9f8M/f8Df/8Ev/8C6/8Kf/8Dv/8L9/8Bf/8F7/8B9/8F93//E/n/Fv3/Dff/GO//B/f/Dv3/Afv7/wP3f/8B/f8Cv/739/8Ff/8C/v8C9/8B3/8J+/8J/v8H/vf/C/f/Gt//Ar//AZ+//wL7/wr3/wb3/wb3f/8D3/8Q7/8B3/8Bf/8B7/7/DfX/Av39/wP+/wXf/xP7f/3/Cf7/B/v/Ad//Bd//Evv/Aut//wLv/wHf/wP79/8Gv/8G8/8B+/8L+/3/Avv/Ed/3/wfv/xDv/wfv/w39/wZf/wff+/8C7/r/D+//BX+//xHv/wV//wX7/wH9/wW//wy//wG//wu//wHz/wH3/wv9/wPf/wLv/xL+/wj8/w3f/wL3/wH9/wb7/wX++/8F9/8C/v8D+f8Ov/8F+/8Ff/8E7/8Cf/8D9/8Hf/8F7/8C7/8H/v7/Bv3/Au//Cr//CP7/Au//Be//De3+/wP3/wG//wHv/wF//wb9/wz+3/8Gv39//wL+v/8L9/8F7/8I/f8Nf/f/Dd//Aff/Ad//DP3/Ae//D/f/Ar/7/xHf/wn9/wL+/wu//wW//wh//wzv/xn+/wnf/wH+/wV//wHv/wL9/wF//v8b7/8F9/8E3/8G7/8G9/8K/v8G3/8H33//Au//AX//CP7/CL9//wq//wj+f9//An//D/f/C/7/A/f/Eu//CPfv/wr2/yp//wPn/wl//v8Ke/8Dn/8Iv/8I3/8D3/8I9/v/Bff/CP3/CO//BO//Bd//Bvf/Cb//A/v/Bff/Dt//A+//BP7/At//Ar//Au/f/xT+f/8P7/8B7v8B3/8Kv/8d+/8I/v8Ef/8F7/8B+/8B/v8Bv/8K3/8D+/8E/f8E3/8Cf/8Bvz/+/wH9/wH7/wnv/wL+/wL5/wX7/wL9/wG//w3f/wX7/v8D9/8B+v8Gf3//Bf3/Avf/AX//Cv3/E7//CH//BPv/Atf/A/v/BO/7/xf3/v8D9/8Fv/8H/f8CvP8V/f8E/f8G/f8Ev/8I/f8G+/8C7/8Hv/8G/v8O/f8L+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wX7/wL3/xDf/wG//wv7/wjv/w3f/wv9/w/+/wb+/wG//xH7/wH7r/8B3/8E1/8H7/8G+/8I3/8N+9//Efb/Dn3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/CH//EfX/Avv/BPf/BOb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wP7/wL+/yDf/f8Ef/8F3/8B7/8F/v8Kv/8P+/8Cf/8F9/8E/f8F5/8E7/8E7/8L3f8H9/8N3/8F9/8Kvv8C+/8K/f8Bf/8Df/8M+9//Avv/C+//A+/9/wL7/wnv/wXf/wdv9/8C3/8B8f8G+/8C/v8D7/8Cf/3+/wX3/w33/wbf/xy//wF//wF//wr+/v8F9/8Fv/8E7/8C9/8L3/8F/v8D/f8Ff/8K3/8BP/8D7/8G+/8C3/8G9/8Wv/8Fvv8U/v8C+/f/Gt/f/wV//wF/+h+fX6tgWMAf+/8C",
	}

	fmt.Printf("%4s   %4s   %4s   %4s   %4s\n", "orig", "1B", "2B", "4B", "8B")
	fmt.Println("==================================")
	for i, dataString := range testData {

		kr := &KnownRounds{}
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			t.Errorf("Failed to decode marshalled known rounds: %+v", err)
		}

		err = kr.Unmarshal(data)
		if err != nil {
			t.Errorf("Failed to unmarshal known rounds: %+v", err)
		}

		buff := kr.bitStream.marshal1Byte()
		u64b := unmarshal1Byte(buff)
		f1bLen := len(buff)
		if !reflect.DeepEqual(kr.bitStream, u64b) {
			t.Errorf("Failed to marshal and unmarshal 1 byte buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, kr.bitStream, u64b)
		}

		buff = kr.bitStream.marshal2Bytes()
		u64b = unmarshal2Bytes(buff)
		f2bLen := len(buff)
		if !reflect.DeepEqual(kr.bitStream, u64b) {
			t.Errorf("Failed to marshal and unmarshal 2 bytes buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, kr.bitStream, u64b)
		}

		buff = kr.bitStream.marshal4Bytes()
		u64b = unmarshal4Bytes(buff)
		f4bLen := len(buff)
		if !reflect.DeepEqual(kr.bitStream, u64b) {
			t.Errorf("Failed to marshal and unmarshal 4 bytes buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, kr.bitStream, u64b)
		}

		buff = kr.bitStream.marshal8Bytes()
		u64b = unmarshal8Bytes(buff)
		f8bLen := len(buff)
		if !reflect.DeepEqual(kr.bitStream, u64b) {
			t.Errorf("Failed to marshal and unmarshal buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, kr.bitStream, u64b)
		}

		origLen := len(kr.bitStream) * 8
		fmt.Printf("%4d   %4d   %4d   %4d   %4d\n", origLen, f1bLen, f2bLen, f4bLen, f8bLen)
		fmt.Printf("       %4.0f%%  %4.0f%%  %4.0f%%  %4.0f%%\n",
			100-float64(f1bLen)/float64(origLen)*100,
			100-float64(f2bLen)/float64(origLen)*100,
			100-float64(f4bLen)/float64(origLen)*100,
			100-float64(f8bLen)/float64(origLen)*100)
		fmt.Println("----------------------------------")
	}
}

// func TestKnownRounds_Marshal2(t *testing.T) {
// 	data := `{"BitStream":"////v6jgxPkU///v////////Af/3/////////wH/3//t//////8B//////////3/Af/////r8////wH/////v//f//8B/////////f3/Af///////////wEB//v/////////Ae/////f////+wH///////3///8B//v////3v///Af7/9//7//7//wH///////2///8B////////////AQG///2////+//8B7/////3/////Afv////////7/wG///////////8B//ef////////Af//3/////3//wH/////9/+//v8B////////////AQH////f////+/cB/f//////3///Af//////////3wH////////3//8B///q/f//////Af///e///////wH///////////8BA///+///7////QH///////7///8B//////3/////Af////3//////wH/7/////7///8B//////3/////Ab///////////wH///////////8BAf/3/////////wH/5/////////8B+/////////+3Af/f/////////wH//9///+////8B/////7////f/Af//v/f//////wH/+/////////8B/7//////////Af/f9////////wH///////////8BAf//v////7///wHf//////////8B/////////7//Af//////////7wH////////f//8B////////////AQH////////3/v8B//+///////f/Af/////+/////wH/7/////f+//8B//+////7//7/Af/////////9/wH/////////9/8Bv////7//////Af//3////////gH////3//+//v8B/7v//////+//Af///////////wEB//7/9/f/////Af////////7//wH//9//9/+///8B/////////+//Af///////////wEB///////19v//Af/97////7///wH//////f///f8B//v///7//9//Af///////////wEB///////////7Af/////////P/wH//////9///7sB////7/v///9e//f////////zAev///+//////wH////+//////8B///v////////Af///////////wEB//////+7////Af//////+///vwH///////////8BAr/////3//v//wH///+///////8B/f////v/////Af////7//9//2gH/////3/////8B////v/v/////Af/////3//f//wH/7/////////8B////+//+7///Af//3////////wH////////7//8B/9/3//////9///3/////////Af/+/////////wH7//f///////8B///7/9//////Af///f//+////wH//7//////7/8B///r///1////Af///////////wEB///f////////Ab///97//////wH///////////8BAf/////7/f///wH//f///9/7/78B//7////7/7//Af///////////wEB/////f+/////Af///////////wEB///+v//////7Af/////+3////wH///////////8BAv+///////fv/wH///7/7////+8B////////////AQH//v////////8B///////en//3Af///////////wEB/////9//v///Af3/////3////wH///////////8BAf//97//u////wH/////+/////8B/7//////////Af///////////wEB//v//////f/7Af//////96///wH/////9/////8B1/v3//f/////Af/////////7/wH///f9//////8B///v/////7//Af//////3/3//wH///////////8BAff////////97wH//9////+//+8B////////////AQH9//////////8B////////7///Af///f/23//9zwH///////////8BAf/3////v////wH///v//f////8B//////7/////Af///////////wEB//////7/////Af///+///////wH//7/7///7//8B/////////f//Af+//f///////wH///////////cB/9//////////Af/7///7////vwH/////9////f8B/+//////////Af//////9///9QH//7////////8B////////////AQH93//3/////P8B9/v/9//////fAf/f////t///9wH///////f///8B//////////v/Af//////v////wH///////////8BAf////////+//gHf/////////f8B/+//////////Af+/////v////wH////+//////8B///9///////vAf///v///////wH3//////////8B//v/////////Af/9/////////wH7/7//3/////8B/////////v//Ad/////v/v2//wH////////3//8B///////////PAf//9///9////wH/////+////v8B//3/////////Ae//97/////3/wH////////9//8B/////////+//Af///////+7/vwHf////9/////8B/////////f//Af/7////////7gH///////f///8B///+//3/////Af//////////7wH///////////8BAf///e///////wH///////f///0B////3///////Af///6+/////9wH///7/////v/sB///v///3////Aff//////f///wH//fv//9////8B////////////AQH/////+v/9//4B////////////AQL+/7////////8B////3///7///Af///f//3////wH/9/////3/9/8B/////////9//Af///f///////wH/9/////+///8B//////3/4///Af//+///9////wG9+////////+8B///////////3Af//9/+//////wH////3//////8B/////7/3////Af/f///X///f7wH///////////8BAf////+//v//9wH//7//////7/8B//7//////f/uAf/////f/////wH//f/////+//8B///////37/v/Af///f///////QH//////////78Bv///////////Af/nv///////vwH//////f/3//8B//3/////////Af///////7///wH7//////////cB3//////v////Af//3////////wG/v//////r//8B////////v///Af///////9///wH/+/7//v/v//8B////////////AQH/3u////////8B////////////AQH//+///f////8B/////////9//Ad///////////wH/3//3//////8B////////////AQH////9///3//8B9///////////Af/////////3/wH///////////8BAfv/+////////wH//////9//v/8B+v//////////Aff//////+//3wH///3f+//9//4B///f////////Ab///f/////z/wH///////////8BAf/////////3/wH///////+///8B////////////AQH//////////78B3////f3///v/Af//////////+wH89//e//////8B////////s///Af///////f//+wH7////v/////8B3/////2/////Af//9////////wH/v///////9/8B//+//7////3/Af//3/7/v//u/wH///////////8BAf/////777///wH7/////v///98B////////////AQH/////7/7//fsB//f///////f/Af/////////f/wG+/////9////8B/v///v7///9//9/////+////Af+///v//////wH///////+///8B////////////AQH/+/////////8B////////3///Af//////3////wH/////+/////8Bv//////////vAf////////3//wHu/9////////8B///////v//9/3/////////7/Af///////////wEB////////7///Af///////////wEB///3////////Af/////////7/wH//////////98B////////////AQH////3//3///8B////////////AQH//////v////8B/b/3////////Af///////////wEB/////////9//Af/+////////f+////v/////+wH///////////8BAv///9/////7/wH////d//////8B/f//////////Af///////////wEC//////////v/Af//9/////v95wH/////////3/8B//////v////fAf/u7v///9//+wH///////////8BAd//////9////wH///////////4B//////3//f27Af///////////wEB///9/+v//f/+Af/////73////wH///////////8BAf/9/////////wH//////v////8B/////v//////Af//+////////wH/v/////73//8B/7///+//////Af//7/v//9//7wH///+///////8B////////////AQH/3//7//////8B////7//+//+/Af//8//7v//7/wH///37//////8B//////////f/Af//////9//v+wHf///+///3//8B////////7//fAf/v/////////wH/////9/////8B////////////AQL////v//3///0B////////////AQH/////+//v//8B/////////+//Af/f//f//////AH+/f3///v/3/4B////////3///Af3///f/5////wH////9//////8B//+////////fAdf//////////wH///////////8BA///////////+wH7/v////////8B/////////v//Af////2//////wH/////3////zf//f/3/////98B/P////////+/Af////f//////wH//////////98B///f////////Af//+////////gHf///////+//8B////////7/P/Af//////+////wH///////////8BAf//v////////wH///////7///8B///v/f/////9Af///////9///QH//////////vsB/9//v///////Af///+/9///9/wH///////////cB/t7/3v//////Af/////////v/wH///////7///0B/////f//////Af//////+////wH///+7////3/8B/7/3+//93///Ab/v////v7///wH//v////////4B///f/9//////Af7v/7///f///wH7/////////v8B///7////////Af//9////////wH///////+///8B////////////AQH///+///3///8B////////////AQH/7////9f///8B///99////f/9Af///////+///wH///////////8BAf///9////77/wH///////////8BAf/f/////////wH/vb////3//78B////m////v3/Af//7/////7//wH///f///////8B9///7///////Af///////////wEB//////7/////Af/7///////v/QH/+//7//////8B//+9///9+ve/Af////f9/////wH///////////8BAf/f//+//9///gH+/9//3/////8B//+///v////3Af///v///f//nwH/+/////////4B/+/////////vAf///////////wEB9/+f/////+//Af/////////7+gH///////////8BAf///d/7///v/wH///////////8BAf/e/////9///QH///////////8BAv///+//v/7//wH+///7//////8B////////9v//Af///////////wEB///f///7///vAf////////+//wH/////v////+8B//////////f/Ab//7/7///+//wG///////7/+/8B///9////////Af/f//f////p/wGf/9//v/////8B/7//////////Af/f/v//+////wH7////////9/8B////////////AQH///////7///8B3e///f//////Af/7/////////wH///////////8BAf//n/7//////wH/vf////7///8B/////////f//Af///////////wEC/7/9/v//////Af//3//fvv/e9wH+////+/////8B//////7/////Af///f///////wG/////v//77/8B//7/////v///Af////+9///7/wH/////////+/8B///8////7/9////7/+//9///Af///////////wEB///////+/9//Af///////f///wH///f//+////cB////////////AQH9///////+//8Bv//7//////r8Af/+///7/////wH/3/////b///8B//////+/////Af/f/////////wH////9//7/3f8B/////9///9d///////////9/////3/////z9Af//7//3///7/wH///////////0Bv///////////Af///////////wEB3////////7//Af3///+//9/9X//3/////////wH//e///////f8B////v//9//+/Aff////d///9/wH//////v37//8B7///////v///Af//////////f/////f////+/wHv///f///9//8B//////7/////Af///v/v/////wH///////////8BAv/////+/////wH///////////4B////////////AQL/+/v/////v/8Bv//x9//77///Af7//7///9/7/wG//t///////3/+///////9//8B///v//fv////Af///////////wEC///////9////Af/+/8f//////wH///////////8BAf3//v/+//3/7wH///////////8BAf///////7/99wH////f//////8B///9/////v//Af//v/3//v///wH/////v/////8B///////+////Afz/////v////wH/////+/7///8B/////////7//Af///////v///gH//+/f//3///8B/9/////9///6Af///////////wEB2///////////Af//v7///////wH//8eZ//3///8B//3/+v//////Af//+////////wH///v7/////v8B/v/7//////9///v//7//////Af/////////3/wH//f////////8B/////+/9////Af//v///3////wH//9///7////8B////////////AQH/////9/////8B//3//////+//Af///////////wEB//39////////Af//////////7wH///////////8BAf////v//////wH+////////z/8B////////////AQH////+/v////8B///9///7////Af/Y/v///+///wH////////f//8B////////////AQH/3//3v//+//8B3//////+///vAf/////7///7/wH///////////8BAf/7///d/////wG///////////8B/////f//////Af///////////wEB//+///vf////Af///v/+///f/wH//9////f///8B///s///93///Ab/////////3/wH///////+///8Bv///////////Af///////////wEB/9///f//////Af/////////+/wH//////r////8B////////v///Af/f/////////wH//////+f/3/8B////7//7//+/Af//////////+gH////////3//8B////v//f////Af///////////wEB///////7//+7Af///////////wEB3//9////////Af/////f//+//wH///////////sB9/v/////////Af/////+/////wHv/f////v///8B////////////AQH////7/v////sB/f//////////Af+/////////7wH/+/////////8B///+//////f/Af7X/////////wH//////v/u//8B/////////7//Af////f//////wH//////////98B////////vf//Af//////7//3/wH3//////////cB/9///d77////Af////+///v//wH////////f//8B////v//f////Af3/9//d/////wHr///3/9/v//8B/////9//////Af///////////wEB9//9///9////Ae///9//////3wH///////////8BAf/////////+/wH3//7//////v8B///9////77//Af/+/////////wH3///777//3/8B///////+7/9///////////f/Af//37f///+//wH///////////8BAf//u//////v/wH/////29//+78B//f//////v//Af/////+/////wH/9////////3////////7//3///////v////8B+/33////////Ab///+/7/v///wH///77//////8Bvv////3f///3Af//9////////wH/v/////////cB3///////////Af//+/////39/wH7//////////cB/////////9v/Af//9////P///wH////9v///+/8B///z/////v//Af/////////v/wH////fv/////4B//n////////7Af//vvv////3/wH//////P//9/8B///////v///+Af+////////v/wH///////////8BAf7v/////////wH9/////v//v78B/////+//////Af/v//z+///39wH///////////8BAf+////////v7wH/////9/////8B/////7//3//9Af/////////v7wH///////////8BAf/33///9////wH//////7///74B/////7f/////Af/3///////+/wH///////////8BAf/////9/v///wH////y//////8B///+///n////Af//9v////f+/wH////////v3/8B/////////7//Af/++7///f///wH///////////8BAvf3/////////wH////v//////8B9///5///v//fAf/////f/f///wH///////////8BAf//+///v///7wH///ff9////3//////////3/4B////v/+/////Af/f/////+/3/wH///////////sB/////////9//Af///v//+v///wH3//////////8B///X//v///f/Aff/////////f/////26/////wH////////+//8B/f//////v///Af////f//+///wH2/////9////0B//7//////9++Af///////t///wH///////2///8B///e///f/v//Af///////////wEB/7//////////Af////////+//wHf//////////8B/v//v+//////Af//////+////wH7//////vv/78B////////////AQH////vv+//13/////P3/////8B+///////+///Af///////////wEB///f//+/////Afn///////v+/wH////X/f//v78B/9////7/////Af6///////f/9wH//////////33//u///7////8B//vz//u/////Af//////7////wH///////7/+/8B//////////7/Af//v//////9/wH///////v///8B3///7///3//5Af+//////////wH////9/7////8B/////////v//Af///f3/////vwH7//////////8B////v/3///79Af///////////wEB////x//////zAf//////v////wH7/////7////8B////////////AQH///7v9//fv98B/////////v9/7///////v///Af7//f///9/+/wH/nffZ////9/oBlu+Yvvvl6eo76///v8CogJB0","FirstUnchecked":94919909,"LastChecked":94969701}`
// 	kr := &KnownRounds{}
//
// 	err := kr.Unmarshal([]byte(data))
// 	if err != nil {
// 		t.Errorf("Unmarshal() returned an error: %+v", err)
// 	}
//
// 	t.Log(kr)
// 	t.Logf("%064b", kr.bitStream)
//
// 	t.Log(kr.Checked(94969696))
// }

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
