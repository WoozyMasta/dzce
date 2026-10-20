// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// TerritoryBlock groups zones for one territory profile and supplies their shared display color for Central Economy tooling.
type TerritoryBlock struct {
	// Color is the ARGB color used to display this territory in tooling.
	Color string `json:"color,omitempty" yaml:"color,omitempty" xml:"color,attr,omitempty"`
	// Zones are circular spawn areas used by territory behaviors.
	Zones []TerritoryZone `json:"zone,omitempty" yaml:"zone,omitempty" xml:"zone,omitempty"`
}

// TerritoryZone defines one circular spawn area and its static and dynamic entity count limits.
type TerritoryZone struct {
	// Name identifies the zone for references and diagnostics.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty"`
	// SMin is the lower bound for the static entity count when the territory behavior supports static spawning.
	SMin string `json:"smin,omitempty" yaml:"smin,omitempty" xml:"smin,attr,omitempty"`
	// SMax is the upper bound for the static entity count when the territory behavior supports static spawning.
	SMax string `json:"smax,omitempty" yaml:"smax,omitempty" xml:"smax,attr,omitempty"`
	// DMin is the lower bound for the dynamic entity count in this zone.
	// For ambient spawners, the count is chosen from DMin through DMax, inclusive, when the zone activates.
	DMin string `json:"dmin,omitempty" yaml:"dmin,omitempty" xml:"dmin,attr,omitempty"`
	// DMax is the upper bound for the dynamic entity count in this zone.
	// For an ambient event, its global max can still cap the count across all zones.
	DMax string `json:"dmax,omitempty" yaml:"dmax,omitempty" xml:"dmax,attr,omitempty"`
	// X is the map's horizontal X coordinate for the center of the zone.
	X string `json:"x,omitempty" yaml:"x,omitempty" xml:"x,attr,omitempty"`
	// Z is the map's horizontal Z coordinate for the center of the zone.
	Z string `json:"z,omitempty" yaml:"z,omitempty" xml:"z,attr,omitempty"`
	// R is the radius of the area where the zone's entities may spawn, in meters.
	R string `json:"r,omitempty" yaml:"r,omitempty" xml:"r,attr,omitempty"`
}
