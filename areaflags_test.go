// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestAreaFlagsEncodeDecodeRoundtrip(t *testing.T) {
	file := &AreaFlagsMapFile{
		Header: AreaFlagsHeader{
			LayerWidth:  8,
			LayerHeight: 8,
			WorldWidth:  80,
			WorldHeight: 80,
			UsageBits:   16,
			Reserved:    0,
		},
		Layers: make([]AreaFlagsLayer, 20),
	}

	layerByteSize, err := file.Header.LayerByteSize()
	if err != nil {
		t.Fatalf("LayerByteSize() error: %v", err)
	}

	for index := range file.Layers {
		file.Layers[index].Bits = make([]byte, layerByteSize)
	}

	if err = setPackedBit(file.Layers[3].Bits, 8, 2, 1); err != nil {
		t.Fatalf("setPackedBit() error: %v", err)
	}

	data, err := EncodeAreaFlagsMap(file)
	if err != nil {
		t.Fatalf("EncodeAreaFlagsMap() error: %v", err)
	}

	decoded, err := DecodeAreaFlagsMap(data)
	if err != nil {
		t.Fatalf("DecodeAreaFlagsMap() error: %v", err)
	}

	if !reflect.DeepEqual(file, decoded) {
		t.Fatal("decoded area flags payload mismatch")
	}
}

func TestAreaFlagsDecodeEncodeInterleavedLayout(t *testing.T) {
	payload := make([]byte, 16+4+16+4+4)
	binary.LittleEndian.PutUint32(payload[0:4], 8)
	binary.LittleEndian.PutUint32(payload[4:8], 1)
	binary.LittleEndian.PutUint32(payload[8:12], 80)
	binary.LittleEndian.PutUint32(payload[12:16], 10)
	binary.LittleEndian.PutUint32(payload[16:20], 16)
	binary.LittleEndian.PutUint32(payload[36:40], 4)

	// usage section (8 pixels x uint16 packed by 32-bit slots):
	// x0=1, x1=2, x2..x7=0 -> first word bytes: 02 00 01 00
	payload[20] = 0x02
	payload[21] = 0x00
	payload[22] = 0x01
	payload[23] = 0x00
	// value section (8 pixels x 4-bit packed in one word):
	// x0=3, x1=4, x2..x7=0 -> bytes: 00 00 00 34
	payload[40] = 0x00
	payload[41] = 0x00
	payload[42] = 0x00
	payload[43] = 0x34

	file, err := DecodeAreaFlagsMap(payload)
	if err != nil {
		t.Fatalf("DecodeAreaFlagsMap() error: %v", err)
	}

	if len(file.Layers) != 20 {
		t.Fatalf("layer count = %d, want 20", len(file.Layers))
	}

	u0, err := getPackedBit(file.Layers[0].Bits, 8, 0, 0)
	if err != nil {
		t.Fatalf("getPackedBit() usage bit0 error: %v", err)
	}

	if !u0 {
		t.Fatal("expected usage bit0 at pixel0")
	}

	u1, err := getPackedBit(file.Layers[1].Bits, 8, 1, 0)
	if err != nil {
		t.Fatalf("getPackedBit() usage bit1 error: %v", err)
	}

	if !u1 {
		t.Fatal("expected usage bit1 at pixel1")
	}

	v0Bit0, err := getPackedBit(file.Layers[16].Bits, 8, 0, 0)
	if err != nil {
		t.Fatalf("getPackedBit() value bit0 error: %v", err)
	}

	v0Bit1, err := getPackedBit(file.Layers[17].Bits, 8, 0, 0)
	if err != nil {
		t.Fatalf("getPackedBit() value bit1 error: %v", err)
	}

	if !v0Bit0 || !v0Bit1 {
		t.Fatal("expected value bits 0 and 1 at pixel0")
	}

	v1Bit2, err := getPackedBit(file.Layers[18].Bits, 8, 1, 0)
	if err != nil {
		t.Fatalf("getPackedBit() value bit2 error: %v", err)
	}

	if !v1Bit2 {
		t.Fatal("expected value bit2 at pixel1")
	}

	encoded, err := EncodeAreaFlagsMap(file)
	if err != nil {
		t.Fatalf("EncodeAreaFlagsMap() error: %v", err)
	}

	if !reflect.DeepEqual(payload, encoded) {
		t.Fatal("interleaved payload roundtrip mismatch")
	}
}

func TestAreaFlagsMaskAndTGARoundtrip(t *testing.T) {
	mask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = mask.Set(1, 2, 255); err != nil {
		t.Fatalf("mask.Set() error: %v", err)
	}

	if err = mask.Set(6, 7, 255); err != nil {
		t.Fatalf("mask.Set() error: %v", err)
	}

	tgaData, err := EncodeTGAMask(mask)
	if err != nil {
		t.Fatalf("EncodeTGAMask() error: %v", err)
	}

	if len(tgaData) < 3 {
		t.Fatalf("tga payload too short: %d", len(tgaData))
	}

	// TGA type 11 means RLE-compressed grayscale.
	if tgaData[2] != 11 {
		t.Fatalf("tga image type = %d, want 11", tgaData[2])
	}

	decodedMask, err := DecodeTGAMask(tgaData)
	if err != nil {
		t.Fatalf("DecodeTGAMask() error: %v", err)
	}

	if !reflect.DeepEqual(mask, decodedMask) {
		t.Fatal("decoded mask mismatch")
	}
}

func TestMaskImageCustomCodecRoundtrip(t *testing.T) {
	const codecFormat = "rawmask"

	err := RegisterMaskImageCodec(codecFormat, MaskImageCodec{
		Decode: decodeRawMaskCodec,
		Encode: encodeRawMaskCodec,
	})
	if err != nil {
		t.Fatalf("RegisterMaskImageCodec() error: %v", err)
	}

	t.Cleanup(func() {
		UnregisterMaskImageCodec(codecFormat)
	})

	mask, err := NewMaskImage(4, 3)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = mask.Set(1, 1, 255); err != nil {
		t.Fatalf("mask.Set() error: %v", err)
	}

	if err = mask.Set(3, 2, 64); err != nil {
		t.Fatalf("mask.Set() error: %v", err)
	}

	raw, err := EncodeMaskImage(mask, codecFormat)
	if err != nil {
		t.Fatalf("EncodeMaskImage() error: %v", err)
	}

	decodedAs, err := DecodeMaskImageAs(raw, codecFormat)
	if err != nil {
		t.Fatalf("DecodeMaskImageAs() error: %v", err)
	}

	if !reflect.DeepEqual(mask, decodedAs) {
		t.Fatal("DecodeMaskImageAs() mismatch")
	}

	decodedAuto, detectedFormat, err := DecodeMaskImage(raw)
	if err != nil {
		t.Fatalf("DecodeMaskImage() error: %v", err)
	}

	if detectedFormat != codecFormat {
		t.Fatalf("detected format = %q, want %q", detectedFormat, codecFormat)
	}

	if !reflect.DeepEqual(mask, decodedAuto) {
		t.Fatal("DecodeMaskImage() mismatch")
	}
}

func TestSaveAndLoadMaskImageFileWithCustomCodec(t *testing.T) {
	const codecFormat = "rawmask-file"

	err := RegisterMaskImageCodec(codecFormat, MaskImageCodec{
		Decode: decodeRawMaskCodec,
		Encode: encodeRawMaskCodec,
	})
	if err != nil {
		t.Fatalf("RegisterMaskImageCodec() error: %v", err)
	}

	t.Cleanup(func() {
		UnregisterMaskImageCodec(codecFormat)
	})

	mask, err := NewMaskImage(2, 2)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = mask.Set(0, 1, 255); err != nil {
		t.Fatalf("mask.Set() error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "custom."+codecFormat)
	if err = SaveMaskImageFile(path, mask); err != nil {
		t.Fatalf("SaveMaskImageFile() error: %v", err)
	}

	loaded, format, err := LoadMaskImageFile(path)
	if err != nil {
		t.Fatalf("LoadMaskImageFile() error: %v", err)
	}

	if format != codecFormat {
		t.Fatalf("format = %q, want %q", format, codecFormat)
	}

	if !reflect.DeepEqual(mask, loaded) {
		t.Fatal("loaded mask mismatch")
	}
}

func TestEncodeMaskImageUnknownFormat(t *testing.T) {
	mask, err := NewMaskImage(1, 1)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	_, err = EncodeMaskImage(mask, "dds")
	if !errors.Is(err, ErrUnsupportedImageFormat) {
		t.Fatalf("EncodeMaskImage() error = %v, want ErrUnsupportedImageFormat", err)
	}
}

func TestBuildAreaFlagsFromCEProjectMasks(t *testing.T) {
	config, err := DecodeCEProjectConfig(readFixture(t, filepath.Join("positive", "ceproject-config.xml")))
	if err != nil {
		t.Fatalf("DecodeCEProjectConfig() error: %v", err)
	}

	militaryMask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	policeMask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = militaryMask.Set(1, 1, 255); err != nil {
		t.Fatalf("militaryMask.Set() error: %v", err)
	}

	if err = policeMask.Set(2, 2, 255); err != nil {
		t.Fatalf("policeMask.Set() error: %v", err)
	}

	file, err := BuildAreaFlagsFromCEProjectMasks(config, map[string]*MaskImage{
		"Military": militaryMask,
		"Police":   policeMask,
	})
	if err != nil {
		t.Fatalf("BuildAreaFlagsFromCEProjectMasks() error: %v", err)
	}

	if file.Header.UsageBits != 16 {
		t.Fatalf("UsageBits = %d, want 16", file.Header.UsageBits)
	}

	if len(file.Layers) != 20 {
		t.Fatalf("layer count = %d, want 20", len(file.Layers))
	}

	// CEProject mapping applies Y flip without additional X offset.
	usageMilitary, err := getPackedBit(file.Layers[1].Bits, 8, 1, 6)
	if err != nil {
		t.Fatalf("getPackedBit() error: %v", err)
	}

	if !usageMilitary {
		t.Fatal("expected usage bit 1 to be set at (1,1)")
	}

	valueTier1, err := getPackedBit(file.Layers[16].Bits, 8, 1, 6)
	if err != nil {
		t.Fatalf("getPackedBit() error: %v", err)
	}

	if !valueTier1 {
		t.Fatal("expected value bit 0 to be set at (1,1)")
	}

	usagePolice, err := getPackedBit(file.Layers[2].Bits, 8, 2, 5)
	if err != nil {
		t.Fatalf("getPackedBit() error: %v", err)
	}

	if !usagePolice {
		t.Fatal("expected usage bit 2 to be set at (2,2)")
	}

	valueTier2, err := getPackedBit(file.Layers[17].Bits, 8, 2, 5)
	if err != nil {
		t.Fatalf("getPackedBit() error: %v", err)
	}

	if !valueTier2 {
		t.Fatal("expected value bit 1 to be set at (2,2)")
	}
}

func TestBuildAreaFlagsFromCEProjectMasksSplitUsageAndValueOffsets(t *testing.T) {
	config := &CEProjectConfigFile{
		Global: &CEProjectGlobal{
			Layer: &CEProjectDimension{Size: "8"},
			World: &CEProjectDimension{Size: "80"},
		},
		Areas: &CEProjectAreas{
			Usages: &CEProjectAreaUsageList{
				Items: []NamedRef{
					{Name: "UsageA"},
				},
			},
			Values: &CEProjectAreaValueList{
				Items: []NamedRef{
					{Name: "ValueA"},
				},
			},
		},
		Layers: &CEProjectLayers{
			Layers: []CEProjectLayer{
				{
					Name:       "UsageOnly",
					UsageFlags: "1",
					ValueFlags: "0",
				},
				{
					Name:       "ValueOnly",
					UsageFlags: "0",
					ValueFlags: "1",
				},
			},
		},
	}

	usageMask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	valueMask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = usageMask.Set(1, 1, 255); err != nil {
		t.Fatalf("usageMask.Set() error: %v", err)
	}

	if err = valueMask.Set(1, 1, 255); err != nil {
		t.Fatalf("valueMask.Set() error: %v", err)
	}

	file, err := BuildAreaFlagsFromCEProjectMasks(config, map[string]*MaskImage{
		"UsageOnly": usageMask,
		"ValueOnly": valueMask,
	})
	if err != nil {
		t.Fatalf("BuildAreaFlagsFromCEProjectMasks() error: %v", err)
	}

	usageSet, err := getPackedBit(file.Layers[0].Bits, 8, 1, 6)
	if err != nil {
		t.Fatalf("getPackedBit() usage error: %v", err)
	}

	if !usageSet {
		t.Fatal("expected usage-only layer at (1,6)")
	}

	valueSet, err := getPackedBit(file.Layers[16].Bits, 8, 1, 6)
	if err != nil {
		t.Fatalf("getPackedBit() value error: %v", err)
	}

	if !valueSet {
		t.Fatal("expected value-only layer at (1,6)")
	}

	exported, err := ExportCEProjectLayerMasksFromAreaFlags(config, file)
	if err != nil {
		t.Fatalf("ExportCEProjectLayerMasksFromAreaFlags() error: %v", err)
	}

	usageMaskOut := exported["UsageOnly"]
	if usageMaskOut == nil {
		t.Fatal("expected UsageOnly mask in export output")
	}

	usagePixel, err := usageMaskOut.At(1, 1)
	if err != nil {
		t.Fatalf("usageMaskOut.At() error: %v", err)
	}

	if usagePixel == 0 {
		t.Fatal("expected usage-only export pixel at (1,1)")
	}

	valueMaskOut := exported["ValueOnly"]
	if valueMaskOut == nil {
		t.Fatal("expected ValueOnly mask in export output")
	}

	valuePixel, err := valueMaskOut.At(1, 1)
	if err != nil {
		t.Fatalf("valueMaskOut.At() error: %v", err)
	}

	if valuePixel == 0 {
		t.Fatal("expected value-only export pixel at (1,1)")
	}
}

func TestExportCEProjectLayerMasksFromAreaFlagsZeroFlagsLayerIsEmpty(t *testing.T) {
	config := &CEProjectConfigFile{
		Global: &CEProjectGlobal{
			Layer: &CEProjectDimension{Size: "8"},
			World: &CEProjectDimension{Size: "80"},
		},
		Areas: &CEProjectAreas{
			Usages: &CEProjectAreaUsageList{
				Items: []NamedRef{{Name: "UsageA"}},
			},
			Values: &CEProjectAreaValueList{
				Items: []NamedRef{{Name: "ValueA"}},
			},
		},
		Layers: &CEProjectLayers{
			Layers: []CEProjectLayer{
				{Name: "ZeroZero", UsageFlags: "0", ValueFlags: "0"},
				{Name: "Usage", UsageFlags: "1", ValueFlags: "0"},
			},
		},
	}

	usageMask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = usageMask.Set(2, 3, 255); err != nil {
		t.Fatalf("usageMask.Set() error: %v", err)
	}

	file, err := BuildAreaFlagsFromCEProjectMasks(config, map[string]*MaskImage{
		"Usage": usageMask,
	})
	if err != nil {
		t.Fatalf("BuildAreaFlagsFromCEProjectMasks() error: %v", err)
	}

	exported, err := ExportCEProjectLayerMasksFromAreaFlags(config, file)
	if err != nil {
		t.Fatalf("ExportCEProjectLayerMasksFromAreaFlags() error: %v", err)
	}

	zeroMask := exported["ZeroZero"]
	if zeroMask == nil {
		t.Fatal("zero/zero layer mask missing")
	}

	for index, value := range zeroMask.Pixels {
		if value != 0 {
			t.Fatalf("zero/zero layer is not empty at pixel %d: %d", index, value)
		}
	}
}

func TestBuildAreaFlagsFromCEProjectMasksMissingLayerDefaultsToEmpty(t *testing.T) {
	config, err := DecodeCEProjectConfig(
		readFixture(t, filepath.Join("positive", "ceproject-config.xml")),
	)
	if err != nil {
		t.Fatalf("DecodeCEProjectConfig() error: %v", err)
	}

	militaryMask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = militaryMask.Set(1, 1, 255); err != nil {
		t.Fatalf("militaryMask.Set() error: %v", err)
	}

	file, err := BuildAreaFlagsFromCEProjectMasks(config, map[string]*MaskImage{
		"Military": militaryMask,
	})
	if err != nil {
		t.Fatalf("BuildAreaFlagsFromCEProjectMasks() error: %v", err)
	}

	policeUsage, err := getPackedBit(file.Layers[2].Bits, 8, 1, 5)
	if err != nil {
		t.Fatalf("getPackedBit() error: %v", err)
	}

	if policeUsage {
		t.Fatal("expected missing Police mask to produce empty layer")
	}
}

func TestBuildAreaFlagsFromCEProjectLayerDirMissingLayerDefaultsToEmpty(t *testing.T) {
	config, err := DecodeCEProjectConfig(
		readFixture(t, filepath.Join("positive", "ceproject-config.xml")),
	)
	if err != nil {
		t.Fatalf("DecodeCEProjectConfig() error: %v", err)
	}

	militaryMask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = militaryMask.Set(1, 1, 255); err != nil {
		t.Fatalf("militaryMask.Set() error: %v", err)
	}

	layerDir := t.TempDir()
	if err = SaveTGAMaskFile(filepath.Join(layerDir, "Military.tga"), militaryMask); err != nil {
		t.Fatalf("SaveTGAMaskFile() error: %v", err)
	}

	file, err := BuildAreaFlagsFromCEProjectLayerDir(config, layerDir)
	if err != nil {
		t.Fatalf("BuildAreaFlagsFromCEProjectLayerDir() error: %v", err)
	}

	policeUsage, err := getPackedBit(file.Layers[2].Bits, 8, 1, 5)
	if err != nil {
		t.Fatalf("getPackedBit() error: %v", err)
	}

	if policeUsage {
		t.Fatal("expected missing Police.tga to produce empty layer")
	}
}

func TestBuildAreaFlagsFromCEProjectLayerDirSupportsCustomMaskCodec(t *testing.T) {
	const codecFormat = "rawmasklayer"

	err := RegisterMaskImageCodec(codecFormat, MaskImageCodec{
		Decode: decodeRawMaskCodec,
		Encode: encodeRawMaskCodec,
	})
	if err != nil {
		t.Fatalf("RegisterMaskImageCodec() error: %v", err)
	}

	t.Cleanup(func() {
		UnregisterMaskImageCodec(codecFormat)
	})

	config, err := DecodeCEProjectConfig(
		readFixture(t, filepath.Join("positive", "ceproject-config.xml")),
	)
	if err != nil {
		t.Fatalf("DecodeCEProjectConfig() error: %v", err)
	}

	militaryMask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = militaryMask.Set(3, 2, 255); err != nil {
		t.Fatalf("militaryMask.Set() error: %v", err)
	}

	layerDir := t.TempDir()
	if err = SaveMaskImageFileWithOptions(
		filepath.Join(layerDir, "Military."+codecFormat),
		militaryMask,
		MaskImageEncodeOptions{Format: codecFormat},
	); err != nil {
		t.Fatalf("SaveMaskImageFileWithOptions() error: %v", err)
	}

	file, err := BuildAreaFlagsFromCEProjectLayerDir(config, layerDir)
	if err != nil {
		t.Fatalf("BuildAreaFlagsFromCEProjectLayerDir() error: %v", err)
	}

	militaryUsage, err := getPackedBit(file.Layers[1].Bits, 8, 3, 5)
	if err != nil {
		t.Fatalf("getPackedBit() error: %v", err)
	}

	if !militaryUsage {
		t.Fatal("expected Military layer from custom codec file to be applied")
	}
}

func TestExportCEProjectLayerMasksFromAreaFlagsWithLayerDirRestoresZeroZero(t *testing.T) {
	const codecFormat = "rawmaskzero"

	err := RegisterMaskImageCodec(codecFormat, MaskImageCodec{
		Decode: decodeRawMaskCodec,
		Encode: encodeRawMaskCodec,
	})
	if err != nil {
		t.Fatalf("RegisterMaskImageCodec() error: %v", err)
	}

	t.Cleanup(func() {
		UnregisterMaskImageCodec(codecFormat)
	})

	config := &CEProjectConfigFile{
		Global: &CEProjectGlobal{
			Layer: &CEProjectDimension{Size: "8"},
			World: &CEProjectDimension{Size: "80"},
		},
		Areas: &CEProjectAreas{
			Usages: &CEProjectAreaUsageList{
				Items: []NamedRef{{Name: "UsageA"}},
			},
			Values: &CEProjectAreaValueList{
				Items: []NamedRef{{Name: "ValueA"}},
			},
		},
		Layers: &CEProjectLayers{
			Layers: []CEProjectLayer{
				{Name: "ZeroZero", UsageFlags: "0", ValueFlags: "0"},
				{Name: "Usage", UsageFlags: "1", ValueFlags: "0"},
			},
		},
	}

	usageMask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = usageMask.Set(4, 3, 255); err != nil {
		t.Fatalf("usageMask.Set() error: %v", err)
	}

	file, err := BuildAreaFlagsFromCEProjectMasks(config, map[string]*MaskImage{
		"Usage": usageMask,
	})
	if err != nil {
		t.Fatalf("BuildAreaFlagsFromCEProjectMasks() error: %v", err)
	}

	zeroMaskSource, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = zeroMaskSource.Set(2, 1, 255); err != nil {
		t.Fatalf("zeroMaskSource.Set() error: %v", err)
	}

	layerDir := t.TempDir()
	if err = SaveMaskImageFileWithOptions(
		filepath.Join(layerDir, "ZeroZero."+codecFormat),
		zeroMaskSource,
		MaskImageEncodeOptions{Format: codecFormat},
	); err != nil {
		t.Fatalf("SaveMaskImageFileWithOptions() error: %v", err)
	}

	exported, err := ExportCEProjectLayerMasksFromAreaFlagsWithLayerDir(
		config,
		file,
		layerDir,
	)
	if err != nil {
		t.Fatalf("ExportCEProjectLayerMasksFromAreaFlagsWithLayerDir() error: %v", err)
	}

	zeroMask := exported["ZeroZero"]
	if zeroMask == nil {
		t.Fatal("zero/zero layer mask missing")
	}

	zeroPixel, err := zeroMask.At(2, 1)
	if err != nil {
		t.Fatalf("zeroMask.At() error: %v", err)
	}

	if zeroPixel == 0 {
		t.Fatal("expected zero/zero layer to be restored from layer directory")
	}
}

func TestExtractApplyTerritoriesFromCEProjectConfig(t *testing.T) {
	config, err := DecodeCEProjectConfig(readFixture(t, filepath.Join("positive", "ceproject-config.xml")))
	if err != nil {
		t.Fatalf("DecodeCEProjectConfig() error: %v", err)
	}

	files := ExtractTerritoryFilesFromCEProjectConfig(config)
	territory, ok := files["zombie_territories"]
	if !ok {
		t.Fatal("zombie_territories file missing")
	}

	if len(territory.Territories) != 1 || len(territory.Territories[0].Zones) != 1 {
		t.Fatalf("unexpected territory export shape")
	}

	target := &CEProjectConfigFile{}
	ApplyTerritoryFilesToCEProjectConfig(target, files)

	if target.TerritoryTypeList == nil || len(target.TerritoryTypeList.Types) != 1 {
		t.Fatalf("unexpected territory type count after apply")
	}

	if target.TerritoryTypeList.Types[0].Name != "zombie_territories" {
		t.Fatalf("applied territory type name mismatch")
	}
}

func TestApplyTerritoryFilesToCEProjectConfigAssignsColorForEmptySource(t *testing.T) {
	target := &CEProjectConfigFile{}
	ApplyTerritoryFilesToCEProjectConfig(target, map[string]*TerritoryFile{
		"custom_territories": {
			Territories: []TerritoryBlock{
				{
					Color: "",
					Zones: []TerritoryZone{
						{Name: "ZoneA", X: "100", Z: "200", R: "20"},
					},
				},
			},
		},
	})

	if target.TerritoryTypeList == nil || len(target.TerritoryTypeList.Types) != 1 {
		t.Fatalf("unexpected territory types after apply")
	}

	color := target.TerritoryTypeList.Types[0].Territories[0].Color
	if color == "" {
		t.Fatal("expected non-empty generated territory color")
	}

	expectedColor := syntheticTerritoryColor("custom_territories", 0)
	if color != expectedColor {
		t.Fatalf("generated territory color = %q, want %q", color, expectedColor)
	}
}

func TestSaveCEProjectLayerMasksToDirWithOptionsReuseBlankLayerFile(t *testing.T) {
	emptyA, err := NewMaskImage(4, 4)
	if err != nil {
		t.Fatalf("NewMaskImage(emptyA) error: %v", err)
	}

	filled, err := NewMaskImage(4, 4)
	if err != nil {
		t.Fatalf("NewMaskImage(filled) error: %v", err)
	}

	emptyB, err := NewMaskImage(4, 4)
	if err != nil {
		t.Fatalf("NewMaskImage(emptyB) error: %v", err)
	}

	if err = filled.Set(1, 1, 255); err != nil {
		t.Fatalf("filled.Set() error: %v", err)
	}

	dir := t.TempDir()
	err = SaveCEProjectLayerMasksToDirWithOptions(
		map[string]*MaskImage{
			"LayerA": emptyA,
			"LayerB": filled,
			"LayerC": emptyB,
		},
		dir,
		MaskImageEncodeOptions{
			Format:              MaskImageFormatTGA,
			ReuseBlankLayerFile: true,
			BlankLayerFileName:  "_blank_layer",
		},
	)
	if err != nil {
		t.Fatalf("SaveCEProjectLayerMasksToDirWithOptions() error: %v", err)
	}

	blankData, err := os.ReadFile(filepath.Join(dir, "_blank_layer.tga"))
	if err != nil {
		t.Fatalf("ReadFile(_blank_layer.tga) error: %v", err)
	}

	layerAData, err := os.ReadFile(filepath.Join(dir, "LayerA.tga"))
	if err != nil {
		t.Fatalf("ReadFile(LayerA.tga) error: %v", err)
	}

	layerBData, err := os.ReadFile(filepath.Join(dir, "LayerB.tga"))
	if err != nil {
		t.Fatalf("ReadFile(LayerB.tga) error: %v", err)
	}

	layerCData, err := os.ReadFile(filepath.Join(dir, "LayerC.tga"))
	if err != nil {
		t.Fatalf("ReadFile(LayerC.tga) error: %v", err)
	}

	if !reflect.DeepEqual(layerAData, blankData) {
		t.Fatal("LayerA blank mask bytes should match shared blank layer bytes")
	}

	if !reflect.DeepEqual(layerCData, blankData) {
		t.Fatal("LayerC blank mask bytes should match shared blank layer bytes")
	}

	if reflect.DeepEqual(layerBData, blankData) {
		t.Fatal("LayerB non-empty mask should not match shared blank layer bytes")
	}
}

func TestAreaFlagsLoadFileByKind(t *testing.T) {
	mask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = mask.Set(0, 0, 255); err != nil {
		t.Fatalf("mask.Set() error: %v", err)
	}

	file := &AreaFlagsMapFile{
		Header: AreaFlagsHeader{
			LayerWidth:  8,
			LayerHeight: 8,
			WorldWidth:  80,
			WorldHeight: 80,
			UsageBits:   16,
		},
		Layers: make([]AreaFlagsLayer, 20),
	}

	layerByteSize, err := file.Header.LayerByteSize()
	if err != nil {
		t.Fatalf("LayerByteSize() error: %v", err)
	}

	for index := range file.Layers {
		file.Layers[index].Bits = make([]byte, layerByteSize)
	}

	if err = file.SetLayerFromMask(0, mask); err != nil {
		t.Fatalf("SetLayerFromMask() error: %v", err)
	}

	data, err := EncodeAreaFlagsMap(file)
	if err != nil {
		t.Fatalf("EncodeAreaFlagsMap() error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "areaflags.map")
	if err = os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", path, err)
	}

	kind, value, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile(%s) error: %v", path, err)
	}

	if kind != KindAreaFlagsMap {
		t.Fatalf("kind = %q, want %q", kind, KindAreaFlagsMap)
	}

	if _, ok := value.(*AreaFlagsMapFile); !ok {
		t.Fatalf("value type = %T, want *AreaFlagsMapFile", value)
	}
}

func encodeRawMaskCodec(mask *MaskImage, _ MaskImageEncodeOptions) ([]byte, error) {
	if mask == nil {
		return nil, errors.New("nil mask")
	}

	if err := validateMaskPixelBuffer(mask); err != nil {
		return nil, err
	}

	data := make([]byte, 12+len(mask.Pixels))
	copy(data[0:4], []byte("RMK1"))
	binary.LittleEndian.PutUint32(data[4:8], mask.Width)
	binary.LittleEndian.PutUint32(data[8:12], mask.Height)
	copy(data[12:], mask.Pixels)
	return data, nil
}

func decodeRawMaskCodec(data []byte) (*MaskImage, error) {
	if len(data) < 12 {
		return nil, errors.New("payload too short")
	}

	if string(data[0:4]) != "RMK1" {
		return nil, errors.New("invalid codec marker")
	}

	width := binary.LittleEndian.Uint32(data[4:8])
	height := binary.LittleEndian.Uint32(data[8:12])
	mask, err := NewMaskImage(width, height)
	if err != nil {
		return nil, err
	}

	if len(data[12:]) != len(mask.Pixels) {
		return nil, errors.New("payload size mismatch")
	}

	copy(mask.Pixels, data[12:])
	return mask, nil
}
