// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCEPipelineForward(t *testing.T) {
	configData := readFixture(t, filepath.Join("positive", "ceproject-config.xml"))
	configPath := filepath.Join(t.TempDir(), "project.xml")
	if err := os.WriteFile(configPath, configData, 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", configPath, err)
	}

	layerDir := filepath.Join(t.TempDir(), "layers")
	if err := os.MkdirAll(layerDir, 0o700); err != nil {
		t.Fatalf("MkdirAll(%s) error: %v", layerDir, err)
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

	if err = SaveTGAMaskFile(filepath.Join(layerDir, "Military.tga"), militaryMask); err != nil {
		t.Fatalf("SaveTGAMaskFile(Military.tga) error: %v", err)
	}

	if err = SaveTGAMaskFile(filepath.Join(layerDir, "Police.tga"), policeMask); err != nil {
		t.Fatalf("SaveTGAMaskFile(Police.tga) error: %v", err)
	}

	mapPath := filepath.Join(t.TempDir(), "out", "areaflags.map")
	territoryDir := filepath.Join(t.TempDir(), "out", "territories")
	pipeline := CEPipeline{
		ConfigPath:       configPath,
		LayersDir:        layerDir,
		TerritoriesDir:   territoryDir,
		AreaFlagsMapPath: mapPath,
	}

	result, err := pipeline.Forward()
	if err != nil {
		t.Fatalf("pipeline.Forward() error: %v", err)
	}

	if result.LayerMaskCount != 2 {
		t.Fatalf("LayerMaskCount = %d, want 2", result.LayerMaskCount)
	}

	if result.TerritoryFileCount != 1 {
		t.Fatalf("TerritoryFileCount = %d, want 1", result.TerritoryFileCount)
	}

	mapData, err := os.ReadFile(mapPath)
	if err != nil {
		t.Fatalf("ReadFile(%s) error: %v", mapPath, err)
	}

	mapFile, err := DecodeAreaFlagsMap(mapData)
	if err != nil {
		t.Fatalf("DecodeAreaFlagsMap() error: %v", err)
	}

	config, err := DecodeCEProjectConfig(configData)
	if err != nil {
		t.Fatalf("DecodeCEProjectConfig() error: %v", err)
	}

	exported, err := ExportCEProjectLayerMasksFromAreaFlags(config, mapFile)
	if err != nil {
		t.Fatalf("ExportCEProjectLayerMasksFromAreaFlags() error: %v", err)
	}

	pixel, err := exported["Military"].At(1, 1)
	if err != nil {
		t.Fatalf("Military.At() error: %v", err)
	}

	if pixel == 0 {
		t.Fatal("Military mask pixel (1,1) is empty")
	}

	pixel, err = exported["Police"].At(2, 2)
	if err != nil {
		t.Fatalf("Police.At() error: %v", err)
	}

	if pixel == 0 {
		t.Fatal("Police mask pixel (2,2) is empty")
	}

	territoryPath := filepath.Join(territoryDir, "zombie_territories.xml")
	territoryData, err := os.ReadFile(territoryPath)
	if err != nil {
		t.Fatalf("ReadFile(%s) error: %v", territoryPath, err)
	}

	territory, err := DecodeTerritory(territoryData)
	if err != nil {
		t.Fatalf("DecodeTerritory() error: %v", err)
	}

	if len(territory.Territories) != 1 {
		t.Fatalf("territory count = %d, want 1", len(territory.Territories))
	}
}

func TestCEPipelineRestoreWithTemplate(t *testing.T) {
	templateData := readFixture(t, filepath.Join("positive", "ceproject-config.xml"))
	templatePath := filepath.Join(t.TempDir(), "template.xml")
	if err := os.WriteFile(templatePath, templateData, 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", templatePath, err)
	}

	templateConfig, err := DecodeCEProjectConfig(templateData)
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

	if err = militaryMask.Set(3, 1, 255); err != nil {
		t.Fatalf("militaryMask.Set() error: %v", err)
	}

	if err = policeMask.Set(5, 2, 255); err != nil {
		t.Fatalf("policeMask.Set() error: %v", err)
	}

	mapFile, err := BuildAreaFlagsFromCEProjectMasks(templateConfig, map[string]*MaskImage{
		"Military": militaryMask,
		"Police":   policeMask,
	})
	if err != nil {
		t.Fatalf("BuildAreaFlagsFromCEProjectMasks() error: %v", err)
	}

	mapData, err := EncodeAreaFlagsMap(mapFile)
	if err != nil {
		t.Fatalf("EncodeAreaFlagsMap() error: %v", err)
	}

	mapPath := filepath.Join(t.TempDir(), "in", "areaflags.map")
	if err = os.MkdirAll(filepath.Dir(mapPath), 0o700); err != nil {
		t.Fatalf("MkdirAll(%s) error: %v", filepath.Dir(mapPath), err)
	}

	if err = os.WriteFile(mapPath, mapData, 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", mapPath, err)
	}

	territoryDir := filepath.Join(t.TempDir(), "territories")
	if err = SaveTerritoryFilesToDir(map[string]*TerritoryFile{
		"zombie_territories": {
			Territories: []TerritoryBlock{
				{
					Color: "4281913127",
					Zones: []TerritoryZone{
						{
							Name: "City",
							X:    "123",
							Z:    "456",
							R:    "30",
						},
					},
				},
			},
		},
	}, territoryDir); err != nil {
		t.Fatalf("SaveTerritoryFilesToDir() error: %v", err)
	}

	outConfigPath := filepath.Join(t.TempDir(), "restore", "project.xml")
	outLayersDir := filepath.Join(t.TempDir(), "restore", "layers")
	pipeline := CEPipeline{
		ConfigPath:       outConfigPath,
		LayersDir:        outLayersDir,
		TerritoriesDir:   territoryDir,
		AreaFlagsMapPath: mapPath,
	}

	result, err := pipeline.Restore(CERestoreOptions{
		TemplateConfigPath: templatePath,
		RequireTerritories: true,
	})
	if err != nil {
		t.Fatalf("pipeline.Restore() error: %v", err)
	}

	if result.LayerMaskCount != 2 {
		t.Fatalf("LayerMaskCount = %d, want 2", result.LayerMaskCount)
	}

	if result.TerritoryFileCount != 1 {
		t.Fatalf("TerritoryFileCount = %d, want 1", result.TerritoryFileCount)
	}

	restoredConfigData, err := os.ReadFile(outConfigPath)
	if err != nil {
		t.Fatalf("ReadFile(%s) error: %v", outConfigPath, err)
	}

	restoredConfig, err := DecodeCEProjectConfig(restoredConfigData)
	if err != nil {
		t.Fatalf("DecodeCEProjectConfig(restored) error: %v", err)
	}

	if restoredConfig.TerritoryTypeList == nil {
		t.Fatal("restored config has nil TerritoryTypeList")
	}

	if len(restoredConfig.TerritoryTypeList.Types) != 1 {
		t.Fatalf(
			"restored territory types = %d, want 1",
			len(restoredConfig.TerritoryTypeList.Types),
		)
	}

	if len(restoredConfig.TerritoryTypeList.Types[0].Territories) != 1 {
		t.Fatalf("restored territory count mismatch")
	}

	if len(restoredConfig.TerritoryTypeList.Types[0].Territories[0].Zones) != 1 {
		t.Fatalf("restored zone count mismatch")
	}

	if restoredConfig.TerritoryTypeList.Types[0].Territories[0].Zones[0].X != "123" {
		t.Fatalf("restored zone X mismatch")
	}

	militaryRestored, err := LoadTGAMaskFile(filepath.Join(outLayersDir, "Military.tga"))
	if err != nil {
		t.Fatalf("LoadTGAMaskFile(Military.tga) error: %v", err)
	}

	value, err := militaryRestored.At(3, 1)
	if err != nil {
		t.Fatalf("Military.At() error: %v", err)
	}

	if value == 0 {
		t.Fatal("restored Military mask pixel (3,1) is empty")
	}

	policeRestored, err := LoadTGAMaskFile(filepath.Join(outLayersDir, "Police.tga"))
	if err != nil {
		t.Fatalf("LoadTGAMaskFile(Police.tga) error: %v", err)
	}

	value, err = policeRestored.At(5, 2)
	if err != nil {
		t.Fatalf("Police.At() error: %v", err)
	}

	if value == 0 {
		t.Fatal("restored Police mask pixel (5,2) is empty")
	}
}

func TestCEPipelineRestoreWithoutTemplateBuildsSyntheticConfig(t *testing.T) {
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

	layerSize, err := file.Header.LayerByteSize()
	if err != nil {
		t.Fatalf("LayerByteSize() error: %v", err)
	}

	for index := range file.Layers {
		file.Layers[index].Bits = make([]byte, layerSize)
	}

	usageMask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	valueMask, err := NewMaskImage(8, 8)
	if err != nil {
		t.Fatalf("NewMaskImage() error: %v", err)
	}

	if err = usageMask.Set(2, 3, 255); err != nil {
		t.Fatalf("usageMask.Set() error: %v", err)
	}

	if err = valueMask.Set(6, 1, 255); err != nil {
		t.Fatalf("valueMask.Set() error: %v", err)
	}

	if err = file.SetLayerFromMask(0, usageMask); err != nil {
		t.Fatalf("SetLayerFromMask(usage) error: %v", err)
	}

	if err = file.SetLayerFromMask(16, valueMask); err != nil {
		t.Fatalf("SetLayerFromMask(value) error: %v", err)
	}

	data, err := EncodeAreaFlagsMap(file)
	if err != nil {
		t.Fatalf("EncodeAreaFlagsMap() error: %v", err)
	}

	mapPath := filepath.Join(t.TempDir(), "in", "areaflags.map")
	if err = os.MkdirAll(filepath.Dir(mapPath), 0o700); err != nil {
		t.Fatalf("MkdirAll(%s) error: %v", filepath.Dir(mapPath), err)
	}

	if err = os.WriteFile(mapPath, data, 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", mapPath, err)
	}

	configPath := filepath.Join(t.TempDir(), "out", "project.xml")
	layersDir := filepath.Join(t.TempDir(), "out", "layers")
	pipeline := CEPipeline{
		ConfigPath:       configPath,
		LayersDir:        layersDir,
		AreaFlagsMapPath: mapPath,
	}

	result, err := pipeline.Restore(CERestoreOptions{})
	if err != nil {
		t.Fatalf("pipeline.Restore() error: %v", err)
	}

	if result.LayerMaskCount != 20 {
		t.Fatalf("LayerMaskCount = %d, want 20", result.LayerMaskCount)
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(%s) error: %v", configPath, err)
	}

	config, err := DecodeCEProjectConfig(configData)
	if err != nil {
		t.Fatalf("DecodeCEProjectConfig() error: %v", err)
	}

	if config.Areas == nil || config.Areas.Usages == nil || config.Areas.Values == nil {
		t.Fatalf("synthetic config has missing areas")
	}

	if len(config.Areas.Usages.Items) != 16 {
		t.Fatalf("usage names count = %d, want 16", len(config.Areas.Usages.Items))
	}

	if len(config.Areas.Values.Items) != 4 {
		t.Fatalf("value names count = %d, want 4", len(config.Areas.Values.Items))
	}

	if config.Layers == nil || len(config.Layers.Layers) != 20 {
		t.Fatalf("synthetic layers count = %d, want 20", len(config.Layers.Layers))
	}

	if config.Layers.Layers[0].Name != "usage_bit_00" {
		t.Fatalf("first synthetic layer name = %q, want usage_bit_00", config.Layers.Layers[0].Name)
	}

	if config.Layers.Layers[16].Name != "value_bit_00" {
		t.Fatalf("synthetic value layer name = %q, want value_bit_00", config.Layers.Layers[16].Name)
	}

	if config.Layers.Layers[0].Color == "4294967295" {
		t.Fatal("synthetic usage layer color should not be plain white")
	}

	if config.Layers.Layers[16].Color == "4294967295" {
		t.Fatal("synthetic value layer color should not be plain white")
	}

	usageLayer, err := LoadTGAMaskFile(filepath.Join(layersDir, "usage_bit_00.tga"))
	if err != nil {
		t.Fatalf("LoadTGAMaskFile(usage_bit_00.tga) error: %v", err)
	}

	pixel, err := usageLayer.At(2, 3)
	if err != nil {
		t.Fatalf("usage layer At() error: %v", err)
	}

	if pixel == 0 {
		t.Fatal("usage_bit_00 pixel (2,3) is empty")
	}

	valueLayer, err := LoadTGAMaskFile(filepath.Join(layersDir, "value_bit_00.tga"))
	if err != nil {
		t.Fatalf("LoadTGAMaskFile(value_bit_00.tga) error: %v", err)
	}

	pixel, err = valueLayer.At(6, 1)
	if err != nil {
		t.Fatalf("value layer At() error: %v", err)
	}

	if pixel == 0 {
		t.Fatal("value_bit_00 pixel (6,1) is empty")
	}
}
