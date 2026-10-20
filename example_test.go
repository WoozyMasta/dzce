// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/woozymasta/dzce"
)

func ExampleDecode() {
	kind := dzce.DetectKind("db/types.xml")
	value, err := dzce.Decode(kind, []byte(`<types/>`))
	if err != nil {
		panic(err)
	}

	typesFile := value.(*dzce.TypesFile)
	encoded, err := dzce.Encode(kind, typesFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(kind, len(typesFile.Types), len(encoded) > 0)
	// Output: bi.ce.db.types 0 true
}

func ExampleLoadFile() {
	dir, err := os.MkdirTemp("", "dzce-example-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "db", "types.xml")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, []byte(`<types/>`), 0o600); err != nil {
		panic(err)
	}

	kind, value, err := dzce.LoadFile(path)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s %T\n", kind, value)
	// Output: bi.ce.db.types *dzce.TypesFile
}

func ExampleLoadFileAs() {
	dir, err := os.MkdirTemp("", "dzce-example-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "custom-presets.json")
	if err := os.WriteFile(path, []byte(`[{"name":"Survivor","spawnWeight":1}]`), 0o600); err != nil {
		panic(err)
	}

	value, err := dzce.LoadFileAs(dzce.KindGameplayGearPresets, path)
	if err != nil {
		panic(err)
	}

	presets := value.(*dzce.GameplayGearPresetsFile)
	fmt.Println(len(*presets), (*presets)[0].Name)
	// Output: 1 Survivor
}

func ExampleBuildAreaFlagsFromCEProjectMasks() {
	config := &dzce.CEProjectConfigFile{
		Global: &dzce.CEProjectGlobal{
			Layer: &dzce.CEProjectDimension{Size: "8"},
			World: &dzce.CEProjectDimension{Size: "100"},
		},
		Layers: &dzce.CEProjectLayers{},
	}

	mapFile, err := dzce.BuildAreaFlagsFromCEProjectMasks(config, nil)
	if err != nil {
		panic(err)
	}

	encoded, err := dzce.EncodeAreaFlagsMap(mapFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(mapFile.Header.LayerWidth, len(encoded) > 0)
	// Output: 8 true
}

func ExampleRegisterMaskImageCodec() {
	const format = "example-mask"
	if err := dzce.RegisterMaskImageCodec(format, dzce.MaskImageCodec{
		Decode: func([]byte) (*dzce.MaskImage, error) {
			return &dzce.MaskImage{
				Pixels: []byte{255},
				Width:  1,
				Height: 1,
			}, nil
		},
	}); err != nil {
		panic(err)
	}
	defer dzce.UnregisterMaskImageCodec(format)

	mask, detectedFormat, err := dzce.DecodeMaskImage([]byte("custom image data"))
	if err != nil {
		panic(err)
	}

	fmt.Println(detectedFormat, mask.Width, mask.Height)
	// Output: example-mask 1 1
}

func ExampleLoadMergedEconomyCore() {
	dir, err := os.MkdirTemp("", "dzce-example-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	coreDir := filepath.Join(dir, "db")
	if err := os.MkdirAll(coreDir, 0o700); err != nil {
		panic(err)
	}
	if err := os.WriteFile(
		filepath.Join(dir, "cfgeconomycore.xml"),
		[]byte(`<economycore><ce folder="db"><file name="types.xml" type="types"/></ce></economycore>`),
		0o600,
	); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(coreDir, "types.xml"), []byte(`<types/>`), 0o600); err != nil {
		panic(err)
	}

	merged, err := dzce.LoadMergedEconomyCore(filepath.Join(dir, "cfgeconomycore.xml"))
	if err != nil {
		panic(err)
	}

	value, ok := merged.Get(dzce.KindTypes)
	if !ok {
		panic("types document not found")
	}

	fmt.Println(len(merged.Sources), len(value.(*dzce.TypesFile).Types))
	// Output: 1 0
}
