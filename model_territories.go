// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// TerritoryBlock stores one territory with shared color value.
type TerritoryBlock struct {
	// Color is ARGB color used by tooling/visualization.
	Color string `xml:"color,attr,omitempty" json:"color,omitempty" yaml:"color,omitempty"`
	// Zones is list of zone circles in this block.
	Zones []TerritoryZone `xml:"zone,omitempty" json:"zone,omitempty" yaml:"zone,omitempty"`
}

// TerritoryZone stores one zone entry in territory files.
type TerritoryZone struct {
	// Name is zone identifier.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	// SMin is minimum slope for zone logic.
	SMin string `xml:"smin,attr,omitempty" json:"smin,omitempty" yaml:"smin,omitempty"`
	// SMax is maximum slope for zone logic.
	SMax string `xml:"smax,attr,omitempty" json:"smax,omitempty" yaml:"smax,omitempty"`
	// DMin is minimum distance threshold.
	DMin string `xml:"dmin,attr,omitempty" json:"dmin,omitempty" yaml:"dmin,omitempty"`
	// DMax is maximum distance threshold.
	DMax string `xml:"dmax,attr,omitempty" json:"dmax,omitempty" yaml:"dmax,omitempty"`
	// X is world X coordinate of zone center.
	X string `xml:"x,attr,omitempty" json:"x,omitempty" yaml:"x,omitempty"`
	// Z is world Z coordinate of zone center.
	Z string `xml:"z,attr,omitempty" json:"z,omitempty" yaml:"z,omitempty"`
	// R is zone radius.
	R string `xml:"r,attr,omitempty" json:"r,omitempty" yaml:"r,omitempty"`
}
