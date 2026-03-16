// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"errors"
	"fmt"
	"image"
	"sync"

	"github.com/woozymasta/tga"
)

const (
	areaFlagsFixedHeaderSize = 16
	areaFlagsSectionBitsSize = 4
	areaFlagsUsageBitsOffset = areaFlagsFixedHeaderSize
)

const (
	// MaskImageFormatTGA is the TGA image format label.
	MaskImageFormatTGA = "tga"
)

var (
	// defaultTGAMaskEncodeOptions keeps TGA output compact for mostly binary masks.
	defaultTGAMaskEncodeOptions = tga.EncodeOptions{
		RLE: true,
	}

	// maskImageCodecs stores optional user-defined image codecs.
	maskImageCodecs = map[string]MaskImageCodec{}
	// maskImageCodecOrder keeps deterministic decoder trial order.
	maskImageCodecOrder []string
	// maskImageCodecsMu protects mask image codec registry.
	maskImageCodecsMu sync.RWMutex
)

// MaskImageDecodeFunc decodes raw image bytes into a CE mask.
type MaskImageDecodeFunc func(data []byte) (*MaskImage, error)

// MaskImageEncodeFunc encodes CE mask payload into raw image bytes.
type MaskImageEncodeFunc func(mask *MaskImage, options MaskImageEncodeOptions) ([]byte, error)

// MaskImageCodec describes optional custom format codec.
type MaskImageCodec struct {
	// Decode decodes image bytes into a mask.
	Decode MaskImageDecodeFunc
	// Encode encodes mask payload into image bytes.
	Encode MaskImageEncodeFunc
}

// MaskImageEncodeOptions defines optional mask image encoding behavior.
type MaskImageEncodeOptions struct {
	// CodecOptions is optional payload passed through to custom encoders.
	CodecOptions any
	// TGAOptions controls TGA-specific encode behavior.
	// Nil value uses package defaults (RLE on).
	TGAOptions *tga.EncodeOptions
	// Format selects target format. Empty value means auto from output path
	// extension in file helpers and TGA in byte helpers.
	Format string
	// BlankLayerFileName is optional base name (without extension) for shared
	// blank layer file used when ReuseBlankLayerFile is enabled.
	BlankLayerFileName string
	// ReuseBlankLayerFile enables deduplicated output for fully empty masks in
	// CEProject layer directory writers by emitting one blank file and linking
	// empty layer files to it when possible.
	ReuseBlankLayerFile bool
}

// AreaFlagsHeader stores fixed `areaflags.map` header values.
type AreaFlagsHeader struct {
	// LayerWidth is raster layer width in pixels.
	LayerWidth uint32
	// LayerHeight is raster layer height in pixels.
	LayerHeight uint32
	// WorldWidth is map world width.
	WorldWidth uint32
	// WorldHeight is map world height.
	WorldHeight uint32
	// UsageBits is per-pixel usage field width in bits (usually 16 or 32).
	// It defines usage-section bytes per pixel and usage-layer count.
	UsageBits uint32
	// Reserved stores value section bit depth marker for compatibility.
	// On disk this marker is written between usage and value sections.
	Reserved uint32
}

// AreaFlagsLayer stores one packed bit-plane.
type AreaFlagsLayer struct {
	// Bits is packed row-major bitmap payload.
	Bits []byte
}

// AreaFlagsMapFile stores decoded `areaflags.map` data.
type AreaFlagsMapFile struct {
	// Layers stores packed bit-planes in logical usage/value bit order.
	// Binary codec expands/collapses native per-pixel packed payload.
	Layers []AreaFlagsLayer
	// Header stores file header values.
	Header AreaFlagsHeader
}

// MaskImage stores one 8-bit top-left grayscale mask.
type MaskImage struct {
	// Pixels is row-major grayscale payload.
	Pixels []byte
	// Width is image width in pixels.
	Width uint32
	// Height is image height in pixels.
	Height uint32
}

// LayerByteSize returns payload byte size of one packed bit-plane.
func (header AreaFlagsHeader) LayerByteSize() (int, error) {
	if header.LayerWidth == 0 || header.LayerHeight == 0 {
		return 0, fmt.Errorf("invalid layer dimensions %dx%d", header.LayerWidth, header.LayerHeight)
	}

	totalPixels := uint64(header.LayerWidth) * uint64(header.LayerHeight)
	totalBytes := (totalPixels + 7) / 8
	if totalBytes > uint64(^uint(0)>>1) {
		return 0, errors.New("layer byte size overflow")
	}

	return int(totalBytes), nil
}

// LayerCount returns number of packed layers.
func (file *AreaFlagsMapFile) LayerCount() int {
	if file == nil {
		return 0
	}

	return len(file.Layers)
}

// ValueBits returns number of value bit-planes inferred from payload.
func (file *AreaFlagsMapFile) ValueBits() (uint32, error) {
	if file == nil {
		return 0, ErrNilAreaFlagsFile
	}

	layerCount := uint64(len(file.Layers))
	if layerCount < uint64(file.Header.UsageBits) {
		return 0, fmt.Errorf(
			"layer count %d is lower than usage bits %d",
			len(file.Layers),
			file.Header.UsageBits,
		)
	}

	valueBits := layerCount - uint64(file.Header.UsageBits)
	if valueBits > uint64(^uint32(0)) {
		return 0, fmt.Errorf("value bit count overflow: %d", valueBits)
	}

	return uint32(valueBits), nil
}

// Validate performs basic internal consistency checks.
func (file *AreaFlagsMapFile) Validate() error {
	if file == nil {
		return ErrNilAreaFlagsFile
	}

	layerByteSize, err := file.Header.LayerByteSize()
	if err != nil {
		return err
	}

	for index := range file.Layers {
		if len(file.Layers[index].Bits) != layerByteSize {
			return fmt.Errorf(
				"layer %d byte size %d mismatch, want %d",
				index,
				len(file.Layers[index].Bits),
				layerByteSize,
			)
		}
	}

	return nil
}

// PixelIndex returns flat pixel index for coordinates.
func (img *MaskImage) PixelIndex(x uint32, y uint32) (int, error) {
	if img == nil {
		return 0, ErrNilMaskImage
	}

	if x >= img.Width || y >= img.Height {
		return 0, fmt.Errorf(
			"pixel coordinates out of bounds (%d,%d) in %dx%d",
			x,
			y,
			img.Width,
			img.Height,
		)
	}

	index := uint64(y)*uint64(img.Width) + uint64(x)
	if index > uint64(^uint(0)>>1) {
		return 0, errors.New("pixel index overflow")
	}

	return int(index), nil
}

// At returns pixel value at coordinates.
func (img *MaskImage) At(x uint32, y uint32) (byte, error) {
	index, err := img.PixelIndex(x, y)
	if err != nil {
		return 0, err
	}

	if index >= len(img.Pixels) {
		return 0, fmt.Errorf("pixel index %d out of range", index)
	}

	return img.Pixels[index], nil
}

// Set sets pixel value at coordinates.
func (img *MaskImage) Set(x uint32, y uint32, value byte) error {
	index, err := img.PixelIndex(x, y)
	if err != nil {
		return err
	}

	if index >= len(img.Pixels) {
		return fmt.Errorf("pixel index %d out of range", index)
	}

	img.Pixels[index] = value
	return nil
}

// NewMaskImage allocates a new grayscale mask.
func NewMaskImage(width uint32, height uint32) (*MaskImage, error) {
	if width == 0 || height == 0 {
		return nil, fmt.Errorf("invalid mask size %dx%d", width, height)
	}

	size := uint64(width) * uint64(height)
	if size > uint64(^uint(0)>>1) {
		return nil, errors.New("mask size overflow")
	}

	return &MaskImage{
		Width:  width,
		Height: height,
		Pixels: make([]byte, int(size)),
	}, nil
}

// ToGrayImage converts mask payload to `*image.Gray`.
func (img *MaskImage) ToGrayImage() (*image.Gray, error) {
	if img == nil {
		return nil, ErrNilMaskImage
	}

	if err := validateMaskPixelBuffer(img); err != nil {
		return nil, err
	}

	gray := image.NewGray(image.Rect(0, 0, int(img.Width), int(img.Height)))
	copy(gray.Pix, img.Pixels)
	return gray, nil
}

// MaskFromImage converts any decoded image into an 8-bit grayscale CE mask.
func MaskFromImage(source image.Image) (*MaskImage, error) {
	if source == nil {
		return nil, errors.New("nil source image")
	}

	bounds := source.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid decoded image size %dx%d", width, height)
	}

	if width > int(^uint32(0)) || height > int(^uint32(0)) {
		return nil, fmt.Errorf("decoded image size overflow %dx%d", width, height)
	}

	mask, err := NewMaskImage(uint32(width), uint32(height))
	if err != nil {
		return nil, err
	}

	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := source.At(x, y).RGBA()
			mask.Pixels[index] = maskValueFromRGBA(r, g, b)
			index++
		}
	}

	return mask, nil
}

// validateMaskPixelBuffer checks pixel payload size against mask dimensions.
func validateMaskPixelBuffer(mask *MaskImage) error {
	expectedSize := uint64(mask.Width) * uint64(mask.Height)
	if uint64(len(mask.Pixels)) != expectedSize {
		return fmt.Errorf(
			"mask pixel size %d mismatch, want %d",
			len(mask.Pixels),
			expectedSize,
		)
	}

	return nil
}

// maskValueFromRGBA converts decoded color channels to CE mask intensity.
func maskValueFromRGBA(r16 uint32, g16 uint32, b16 uint32) byte {
	r := (r16 >> 8) & 0xFF
	g := (g16 >> 8) & 0xFF
	b := (b16 >> 8) & 0xFF
	maxChannel := min(max(max(r, g), b), uint32(0xFF))

	return byte(maxChannel)
}
