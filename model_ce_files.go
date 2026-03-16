// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "encoding/xml"

// TypesFile is the root of `db/types.xml`.
type TypesFile struct {
	// XMLName is XML root marker for `<types>`.
	XMLName xml.Name `xml:"types" json:"types" yaml:"types"`
	// Types is full item type list managed by CE.
	Types []TypeDef `xml:"type" json:"type,omitempty" yaml:"type,omitempty"`
}

// EventsFile is the root of `db/events.xml`.
type EventsFile struct {
	// XMLName is XML root marker for `<events>`.
	XMLName xml.Name `xml:"events" json:"events" yaml:"events"`
	// Events is dynamic event definitions used by CE.
	Events []EventDef `xml:"event" json:"event,omitempty" yaml:"event,omitempty"`
}

// EconomyFile is the root of `db/economy.xml`.
type EconomyFile struct {
	// Dynamic controls dynamic loot entity persistence.
	Dynamic *EconomySection `xml:"dynamic,omitempty" json:"dynamic,omitempty" yaml:"dynamic,omitempty"`
	// Animals controls wildlife persistence behavior.
	Animals *EconomySection `xml:"animals,omitempty" json:"animals,omitempty" yaml:"animals,omitempty"`
	// Zombies controls infected persistence behavior.
	Zombies *EconomySection `xml:"zombies,omitempty" json:"zombies,omitempty" yaml:"zombies,omitempty"`
	// Vehicles controls vehicle persistence behavior.
	Vehicles *EconomySection `xml:"vehicles,omitempty" json:"vehicles,omitempty" yaml:"vehicles,omitempty"`
	// Randoms controls random-event persistence behavior.
	Randoms *EconomySection `xml:"randoms,omitempty" json:"randoms,omitempty" yaml:"randoms,omitempty"`
	// Custom controls custom CE section persistence.
	Custom *EconomySection `xml:"custom,omitempty" json:"custom,omitempty" yaml:"custom,omitempty"`
	// Building controls building-state persistence.
	Building *EconomySection `xml:"building,omitempty" json:"building,omitempty" yaml:"building,omitempty"`
	// Player controls player-state persistence.
	Player *EconomySection `xml:"player,omitempty" json:"player,omitempty" yaml:"player,omitempty"`
	// XMLName is XML root marker for `<economy>`.
	XMLName xml.Name `xml:"economy" json:"economy" yaml:"economy"`
}

// GlobalsFile is the root of `db/globals.xml`.
type GlobalsFile struct {
	// XMLName is XML root marker for `<variables>`.
	XMLName xml.Name `xml:"variables" json:"variables" yaml:"variables"`
	// Vars is the global CE variable set.
	Vars []GlobalVar `xml:"var" json:"var,omitempty" yaml:"var,omitempty"`
}

// MessagesFile is the root of `db/messages.xml`.
type MessagesFile struct {
	// XMLName is XML root marker for `<messages>`.
	XMLName xml.Name `xml:"messages" json:"messages" yaml:"messages"`
	// Messages is server message schedule entries.
	Messages []MessageDef `xml:"message,omitempty" json:"message,omitempty" yaml:"message,omitempty"`
}

// SpawnableTypesFile is the root of `cfgspawnabletypes.xml`.
type SpawnableTypesFile struct {
	// XMLName is XML root marker for `<spawnabletypes>`.
	XMLName xml.Name `xml:"spawnabletypes" json:"spawnabletypes" yaml:"spawnabletypes"`
	// Types is attachment/cargo spawn rules by item class.
	Types []SpawnableTypeDef `xml:"type" json:"type,omitempty" yaml:"type,omitempty"`
}

// RandomPresetsFile is the root of `cfgrandompresets.xml`.
type RandomPresetsFile struct {
	// XMLName is XML root marker for `<randompresets>`.
	XMLName xml.Name `xml:"randompresets" json:"randompresets" yaml:"randompresets"`
	// Cargo is named random preset set for cargo.
	Cargo []RandomPreset `xml:"cargo,omitempty" json:"cargo,omitempty" yaml:"cargo,omitempty"`
	// Attachments is named random preset set for attachments.
	Attachments []RandomPreset `xml:"attachments,omitempty" json:"attachments,omitempty" yaml:"attachments,omitempty"`
}

// EconomyCoreFile is the root of `cfgeconomycore.xml`.
type EconomyCoreFile struct {
	// XMLName is XML root marker for `<economycore>`.
	XMLName xml.Name `xml:"economycore" json:"economycore" yaml:"economycore"`
	// Classes lists root classes known to CE.
	Classes *EconomyCoreClasses `xml:"classes,omitempty" json:"classes,omitempty" yaml:"classes,omitempty"`
	// Defaults stores default CE variables.
	Defaults *EconomyCoreDefaults `xml:"defaults,omitempty" json:"defaults,omitempty" yaml:"defaults,omitempty"`
	// CE lists external CE config folders and files.
	CE []EconomyCoreCE `xml:"ce,omitempty" json:"ce,omitempty" yaml:"ce,omitempty"`
}

// EnvironmentFile is the root of `cfgenvironment.xml`.
type EnvironmentFile struct {
	// XMLName is XML root marker for `<env>`.
	XMLName struct{} `xml:"env" json:"env" yaml:"env"`
	// Territories stores ambient life territories and sources.
	Territories *EnvironmentTerritories `xml:"territories,omitempty" json:"territories,omitempty" yaml:"territories,omitempty"`
}

// EventSpawnsFile is the root of `cfgeventspawns.xml`.
type EventSpawnsFile struct {
	// XMLName is XML root marker for `<eventposdef>`.
	XMLName xml.Name `xml:"eventposdef" json:"eventposdef" yaml:"eventposdef"`
	// Events is list of event spawn definitions.
	Events []EventSpawnEntry `xml:"event,omitempty" json:"event,omitempty" yaml:"event,omitempty"`
}

// EventGroupsFile is the root of `cfgeventgroups.xml`.
type EventGroupsFile struct {
	// XMLName is XML root marker for `<eventgroupdef>`.
	XMLName xml.Name `xml:"eventgroupdef" json:"eventgroupdef" yaml:"eventgroupdef"`
	// Groups is list of reusable static object groups.
	Groups []EventGroup `xml:"group,omitempty" json:"group,omitempty" yaml:"group,omitempty"`
}

// PlayerSpawnPointsFile is the root of `cfgplayerspawnpoints.xml`.
type PlayerSpawnPointsFile struct {
	// XMLName is XML root marker for `<playerspawnpoints>`.
	XMLName struct{} `xml:"playerspawnpoints" json:"playerspawnpoints" yaml:"playerspawnpoints"`
	// Fresh stores parameters for fresh character spawns.
	Fresh *PlayerSpawnProfile `xml:"fresh,omitempty" json:"fresh,omitempty" yaml:"fresh,omitempty"`
	// Hop stores parameters for server-hop spawns.
	Hop *PlayerSpawnProfile `xml:"hop,omitempty" json:"hop,omitempty" yaml:"hop,omitempty"`
	// Travel stores parameters for travel transfer spawns.
	Travel *PlayerSpawnProfile `xml:"travel,omitempty" json:"travel,omitempty" yaml:"travel,omitempty"`
}

// WeatherFile is the root of `cfgweather.xml`.
type WeatherFile struct {
	// Overcast configures cloudiness channel behavior.
	Overcast *WeatherSection `xml:"overcast,omitempty" json:"overcast,omitempty" yaml:"overcast,omitempty"`
	// Fog configures fog channel behavior.
	Fog *WeatherSection `xml:"fog,omitempty" json:"fog,omitempty" yaml:"fog,omitempty"`
	// Rain configures rain channel behavior.
	Rain *WeatherSection `xml:"rain,omitempty" json:"rain,omitempty" yaml:"rain,omitempty"`
	// WindMagnitude configures wind speed behavior.
	WindMagnitude *WeatherSection `xml:"windMagnitude,omitempty" json:"windMagnitude,omitempty" yaml:"windMagnitude,omitempty"`
	// WindDirection configures wind direction behavior.
	WindDirection *WeatherSection `xml:"windDirection,omitempty" json:"windDirection,omitempty" yaml:"windDirection,omitempty"`
	// Snowfall configures snowfall behavior.
	Snowfall *WeatherSection `xml:"snowfall,omitempty" json:"snowfall,omitempty" yaml:"snowfall,omitempty"`

	// Storm configures storm/thunder behavior.
	Storm *WeatherStorm `xml:"storm,omitempty" json:"storm,omitempty" yaml:"storm,omitempty"`
	// Reset forces weather reset on mission start.
	Reset WeatherToggle `xml:"reset,attr" json:"reset" yaml:"reset"`
	// Enable toggles weather simulation globally.
	Enable WeatherToggle `xml:"enable,attr" json:"enable" yaml:"enable"`
	// XMLName is XML root marker for `<weather>`.
	XMLName struct{} `xml:"weather" json:"weather" yaml:"weather"`
}

// LimitsDefinitionFile is the root of `cfglimitsdefinition.xml`.
type LimitsDefinitionFile struct {
	// XMLName is XML root marker for `<lists>`.
	XMLName struct{} `xml:"lists" json:"lists" yaml:"lists"`
	// Categories is set of allowed item categories.
	Categories *LimitsCategoryList `xml:"categories,omitempty" json:"categories,omitempty" yaml:"categories,omitempty"`
	// Tags is set of allowed item tags.
	Tags *LimitsTagList `xml:"tags,omitempty" json:"tags,omitempty" yaml:"tags,omitempty"`
	// UsageFlags is set of allowed usage flags.
	UsageFlags *LimitsUsageList `xml:"usageflags,omitempty" json:"usageflags,omitempty" yaml:"usageflags,omitempty"`
	// ValueFlags is set of allowed value flags.
	ValueFlags *LimitsValueList `xml:"valueflags,omitempty" json:"valueflags,omitempty" yaml:"valueflags,omitempty"`
}

// LimitsDefinitionUserFile is the root of `cfglimitsdefinitionuser.xml`.
type LimitsDefinitionUserFile struct {
	// XMLName is XML root marker for `<user_lists>`.
	XMLName struct{} `xml:"user_lists" json:"user_lists" yaml:"user_lists"`
	// UsageFlags is alias mapping for usage flags.
	UsageFlags *LimitsUserBindingList `xml:"usageflags,omitempty" json:"usageflags,omitempty" yaml:"usageflags,omitempty"`
	// ValueFlags is alias mapping for value flags.
	ValueFlags *LimitsUserBindingList `xml:"valueflags,omitempty" json:"valueflags,omitempty" yaml:"valueflags,omitempty"`
}

// IgnoreListFile is the root of `cfgignorelist.xml`.
type IgnoreListFile struct {
	// XMLName is XML root marker for `<ignore>`.
	XMLName xml.Name `xml:"ignore" json:"ignore" yaml:"ignore"`
	// Types is list of item class names ignored by CE.
	Types []NamedRef `xml:"type,omitempty" json:"type,omitempty" yaml:"type,omitempty"`
}

// TerritoryFile is the root of `env/*_territories.xml`.
type TerritoryFile struct {
	// XMLName is XML root marker for `<territory-type>`.
	XMLName xml.Name `xml:"territory-type" json:"territory-type" yaml:"territory-type"`
	// Territories is all territory blocks in the file.
	Territories []TerritoryBlock `xml:"territory,omitempty" json:"territory,omitempty" yaml:"territory,omitempty"`
}

// MapGroupProtoFile is the root of `mapgroupproto.xml`.
type MapGroupProtoFile struct {
	// XMLName is XML root marker for `<prototype>`.
	XMLName xml.Name `xml:"prototype" json:"prototype" yaml:"prototype"`
	// Defaults stores shared defaults for groups/containers.
	Defaults *MapPrototypeDefaults `xml:"defaults,omitempty" json:"defaults,omitempty" yaml:"defaults,omitempty"`
	// Clusters stores exported cluster model mappings.
	Clusters *MapPrototypeExports `xml:"clusters,omitempty" json:"clusters,omitempty" yaml:"clusters,omitempty"`
	// Groups stores normal loot map groups.
	Groups []MapPrototypeEntry `xml:"group,omitempty" json:"group,omitempty" yaml:"group,omitempty"`
	// ClusterGroups stores cluster loot groups.
	ClusterGroups []MapPrototypeEntry `xml:"cluster,omitempty" json:"cluster,omitempty" yaml:"cluster,omitempty"`
}

// MapClusterProtoFile is the root of `mapclusterproto.xml`.
type MapClusterProtoFile struct {
	// XMLName is XML root marker for `<prototype>`.
	XMLName xml.Name `xml:"prototype" json:"prototype" yaml:"prototype"`
	// Defaults stores shared defaults for groups/containers.
	Defaults *MapPrototypeDefaults `xml:"defaults,omitempty" json:"defaults,omitempty" yaml:"defaults,omitempty"`
	// Clusters stores exported cluster model mappings.
	Clusters *MapPrototypeExports `xml:"clusters,omitempty" json:"clusters,omitempty" yaml:"clusters,omitempty"`
	// Groups stores normal loot map groups.
	Groups []MapPrototypeEntry `xml:"group,omitempty" json:"group,omitempty" yaml:"group,omitempty"`
	// ClusterGroups stores cluster loot groups.
	ClusterGroups []MapPrototypeEntry `xml:"cluster,omitempty" json:"cluster,omitempty" yaml:"cluster,omitempty"`
}

// MapGroupPosFile is the root of `mapgrouppos.xml`.
type MapGroupPosFile struct {
	// XMLName is XML root marker for `<map>`.
	XMLName xml.Name `xml:"map" json:"map" yaml:"map"`
	// Groups stores map object instances with transform data.
	Groups []MapGroupInstance `xml:"group,omitempty" json:"group,omitempty" yaml:"group,omitempty"`
}

// MapGroupDirtFile is the root of `mapgroupdirt.xml`.
type MapGroupDirtFile struct {
	// XMLName is XML root marker for `<map>`.
	XMLName xml.Name `xml:"map" json:"map" yaml:"map"`
	// Groups stores map object instances with transform data.
	Groups []MapGroupInstance `xml:"group,omitempty" json:"group,omitempty" yaml:"group,omitempty"`
}

// MapGroupClusterFile is the root of `mapgroupcluster*.xml`.
type MapGroupClusterFile struct {
	// XMLName is XML root marker for `<map>`.
	XMLName xml.Name `xml:"map" json:"map" yaml:"map"`
	// Groups stores map cluster instances with transform data.
	Groups []MapGroupInstance `xml:"group,omitempty" json:"group,omitempty" yaml:"group,omitempty"`
}
