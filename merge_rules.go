// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "fmt"

// mergeTypesValue merges `types.xml` by type name with override semantics.
func mergeTypesValue(current any, incoming any) (any, error) {
	left, ok := castValue[TypesFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindTypes)
	}

	right, ok := castValue[TypesFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindTypes)
	}

	for _, inType := range right.Types {
		index := findTypeDefByName(left.Types, inType.Name)
		if index == -1 {
			left.Types = append(left.Types, inType)

			continue
		}

		left.Types[index] = mergeTypeDef(left.Types[index], inType)
	}

	return left, nil
}

// mergeEventsValue merges `events.xml` by event name.
func mergeEventsValue(current any, incoming any) (any, error) {
	left, ok := castValue[EventsFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEvents)
	}

	right, ok := castValue[EventsFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEvents)
	}

	for _, inEvent := range right.Events {
		index := findEventDefByName(left.Events, inEvent.Name)
		if index == -1 {
			left.Events = append(left.Events, inEvent)

			continue
		}

		left.Events[index] = mergeEventDef(left.Events[index], inEvent)
	}

	return left, nil
}

// mergeEconomyValue overlays `economy.xml` section pointers.
func mergeEconomyValue(current any, incoming any) (any, error) {
	left, ok := castValue[EconomyFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEconomy)
	}

	right, ok := castValue[EconomyFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEconomy)
	}

	replaceIfNonNil(&left.Dynamic, right.Dynamic)
	replaceIfNonNil(&left.Animals, right.Animals)
	replaceIfNonNil(&left.Zombies, right.Zombies)
	replaceIfNonNil(&left.Vehicles, right.Vehicles)
	replaceIfNonNil(&left.Randoms, right.Randoms)
	replaceIfNonNil(&left.Custom, right.Custom)
	replaceIfNonNil(&left.Building, right.Building)
	replaceIfNonNil(&left.Player, right.Player)

	return left, nil
}

// mergeGlobalsValue merges globals variables by name.
func mergeGlobalsValue(current any, incoming any) (any, error) {
	left, ok := castValue[GlobalsFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindGlobals)
	}

	right, ok := castValue[GlobalsFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindGlobals)
	}

	for _, inVar := range right.Vars {
		index := findGlobalVarByName(left.Vars, inVar.Name)
		if index == -1 {
			left.Vars = append(left.Vars, inVar)

			continue
		}

		if left.Vars[index].Type != inVar.Type {
			return nil, fmt.Errorf(
				"%w: globals variable %q type mismatch (%d != %d)",
				ErrMergeConflict,
				inVar.Name,
				left.Vars[index].Type,
				inVar.Type,
			)
		}

		left.Vars[index].Value = inVar.Value
	}

	return left, nil
}

// mergeMessagesValue appends messages entries.
func mergeMessagesValue(current any, incoming any) (any, error) {
	left, ok := castValue[MessagesFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMessages)
	}

	right, ok := castValue[MessagesFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMessages)
	}

	left.Messages = append(left.Messages, right.Messages...)

	return left, nil
}

// mergeSpawnableTypesValue merges `cfgspawnabletypes.xml` by type name.
func mergeSpawnableTypesValue(current any, incoming any) (any, error) {
	left, ok := castValue[SpawnableTypesFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindSpawnableTypes)
	}

	right, ok := castValue[SpawnableTypesFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindSpawnableTypes)
	}

	for _, inType := range right.Types {
		index := findSpawnableTypeDefByName(left.Types, inType.Name)
		if index == -1 {
			left.Types = append(left.Types, inType)

			continue
		}

		left.Types[index] = mergeSpawnableTypeDef(left.Types[index], inType)
	}

	return left, nil
}

// findTypeDefByName returns index of type entry by name.
func findTypeDefByName(items []TypeDef, name string) int {
	for index := range items {
		if items[index].Name == name {
			return index
		}
	}

	return -1
}

// mergeTypeDef overlays one type definition with include payload values.
func mergeTypeDef(current TypeDef, incoming TypeDef) TypeDef {
	merged := current
	merged.Name = incoming.Name

	replaceIfNonNil(&merged.Nominal, incoming.Nominal)
	replaceIfNonNil(&merged.Lifetime, incoming.Lifetime)
	replaceIfNonNil(&merged.Restock, incoming.Restock)
	replaceIfNonNil(&merged.Min, incoming.Min)
	replaceIfNonNil(&merged.QuantityMin, incoming.QuantityMin)
	replaceIfNonNil(&merged.QuantityMax, incoming.QuantityMax)
	replaceIfNonNil(&merged.Cost, incoming.Cost)
	replaceIfNonNil(&merged.Flags, incoming.Flags)

	if len(incoming.Categories) > 0 {
		merged.Categories = incoming.Categories
	}

	if len(incoming.Usages) > 0 {
		merged.Usages = incoming.Usages
	}

	if len(incoming.Tags) > 0 {
		merged.Tags = incoming.Tags
	}

	if len(incoming.Values) > 0 {
		merged.Values = incoming.Values
	}

	return merged
}

// findEventDefByName returns index of event entry by name.
func findEventDefByName(items []EventDef, name string) int {
	for index := range items {
		if items[index].Name == name {
			return index
		}
	}

	return -1
}

// mergeEventDef overlays one event definition with include payload values.
func mergeEventDef(current EventDef, incoming EventDef) EventDef {
	merged := current
	merged.Name = incoming.Name

	replaceIfNonNil(&merged.Nominal, incoming.Nominal)
	replaceIfNonNil(&merged.Min, incoming.Min)
	replaceIfNonNil(&merged.Max, incoming.Max)
	replaceIfNonNil(&merged.Lifetime, incoming.Lifetime)
	replaceIfNonNil(&merged.Restock, incoming.Restock)
	replaceIfNonNil(&merged.SafeRadius, incoming.SafeRadius)
	replaceIfNonNil(&merged.DistanceRadius, incoming.DistanceRadius)
	replaceIfNonNil(&merged.CleanupRadius, incoming.CleanupRadius)
	replaceIfNonNil(&merged.Secondary, incoming.Secondary)
	replaceIfNonNil(&merged.Flags, incoming.Flags)
	replaceIfNonNil(&merged.Position, incoming.Position)
	replaceIfNonNil(&merged.Limit, incoming.Limit)
	replaceIfNonNil(&merged.Active, incoming.Active)

	if incoming.Children == nil {
		return merged
	}

	if merged.Children == nil {
		merged.Children = &EventChildren{}
	}

	for _, inChild := range incoming.Children.Children {
		childIndex := findEventChildByType(merged.Children.Children, inChild.Type)
		if childIndex == -1 {
			merged.Children.Children = append(merged.Children.Children, inChild)

			continue
		}

		merged.Children.Children[childIndex] = mergeEventChild(
			merged.Children.Children[childIndex],
			inChild,
		)
	}

	return merged
}

// findEventChildByType returns child index by type name.
func findEventChildByType(items []EventChild, name string) int {
	for index := range items {
		if items[index].Type == name {
			return index
		}
	}

	return -1
}

// mergeEventChild overlays one event child entry.
func mergeEventChild(current EventChild, incoming EventChild) EventChild {
	merged := current
	merged.Type = incoming.Type

	replaceIfNonNil(&merged.LootMax, incoming.LootMax)
	replaceIfNonNil(&merged.LootMin, incoming.LootMin)
	replaceIfNonNil(&merged.Max, incoming.Max)
	replaceIfNonNil(&merged.Min, incoming.Min)

	return merged
}

// findGlobalVarByName returns index of globals variable by name.
func findGlobalVarByName(items []GlobalVar, name string) int {
	for index := range items {
		if items[index].Name == name {
			return index
		}
	}

	return -1
}

// findSpawnableTypeDefByName returns index of spawnable type entry by name.
func findSpawnableTypeDefByName(items []SpawnableTypeDef, name string) int {
	for index := range items {
		if items[index].Name == name {
			return index
		}
	}

	return -1
}

// mergeSpawnableTypeDef overlays one spawnable type definition.
func mergeSpawnableTypeDef(
	current SpawnableTypeDef,
	incoming SpawnableTypeDef,
) SpawnableTypeDef {
	merged := current
	merged.Name = incoming.Name

	replaceIfNonNil(&merged.Damage, incoming.Damage)

	if incoming.Hoarder != nil || incoming.Unique != nil {
		merged.Hoarder = nil
		merged.Unique = nil
		merged.Hoarder = incoming.Hoarder
		merged.Unique = incoming.Unique
	}

	if len(incoming.Tags) > 0 {
		merged.Tags = incoming.Tags
	}

	if len(incoming.Cargo) > 0 {
		merged.Cargo = incoming.Cargo
	}

	if len(incoming.Attachments) > 0 {
		merged.Attachments = incoming.Attachments
	}

	return merged
}

// mergeRandomPresetsValue appends random preset entries.
func mergeRandomPresetsValue(current any, incoming any) (any, error) {
	left, ok := castValue[RandomPresetsFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindRandomPresets)
	}

	right, ok := castValue[RandomPresetsFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindRandomPresets)
	}

	left.Cargo = append(left.Cargo, right.Cargo...)
	left.Attachments = append(left.Attachments, right.Attachments...)

	return left, nil
}

// mergeEnvironmentValue appends environment territories and file refs.
func mergeEnvironmentValue(current any, incoming any) (any, error) {
	left, ok := castValue[EnvironmentFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEnvironment)
	}

	right, ok := castValue[EnvironmentFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEnvironment)
	}

	if right.Territories == nil {
		return left, nil
	}

	if left.Territories == nil {
		left.Territories = &EnvironmentTerritories{}
	}

	left.Territories.Files = append(left.Territories.Files, right.Territories.Files...)
	left.Territories.Territories = append(
		left.Territories.Territories,
		right.Territories.Territories...,
	)

	return left, nil
}

// mergeEventSpawnsValue appends event position entries.
func mergeEventSpawnsValue(current any, incoming any) (any, error) {
	left, ok := castValue[EventSpawnsFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEventSpawns)
	}

	right, ok := castValue[EventSpawnsFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEventSpawns)
	}

	left.Events = append(left.Events, right.Events...)

	return left, nil
}

// mergeEventGroupsValue appends event group entries.
func mergeEventGroupsValue(current any, incoming any) (any, error) {
	left, ok := castValue[EventGroupsFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEventGroups)
	}

	right, ok := castValue[EventGroupsFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEventGroups)
	}

	left.Groups = append(left.Groups, right.Groups...)

	return left, nil
}

// mergePlayerSpawnPointsValue overlays player spawn profiles.
func mergePlayerSpawnPointsValue(current any, incoming any) (any, error) {
	left, ok := castValue[PlayerSpawnPointsFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindPlayerSpawnPoints)
	}

	right, ok := castValue[PlayerSpawnPointsFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindPlayerSpawnPoints)
	}

	if right.Fresh != nil {
		left.Fresh = right.Fresh
	}

	if right.Hop != nil {
		left.Hop = right.Hop
	}

	if right.Travel != nil {
		left.Travel = right.Travel
	}

	return left, nil
}

// mergeWeatherValue overlays weather branches.
func mergeWeatherValue(current any, incoming any) (any, error) {
	left, ok := castValue[WeatherFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindWeather)
	}

	right, ok := castValue[WeatherFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindWeather)
	}

	if right.Overcast != nil {
		left.Overcast = right.Overcast
	}

	if right.Fog != nil {
		left.Fog = right.Fog
	}

	if right.Rain != nil {
		left.Rain = right.Rain
	}

	if right.WindMagnitude != nil {
		left.WindMagnitude = right.WindMagnitude
	}

	if right.WindDirection != nil {
		left.WindDirection = right.WindDirection
	}

	if right.Snowfall != nil {
		left.Snowfall = right.Snowfall
	}

	if right.Storm != nil {
		left.Storm = right.Storm
	}

	left.Reset = right.Reset
	left.Enable = right.Enable

	return left, nil
}

// mergeLimitsDefinitionValue appends limits categories/tags/usages/values.
func mergeLimitsDefinitionValue(current any, incoming any) (any, error) {
	left, ok := castValue[LimitsDefinitionFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindLimitsDefinition)
	}

	right, ok := castValue[LimitsDefinitionFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindLimitsDefinition)
	}

	if right.Categories != nil {
		if left.Categories == nil {
			left.Categories = &LimitsCategoryList{}
		}

		left.Categories.Categories = append(
			left.Categories.Categories,
			right.Categories.Categories...,
		)
	}

	if right.Tags != nil {
		if left.Tags == nil {
			left.Tags = &LimitsTagList{}
		}

		left.Tags.Tags = append(left.Tags.Tags, right.Tags.Tags...)
	}

	if right.UsageFlags != nil {
		if left.UsageFlags == nil {
			left.UsageFlags = &LimitsUsageList{}
		}

		left.UsageFlags.Usages = append(
			left.UsageFlags.Usages,
			right.UsageFlags.Usages...,
		)
	}

	if right.ValueFlags != nil {
		if left.ValueFlags == nil {
			left.ValueFlags = &LimitsValueList{}
		}

		left.ValueFlags.Values = append(
			left.ValueFlags.Values,
			right.ValueFlags.Values...,
		)
	}

	return left, nil
}

// mergeLimitsDefinitionUserValue appends user alias usage/value entries.
func mergeLimitsDefinitionUserValue(current any, incoming any) (any, error) {
	left, ok := castValue[LimitsDefinitionUserFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindLimitsDefinitionUser)
	}

	right, ok := castValue[LimitsDefinitionUserFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindLimitsDefinitionUser)
	}

	if right.UsageFlags != nil {
		if left.UsageFlags == nil {
			left.UsageFlags = &LimitsUserBindingList{}
		}

		left.UsageFlags.Users = append(left.UsageFlags.Users, right.UsageFlags.Users...)
	}

	if right.ValueFlags != nil {
		if left.ValueFlags == nil {
			left.ValueFlags = &LimitsUserBindingList{}
		}

		left.ValueFlags.Users = append(left.ValueFlags.Users, right.ValueFlags.Users...)
	}

	return left, nil
}

// mergeIgnoreListValue appends ignored class names.
func mergeIgnoreListValue(current any, incoming any) (any, error) {
	left, ok := castValue[IgnoreListFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindIgnoreList)
	}

	right, ok := castValue[IgnoreListFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindIgnoreList)
	}

	left.Types = append(left.Types, right.Types...)

	return left, nil
}

// mergeTerritoryValue appends territory blocks.
func mergeTerritoryValue(current any, incoming any) (any, error) {
	left, ok := castValue[TerritoryFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindTerritories)
	}

	right, ok := castValue[TerritoryFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindTerritories)
	}

	left.Territories = append(left.Territories, right.Territories...)

	return left, nil
}

// mergeUndergroundTriggersValue appends underground trigger entries.
func mergeUndergroundTriggersValue(current any, incoming any) (any, error) {
	left, ok := castValue[UndergroundTriggersFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindUndergroundTriggers)
	}

	right, ok := castValue[UndergroundTriggersFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindUndergroundTriggers)
	}

	left.Triggers = append(left.Triggers, right.Triggers...)

	return left, nil
}

// mergeEffectAreaValue appends effect areas and safe positions.
func mergeEffectAreaValue(current any, incoming any) (any, error) {
	left, ok := castValue[EffectAreaFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEffectArea)
	}

	right, ok := castValue[EffectAreaFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindEffectArea)
	}

	left.Areas = append(left.Areas, right.Areas...)
	left.SafePositions = append(left.SafePositions, right.SafePositions...)

	return left, nil
}

// mergeGameplayValue overlays top-level cfggameplay branches.
func mergeGameplayValue(current any, incoming any) (any, error) {
	left, ok := castValue[GameplayFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindGameplay)
	}

	right, ok := castValue[GameplayFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindGameplay)
	}

	replaceIfNonNil(&left.Version, right.Version)
	replaceIfNonNil(&left.GeneralData, right.GeneralData)
	replaceIfNonNil(&left.PlayerData, right.PlayerData)
	replaceIfNonNil(&left.WorldsData, right.WorldsData)
	replaceIfNonNil(&left.BaseBuildingData, right.BaseBuildingData)
	replaceIfNonNil(&left.UIData, right.UIData)
	replaceIfNonNil(&left.MapData, right.MapData)
	replaceIfNonNil(&left.VehicleData, right.VehicleData)

	return left, nil
}

// mergeGameplayGearPresetsValue appends gameplay gear presets.
func mergeGameplayGearPresetsValue(current any, incoming any) (any, error) {
	left, ok := castValue[GameplayGearPresetsFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindGameplayGearPresets)
	}

	right, ok := castValue[GameplayGearPresetsFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindGameplayGearPresets)
	}

	*left = append(*left, (*right)...)
	return left, nil
}

// mergeObjectSpawnerValue appends object spawner entries.
func mergeObjectSpawnerValue(current any, incoming any) (any, error) {
	left, ok := castValue[ObjectSpawnerFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindObjectSpawner)
	}

	right, ok := castValue[ObjectSpawnerFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindObjectSpawner)
	}

	left.Objects = append(left.Objects, right.Objects...)
	return left, nil
}

// mergeCEProjectConfigValue replaces previous CEProject config with incoming payload.
func mergeCEProjectConfigValue(current any, incoming any) (any, error) {
	_, ok := castValue[CEProjectConfigFile](current)
	if !ok {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindCEProjectConfig)
	}

	right, ok := castValue[CEProjectConfigFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindCEProjectConfig)
	}

	return right, nil
}

// mergeAreaFlagsMapValue replaces previous `areaflags.map` payload.
func mergeAreaFlagsMapValue(current any, incoming any) (any, error) {
	_, ok := castValue[AreaFlagsMapFile](current)
	if !ok {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindAreaFlagsMap)
	}

	right, ok := castValue[AreaFlagsMapFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindAreaFlagsMap)
	}

	return right, nil
}

// mergeMapGroupProtoValue appends map group prototype data.
func mergeMapGroupProtoValue(current any, incoming any) (any, error) {
	left, ok := castValue[MapGroupProtoFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMapGroupProto)
	}

	right, ok := castValue[MapGroupProtoFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMapGroupProto)
	}

	mergeMapPrototypeContent(left, right)

	return left, nil
}

// mergeMapClusterProtoValue appends map cluster prototype data.
func mergeMapClusterProtoValue(current any, incoming any) (any, error) {
	left, ok := castValue[MapClusterProtoFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMapClusterProto)
	}

	right, ok := castValue[MapClusterProtoFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMapClusterProto)
	}

	mergeMapPrototypeContent(left, right)

	return left, nil
}

// mergeMapGroupPosValue appends map group position entries.
func mergeMapGroupPosValue(current any, incoming any) (any, error) {
	left, ok := castValue[MapGroupPosFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMapGroupPos)
	}

	right, ok := castValue[MapGroupPosFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMapGroupPos)
	}

	left.Groups = append(left.Groups, right.Groups...)

	return left, nil
}

// mergeMapGroupDirtValue appends map group dirt entries.
func mergeMapGroupDirtValue(current any, incoming any) (any, error) {
	left, ok := castValue[MapGroupDirtFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMapGroupDirt)
	}

	right, ok := castValue[MapGroupDirtFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMapGroupDirt)
	}

	left.Groups = append(left.Groups, right.Groups...)

	return left, nil
}

// mergeMapGroupClusterValue appends map group cluster entries.
func mergeMapGroupClusterValue(current any, incoming any) (any, error) {
	left, ok := castValue[MapGroupClusterFile](current)
	if !ok || left == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMapGroupCluster)
	}

	right, ok := castValue[MapGroupClusterFile](incoming)
	if !ok || right == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindMapGroupCluster)
	}

	left.Groups = append(left.Groups, right.Groups...)

	return left, nil
}

// mapPrototypeAccess provides a shared interface for map prototype merge.
type mapPrototypeAccess interface {
	getDefaults() *MapPrototypeDefaults
	setDefaults(*MapPrototypeDefaults)
	getClusters() *MapPrototypeExports
	setClusters(*MapPrototypeExports)
	getGroups() []MapPrototypeEntry
	setGroups([]MapPrototypeEntry)
	getClusterGroups() []MapPrototypeEntry
	setClusterGroups([]MapPrototypeEntry)
}

// getDefaults returns defaults pointer.
func (file *MapGroupProtoFile) getDefaults() *MapPrototypeDefaults {
	return file.Defaults
}

// setDefaults updates defaults pointer.
func (file *MapGroupProtoFile) setDefaults(value *MapPrototypeDefaults) {
	file.Defaults = value
}

// getClusters returns cluster export mappings.
func (file *MapGroupProtoFile) getClusters() *MapPrototypeExports {
	return file.Clusters
}

// setClusters updates cluster export mappings.
func (file *MapGroupProtoFile) setClusters(value *MapPrototypeExports) {
	file.Clusters = value
}

// getGroups returns prototype groups.
func (file *MapGroupProtoFile) getGroups() []MapPrototypeEntry {
	return file.Groups
}

// setGroups updates prototype groups.
func (file *MapGroupProtoFile) setGroups(value []MapPrototypeEntry) {
	file.Groups = value
}

// getClusterGroups returns prototype cluster groups.
func (file *MapGroupProtoFile) getClusterGroups() []MapPrototypeEntry {
	return file.ClusterGroups
}

// setClusterGroups updates prototype cluster groups.
func (file *MapGroupProtoFile) setClusterGroups(value []MapPrototypeEntry) {
	file.ClusterGroups = value
}

// getDefaults returns defaults pointer.
func (file *MapClusterProtoFile) getDefaults() *MapPrototypeDefaults {
	return file.Defaults
}

// setDefaults updates defaults pointer.
func (file *MapClusterProtoFile) setDefaults(value *MapPrototypeDefaults) {
	file.Defaults = value
}

// getClusters returns cluster export mappings.
func (file *MapClusterProtoFile) getClusters() *MapPrototypeExports {
	return file.Clusters
}

// setClusters updates cluster export mappings.
func (file *MapClusterProtoFile) setClusters(value *MapPrototypeExports) {
	file.Clusters = value
}

// getGroups returns prototype groups.
func (file *MapClusterProtoFile) getGroups() []MapPrototypeEntry {
	return file.Groups
}

// setGroups updates prototype groups.
func (file *MapClusterProtoFile) setGroups(value []MapPrototypeEntry) {
	file.Groups = value
}

// getClusterGroups returns prototype cluster groups.
func (file *MapClusterProtoFile) getClusterGroups() []MapPrototypeEntry {
	return file.ClusterGroups
}

// setClusterGroups updates prototype cluster groups.
func (file *MapClusterProtoFile) setClusterGroups(value []MapPrototypeEntry) {
	file.ClusterGroups = value
}

// mergeMapPrototypeContent merges fields shared by map prototype files.
func mergeMapPrototypeContent[T mapPrototypeAccess](left T, right T) {
	if right.getDefaults() != nil {
		if left.getDefaults() == nil {
			left.setDefaults(&MapPrototypeDefaults{})
		}

		leftDefaults := left.getDefaults()
		leftDefaults.Defaults = append(
			leftDefaults.Defaults,
			right.getDefaults().Defaults...,
		)
	}

	if right.getClusters() != nil {
		if left.getClusters() == nil {
			left.setClusters(&MapPrototypeExports{})
		}

		leftClusters := left.getClusters()
		leftClusters.Exports = append(
			leftClusters.Exports,
			right.getClusters().Exports...,
		)
	}

	left.setGroups(append(left.getGroups(), right.getGroups()...))
	left.setClusterGroups(append(left.getClusterGroups(), right.getClusterGroups()...))
}

// replaceIfNonNil overwrites target pointer when incoming pointer is non-nil.
func replaceIfNonNil[T any](target **T, incoming *T) {
	if incoming == nil {
		return
	}

	*target = incoming
}
