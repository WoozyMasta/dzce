// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"encoding/binary"
	"reflect"
	"testing"
)

// maxFuzzTGAPixels bounds image allocations requested by mutated headers.
const maxFuzzTGAPixels = 1 << 22

// FuzzDecodeTGAMask exercises the TGA decoder and mask conversion.
func FuzzDecodeTGAMask(f *testing.F) {
	mask, err := NewMaskImage(2, 2)
	if err != nil {
		f.Fatalf("create TGA fuzz seed: %v", err)
	}
	mask.Pixels = []byte{0, 64, 128, 255}

	seed, err := EncodeTGAMask(mask)
	if err != nil {
		f.Fatalf("encode TGA fuzz seed: %v", err)
	}
	f.Add(seed)
	f.Add([]byte{})
	f.Add(make([]byte, 18))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1<<20 || tgaFuzzInputTooLarge(data) {
			return
		}

		decoded, err := DecodeTGAMask(data)
		if err != nil {
			return
		}
		if err := validateMaskPixelBuffer(decoded); err != nil {
			t.Fatalf("decoded TGA mask has invalid pixel buffer: %v", err)
		}

		reencoded, err := EncodeTGAMask(decoded)
		if err != nil {
			t.Fatalf("encode decoded TGA mask: %v", err)
		}

		roundTripped, err := DecodeTGAMask(reencoded)
		if err != nil {
			t.Fatalf("decode encoded TGA mask: %v", err)
		}
		if !reflect.DeepEqual(decoded, roundTripped) {
			t.Fatal("TGA mask changed after encode/decode round trip")
		}
	})
}

// FuzzDecodeAreaFlagsMap exercises malformed headers and packed sections.
func FuzzDecodeAreaFlagsMap(f *testing.F) {
	seed := &AreaFlagsMapFile{
		Header: AreaFlagsHeader{
			LayerWidth:  8,
			LayerHeight: 1,
			WorldWidth:  80,
			WorldHeight: 10,
			UsageBits:   16,
		},
		Layers: make([]AreaFlagsLayer, 20),
	}
	for index := range seed.Layers {
		seed.Layers[index].Bits = make([]byte, 1)
	}

	seedData, err := EncodeAreaFlagsMap(seed)
	if err != nil {
		f.Fatalf("encode area flags fuzz seed: %v", err)
	}
	f.Add(seedData)
	f.Add([]byte{})
	f.Add(make([]byte, areaFlagsFixedHeaderSize+areaFlagsSectionBitsSize))
	f.Add(areaFlagsOverflowHeader())

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1<<20 {
			return
		}

		decoded, err := DecodeAreaFlagsMap(data)
		if err != nil {
			return
		}

		encoded, err := EncodeAreaFlagsMap(decoded)
		if err != nil {
			t.Fatalf("encode decoded area flags map: %v", err)
		}

		roundTripped, err := DecodeAreaFlagsMap(encoded)
		if err != nil {
			t.Fatalf("decode encoded area flags map: %v", err)
		}
		if !reflect.DeepEqual(decoded, roundTripped) {
			t.Fatal("area flags map changed after encode/decode round trip")
		}
	})
}

// TestDecodeAreaFlagsMapRejectsOverflowingSectionSize guards header arithmetic.
func TestDecodeAreaFlagsMapRejectsOverflowingSectionSize(t *testing.T) {
	if _, err := DecodeAreaFlagsMap(areaFlagsOverflowHeader()); err == nil {
		t.Fatal("DecodeAreaFlagsMap() error = nil, want packed section size error")
	}
}

// areaFlagsOverflowHeader creates a small payload with dimensions that
// overflow the packed-section byte-size calculation if multiplication is unchecked.
func areaFlagsOverflowHeader() []byte {
	data := make([]byte, areaFlagsFixedHeaderSize+2*areaFlagsSectionBitsSize)
	binary.LittleEndian.PutUint32(data[0:4], 1<<31)
	binary.LittleEndian.PutUint32(data[4:8], 1<<31)
	binary.LittleEndian.PutUint32(
		data[areaFlagsUsageBitsOffset:areaFlagsUsageBitsOffset+areaFlagsSectionBitsSize],
		32,
	)
	return data
}

// tgaFuzzInputTooLarge skips headers that could request excessive allocation.
func tgaFuzzInputTooLarge(data []byte) bool {
	if len(data) < 18 {
		return false
	}

	width := uint64(binary.LittleEndian.Uint16(data[12:14]))
	height := uint64(binary.LittleEndian.Uint16(data[14:16]))
	return width*height > maxFuzzTGAPixels
}
