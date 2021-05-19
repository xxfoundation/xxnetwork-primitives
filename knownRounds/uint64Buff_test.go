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

		buff := data.marshal1ByteVer2()
		u64b := unmarshal1ByteVer2(buff)
		f1bLen := len(buff)
		if !reflect.DeepEqual(data, u64b) {
			t.Errorf("Failed to marshal and unmarshal 1 byte buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, data, u64b)
		}

		buff = data.marshal2BytesVer2()
		u64b = unmarshal2BytesVer2(buff)
		f2bLen := len(buff)
		if !reflect.DeepEqual(data, u64b) {
			t.Errorf("Failed to marshal and unmarshal 2 bytes buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, data, u64b)
		}

		buff = data.marshal4BytesVer2()
		u64b = unmarshal4BytesVer2(buff)
		f4bLen := len(buff)
		if !reflect.DeepEqual(data, u64b) {
			t.Errorf("Failed to marshal and unmarshal 4 bytes buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, data, u64b)
		}

		buff = data.marshal8BytesVer2()
		u64b = unmarshal8BytesVer2(buff)
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
		"c21ABgAAAAC7DUEGAAAAAAIIAgj/Be//CL//Aff/A3//Ab//Dbf/Eff+/wP6/wv7/wv8/wf3/wf+/wi//wT3/wWe3/8Dv/8B8/8Bv/8H9/8B/v8B/v8H7/8I3/8T9/8P7/8I/v8Yr/8Fv/8C/v8Ff/8Ef/8B/X//Cb//Aff/Bu//Bfv/Cb//Dnf3/wd//wN//wP9/xf+/wO//wH9/wr3/wvv/wH3/wb3/wH7/f8Kv/8C7/8Bb/8H3/8D7/8E3/8L/v8Fr/8F7/8F/f8K/fv/CO//Hn/+/wTf/wf3/wHv/wXf/wi//wP9/f8J7/8Bf3//Fvf/At//A/7/A+3/Au//Ar/+7/8B7/8C/v8O+f8Jf/8Lf/8Dv3//At//Bff/At//Hvff9/8B5/8C9/8F3+/+/wH7/xDf/wP+/wPv/yPv/wG//wj+/wP3/wXf/wW//yj+7/8B/v8E7/f9/wK//wH7f/8Gv/8a/t/f/wO//wvvf/8B+/8J+/8I9/8Df/v/Af3/A/v/Ar//BP7/Cff/BO//Dn//Ar//C79//wLz/wz9/wrv/v8Gv7/f/wP9/wH3/wXf/wO//wa+/wj7/wj77/f/Cf3/BO+//wP9/wv7/wj3/wf+/wr+/wP3/wb6/wHf/wH3/w3f/f8F/v8J/f8E9/8Pv/3/B7//Be//BPf/Ffv/A7//Avv/B+//Dff9/wJ//wt//wLf/wZ3/wJ//yV//wF/f/8P/v8Bv/8Z/d//Ab9//wP7y/8D/f8Z77//Aff/A9v/BP3/Bvf/Ajb/AX//Bt//CP3/Ar//Ad//C+//Dv7/BO/7/wH3/wb3/wb9/wK//wH+/wH7/wvv/wH+/wT3z/8K3/8F9/8Fv/8D+/8Fu/8K9/8B9/8C3z//D/7/Bb//Af7/BN//A/v/Ce//A/f/CX9//wd//wT73/8O/P7/Bff/A+//F3//Avv/B/3/Avv/BL//BN//A3//An//Ar//A3/79/8d9/8Cf+//Bn//Bff/B/v99/8Zf/8J3/8C/v7/Bfv/AX//CPX/Be/+/wL3/wb+/wLP/wKf/wnv/wTv/wb9/wG//wF//wX9/wL7/xfv/wG//wH3v/8E+/8Hv/8N3/8B/f8C9v7/Au//Bd//Ab//Avvn9/8C/f8If/8Gffv/Bu//BPP/Be//EH//EP7/CM77/wT+/wL3/wT33/8D3/8F+/8Uv/8D/f8M5v8L+/8Se/8M5/8R3/8C/v8B/fv/B/7/Avz/Ff3/CH//B7//Ae//Ev3/DP3/A3//BL//Auv/Cn//A7//C/f/AX//Be//Aff/Bn//E/n/D7//Bv3/Lvf/Dv3/Afv7/wP3f/8B/f8Cv/739/8Ff/8C/v8C9/8B3/8J+/8R/vf/C/f/Gt//Ar//AZ+//wL79/8J9/8G9/8E/f8B93//A9//EO//A3//Ae/+/w31/wL9/f8D/v8F3/8F9/8N+3f/Cv7/Cd//Bd//Dv7/A/v/Aut//wLv/wHf/wP7/we//wTf/wHz/we//wX7/f8C+/8D/v8N3/b/B+//Dn//Ae//B+/7/wz9/wZf/wff+/8C7/r/D+//BX+//wT7/wzv/wv7/wH9/wW//ww//wG//wr+/wLz/wH3/wv9/wPf/wLv/xL+/wj8/wz33/8C9/8B/f8G+/8F/vv/Bff/Av7/A/n/Ad//Cv7/Ab//C3//BO//An//A/f/B3//Af3/A+//Au//B/7+/wb9f/8B7/8T/v8Bv+//E+3+/wP3/wG//wHv/X//Bv3/Bb//B9//Br9/f/8Dv/8G3/8E9/8O/f8Nf/f/Dd//A9//Du//CO//Bvf/A/v/G/3/Av7/C7//Bb//CH//DO//Gf7/Cd//Af7/BX//Ae//Av3/AX/+/xC//wrv/wrf/wbv/wb3/wr+/w7ff/8C7/8Bf/8I/vv/B7//C7//CP5//wN//wL7/wz3/wv+/wP39/8R7/8I9+//Cvb/Bt//J+f/Cv7/Cnv/A5//CLf/CN//A9//CN/7/wX3/wj9/wS//wPv/wTv/wJ//wLf/wb3/wm//wP7/wX3/w7f/wPv/wT+/wLf/wK//wLv3/8E3/8P/n/7/w2/7/8B7v8I7/8Dv/8d+/8I/v8Ef/8F7/8B+/8B/v8Bv/8J39//CP3/BN7/An//Ab8//v8B/f8B+/8J7/8C/v8B7/n/Bfv/Av3+v/8N3/8F+/7/A/f/Afr/Bn9//wj3/wF//wr9/xx//wT7/wLX/wP7/wTv+/8X9/7/Cb//B/3/Av3/Gfv9/wb9/wS//v8H/f8G+/8C7/8Hv/8G/v8O/f8L+/3/Bff/Bb//Ar//Cb9//wH9/wH3/wj7/wT7/wP+/v8I+9//CvX/Au//Av3/CPf/EN//Ab//C/v/B/7v/we//wXf/wv9/w/+/wb+/wG//xP7r/8B3/8El/8H7/8G+/8Q9/8F+9//Df7/A/7/Dv3/EN//Afv3/wXu3/8C+/8D9/8H3/f/Af7/CH//Ef3/Avv/BPf/BOb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/yDf/wV//wXf/wHv/xC//w/7/wj3/wT9/wXn/wTv/xDd/wf3/w3f/wS/9/8Kvv8Dv/8J/f8Bf/8L3/8E+/8D+/8B/v8J7/8D7/3/Avv/Cev/Bd//B3/3/wLf/wHx/wh//v8D7/8Cf/3+/wF//wP3/w33/wbf/x5//wF//wr+/wb3/wW//wH3/wLv/w7f/wX+/wP9/xDf/wE//wOv/wb7/wLf7/8F9/8F/f8Qv/8Fvv3/C/f/B/7/Avv3/xD+/wL7/n3/AX//AWpX3uy0HZ8=",
		"zZtABgAAAAD8DUEGAAAAAAIIAgj7/wH3/wb3/wb9/wK//wH+/wH7/wvv/wH+/wT3j/8K3/8C/v8Iv/8D+/8Fu/8K9/8B9/8C3z//B/3/B/7/Bb//Af7/BN//A/v/Ce//A/f/CX9/v/8Gf/8Ee9//Dvz+/wnv/xd//wL7/wf9/wL7/wL7/wG//wTf/wN//wJ//wK//wN//wH3/x33/wJ/7/8Gf/8F9/8H+/33/w/7/wl//wHv/wff/wL+/wb7/wF//wTf/wP1/wb+/wL3/wb+/wLf/wKf/wnv/wTv/wb9/wG//wF//wL+/wL9/wL7/xbf7/8Bv/8B9/8Fu/8Hv/8N3/8B/f8C/v7/Au//Bd//Ab//Auvn9/8C/f8Pffv/Bu//AX//Avv/Be//EH//EP7/CM77/wT+/wL3/wT33/8D3/8F+/8Uv/8D/f8M5v8ee/8M5/8R3/8E/fv/Cvz/Ff3/CH//Bt+//wHv/xL1/wz9/wN//wS//wLr/wp//wH7/wG//wl//wH3/wF//wXv/wH3/wZ//wL9/xD5/xb9/wfv/xn9/wz3/wH9/wz9/wH7+/8D93//Af3/Ar/+9/f/BXf/Av7/Avf/Ad//Cfv/C/3/Bf73/wv3/wTv/xXf/wK//wGfv/8C+/8K9/8E+/8B9/8G93//A9//Av7/De//Ae//AX//Ae/+/w31/wL9/f8D/v8F3/8N7/8F+3//Cv7/Cd//Avf/At//Evv/Aut//wLv/wHf/wP7/we//wbz/w37/f8C+/8R3/f/B+//EO//B+//A/v/Cf3/Bl//B9/73/8B7/r/D+//BX+//xHn/wv7/wH9/wW//wy//wG//w3z/wH3/wv9/wPf/wLv/xL+/wj8/w3f/wL3/wH9/wb7/wX++/8F9/8C/v8D+f8Ov/8Lf/8E7/8Cf/8D9/8Hf/8F7/8C7/8H/v7/Bv3/Au//E/7/Au//E+3+/wP3/wG//wHv/wF//wb9/wT9/wjf/wa/f3//A7//C/f/Dv3/DX/3/w3f/wPf/wN//wrv/w/3/wP7/wP3/xf9/wL+/wb9/wS//wW//wPv/wR//wzv/wHv/xf+/wnf/wH+/wV//wHv/wL9/wF//v8b7/8K3/8G7/8G9/8K/v8O33//Au//AX//CP7/CL//C7//CP5//wN//wP7/wvn/wv+/wP3/wL7/wrf/wTv/wj37/8K9v8G/v8c7/8K5/8C/f8H/v8Df/8Ge/8Dn/8Iv/8D+/8E3/8D3/8F3/8D+/8F9/8I/f8I7/8E7/8F3/8G9/8Jv/8D+/8F9/8O3/8D7/8E/v8C3/8Cv/8C79//FP5//w/v/wHu/wv3v/8d+/8I/v8Ef/8F7/8B+/8B/v8Bv/8K3/8I/f8E3/8Cf/8Bvz/+/wH9/wH7/wnv/wL+/wL5/wX7/wL9/wG//w3f/wX7/v8D9/8B+v8Gf3//CPf/AX//Cv3/Dr//C3//AX//BPv/Atf/A/v/BO/7/w/3/wf3/v8Jv/8G/f3/Ar3/Ff3/BP3/Bv3/BL//CP3/AX//BPv/Au//B7//Bf3+/w79/wv7/f8F9/v/BL//Ar/7/wi/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wa//wH3/xDf/wG//we//wP7/wjv/wT+/wK//wXf/wP+/wf9/wP9/wv+/v8F/v8Bv/8T+6//Ad//BNf/B+//Bf37/xb73/8R/v8O/f8S+/f/Be7f/wL7/wP3/wff9/8B/v8H93//Ef3/Avv/BPf/BOb/Cvf/Af7+f/8B9/3/A2/1/bv/BO//Bv7/Db//EL//Ad//BX//Bd//Ae//EL//D/v/CPf/BP3/Bef/BO//EN3/B/f/Dd//Bff/Ar//B77/Df3/AX//EPv/Ab//Afv/CX//Ae//A+/9/wL7/wnv/wXf/wd/9/8C3/8B8f8J/v8D7/8Cf/3+7/8E9/8N9/8G3/8ef/8Bf/8K/v8G9/8Fv/8E7/8O3/8F/v8D/f8Q338//wPv/wb7/wLf/wG//wT3/wjf/w2//wW+/xT+/wL79/8Q/v8Ef/8E3977/wG9/wFtPzol2B2QDw==",
		"+T1ABgAAAAD8DUEGAAAAAAIIAgi7lhzSAQ+//wj9/wF//wW//wT9/wH3/wX7/wf+f/8Hf/8Lf/8Rv/3/BP3f/we//wL3/wK/99//A/3/A+//Dv7/Cn//A/3/BPf/BK//Avv/BLv/BH//Av3/Be//Avf/Bt//Eu//Bff3/wr3/wG//wT3/wT9/wf3/wL1/wb9/wL8/wH3/xD73/8I+9//Ab//Cvf3/wT7/wb7/wtv/wL7/wL3/wP9/wLf/wTf/wPf/wr3/wO//wG+/wP+/wr9+/8D/f8D73f/EP7/CM//Af1/f/8D9/8H+/8G7/8Bv+//BP3/At//Ae3/A9//CO//Av3/BP3/A/f/Ab//B/7/Af3/A+//Cz//BP7/A9//AX//Aez3f/7/Av3/Avf/Avv/Dnv/CL//A3P/Bn7/BP3/Bd//Avf/An//B9//Ae//Bvf/CH//Aff/Af3/Bv7/Ae/9/wZv/we/+3//A/3/CP7/BX//BP7/Bu//Dvf/Avf9/wG//we//f8B7/8Cb799/f3/CP7/Bf7/BP3/BP7/Bs//CO//Ar//Be/f/wr7/v8Bf/8Q+/8C/v8Rv/8Y+/8H/f8Lf/8Rf3/v/wZ//wK//wTvvf8Dv/8B9/8Kv/8R+v8Fv/8D/v8I/v8C+/8H/f8C/v7/Cff+/wHv/wW/+/8B3v6//wLv/wbf/wj5/wTv/wvv/wH3/wK//wnf/wf7/wf9/wJ//wF7/wO//wb9/wL9/wO//wX3/wPH+f8B3/8C7/8F3/8G/f8Ov/8Bd+//Af77/wF//wX5f/8Cf+7/A7//A/v/BH//Av67/wN//wP7/wT3+99//wL3/wX+73//An9//wG//wPP/wOe/wH7/wl//wK//wTb/wP9/wHf/xS//wPf/wG//wf7/wb3/wK//xv3/wH3/wf7/xd//wLX/xHv/wHv/wG//f7/Au//Bt//AX//Af3/C/f/A9f/Av3/A/3/B+//Aq/f/wF//wL9/wjX/wLv/wbv/wb8/xH3/wTv+/8E/v8B+/8Df/7/At//Av7/Au//Ad//B+//FO//Cvv/C7//Aff/A3//Ab//BP3/CLf/Bb//C/f+f/8C+v8L+/8L/P8H9/8H/v8Iv/8E9/8Fnt//A7//Afv/Ab//B/f/Af7/Af7/B+//Avv/Bd//Be//Dff/D+//CP7/Bv3/C7//Ba//Bb//Av7/BX//Av3/AX//Af1//wT9/wS//wH3/wN//wLv/wX7/wm//wm//wR/9/8Hf/8Df/8b/v8F/f8K9/8N9/8G9/8C/f8C3/8Hv/8C7/8Bb/8F3/8B3/8I3/8L/v8Fr/8F7/8F/f8K/fv/CO//Bfv/GH/+/wTf/wf3/wHv/wVf/wz9/f8Lf3//Fvf/At//A/7/A+3/Au//Ar/+7/8B7/8E3/8M+f8Pf/8Ff/8Dv3//Aff/Bvf/At//Hvff9/8B5/8C9/8F3+/+/wH7/xDf/wP+/xL3/xTv/wr+/wP3/wXf/wT3v/8o/u//Af7/BO/3/f8Cv/8B+3//Avv/E3//C9/f/wO//wvvf/8F+/8F+/8I9/8Df/v/Af3/A/v9/wG//wT+/w7v/wV//wh//wK//wu/f/8C8/8Y9v8Gv7/f/wP9/wHz/wm//wa+/xH77/8P77//A/3v/wr7/wj3/wf+/wL9/wf+/wP3/wb+/wHf/wH3/w3f/f8F/v8O9/8J/f8Fv/3/De//Gvv/Bvv/B+//Dff9/wJ//wt//wLf/wN//wJ3/wJ//yd/f/8P/v8Bv/8Q+/8I/d//Ab9//wP7z/8D/f8Z77//Aff/A9v/BP3/Bvf/Ajb/AX//Bt//CP3f/wG//wHf/wvv/wTv/w7v+/8B9/8G9/8G/f8Cv/8B/v8B+/8Cv/8I7/8B/v8E98//Ct//C7//A/v/Bbv/Cvf/Aff/At8//wv+/wP+/wW//wH6/wTf/wP7/wnv/wP3/wl/f/8Hf/8E+9//Dvz+/wnv/xd//wL7/wf9/wL7/wPfv/8E3/8Df/8Cf/8Cv/8Df/8B9/8d9/8Cf+//Bn//Bff/B/v99/8Sf/8Gf/8E3/8E3/8C/v8G+/8Bf/8I9f7/Bf7/Avf/Bv7/At//Ap/v/wjv/wTv/wb9/wG//wF//wH+/wP9/wL7/wa//xDv/wG//wH3/wX7/we//w3f/wH9/wH9/v7/Au//Bd7/Ab//Avvn9/8C/f8Pffv/Bu//BPv/Be//EH//Dv7/Af7/CM77/wT+/wL3/wT33/8C79//Bfv/FL//A/3/DOb/Hnv/DOf/Ed//BP37/wr8/xX9/wh//we//wHv/xL1/wz9/wN//wS//wLr/wp//wO//wP9/wf3/wF//wH7/wPv/wH3/wZ//xPp/xD9/wX9/y73/w79/wH7+/8D93//Af3/Ar/+9/f/BX//Av7/Avf/AZ//Cfv/Ef73/wT7/wb3/wZ//xPf/wK//wGfv/8C+/8K9/8Df/8C9/8G93//A9//EO//A3//Ae/+/wH9/wv1/wL9/f8D/v8E/d//E/t//wf+/wL+/wnf/wLf/wLf/wy//wX7/wLrf/8C7/8B3/8D+f8Hv/8G8/8N+/3/Avv/Ed/3/wfv/xDv/wfv/wTv/wj9/wZf/wff+/8C7/r/D+//BX+//xHv/wv7/wH1/wW//wy+/wG//w3z/wH3/wv9/wPf/wLv/xL+/wj8/wXf/wff/wL3/wH9/wb7/wX++/8F9/8C/v8D+f8Ov/8Lf/8E7/8Cf/8D9/8Hf/8F7/8C7/8H/v7/Bv3/Au//E/7/Au//E+3+/wP3/wG//wHv/wF//wb9/wff/wXf/wa/f3//A7//Ab//Cff/Dv3/DX/3/w3f/wPf/w7v/w/3/wP7/xv9/wL+/wX+/wW//wW//wh//wzv/xn+/wO//wXf/wH+/wV//wHt/wL9/wF//v8Gv/8U7/8K3/8G7/8G9/8K/v8O33//Au//AX//CP7/Bf3/Ar//C7//CP5//wN//w/3/wv+/wP3/wy//wXv/wj37/8If/8B9v8J/f8h/f8C5/8J/f7/Cnv/A5//A7//BL//CN//A9//Cfv9/wT3/wj9/wjv/wTv/wXf/wb3/wm//wP7/wX3/w7f/wPv/wT+/wLf/wK//wLv3/8U/n//D+//Ae7/DL//BL//GPv/CP7/BH//Be//Afv/Afr/Ab//Cs//CP3/A7/f/wJ//wG/P/7/Af3/AfP/Ce//Av7/Avn/Bfv/Av3/Ab//Ce//A9//Bev+/wP3/wH6/wZ/f/8I9/8Bf/8K/f8I/f8Tf/8E+/8C1/8D+/8E7/v/F/f+/wm//wf8/wK9/xr9/wb9/wS//wj9/wb7/wLv/wT3/wK//wb+/w79/wv7/f8F9/v/BL//Ar//Cb9//wH9/wH3/wj7/wT7/wP+/v8I+9//CvX/Au//Av3/CPf/EN//Ab//C/v/CO//B7//Bd//C/3/D/7/Bv7/Ab//Avv/EPuv/wHf/wTX/wfv/wb7/wG//w73/wX73/8R/v8O/f8S+/f/Be7f/wL7/wP3/wff9/8B/v8G3/8Bf/8R/f8C+/3/A/f/BOb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/wjf/xff/wV//wXf/wHv/xC//w/7/wj3/wT9/wXn/wTv/xDd/wf3/w3f/wX3/wq+/w39/wF//xD7/wP7/wvv/wPv/f8C+f8J7/8F3/8Hf/f/At//AfH/Cf7/A+//An/9/v8F9/8N9/8G3/8ef/8Bf/8K/v8G9/8Fv/8E7/8H+/8G3/8F/v8D/f8Ev/8L3/8BP/8D7/8G+/8C3/8G9/7/Fb//Bb7/Dfv/Bv7/Avv3/xD+/wRv/wTf3v8C/f8B7T++d/gdkA8=",
		"yJtABgAAAAAEDkEGAAAAAAIIAgh7/wH3/wb3/wb9/wK//wH+/wH7/wvv/wH+/wT3z/8K3/8B+/8Jv/8D+/8Fu/8K9/8B9/8C3z//D/7/Bb//Af7/BN//A/v/Ce//A/f/CX9//wd//wT73/8O/Pb/Ce//F3//Avv/B/3/Ad/7/wS//wTf/wN//wJ//wK//wN//wH3/xF//wv3/wJ/7/8Gf/8F9/8H+/33/xl//wnf/wL+/wb7/wF//wH7/wb1/wb+/wL3/wZ+/wLf/wKf/wnv/wTv/wb9/wG//wF//wX9/wL7/xfv/wG//wH3/wX7/we//w3f/wH9/wL+/v8C7/8F3/8Bv/8C++f33/8B/f8H+/8Hffv/Bu//BPv/Be//EH//EP7/CM77/wT+/wL3/wT33/8D3/8F+/8Uv/8D/f8M5v8ee/8M5/8R3/8E/fv/A+//Bvz/Ff3/CH//B7//Ae//Bu//C/X/CN//A/3/A3//BL9//wHr/we//wJ//wO//wv3/wF//wXv/wH3/wZ//xDv/wL5/wb9/w/9/x/3/w73/wjf/wX9/wH7+/8D93//Af3/Ar/+9/f/BX//Av7/Avf/Ad//Cfv/Ef73/wb9/wT3/wXv/xTf/wK//wGfv/8C+/8K9/8G9/8G93//A9//Bf3/Cu//A3//Ae/+/w31/wL9/f8D/v8F3/8T+3/+/wn+/wnf/wXf/xL7/wLrf/8C7/8B3/8D+/8Hv/8G8/8N+/3/Avv/Ed/3/wfv/xDv/wfv/w39/wZf/wff+/8C7/r/D+//BX+//xHv/wv7/wH9/wW//wy//wG//wO//wnz/wH3/wv9/wPf/wLv/xL+/wj8/w3f/wL3/wH9/wb7/wXu+/8F9/8C/vv/Avn/CX//BL//C3//BO//An//A/f/B3//Be//Au//B/7+/wb9/wLv/xP+/wLr/xPt/v8D9/8Bv/8B7/8Bf/8G/f8N3/8Gv39//wO//wm//wH3/wf9/wb9/wX+/wd/9/8N3/8D3/8Kf/8D7/8K/v8E9/8D+/8b/f8C/v8G+/8Ev/8Fv/8If/8M7/8Z/v8J3/8B/v8Ff/8B7/8C/f8Bf/7/Ae//CP3/EO//Ct//A/v/Au//Bvf/Cv7/Dt9//wLv/wF//wj+3/8Hv/8E3/8Gv/8I/n//A3//D/f/C/7/A/f/Eu//CPfv/wT3/wX2/wjf/yXn/wLf/wf+/wp7/wOf/wZ//wG//wjf/wG//wHf/wn7/v8E9/8I/f8I7/8C+/8B7/8F3/8G9/8H3/8Bv/8D+/8F9/8O3/8D7/8E7v8C3/8Cv/8C79//DP3/B/5//w/v/wHu/wy//xS//wj7/wj+/wR+/wXv/wH7/wH+/wG//wfv/wLf/wj9/wTf/wJ//wG/P/7/Af3/Afv/Ce//Av7/Avn/Bfv/Av3/Ab//Cv7/At//Bfv+/wP3/wH6/wZ/f/8I9/8Bf/8I/v8B/f8cf/8E+/8C1/8D+/8E7/v/D/7/B/f+/wm//wf9/wK9/wd//wp//wf9/wb9/wS//wj9/wb7/wLv/we//wb+/w79/wv7/f8F97v/BL//Ar//Cb9//wH9/wH3/wH3/wb7/wPf+/8D/v7/CPvf/wl/9f8C7/8C/f8I9/8Q3/8Bv/8J/f8B+/8I7/8Hv/8F3/8J3/8B/f8P/v8G/v8Bv/8Iv/8K+6//Ad//BNf/B+//Bvv/EPf/BP773/8O7/8C/v8O/f8S+/f/Be7f/wL7/wH+/wH3/wff9/8B/v8If/8R/f8C+/8E9/8Dv+b/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/w79/xHf/wV//wXf/wHv/wT9/wu//wX+/wn7/wj3/wT9/wXn/wTv/wvf/wTd/wf3/w3f/wK//wL3/wq+/w39/wF//xD7/wP7/wvv/wPv/f8C+/8J7/8F3/8F9/8Bf/f/At//AfH/Cf7/A+//An/9/v8F9/8K+/8C9/8G3+//HX//AX//Cv7/Bvf/Bb//BO//Dt//Bf7/A/3/EN//AT//A+//Bvv/At//Bvf/FN//Ab//Bb7/FP7/Avv3/wm//wb+/wR//wTf3/8C/f8CP75/+B2YCI//BPv/Ag==",
		"zZtABgAAAAAEDkEGAAAAAAIIAgj7/wH3/wb3/wa9/wK//wH+/wH7/wb7/wTv/wH+/wT3z/8K3/8Lv/8D+/8Bf/8Du/8K9/8B9/8C3z//D/7/Bb//Af7/BN//A/v/Ce//A/f/CX9/7/8Gf/8E+9//Dvz+/wnv/xd//wL7/wf9/wL7/wS//wTf/wN/7/8Bf/8Cv/8Df/8B9/8B/f8b9/8B/X/v/wZ//wX3/wf7/ff/GX//Bfv/A9//Av7/Bvv/AX/f/wf1/wb+/wL3/wb+/wLf/wKf/wnv/wK//wHv/wb9/wG+/wF//wX9/wL7/xfv/wG//wH3/wX7/we//w3f/wH9/wL+/v8C7/8F3/8Bv/8C++f3/wL9/w99+/8G7/8E+/8F7/8K/v8Ff/8Q/v8ITvv/BP7/Avf/BPff/wPf/wX7/xS//wP9/wu/5v8ee/8M5/8R3/8E/fv/Cvz/Ae//E/3/CH//B7//Ae//EvX/DP3/A3//BL/+/wHr/wp/9/8Cv/8L9/8Bf/8F7/8B9/8B+/8Ef/8T+f8E/v8Q9/3/Lvf/Dv3/Afv7/wP3f/8B/f8Cv/739/8Ff/8C/v8C9/8B3/8J+/8R/vf/C/f/Gt//Ar//AZ+//wH++/8K9/8G9/8G93//A9//Afv/Ds//A3//Ae/+/wf+/wX1/wL9/f8D/v8F3/8T+3//Cv7/Cd//Bd//Cf7/CPv/Aut//wLv/wHf/wP7/we+/wbz/w37/f8C+/8R3/f/B+//EO//B+//Df3/Bl//B9/7/wLv+v8I7/8G7/8E33+//wj+/wjv/wS//wb7/wH9/wW//wy//wG//wXf/wfz/wH3/wH9/wn9/wPf/wLv/wbf/wv+/wj8/w3f/wL3/wH9/wTv/wH7/wX++/8F9/8C/v8B7/8B+f8N77//Cvd//wTv/wJ//wL39/8Hf/8F7/8C7/8H/v7/Bv3/Au//E/7/Au//E+3+/wP3/wG//wHv/wF//wb9/w3f/wa/f3//A7//C/f/Dv3/DX/3/w3f/wPf/wXv/wjv/w/3/wP7/w+//wv9+/8B/v8Lv/8Fv/8If/8M7/8C/f8W/v8G+/8C3/8B/v8Ff/8B7/8Bf/3/AX/+/xvv/wrf/wbv/wb3/wr+/w7ff/8C7/8Bf/8I/v8Iv7//Cr//CP5//wN//wq//wT3/wf7/wP+/wP3/wZ//wvv/wG//wb37/8E/f8F9v8Yv/8V5/8K/v8Ke/8Dn/8Iv/8I3/8D3/8J+/8F9/8I/f8I7/8C/f8B7/8F3/8G9/8Jv/8D+/8F9/8If/8F3/8D7/8E/v8B/d//Ar//Au/f/wvv/wj+f/8D/v8L7/8B7t//C7//HP37/wj+/wR//v8E7/8B+/8B/v8Bv/8K3/8H+/3/BN//An//Ab8//v8B/f8B+/8J7/8C+v8C+f8F+/8C/f8Bv/3/C/ff/wX7/v8D8/8B+v8Gf3//CPf/AX//CX/9/w33/w5//wT7/wLX/wP7v/8D7/v/F/f+/wm//wf9/wK9/wHf/xD7/wL9/wT9/wX9/f8Ev/8I/f8C7/8D+/8C7/8Hv/8G/v8O/f8L+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wj3/xDf/wG//wv7/wjv/we//wXf/wv9/w/+/wb+/wG//wH3/w39/wP7r/8B3/8E1/8H7/8G2/8Q9/8F+9//Ef7/B/f/Bv3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/CH//Ef3/Avv/BPf/Ar//Aeb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/wy//xPf/wV//wXf/wHv/xC//w/7/wj3/wT9/wXn/wTv/wb9/wnd/wf3+/8M3/8D7/8B9/8Kvv8C/f8K/f8Bf/8I/v8H+/8D+/8L7/8D7/3/Avv/Ce//Bd//Af3/BX/3/wLf/wHx/wb3/wL+/wPv/wJ//f7/Bff/Dff/Bt//Hn//AX//Cv7/Bvf/Bb3/BO//C9//At//Bf7/A/3/EN//AT//A+//Bvv/At//Bvf/C9//Cr//Bb7/FP7/Avv3/xD+/wR//wTf3/8C/f8CP75/+B+YCI//BPv/Ag==",
		"/Z1ABgAAAAAGDkEGAAAAAAIIAgjKAAECP/8C+/8Fuf8K9/8B9/8C3z//D/77/wS//wH+/wTf/wP7/wnv/wP3/f8If3//B3//BPvf/w78/v8D/f8F7/8V/v8Bf/8C+/8H/e//Afv/BL//BN//A3//An//Ar//A3+/9/8Z/f8D9/8Cf+//Bn//Bff/B/v99/8M3/8Mf/8J3/8C/v8G+/8Bf/8I9f8G/v8C9/8G/v8C3/8Cn/8I/e//BO//Bv3/Ab//AX//Bf3/Avv/BL//Db//BO//Ab//Aff/Ae//A/v/B7//Dd//Af3/Av7+/wH+7/8F3/8Bv/8B/vvn9/8C/f8G7/8Iffv/Bu//BPv/Be//EH//EP7/Avf/Bc77/wT+/wL3/wT33/8D3/8F+/8G9/8Nv/8D/f8M5v8ee/8M5/8R3/8E/fv/Cvz/Ff3/CH//B7//Ae//EvX/DP3/A3//BL//Auv/Cn//A7//BPf/Bvf/AX//Be//Aff/Bn//Dv3/BPn/Et//A/3/Lvf/De/9/wH7+/8D93//Af3/Ar/+9/f/BX//Av7/Avf/Ad//Cfv/CP7/CP73/wG//wn3/f8Z3/8Cv/8Bn7//Avv/Cvf/Bvf/Bvd//wPf/wPv/wzv/wN//wHv/v8N9f8C/f3/A/7/Bd//Cf3/Cft//wr+/wT3/wTf/wXf/xL7/wLrf/8C7/8B3/8D+/8Hv/8G8/8N+/3/Avv/BO//DN/3/wfv/xDv/wfv/w39/wZf/wff+/8C7/r/D+//BX+//wr9/wbv/wv7/wH9/wW//wb3/wW//wG//w3z/wH3/wv9/wPf/wLv/xL+/wj8/wJ//wrf/wL37/3/Bvv/Bf77/wX3/wL+/wP5/w6//wPv/wd//wTv/wJ//wP3/wd//wXv/wLv/wHf/wX+/v8G/f8C7/8T/v8C7/8T7f7/A/f/Ab//Ae//AX//Bv3/Dd//Br9/f/8Dv/8I/f8C9/8O/f8Nf/f/Dd//A9//Du//D/e//wL7/xv9/wL+/wu//wW//wh//wzv/xn+/wnf/wH+/v8Ef/8B7/8C/f8Bf/7/G+//Avf/B9//Bu//Bvf/Cv7/Dt9//wLv/wF//wH9/wb+/wi//wu/f/8H/n//Au9//w/3/wv+/wP3/wr7/wfv/wj37/8K9v8u5/8Bf/8I/v8Ke/8Dn/8Iv/8I3/8D3/8H/f8B+/8F9/8G+/8B/f8I7/8E7/7/BN//Bvf/Cb//A/v/Bff/Dt//A+//BP7/At//Ar//Au/f/xD3/wP+f/8P7/8B7v8D3/8Iv/8Zf/8D+/8I/v8Ef/8F7/8B+/8B/v8Bv/8K3/8F/v8C/f8E3/8Cf/8Bvz/+/wH9/wH7/wnv/wL+/wL5/wX7/wL9/wG//w3f/wX7+v8D9/8B+v8Gf33/CPf/AX//Ab//CP3/A3//GH//BPv/Atf/A/v/BO/7/xf3/v8Jv/8H/f8Cvf8V/f8E/f8E/f8B/f8Ev/8I/f8G+/8C7/8Hv/8G/v8O/f8L+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wj3/xDf/wG//wv7/wjv/wH9/wW//wXf/wnv/wH9/w/+/wb+/wG//xP7r/8B3/8E1/8H7/8Ge/8Mf/8D9/8F+9//Ef7/Dv3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/CH//Bv3/Cv3/Af37/wT3/wTm/wr3/wH+/wF//wH3/f8Db/X1u/8E7/8G/v8Ov/8R3/8Ff/8F3/8B7/8Gf/8Jv/v/Df37/wj3/wT9/wXn/wTv/xDd/wf3/w3P/wX3/wb+/wO+/w39/wF//xDz/wP7/wvv+/8C7/3/Avv/Ce//Bd//B3+3/wLf/wHx/wn+/wPv/wJ//f7/Bff/C/v/Aff/Bt//Hn//AX//Cv7/Bvf/Bb//BO//Dt//Bf7/A/3/EN//AT//A+//A/f/Avv/At//Bvf/Bn//Ct//BL//Bb7v/xP+/wL79/8Q/v8Ef/8E39//Av3/An++f/gfmgyr/wT7/wI=",
		"x8BABgAAAAAJDkEGAAAAAAIIAgj/BN//BN//Bd/v/xH7/wLrf/8C7/8B3/8D+/8C3/8Ev/8G8/8N+/3/Avv/Ed/3/wfv/wv+/wTv/wfv/w39/wZf3/8G3/v/Au/6/w/P/wV/v/8R7/8L+/8B/f8Fv/8Mv/8Bv/8N8/8B9/8L/f8D3/8C7/8S/v8I/P8N3/8C9/8B/f8G+/8F/vv/Bff/Av7/A/n/Dr//C3//Af7/Au//An//A/f/B3//Be//Au//B/7+/wb9/wLv/xP+/wLv/xPt/v8D9/8Bv/8B7/8Bf/8G/f8N3/8Gv39//wO//wv3/w79/wj7/wR/9/8N3/8D3/8O7/8P9/8D+/8b/f8C/v8Lv/8Fv/8F+/8Cf/8Jv/8C7/8H9/8R/v8J3/8B/v8Ff/8B7/8C/d9//v8I+/8K+/8H7/8K37//Be//Bvf/Cv7/Dt9//wLv/wF//wj+/wi//wu//wj+f3//An//DL//Avf/C/7/A/f/Eu//CPfv/wr2/y7n/wr+/wp7/wOf/wi//wjf/wPf/wn7/wX3/wF//wb9/wT9/wPv/wTv/wXb/wb3/wm//wP7/wX3/w7f/wPv/wT+/wLf/wK//wLv3/8Ff/8O+n//D+//Ae7/Cf3/Ar//C3//Cvf/Bvv/A9//BP7/BH//Be//Afv/Af7/Ab//Ct//CP3/BN//An//Ab8//v8B/f8B+/8J7/8C/v8C+f8F+/8C/f8Bv/8E+/8I3/8F+/7/A/f/Afr/Bn9//wj3f3//Cv3/Ab//E/v/Bn//BPt//wHX/wP7/wP+7/v/C9//C/f+/wm//wf9/wK9/wH7/wz3/wb9/wT9/wb9/wS//wj9/wb7/wLv/we//wb+/w79/wf3/wP7/f8F9/v/BL//Ar//Cb9//wH9/wH3/wj7/wT7/wP+/v8I+9//CvX/Au//Av3/A/7/BPf/EN//Ab//C/v/BL//A+//Dd//C/3/D/7/Bv7/Ab//E/uv/wHf/wTX/wfv/wb7/wrf/wX3/wX73/8R/v8O/f8Nv/8E+/f/Be7f/wL7/wP3/wL+/wTf9/8B/v8D/f8Ef/8R/ff/Afv/BPf/BOb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/yDf/wV/3/8E3/8B7/8Qv/8I+/8G+/8F3/8C9/8E/f8F5/8E7/8Q3f8H9/8N3/8Fd/8Kvv8N/f8Bd/8K7/8F+/8D+/8L7/8B/f8B7/3/Avv/B3//Ae//Bd//B3/3/wLf/wHx/wK//wb+/wPv/wJf/f7/Bff/Dff/Bt//Hn//AX//Cv7/Bd/3/wW//wTv/w7f/wX+/wP9/xDf/wE//wPr/wb7/wLf/wb3/wnv/wy//wW+/xT+/wL79/8Q/v8Ef/8E39//Av3/An++f/gfmgyqf/8D+/8C",
		"zZtABgAAAAARDkEGAAAAAAIIAgj7/wH3/wb3/wb9/wK//wH+/wH7/wvv/wH+/wT3j/8K3/8C/v8Iv/8D+/8Fu/8K9/8B9/8C3z//B/3/B/7/Bb//Af7/BN//A/v/Ce//A/f/CX9/v/8Gf/8Ee9//Dvz+/wnv/xd//wL7/wf9/wL7/wL7/wG//wTf/wN//wJ//wK//wN//wH3/x33/wJ/7/8Gf/8F9/8H+/33/w/7/wl//wHv/wff/wL+/wb7/wF//wTf/wP1/wb+/wL3/wb+/wLf/wKf/wnv/wTv/wb9/wG//wF//wL+/wL9/wL7/xbf7/8Bv/8B9/8Fu/8Hv/8N3/8B/f8C/v7/Au//Bd//Ab//Auvn9/8C/f8Pffv/Bu//AX//Avv/Be//EH//EP7/CM77/wT+/wL3/wT33/8D3/8F+/8Uv/8D/f8M5v8ee/8M5/8R3/8E/fv/Cvz/Ff3/CH//Bt+//wHv/xL1/wz9/wN//wS//wLr/wp//wH7/wG//wl//wH3/wF//wXv/wH3/wZ//wL9/xD5/xb9/wfv/xn9/wz3/wH9/wz9/wH7+/8D93//Af3/Ar/+9/f/BXf/Av7/Avf/Ad//Cfv/C/3/Bf73/wv3/wTv/xXf/wK//wGfv/8C+/8K9/8E+/8B9/8G93//A9//Av7/De//Ae//AX//Ae/+/w31/wL9/f8D/v8F3/8N7/8F+3//Cv7/Cd//Avf/At//Evv/Aut//wLv/wHf/wP7/we//wbz/w37/f8C+/8R3/f/B+//EO//B+//A/v/Cf3/Bl//B9/73/8B7/r/D+//BX+//xHn/wv7/wH9/wW//wy//wG//w3z/wH3/wv9/wPf/wLv/xL+/wj8/w3f/wL3/wH9/wb7/wX++/8F9/8C/v8D+f8Ov/8Lf/8E7/8Cf/8D9/8Hf/8F7/8C7/8H/v7/Bv3/Au//E/7/Au//E+3+/wP3/wG//wHv/wF//wb9/wT9/wjf/wa/f3//A7//C/f/Dv3/DX/3/w3f/wPf/wN//wrv/w/3/wP7/wP3/xf9/wL+/wb9/wS//wW//wPv/wR//wzv/wHv/xf+/wnf/wH+/wV//wHv/wL9/wF//v8b7/8K3/8G7/8G9/8K/v8O33//Au//AX//CP7/CL//C7//CP5//wN//wP7/wvn/wv+/wP3/wL7/wrf/wTv/wj37/8K9v8G/v8c7/8K5/8C/f8H/v8Df/8Ge/8Dn/8Iv/8D+/8E3/8D3/8F3/8D+/8F9/8I/f8I7/8E7/8F3/8G9/8Jv/8D+/8F9/8O3/8D7/8E/v8C3/8Cv/8C79//FP5//w/v/wHu/wv3v/8d+/8I/v8Ef/8F7/8B+/8B/v8Bv/8K3/8I/f8E3/8Cf/8Bvz/+/wH9/wH7/wnv/wL+/wL5/wX7/wL9/wG//w3f/wX7/v8D9/8B+v8Gf3//CPf/AX//Cv3/Dr//C3//AX//BPv/Atf/A/v/BO/7/w/3/wf3/v8Jv/8G/f3/Ar3/Ff3/BP3/Bv3/BL//CP3/AX//BPv/Au//B7//Bf3+/w79/wv7/f8F9/v/BL//Ar/7/wi/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wa//wH3/xDf/wG//we//wP7/wjv/wT+/wK//wXf/wP+/wf9/wP9/wv+/v8F/v8Bv/8T+6//Ad//BNf/B+//Bf37/xb73/8R/v8O/f8S+/f/Be7f/wL7/wP3/wff9/8B/v8H93//Ef3/Avv/BPf/BOb/Cvf/Af7+f/8B9/3/A2/1/bv/BO//Bv7/Db//EL//Ad//BX//Bd//Ae//EL//D/v/CPf/BP3/Bef/BO//EN3/B/f/Dd//Bff/Ar//B77/Df3/AX//EPv/Ab//Afv/CX//Ae//A+/9/wL7/wnv/wXf/wd/9/8C3/8B8f8J/v8D7/8Cf/3+7/8E9/8N9/8G3/8ef/8Bf/8K/v8G9/8Fv/8E7/8O3/8F/v8D/f8Q338//wPv/wb7/wLf/wG//wT3/wjf/w2//wW+/xT+/wL79/8Q/v8Ef/8E39/7/wH9/wJ/vn/4H5pOqkB//wL7/wI=",
		"r9lABgAAAAATDkEGAAAAAAIIAggF/v8O33//Au//AX//CP7/CJ//C7//A9//BP5//wN//w/3/wv+/wP3/xLv/wj37/8K9v8u5/8K/v8Ke/8Dn/8Iv/8I3/8D3/8J+/8F9/8I/f8I7/8E7/8F3/8G9/8Jv/8De/8F9/8O3/8D7/8E/v8C3/8Cv/8C79//FP5//wK//wzv/wHu/wy//x37/wj+/wR//wP9/wHv/wH7/wH+/wG//wrf/wj9/wTfv/8Bf/8Bvz/+/wH9/wH7/wX7/wPv/wL+/wL5/wX7/wL9/wG//w3f/wX7/v8D9/8B+v8Gf3//CPf/AX//Cv3/GH//A3//BPv/Atf/A/v/BO/7/xP7/wP3/v8Jv/8H/f8Cvf8V/f8E/f8E+/8B/f8Ev/8I/f8G+/8C7/8Hv/8G/v8O/f8L+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D9v7/CPvf/wr1/wLv/wL9/wj3/xDf/wG//wH9/wn7/wZ//wHv/we//wPv/wHf/wv9/w/+/wb+/wG//xP7r/8B3/8E1/8H7/8G+/8Q9/8F+9//Ef7/Dv3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/CH//Cff/B/3/Avv/BPf/BOb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/w/+/xDf/wV//wXf/wHv/wrv/wW//w/7/wW//wL3/wT9/wXn7/8D7/8Q3f8H9/8Dv/8J3/8F9/8Kvv8N/f8Bf/8Q+/8D+/8L7/8D7/3/Avv/Ce//Bd//B3/3/wLf/wHx/wn+/wPv3/8Bf/3+/wX3/wT7/wj3/wbf/x5//wF//wr+/wb3/wW//wTv/w7f/wX+/wP9/xDf/wE//wH3/wHv/wb7/wLf/wb3/xa//wW+/xT+/wL79/8J+/8G/v8Ef/8E39//Av3/An/+f/ofml6qQF//Avv/Ag==",
		"+nlABgAAAAATDkEGAAAAAAIIAggB+rINAAFKAAHf/wjf/wV//wX+/wWv/wXv/wX9/wfv/wL9+/8I7/8T3/8Kf/7/BN//B/f/Ae//Bd//DP39/wt/f/8W9/8C3/8D/v8D7f8C7/8Cv/7v/wHv/xH5/xV//wO/f/8I9/v/Ad//Hvff9/8B5/8C9/8F3+/+/wH7/xDf/wP+/yfv/wr+/wP3/wXf/wW//yj+7/8B/v8E7/f9/wK//wH7f/8i39//A7//C+9//wv7/wj3/wN/+/8B/f8D+/8Cu/8E/v8O7/8Of/8Cv/8Lv3//AvP/GP7/Br+/3/8D/f8B9/8Jv/8Gvv8R++//D++//wP9/wv7/wj3/wf+/wPf/wb+/wP3/wb+/wHf/wH3/w3f/f8F/v8Gf/8H9/8Pv/3/B/3/Be//Gvv/Bvv/B+//Dff9/wJ//wt//wLf/wZ3/wJ//wv7/w/7/wt/f/8Hv/8H/v8Bv/8Z/d//Ab9//wP7z/8D/f8D/v8Sf/8C77//Aff/A9v/BP3/Bvf/Ajb/AX//Bt//CP3/Ar//Ad//C+//E+/7/wH3/wb3/wb9/wK//wH++/v/C+//Af7/BPfP/wrf/wu//wP7/wW7/wr3/wH3/wLfP/8P/v8Ft/8B/v8E3/8D+/8J7/8D9/8Jf3//B3//BPvf/v8N/P79/wjv/wW//xF//wL7/wf9/wL7/wS//wTf3/8Cf/8Cf/8Cv/8Df/8B9/8M/f8Q9/8Cf+//DPf/B/v99/8Zf/8Dv/8F3/8C/v8G+/8Bf/8I9f8G/v8C9/8E9/8B/v8C3/8Cn/8J7/8E7/8G/f8Bv/8Bf/8F/f8C+/8X7/8Bv/8B9/8F+/8Hv/8N3/8B3f8C/v7/Au//Bd//Ab//Avvn9/8C/f8Pffv/Bu//BPv/Be//EH//EP7/CM77/wT+/wL3/wT33/8D3/8F+/8Uv+//Av3/DOb/HP7/AXv/DOf/Fv37/wr8/wXv/w/9/wh//we//wHv/xL1/wz9/wN//wS//wLr/wp//wO//wv3/wF//wXv/wH3/wX3f/8T+f8W/f8N9/8Y7/8H9/8O/f8B+/v/A/d//wH9/wK//vf3/wV//wL+/wL3/wHf/wn7/wn+/wf+9/8L9/8a3/8Cv/8Bn7//Avv/Cvf/Bvf/Bvd//wPf/xDv/wHf/wF//wHv/v8N9f8C/f3/A/7/Bd//E/t//f8J/v8H+/8B3/8F3/8S+/8C63//Au//Ad//A/v3/wa//wbz/wH7/wv7/f8C+/8R3/f/B+//EO//B+//Df3/Bl//B9/7/wLv+v8P7/8Ff7//Ee//BX//Bfv/Af3/Bb//DL//Ab//C7//AfP/Aff/C/3/A9//Au//Ev7/CPz/Dd//Avf/Af3/Bvv/Bf77/wX3/wL+/wP5/w6//wX7/wV//wTv/wJ//wP3/wd//wXv/wLv/wf+/v8G/f8C7/8Kv/8I/v8C7/8F7/8N7f7/A/f/Ab//Ae//AX//Bv3/DP7f/wa/f3//Av6//wv3/wXv/wj9/w1/9/8N3/8B9/8B3/8M/f8B7/8P9/8Cv/v/Ed//Cf3/Av7/C7//Bb//CH//DO//Gf7/Cd//Af7/BX//Ae//Av3/AX/+/xvv/wX3/wTf/wbv/wb3/wr+/wbf/wfff/8C7/8Bf/8I/v8Iv3//Cr//CP5/3/8Cf/8P9/8L/v8D9/8S7/8I9+//Cvb/Kn//A+f/CX/+/wp7/wOf/wi//wjf/wPf/wj3+/8F9/8I/f8I7/8E7/8F3/8G9/8Jv/8D+/8F9/8O3/8D7/8E/v8C3/8Cv/8C79//FP5//w/v/wHu/wHf/wq//x37/wj+/wR//wXv/wH7/wH+/wG//wrf/wP7/wT9/wTf/wJ//wG/P/7/Af3/Afv/Ce//Av7/Avn/Bfv/Av3/Ab//Dd//Bfv+/wP3/wH6/wZ/f/8F/f8C9/8Bf/8K/f8Tv/8If/8E+/8C1/8D+/8E7/v/F/f+/wP3/wW//wf9/wK8/xX9/wT9/wb9/wS//wj9/wb7/wLv/we//wb+/w79/wv7/f8F9/v/BL//Ar//Cb9//wH9/wH3/wj7/wT7/wP+/v8I+9//CvX/Au//Av3/Bfv/Avf/EN//Ab//C/v/CO//Dd//C/3/D/7/Bv7/Ab//Efv/Afuv/wHf/wTX/wfv/wb7/wjf/w373/8R9v8Off8S+/f/Be7f/wL7/wP3/wff9/8B/v8If/8R9f8C+/8E9/8E5v8K9/8B/v8Bf/8B9/3/A2/1/bv/BO//A/v/Av7/IN/9/wR//wXf/wHv/wX+/wq//w/7/wJ//wX3/wT9/wXn/wTv/wTv/wvd/wf3/w3f/wX3/wq+/wL7/wr9/wF//wN//wz73/8C+/8L7/8D7/3/Avv/Ce//Bd//B2/3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/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL5/wj3/xDf/wG//wv7/wjv/we//wXf/wv9/w/+/wb+/wG//wn3/wn7r/8B3/8E1/8H7/8G+/8J3/8G9/8F+9//Ef7/Dv3/Evv3/wXu3/8C+/8D9/8D/f8D3/f/Af7/CH//Ef3/Avv/BPf/BOb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/f8f3/8D9/8Bf/8F3/8B7/8Ef/8Lv/8P+/8I5/8E/f8F5/8E7/8Q3f8H9/8Nn/8F9/8F/f8Evv8H9/8F/f8Bf/8Q+/8D+/8L7/8D7/3/Avv/Ce//Bd//B3/3/wLf/wHx/wn+/wPv/wJ//f6//wT3/w33/wbf/x5//wF//wr+/wb3/wW//wTv/w7f/wX+/wP9/wH7/wvf/wLf/wE//wPv/wb7/wLf/wb3/xa//wW+/wj3/wv+/wL79/8Vf/8E39//BX//AX/6H55eqmBYf/8B+/8C",
		"zZtABgAAAAAZDkEGAAAAAAIIAgj7/wH3/wb3/wa9/wK//wH+/wH7/wb7/wTv/wH+/wT3z/8K3/8Lv/8D+/8Bf/8Du/8K9/8B9/8C3z//D/7/Bb//Af7/BN//A/v/Ce//A/f/CX9/7/8Gf/8E+9//Dvz+/wnv/xd//wL7/wf9/wL7/wS//wTf/wN/7/8Bf/8Cv/8Df/8B9/8B/f8b9/8B/X/v/wZ//wX3/wf7/ff/GX//Bfv/A9//Av7/Bvv/AX/f/wf1/wb+/wL3/wb+/wLf/wKf/wnv/wK//wHv/wb9/wG+/wF//wX9/wL7/xfv/wG//wH3/wX7/we//w3f/wH9/wL+/v8C7/8F3/8Bv/8C++f3/wL9/w99+/8G7/8E+/8F7/8K/v8Ff/8Q/v8ITvv/BP7/Avf/BPff/wPf/wX7/xS//wP9/wu/5v8ee/8M5/8R3/8E/fv/Cvz/Ae//E/3/CH//B7//Ae//EvX/DP3/A3//BL/+/wHr/wp/9/8Cv/8L9/8Bf/8F7/8B9/8B+/8Ef/8T+f8E/v8Q9/3/Lvf/Dv3/Afv7/wP3f/8B/f8Cv/739/8Ff/8C/v8C9/8B3/8J+/8R/vf/C/f/Gt//Ar//AZ+//wH++/8K9/8G9/8G93//A9//Afv/Ds//A3//Ae/+/wf+/wX1/wL9/f8D/v8F3/8T+3//Cv7/Cd//Bd//Cf7/CPv/Aut//wLv/wHf/wP7/we+/wbz/w37/f8C+/8R3/f/B+//EO//B+//Df3/Bl//B9/7/wLv+v8I7/8G7/8E33+//wj+/wjv/wS//wb7/wH9/wW//wy//wG//wXf/wfz/wH3/wH9/wn9/wPf/wLv/wbf/wv+/wj8/w3f/wL3/wH9/wTv/wH7/wX++/8F9/8C/v8B7/8B+f8N77//Cvd//wTv/wJ//wL39/8Hf/8F7/8C7/8H/v7/Bv3/Au//E/7/Au//E+3+/wP3/wG//wHv/wF//wb9/w3f/wa/f3//A7//C/f/Dv3/DX/3/w3f/wPf/wXv/wjv/w/3/wP7/w+//wv9+/8B/v8Lv/8Fv/8If/8M7/8C/f8W/v8G+/8C3/8B/v8Ff/8B7/8Bf/3/AX/+/xvv/wrf/wbv/wb3/wr+/w7ff/8C7/8Bf/8I/v8Iv7//Cr//CP5//wN//wq//wT3/wf7/wP+/wP3/wZ//wvv/wG//wb37/8E/f8F9v8Yv/8V5/8K/v8Ke/8Dn/8Iv/8I3/8D3/8J+/8F9/8I/f8I7/8C/f8B7/8F3/8G9/8Jv/8D+/8F9/8If/8F3/8D7/8E/v8B/d//Ar//Au/f/wvv/wj+f/8D/v8L7/8B7t//C7//HP37/wj+/wR//v8E7/8B+/8B/v8Bv/8K3/8H+/3/BN//An//Ab8//v8B/f8B+/8J7/8C+v8C+f8F+/8C/f8Bv/3/C/ff/wX7/v8D8/8B+v8Gf3//CPf/AX//CX/9/w33/w5//wT7/wLX/wP7v/8D7/v/F/f+/wm//wf9/wK9/wHf/xD7/wL9/wT9/wX9/f8Ev/8I/f8C7/8D+/8C7/8Hv/8G/v8O/f8L+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wj3/xDf/wG//wv7/wjv/we//wXf/wv9/w/+/wb+/wG//wH3/w39/wP7r/8B3/8E1/8H7/8G2/8Q9/8F+9//Ef7/B/f/Bv3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/CH//Ef3/Avv/BPf/Ar//Aeb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/wy//xPf/wV//wXf/wHv/xC//w/7/wj3/wT9/wXn/wTv/wb9/wnd/wf3+/8M3/8D7/8B9/8Kvv8C/f8K/f8Bf/8I/v8H+/8D+/8L7/8D7/3/Avv/Ce//Bd//Af3/BX/3/wLf/wHx/wb3/wL+/wPv/wJ//f7/Bff/Dff/Bt//Hn//AX//Cv7/Bvf/Bb3/BO//C9//At//Bf7/A/3/EN//AT//A+//Bvv/At//Bvf/C9//Cr//Bb7/FP7/Avv3/xrf3/8Ff/8Bf/ofnl6qYFh//wH7/wI=",
		"+T1ABgAAAAAZDkEGAAAAAAIIAgi7lhzSAQ+//wj9/wF//wW//wT9/wH3/wX7/wf+f/8Hf/8Lf/8Rv/3/BP3f/we//wL3/wK/99//A/3/A+//Dv7/Cn//A/3/BPf/BK//Avv/BLv/BH//Av3/Be//Avf/Bt//Eu//Bff3/wr3/wG//wT3/wT9/wf3/wL1/wb9/wL8/wH3/xD73/8I+9//Ab//Cvf3/wT7/wb7/wtv/wL7/wL3/wP9/wLf/wTf/wPf/wr3/wO//wG+/wP+/wr9+/8D/f8D73f/EP7/CM//Af1/f/8D9/8H+/8G7/8Bv+//BP3/At//Ae3/A9//CO//Av3/BP3/A/f/Ab//B/7/Af3/A+//Cz//BP7/A9//AX//Aez3f/7/Av3/Avf/Avv/Dnv/CL//A3P/Bn7/BP3/Bd//Avf/An//B9//Ae//Bvf/CH//Aff/Af3/Bv7/Ae/9/wZv/we/+3//A/3/CP7/BX//BP7/Bu//Dvf/Avf9/wG//we//f8B7/8Cb799/f3/CP7/Bf7/BP3/BP7/Bs//CO//Ar//Be/f/wr7/v8Bf/8Q+/8C/v8Rv/8Y+/8H/f8Lf/8Rf3/v/wZ//wK//wTvvf8Dv/8B9/8Kv/8R+v8Fv/8D/v8I/v8C+/8H/f8C/v7/Cff+/wHv/wW/+/8B3v6//wLv/wbf/wj5/wTv/wvv/wH3/wK//wnf/wf7/wf9/wJ//wF7/wO//wb9/wL9/wO//wX3/wPH+f8B3/8C7/8F3/8G/f8Ov/8Bd+//Af77/wF//wX5f/8Cf+7/A7//A/v/BH//Av67/wN//wP7/wT3+99//wL3/wX+73//An9//wG//wPP/wOe/wH7/wl//wK//wTb/wP9/wHf/xS//wPf/wG//wf7/wb3/wK//xv3/wH3/wf7/xd//wLX/xHv/wHv/wG//f7/Au//Bt//AX//Af3/C/f/A9f/Av3/A/3/B+//Aq/f/wF//wL9/wjX/wLv/wbv/wb8/xH3/wTv+/8E/v8B+/8Df/7/At//Av7/Au//Ad//B+//FO//Cvv/C7//Aff/A3//Ab//BP3/CLf/Bb//C/f+f/8C+v8L+/8L/P8H9/8H/v8Iv/8E9/8Fnt//A7//Afv/Ab//B/f/Af7/Af7/B+//Avv/Bd//Be//Dff/D+//CP7/Bv3/C7//Ba//Bb//Av7/BX//Av3/AX//Af1//wT9/wS//wH3/wN//wLv/wX7/wm//wm//wR/9/8Hf/8Df/8b/v8F/f8K9/8N9/8G9/8C/f8C3/8Hv/8C7/8Bb/8F3/8B3/8I3/8L/v8Fr/8F7/8F/f8K/fv/CO//Bfv/GH/+/wTf/wf3/wHv/wVf/wz9/f8Lf3//Fvf/At//A/7/A+3/Au//Ar/+7/8B7/8E3/8M+f8Pf/8Ff/8Dv3//Aff/Bvf/At//Hvff9/8B5/8C9/8F3+/+/wH7/xDf/wP+/xL3/xTv/wr+/wP3/wXf/wT3v/8o/u//Af7/BO/3/f8Cv/8B+3//Avv/E3//C9/f/wO//wvvf/8F+/8F+/8I9/8Df/v/Af3/A/v9/wG//wT+/w7v/wV//wh//wK//wu/f/8C8/8Y9v8Gv7/f/wP9/wHz/wm//wa+/xH77/8P77//A/3v/wr7/wj3/wf+/wL9/wf+/wP3/wb+/wHf/wH3/w3f/f8F/v8O9/8J/f8Fv/3/De//Gvv/Bvv/B+//Dff9/wJ//wt//wLf/wN//wJ3/wJ//yd/f/8P/v8Bv/8Q+/8I/d//Ab9//wP7z/8D/f8Z77//Aff/A9v/BP3/Bvf/Ajb/AX//Bt//CP3f/wG//wHf/wvv/wTv/w7v+/8B9/8G9/8G/f8Cv/8B/v8B+/8Cv/8I7/8B/v8E98//Ct//C7//A/v/Bbv/Cvf/Aff/At8//wv+/wP+/wW//wH6/wTf/wP7/wnv/wP3/wl/f/8Hf/8E+9//Dvz+/wnv/xd//wL7/wf9/wL7/wPfv/8E3/8Df/8Cf/8Cv/8Df/8B9/8d9/8Cf+//Bn//Bff/B/v99/8Sf/8Gf/8E3/8E3/8C/v8G+/8Bf/8I9f7/Bf7/Avf/Bv7/At//Ap/v/wjv/wTv/wb9/wG//wF//wH+/wP9/wL7/wa//xDv/wG//wH3/wX7/we//w3f/wH9/wH9/v7/Au//Bd7/Ab//Avvn9/8C/f8Pffv/Bu//BPv/Be//EH//Dv7/Af7/CM77/wT+/wL3/wT33/8C79//Bfv/FL//A/3/DOb/Hnv/DOf/Ed//BP37/wr8/xX9/wh//we//wHv/xL1/wz9/wN//wS//wLr/wp//wO//wP9/wf3/wF//wH7/wPv/wH3/wZ//xPp/xD9/wX9/y73/w79/wH7+/8D93//Af3/Ar/+9/f/BX//Av7/Avf/AZ//Cfv/Ef73/wT7/wb3/wZ//xPf/wK//wGfv/8C+/8K9/8Df/8C9/8G93//A9//EO//A3//Ae/+/wH9/wv1/wL9/f8D/v8E/d//E/t//wf+/wL+/wnf/wLf/wLf/wy//wX7/wLrf/8C7/8B3/8D+f8Hv/8G8/8N+/3/Avv/Ed/3/wfv/xDv/wfv/wTv/wj9/wZf/wff+/8C7/r/D+//BX+//xHv/wv7/wH1/wW//wy+/wG//w3z/wH3/wv9/wPf/wLv/xL+/wj8/wXf/wff/wL3/wH9/wb7/wX++/8F9/8C/v8D+f8Ov/8Lf/8E7/8Cf/8D9/8Hf/8F7/8C7/8H/v7/Bv3/Au//E/7/Au//E+3+/wP3/wG//wHv/wF//wb9/wff/wXf/wa/f3//A7//Ab//Cff/Dv3/DX/3/w3f/wPf/w7v/w/3/wP7/xv9/wL+/wX+/wW//wW//wh//wzv/xn+/wO//wXf/wH+/wV//wHt/wL9/wF//v8Gv/8U7/8K3/8G7/8G9/8K/v8O33//Au//AX//CP7/Bf3/Ar//C7//CP5//wN//w/3/wv+/wP3/wy//wXv/wj37/8If/8B9v8J/f8h/f8C5/8J/f7/Cnv/A5//A7//BL//CN//A9//Cfv9/wT3/wj9/wjv/wTv/wXf/wb3/wm//wP7/wX3/w7f/wPv/wT+/wLf/wK//wLv3/8U/n//D+//Ae7/DL//BL//GPv/CP7/BH//Be//Afv/Afr/Ab//Cs//CP3/A7/f/wJ//wG/P/7/Af3/AfP/Ce//Av7/Avn/Bfv/Av3/Ab//Ce//A9//Bev+/wP3/wH6/wZ/f/8I9/8Bf/8K/f8I/f8Tf/8E+/8C1/8D+/8E7/v/F/f+/wm//wf8/wK9/xr9/wb9/wS//wj9/wb7/wLv/wT3/wK//wb+/w79/wv7/f8F9/v/BL//Ar//Cb9//wH9/wH3/wj7/wT7/wP+/v8I+9//CvX/Au//Av3/CPf/EN//Ab//C/v/CO//B7//Bd//C/3/D/7/Bv7/Ab//Avv/EPuv/wHf/wTX/wfv/wb7/wG//w73/wX73/8R/v8O/f8S+/f/Be7f/wL7/wP3/wff9/8B/v8G3/8Bf/8R/f8C+/3/A/f/BOb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/wjf/xff/wV//wXf/wHv/xC//w/7/wj3/wT9/wXn/wTv/xDd/wf3/w3f/wX3/wq+/w39/wF//xD7/wP7/wvv/wPv/f8C+f8J7/8F3/8Hf/f/At//AfH/Cf7/A+//An/9/v8F9/8N9/8G3/8ef/8Bf/8K/v8G9/8Fv/8E7/8H+/8G3/8F/v8D/f8Ev/8L3/8BP/8D7/8G+/8C3/8G9/7/Fb//Bb7/Dfv/Bv7/Avv3/xXv/wTf3/8Ff/8Bf/ofnl6qYFh//wH73/8B",
		"r9lABgAAAAAZDkEGAAAAAAIIAggF/v8O33//Au//AX//CP7/CJ//C7//A9//BP5//wN//w/3/wv+/wP3/xLv/wj37/8K9v8u5/8K/v8Ke/8Dn/8Iv/8I3/8D3/8J+/8F9/8I/f8I7/8E7/8F3/8G9/8Jv/8De/8F9/8O3/8D7/8E/v8C3/8Cv/8C79//FP5//wK//wzv/wHu/wy//x37/wj+/wR//wP9/wHv/wH7/wH+/wG//wrf/wj9/wTfv/8Bf/8Bvz/+/wH9/wH7/wX7/wPv/wL+/wL5/wX7/wL9/wG//w3f/wX7/v8D9/8B+v8Gf3//CPf/AX//Cv3/GH//A3//BPv/Atf/A/v/BO/7/xP7/wP3/v8Jv/8H/f8Cvf8V/f8E/f8E+/8B/f8Ev/8I/f8G+/8C7/8Hv/8G/v8O/f8L+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D9v7/CPvf/wr1/wLv/wL9/wj3/xDf/wG//wH9/wn7/wZ//wHv/we//wPv/wHf/wv9/w/+/wb+/wG//xP7r/8B3/8E1/8H7/8G+/8Q9/8F+9//Ef7/Dv3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/CH//Cff/B/3/Avv/BPf/BOb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/w/+/xDf/wV//wXf/wHv/wrv/wW//w/7/wW//wL3/wT9/wXn7/8D7/8Q3f8H9/8Dv/8J3/8F9/8Kvv8N/f8Bf/8Q+/8D+/8L7/8D7/3/Avv/Ce//Bd//B3/3/wLf/wHx/wn+/wPv3/8Bf/3+/wX3/wT7/wj3/wbf/x5//wF//wr+/wb3/wW//wTv/w7f/wX+/wP9/xDf/wE//wH3/wHv/wb7/wLf/wb3/xa//wW+/xT+/wL79/8J+/8Q39//BX//AX/6H55eqmBYf/8B+/8C",
		"wKBABgAAAAAjDkEGAAAAAAIIAgj/BPvf/wH+/wz8/v8G7/8C7/8S7/8Ef/8C+/8H/f8C+/8Ev/8E3/8C+3//An//Ar//A3//Aff/Fd//B/f/An/v/wZ//wS/9/8H+/33/xl//wnf/wL+/wb733//CPX/Bv7/Avf/Af3/BP7/At//Ap//Ce//BO//A/f/Av3/Ab//AX//Ar//Av3/Avv/F+//Ab//Aff/Bfv/B7//Afv/C9//Af3/Av7+/wLv/wXf/wG//wL75/f/Av3/D337/wL9/wPv/wT7/wLf/wLv/wfv/wh/7/8P/v8I3vv/BP7/Avb/BPff/wPf/wP+/wH7/xS//wP9/wzm/x57/wn3/wLn/xHf/wT9+/8K/P8V/f8If/8Hv/8B7/8B/v8Q9f8Hv/8E/f8Df/8Ev/8C6/8H/f8Cf/8Dv/8L9/8Bf/8F7/8B9/8Gf/8B/v8M/v8E+f8W/f8e9/8P9/8O/f8B+/v3/wL3f/8B/f8Cv/739/8B7/8Df/8C/v8C9/8B3/8J+/8R/vf/C/f/Gt//Ar//AR+//wL7/wr3/wb3/wa3f/8D3/8Q7f8Df/8B7/7/BH//CPX/Av39/wP+/wT+3/8L+/8H+3//Cv7/Cd/9/wTf/xL7/wLrf/8C7/8B3/8D+/8Hv/8G8/8N+/3/Avv/Ed/3/wHf/wXv/wf+/wjv/wfv/w39/wZfv/8G3/v/Au/6/w/v/wV/v/8R7/7/Cvv/Af3/Bb//DL//Ab//DfP/Aff/C/3/Afv/Ad//Au//Br//C/7/CPz/Dd//Avf/Af3/Bvv/Bf77/wTv9/8C/v8D+f8Ov/8Lf/8E7/8Cf/8D9/8Hf/8F7/8B/u//B/7+/wb9/wLv/wJ//xD+/wLv/xPt/v8D9/8Bv/8B7/8Bf/8G/f8N3/8Gv39//wO//wv3/w79/w1/9/8N3/8DX/8O7/8J3/8F9/8D+/8P+/8L/f8C/v8Lv/8Fv/8If/8M7/8Z/v8J3/8B/v8Ff/8B7/8C9f8Bf/7/G+//Be//BN//Bu//Ae//BPf/Cv7/Dt9//wLv/wF//wj+/wi//wu//wj+f/8Df/8P8/8L/v8D9/8S7/8I9+//Cvb/Luf/Cv7/Cnv/A5//BH//A7//CN//Avvf/wn7/wP9/wH3/wj9/wG//wbv/wL3/wHv/wXf/wb3/wm//wP7/wX3/w7f/wPv/wT+/wHv3/8Cv/8C79//A/3/EP5//w/v/wHu/wy//x37/wj+/wR//wL3/wLv/wH7/wH+/wG//wP9/wbf/wT3/wP9/wTf/wJ//wG9P/7/Af3/Afv/Ce//Av7/Avn/Bfv/Av3/Ab//Dd//Bfv+/wP3/wH6/wZ/f/8I9/8Bf/8K/f8Mf/8Pf/8E+/8C1/8Bv/8B+/8E7/v/Bd//Eff+/wF//we//wf9/wK9/w/9/wX9/wT9/wb9/wS//wj9/wb7/wLv/we//wb+/wT+/wn9/wv7/f8C/f8C9/v/BL//Ar//Cb9//wH99/f/CPv/BPv/A/7+/wj73/8K9f8C7/8C/f8I9/8Q3/8Bn/8L+/8I7/8Hn/8F3/8L/f8P/v8G/v8Bv/8T+6//Ad//BNf/B+//Bvv/Dfv/Avf/Bfvf/w3v/wP+/wj3/wX9/xL79/8F7t//Avv/A/f/B9/3/v7/CH//Bvf/Cv3/Avv/Af3/Avf/BOb/Cvf/Af7/AX//Aff9v/8Cb/X9u/8E7/8G/v8g3/8Ff/8F3/8B7/8F+/8Kv/8C/v8M+/8Cv/8F9/8E/f8F5/8E7/8H9/8I3f8H9/8N3/8F9/8B+/8Ivv8B/f8L/f8Bf/8Q+/8D+/8I9/8C7/8D7/3/Ab/7/wnP/wXf/wV//wF/9/8C3/8B8f8J/v8D7/8Cf/3+/wX3/w33/wbf/x5//wF7/wq+/wb3/wW//wTv/w7f/wX+/wP9/wZ//wnf/wE//wPv/wb7/wLf/wb3/xXvv/8Fvv8U/v8C+/f/CPf/DPf/BN/f/wV//wF/+h+eXqtgWMAf+/8C",
		"x8BABgAAAAAjDkEGAAAAAAIIAgj/BN//BN//Bd/v/xH7/wLrf/8C7/8B3/8D+/8C3/8Ev/8G8/8N+/3/Avv/Ed/3/wfv/wv+/wTv/wfv/w39/wZf3/8G3/v/Au/6/w/P/wV/v/8R7/8L+/8B/f8Fv/8Mv/8Bv/8N8/8B9/8L/f8D3/8C7/8S/v8I/P8N3/8C9/8B/f8G+/8F/vv/Bff/Av7/A/n/Dr//C3//Af7/Au//An//A/f/B3//Be//Au//B/7+/wb9/wLv/xP+/wLv/xPt/v8D9/8Bv/8B7/8Bf/8G/f8N3/8Gv39//wO//wv3/w79/wj7/wR/9/8N3/8D3/8O7/8P9/8D+/8b/f8C/v8Lv/8Fv/8F+/8Cf/8Jv/8C7/8H9/8R/v8J3/8B/v8Ff/8B7/8C/d9//v8I+/8K+/8H7/8K37//Be//Bvf/Cv7/Dt9//wLv/wF//wj+/wi//wu//wj+f3//An//DL//Avf/C/7/A/f/Eu//CPfv/wr2/y7n/wr+/wp7/wOf/wi//wjf/wPf/wn7/wX3/wF//wb9/wT9/wPv/wTv/wXb/wb3/wm//wP7/wX3/w7f/wPv/wT+/wLf/wK//wLv3/8Ff/8O+n//D+//Ae7/Cf3/Ar//C3//Cvf/Bvv/A9//BP7/BH//Be//Afv/Af7/Ab//Ct//CP3/BN//An//Ab8//v8B/f8B+/8J7/8C/v8C+f8F+/8C/f8Bv/8E+/8I3/8F+/7/A/f/Afr/Bn9//wj3f3//Cv3/Ab//E/v/Bn//BPt//wHX/wP7/wP+7/v/C9//C/f+/wm//wf9/wK9/wH7/wz3/wb9/wT9/wb9/wS//wj9/wb7/wLv/we//wb+/w79/wf3/wP7/f8F9/v/BL//Ar//Cb9//wH9/wH3/wj7/wT7/wP+/v8I+9//CvX/Au//Av3/A/7/BPf/EN//Ab//C/v/BL//A+//Dd//C/3/D/7/Bv7/Ab//E/uv/wHf/wTX/wfv/wb7/wrf/wX3/wX73/8R/v8O/f8Nv/8E+/f/Be7f/wL7/wP3/wL+/wTf9/8B/v8D/f8Ef/8R/ff/Afv/BPf/BOb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/yDf/wV/3/8E3/8B7/8Qv/8I+/8G+/8F3/8C9/8E/f8F5/8E7/8Q3f8H9/8N3/8Fd/8Kvv8N/f8Bd/8K7/8F+/8D+/8L7/8B/f8B7/3/Avv/B3//Ae//Bd//B3/3/wLf/wHx/wK//wb+/wPv/wJf/f7/Bff/Dff/Bt//Hn//AX//Cv7/Bd/3/wW//wTv/w7f/wX+/wP9/xDf/wE//wPr/wb7/wLf/wb3/wnv/wy//wW+/xT+/wL79/8a39//BX//AX/6H59eq2BYwB/7/wI=",
		"zZtABgAAAAAjDkEGAAAAAAIIAgj7/wH3/wb3/wa9/wK//wH+/wH7/wb7/wTv/wH+/wT3z/8K3/8Lv/8D+/8Bf/8Du/8K9/8B9/8C3z//D/7/Bb//Af7/BN//A/v/Ce//A/f/CX9/7/8Gf/8E+9//Dvz+/wnv/xd//wL7/wf9/wL7/wS//wTf/wN/7/8Bf/8Cv/8Df/8B9/8B/f8b9/8B/X/v/wZ//wX3/wf7/ff/GX//Bfv/A9//Av7/Bvv/AX/f/wf1/wb+/wL3/wb+/wLf/wKf/wnv/wK//wHv/wb9/wG+/wF//wX9/wL7/xfv/wG//wH3/wX7/we//w3f/wH9/wL+/v8C7/8F3/8Bv/8C++f3/wL9/w99+/8G7/8E+/8F7/8K/v8Ff/8Q/v8ITvv/BP7/Avf/BPff/wPf/wX7/xS//wP9/wu/5v8ee/8M5/8R3/8E/fv/Cvz/Ae//E/3/CH//B7//Ae//EvX/DP3/A3//BL/+/wHr/wp/9/8Cv/8L9/8Bf/8F7/8B9/8B+/8Ef/8T+f8E/v8Q9/3/Lvf/Dv3/Afv7/wP3f/8B/f8Cv/739/8Ff/8C/v8C9/8B3/8J+/8R/vf/C/f/Gt//Ar//AZ+//wH++/8K9/8G9/8G93//A9//Afv/Ds//A3//Ae/+/wf+/wX1/wL9/f8D/v8F3/8T+3//Cv7/Cd//Bd//Cf7/CPv/Aut//wLv/wHf/wP7/we+/wbz/w37/f8C+/8R3/f/B+//EO//B+//Df3/Bl//B9/7/wLv+v8I7/8G7/8E33+//wj+/wjv/wS//wb7/wH9/wW//wy//wG//wXf/wfz/wH3/wH9/wn9/wPf/wLv/wbf/wv+/wj8/w3f/wL3/wH9/wTv/wH7/wX++/8F9/8C/v8B7/8B+f8N77//Cvd//wTv/wJ//wL39/8Hf/8F7/8C7/8H/v7/Bv3/Au//E/7/Au//E+3+/wP3/wG//wHv/wF//wb9/w3f/wa/f3//A7//C/f/Dv3/DX/3/w3f/wPf/wXv/wjv/w/3/wP7/w+//wv9+/8B/v8Lv/8Fv/8If/8M7/8C/f8W/v8G+/8C3/8B/v8Ff/8B7/8Bf/3/AX/+/xvv/wrf/wbv/wb3/wr+/w7ff/8C7/8Bf/8I/v8Iv7//Cr//CP5//wN//wq//wT3/wf7/wP+/wP3/wZ//wvv/wG//wb37/8E/f8F9v8Yv/8V5/8K/v8Ke/8Dn/8Iv/8I3/8D3/8J+/8F9/8I/f8I7/8C/f8B7/8F3/8G9/8Jv/8D+/8F9/8If/8F3/8D7/8E/v8B/d//Ar//Au/f/wvv/wj+f/8D/v8L7/8B7t//C7//HP37/wj+/wR//v8E7/8B+/8B/v8Bv/8K3/8H+/3/BN//An//Ab8//v8B/f8B+/8J7/8C+v8C+f8F+/8C/f8Bv/3/C/ff/wX7/v8D8/8B+v8Gf3//CPf/AX//CX/9/w33/w5//wT7/wLX/wP7v/8D7/v/F/f+/wm//wf9/wK9/wHf/xD7/wL9/wT9/wX9/f8Ev/8I/f8C7/8D+/8C7/8Hv/8G/v8O/f8L+/3/Bff7/wS//wK//wm/f/8B/f8B9/8I+/8E+/8D/v7/CPvf/wr1/wLv/wL9/wj3/xDf/wG//wv7/wjv/we//wXf/wv9/w/+/wb+/wG//wH3/w39/wP7r/8B3/8E1/8H7/8G2/8Q9/8F+9//Ef7/B/f/Bv3/Evv3/wXu3/8C+/8D9/8H3/f/Af7/CH//Ef3/Avv/BPf/Ar//Aeb/Cvf/Af7/AX//Aff9/wNv9f27/wTv/wb+/wy//xPf/wV//wXf/wHv/xC//w/7/wj3/wT9/wXn/wTv/wb9/wnd/wf3+/8M3/8D7/8B9/8Kvv8C/f8K/f8Bf/8I/v8H+/8D+/8L7/8D7/3/Avv/Ce//Bd//Af3/BX/3/wLf/wHx/wb3/wL+/wPv/wJ//f7/Bff/Dff/Bt//Hn//AX//Cv7/Bvf/Bb3/BO//C9//At//Bf7/A/3/EN//AT//A+//Bvv/At//Bvf/C9//Cr//Bb7/FP7/Avv3/xrf3/8Ff/8Bf/ofn1+rYFjAH/v/Ag==",
		"+nlABgAAAAAjDkEGAAAAAAIIAggB+rINAAFKAAHf/wjf/wV//wX+/wWv/wXv/wX9/wfv/wL9+/8I7/8T3/8Kf/7/BN//B/f/Ae//Bd//DP39/wt/f/8W9/8C3/8D/v8D7f8C7/8Cv/7v/wHv/xH5/xV//wO/f/8I9/v/Ad//Hvff9/8B5/8C9/8F3+/+/wH7/xDf/wP+/yfv/wr+/wP3/wXf/wW//yj+7/8B/v8E7/f9/wK//wH7f/8i39//A7//C+9//wv7/wj3/wN/+/8B/f8D+/8Cu/8E/v8O7/8Of/8Cv/8Lv3//AvP/GP7/Br+/3/8D/f8B9/8Jv/8Gvv8R++//D++//wP9/wv7/wj3/wf+/wPf/wb+/wP3/wb+/wHf/wH3/w3f/f8F/v8Gf/8H9/8Pv/3/B/3/Be//Gvv/Bvv/B+//Dff9/wJ//wt//wLf/wZ3/wJ//wv7/w/7/wt/f/8Hv/8H/v8Bv/8Z/d//Ab9//wP7z/8D/f8D/v8Sf/8C77//Aff/A9v/BP3/Bvf/Ajb/AX//Bt//CP3/Ar//Ad//C+//E+/7/wH3/wb3/wb9/wK//wH++/v/C+//Af7/BPfP/wrf/wu//wP7/wW7/wr3/wH3/wLfP/8P/v8Ft/8B/v8E3/8D+/8J7/8D9/8Jf3//B3//BPvf/v8N/P79/wjv/wW//xF//wL7/wf9/wL7/wS//wTf3/8Cf/8Cf/8Cv/8Df/8B9/8M/f8Q9/8Cf+//DPf/B/v99/8Zf/8Dv/8F3/8C/v8G+/8Bf/8I9f8G/v8C9/8E9/8B/v8C3/8Cn/8J7/8E7/8G/f8Bv/8Bf/8F/f8C+/8X7/8Bv/8B9/8F+/8Hv/8N3/8B3f8C/v7/Au//Bd//Ab//Avvn9/8C/f8Pffv/Bu//BPv/Be//EH//EP7/CM77/wT+/wL3/wT33/8D3/8F+/8Uv+//Av3/DOb/HP7/AXv/DOf/Fv37/wr8/wXv/w/9/wh//we//wHv/xL1/wz9/wN//wS//wLr/wp//wO//wv3/wF//wXv/wH3/wX3f/8T+f8W/f8N9/8Y7/8H9/8O/f8B+/v/A/d//wH9/wK//vf3/wV//wL+/wL3/wHf/wn7/wn+/wf+9/8L9/8a3/8Cv/8Bn7//Avv/Cvf/Bvf/Bvd//wPf/xDv/wHf/wF//wHv/v8N9f8C/f3/A/7/Bd//E/t//f8J/v8H+/8B3/8F3/8S+/8C63//Au//Ad//A/v3/wa//wbz/wH7/wv7/f8C+/8R3/f/B+//EO//B+//Df3/Bl//B9/7/wLv+v8P7/8Ff7//Ee//BX//Bfv/Af3/Bb//DL//Ab//C7//AfP/Aff/C/3/A9//Au//Ev7/CPz/Dd//Avf/Af3/Bvv/Bf77/wX3/wL+/wP5/w6//wX7/wV//wTv/wJ//wP3/wd//wXv/wLv/wf+/v8G/f8C7/8Kv/8I/v8C7/8F7/8N7f7/A/f/Ab//Ae//AX//Bv3/DP7f/wa/f3//Av6//wv3/wXv/wj9/w1/9/8N3/8B9/8B3/8M/f8B7/8P9/8Cv/v/Ed//Cf3/Av7/C7//Bb//CH//DO//Gf7/Cd//Af7/BX//Ae//Av3/AX/+/xvv/wX3/wTf/wbv/wb3/wr+/wbf/wfff/8C7/8Bf/8I/v8Iv3//Cr//CP5/3/8Cf/8P9/8L/v8D9/8S7/8I9+//Cvb/Kn//A+f/CX/+/wp7/wOf/wi//wjf/wPf/wj3+/8F9/8I/f8I7/8E7/8F3/8G9/8Jv/8D+/8F9/8O3/8D7/8E/v8C3/8Cv/8C79//FP5//w/v/wHu/wHf/wq//x37/wj+/wR//wXv/wH7/wH+/wG//wrf/wP7/wT9/wTf/wJ//wG/P/7/Af3/Afv/Ce//Av7/Avn/Bfv/Av3/Ab//Dd//Bfv+/wP3/wH6/wZ/f/8F/f8C9/8Bf/8K/f8Tv/8If/8E+/8C1/8D+/8E7/v/F/f+/wP3/wW//wf9/wK8/xX9/wT9/wb9/wS//wj9/wb7/wLv/we//wb+/w79/wv7/f8F9/v/BL//Ar//Cb9//wH9/wH3/wj7/wT7/wP+/v8I+9//CvX/Au//Av3/Bfv/Avf/EN//Ab//C/v/CO//Dd//C/3/D/7/Bv7/Ab//Efv/Afuv/wHf/wTX/wfv/wb7/wjf/w373/8R9v8Off8S+/f/Be7f/wL7/wP3/wff9/8B/v8If/8R9f8C+/8E9/8E5v8K9/8B/v8Bf/8B9/3/A2/1/bv/BO//A/v/Av7/IN/9/wR//wXf/wHv/wX+/wq//w/7/wJ//wX3/wT9/wXn/wTv/wTv/wvd/wf3/w3f/wX3/wq+/wL7/wr9/wF//wN//wz73/8C+/8L7/8D7/3/Avv/Ce//Bd//B2/3/wLf/wHx/wb7/wL+/wPv/wJ//f7/Bff/Dff/Bt//HL//AX//AX//Cv7+/wX3/wW//wTv/wL3/wvf/wX+/wP9/wV//wrf/wE//wPv/wb7/wLf/wb3/xa//wW+/xT+/wL79/8a39//BX//AX/6H59fq2BYwB/7/wI=",
	}

	fmt.Printf("%4s   %4s   %4s   %4s   %4s\n", "orig", "1B", "2B", "4B", "8B")
	fmt.Println("==================================")
	for i, dataString := range testData {

		kr := &KnownRounds{}
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			t.Errorf("Failed to decode marshalled known rounds (%d): %+v", i, err)
		}

		err = kr.Unmarshal(data)
		if err != nil {
			t.Errorf("Failed to unmarshal known rounds (%d): %+v", i, err)
		}

		buff := kr.bitStream.marshal1ByteVer2()
		u64b := unmarshal1ByteVer2(buff)
		f1bLen := len(buff)
		if !reflect.DeepEqual(kr.bitStream, u64b) {
			t.Errorf("Failed to marshal and unmarshal 1 byte buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, kr.bitStream, u64b)
		}

		buff = kr.bitStream.marshal2BytesVer2()
		u64b = unmarshal2BytesVer2(buff)
		f2bLen := len(buff)
		if !reflect.DeepEqual(kr.bitStream, u64b) {
			t.Errorf("Failed to marshal and unmarshal 2 bytes buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, kr.bitStream, u64b)
		}

		buff = kr.bitStream.marshal4BytesVer2()
		u64b = unmarshal4BytesVer2(buff)
		f4bLen := len(buff)
		if !reflect.DeepEqual(kr.bitStream, u64b) {
			t.Errorf("Failed to marshal and unmarshal 4 bytes buffer (%d)."+
				"\n\texpected: %X\n\treceived: %X", i, kr.bitStream, u64b)
		}

		buff = kr.bitStream.marshal8BytesVer2()
		u64b = unmarshal8BytesVer2(buff)
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
