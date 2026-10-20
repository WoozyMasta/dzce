// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// EventSpawnEntry associates an event from db/events.xml with mission-defined positions
// and optional settings for infected spawned around those positions.
type EventSpawnEntry struct {
	// Zone sets infected-count ranges and their spawn radius for positions in this event.
	// It is used when the event has a secondary infected event.
	Zone *EventSpawnZone `json:"zone,omitempty" yaml:"zone,omitempty" xml:"zone,omitempty"`
	// Name matches the event identifier in db/events.xml.
	Name string `json:"name" yaml:"name" xml:"name,attr"`
	// Pos lists positions where this event may spawn.
	// A position may reference a reusable object group from cfgeventgroups.xml.
	Pos []EventSpawnPos `json:"pos,omitempty" yaml:"pos,omitempty" xml:"pos,omitempty"`
}

// EventSpawnZone controls the number and spread of infected spawned around an event position
// when the event defines a secondary infected event.
type EventSpawnZone struct {
	// SMin is the lower bound for the static infected count.
	SMin string `json:"smin,omitempty" yaml:"smin,omitempty" xml:"smin,attr,omitempty"`
	// SMax is the upper bound for the static infected count.
	SMax string `json:"smax,omitempty" yaml:"smax,omitempty" xml:"smax,attr,omitempty"`
	// DMin is the lower bound for the dynamic infected count.
	DMin string `json:"dmin,omitempty" yaml:"dmin,omitempty" xml:"dmin,attr,omitempty"`
	// DMax is the upper bound for the dynamic infected count.
	DMax string `json:"dmax,omitempty" yaml:"dmax,omitempty" xml:"dmax,attr,omitempty"`
	// R is the radius in meters around the event position where infected may spawn.
	R string `json:"r,omitempty" yaml:"r,omitempty" xml:"r,attr,omitempty"`
}

// EventSpawnPos stores one mission-defined event position in map coordinates.
type EventSpawnPos struct {
	// X is the map's horizontal X coordinate.
	X string `json:"x,omitempty" yaml:"x,omitempty" xml:"x,attr,omitempty"`
	// Z is the map's horizontal Z coordinate.
	Z string `json:"z,omitempty" yaml:"z,omitempty" xml:"z,attr,omitempty"`
	// A is the event's yaw rotation in degrees. Shipped mission data also uses -1,
	// so retain the supplied value instead of normalizing it.
	A string `json:"a,omitempty" yaml:"a,omitempty" xml:"a,attr,omitempty"`
	// Y is the optional world height. It may be omitted when placement uses terrain height.
	Y string `json:"y,omitempty" yaml:"y,omitempty" xml:"y,attr,omitempty"`
	// Group names a cfgeventgroups.xml group to create at this position.
	Group string `json:"group,omitempty" yaml:"group,omitempty" xml:"group,attr,omitempty"`
}
