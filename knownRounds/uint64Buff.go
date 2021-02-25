///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package knownRounds

import (
	"bytes"
	"encoding/binary"
	jww "github.com/spf13/jwalterweatherman"
	"io"
	"math"
)

const (
	ones   = math.MaxUint64
	zeroes = 0
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
//   https://en.wikipedia.org/wiki/Material_conditional
func (u64b uint64Buff) implies(mask uint64Buff) uint64Buff {
	if len(u64b) != len(mask) {
		jww.FATAL.Printf("REPORT THIS ERROR TO JONO ↓")
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
	copy(ext[:len(u64b)], u64b)
	return ext
}

// convertLoc returns the block index and the position of the bit in that block
// for the given position in the buffer.
func (u64b uint64Buff) convertLoc(pos int) (int, int) {
	// Block index in buffer (position / 64)
	bin := pos >> 6 % len(u64b)

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

// TODO: fix licensing for code below. Code below is derived from github.com/tj/go-rle

// marshal encodes the buffer into a byte slice and compresses the data using
// run-length encoding on the integer level. For this implementation, run
// lengths are only included after one or more consecutive integers of all 1s or
// all 0s. All other data is kept in its original form.
func (u64b uint64Buff) marshal() []byte {
	size := len(u64b)

	if size == 0 {
		return nil
	}

	var b = make([]byte, binary.MaxVarintLen64)
	var buf bytes.Buffer
	var cur = u64b[0]
	var run uint64

	for _, next := range u64b {
		if cur != next {
			n := binary.PutUvarint(b, cur)
			buf.Write(b[:n])
			if run > 0 {
				n := binary.PutUvarint(b, run)
				buf.Write(b[:n])
				run = 0
			}
		}
		if next == zeroes || next == ones {
			run++
		}
		cur = next
	}

	n := binary.PutUvarint(b, cur)
	buf.Write(b[:n])
	if run > 0 {
		n := binary.PutUvarint(b, run)
		buf.Write(b[:n])
	}

	return buf.Bytes()
}

// unmarshal decodes the run-length encoded buffer.
func unmarshal(b []byte) uint64Buff {
	buf := bytes.NewBuffer(b)
	buff := uint64Buff{}

	// Reach each uint out of the buffer
	num, err := binary.ReadUvarint(buf)
	for ; err == nil; num, err = binary.ReadUvarint(buf) {
		if num == zeroes || num == ones {
			run, err := binary.ReadUvarint(buf)
			if err != nil {
				jww.FATAL.Panicf("Failed to unmarshal run-length encoded buffer run: %+v", err)
				return nil
			}
			for i := uint64(0); i < run; i++ {
				buff = append(buff, num)
			}
		} else {
			buff = append(buff, num)
		}
	}

	if err != io.EOF {
		jww.FATAL.Panicf("Failed to unmarshal run-length encoded buffer: %+v", err)
	}

	return buff
}
