// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RegisterMaskImageCodec registers custom image codec for arbitrary format.
func RegisterMaskImageCodec(
	format string,
	codec MaskImageCodec,
) error {
	name := normalizeMaskImageFormat(format)
	if name == "" {
		return errors.New("empty mask image codec format")
	}

	if name == MaskImageFormatTGA {
		return errors.New("cannot override built-in tga codec")
	}

	if codec.Decode == nil && codec.Encode == nil {
		return errors.New("mask image codec has no decode and encode callbacks")
	}

	maskImageCodecsMu.Lock()
	defer maskImageCodecsMu.Unlock()

	if _, exists := maskImageCodecs[name]; !exists {
		maskImageCodecOrder = append(maskImageCodecOrder, name)
	}

	maskImageCodecs[name] = codec
	return nil
}

// UnregisterMaskImageCodec removes custom image codec by format.
func UnregisterMaskImageCodec(format string) {
	name := normalizeMaskImageFormat(format)
	if name == "" || name == MaskImageFormatTGA {
		return
	}

	maskImageCodecsMu.Lock()
	defer maskImageCodecsMu.Unlock()

	delete(maskImageCodecs, name)
	for index := range maskImageCodecOrder {
		if maskImageCodecOrder[index] != name {
			continue
		}

		maskImageCodecOrder = append(
			maskImageCodecOrder[:index],
			maskImageCodecOrder[index+1:]...,
		)
		break
	}
}

// DecodeMaskImage decodes image bytes using TGA and registered custom decoders.
func DecodeMaskImage(data []byte) (*MaskImage, string, error) {
	mask, err := DecodeTGAMask(data)
	if err == nil {
		return mask, MaskImageFormatTGA, nil
	}

	maskImageCodecsMu.RLock()
	codecNames := append([]string(nil), maskImageCodecOrder...)
	maskImageCodecsMu.RUnlock()

	for _, name := range codecNames {
		codec, ok := lookupMaskImageCodec(name)
		if !ok || codec.Decode == nil {
			continue
		}

		decoded, decodeErr := codec.Decode(data)
		if decodeErr == nil {
			return decoded, name, nil
		}
	}

	return nil, "", fmt.Errorf("decode mask image: %w", ErrUnsupportedImageFormat)
}

// DecodeMaskImageAs decodes image bytes as explicitly selected format.
func DecodeMaskImageAs(data []byte, format string) (*MaskImage, error) {
	name := normalizeMaskImageFormat(format)
	if name == "" {
		name = MaskImageFormatTGA
	}

	if name == MaskImageFormatTGA {
		return DecodeTGAMask(data)
	}

	codec, ok := lookupMaskImageCodec(name)
	if !ok || codec.Decode == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedImageFormat, name)
	}

	mask, err := codec.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("decode %q image: %w", name, err)
	}

	return mask, nil
}

// EncodeMaskImage encodes mask to requested image format.
// Empty format defaults to TGA.
func EncodeMaskImage(mask *MaskImage, format string) ([]byte, error) {
	return EncodeMaskImageWithOptions(mask, MaskImageEncodeOptions{
		Format: format,
	})
}

// EncodeMaskImageWithOptions encodes mask using selected built-in or
// registered custom codec.
func EncodeMaskImageWithOptions(
	mask *MaskImage,
	options MaskImageEncodeOptions,
) ([]byte, error) {
	if mask == nil {
		return nil, fmt.Errorf("%w: nil", ErrUnsupportedValue)
	}

	format := normalizeMaskImageFormat(options.Format)
	if format == "" {
		format = MaskImageFormatTGA
	}

	if format == MaskImageFormatTGA {
		return encodeTGAMaskWithOptions(mask, options)
	}

	codec, ok := lookupMaskImageCodec(format)
	if !ok || codec.Encode == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedImageFormat, format)
	}

	data, err := codec.Encode(mask, options)
	if err != nil {
		return nil, fmt.Errorf("encode %q image: %w", format, err)
	}

	return data, nil
}

// LoadMaskImageFile loads one image mask file from any supported format.
func LoadMaskImageFile(path string) (*MaskImage, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("read image %q: %w", path, err)
	}

	formatHint := maskImageFormatFromPath(path)
	if formatHint != "" {
		mask, decodeErr := DecodeMaskImageAs(data, formatHint)
		if decodeErr == nil {
			return mask, formatHint, nil
		}
	}

	mask, format, err := DecodeMaskImage(data)
	if err != nil {
		return nil, "", fmt.Errorf("decode image %q: %w", path, err)
	}

	return mask, format, nil
}

// LoadMaskFile loads one image mask file from any supported format.
func LoadMaskFile(path string) (*MaskImage, error) {
	mask, _, err := LoadMaskImageFile(path)
	return mask, err
}

// SaveMaskImageFile stores one image mask file in format inferred from path
// extension. Empty extension defaults to TGA.
func SaveMaskImageFile(path string, mask *MaskImage) error {
	return SaveMaskImageFileWithOptions(path, mask, MaskImageEncodeOptions{})
}

// SaveMaskImageFileWithOptions stores one image mask file with explicit
// encoding options.
func SaveMaskImageFileWithOptions(
	path string,
	mask *MaskImage,
	options MaskImageEncodeOptions,
) error {
	format := normalizeMaskImageFormat(options.Format)
	if format == "" {
		format = maskImageFormatFromPath(path)
	}

	options.Format = format
	data, err := EncodeMaskImageWithOptions(mask, options)
	if err != nil {
		return fmt.Errorf("encode image %q: %w", path, err)
	}

	if err = writeFile600(path, data); err != nil {
		return fmt.Errorf("write image %q: %w", path, err)
	}

	return nil
}

// LoadTGAMaskFile loads one `.tga` mask file.
func LoadTGAMaskFile(path string) (*MaskImage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read tga %q: %w", path, err)
	}

	mask, err := DecodeTGAMask(data)
	if err != nil {
		return nil, fmt.Errorf("decode tga %q: %w", path, err)
	}

	return mask, nil
}

// SaveTGAMaskFile stores one `.tga` mask file.
func SaveTGAMaskFile(path string, mask *MaskImage) error {
	data, err := EncodeTGAMask(mask)
	if err != nil {
		return fmt.Errorf("encode tga %q: %w", path, err)
	}

	if err = writeFile600(path, data); err != nil {
		return fmt.Errorf("write tga %q: %w", path, err)
	}

	return nil
}

// maskImageFormatFromPath infers image format from output extension.
func maskImageFormatFromPath(path string) string {
	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	return normalizeMaskImageFormat(extension)
}

// normalizeMaskImageFormat canonicalizes image format names.
func normalizeMaskImageFormat(format string) string {
	name := strings.TrimSpace(strings.ToLower(format))
	switch name {
	case "":
		return ""
	case "jpg":
		return "jpeg"
	case "tif":
		return "tiff"
	default:
		return name
	}
}

// lookupMaskImageCodec returns registered custom codec by format name.
func lookupMaskImageCodec(format string) (MaskImageCodec, bool) {
	maskImageCodecsMu.RLock()
	defer maskImageCodecsMu.RUnlock()

	codec, ok := maskImageCodecs[format]
	return codec, ok
}
