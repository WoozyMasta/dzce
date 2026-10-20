// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// TerritoryBlock groups zones for one territory profile and supplies their shared display color for Central Economy tooling.
type TerritoryBlock struct {
	// Color is the ARGB color used to display this territory in tooling.
	Color string `json:"color,omitempty" yaml:"color,omitempty" xml:"color,attr,omitempty" jsonschema:"example=4294923520"`
	// Zones are circular spawn areas used by territory behaviors.
	Zones []TerritoryZone `json:"zone,omitempty" yaml:"zone,omitempty" xml:"zone,omitempty"`
}

// TerritoryZone defines one circular spawn area and its static and dynamic entity count limits.
type TerritoryZone struct {
	// Name identifies the zone for references and diagnostics.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty" jsonschema:"example=Graze"`
	// SMin is the lower bound for the static entity count when the territory behavior supports static spawning.
	SMin string `json:"smin,omitempty" yaml:"smin,omitempty" xml:"smin,attr,omitempty" jsonschema:"example=0"`
	// SMax is the upper bound for the static entity count when the territory behavior supports static spawning.
	SMax string `json:"smax,omitempty" yaml:"smax,omitempty" xml:"smax,attr,omitempty" jsonschema:"example=0"`
	// DMin is the lower bound for the dynamic entity count in this zone.
	// For ambient spawners, the count is chosen from DMin through DMax, inclusive, when the zone activates.
	DMin string `json:"dmin,omitempty" yaml:"dmin,omitempty" xml:"dmin,attr,omitempty" jsonschema:"example=0"`
	// DMax is the upper bound for the dynamic entity count in this zone.
	// For an ambient event, its global max can still cap the count across all zones.
	DMax string `json:"dmax,omitempty" yaml:"dmax,omitempty" xml:"dmax,attr,omitempty" jsonschema:"example=0"`
	// X is the map's horizontal X coordinate for the center of the zone.
	X string `json:"x,omitempty" yaml:"x,omitempty" xml:"x,attr,omitempty" jsonschema:"example=1920"`
	// Z is the map's horizontal Z coordinate for the center of the zone.
	Z string `json:"z,omitempty" yaml:"z,omitempty" xml:"z,attr,omitempty" jsonschema:"example=13115"`
	// R is the radius of the area where the zone's entities may spawn, in meters.
	R string `json:"r,omitempty" yaml:"r,omitempty" xml:"r,attr,omitempty" jsonschema:"example=100"`
}
