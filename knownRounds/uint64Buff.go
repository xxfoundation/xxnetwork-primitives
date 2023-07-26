////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package knownRounds

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"

	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
)

const (
	ones = math.MaxUint64
)

type uint64Buff []uint64

// Get returns the value of the bit at the given position.
func (u64b uint64Buff) get(pos int) bool {
	bin, offset := u64b.convertLoc(pos)

	return u64b[bin]>>(63-offset)&1 == 1
}

// set modifies the bit at the specified position to be 1.
func (u64b uint64Buff) set(pos int) {
	bin, offset := u64b.convertLoc(pos)
	u64b[bin] |= 1 << (63 - offset)
}

// set modifies the bit at the specified position to be 1.
func (u64b uint64Buff) clear(pos int) {
	bin, offset := u64b.convertLoc(pos)
	u64b[bin] &= ^(1 << (63 - offset))
}

// clearRange clears all the bits in the buffer between the given range
// (including the start and end bits).
//
// If start is greater than end, then the selection is inverted.
func (u64b uint64Buff) clearRange(start, end int) {

	// Determine the starting positions the buffer
	numBlocks := u64b.delta(start, end)
	firstBlock, firstBit := u64b.convertLoc(start)

	// Loop over every the blocks in u64b that are in the range
	for blockIndex := 0; blockIndex < numBlocks; blockIndex++ {
		// Get index where the block appears in the buffer
		buffBlock := u64b.getBin(firstBlock + blockIndex)

		// Get the position of the last bit in the current block
		lastBit := 64
		if blockIndex == numBlocks-1 {
			_, lastBit = u64b.convertEnd(end)
		}

		// Generate bit mask for the range and apply it
		bm := bitMaskRange(firstBit, lastBit)
		u64b[buffBlock] &= bm

		// Set position to the first bit in the next block
		firstBit = 0
	}
}

func (u64b uint64Buff) clearAll() {
	for i := range u64b {
		u64b[i] = 0
	}
}

// copy returns a copy of the bits from start to end (inclusive) from u64b.
func (u64b uint64Buff) copy(start, end int) uint64Buff {
	startBlock, startPos := u64b.convertLoc(start)

	numBlocks := u64b.delta(start, end)
	copied := make(uint64Buff, numBlocks)

	// Copy all blocks in range
	for i := 0; i < numBlocks; i++ {
		realBlock := u64b.getBin(startBlock + i)
		copied[i] = u64b[realBlock]
	}

	// Set all bits before the start
	copied[0] |= ^bitMaskRange(0, startPos)

	// Clear all bits after end
	_, endPos := u64b.convertEnd(end)
	copied[numBlocks-1] &= ^bitMaskRange(0, endPos)

	return copied
}

// implies applies the material implication of mask and u64b in the given range
// (including the start and end bits) and places the result in masked starting
// at position maskedStart. An error is returned if the range is larger than the
// length of masked.
//
// If u64bStart is greater than u64bEnd, then the selection is inverted.
//
// More info on material implication:
// https://en.wikipedia.org/wiki/Material_conditional
func (u64b uint64Buff) implies(mask uint64Buff) uint64Buff {
	if len(u64b) != len(mask) {
		jww.FATAL.Panicf("Cannot imply two buffers of different lengths "+
			"(%v and %v).\nu64b: %064b\nmask: %064b", len(u64b), len(mask), u64b, mask)
	}
	result := make(uint64Buff, len(u64b))

	for i := 0; i < len(u64b); i++ {
		result[i] = ^mask[i] | u64b[i]
	}
	return result
}

// extend increases the length of the buffer to the given size and fills in the
// values with zeros.
func (u64b uint64Buff) extend(numBlocks int) uint64Buff {
	// The created buffer is all zeroes per go spec
	ext := make(uint64Buff, numBlocks)
	copy(ext[:], u64b[:])
	return ext
}

// convertLoc returns the block index and the position of the bit in that block
// for the given position in the buffer.
func (u64b uint64Buff) convertLoc(pos int) (int, int) {
	// Block index in buffer (position / 64)
	bin := (pos / 64) % len(u64b)

	// Position of bit in block
	offset := pos % 64

	return bin, offset
}

func (u64b uint64Buff) convertEnd(pos int) (int, int) {
	bin := (pos - 1) / 64

	offset := (pos-1)%64 + 1

	return bin, offset
}

// getBin returns the block index in the buffer for the given absolute index.
func (u64b uint64Buff) getBin(block int) int {
	return block % len(u64b)
}

// delta calculates the number of blocks or parts of blocks contained within the
// range between start and end. If the start and end appear in the same block,
// then delta returns 1.
func (u64b uint64Buff) delta(start, end int) int {
	if end == start {
		return 1
	}
	end--
	if end < start {
		return len(u64b) - start/64 + end/64 + 1
	} else {
		return end/64 - start/64 + 1
	}
}

// deepCopy returns a copy of the buffer.
func (u64b uint64Buff) deepCopy() uint64Buff {
	buff := make(uint64Buff, len(u64b))
	copy(buff, u64b)
	return buff
}

// bitMaskRange generates a bit mask that targets the bits in the provided
// range. The resulting value has 0s in that range and 1s everywhere else.
func bitMaskRange(start, end int) uint64 {
	s := uint64(math.MaxUint64 << uint(64-start))
	e := uint64(math.MaxUint64 >> uint(end))
	return (s | e) & (getInvert(end < start) ^ (s ^ e))
}

func getInvert(b bool) uint64 {
	switch b {
	case true:
		return math.MaxUint64
	default:
		return 0
	}
}

// Word sizes for each marshal/unmarshal function.
const (
	u8bLen  = 1
	u16bLen = 2
	u32bLen = 4
	u64bLen = 8
)

// Current encoding version written to the marshaled buffer. Increment this on
// changes.
const currentVersion = 2

// Map used to select correct unmarshal for the data version.
var u64bUnmarshalVersion = map[uint8]map[uint8]func(b []byte) (uint64Buff, error){
	currentVersion: {
		u8bLen:  unmarshal1ByteVer2,
		u16bLen: unmarshal2BytesVer2,
		u32bLen: unmarshal4BytesVer2,
		u64bLen: unmarshal8BytesVer2,
	},
}

// marshal encodes the buffer into a byte slice and compresses the data using
// run-length encoding on the integer level. For this implementation, run
// lengths are only included after one or more consecutive words of all 1s or
// all 0s. All other data is kept in its original form.
//
// The data is encoded along with the word size used and the version in the
// following structure:
// +---------+-----------+----------+
// | version | word size |   data   |
// | 1 byte  |  1 byte   | variable |
// +---------+-----------+----------+
func (u64b uint64Buff) marshal() []byte {
	return append([]byte{currentVersion, u8bLen}, u64b.marshal1ByteVer2()...)
}

// unmarshal decodes the run-length encoded buffer.
func unmarshal(b []byte) (uint64Buff, error) {
	if len(b) < 3 {
		return nil, errors.Errorf("marshaled bytes length %d smaller than "+
			"minimum %d", len(b), 3)
	}

	unmarshalWordMap, exists := u64bUnmarshalVersion[b[0]]
	if !exists {
		return nil, errors.Errorf("encoding version %d unrecognized", b[0])
	}

	unmarshal, exists := unmarshalWordMap[b[1]]
	if !exists {
		return nil, errors.Errorf("encoding word size %d unrecognized", b[1])
	}

	return unmarshal(b[2:])
}

func (u64b uint64Buff) marshal1ByteVer2() []byte {
	if len(u64b) == 0 {
		return nil
	}

	u8b := make([]uint8, 0, len(u64b)*8)
	for _, u64 := range u64b {
		u8b = append(u8b,
			uint8(u64>>56),
			uint8(u64>>48),
			uint8(u64>>40),
			uint8(u64>>32),
			uint8(u64>>24),
			uint8(u64>>16),
			uint8(u64>>8),
			uint8(u64),
		)
	}

	var buf bytes.Buffer
	var cur = u8b[0]
	var run uint8

	if cur == 0 || cur == math.MaxUint8 {
		run = 1
	}
	for _, next := range u8b[1:] {
		if cur != next || run == 0 {
			buf.WriteByte(cur)
			if run > 0 {
				buf.WriteByte(run)
				run = 0
			}
		}
		if next == 0 || next == math.MaxUint8 {
			if run == math.MaxUint8 {
				buf.WriteByte(cur)
				buf.WriteByte(run)
				run = 0
			}
			run++
		}
		cur = next
	}

	buf.WriteByte(cur)
	if run > 0 {
		buf.WriteByte(run)
	}

	return buf.Bytes()
}

func unmarshal1ByteVer2(b []byte) (uint64Buff, error) {
	buf := bytes.NewBuffer(b)
	var u8b []uint8
	var err error
	var num uint8

	// Reach each uint out of the buffer
	for num, err = buf.ReadByte(); err == nil; num, err = buf.ReadByte() {
		if num == 0 || num == 0xFF {
			run, err := buf.ReadByte()
			if err != nil {
				jww.FATAL.Panicf("Failed to read next byte: %+v", err)
			}
			runBuf := make([]uint8, run)
			for i := range runBuf {
				runBuf[i] = num
			}
			u8b = append(u8b, runBuf...)
		} else {
			u8b = append(u8b, num)
		}
	}

	if err != io.EOF {
		return nil, errors.Errorf("failed to get next uint8 from buffer: "+
			"%+v", err)
	}

	if len(u8b)%8 != 0 {
		return nil, errors.Errorf("length of uncompressed data (%d) must be "+
			"divisible by 8", len(u8b))
	}

	var u64b uint64Buff

	for i := 0; i < len(u8b); i += 8 {
		u8P0 := uint64(u8b[i]) << 56
		u8P1 := uint64(u8b[i+1]) << 48
		u8P2 := uint64(u8b[i+2]) << 40
		u8P3 := uint64(u8b[i+3]) << 32
		u8P4 := uint64(u8b[i+4]) << 24
		u8P5 := uint64(u8b[i+5]) << 16
		u8P6 := uint64(u8b[i+6]) << 8
		u8P7 := uint64(u8b[i+7])

		u64b = append(u64b, u8P0|u8P1|u8P2|u8P3|u8P4|u8P5|u8P6|u8P7)
	}

	return u64b, nil
}

func write2Bytes(i uint16) []byte {
	b := make([]byte, u16bLen)
	binary.BigEndian.PutUint16(b, i)
	return b
}

func (u64b uint64Buff) marshal2BytesVer2() []byte {
	if len(u64b) == 0 {
		return nil
	}

	u16b := make([]uint16, 0, len(u64b)*4)
	for _, u64 := range u64b {
		u16b = append(u16b,
			uint16(u64>>48),
			uint16(u64>>32),
			uint16(u64>>16),
			uint16(u64),
		)
	}

	var buf bytes.Buffer
	var cur = u16b[0]
	var run uint16
	if cur == 0 || cur == math.MaxUint16 {
		run = 1
	}
	for _, next := range u16b[1:] {
		if cur != next || run == 0 {
			buf.Write(write2Bytes(cur))
			if run > 0 {
				buf.Write(write2Bytes(run))
				run = 0
			}
		}
		if next == 0 || next == math.MaxUint16 {
			if run == math.MaxUint16 {
				buf.Write(write2Bytes(cur))
				buf.Write(write2Bytes(run))
				run = 0
			}
			run++
		}
		cur = next
	}

	buf.Write(write2Bytes(cur))
	if run > 0 {
		buf.Write(write2Bytes(run))
	}

	return buf.Bytes()
}

func unmarshal2BytesVer2(b []byte) (uint64Buff, error) {
	buf := bytes.NewBuffer(b)
	var u16b []uint16

	// Reach each uint out of the buffer
	bb := buf.Next(u16bLen)
	for ; len(bb) == u16bLen; bb = buf.Next(u16bLen) {
		num := binary.BigEndian.Uint16(bb)
		if num == 0 || num == math.MaxUint16 {
			run := binary.BigEndian.Uint16(buf.Next(u16bLen))
			runBuf := make([]uint16, run)
			for i := range runBuf {
				runBuf[i] = num
			}
			u16b = append(u16b, runBuf...)
		} else {
			u16b = append(u16b, num)
		}
	}

	if len(bb) != 0 {
		return nil, errors.Errorf("extraneous data of length %d found at end "+
			"of buffer", len(bb))
	} else if len(u16b)%4 != 0 {
		return nil, errors.Errorf("length of uncompressed data (%d) must be "+
			"divisible by 4", len(u16b))
	}

	var u64b uint64Buff

	for i := 0; i < len(u16b); i += 4 {
		u16P0 := uint64(u16b[i]) << 48
		u16P1 := uint64(u16b[i+1]) << 32
		u16P2 := uint64(u16b[i+2]) << 16
		u16P3 := uint64(u16b[i+3])

		u64b = append(u64b, u16P0|u16P1|u16P2|u16P3)
	}

	return u64b, nil
}

func write4Bytes(i uint32) []byte {
	b := make([]byte, u32bLen)
	binary.BigEndian.PutUint32(b, i)
	return b
}

func (u64b uint64Buff) marshal4BytesVer2() []byte {
	if len(u64b) == 0 {
		return nil
	}

	u32b := make([]uint32, 0, len(u64b)*2)
	for _, u64 := range u64b {
		u32b = append(u32b,
			uint32(u64>>32),
			uint32(u64),
		)
	}

	var buf bytes.Buffer
	var cur = u32b[0]
	var run uint32
	if cur == 0 || cur == math.MaxUint32 {
		run = 1
	}
	for _, next := range u32b[1:] {
		if cur != next || run == 0 {
			buf.Write(write4Bytes(cur))
			if run > 0 {
				buf.Write(write4Bytes(run))
				run = 0
			}
		}
		if next == 0 || next == math.MaxUint32 {
			if run == math.MaxUint32 {
				buf.Write(write4Bytes(cur))
				buf.Write(write4Bytes(run))
				run = 0
			}
			run++
		}
		cur = next
	}

	buf.Write(write4Bytes(cur))
	if run > 0 {
		buf.Write(write4Bytes(run))
	}

	return buf.Bytes()
}

func unmarshal4BytesVer2(b []byte) (uint64Buff, error) {
	buf := bytes.NewBuffer(b)
	var u32b []uint32

	// Reach each uint out of the buffer
	bb := buf.Next(u32bLen)
	for ; len(bb) == u32bLen; bb = buf.Next(u32bLen) {
		num := binary.BigEndian.Uint32(bb)
		if num == 0 || num == math.MaxUint32 {
			run := binary.BigEndian.Uint32(buf.Next(u32bLen))
			runBuf := make([]uint32, run)
			for i := range runBuf {
				runBuf[i] = num
			}
			u32b = append(u32b, runBuf...)
		} else {
			u32b = append(u32b, num)
		}
	}

	if len(bb) != 0 {
		return nil, errors.Errorf("extraneous data of length %d found at end "+
			"of buffer", len(bb))
	} else if len(u32b)%2 != 0 {
		return nil, errors.Errorf("length of uncompressed data (%d) must be "+
			"divisible by 2", len(u32b))
	}

	var u64b uint64Buff

	for i := 0; i < len(u32b); i += 2 {
		u16P0 := uint64(u32b[i]) << 32
		u16P1 := uint64(u32b[i+1])

		u64b = append(u64b, u16P0|u16P1)
	}

	return u64b, nil
}

func write8Bytes(i uint64) []byte {
	b := make([]byte, u64bLen)
	binary.LittleEndian.PutUint64(b, i)
	return b
}

func (u64b uint64Buff) marshal8BytesVer2() []byte {
	if len(u64b) == 0 {
		return nil
	}

	var buf bytes.Buffer
	var cur = u64b[0]
	var run uint64
	if cur == 0 || cur == math.MaxUint64 {
		run = 1
	}
	for _, next := range u64b[1:] {
		if cur != next || run == 0 {
			buf.Write(write8Bytes(cur))
			if run > 0 {
				buf.Write(write8Bytes(run))
				run = 0
			}
		}
		if next == 0 || next == math.MaxUint64 {
			if run == math.MaxUint64 {
				buf.Write(write8Bytes(cur))
				buf.Write(write8Bytes(run))
				run = 0
			}
			run++
		}
		cur = next
	}

	buf.Write(write8Bytes(cur))
	if run > 0 {
		buf.Write(write8Bytes(run))
	}

	return buf.Bytes()
}

func unmarshal8BytesVer2(b []byte) (uint64Buff, error) {
	buf := bytes.NewBuffer(b)
	buff := uint64Buff{}

	bb := buf.Next(u64bLen)
	for ; len(bb) == u64bLen; bb = buf.Next(u64bLen) {
		num := binary.LittleEndian.Uint64(bb)
		if num == 0 || num == math.MaxUint64 {
			bb = buf.Next(u64bLen)
			if len(bb) != u64bLen {
				return nil, errors.New("failed to get run")
			}
			run := binary.LittleEndian.Uint64(bb)
			runBuf := make(uint64Buff, run)
			for i := range runBuf {
				runBuf[i] = num
			}
			buff = append(buff, runBuf...)
		} else {
			buff = append(buff, num)
		}
	}

	if len(bb) != 0 {
		return nil, errors.Errorf("extraneous data of length %d found at end "+
			"of buffer", len(bb))
	}

	return buff, nil
}
