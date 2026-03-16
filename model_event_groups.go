// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// EventGroup stores grouped static objects for dynamic events.
type EventGroup struct {
	// Name is group identifier referenced by spawns.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
	// Children is objects spawned with relative transforms.
	Children []EventGroupChild `xml:"child,omitempty" json:"child,omitempty" yaml:"child,omitempty"`
}

// EventGroupChild stores one child object in event group.
type EventGroupChild struct {
	// Type is spawned object class.
	Type string `xml:"type,attr" json:"type" yaml:"type"`
	// Deloot controls cleanup rules for dropped loot.
	Deloot string `xml:"deloot,attr,omitempty" json:"deloot,omitempty" yaml:"deloot,omitempty"`
	// LootMax is max loot amount for this object.
	LootMax string `xml:"lootmax,attr,omitempty" json:"lootmax,omitempty" yaml:"lootmax,omitempty"`
	// LootMin is min loot amount for this object.
	LootMin string `xml:"lootmin,attr,omitempty" json:"lootmin,omitempty" yaml:"lootmin,omitempty"`
	// SpawnSecondary allows secondary children spawn.
	SpawnSecondary string `xml:"spawnsecondary,attr,omitempty" json:"spawnsecondary,omitempty" yaml:"spawnsecondary,omitempty"`
	// X is local X offset from group anchor.
	X string `xml:"x,attr,omitempty" json:"x,omitempty" yaml:"x,omitempty"`
	// Z is local Z offset from group anchor.
	Z string `xml:"z,attr,omitempty" json:"z,omitempty" yaml:"z,omitempty"`
	// A is local yaw angle.
	A string `xml:"a,attr,omitempty" json:"a,omitempty" yaml:"a,omitempty"`
	// Y is local vertical offset.
	Y string `xml:"y,attr,omitempty" json:"y,omitempty" yaml:"y,omitempty"`
}
