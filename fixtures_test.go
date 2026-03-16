// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// readFixture reads one test fixture file by relative path.
func readFixture(t *testing.T, rel string) []byte {
	t.Helper()

	path := filepath.Join("testdata", rel)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s) error: %v", path, err)
	}

	return data
}

func TestFixturesPositive(t *testing.T) {
	tests := []struct {
		rel  string
		kind Kind
	}{
		{rel: filepath.Join("positive", "db", "types.xml"), kind: KindTypes},
		{rel: filepath.Join("positive", "db", "events.xml"), kind: KindEvents},
		{rel: filepath.Join("positive", "db", "economy.xml"), kind: KindEconomy},
		{rel: filepath.Join("positive", "db", "globals.xml"), kind: KindGlobals},
		{rel: filepath.Join("positive", "db", "messages.xml"), kind: KindMessages},
		{rel: filepath.Join("positive", "cfgspawnabletypes.xml"), kind: KindSpawnableTypes},
		{rel: filepath.Join("positive", "cfgrandompresets.xml"), kind: KindRandomPresets},
		{rel: filepath.Join("positive", "cfgeconomycore.xml"), kind: KindEconomyCore},
		{rel: filepath.Join("positive", "cfgenvironment.xml"), kind: KindEnvironment},
		{rel: filepath.Join("positive", "cfgeventspawns.xml"), kind: KindEventSpawns},
		{rel: filepath.Join("positive", "cfgeventgroups.xml"), kind: KindEventGroups},
		{rel: filepath.Join("positive", "cfgplayerspawnpoints.xml"), kind: KindPlayerSpawnPoints},
		{rel: filepath.Join("positive", "cfgweather.xml"), kind: KindWeather},
		{rel: filepath.Join("positive", "cfglimitsdefinition.xml"), kind: KindLimitsDefinition},
		{rel: filepath.Join("positive", "cfglimitsdefinitionuser.xml"), kind: KindLimitsDefinitionUser},
		{rel: filepath.Join("positive", "cfgignorelist.xml"), kind: KindIgnoreList},
		{rel: filepath.Join("positive", "env", "hare_territories.xml"), kind: KindTerritories},
		{rel: filepath.Join("positive", "cfgundergroundtriggers.json"), kind: KindUndergroundTriggers},
		{rel: filepath.Join("positive", "cfgeffectarea.json"), kind: KindEffectArea},
		{rel: filepath.Join("positive", "cfggameplay.json"), kind: KindGameplay},
		{rel: filepath.Join("positive", "mapgroupproto.xml"), kind: KindMapGroupProto},
		{rel: filepath.Join("positive", "mapclusterproto.xml"), kind: KindMapClusterProto},
		{rel: filepath.Join("positive", "mapgrouppos.xml"), kind: KindMapGroupPos},
		{rel: filepath.Join("positive", "mapgroupdirt.xml"), kind: KindMapGroupDirt},
		{rel: filepath.Join("positive", "mapgroupcluster.xml"), kind: KindMapGroupCluster},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.rel, func(t *testing.T) {
			data := readFixture(t, tc.rel)
			path := filepath.Join("testdata", tc.rel)

			kind := DetectKind(path)
			if kind != tc.kind {
				t.Fatalf("DetectKind(%s) = %q, want %q", path, kind, tc.kind)
			}

			value, err := Decode(kind, data)
			if err != nil {
				t.Fatalf("Decode(%s) error: %v", path, err)
			}

			if _, err = Encode(kind, value); err != nil {
				t.Fatalf("Encode(%s) error: %v", path, err)
			}

			outPath := filepath.Join(t.TempDir(), filepath.Base(path))
			if err = SaveFile(outPath, value); err != nil {
				t.Fatalf("SaveFile(%s) error: %v", outPath, err)
			}

			reloadedKind, reloadedValue, err := LoadFile(outPath)
			if err != nil {
				t.Fatalf("LoadFile(%s) error: %v", outPath, err)
			}

			if reloadedKind != kind {
				t.Fatalf("LoadFile(%s) kind = %q, want %q", outPath, reloadedKind, kind)
			}

			if !reflect.DeepEqual(value, reloadedValue) {
				t.Fatalf("LoadFile(%s) value mismatch after SaveFile", outPath)
			}
		})
	}
}

func TestFixturesNegative(t *testing.T) {
	tests := []struct {
		rel  string
		kind Kind
	}{
		{rel: filepath.Join("negative", "db", "types.xml"), kind: KindTypes},
		{rel: filepath.Join("negative", "db", "events.xml"), kind: KindEvents},
		{rel: filepath.Join("negative", "db", "messages.xml"), kind: KindMessages},
		{rel: filepath.Join("negative", "cfgenvironment.xml"), kind: KindEnvironment},
		{rel: filepath.Join("negative", "cfgundergroundtriggers.json"), kind: KindUndergroundTriggers},
		{rel: filepath.Join("negative", "cfgeffectarea.json"), kind: KindEffectArea},
		{rel: filepath.Join("negative", "cfggameplay.json"), kind: KindGameplay},
		{rel: filepath.Join("negative", "mapgroupproto.xml"), kind: KindMapGroupProto},
		{rel: filepath.Join("negative", "mapgroupcluster.xml"), kind: KindMapGroupCluster},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.rel, func(t *testing.T) {
			data := readFixture(t, tc.rel)
			path := filepath.Join("testdata", tc.rel)

			kind := DetectKind(path)
			if kind != tc.kind {
				t.Fatalf("DetectKind(%s) = %q, want %q", path, kind, tc.kind)
			}

			if _, err := Decode(kind, data); err == nil {
				t.Fatalf("Decode(%s) expected error", path)
			}

			tmpPath := filepath.Join(t.TempDir(), filepath.Base(path))
			if err := os.WriteFile(tmpPath, data, 0o600); err != nil {
				t.Fatalf("WriteFile(%s) error: %v", tmpPath, err)
			}

			if _, _, err := LoadFile(tmpPath); err == nil {
				t.Fatalf("LoadFile(%s) expected decode error", tmpPath)
			}
		})
	}
}
