// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"errors"
	"fmt"
	"math/bits"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	// ceProjectUsageRasterXOffset aligns usage mask pixels with areaflags raster.
	ceProjectUsageRasterXOffset = 0
	// ceProjectValueOnlyRasterXOffset aligns value-only mask pixels.
	ceProjectValueOnlyRasterXOffset = 0
)

// BuildAreaFlagsFromCEProjectLayerDir builds `areaflags.map` from CEProject config
// and `layers/*.tga` masks (mapped by `<layer name>.tga`).
func BuildAreaFlagsFromCEProjectLayerDir(
	config *CEProjectConfigFile,
	layerDir string,
) (*AreaFlagsMapFile, error) {
	if config == nil {
		return nil, ErrNilCEProjectConfig
	}

	if config.Layers == nil {
		return nil, ErrMissingCEProjectLayers
	}

	masks, err := loadCEProjectLayerMasksFromDir(config, layerDir)
	if err != nil {
		return nil, err
	}

	return BuildAreaFlagsFromCEProjectMasks(config, masks)
}

// BuildAreaFlagsFromCEProjectMasks builds `areaflags.map` from CEProject config and
// already loaded masks keyed by layer name.
func BuildAreaFlagsFromCEProjectMasks(
	config *CEProjectConfigFile,
	masks map[string]*MaskImage,
) (*AreaFlagsMapFile, error) {
	if config == nil {
		return nil, ErrNilCEProjectConfig
	}

	if config.Global == nil || config.Global.Layer == nil || config.Global.World == nil {
		return nil, errors.New("missing CEProject global settings")
	}

	if config.Layers == nil {
		return nil, ErrMissingCEProjectLayers
	}

	layerSize, err := parseUint32Text(config.Global.Layer.Size)
	if err != nil {
		return nil, fmt.Errorf("parse global layer size: %w", err)
	}

	worldSize, err := parseUint32Text(config.Global.World.Size)
	if err != nil {
		return nil, fmt.Errorf("parse global world size: %w", err)
	}

	usageBits := inferUsageBitCount(config)
	valueBits := inferValueBitCount(config)
	layerByteSize, err := (AreaFlagsHeader{
		LayerWidth:  layerSize,
		LayerHeight: layerSize,
	}).LayerByteSize()
	if err != nil {
		return nil, err
	}

	layerCount := int(usageBits + valueBits)
	file := &AreaFlagsMapFile{
		Header: AreaFlagsHeader{
			LayerWidth:  layerSize,
			LayerHeight: layerSize,
			WorldWidth:  worldSize,
			WorldHeight: worldSize,
			UsageBits:   usageBits,
			Reserved:    0,
		},
		Layers: make([]AreaFlagsLayer, layerCount),
	}

	for index := range file.Layers {
		file.Layers[index].Bits = make([]byte, layerByteSize)
	}

	for _, layer := range config.Layers.Layers {
		name := strings.TrimSpace(layer.Name)

		usageFlags, parseErr := parseUint32Text(layer.UsageFlags)
		if parseErr != nil {
			return nil, fmt.Errorf("parse usage_flags for layer %q: %w", name, parseErr)
		}

		valueFlags, parseErr := parseUint32Text(layer.ValueFlags)
		if parseErr != nil {
			return nil, fmt.Errorf("parse value_flags for layer %q: %w", name, parseErr)
		}

		if usageBits < 32 && usageFlags>>usageBits != 0 {
			return nil, fmt.Errorf(
				"usage_flags %d exceed usage bit count %d for layer %q",
				usageFlags,
				usageBits,
				name,
			)
		}

		mask, ok := masks[name]
		if !ok {
			continue
		}

		if mask == nil {
			return nil, fmt.Errorf("nil mask for layer %q", name)
		}

		if mask.Width != layerSize || mask.Height != layerSize {
			return nil, fmt.Errorf(
				"mask size mismatch for layer %q: %dx%d, want %dx%d",
				name,
				mask.Width,
				mask.Height,
				layerSize,
				layerSize,
			)
		}

		if valueBits < 32 && valueFlags>>valueBits != 0 {
			return nil, fmt.Errorf(
				"value_flags %d exceed value bit count %d for layer %q",
				valueFlags,
				valueBits,
				name,
			)
		}

		usageOffset, valueOffset := ceProjectRasterOffsets(usageFlags, valueFlags)
		for y := range layerSize {
			for x := range layerSize {
				pixel, readErr := mask.At(x, y)
				if readErr != nil {
					return nil, readErr
				}

				if pixel == 0 {
					continue
				}

				mapXUsage, usageOK := ceProjectMaskXToAreaFlagsX(
					x,
					layerSize,
					usageOffset,
				)
				mapXValue, valueOK := ceProjectMaskXToAreaFlagsX(
					x,
					layerSize,
					valueOffset,
				)

				mapY := layerSize - 1 - y
				for bit := range usageBits {
					if usageFlags&(uint32(1)<<bit) == 0 {
						continue
					}

					if !usageOK {
						continue
					}

					if setErr := setPackedBit(
						file.Layers[bit].Bits,
						layerSize,
						mapXUsage,
						mapY,
					); setErr != nil {
						return nil, setErr
					}
				}

				for bit := range valueBits {
					if valueFlags&(uint32(1)<<bit) == 0 {
						continue
					}

					if !valueOK {
						continue
					}

					index := usageBits + bit
					if setErr := setPackedBit(
						file.Layers[index].Bits,
						layerSize,
						mapXValue,
						mapY,
					); setErr != nil {
						return nil, setErr
					}
				}
			}
		}
	}

	return file, nil
}

// ExportAreaFlagsBitMasksToTGADir exports all usage/value bit-planes into
// individual `.tga` files.
func ExportAreaFlagsBitMasksToTGADir(
	file *AreaFlagsMapFile,
	dir string,
) error {
	if file == nil {
		return ErrNilAreaFlagsFile
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("mkdir %q: %w", dir, err)
	}

	valueBits, err := file.ValueBits()
	if err != nil {
		return err
	}

	for bit := uint32(0); bit < file.Header.UsageBits; bit++ {
		mask, maskErr := file.MaskFromLayer(int(bit))
		if maskErr != nil {
			return maskErr
		}

		name := fmt.Sprintf("usage_bit_%02d.tga", bit)
		if writeErr := SaveTGAMaskFile(filepath.Join(dir, name), mask); writeErr != nil {
			return writeErr
		}
	}

	for bit := range valueBits {
		mask, maskErr := file.MaskFromLayer(int(file.Header.UsageBits + bit))
		if maskErr != nil {
			return maskErr
		}

		name := fmt.Sprintf("value_bit_%02d.tga", bit)
		if writeErr := SaveTGAMaskFile(filepath.Join(dir, name), mask); writeErr != nil {
			return writeErr
		}
	}

	return nil
}

// ExportCEProjectLayerMasksFromAreaFlags projects bit-planes back into masks by
// layer `usage_flags`/`value_flags` definitions from CEProject config.
func ExportCEProjectLayerMasksFromAreaFlags(
	config *CEProjectConfigFile,
	file *AreaFlagsMapFile,
) (map[string]*MaskImage, error) {
	if config == nil {
		return nil, ErrNilCEProjectConfig
	}

	if config.Layers == nil {
		return nil, ErrMissingCEProjectLayers
	}

	if file == nil {
		return nil, ErrNilAreaFlagsFile
	}

	valueBits, err := file.ValueBits()
	if err != nil {
		return nil, err
	}

	output := make(map[string]*MaskImage, len(config.Layers.Layers))
	for _, layer := range config.Layers.Layers {
		name := strings.TrimSpace(layer.Name)
		if name == "" {
			return nil, ErrEmptyLayerName
		}

		usageFlags, parseErr := parseUint32Text(layer.UsageFlags)
		if parseErr != nil {
			return nil, fmt.Errorf("parse usage_flags for layer %q: %w", name, parseErr)
		}

		valueFlags, parseErr := parseUint32Text(layer.ValueFlags)
		if parseErr != nil {
			return nil, fmt.Errorf("parse value_flags for layer %q: %w", name, parseErr)
		}

		mask, maskErr := NewMaskImage(file.Header.LayerWidth, file.Header.LayerHeight)
		if maskErr != nil {
			return nil, maskErr
		}

		// CEProject does not store zero/zero layers in areaflags payload.
		// Keep them empty on projection.
		if usageFlags == 0 && valueFlags == 0 {
			output[name] = mask
			continue
		}

		usageOffset, valueOffset := ceProjectRasterOffsets(usageFlags, valueFlags)
		for y := range file.Header.LayerHeight {
			for x := range file.Header.LayerWidth {
				usageOK := usageFlags == 0
				valueOK := valueFlags == 0

				mapXUsage, usageXOK := ceProjectMaskXToAreaFlagsX(
					x,
					file.Header.LayerWidth,
					usageOffset,
				)
				mapXValue, valueXOK := ceProjectMaskXToAreaFlagsX(
					x,
					file.Header.LayerWidth,
					valueOffset,
				)

				mapY := file.Header.LayerHeight - 1 - y
				if usageFlags != 0 {
					if !usageXOK {
						usageOK = false
					} else {
						for bit := range file.Header.UsageBits {
							if usageFlags&(uint32(1)<<bit) == 0 {
								continue
							}

							set, readErr := getPackedBit(
								file.Layers[bit].Bits,
								file.Header.LayerWidth,
								mapXUsage,
								mapY,
							)
							if readErr != nil {
								return nil, readErr
							}

							if !set {
								usageOK = false
								break
							}

							usageOK = true
						}
					}
				}

				if valueFlags != 0 {
					if !valueXOK {
						valueOK = false
					} else {
						for bit := range valueBits {
							if valueFlags&(uint32(1)<<bit) == 0 {
								continue
							}

							index := file.Header.UsageBits + bit
							set, readErr := getPackedBit(
								file.Layers[index].Bits,
								file.Header.LayerWidth,
								mapXValue,
								mapY,
							)
							if readErr != nil {
								return nil, readErr
							}

							if !set {
								valueOK = false
								break
							}

							valueOK = true
						}
					}
				}

				if usageOK && valueOK {
					if setErr := mask.Set(x, y, 255); setErr != nil {
						return nil, setErr
					}
				}
			}
		}

		output[name] = mask
	}

	return output, nil
}

// ExportCEProjectLayerMasksFromAreaFlagsWithLayerDir projects masks from
// `areaflags.map` and restores zero/zero layers from CEProject source layer
// directory when available.
func ExportCEProjectLayerMasksFromAreaFlagsWithLayerDir(
	config *CEProjectConfigFile,
	file *AreaFlagsMapFile,
	layerDir string,
) (map[string]*MaskImage, error) {
	output, err := ExportCEProjectLayerMasksFromAreaFlags(config, file)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(layerDir) == "" {
		return output, nil
	}

	maskByLayer, err := loadCEProjectLayerMasksFromDir(config, layerDir)
	if err != nil {
		return nil, err
	}

	for _, layer := range config.Layers.Layers {
		name := strings.TrimSpace(layer.Name)
		if name == "" {
			return nil, ErrEmptyLayerName
		}

		usageFlags, parseErr := parseUint32Text(layer.UsageFlags)
		if parseErr != nil {
			return nil, fmt.Errorf("parse usage_flags for layer %q: %w", name, parseErr)
		}

		valueFlags, parseErr := parseUint32Text(layer.ValueFlags)
		if parseErr != nil {
			return nil, fmt.Errorf("parse value_flags for layer %q: %w", name, parseErr)
		}

		if usageFlags != 0 || valueFlags != 0 {
			continue
		}

		mask, ok := maskByLayer[name]
		if !ok || mask == nil {
			continue
		}

		if mask.Width != file.Header.LayerWidth || mask.Height != file.Header.LayerHeight {
			return nil, fmt.Errorf(
				"zero/zero layer mask size mismatch for layer %q: %dx%d, want %dx%d",
				name,
				mask.Width,
				mask.Height,
				file.Header.LayerWidth,
				file.Header.LayerHeight,
			)
		}

		output[name] = mask
	}

	return output, nil
}

// SaveCEProjectLayerMasksToDir writes projected CEProject masks to directory.
func SaveCEProjectLayerMasksToDir(
	masks map[string]*MaskImage,
	dir string,
) error {
	return SaveCEProjectLayerMasksToDirWithOptions(
		masks,
		dir,
		MaskImageEncodeOptions{Format: MaskImageFormatTGA},
	)
}

// SaveCEProjectLayerMasksToDirWithOptions writes projected CEProject masks to
// directory using selected image format.
func SaveCEProjectLayerMasksToDirWithOptions(
	masks map[string]*MaskImage,
	dir string,
	options MaskImageEncodeOptions,
) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("mkdir %q: %w", dir, err)
	}

	names := make([]string, 0, len(masks))
	for name := range masks {
		names = append(names, name)
	}

	sort.Strings(names)

	format := normalizeMaskImageFormat(options.Format)
	if format == "" {
		format = MaskImageFormatTGA
	}

	extension := "." + format

	options.Format = format
	if !options.ReuseBlankLayerFile {
		for _, name := range names {
			path := filepath.Join(dir, name+extension)
			if err := SaveMaskImageFileWithOptions(path, masks[name], options); err != nil {
				return err
			}
		}

		return nil
	}

	blankByName := make(map[string]bool, len(masks))
	for _, name := range names {
		isBlank, err := isMaskImageBlank(masks[name])
		if err != nil {
			return err
		}

		blankByName[name] = isBlank
	}

	blankLayerName := uniqueBlankLayerBaseName(
		names,
		options.BlankLayerFileName,
	)
	blankPath := filepath.Join(dir, blankLayerName+extension)
	blankSourceName := ""
	for _, name := range names {
		if blankByName[name] {
			blankSourceName = name
			break
		}
	}

	if blankSourceName != "" {
		if err := SaveMaskImageFileWithOptions(
			blankPath,
			masks[blankSourceName],
			options,
		); err != nil {
			return err
		}
	}

	for _, name := range names {
		path := filepath.Join(dir, name+extension)
		if blankSourceName != "" && blankByName[name] {
			if err := linkOrCopyFile(blankPath, path); err != nil {
				return err
			}

			continue
		}

		if err := SaveMaskImageFileWithOptions(path, masks[name], options); err != nil {
			return err
		}
	}

	return nil
}

// ExtractTerritoryFilesFromCEProjectConfig converts CEProject territory list to
// `*_territories.xml` payloads keyed by territory type name.
func ExtractTerritoryFilesFromCEProjectConfig(
	config *CEProjectConfigFile,
) map[string]*TerritoryFile {
	output := map[string]*TerritoryFile{}
	if config == nil || config.TerritoryTypeList == nil {
		return output
	}

	for _, item := range config.TerritoryTypeList.Types {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}

		file := &TerritoryFile{}
		for _, territory := range item.Territories {
			block := TerritoryBlock{
				Color: strings.TrimSpace(territory.Color),
			}

			for _, zone := range territory.Zones {
				block.Zones = append(block.Zones, TerritoryZone{
					Name: zone.Name,
					SMin: zone.SMin,
					SMax: zone.SMax,
					DMin: zone.DMin,
					DMax: zone.DMax,
					X:    zone.X,
					Z:    zone.Z,
					R:    zone.R,
				})
			}

			file.Territories = append(file.Territories, block)
		}

		output[name] = file
	}

	return output
}

// ApplyTerritoryFilesToCEProjectConfig writes `*_territories.xml` payloads into
// CEProject territory list.
func ApplyTerritoryFilesToCEProjectConfig(
	config *CEProjectConfigFile,
	files map[string]*TerritoryFile,
) {
	if config == nil {
		return
	}

	if config.TerritoryTypeList == nil {
		config.TerritoryTypeList = &CEProjectTerritoryTypeList{}
	}

	keys := make([]string, 0, len(files))
	for name := range files {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	config.TerritoryTypeList.Types = config.TerritoryTypeList.Types[:0]
	for _, key := range keys {
		name := strings.TrimSuffix(strings.TrimSpace(key), ".xml")
		if name == "" {
			continue
		}

		source := files[key]
		if source == nil {
			continue
		}

		item := CEProjectTerritoryType{Name: name}
		for territoryIndex, block := range source.Territories {
			color := strings.TrimSpace(block.Color)
			if color == "" {
				color = syntheticTerritoryColor(name, territoryIndex)
			}

			territory := CEProjectTerritory{
				Name:    inferCEProjectTerritoryName(block, territoryIndex),
				Visible: "0",
				Color:   color,
			}

			for _, zone := range block.Zones {
				territory.Zones = append(territory.Zones, CEProjectTerritoryZone{
					Name: zone.Name,
					SMin: zone.SMin,
					SMax: zone.SMax,
					DMin: zone.DMin,
					DMax: zone.DMax,
					X:    zone.X,
					Z:    zone.Z,
					R:    zone.R,
				})
			}

			item.Territories = append(item.Territories, territory)
		}

		config.TerritoryTypeList.Types = append(config.TerritoryTypeList.Types, item)
	}
}

// inferUsageBitCount infers usage bit-plane count from usages list length.
func inferUsageBitCount(config *CEProjectConfigFile) uint32 {
	count := 0
	if config != nil && config.Areas != nil && config.Areas.Usages != nil {
		count = len(config.Areas.Usages.Items)
	}

	if count > 16 {
		return 32
	}

	return 16
}

// inferValueBitCount infers value bit-plane count from values list length.
func inferValueBitCount(config *CEProjectConfigFile) uint32 {
	count := 0
	if config != nil && config.Areas != nil && config.Areas.Values != nil {
		count = len(config.Areas.Values.Items)
	}

	if count <= 0 {
		return 4
	}

	// value planes are grouped in nibble/power-of-two units (4, 8, 16...).
	value := uint32(1 << bits.Len(uint(count-1)))
	if value < 4 {
		return 4
	}

	return value
}

// inferCEProjectTerritoryName picks a stable territory name when not preserved.
func inferCEProjectTerritoryName(block TerritoryBlock, index int) string {
	if len(block.Zones) > 0 && strings.TrimSpace(block.Zones[0].Name) != "" {
		return block.Zones[0].Name
	}

	return fmt.Sprintf("Territory_%02d", index+1)
}

// parseUint32Text parses uint32 text and returns zero for empty input.
func parseUint32Text(value string) (uint32, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}

	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(parsed), nil
}

// isMaskImageBlank reports whether mask has no non-zero pixels.
func isMaskImageBlank(mask *MaskImage) (bool, error) {
	if mask == nil {
		return false, fmt.Errorf("%w: nil", ErrUnsupportedValue)
	}

	if err := validateMaskPixelBuffer(mask); err != nil {
		return false, err
	}

	for _, pixel := range mask.Pixels {
		if pixel != 0 {
			return false, nil
		}
	}

	return true, nil
}

// uniqueBlankLayerBaseName returns non-conflicting blank file base name.
func uniqueBlankLayerBaseName(existing []string, requested string) string {
	name := strings.TrimSpace(requested)
	if name == "" {
		name = "_blank"
	}

	existingSet := make(map[string]struct{}, len(existing))
	for _, item := range existing {
		existingSet[strings.ToLower(strings.TrimSpace(item))] = struct{}{}
	}

	base := name
	for suffix := 0; ; suffix++ {
		candidate := base
		if suffix > 0 {
			candidate = fmt.Sprintf("%s_%d", base, suffix)
		}

		if _, exists := existingSet[strings.ToLower(candidate)]; !exists {
			return candidate
		}
	}
}

// linkOrCopyFile links destination to source, copying when hard links are unavailable.
func linkOrCopyFile(source string, destination string) error {
	if source == destination {
		return nil
	}

	if removeErr := os.Remove(destination); removeErr != nil &&
		!errors.Is(removeErr, os.ErrNotExist) {
		return fmt.Errorf("remove %q before link: %w", destination, removeErr)
	}

	if err := os.Link(source, destination); err == nil {
		return nil
	}

	data, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("read blank layer %q: %w", source, err)
	}

	if err = writeFile600(destination, data); err != nil {
		return fmt.Errorf("write blank layer %q: %w", destination, err)
	}

	return nil
}

// loadCEProjectLayerMasksFromDir loads mask files mapped by CEProject layer names.
// Mask file extension is arbitrary and resolved by registered image codecs.
func loadCEProjectLayerMasksFromDir(
	config *CEProjectConfigFile,
	layerDir string,
) (map[string]*MaskImage, error) {
	layerFiles, err := indexMaskFilesByBaseName(layerDir)
	if err != nil {
		return nil, err
	}

	output := make(map[string]*MaskImage, len(config.Layers.Layers))
	for _, layer := range config.Layers.Layers {
		name := strings.TrimSpace(layer.Name)
		if name == "" {
			return nil, ErrEmptyLayerName
		}

		path, ok := layerFiles[strings.ToLower(name)]
		if !ok {
			// Missing layer file is treated as empty mask.
			continue
		}

		mask, _, loadErr := LoadMaskImageFile(path)
		if loadErr != nil {
			return nil, fmt.Errorf("load layer mask %q from %q: %w", name, path, loadErr)
		}

		output[name] = mask
	}

	return output, nil
}

// indexMaskFilesByBaseName indexes files in directory by lowercase basename.
// If multiple extensions exist for one layer, `.tga` takes precedence.
func indexMaskFilesByBaseName(dir string) (map[string]string, error) {
	names, err := readDirNamesSorted(dir)
	if err != nil {
		return nil, fmt.Errorf("read CEProject layer directory %q: %w", dir, err)
	}

	output := make(map[string]string, len(names))
	for _, name := range names {
		path := filepath.Join(dir, name)
		info, statErr := os.Lstat(path)
		if statErr != nil {
			return nil, fmt.Errorf("stat CEProject layer path %q: %w", path, statErr)
		}

		if info.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(name))
		if ext == "" {
			continue
		}

		base := strings.TrimSuffix(name, ext)
		base = strings.TrimSpace(base)
		if base == "" {
			continue
		}

		key := strings.ToLower(base)
		currentPath, exists := output[key]
		if !exists {
			output[key] = path
			continue
		}

		currentExt := strings.ToLower(filepath.Ext(currentPath))
		if currentExt != ".tga" && ext == ".tga" {
			output[key] = path
		}
	}

	return output, nil
}

// ceProjectRasterOffsets returns X offsets for usage and value sections.
func ceProjectRasterOffsets(usageFlags uint32, valueFlags uint32) (int, int) {
	usageOffset := ceProjectUsageRasterXOffset
	valueOffset := ceProjectUsageRasterXOffset

	// Value-only layers align directly with value section coordinates.
	if usageFlags == 0 && valueFlags != 0 {
		valueOffset = ceProjectValueOnlyRasterXOffset
	}

	return usageOffset, valueOffset
}

// ceProjectMaskXToAreaFlagsX maps CEProject mask X coordinate to areaflags X.
func ceProjectMaskXToAreaFlagsX(maskX uint32, width uint32, offset int) (uint32, bool) {
	const maxUint32Int64 = int64(^uint32(0))

	if offset < 0 {
		shift64 := -int64(offset)
		if shift64 > maxUint32Int64 {
			return 0, false
		}

		shift := uint32(shift64)
		if maskX < shift {
			return 0, false
		}

		mapped := maskX - shift
		if mapped >= width {
			return 0, false
		}

		return mapped, true
	}

	shift64 := int64(offset)
	if shift64 > maxUint32Int64 {
		return 0, false
	}

	shift := uint32(shift64)
	if maskX > ^uint32(0)-shift {
		return 0, false
	}

	mapped := maskX + shift
	if mapped >= width {
		return 0, false
	}

	return mapped, true
}
