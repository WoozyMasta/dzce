// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

/*
Package dzce provides typed models, codecs, and merge helpers for
DayZ Central Economy (CE) mission configuration files.

Core flow by detected kind:

	kind := dzce.DetectKind("db/types.xml")

	value, err := dzce.Decode(kind, input)
	if err != nil {
		// handle decode error
	}

	out, err := dzce.Encode(kind, value)
	if err != nil {
		// handle encode error
	}

	_ = out

File flow:

	kind, value, err := dzce.LoadFile("path/to/db/types.xml")
	if err != nil {
		// handle load error
	}

	_ = kind
	_ = value

Arbitrary file name flow (for formats without stable basename):

	value, err := dzce.LoadFileAs(
		dzce.KindGameplayGearPresets,
		"path/to/custom/presets.json",
	)
	if err != nil {
		// handle load error
	}

	_ = value

LoadFile also auto-detects CEProject XML by root element (`<zg-config>`). SaveFile
can infer kind from value type for `*CEProjectConfigFile` and `*AreaFlagsMapFile`.

CE project and areaflags flow:

	cfg, err := dzce.DecodeCEProjectConfig(ceProjectXML)
	if err != nil {
		// handle decode error
	}

	mapFile, err := dzce.BuildAreaFlagsFromCEProjectLayerDir(
		cfg,
		"path/to/layers",
	)
	if err != nil {
		// handle build error
	}

	rawMap, err := dzce.EncodeAreaFlagsMap(mapFile)
	if err != nil {
		// handle encode error
	}

	_ = rawMap

When projecting `areaflags.map` back to CEProject layer masks, layers with both
`usage_flags=0` and `value_flags=0` are returned as empty masks because CEProject
does not persist them in map payload.

If CEProject source `layers/` directory is available, zero/zero layers can be
restored from it:

	layers, err := dzce.ExportCEProjectLayerMasksFromAreaFlagsWithLayerDir(
		cfg,
		mapFile,
		"path/to/layers",
	)
	if err != nil {
		// handle export error
	}

	_ = layers

Mask image codecs:

Built-in mask image codec is TGA. Any custom format can be added by
registering callbacks:

	err := dzce.RegisterMaskImageCodec("dds", dzce.MaskImageCodec{
		Decode: func(data []byte) (*dzce.MaskImage, error) {
			return nil, nil
		},
		Encode: func(
			mask *dzce.MaskImage,
			opts dzce.MaskImageEncodeOptions,
		) ([]byte, error) {
			return nil, nil
		},
	})
	if err != nil {
		// handle registration error
	}

	defer dzce.UnregisterMaskImageCodec("dds")

High-level CE project pipeline:

	pipeline := dzce.CEPipeline{
		ConfigPath:       "CEProject/MyMap/mymap.xml",
		LayersDir:        "CEProject/MyMap/layers",
		TerritoriesDir:   "CEProject/MyMap/territoryTypes",
		AreaFlagsMapPath: "CEProject/MyMap/map/mymap.map",
	}

	_, err = pipeline.Forward()
	if err != nil {
		// handle forward error
	}

	_, err = pipeline.Restore(dzce.CERestoreOptions{
		TemplateConfigPath: "backup/mymap.xml",
	})
	if err != nil {
		// handle restore error
	}

Economycore include merge:

	merged, err := dzce.LoadMergedEconomyCore("path/to/cfgeconomycore.xml")
	if err != nil {
		// handle merge error
	}

	rawTypes, ok := merged.Get(dzce.KindTypes)
	if ok {
		typesDoc := rawTypes.(*dzce.TypesFile)
		_ = typesDoc
	}

Strict include mode (default) accepts wiki-compatible include types:
types, spawnabletypes, globals, economy, events, messages; economycore is
supported for recursive include traversal.

Relaxed include mode (opt-in):

	merged, err := dzce.LoadMergedEconomyCoreWithOptions(
		"path/to/cfgeconomycore.xml",
		dzce.MergeOptions{RelaxedIncludeTypes: true},
	)
	if err != nil {
		// handle merge error
	}

	_ = merged

Merge priority is load-order based. Later include files override earlier ones
where the merge rule for that kind allows overriding. Traversal order is
declaration order with depth-first recursion for nested economycore includes.

LoadMergedEconomyCore* merges include tree content only. It does not preload
vanilla/base CE datasets.
*/
package dzce
