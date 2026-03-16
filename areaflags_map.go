// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
)

// DecodeAreaFlagsMap decodes binary `areaflags.map` payload.
func DecodeAreaFlagsMap(data []byte) (*AreaFlagsMapFile, error) {
	minSize := areaFlagsFixedHeaderSize + areaFlagsSectionBitsSize
	if len(data) < minSize {
		return nil, fmt.Errorf(
			"area flags payload too short: %d bytes",
			len(data),
		)
	}

	header := AreaFlagsHeader{
		LayerWidth:  binary.LittleEndian.Uint32(data[0:4]),
		LayerHeight: binary.LittleEndian.Uint32(data[4:8]),
		WorldWidth:  binary.LittleEndian.Uint32(data[8:12]),
		WorldHeight: binary.LittleEndian.Uint32(data[12:16]),
		UsageBits:   binary.LittleEndian.Uint32(data[areaFlagsUsageBitsOffset : areaFlagsUsageBitsOffset+areaFlagsSectionBitsSize]),
	}

	layerByteSize, err := header.LayerByteSize()
	if err != nil {
		return nil, err
	}

	pixelCount := uint64(header.LayerWidth) * uint64(header.LayerHeight)
	if pixelCount == 0 {
		return nil, errors.New("invalid layer pixel count")
	}

	if pixelCount > uint64(^uint(0)>>1) {
		return nil, fmt.Errorf("pixel count overflow: %d", pixelCount)
	}

	if header.UsageBits == 0 || header.UsageBits > 32 {
		return nil, fmt.Errorf("unsupported usage bit count: %d", header.UsageBits)
	}

	usageLayout, err := newAreaFlagsWordLayout(
		header.LayerWidth,
		header.LayerHeight,
		header.UsageBits,
	)
	if err != nil {
		return nil, err
	}

	usageDataStart := areaFlagsFixedHeaderSize + areaFlagsSectionBitsSize
	usageDataEnd := usageDataStart + usageLayout.sectionByteSize
	if usageDataEnd > len(data) {
		return nil, fmt.Errorf(
			"area flags payload too short for usage section: payload=%d need=%d",
			len(data)-usageDataStart,
			usageLayout.sectionByteSize,
		)
	}

	usageData := data[usageDataStart:usageDataEnd]
	valueBitsStart := usageDataEnd
	valueBitsEnd := valueBitsStart + areaFlagsSectionBitsSize
	if valueBitsEnd > len(data) {
		return nil, fmt.Errorf(
			"area flags payload too short for value bit depth marker: payload=%d need=%d",
			len(data)-valueBitsStart,
			areaFlagsSectionBitsSize,
		)
	}

	valueBits := binary.LittleEndian.Uint32(data[valueBitsStart:valueBitsEnd])
	if valueBits > 32 {
		return nil, fmt.Errorf("unsupported value bit count: %d", valueBits)
	}

	var valueLayout areaFlagsWordLayout
	valueDataStart := valueBitsEnd
	valueDataEnd := valueDataStart
	if valueBits != 0 {
		valueLayout, err = newAreaFlagsWordLayout(
			header.LayerWidth,
			header.LayerHeight,
			valueBits,
		)
		if err != nil {
			return nil, err
		}

		valueDataEnd += valueLayout.sectionByteSize
	} else {
		valueLayout.sectionByteSize = 0
	}

	if valueDataEnd > len(data) {
		return nil, fmt.Errorf(
			"area flags payload too short for value section: payload=%d need=%d",
			len(data)-valueDataStart,
			valueLayout.sectionByteSize,
		)
	}

	if valueDataEnd != len(data) {
		return nil, fmt.Errorf(
			"unexpected trailing bytes in area flags payload: %d",
			len(data)-valueDataEnd,
		)
	}

	valueData := data[valueDataStart:valueDataEnd]

	layerCount64 := uint64(header.UsageBits) + uint64(valueBits)
	if layerCount64 > uint64(^uint(0)>>1) {
		return nil, fmt.Errorf("layer count overflow: %d", layerCount64)
	}

	layerCount := int(layerCount64)
	layers := make([]AreaFlagsLayer, layerCount)
	for index := range layers {
		layers[index].Bits = make([]byte, layerByteSize)
	}

	for y := uint32(0); y < header.LayerHeight; y++ {
		for x := uint32(0); x < header.LayerWidth; x++ {
			usageValue, readErr := usageLayout.read(usageData, x, y)
			if readErr != nil {
				return nil, readErr
			}

			for usageValue != 0 {
				bit := bits.TrailingZeros32(usageValue)
				usageValue &= usageValue - 1

				if bit >= int(header.UsageBits) {
					continue
				}

				if setErr := setPackedBit(
					layers[bit].Bits,
					header.LayerWidth,
					x,
					y,
				); setErr != nil {
					return nil, setErr
				}
			}

			if valueBits == 0 {
				continue
			}

			valueValue, readErr := valueLayout.read(valueData, x, y)
			if readErr != nil {
				return nil, readErr
			}

			for valueValue != 0 {
				bit := bits.TrailingZeros32(valueValue)
				valueValue &= valueValue - 1

				if bit >= int(valueBits) {
					continue
				}

				layerIndex := int(header.UsageBits) + bit
				if setErr := setPackedBit(
					layers[layerIndex].Bits,
					header.LayerWidth,
					x,
					y,
				); setErr != nil {
					return nil, setErr
				}
			}
		}
	}

	file := &AreaFlagsMapFile{
		Header: header,
		Layers: layers,
	}

	if err = file.Validate(); err != nil {
		return nil, err
	}

	return file, nil
}

// EncodeAreaFlagsMap encodes binary `areaflags.map` payload.
func EncodeAreaFlagsMap(value *AreaFlagsMapFile) ([]byte, error) {
	if value == nil {
		return nil, fmt.Errorf("%w: nil", ErrUnsupportedValue)
	}

	if err := value.Validate(); err != nil {
		return nil, err
	}

	if value.Header.UsageBits == 0 || value.Header.UsageBits > 32 {
		return nil, fmt.Errorf("unsupported usage bit count: %d", value.Header.UsageBits)
	}

	valueBits, err := value.ValueBits()
	if err != nil {
		return nil, err
	}

	if valueBits > 32 {
		return nil, fmt.Errorf("unsupported value bit count: %d", valueBits)
	}

	if value.Header.Reserved != 0 && value.Header.Reserved != valueBits {
		return nil, fmt.Errorf(
			"header reserved/value marker mismatch: header=%d payload=%d",
			value.Header.Reserved,
			valueBits,
		)
	}

	usageLayout, err := newAreaFlagsWordLayout(
		value.Header.LayerWidth,
		value.Header.LayerHeight,
		value.Header.UsageBits,
	)
	if err != nil {
		return nil, err
	}

	var valueLayout areaFlagsWordLayout
	if valueBits != 0 {
		valueLayout, err = newAreaFlagsWordLayout(
			value.Header.LayerWidth,
			value.Header.LayerHeight,
			valueBits,
		)
		if err != nil {
			return nil, err
		}
	}

	usageSectionSize64 := uint64(usageLayout.wordsPerRow) *
		uint64(value.Header.LayerHeight) * 4
	valueSectionSize64 := uint64(valueLayout.wordsPerRow) *
		uint64(value.Header.LayerHeight) * 4

	totalSize64 := uint64(areaFlagsFixedHeaderSize+areaFlagsSectionBitsSize) +
		usageSectionSize64 +
		uint64(areaFlagsSectionBitsSize) +
		valueSectionSize64
	if totalSize64 > uint64(^uint(0)>>1) {
		return nil, fmt.Errorf("encoded area flags size overflow: %d", totalSize64)
	}

	totalSize := int(totalSize64)
	output := make([]byte, totalSize)

	binary.LittleEndian.PutUint32(output[0:4], value.Header.LayerWidth)
	binary.LittleEndian.PutUint32(output[4:8], value.Header.LayerHeight)
	binary.LittleEndian.PutUint32(output[8:12], value.Header.WorldWidth)
	binary.LittleEndian.PutUint32(output[12:16], value.Header.WorldHeight)
	binary.LittleEndian.PutUint32(
		output[areaFlagsUsageBitsOffset:areaFlagsUsageBitsOffset+areaFlagsSectionBitsSize],
		value.Header.UsageBits,
	)

	usageDataStart := areaFlagsFixedHeaderSize + areaFlagsSectionBitsSize
	usageDataEnd := usageDataStart + usageLayout.sectionByteSize
	usageData := output[usageDataStart:usageDataEnd]

	valueBitsStart := usageDataEnd
	valueBitsEnd := valueBitsStart + areaFlagsSectionBitsSize
	binary.LittleEndian.PutUint32(output[valueBitsStart:valueBitsEnd], valueBits)
	valueData := output[valueBitsEnd:]

	for y := uint32(0); y < value.Header.LayerHeight; y++ {
		for x := uint32(0); x < value.Header.LayerWidth; x++ {
			usageValue := uint32(0)
			for bit := uint32(0); bit < value.Header.UsageBits; bit++ {
				set, readErr := getPackedBit(
					value.Layers[bit].Bits,
					value.Header.LayerWidth,
					x,
					y,
				)
				if readErr != nil {
					return nil, readErr
				}

				if set {
					usageValue |= uint32(1) << bit
				}
			}

			if err = usageLayout.write(usageData, x, y, usageValue); err != nil {
				return nil, err
			}

			if valueBits == 0 {
				continue
			}

			valueValue := uint32(0)
			for bit := range valueBits {
				layerIndex := value.Header.UsageBits + bit
				set, readErr := getPackedBit(
					value.Layers[layerIndex].Bits,
					value.Header.LayerWidth,
					x,
					y,
				)
				if readErr != nil {
					return nil, readErr
				}

				if set {
					valueValue |= uint32(1) << bit
				}
			}

			if err = valueLayout.write(valueData, x, y, valueValue); err != nil {
				return nil, err
			}
		}
	}

	return output, nil
}

// MaskFromLayer exports one packed layer into 8-bit mask (0 or 255).
func (file *AreaFlagsMapFile) MaskFromLayer(layerIndex int) (*MaskImage, error) {
	if file == nil {
		return nil, ErrNilAreaFlagsFile
	}

	if layerIndex < 0 || layerIndex >= len(file.Layers) {
		return nil, fmt.Errorf(
			"layer index %d out of bounds (count=%d)",
			layerIndex,
			len(file.Layers),
		)
	}

	mask, err := NewMaskImage(file.Header.LayerWidth, file.Header.LayerHeight)
	if err != nil {
		return nil, err
	}

	for y := uint32(0); y < file.Header.LayerHeight; y++ {
		for x := uint32(0); x < file.Header.LayerWidth; x++ {
			mapY := file.Header.LayerHeight - 1 - y
			enabled, readErr := getPackedBit(
				file.Layers[layerIndex].Bits,
				file.Header.LayerWidth,
				x,
				mapY,
			)
			if readErr != nil {
				return nil, readErr
			}

			if enabled {
				if setErr := mask.Set(x, y, 255); setErr != nil {
					return nil, setErr
				}
			}
		}
	}

	return mask, nil
}

// SetLayerFromMask overwrites one packed layer from 8-bit mask values.
func (file *AreaFlagsMapFile) SetLayerFromMask(
	layerIndex int,
	mask *MaskImage,
) error {
	if file == nil {
		return ErrNilAreaFlagsFile
	}

	if layerIndex < 0 || layerIndex >= len(file.Layers) {
		return fmt.Errorf(
			"layer index %d out of bounds (count=%d)",
			layerIndex,
			len(file.Layers),
		)
	}

	if mask == nil {
		return ErrNilMaskImage
	}

	if mask.Width != file.Header.LayerWidth || mask.Height != file.Header.LayerHeight {
		return fmt.Errorf(
			"mask size %dx%d mismatch, want %dx%d",
			mask.Width,
			mask.Height,
			file.Header.LayerWidth,
			file.Header.LayerHeight,
		)
	}

	for index := range file.Layers[layerIndex].Bits {
		file.Layers[layerIndex].Bits[index] = 0
	}

	for y := uint32(0); y < file.Header.LayerHeight; y++ {
		for x := uint32(0); x < file.Header.LayerWidth; x++ {
			value, err := mask.At(x, y)
			if err != nil {
				return err
			}

			if value == 0 {
				continue
			}

			mapY := file.Header.LayerHeight - 1 - y
			if err = setPackedBit(
				file.Layers[layerIndex].Bits,
				file.Header.LayerWidth,
				x,
				mapY,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
