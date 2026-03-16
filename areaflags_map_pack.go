// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// areaFlagsWordLayout describes packed 32-bit slot layout used by CEProject map.
type areaFlagsWordLayout struct {
	// bitsPerPixel is section element width in bits.
	bitsPerPixel uint32
	// slotsPerWord is number of elements packed into one uint32 word.
	slotsPerWord uint32
	// wordsPerRow is packed word count per row.
	wordsPerRow uint32
	// width is raster width in pixels.
	width uint32
	// height is raster height in pixels.
	height uint32
	// mask keeps low `bitsPerPixel` bits.
	mask uint32
	// sectionByteSize is encoded section size in bytes.
	sectionByteSize int
}

// newAreaFlagsWordLayout validates and prepares CEProject packed slot layout.
func newAreaFlagsWordLayout(
	width uint32,
	height uint32,
	bitsPerPixel uint32,
) (areaFlagsWordLayout, error) {
	layout := areaFlagsWordLayout{}
	if bitsPerPixel == 0 || bitsPerPixel > 32 {
		return layout, fmt.Errorf("unsupported packed bit count: %d", bitsPerPixel)
	}

	// CEProject stores section values in 32-bit words where pixel slots align by
	// power-of-two bit widths.
	if bitsPerPixel&(bitsPerPixel-1) != 0 {
		return layout, fmt.Errorf("packed bit count must be power of two: %d", bitsPerPixel)
	}

	if 32%bitsPerPixel != 0 {
		return layout, fmt.Errorf("packed bit count must divide 32: %d", bitsPerPixel)
	}

	slotsPerWord := 32 / bitsPerPixel
	if slotsPerWord == 0 {
		return layout, fmt.Errorf("invalid slots per word for %d bits", bitsPerPixel)
	}

	if width%slotsPerWord != 0 {
		return layout, fmt.Errorf(
			"layer width %d must be divisible by slots per word %d",
			width,
			slotsPerWord,
		)
	}

	wordsPerRow := width / slotsPerWord
	wordCount := uint64(wordsPerRow) * uint64(height)
	byteSize64 := wordCount * 4
	if byteSize64 > uint64(^uint(0)>>1) {
		return layout, errors.New("packed section size overflow")
	}

	mask := ^uint32(0)
	if bitsPerPixel < 32 {
		mask = (uint32(1) << bitsPerPixel) - 1
	}

	layout.bitsPerPixel = bitsPerPixel
	layout.slotsPerWord = slotsPerWord
	layout.wordsPerRow = wordsPerRow
	layout.width = width
	layout.height = height
	layout.mask = mask
	layout.sectionByteSize = int(byteSize64)
	return layout, nil
}

// read reads packed section value for one pixel coordinate.
func (layout areaFlagsWordLayout) read(
	section []byte,
	x uint32,
	y uint32,
) (uint32, error) {
	if x >= layout.width || y >= layout.height {
		return 0, fmt.Errorf(
			"pixel coordinates out of bounds (%d,%d) in %dx%d",
			x,
			y,
			layout.width,
			layout.height,
		)
	}

	wordIndex := uint64(layout.wordsPerRow)*uint64(y) + uint64(x/layout.slotsPerWord)
	byteOffset := wordIndex * 4
	if byteOffset+4 > uint64(len(section)) {
		return 0, io.ErrUnexpectedEOF
	}

	word := binary.LittleEndian.Uint32(section[byteOffset : byteOffset+4])
	shift := ((layout.slotsPerWord - (x % layout.slotsPerWord) - 1) * layout.bitsPerPixel) & 0x1f
	return (word >> shift) & layout.mask, nil
}

// write writes packed section value for one pixel coordinate.
func (layout areaFlagsWordLayout) write(
	section []byte,
	x uint32,
	y uint32,
	value uint32,
) error {
	if x >= layout.width || y >= layout.height {
		return fmt.Errorf(
			"pixel coordinates out of bounds (%d,%d) in %dx%d",
			x,
			y,
			layout.width,
			layout.height,
		)
	}

	wordIndex := uint64(layout.wordsPerRow)*uint64(y) + uint64(x/layout.slotsPerWord)
	byteOffset := wordIndex * 4
	if byteOffset+4 > uint64(len(section)) {
		return io.ErrUnexpectedEOF
	}

	shift := ((layout.slotsPerWord - (x % layout.slotsPerWord) - 1) * layout.bitsPerPixel) & 0x1f
	word := binary.LittleEndian.Uint32(section[byteOffset : byteOffset+4])
	fieldMask := layout.mask << shift
	word = (word &^ fieldMask) | ((value & layout.mask) << shift)
	binary.LittleEndian.PutUint32(section[byteOffset:byteOffset+4], word)
	return nil
}

// getPackedBitByIndex reads one LSB-packed bit by flat pixel index.
func getPackedBitByIndex(bits []byte, index int) (bool, error) {
	if index < 0 {
		return false, io.ErrUnexpectedEOF
	}

	byteIndex := index / 8
	if byteIndex >= len(bits) {
		return false, io.ErrUnexpectedEOF
	}

	mask := byte(1 << (index % 8))
	return bits[byteIndex]&mask != 0, nil
}

// setPackedBitByIndex writes one LSB-packed bit by flat pixel index.
func setPackedBitByIndex(bits []byte, index int) error {
	if index < 0 {
		return io.ErrUnexpectedEOF
	}

	byteIndex := index / 8
	if byteIndex >= len(bits) {
		return io.ErrUnexpectedEOF
	}

	mask := byte(1 << (index % 8))
	bits[byteIndex] |= mask

	return nil
}

// getPackedBit reads one LSB-packed row-major bit.
func getPackedBit(
	bits []byte,
	width uint32,
	x uint32,
	y uint32,
) (bool, error) {
	index := uint64(y)*uint64(width) + uint64(x)
	if index > uint64(^uint(0)>>1) {
		return false, io.ErrUnexpectedEOF
	}

	return getPackedBitByIndex(bits, int(index))
}

// setPackedBit writes one LSB-packed row-major bit.
func setPackedBit(
	bits []byte,
	width uint32,
	x uint32,
	y uint32,
) error {
	index := uint64(y)*uint64(width) + uint64(x)
	if index > uint64(^uint(0)>>1) {
		return io.ErrUnexpectedEOF
	}

	return setPackedBitByIndex(bits, int(index))
}
