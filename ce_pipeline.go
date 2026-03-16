// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// CEPipeline configures high-level CEProject forward and restore flows.
type CEPipeline struct {
	// ConfigPath is CEProject `zg-config` path (input for Forward, output for
	// Restore).
	ConfigPath string
	// LayersDir is layer mask directory (input for Forward, output for
	// Restore).
	LayersDir string
	// TerritoriesDir is `*_territories.xml` directory
	// (output for Forward, input for Restore).
	TerritoriesDir string
	// AreaFlagsMapPath is `areaflags.map` path
	// (output for Forward, input for Restore).
	AreaFlagsMapPath string
}

// CEForwardResult reports forward pipeline output stats.
type CEForwardResult struct {
	// LayerMaskCount is number of CEProject layers in config.
	LayerMaskCount int
	// TerritoryFileCount is number of written `*_territories.xml` files.
	TerritoryFileCount int
	// AreaFlagsMapBytes is encoded `areaflags.map` size in bytes.
	AreaFlagsMapBytes int
}

// CERestoreOptions controls reverse restoration behavior.
type CERestoreOptions struct {
	// TemplateConfigPath is optional CEProject config template path used to
	// preserve custom layer names/flags and UI metadata.
	// If empty, Restore tries existing pipeline ConfigPath, then generates
	// synthetic config from map bit planes.
	TemplateConfigPath string
	// LayerImageOptions controls output mask image format and options.
	LayerImageOptions MaskImageEncodeOptions
	// RequireTerritories makes restore fail when TerritoriesDir is missing.
	RequireTerritories bool
}

// CERestoreResult reports reverse pipeline output stats.
type CERestoreResult struct {
	// LayerMaskCount is number of exported layer mask files.
	LayerMaskCount int
	// TerritoryFileCount is number of loaded `*_territories.xml` files.
	TerritoryFileCount int
	// AreaFlagsMapBytes is source `areaflags.map` size in bytes.
	AreaFlagsMapBytes int
}

// Forward builds `areaflags.map` and `*_territories.xml` outputs from CEProject
// `zg-config` and source layer masks.
func (pipeline CEPipeline) Forward() (*CEForwardResult, error) {
	if err := pipeline.validateForward(); err != nil {
		return nil, err
	}

	configData, err := os.ReadFile(pipeline.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("read CEProject config %q: %w", pipeline.ConfigPath, err)
	}

	config, err := DecodeCEProjectConfig(configData)
	if err != nil {
		return nil, fmt.Errorf("decode CEProject config %q: %w", pipeline.ConfigPath, err)
	}

	mapFile, err := BuildAreaFlagsFromCEProjectLayerDir(config, pipeline.LayersDir)
	if err != nil {
		return nil, err
	}

	encodedMap, err := EncodeAreaFlagsMap(mapFile)
	if err != nil {
		return nil, err
	}

	if err = ensureParentDirForFile(pipeline.AreaFlagsMapPath); err != nil {
		return nil, err
	}

	if err = writeFile600(pipeline.AreaFlagsMapPath, encodedMap); err != nil {
		return nil, fmt.Errorf("write area flags %q: %w", pipeline.AreaFlagsMapPath, err)
	}

	territoryFiles := ExtractTerritoryFilesFromCEProjectConfig(config)
	if pipeline.TerritoriesDir != "" {
		if err = SaveTerritoryFilesToDir(territoryFiles, pipeline.TerritoriesDir); err != nil {
			return nil, err
		}
	}

	layerMaskCount := 0
	if config != nil && config.Layers != nil {
		layerMaskCount = len(config.Layers.Layers)
	}

	return &CEForwardResult{
		LayerMaskCount:     layerMaskCount,
		TerritoryFileCount: len(territoryFiles),
		AreaFlagsMapBytes:  len(encodedMap),
	}, nil
}

// Restore reconstructs CEProject `zg-config` and layer masks from
// `areaflags.map` and optional `*_territories.xml` files.
func (pipeline CEPipeline) Restore(
	options CERestoreOptions,
) (*CERestoreResult, error) {
	if err := pipeline.validateRestore(); err != nil {
		return nil, err
	}

	mapData, err := os.ReadFile(pipeline.AreaFlagsMapPath)
	if err != nil {
		return nil, fmt.Errorf("read area flags %q: %w", pipeline.AreaFlagsMapPath, err)
	}

	mapFile, err := DecodeAreaFlagsMap(mapData)
	if err != nil {
		return nil, err
	}

	config, err := loadRestoreConfigTemplate(
		mapFile,
		options.TemplateConfigPath,
		pipeline.ConfigPath,
	)
	if err != nil {
		return nil, err
	}

	territoryFiles, err := loadTerritoryFilesForRestore(
		pipeline.TerritoriesDir,
		options.RequireTerritories,
	)
	if err != nil {
		return nil, err
	}

	ApplyTerritoryFilesToCEProjectConfig(config, territoryFiles)

	masks, err := ExportCEProjectLayerMasksFromAreaFlags(config, mapFile)
	if err != nil {
		return nil, err
	}

	if err = SaveCEProjectLayerMasksToDirWithOptions(
		masks,
		pipeline.LayersDir,
		options.LayerImageOptions,
	); err != nil {
		return nil, err
	}

	encodedConfig, err := EncodeCEProjectConfig(config)
	if err != nil {
		return nil, err
	}

	if err = ensureParentDirForFile(pipeline.ConfigPath); err != nil {
		return nil, err
	}

	if err = writeFile600(pipeline.ConfigPath, encodedConfig); err != nil {
		return nil, fmt.Errorf("write CEProject config %q: %w", pipeline.ConfigPath, err)
	}

	return &CERestoreResult{
		LayerMaskCount:     len(masks),
		TerritoryFileCount: len(territoryFiles),
		AreaFlagsMapBytes:  len(mapData),
	}, nil
}

// LoadTerritoryFilesFromDir loads all `*.xml` files as `*_territories.xml`
// payloads keyed by basename without extension.
func LoadTerritoryFilesFromDir(dir string) (map[string]*TerritoryFile, error) {
	output := map[string]*TerritoryFile{}
	if strings.TrimSpace(dir) == "" {
		return output, ErrEmptyTerritoryDirPath
	}

	names, err := readDirNamesSorted(dir)
	if err != nil {
		return nil, fmt.Errorf("read territory directory %q: %w", dir, err)
	}

	for _, fileName := range names {
		path := filepath.Join(dir, fileName)
		info, statErr := os.Lstat(path)
		if statErr != nil {
			return nil, fmt.Errorf("stat territory path %q: %w", path, statErr)
		}

		if info.IsDir() {
			continue
		}

		if !strings.EqualFold(filepath.Ext(fileName), ".xml") {
			continue
		}

		base := strings.TrimSpace(strings.TrimSuffix(fileName, filepath.Ext(fileName)))
		if base == "" {
			continue
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil, fmt.Errorf("read territory file %q: %w", path, readErr)
		}

		file, decodeErr := DecodeTerritory(data)
		if decodeErr != nil {
			return nil, fmt.Errorf("decode territory file %q: %w", path, decodeErr)
		}

		output[base] = file
	}

	return output, nil
}

// SaveTerritoryFilesToDir writes territory payloads into
// `<name>.xml` files in deterministic order.
func SaveTerritoryFilesToDir(
	files map[string]*TerritoryFile,
	dir string,
) error {
	if strings.TrimSpace(dir) == "" {
		return ErrEmptyTerritoryDirPath
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("mkdir %q: %w", dir, err)
	}

	keys := make([]string, 0, len(files))
	for name := range files {
		keys = append(keys, name)
	}

	sort.Strings(keys)
	for _, key := range keys {
		file := files[key]
		if file == nil {
			continue
		}

		fileName := normalizeTerritoryFileName(key)
		if fileName == "" {
			continue
		}

		encoded, err := EncodeTerritory(file)
		if err != nil {
			return err
		}

		path := filepath.Join(dir, fileName)
		if err = writeFile600(path, encoded); err != nil {
			return fmt.Errorf("write territory file %q: %w", path, err)
		}
	}

	return nil
}

// validateForward validates required paths for forward flow.
func (pipeline CEPipeline) validateForward() error {
	if strings.TrimSpace(pipeline.ConfigPath) == "" {
		return errors.New("empty CEProject config path")
	}

	if strings.TrimSpace(pipeline.LayersDir) == "" {
		return errors.New("empty CEProject layers directory path")
	}

	if strings.TrimSpace(pipeline.AreaFlagsMapPath) == "" {
		return errors.New("empty area flags output path")
	}

	return nil
}

// validateRestore validates required paths for reverse flow.
func (pipeline CEPipeline) validateRestore() error {
	if strings.TrimSpace(pipeline.AreaFlagsMapPath) == "" {
		return errors.New("empty area flags input path")
	}

	if strings.TrimSpace(pipeline.ConfigPath) == "" {
		return errors.New("empty CEProject config output path")
	}

	if strings.TrimSpace(pipeline.LayersDir) == "" {
		return errors.New("empty CEProject layers output directory path")
	}

	return nil
}

// ensureParentDirForFile creates parent directory for file path.
func ensureParentDirForFile(path string) error {
	parent := filepath.Dir(path)
	if parent == "." || parent == "" {
		return nil
	}

	if err := os.MkdirAll(parent, 0o700); err != nil {
		return fmt.Errorf("mkdir %q: %w", parent, err)
	}

	return nil
}

// normalizeTerritoryFileName normalizes territory output file name.
func normalizeTerritoryFileName(name string) string {
	base := strings.TrimSpace(name)
	if base == "" {
		return ""
	}

	if strings.EqualFold(filepath.Ext(base), ".xml") {
		return base
	}

	return base + ".xml"
}

// loadRestoreConfigTemplate resolves config template for reverse flow.
func loadRestoreConfigTemplate(
	mapFile *AreaFlagsMapFile,
	templatePath string,
	outputPath string,
) (*CEProjectConfigFile, error) {
	path := strings.TrimSpace(templatePath)
	if path == "" {
		path = strings.TrimSpace(outputPath)
	}

	if path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			config, decodeErr := DecodeCEProjectConfig(data)
			if decodeErr != nil {
				return nil, fmt.Errorf("decode restore template %q: %w", path, decodeErr)
			}

			return config, nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("read restore template %q: %w", path, err)
		}
	}

	return newSyntheticCEProjectConfigFromAreaFlags(mapFile)
}

// loadTerritoryFilesForRestore loads territory directory for reverse flow.
func loadTerritoryFilesForRestore(
	dir string,
	required bool,
) (map[string]*TerritoryFile, error) {
	if strings.TrimSpace(dir) == "" {
		if required {
			return nil, ErrEmptyTerritoryDirPath
		}

		return map[string]*TerritoryFile{}, nil
	}

	files, err := LoadTerritoryFilesFromDir(dir)
	if err == nil {
		return files, nil
	}

	if errors.Is(err, os.ErrNotExist) && !required {
		return map[string]*TerritoryFile{}, nil
	}

	return nil, err
}

// newSyntheticCEProjectConfigFromAreaFlags builds minimal CEProject config from map.
func newSyntheticCEProjectConfigFromAreaFlags(
	mapFile *AreaFlagsMapFile,
) (*CEProjectConfigFile, error) {
	if mapFile == nil {
		return nil, ErrNilAreaFlagsFile
	}

	valueBits, err := mapFile.ValueBits()
	if err != nil {
		return nil, err
	}

	config := &CEProjectConfigFile{
		Selected: &CEProjectSelected{
			Row: "0",
		},
		Global: &CEProjectGlobal{
			Background: &CEProjectBackground{
				File: "map.png",
				RGBA: "16777215",
			},
			Layer: &CEProjectDimension{
				Size: strconv.FormatUint(uint64(mapFile.Header.LayerWidth), 10),
			},
			World: &CEProjectDimension{
				Size: strconv.FormatUint(uint64(mapFile.Header.WorldWidth), 10),
			},
		},
		Areas: &CEProjectAreas{
			Usages: &CEProjectAreaUsageList{
				Items: make([]NamedRef, 0, mapFile.Header.UsageBits),
			},
			Values: &CEProjectAreaValueList{
				Items: make([]NamedRef, 0, valueBits),
			},
		},
		Layers: &CEProjectLayers{
			Layers: make([]CEProjectLayer, 0, mapFile.Header.UsageBits+valueBits),
		},
		TerritoryTypeList: &CEProjectTerritoryTypeList{},
	}

	for bit := uint32(0); bit < mapFile.Header.UsageBits; bit++ {
		name := fmt.Sprintf("usage_bit_%02d", bit)
		mask := uint32(1) << bit
		config.Areas.Usages.Items = append(config.Areas.Usages.Items, NamedRef{
			Name: name,
		})
		config.Layers.Layers = append(config.Layers.Layers, CEProjectLayer{
			Name:       name,
			UsageFlags: strconv.FormatUint(uint64(mask), 10),
			ValueFlags: "0",
			Color:      syntheticLayerColor(name, mask, 0),
			Visible:    "1",
		})
	}

	for bit := range valueBits {
		name := fmt.Sprintf("value_bit_%02d", bit)
		mask := uint32(1) << bit
		config.Areas.Values.Items = append(config.Areas.Values.Items, NamedRef{
			Name: name,
		})
		config.Layers.Layers = append(config.Layers.Layers, CEProjectLayer{
			Name:       name,
			UsageFlags: "0",
			ValueFlags: strconv.FormatUint(uint64(mask), 10),
			Color:      syntheticLayerColor(name, 0, mask),
			Visible:    "1",
		})
	}

	return config, nil
}
