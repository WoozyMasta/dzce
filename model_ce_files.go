// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "encoding/xml"

// TypesFile is the Central Economy item catalog in `db/types.xml`.
// Each entry sets one class's target counts, spawn filters, and replenishment behavior.
type TypesFile struct {
	// XMLName identifies the required `<types>` document root.
	XMLName xml.Name `json:"types" yaml:"types" xml:"types"`
	// Types are the item classes whose loot economy rules are defined here.
	Types []TypeDef `json:"type,omitempty" yaml:"type,omitempty" xml:"type"`
}

// EventsFile configures Central Economy dynamic events such as vehicle spawns, crash sites, and ambient animal zones.
type EventsFile struct {
	// XMLName identifies the required `<events>` document root.
	XMLName xml.Name `json:"events" yaml:"events" xml:"events"`
	// Events are the dynamic event definitions processed by CE.
	Events []EventDef `json:"event,omitempty" yaml:"event,omitempty" xml:"event"`
}

// EconomyFile controls initialization, persistence, and respawning for Central Economy subsystems in `db/economy.xml`.
type EconomyFile struct {
	// Dynamic controls CE initialization, persistence, and respawning for dynamic economy entities such as world loot.
	Dynamic *EconomySection `json:"dynamic,omitempty" yaml:"dynamic,omitempty" xml:"dynamic,omitempty"`
	// Animals controls initialization and persistence of non-ambient wildlife.
	Animals *EconomySection `json:"animals,omitempty" yaml:"animals,omitempty" xml:"animals,omitempty"`
	// Zombies controls initialization and persistence of infected entities.
	Zombies *EconomySection `json:"zombies,omitempty" yaml:"zombies,omitempty" xml:"zombies,omitempty"`
	// Vehicles controls initialization and persistence of CE-managed vehicles.
	Vehicles *EconomySection `json:"vehicles,omitempty" yaml:"vehicles,omitempty" xml:"vehicles,omitempty"`
	// Randoms controls initialization and persistence for random economy entities.
	Randoms *EconomySection `json:"randoms,omitempty" yaml:"randoms,omitempty" xml:"randoms,omitempty"`
	// Custom controls initialization and persistence for custom CE entities.
	Custom *EconomySection `json:"custom,omitempty" yaml:"custom,omitempty" xml:"custom,omitempty"`
	// Building controls initialization and persistence of building state, such as randomized or saved door states.
	Building *EconomySection `json:"building,omitempty" yaml:"building,omitempty" xml:"building,omitempty"`
	// Player controls initialization and persistence of player CE state.
	Player *EconomySection `json:"player,omitempty" yaml:"player,omitempty" xml:"player,omitempty"`
	// XMLName identifies the required `<economy>` document root.
	XMLName xml.Name `json:"economy" yaml:"economy" xml:"economy"`
}

// GlobalsFile sets mission-wide Central Economy variables in `db/globals.xml`,
// including spawn limits, cleanup timing, and persistence-related behavior.
type GlobalsFile struct {
	// XMLName identifies the required `<variables>` document root.
	XMLName xml.Name `json:"variables" yaml:"variables" xml:"variables"`
	// Vars are the CE variables overridden by this mission.
	Vars []GlobalVar `json:"var,omitempty" yaml:"var,omitempty" xml:"var"`
}

// MessagesFile is the root of `db/messages.xml`.
type MessagesFile struct {
	// XMLName is XML root marker for `<messages>`.
	XMLName xml.Name `json:"messages" yaml:"messages" xml:"messages"`
	// Messages is server message schedule entries.
	Messages []MessageDef `json:"message,omitempty" yaml:"message,omitempty" xml:"message,omitempty"`
}

// SpawnableTypesFile defines which cargo items and attachments may be created when CE spawns an item, vehicle, or other configured class.
type SpawnableTypesFile struct {
	// XMLName identifies the required `<spawnabletypes>` document root.
	XMLName xml.Name `json:"spawnabletypes" yaml:"spawnabletypes" xml:"spawnabletypes"`
	// Types map game classes to their cargo and attachment generation rules.
	Types []SpawnableTypeDef `json:"type,omitempty" yaml:"type,omitempty" xml:"type"`
}

// RandomPresetsFile defines reusable cargo and attachment candidate lists for `cfgspawnabletypes.xml`.
type RandomPresetsFile struct {
	// XMLName identifies the required `<randompresets>` document root.
	XMLName xml.Name `json:"randompresets" yaml:"randompresets" xml:"randompresets"`
	// Cargo contains presets referenced by cargo rules.
	Cargo []RandomPreset `json:"cargo,omitempty" yaml:"cargo,omitempty" xml:"cargo,omitempty"`
	// Attachments contains presets referenced by attachment rules.
	Attachments []RandomPreset `json:"attachments,omitempty" yaml:"attachments,omitempty" xml:"attachments,omitempty"`
}

// EconomyCoreFile configures the CE class registry, default variables, and mission-provided CE file additions or overrides.
type EconomyCoreFile struct {
	// Classes declares the game class hierarchies tracked by CE.
	Classes *EconomyCoreClasses `json:"classes,omitempty" yaml:"classes,omitempty" xml:"classes,omitempty"`
	// Defaults overrides built-in CE default variables for this mission.
	Defaults *EconomyCoreDefaults `json:"defaults,omitempty" yaml:"defaults,omitempty" xml:"defaults,omitempty"`
	// XMLName identifies the required `<economycore>` document root.
	XMLName xml.Name `json:"economycore" yaml:"economycore" xml:"economycore"`
	// CE registers mission-relative folders containing CE files to append or use as overrides.
	CE []EconomyCoreCE `json:"ce,omitempty" yaml:"ce,omitempty" xml:"ce,omitempty"`
}

// EnvironmentFile configures territory profiles, their spawn behaviors, and the files that provide their zones in `cfgenvironment.xml`.
type EnvironmentFile struct {
	// XMLName identifies the required `<env>` document root.
	XMLName struct{} `json:"env" yaml:"env" xml:"env"`
	// Territories declares territory behaviors, zone sources, and the agent classes that can populate those territories.
	Territories *EnvironmentTerritories `json:"territories,omitempty" yaml:"territories,omitempty" xml:"territories,omitempty"`
}

// EventSpawnsFile is the root of `cfgeventspawns.xml`.
type EventSpawnsFile struct {
	// XMLName is XML root marker for `<eventposdef>`.
	XMLName xml.Name `json:"eventposdef" yaml:"eventposdef" xml:"eventposdef"`
	// Events is list of event spawn definitions.
	Events []EventSpawnEntry `json:"event,omitempty" yaml:"event,omitempty" xml:"event,omitempty"`
}

// EventGroupsFile is the root of `cfgeventgroups.xml`.
type EventGroupsFile struct {
	// XMLName is XML root marker for `<eventgroupdef>`.
	XMLName xml.Name `json:"eventgroupdef" yaml:"eventgroupdef" xml:"eventgroupdef"`
	// Groups is list of reusable static object groups.
	Groups []EventGroup `json:"group,omitempty" yaml:"group,omitempty" xml:"group,omitempty"`
}

// PlayerSpawnPointsFile configures candidate locations and selection rules for fresh-character, server-hop, and travel spawns.
type PlayerSpawnPointsFile struct {
	// XMLName identifies the required `<playerspawnpoints>` document root.
	XMLName struct{} `json:"playerspawnpoints" yaml:"playerspawnpoints" xml:"playerspawnpoints"`
	// Fresh configures locations and selection rules for newly created characters.
	Fresh *PlayerSpawnProfile `json:"fresh,omitempty" yaml:"fresh,omitempty" xml:"fresh,omitempty"`
	// Hop configures locations and selection rules for characters returning after connecting to another server.
	Hop *PlayerSpawnProfile `json:"hop,omitempty" yaml:"hop,omitempty" xml:"hop,omitempty"`
	// Travel configures locations and selection rules for characters moved between servers in a travel connection.
	Travel *PlayerSpawnProfile `json:"travel,omitempty" yaml:"travel,omitempty" xml:"travel,omitempty"`
}

// WeatherFile configures global weather simulation and the target ranges
// and transition behavior for each weather channel in `cfgweather.xml`.
type WeatherFile struct {
	// Overcast configures cloudiness channel behavior.
	Overcast *WeatherSection `json:"overcast,omitempty" yaml:"overcast,omitempty" xml:"overcast,omitempty"`
	// Fog configures fog channel behavior.
	Fog *WeatherSection `json:"fog,omitempty" yaml:"fog,omitempty" xml:"fog,omitempty"`
	// Rain configures rain channel behavior.
	Rain *WeatherSection `json:"rain,omitempty" yaml:"rain,omitempty" xml:"rain,omitempty"`
	// WindMagnitude configures wind speed behavior.
	WindMagnitude *WeatherSection `json:"windMagnitude,omitempty" yaml:"windMagnitude,omitempty" xml:"windMagnitude,omitempty"`
	// WindDirection configures wind direction behavior.
	WindDirection *WeatherSection `json:"windDirection,omitempty" yaml:"windDirection,omitempty" xml:"windDirection,omitempty"`
	// Snowfall configures snowfall behavior.
	Snowfall *WeatherSection `json:"snowfall,omitempty" yaml:"snowfall,omitempty" xml:"snowfall,omitempty"`
	// Storm configures when storms can occur and how long their thunder effects remain active.
	Storm *WeatherStorm `json:"storm,omitempty" yaml:"storm,omitempty" xml:"storm,omitempty"`
	// Reset determines whether mission startup replaces persisted weather with the values configured in this file.
	Reset WeatherToggle `json:"reset" yaml:"reset" xml:"reset,attr"`
	// Enable determines whether the server runs the configured weather simulation.
	Enable WeatherToggle `json:"enable" yaml:"enable" xml:"enable,attr"`
	// XMLName identifies the required `<weather>` document root.
	XMLName struct{} `json:"weather" yaml:"weather" xml:"weather"`
}

// LimitsDefinitionFile declares the category, tag, usage, and value names referenced by Central Economy item filters.
type LimitsDefinitionFile struct {
	// XMLName identifies the required `<lists>` document root.
	XMLName struct{} `json:"lists" yaml:"lists" xml:"lists"`
	// Categories defines the item categories recognized by the economy.
	Categories *LimitsCategoryList `json:"categories,omitempty" yaml:"categories,omitempty" xml:"categories,omitempty"`
	// Tags defines the item tags recognized by the economy.
	Tags *LimitsTagList `json:"tags,omitempty" yaml:"tags,omitempty" xml:"tags,omitempty"`
	// UsageFlags defines location-purpose flags used to filter spawn locations.
	UsageFlags *LimitsUsageList `json:"usageflags,omitempty" yaml:"usageflags,omitempty" xml:"usageflags,omitempty"`
	// ValueFlags defines loot-tier flags used to filter spawn locations.
	ValueFlags *LimitsValueList `json:"valueflags,omitempty" yaml:"valueflags,omitempty" xml:"valueflags,omitempty"`
}

// LimitsDefinitionUserFile defines named shorthand groups that combine existing usage or value flags for reuse in item filters.
type LimitsDefinitionUserFile struct {
	// XMLName identifies the required `<user_lists>` document root.
	XMLName struct{} `json:"user_lists" yaml:"user_lists" xml:"user_lists"`
	// UsageFlags maps shorthand names to groups of usage flags.
	UsageFlags *LimitsUserBindingList `json:"usageflags,omitempty" yaml:"usageflags,omitempty" xml:"usageflags,omitempty"`
	// ValueFlags maps shorthand names to groups of value flags.
	ValueFlags *LimitsUserBindingList `json:"valueflags,omitempty" yaml:"valueflags,omitempty" xml:"valueflags,omitempty"`
}

// IgnoreListFile is the root of `cfgignorelist.xml`.
type IgnoreListFile struct {
	// XMLName is XML root marker for `<ignore>`.
	XMLName xml.Name `json:"ignore" yaml:"ignore" xml:"ignore"`
	// Types is list of item class names ignored by CE.
	Types []NamedRef `json:"type,omitempty" yaml:"type,omitempty" xml:"type,omitempty"`
}

// TerritoryFile defines named spawn zones used by territory profiles in `cfgenvironment.xml`.
type TerritoryFile struct {
	// XMLName identifies the required `<territory-type>` document root.
	XMLName xml.Name `json:"territory-type" yaml:"territory-type" xml:"territory-type"`
	// Territories defines zone positions and static and dynamic entity count
	// limits for territory behaviors.
	Territories []TerritoryBlock `json:"territory,omitempty" yaml:"territory,omitempty" xml:"territory,omitempty"`
}

// MapGroupProtoFile defines reusable building and cluster prototypes for loot placement,
// including containers, spawn points, and prototype defaults.
type MapGroupProtoFile struct {
	// Defaults provides shared property values inherited by prototypes that do not specify those values directly.
	Defaults *MapPrototypeDefaults `json:"defaults,omitempty" yaml:"defaults,omitempty" xml:"defaults,omitempty"`
	// Clusters maps terrain model shapes to logical cluster prototype names.
	Clusters *MapPrototypeExports `json:"clusters,omitempty" yaml:"clusters,omitempty" xml:"clusters,omitempty"`
	// XMLName identifies the required `<prototype>` document root.
	XMLName xml.Name `json:"prototype" yaml:"prototype" xml:"prototype"`
	// Groups defines building prototypes and their lootable containers.
	Groups []MapPrototypeEntry `json:"group,omitempty" yaml:"group,omitempty" xml:"group,omitempty"`
	// ClusterGroups defines reusable groups of nearby objects and their loot containers.
	ClusterGroups []MapPrototypeEntry `json:"cluster,omitempty" yaml:"cluster,omitempty" xml:"cluster,omitempty"`
}

// MapClusterProtoFile defines reusable cluster prototypes for grouped map objects and their loot placement.
type MapClusterProtoFile struct {
	// Defaults provides shared property values inherited by prototypes that do not specify those values directly.
	Defaults *MapPrototypeDefaults `json:"defaults,omitempty" yaml:"defaults,omitempty" xml:"defaults,omitempty"`
	// Clusters maps terrain model shapes to logical cluster prototype names.
	Clusters *MapPrototypeExports `json:"clusters,omitempty" yaml:"clusters,omitempty" xml:"clusters,omitempty"`
	// XMLName identifies the required `<prototype>` document root.
	XMLName xml.Name `json:"prototype" yaml:"prototype" xml:"prototype"`
	// Groups defines building prototypes and their lootable containers.
	Groups []MapPrototypeEntry `json:"group,omitempty" yaml:"group,omitempty" xml:"group,omitempty"`
	// ClusterGroups defines reusable groups of nearby objects and their loot containers.
	ClusterGroups []MapPrototypeEntry `json:"cluster,omitempty" yaml:"cluster,omitempty" xml:"cluster,omitempty"`
}

// MapGroupPosFile places loot prototypes at world positions for buildings and other terrain objects in `mapgrouppos.xml`.
type MapGroupPosFile struct {
	// XMLName identifies the required `<map>` document root.
	XMLName xml.Name `json:"map" yaml:"map" xml:"map"`
	// Groups lists prototype names and their world positions and rotations.
	Groups []MapGroupInstance `json:"group,omitempty" yaml:"group,omitempty" xml:"group,omitempty"`
}

// MapGroupDirtFile places dirt and other terrain-specific map groups at world positions in `mapgroupdirt.xml`.
type MapGroupDirtFile struct {
	// XMLName identifies the required `<map>` document root.
	XMLName xml.Name `json:"map" yaml:"map" xml:"map"`
	// Groups lists prototype names and their world positions and rotations.
	Groups []MapGroupInstance `json:"group,omitempty" yaml:"group,omitempty" xml:"group,omitempty"`
}

// MapGroupClusterFile places cluster prototypes at world positions and can include additional cluster placement files.
type MapGroupClusterFile struct {
	// XMLName identifies the required `<map>` document root.
	XMLName xml.Name `json:"map" yaml:"map" xml:"map"`
	// Groups lists cluster prototype names and their world positions and rotations.
	Groups []MapGroupInstance `json:"group,omitempty" yaml:"group,omitempty" xml:"group,omitempty"`
	// Includes references additional map group cluster files loaded with this file.
	Includes []MapGroupInclude `json:"include,omitempty" yaml:"include,omitempty" xml:"include,omitempty"`
}

// MapGroupInclude identifies another map group cluster file to load.
type MapGroupInclude struct {
	// File is the referenced file name without its `.xml` extension.
	File string `json:"file" yaml:"file" xml:"file,attr"`
}
