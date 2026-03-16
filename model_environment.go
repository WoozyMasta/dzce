// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// EnvironmentTerritories groups territory files and bindings.
type EnvironmentTerritories struct {
	// Files is reusable territory file includes.
	Files []EnvironmentTerritoryFile `xml:"file,omitempty" json:"file,omitempty" yaml:"file,omitempty"`
	// Territories is inline territory behavior definitions.
	Territories []EnvironmentTerritory `xml:"territory,omitempty" json:"territory,omitempty" yaml:"territory,omitempty"`
}

// EnvironmentTerritory binds territory behavior and agents.
type EnvironmentTerritory struct {
	// Type defines territory system category.
	Type string `xml:"type,attr,omitempty" json:"type,omitempty" yaml:"type,omitempty"`
	// Name is territory profile identifier.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	// Behavior is script behavior class name.
	Behavior string `xml:"behavior,attr,omitempty" json:"behavior,omitempty" yaml:"behavior,omitempty"`
	// Files references territory zone data.
	Files []EnvironmentTerritoryFile `xml:"file,omitempty" json:"file,omitempty" yaml:"file,omitempty"`
	// Agents configures ambient agents for this territory.
	Agents []EnvironmentAgent `xml:"agent,omitempty" json:"agent,omitempty" yaml:"agent,omitempty"`
	// Items stores auxiliary numeric/string params.
	Items []EnvironmentItem `xml:"item,omitempty" json:"item,omitempty" yaml:"item,omitempty"`
}

// EnvironmentAgent describes one agent block in environment territory.
type EnvironmentAgent struct {
	// Type is agent variant (male/female/etc).
	Type string `xml:"type,attr,omitempty" json:"type,omitempty" yaml:"type,omitempty"`
	// Chance is spawn chance for this agent block.
	Chance string `xml:"chance,attr,omitempty" json:"chance,omitempty" yaml:"chance,omitempty"`
	// Spawns is concrete config classes spawned by this agent.
	Spawns []EnvironmentSpawn `xml:"spawn,omitempty" json:"spawn,omitempty" yaml:"spawn,omitempty"`
	// Items is extra settings attached to this agent.
	Items []EnvironmentItem `xml:"item,omitempty" json:"item,omitempty" yaml:"item,omitempty"`
}

// EnvironmentSpawn is one spawn entry in environment agent.
type EnvironmentSpawn struct {
	// ConfigName is class name to spawn.
	ConfigName string `xml:"configName,attr,omitempty" json:"configName,omitempty" yaml:"configName,omitempty"`
	// Chance is per-entry spawn chance.
	Chance string `xml:"chance,attr,omitempty" json:"chance,omitempty" yaml:"chance,omitempty"`
}

// EnvironmentItem stores one named item/value pair in environment blocks.
type EnvironmentItem struct {
	// Name is setting key.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	// Value is setting value as text.
	Value string `xml:"val,attr,omitempty" json:"val,omitempty" yaml:"val,omitempty"`
}
