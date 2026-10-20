// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// EnvironmentTerritories configures territory spawning systems
// by including territory data files and binding those files to spawn behaviors.
type EnvironmentTerritories struct {
	// Files declares territory files that can be referenced by a territory's usable identifier.
	Files []EnvironmentTerritoryFile `json:"file,omitempty" yaml:"file,omitempty" xml:"file,omitempty"`
	// Territories bind a territory profile to its script behavior, zone data, agent classes, and behavior-specific settings.
	Territories []EnvironmentTerritory `json:"territory,omitempty" yaml:"territory,omitempty" xml:"territory,omitempty"`
}

// EnvironmentTerritory connects a named territory profile to the behavior that consumes its zones and agent definitions.
type EnvironmentTerritory struct {
	// Type selects the territory system category, such as `Herd` or `Ambient`.
	Type string `json:"type,omitempty" yaml:"type,omitempty" xml:"type,attr,omitempty"`
	// Name identifies the profile used by the corresponding Central Economy spawner or event configuration.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty"`
	// Behavior is the Enforce Script behavior class that implements this territory's runtime logic.
	Behavior string `json:"behavior,omitempty" yaml:"behavior,omitempty" xml:"behavior,attr,omitempty"`
	// Files select the territory-zone data used by this profile.
	// Entries can refer to paths directly or to usable identifiers declared above.
	Files []EnvironmentTerritoryFile `json:"file,omitempty" yaml:"file,omitempty" xml:"file,omitempty"`
	// Agents define the agent variants and game classes that may spawn in the territory.
	Agents []EnvironmentAgent `json:"agent,omitempty" yaml:"agent,omitempty" xml:"agent,omitempty"`
	// Items provide named, behavior-specific settings as text values.
	Items []EnvironmentItem `json:"item,omitempty" yaml:"item,omitempty" xml:"item,omitempty"`
}

// EnvironmentAgent defines one agent variant and the classes available for that variant in a territory behavior.
type EnvironmentAgent struct {
	// Type is the agent variant used by this block, such as `Male` or `Female`.
	Type string `json:"type,omitempty" yaml:"type,omitempty" xml:"type,attr,omitempty"`
	// Chance is the selection weight for this agent variant.
	Chance string `json:"chance,omitempty" yaml:"chance,omitempty" xml:"chance,attr,omitempty"`
	// Spawns are the game classes this agent variant can create.
	Spawns []EnvironmentSpawn `json:"spawn,omitempty" yaml:"spawn,omitempty" xml:"spawn,omitempty"`
	// Items provide agent-specific settings consumed by the selected behavior.
	Items []EnvironmentItem `json:"item,omitempty" yaml:"item,omitempty" xml:"item,omitempty"`
}

// EnvironmentSpawn assigns a game class and selection weight to a territory agent variant.
type EnvironmentSpawn struct {
	// ConfigName is the game class name that may be spawned.
	ConfigName string `json:"configName,omitempty" yaml:"configName,omitempty" xml:"configName,attr,omitempty"`
	// Chance is the selection weight for this class within its agent variant.
	Chance string `json:"chance,omitempty" yaml:"chance,omitempty" xml:"chance,attr,omitempty"`
}

// EnvironmentItem stores one behavior-specific setting as a named text value.
type EnvironmentItem struct {
	// Name is the setting key interpreted by the selected territory behavior.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty"`
	// Value is the setting value encoded as text.
	Value string `json:"val,omitempty" yaml:"val,omitempty" xml:"val,attr,omitempty"`
}
