// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/woozymasta/tga"
)

// DecodeTGAMask decodes TGA grayscale/truecolor payload into 8-bit mask.
func DecodeTGAMask(data []byte) (*MaskImage, error) {
	decoded, err := tga.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode tga image: %w", err)
	}

	mask, err := MaskFromImage(decoded)
	if err != nil {
		return nil, fmt.Errorf("convert tga to mask: %w", err)
	}

	return mask, nil
}

// EncodeTGAMask encodes 8-bit mask into RLE-compressed grayscale TGA.
func EncodeTGAMask(mask *MaskImage) ([]byte, error) {
	return encodeTGAMaskWithOptions(mask, MaskImageEncodeOptions{
		Format: MaskImageFormatTGA,
	})
}

// encodeTGAMaskWithOptions encodes mask with built-in TGA codec.
func encodeTGAMaskWithOptions(
	mask *MaskImage,
	options MaskImageEncodeOptions,
) ([]byte, error) {
	if err := validateMaskPixelBuffer(mask); err != nil {
		return nil, err
	}

	if mask.Width > 0xFFFF || mask.Height > 0xFFFF {
		return nil, errors.New("mask dimensions exceed tga uint16 bounds")
	}

	gray, err := mask.ToGrayImage()
	if err != nil {
		return nil, err
	}

	tgaOptions := defaultTGAMaskEncodeOptions
	if options.TGAOptions != nil {
		tgaOptions = *options.TGAOptions
	}

	var output bytes.Buffer
	if err = tga.EncodeWithOptions(&output, gray, &tgaOptions); err != nil {
		return nil, fmt.Errorf("encode tga image: %w", err)
	}

	return output.Bytes(), nil
}
