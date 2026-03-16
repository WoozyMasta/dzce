// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var sampleFixtureByKind = map[Kind]string{
	KindTypes:                filepath.Join("positive", "db", "types.xml"),
	KindEvents:               filepath.Join("positive", "db", "events.xml"),
	KindEconomy:              filepath.Join("positive", "db", "economy.xml"),
	KindGlobals:              filepath.Join("positive", "db", "globals.xml"),
	KindMessages:             filepath.Join("positive", "db", "messages.xml"),
	KindSpawnableTypes:       filepath.Join("positive", "cfgspawnabletypes.xml"),
	KindRandomPresets:        filepath.Join("positive", "cfgrandompresets.xml"),
	KindEconomyCore:          filepath.Join("positive", "cfgeconomycore.xml"),
	KindEnvironment:          filepath.Join("positive", "cfgenvironment.xml"),
	KindEventSpawns:          filepath.Join("positive", "cfgeventspawns.xml"),
	KindEventGroups:          filepath.Join("positive", "cfgeventgroups.xml"),
	KindPlayerSpawnPoints:    filepath.Join("positive", "cfgplayerspawnpoints.xml"),
	KindWeather:              filepath.Join("positive", "cfgweather.xml"),
	KindLimitsDefinition:     filepath.Join("positive", "cfglimitsdefinition.xml"),
	KindLimitsDefinitionUser: filepath.Join("positive", "cfglimitsdefinitionuser.xml"),
	KindIgnoreList:           filepath.Join("positive", "cfgignorelist.xml"),
	KindTerritories:          filepath.Join("positive", "env", "hare_territories.xml"),
	KindUndergroundTriggers:  filepath.Join("positive", "cfgundergroundtriggers.json"),
	KindEffectArea:           filepath.Join("positive", "cfgeffectarea.json"),
	KindGameplay:             filepath.Join("positive", "cfggameplay.json"),
	KindGameplayGearPresets:  filepath.Join("positive", "gameplay-gear-presets.json"),
	KindObjectSpawner:        filepath.Join("positive", "object-spawner.json"),
	KindCEProjectConfig:      filepath.Join("positive", "ceproject-config.xml"),
	KindMapGroupProto:        filepath.Join("positive", "mapgroupproto.xml"),
	KindMapClusterProto:      filepath.Join("positive", "mapclusterproto.xml"),
	KindMapGroupPos:          filepath.Join("positive", "mapgrouppos.xml"),
	KindMapGroupDirt:         filepath.Join("positive", "mapgroupdirt.xml"),
	KindMapGroupCluster:      filepath.Join("positive", "mapgroupcluster.xml"),
}

// readSampleByKind loads one positive fixture payload for a known Kind.
func readSampleByKind(t *testing.T, kind Kind) []byte {
	t.Helper()

	rel, ok := sampleFixtureByKind[kind]
	if !ok {
		t.Fatalf("sample fixture path is not defined for kind %q", kind)
	}

	return readFixture(t, rel)
}

func TestDetectKind(t *testing.T) {
	tests := []struct {
		path string
		want Kind
	}{
		{path: `db\types.xml`, want: KindTypes},
		{path: `db\events.xml`, want: KindEvents},
		{path: `db\economy.xml`, want: KindEconomy},
		{path: `db\globals.xml`, want: KindGlobals},
		{path: `db\messages.xml`, want: KindMessages},
		{path: `cfgspawnabletypes.xml`, want: KindSpawnableTypes},
		{path: `cfgrandompresets.xml`, want: KindRandomPresets},
		{path: `cfgeconomycore.xml`, want: KindEconomyCore},
		{path: `cfgenvironment.xml`, want: KindEnvironment},
		{path: `cfgeventspawns.xml`, want: KindEventSpawns},
		{path: `cfgeventgroups.xml`, want: KindEventGroups},
		{path: `cfgplayerspawnpoints.xml`, want: KindPlayerSpawnPoints},
		{path: `cfgweather.xml`, want: KindWeather},
		{path: `cfglimitsdefinition.xml`, want: KindLimitsDefinition},
		{path: `cfglimitsdefinitionuser.xml`, want: KindLimitsDefinitionUser},
		{path: `cfgignorelist.xml`, want: KindIgnoreList},
		{path: `env\hare_territories.xml`, want: KindTerritories},
		{path: `cfgundergroundtriggers.json`, want: KindUndergroundTriggers},
		{path: `cfgeffectarea.json`, want: KindEffectArea},
		{path: `cfggameplay.json`, want: KindGameplay},
		{path: `areaflags.map`, want: KindAreaFlagsMap},
		{path: `mapgroupproto.xml`, want: KindMapGroupProto},
		{path: `mapclusterproto.xml`, want: KindMapClusterProto},
		{path: `mapgrouppos.xml`, want: KindMapGroupPos},
		{path: `mapgroupdirt.xml`, want: KindMapGroupDirt},
		{path: `mapgroupcluster.xml`, want: KindMapGroupCluster},
		{path: `mapgroupcluster02.xml`, want: KindMapGroupCluster},
		{path: `unknown.xml`, want: KindUnknown},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.path, func(t *testing.T) {
			if got := DetectKind(tc.path); got != tc.want {
				t.Fatalf("DetectKind(%q) = %q, want %q", tc.path, got, tc.want)
			}
		})
	}
}

func TestDecodeEncodeRoundtripByKind(t *testing.T) {
	for kind := range sampleFixtureByKind {
		kind := kind

		t.Run(string(kind), func(t *testing.T) {
			first, err := Decode(kind, readSampleByKind(t, kind))
			if err != nil {
				t.Fatalf("Decode(%s) error: %v", kind, err)
			}

			encodedFirst, err := Encode(kind, first)
			if err != nil {
				t.Fatalf("Encode(%s) error: %v", kind, err)
			}

			if kind != KindUndergroundTriggers && kind != KindEffectArea &&
				kind != KindGameplay && kind != KindGameplayGearPresets &&
				kind != KindObjectSpawner &&
				!strings.HasPrefix(string(encodedFirst), xml.Header) {
				t.Fatalf("Encode(%s) does not contain xml header", kind)
			}

			second, err := Decode(kind, encodedFirst)
			if err != nil {
				t.Fatalf("Decode(%s) on encoded payload error: %v", kind, err)
			}

			encodedSecond, err := Encode(kind, second)
			if err != nil {
				t.Fatalf("Encode(%s) second pass error: %v", kind, err)
			}

			if string(encodedFirst) != string(encodedSecond) {
				t.Fatalf("Encode(%s) output is not stable after roundtrip", kind)
			}
		})
	}
}

func TestDecodeNegativeByKind(t *testing.T) {
	tests := []struct {
		name    string
		kind    Kind
		fixture string
	}{
		{
			name:    "gameplay gear presets",
			kind:    KindGameplayGearPresets,
			fixture: filepath.Join("negative", "gameplay-gear-presets.json"),
		},
		{
			name:    "object spawner",
			kind:    KindObjectSpawner,
			fixture: filepath.Join("negative", "object-spawner.json"),
		},
		{
			name:    "ce project config",
			kind:    KindCEProjectConfig,
			fixture: filepath.Join("negative", "ceproject-config.xml"),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			if _, err := Decode(tc.kind, readFixture(t, tc.fixture)); err == nil {
				t.Fatalf("Decode(%q) expected error", tc.kind)
			}
		})
	}
}

func TestTypedDecodeFunctions(t *testing.T) {
	typesFile, err := DecodeTypes(readSampleByKind(t, KindTypes))
	if err != nil {
		t.Fatalf("DecodeTypes() error: %v", err)
	}

	if len(typesFile.Types) != 1 {
		t.Fatalf("DecodeTypes() type count = %d, want 1", len(typesFile.Types))
	}

	eventsFile, err := DecodeEvents(readSampleByKind(t, KindEvents))
	if err != nil {
		t.Fatalf("DecodeEvents() error: %v", err)
	}

	if len(eventsFile.Events) != 1 {
		t.Fatalf("DecodeEvents() event count = %d, want 1", len(eventsFile.Events))
	}

	globalsFile, err := DecodeGlobals(readSampleByKind(t, KindGlobals))
	if err != nil {
		t.Fatalf("DecodeGlobals() error: %v", err)
	}

	if len(globalsFile.Vars) != 3 {
		t.Fatalf("DecodeGlobals() variable count = %d, want 3", len(globalsFile.Vars))
	}

	messagesFile, err := DecodeMessages(readSampleByKind(t, KindMessages))
	if err != nil {
		t.Fatalf("DecodeMessages() error: %v", err)
	}

	if len(messagesFile.Messages) != 1 {
		t.Fatalf("DecodeMessages() messages count = %d, want 1", len(messagesFile.Messages))
	}

	environment, err := DecodeEnvironment(readSampleByKind(t, KindEnvironment))
	if err != nil {
		t.Fatalf("DecodeEnvironment() error: %v", err)
	}

	if environment.Territories == nil || len(environment.Territories.Territories) != 1 {
		t.Fatalf("DecodeEnvironment() territory count mismatch")
	}

	spawns, err := DecodeEventSpawns(readSampleByKind(t, KindEventSpawns))
	if err != nil {
		t.Fatalf("DecodeEventSpawns() error: %v", err)
	}

	if len(spawns.Events) != 1 {
		t.Fatalf("DecodeEventSpawns() event count = %d, want 1", len(spawns.Events))
	}

	territory, err := DecodeTerritory(readSampleByKind(t, KindTerritories))
	if err != nil {
		t.Fatalf("DecodeTerritory() error: %v", err)
	}

	if len(territory.Territories) != 1 {
		t.Fatalf("DecodeTerritory() territories count = %d, want 1", len(territory.Territories))
	}

	mapGroupProto, err := DecodeMapGroupProto(readSampleByKind(t, KindMapGroupProto))
	if err != nil {
		t.Fatalf("DecodeMapGroupProto() error: %v", err)
	}

	if len(mapGroupProto.Groups) != 1 {
		t.Fatalf("DecodeMapGroupProto() groups count = %d, want 1", len(mapGroupProto.Groups))
	}

	mapClusterProto, err := DecodeMapClusterProto(readSampleByKind(t, KindMapClusterProto))
	if err != nil {
		t.Fatalf("DecodeMapClusterProto() error: %v", err)
	}

	if len(mapClusterProto.ClusterGroups) != 1 {
		t.Fatalf("DecodeMapClusterProto() clusters count = %d, want 1", len(mapClusterProto.ClusterGroups))
	}

	mapGroupPos, err := DecodeMapGroupPos(readSampleByKind(t, KindMapGroupPos))
	if err != nil {
		t.Fatalf("DecodeMapGroupPos() error: %v", err)
	}

	if len(mapGroupPos.Groups) != 1 {
		t.Fatalf("DecodeMapGroupPos() groups count = %d, want 1", len(mapGroupPos.Groups))
	}

	mapGroupDirt, err := DecodeMapGroupDirt(readSampleByKind(t, KindMapGroupDirt))
	if err != nil {
		t.Fatalf("DecodeMapGroupDirt() error: %v", err)
	}

	if len(mapGroupDirt.Groups) != 1 {
		t.Fatalf("DecodeMapGroupDirt() groups count = %d, want 1", len(mapGroupDirt.Groups))
	}

	mapGroupCluster, err := DecodeMapGroupCluster(readSampleByKind(t, KindMapGroupCluster))
	if err != nil {
		t.Fatalf("DecodeMapGroupCluster() error: %v", err)
	}

	if len(mapGroupCluster.Groups) != 1 {
		t.Fatalf("DecodeMapGroupCluster() groups count = %d, want 1", len(mapGroupCluster.Groups))
	}

	underground, err := DecodeUndergroundTriggers(
		readSampleByKind(t, KindUndergroundTriggers),
	)
	if err != nil {
		t.Fatalf("DecodeUndergroundTriggers() error: %v", err)
	}

	if len(underground.Triggers) != 1 {
		t.Fatalf(
			"DecodeUndergroundTriggers() triggers count = %d, want 1",
			len(underground.Triggers),
		)
	}

	effectAreas, err := DecodeEffectArea(readSampleByKind(t, KindEffectArea))
	if err != nil {
		t.Fatalf("DecodeEffectArea() error: %v", err)
	}

	if len(effectAreas.Areas) != 2 {
		t.Fatalf("DecodeEffectArea() areas count = %d, want 2", len(effectAreas.Areas))
	}

	gameplay, err := DecodeGameplay(readSampleByKind(t, KindGameplay))
	if err != nil {
		t.Fatalf("DecodeGameplay() error: %v", err)
	}

	if gameplay.Version == nil || *gameplay.Version != 123 {
		t.Fatalf("DecodeGameplay() version mismatch")
	}

	gearPresets, err := DecodeGameplayGearPresets(
		readSampleByKind(t, KindGameplayGearPresets),
	)
	if err != nil {
		t.Fatalf("DecodeGameplayGearPresets() error: %v", err)
	}

	if len(*gearPresets) != 1 {
		t.Fatalf(
			"DecodeGameplayGearPresets() presets count = %d, want 1",
			len(*gearPresets),
		)
	}

	objectSpawner, err := DecodeObjectSpawner(
		readSampleByKind(t, KindObjectSpawner),
	)
	if err != nil {
		t.Fatalf("DecodeObjectSpawner() error: %v", err)
	}

	if len(objectSpawner.Objects) != 1 {
		t.Fatalf(
			"DecodeObjectSpawner() objects count = %d, want 1",
			len(objectSpawner.Objects),
		)
	}

	ceProjectConfig, err := DecodeCEProjectConfig(
		readSampleByKind(t, KindCEProjectConfig),
	)
	if err != nil {
		t.Fatalf("DecodeCEProjectConfig() error: %v", err)
	}

	if ceProjectConfig.Global == nil || ceProjectConfig.Layers == nil {
		t.Fatalf("DecodeCEProjectConfig() expected global and layers sections")
	}
}

func TestEncodeWrongValueType(t *testing.T) {
	_, err := Encode(KindTypes, &EventsFile{})
	if err == nil {
		t.Fatal("Encode(KindTypes, EventsFile) expected error")
	}

	if !errors.Is(err, ErrUnsupportedValue) {
		t.Fatalf("Encode wrong value error = %v, want ErrUnsupportedValue", err)
	}
}

func TestLoadSaveFile(t *testing.T) {
	tests := []struct {
		fileName string
		kind     Kind
		fixture  string
	}{
		{fileName: "types.xml", kind: KindTypes, fixture: filepath.Join("positive", "db", "types.xml")},
		{fileName: "events.xml", kind: KindEvents, fixture: filepath.Join("positive", "db", "events.xml")},
		{fileName: "economy.xml", kind: KindEconomy, fixture: filepath.Join("positive", "db", "economy.xml")},
		{fileName: "globals.xml", kind: KindGlobals, fixture: filepath.Join("positive", "db", "globals.xml")},
		{fileName: "messages.xml", kind: KindMessages, fixture: filepath.Join("positive", "db", "messages.xml")},
		{fileName: "cfgspawnabletypes.xml", kind: KindSpawnableTypes, fixture: filepath.Join("positive", "cfgspawnabletypes.xml")},
		{fileName: "cfgrandompresets.xml", kind: KindRandomPresets, fixture: filepath.Join("positive", "cfgrandompresets.xml")},
		{fileName: "cfgeconomycore.xml", kind: KindEconomyCore, fixture: filepath.Join("positive", "cfgeconomycore.xml")},
		{fileName: "cfgenvironment.xml", kind: KindEnvironment, fixture: filepath.Join("positive", "cfgenvironment.xml")},
		{fileName: "cfgeventspawns.xml", kind: KindEventSpawns, fixture: filepath.Join("positive", "cfgeventspawns.xml")},
		{fileName: "cfgeventgroups.xml", kind: KindEventGroups, fixture: filepath.Join("positive", "cfgeventgroups.xml")},
		{fileName: "cfgplayerspawnpoints.xml", kind: KindPlayerSpawnPoints, fixture: filepath.Join("positive", "cfgplayerspawnpoints.xml")},
		{fileName: "cfgweather.xml", kind: KindWeather, fixture: filepath.Join("positive", "cfgweather.xml")},
		{fileName: "cfglimitsdefinition.xml", kind: KindLimitsDefinition, fixture: filepath.Join("positive", "cfglimitsdefinition.xml")},
		{fileName: "cfglimitsdefinitionuser.xml", kind: KindLimitsDefinitionUser, fixture: filepath.Join("positive", "cfglimitsdefinitionuser.xml")},
		{fileName: "cfgignorelist.xml", kind: KindIgnoreList, fixture: filepath.Join("positive", "cfgignorelist.xml")},
		{fileName: "hare_territories.xml", kind: KindTerritories, fixture: filepath.Join("positive", "env", "hare_territories.xml")},
		{fileName: "cfgundergroundtriggers.json", kind: KindUndergroundTriggers, fixture: filepath.Join("positive", "cfgundergroundtriggers.json")},
		{fileName: "cfgeffectarea.json", kind: KindEffectArea, fixture: filepath.Join("positive", "cfgeffectarea.json")},
		{fileName: "cfggameplay.json", kind: KindGameplay, fixture: filepath.Join("positive", "cfggameplay.json")},
		{fileName: "mapgroupproto.xml", kind: KindMapGroupProto, fixture: filepath.Join("positive", "mapgroupproto.xml")},
		{fileName: "mapclusterproto.xml", kind: KindMapClusterProto, fixture: filepath.Join("positive", "mapclusterproto.xml")},
		{fileName: "mapgrouppos.xml", kind: KindMapGroupPos, fixture: filepath.Join("positive", "mapgrouppos.xml")},
		{fileName: "mapgroupdirt.xml", kind: KindMapGroupDirt, fixture: filepath.Join("positive", "mapgroupdirt.xml")},
		{fileName: "mapgroupcluster.xml", kind: KindMapGroupCluster, fixture: filepath.Join("positive", "mapgroupcluster.xml")},
	}

	rootDir := t.TempDir()
	inDir := filepath.Join(rootDir, "in")
	outDir := filepath.Join(rootDir, "out")

	if err := os.MkdirAll(inDir, 0o700); err != nil {
		t.Fatalf("MkdirAll(%s) error: %v", inDir, err)
	}

	if err := os.MkdirAll(outDir, 0o700); err != nil {
		t.Fatalf("MkdirAll(%s) error: %v", outDir, err)
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.fileName, func(t *testing.T) {
			inPath := filepath.Join(inDir, tc.fileName)
			if err := os.WriteFile(inPath, readFixture(t, tc.fixture), 0o600); err != nil {
				t.Fatalf("WriteFile(%s) error: %v", inPath, err)
			}

			kind, value, err := LoadFile(inPath)
			if err != nil {
				t.Fatalf("LoadFile(%s) error: %v", inPath, err)
			}

			if kind != tc.kind {
				t.Fatalf("LoadFile(%s) kind = %q, want %q", inPath, kind, tc.kind)
			}

			outPath := filepath.Join(outDir, tc.fileName)
			if err := SaveFile(outPath, value); err != nil {
				t.Fatalf("SaveFile(%s) error: %v", outPath, err)
			}

			reloadedKind, reloadedValue, err := LoadFile(outPath)
			if err != nil {
				t.Fatalf("LoadFile(%s) after SaveFile error: %v", outPath, err)
			}

			if reloadedKind != tc.kind {
				t.Fatalf("Reloaded kind = %q, want %q", reloadedKind, tc.kind)
			}

			if !reflect.DeepEqual(value, reloadedValue) {
				t.Fatalf("Reloaded value differs for %s", tc.fileName)
			}
		})
	}
}

func TestLoadFileDetectCEProjectByContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "my-custom-map.xml")
	data := readFixture(t, filepath.Join("positive", "ceproject-config.xml"))
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", path, err)
	}

	kind, value, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile(%s) error: %v", path, err)
	}

	if kind != KindCEProjectConfig {
		t.Fatalf("LoadFile(%s) kind = %q, want %q", path, kind, KindCEProjectConfig)
	}

	if _, ok := value.(*CEProjectConfigFile); !ok {
		t.Fatalf("LoadFile(%s) value type = %T, want *CEProjectConfigFile", path, value)
	}
}

func TestSaveFileDetectByValueType(t *testing.T) {
	cfg, err := DecodeCEProjectConfig(
		readFixture(t, filepath.Join("positive", "ceproject-config.xml")),
	)
	if err != nil {
		t.Fatalf("DecodeCEProjectConfig() error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "custom-output.xml")
	if err = SaveFile(path, cfg); err != nil {
		t.Fatalf("SaveFile(%s) error: %v", path, err)
	}

	kind, value, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile(%s) error: %v", path, err)
	}

	if kind != KindCEProjectConfig {
		t.Fatalf("LoadFile(%s) kind = %q, want %q", path, kind, KindCEProjectConfig)
	}

	if _, ok := value.(*CEProjectConfigFile); !ok {
		t.Fatalf("LoadFile(%s) value type = %T, want *CEProjectConfigFile", path, value)
	}
}

func TestLoadSaveFileAs(t *testing.T) {
	tests := []struct {
		fileName string
		kind     Kind
		fixture  string
	}{
		{
			fileName: "my-gears.json",
			kind:     KindGameplayGearPresets,
			fixture:  filepath.Join("positive", "gameplay-gear-presets.json"),
		},
		{
			fileName: "spawnerData.json",
			kind:     KindObjectSpawner,
			fixture:  filepath.Join("positive", "object-spawner.json"),
		},
	}

	rootDir := t.TempDir()
	inDir := filepath.Join(rootDir, "in")
	outDir := filepath.Join(rootDir, "out")

	if err := os.MkdirAll(inDir, 0o700); err != nil {
		t.Fatalf("MkdirAll(%s) error: %v", inDir, err)
	}

	if err := os.MkdirAll(outDir, 0o700); err != nil {
		t.Fatalf("MkdirAll(%s) error: %v", outDir, err)
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.fileName, func(t *testing.T) {
			inPath := filepath.Join(inDir, tc.fileName)
			if err := os.WriteFile(inPath, readFixture(t, tc.fixture), 0o600); err != nil {
				t.Fatalf("WriteFile(%s) error: %v", inPath, err)
			}

			value, err := LoadFileAs(tc.kind, inPath)
			if err != nil {
				t.Fatalf("LoadFileAs(%q, %s) error: %v", tc.kind, inPath, err)
			}

			outPath := filepath.Join(outDir, tc.fileName)
			if err := SaveFileAs(tc.kind, outPath, value); err != nil {
				t.Fatalf("SaveFileAs(%q, %s) error: %v", tc.kind, outPath, err)
			}

			reloadedValue, err := LoadFileAs(tc.kind, outPath)
			if err != nil {
				t.Fatalf("LoadFileAs(%q, %s) error: %v", tc.kind, outPath, err)
			}

			if !reflect.DeepEqual(value, reloadedValue) {
				t.Fatalf("Reloaded value differs for %s", tc.fileName)
			}
		})
	}
}

func TestLoadFileUnknownKind(t *testing.T) {
	path := filepath.Join(t.TempDir(), "unknown.xml")
	if err := os.WriteFile(path, []byte("<root/>"), 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", path, err)
	}

	_, _, err := LoadFile(path)
	if err == nil {
		t.Fatal("LoadFile(unknown.xml) expected error")
	}

	if !errors.Is(err, ErrUnknownFileKind) {
		t.Fatalf("LoadFile(unknown.xml) error = %v, want ErrUnknownFileKind", err)
	}
}
