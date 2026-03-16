// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "fmt"

// mergeValueFunc is one per-kind merge implementation.
type mergeValueFunc func(current any, incoming any) (any, error)

var mergeHandlers = map[Kind]mergeValueFunc{
	KindTypes:                mergeTypesValue,
	KindEvents:               mergeEventsValue,
	KindEconomy:              mergeEconomyValue,
	KindGlobals:              mergeGlobalsValue,
	KindMessages:             mergeMessagesValue,
	KindSpawnableTypes:       mergeSpawnableTypesValue,
	KindRandomPresets:        mergeRandomPresetsValue,
	KindEnvironment:          mergeEnvironmentValue,
	KindEventSpawns:          mergeEventSpawnsValue,
	KindEventGroups:          mergeEventGroupsValue,
	KindPlayerSpawnPoints:    mergePlayerSpawnPointsValue,
	KindWeather:              mergeWeatherValue,
	KindLimitsDefinition:     mergeLimitsDefinitionValue,
	KindLimitsDefinitionUser: mergeLimitsDefinitionUserValue,
	KindIgnoreList:           mergeIgnoreListValue,
	KindTerritories:          mergeTerritoryValue,
	KindUndergroundTriggers:  mergeUndergroundTriggersValue,
	KindEffectArea:           mergeEffectAreaValue,
	KindGameplay:             mergeGameplayValue,
	KindGameplayGearPresets:  mergeGameplayGearPresetsValue,
	KindObjectSpawner:        mergeObjectSpawnerValue,
	KindCEProjectConfig:      mergeCEProjectConfigValue,
	KindAreaFlagsMap:         mergeAreaFlagsMapValue,
	KindMapGroupProto:        mergeMapGroupProtoValue,
	KindMapClusterProto:      mergeMapClusterProtoValue,
	KindMapGroupPos:          mergeMapGroupPosValue,
	KindMapGroupDirt:         mergeMapGroupDirtValue,
	KindMapGroupCluster:      mergeMapGroupClusterValue,
}

// mergeKindValue merges source payload into already merged payload by kind.
func mergeKindValue(kind Kind, current any, incoming any) (any, error) {
	if current == nil {
		return incoming, nil
	}

	handler, ok := mergeHandlers[kind]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedKind, kind)
	}

	return handler(current, incoming)
}
