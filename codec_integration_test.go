//go:build integration
// +build integration

// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDecodeRealWorldFixtures validates codec compatibility against external
// CE directories listed in DZCE_REAL_CE_DIRS.
func TestDecodeRealWorldFixtures(t *testing.T) {
	raw := strings.TrimSpace(os.Getenv("DZCE_REAL_CE_DIRS"))
	if raw == "" {
		t.Skip("DZCE_REAL_CE_DIRS is empty")
	}

	kinds := []struct {
		kind Kind
		rel  string
	}{
		{kind: KindTypes, rel: filepath.Join("db", "types.xml")},
		{kind: KindEvents, rel: filepath.Join("db", "events.xml")},
		{kind: KindEconomy, rel: filepath.Join("db", "economy.xml")},
		{kind: KindGlobals, rel: filepath.Join("db", "globals.xml")},
		{kind: KindSpawnableTypes, rel: "cfgspawnabletypes.xml"},
		{kind: KindRandomPresets, rel: "cfgrandompresets.xml"},
		{kind: KindEconomyCore, rel: "cfgeconomycore.xml"},
		{kind: KindEnvironment, rel: "cfgenvironment.xml"},
		{kind: KindEventSpawns, rel: "cfgeventspawns.xml"},
		{kind: KindEventGroups, rel: "cfgeventgroups.xml"},
		{kind: KindPlayerSpawnPoints, rel: "cfgplayerspawnpoints.xml"},
		{kind: KindWeather, rel: "cfgweather.xml"},
		{kind: KindLimitsDefinition, rel: "cfglimitsdefinition.xml"},
		{kind: KindLimitsDefinitionUser, rel: "cfglimitsdefinitionuser.xml"},
		{kind: KindIgnoreList, rel: "cfgignorelist.xml"},
		{kind: KindTerritories, rel: filepath.Join("env", "hare_territories.xml")},
		{kind: KindUndergroundTriggers, rel: "cfgundergroundtriggers.json"},
		{kind: KindEffectArea, rel: "cfgeffectarea.json"},
		{kind: KindGameplay, rel: "cfggameplay.json"},
		{kind: KindMapGroupProto, rel: "mapgroupproto.xml"},
		{kind: KindMapClusterProto, rel: "mapclusterproto.xml"},
		{kind: KindMapGroupPos, rel: "mapgrouppos.xml"},
		{kind: KindMapGroupDirt, rel: "mapgroupdirt.xml"},
		{kind: KindMapGroupCluster, rel: "mapgroupcluster.xml"},
	}

	roots := strings.Split(raw, string(os.PathListSeparator))

	for _, root := range roots {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}

		for _, target := range kinds {
			target := target
			path := filepath.Join(root, target.rel)

			t.Run(path, func(t *testing.T) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("ReadFile(%s) error: %v", path, err)
				}

				value, err := Decode(target.kind, data)
				if err != nil {
					t.Fatalf("Decode(%s) error: %v", path, err)
				}

				if _, err = Encode(target.kind, value); err != nil {
					t.Fatalf("Encode(%s) error: %v", path, err)
				}
			})
		}
	}
}
