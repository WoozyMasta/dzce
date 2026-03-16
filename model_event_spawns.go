// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// EventSpawnEntry is one event definition in event spawns config.
type EventSpawnEntry struct {
	// Name is event name this position set belongs to.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
	// Zone is optional randomized zone around anchor.
	Zone *EventSpawnZone `xml:"zone,omitempty" json:"zone,omitempty" yaml:"zone,omitempty"`
	// Pos is explicit fixed spawn positions.
	Pos []EventSpawnPos `xml:"pos,omitempty" json:"pos,omitempty" yaml:"pos,omitempty"`
}

// EventSpawnZone defines spawn zone limits for event positions.
type EventSpawnZone struct {
	// SMin is minimum slope value for zone placement.
	SMin string `xml:"smin,attr,omitempty" json:"smin,omitempty" yaml:"smin,omitempty"`
	// SMax is maximum slope value for zone placement.
	SMax string `xml:"smax,attr,omitempty" json:"smax,omitempty" yaml:"smax,omitempty"`
	// DMin is minimum distance from center.
	DMin string `xml:"dmin,attr,omitempty" json:"dmin,omitempty" yaml:"dmin,omitempty"`
	// DMax is maximum distance from center.
	DMax string `xml:"dmax,attr,omitempty" json:"dmax,omitempty" yaml:"dmax,omitempty"`
	// R is optional zone radius.
	R string `xml:"r,attr,omitempty" json:"r,omitempty" yaml:"r,omitempty"`
}

// EventSpawnPos stores one event position in cfgeventspawns.xml.
type EventSpawnPos struct {
	// X is world X coordinate.
	X string `xml:"x,attr,omitempty" json:"x,omitempty" yaml:"x,omitempty"`
	// Z is world Z coordinate.
	Z string `xml:"z,attr,omitempty" json:"z,omitempty" yaml:"z,omitempty"`
	// A is yaw angle in degrees.
	A string `xml:"a,attr,omitempty" json:"a,omitempty" yaml:"a,omitempty"`
	// Y is optional vertical offset.
	Y string `xml:"y,attr,omitempty" json:"y,omitempty" yaml:"y,omitempty"`
	// Group binds this pos to event group name.
	Group string `xml:"group,attr,omitempty" json:"group,omitempty" yaml:"group,omitempty"`
}
