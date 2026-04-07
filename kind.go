// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"strings"

	"github.com/woozymasta/bimime"
)

const (
	// KindUnknown is used when file kind cannot be detected.
	KindUnknown Kind = ""

	// KindTypes maps to db/types.xml.
	KindTypes Kind = "bi.ce.db.types"

	// KindEvents maps to db/events.xml.
	KindEvents Kind = "bi.ce.db.events"

	// KindEconomy maps to db/economy.xml.
	KindEconomy Kind = "bi.ce.db.economy"

	// KindGlobals maps to db/globals.xml.
	KindGlobals Kind = "bi.ce.db.globals"

	// KindMessages maps to db/messages.xml.
	KindMessages Kind = "bi.ce.db.messages"

	// KindSpawnableTypes maps to cfgspawnabletypes.xml.
	KindSpawnableTypes Kind = "bi.ce.cfgspawnabletypes"

	// KindRandomPresets maps to cfgrandompresets.xml.
	KindRandomPresets Kind = "bi.ce.cfgrandompresets"

	// KindEconomyCore maps to cfgeconomycore.xml.
	KindEconomyCore Kind = "bi.ce.cfgeconomycore"

	// KindEnvironment maps to cfgenvironment.xml.
	KindEnvironment Kind = "bi.ce.cfgenvironment"

	// KindEventSpawns maps to cfgeventspawns.xml.
	KindEventSpawns Kind = "bi.ce.cfgeventspawns"

	// KindEventGroups maps to cfgeventgroups.xml.
	KindEventGroups Kind = "bi.ce.cfgeventgroups"

	// KindPlayerSpawnPoints maps to cfgplayerspawnpoints.xml.
	KindPlayerSpawnPoints Kind = "bi.ce.cfgplayerspawnpoints"

	// KindWeather maps to cfgweather.xml.
	KindWeather Kind = "bi.ce.cfgweather"

	// KindLimitsDefinition maps to cfglimitsdefinition.xml.
	KindLimitsDefinition Kind = "bi.ce.cfglimitsdefinition"

	// KindLimitsDefinitionUser maps to cfglimitsdefinitionuser.xml.
	KindLimitsDefinitionUser Kind = "bi.ce.cfglimitsdefinitionuser"

	// KindIgnoreList maps to cfgignorelist.xml.
	KindIgnoreList Kind = "bi.ce.cfgignorelist"

	// KindTerritories maps to env/*_territories.xml files.
	KindTerritories Kind = "bi.ce.env.territories"

	// KindUndergroundTriggers maps to cfgundergroundtriggers.json.
	KindUndergroundTriggers Kind = "bi.ce.cfgundergroundtriggers"

	// KindEffectArea maps to cfgeffectarea.json.
	KindEffectArea Kind = "bi.ce.cfgeffectarea"

	// KindGameplay maps to cfggameplay.json.
	KindGameplay Kind = "bi.ce.cfggameplay"

	// KindGameplayGearPresets maps to JSON payloads listed in
	// `cfggameplay.json -> PlayerData.spawnGearPresetFiles`.
	KindGameplayGearPresets Kind = "bi.ce.gameplay-gear-presets"

	// KindObjectSpawner maps to JSON payloads listed in
	// `cfggameplay.json -> WorldsData.objectSpawnersArr`.
	KindObjectSpawner Kind = "bi.ce.object-spawner"

	// KindCEProjectConfig maps to CEProject `mapname.xml` (`<zg-config>`).
	KindCEProjectConfig Kind = "bi.ce.ceproject-config"

	// KindAreaFlagsMap maps to `areaflags.map` binary payload.
	KindAreaFlagsMap Kind = "bi.world.areaflags-map"

	// KindMapGroupProto maps to mapgroupproto.xml.
	KindMapGroupProto Kind = "bi.ce.mapgroupproto"

	// KindMapClusterProto maps to mapclusterproto.xml.
	KindMapClusterProto Kind = "bi.ce.mapclusterproto"

	// KindMapGroupPos maps to mapgrouppos.xml.
	KindMapGroupPos Kind = "bi.ce.mapgrouppos"

	// KindMapGroupDirt maps to mapgroupdirt.xml.
	KindMapGroupDirt Kind = "bi.ce.mapgroupdirt"

	// KindMapGroupCluster maps to mapgroupcluster*.xml.
	KindMapGroupCluster Kind = "bi.ce.mapgroupcluster"
)

// Kind describes supported CE configuration file kind.
type Kind string

var (
	// economyCoreIncludeKinds stores strict include type mapping from economycore.
	economyCoreIncludeKinds = map[string]Kind{
		"types":             KindTypes,
		"events":            KindEvents,
		"economy":           KindEconomy,
		"globals":           KindGlobals,
		"messages":          KindMessages,
		"spawnabletypes":    KindSpawnableTypes,
		"cfgspawnabletypes": KindSpawnableTypes,
		"economycore":       KindEconomyCore,
		"cfgeconomycore":    KindEconomyCore,
	}

	// kindAliases stores relaxed user-facing kind aliases to canonical kind ids.
	kindAliases = map[string]Kind{
		"types":                KindTypes,
		"events":               KindEvents,
		"economy":              KindEconomy,
		"globals":              KindGlobals,
		"messages":             KindMessages,
		"spawnabletypes":       KindSpawnableTypes,
		"cfgspawnabletypes":    KindSpawnableTypes,
		"randompresets":        KindRandomPresets,
		"economycore":          KindEconomyCore,
		"cfgeconomycore":       KindEconomyCore,
		"environment":          KindEnvironment,
		"eventspawns":          KindEventSpawns,
		"eventgroups":          KindEventGroups,
		"playerspawnpoints":    KindPlayerSpawnPoints,
		"weather":              KindWeather,
		"limitsdefinition":     KindLimitsDefinition,
		"limitsdefinitionuser": KindLimitsDefinitionUser,
		"ignorelist":           KindIgnoreList,
		"territories":          KindTerritories,
		"undergroundtriggers":  KindUndergroundTriggers,
		"effectarea":           KindEffectArea,
		"gameplay":             KindGameplay,
		"gameplaygearpresets":  KindGameplayGearPresets,
		"objectspawner":        KindObjectSpawner,
		"ceprojectconfig":      KindCEProjectConfig,
		"areaflagsmap":         KindAreaFlagsMap,
		"mapgroupproto":        KindMapGroupProto,
		"mapclusterproto":      KindMapClusterProto,
		"mapgrouppos":          KindMapGroupPos,
		"mapgroupdirt":         KindMapGroupDirt,
		"mapgroupcluster":      KindMapGroupCluster,
	}
)

// KindFromEconomyCoreType maps `<file type="...">` from economycore to Kind.
//
// According to DayZ CE wiki mission file modding docs, include `type` values
// are limited to: types, spawnabletypes, globals, economy, events, messages.
// `economycore` is additionally accepted for recursive include expansion.
func KindFromEconomyCoreType(value string) Kind {
	key := lowerTrim(value)
	kind, ok := economyCoreIncludeKinds[key]
	if !ok {
		return KindUnknown
	}

	return kind
}

// DetectKind returns CE file kind by file base name.
func DetectKind(path string) Kind {
	typ, ok := bimime.DetectByExtension(path)
	if !ok {
		return KindUnknown
	}

	return kindFromTypeID(typ.ID)
}

// kindFromTypeID resolves dzce Kind from bimime type id.
func kindFromTypeID(typeID string) Kind {
	kind := normalizeKindAlias(typeID)
	if kind == KindUnknown {
		return KindUnknown
	}
	if _, ok := DefaultRegistry().Get(kind); !ok {
		return KindUnknown
	}

	return kind
}

// normalizeKindAlias resolves one raw token to canonical kind id or unknown.
func normalizeKindAlias(value string) Kind {
	key := lowerTrim(value)
	if key == "" {
		return KindUnknown
	}

	kind, ok := kindAliases[key]
	if ok {
		return kind
	}

	return Kind(key)
}

// lowerTrim normalizes string token for maps and comparisons.
func lowerTrim(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
