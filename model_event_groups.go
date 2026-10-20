// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// EventGroup defines the objects placed together when a position in cfgeventspawns.xml references this group.
type EventGroup struct {
	// Name is the identifier matched by the event spawn's `group` attribute.
	Name string `json:"name" yaml:"name" xml:"name,attr"`
	// Children are the objects in the composition. Their positions are offsets from the event spawn position.
	Children []EventGroupChild `json:"child,omitempty" yaml:"child,omitempty" xml:"child,omitempty"`
}

// EventGroupChild defines one object and its placement or loot behavior within an event group.
type EventGroupChild struct {
	// Type is the game class spawned as this object in the composition.
	Type string `json:"type" yaml:"type" xml:"type,attr"`
	// Deloot sets the amount of dynamic-event loot requested for this object, overriding the default amount.
	Deloot string `json:"deloot,omitempty" yaml:"deloot,omitempty" xml:"deloot,attr,omitempty"`
	// LootMax is the maximum number of random items requested for this event child.
	LootMax string `json:"lootmax,omitempty" yaml:"lootmax,omitempty" xml:"lootmax,attr,omitempty"`
	// LootMin is the minimum number of random items requested for this event child.
	LootMin string `json:"lootmin,omitempty" yaml:"lootmin,omitempty" xml:"lootmin,attr,omitempty"`
	// SpawnSecondary controls whether spawning this object starts its secondary infected event.
	SpawnSecondary string `json:"spawnsecondary,omitempty" yaml:"spawnsecondary,omitempty" xml:"spawnsecondary,attr,omitempty"`
	// X is the horizontal X offset from the event spawn position.
	X string `json:"x,omitempty" yaml:"x,omitempty" xml:"x,attr,omitempty"`
	// Z is the horizontal Z offset from the event spawn position.
	Z string `json:"z,omitempty" yaml:"z,omitempty" xml:"z,attr,omitempty"`
	// A is the object's yaw rotation in degrees.
	A string `json:"a,omitempty" yaml:"a,omitempty" xml:"a,attr,omitempty"`
	// Y is the vertical offset from the event spawn position.
	Y string `json:"y,omitempty" yaml:"y,omitempty" xml:"y,attr,omitempty"`
}
