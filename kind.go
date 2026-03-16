// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"path/filepath"
	"strings"
)

const (
	// KindUnknown is used when file kind cannot be detected.
	KindUnknown Kind = ""

	// KindTypes maps to db/types.xml.
	KindTypes Kind = "types"

	// KindEvents maps to db/events.xml.
	KindEvents Kind = "events"

	// KindEconomy maps to db/economy.xml.
	KindEconomy Kind = "economy"

	// KindGlobals maps to db/globals.xml.
	KindGlobals Kind = "globals"

	// KindMessages maps to db/messages.xml.
	KindMessages Kind = "messages"

	// KindSpawnableTypes maps to cfgspawnabletypes.xml.
	KindSpawnableTypes Kind = "spawnabletypes"

	// KindRandomPresets maps to cfgrandompresets.xml.
	KindRandomPresets Kind = "randompresets"

	// KindEconomyCore maps to cfgeconomycore.xml.
	KindEconomyCore Kind = "economycore"

	// KindEnvironment maps to cfgenvironment.xml.
	KindEnvironment Kind = "environment"

	// KindEventSpawns maps to cfgeventspawns.xml.
	KindEventSpawns Kind = "eventspawns"

	// KindEventGroups maps to cfgeventgroups.xml.
	KindEventGroups Kind = "eventgroups"

	// KindPlayerSpawnPoints maps to cfgplayerspawnpoints.xml.
	KindPlayerSpawnPoints Kind = "playerspawnpoints"

	// KindWeather maps to cfgweather.xml.
	KindWeather Kind = "weather"

	// KindLimitsDefinition maps to cfglimitsdefinition.xml.
	KindLimitsDefinition Kind = "limitsdefinition"

	// KindLimitsDefinitionUser maps to cfglimitsdefinitionuser.xml.
	KindLimitsDefinitionUser Kind = "limitsdefinitionuser"

	// KindIgnoreList maps to cfgignorelist.xml.
	KindIgnoreList Kind = "ignorelist"

	// KindTerritories maps to env/*_territories.xml files.
	KindTerritories Kind = "territories"

	// KindUndergroundTriggers maps to cfgundergroundtriggers.json.
	KindUndergroundTriggers Kind = "undergroundtriggers"

	// KindEffectArea maps to cfgeffectarea.json.
	KindEffectArea Kind = "effectarea"

	// KindGameplay maps to cfggameplay.json.
	KindGameplay Kind = "gameplay"

	// KindGameplayGearPresets maps to JSON payloads listed in
	// `cfggameplay.json -> PlayerData.spawnGearPresetFiles`.
	KindGameplayGearPresets Kind = "gameplaygearpresets"

	// KindObjectSpawner maps to JSON payloads listed in
	// `cfggameplay.json -> WorldsData.objectSpawnersArr`.
	KindObjectSpawner Kind = "objectspawner"

	// KindCEProjectConfig maps to CEProject `mapname.xml` (`<zg-config>`).
	KindCEProjectConfig Kind = "ceprojectconfig"

	// KindAreaFlagsMap maps to `areaflags.map` binary payload.
	KindAreaFlagsMap Kind = "areaflagsmap"

	// KindMapGroupProto maps to mapgroupproto.xml.
	KindMapGroupProto Kind = "mapgroupproto"

	// KindMapClusterProto maps to mapclusterproto.xml.
	KindMapClusterProto Kind = "mapclusterproto"

	// KindMapGroupPos maps to mapgrouppos.xml.
	KindMapGroupPos Kind = "mapgrouppos"

	// KindMapGroupDirt maps to mapgroupdirt.xml.
	KindMapGroupDirt Kind = "mapgroupdirt"

	// KindMapGroupCluster maps to mapgroupcluster*.xml.
	KindMapGroupCluster Kind = "mapgroupcluster"
)

// Kind describes supported CE configuration file kind.
type Kind string

// KindFromEconomyCoreType maps `<file type="...">` from economycore to Kind.
//
// According to DayZ CE wiki mission file modding docs, include `type` values
// are limited to: types, spawnabletypes, globals, economy, events, messages.
// `economycore` is additionally accepted for recursive include expansion.
func KindFromEconomyCoreType(value string) Kind {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "types":
		return KindTypes
	case "events":
		return KindEvents
	case "economy":
		return KindEconomy
	case "globals":
		return KindGlobals
	case "messages":
		return KindMessages
	case "spawnabletypes", "cfgspawnabletypes":
		return KindSpawnableTypes
	case "economycore", "cfgeconomycore":
		return KindEconomyCore
	default:
		return KindUnknown
	}
}

// DetectKind returns CE file kind by file base name.
func DetectKind(path string) Kind {
	base := strings.ToLower(filepath.Base(path))

	switch base {
	case "types.xml":
		return KindTypes
	case "events.xml":
		return KindEvents
	case "economy.xml":
		return KindEconomy
	case "globals.xml":
		return KindGlobals
	case "messages.xml":
		return KindMessages
	case "cfgspawnabletypes.xml":
		return KindSpawnableTypes
	case "cfgrandompresets.xml":
		return KindRandomPresets
	case "cfgeconomycore.xml":
		return KindEconomyCore
	case "cfgenvironment.xml":
		return KindEnvironment
	case "cfgeventspawns.xml":
		return KindEventSpawns
	case "cfgeventgroups.xml":
		return KindEventGroups
	case "cfgplayerspawnpoints.xml":
		return KindPlayerSpawnPoints
	case "cfgweather.xml":
		return KindWeather
	case "cfglimitsdefinition.xml":
		return KindLimitsDefinition
	case "cfglimitsdefinitionuser.xml":
		return KindLimitsDefinitionUser
	case "cfgignorelist.xml":
		return KindIgnoreList
	case "cfgundergroundtriggers.json":
		return KindUndergroundTriggers
	case "cfgeffectarea.json":
		return KindEffectArea
	case "cfggameplay.json":
		return KindGameplay
	case "areaflags.map":
		return KindAreaFlagsMap
	case "mapgroupproto.xml":
		return KindMapGroupProto
	case "mapclusterproto.xml":
		return KindMapClusterProto
	case "mapgrouppos.xml":
		return KindMapGroupPos
	case "mapgroupdirt.xml":
		return KindMapGroupDirt
	case "mapgroupcluster.xml":
		return KindMapGroupCluster
	default:
		if strings.HasPrefix(base, "mapgroupcluster") &&
			strings.HasSuffix(base, ".xml") {
			return KindMapGroupCluster
		}

		if strings.HasSuffix(base, "_territories.xml") {
			return KindTerritories
		}

		return KindUnknown
	}
}
