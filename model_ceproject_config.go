// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "encoding/xml"

// CEProjectConfigFile is the root of CEProject `mapname.xml` (`<zg-config>`).
type CEProjectConfigFile struct {
	// Selected stores current GUI row selection.
	Selected *CEProjectSelected `xml:"selected,omitempty" json:"selected,omitempty" yaml:"selected,omitempty"`
	// Global stores map/world and background settings.
	Global *CEProjectGlobal `xml:"global,omitempty" json:"global,omitempty" yaml:"global,omitempty"`
	// Areas stores usage/value flag names.
	Areas *CEProjectAreas `xml:"areas,omitempty" json:"areas,omitempty" yaml:"areas,omitempty"`
	// Brush stores active brush settings.
	Brush *CEProjectBrush `xml:"brush,omitempty" json:"brush,omitempty" yaml:"brush,omitempty"`
	// Layers stores paint layer definitions.
	Layers *CEProjectLayers `xml:"layers,omitempty" json:"layers,omitempty" yaml:"layers,omitempty"`
	// TerritoryTypeList stores all territory type definitions.
	TerritoryTypeList *CEProjectTerritoryTypeList `xml:"territory-type-list,omitempty" json:"territory-type-list,omitempty" yaml:"territory-type-list,omitempty"`
	// XMLName is XML root marker for `<zg-config>`.
	XMLName xml.Name `xml:"zg-config" json:"zg-config" yaml:"zg-config"`
}

// CEProjectSelected stores selected row index in GUI state.
type CEProjectSelected struct {
	// Row is selected row index.
	Row string `xml:"row,attr,omitempty" json:"row,omitempty" yaml:"row,omitempty"`
}

// CEProjectGlobal stores global map/world settings.
type CEProjectGlobal struct {
	// Background stores preview map texture info.
	Background *CEProjectBackground `xml:"background,omitempty" json:"background,omitempty" yaml:"background,omitempty"`
	// Layer stores raster layer size.
	Layer *CEProjectDimension `xml:"layer,omitempty" json:"layer,omitempty" yaml:"layer,omitempty"`
	// World stores world size.
	World *CEProjectDimension `xml:"world,omitempty" json:"world,omitempty" yaml:"world,omitempty"`
}

// CEProjectBackground stores background texture settings.
type CEProjectBackground struct {
	// File is map image path used by CEProject.
	File string `xml:"file,attr,omitempty" json:"file,omitempty" yaml:"file,omitempty"`
	// RGBA is packed background tint color.
	RGBA string `xml:"rgba,attr,omitempty" json:"rgba,omitempty" yaml:"rgba,omitempty"`
}

// CEProjectDimension stores one numeric `size` attribute node.
type CEProjectDimension struct {
	// Size is numeric size value as text.
	Size string `xml:"size,attr,omitempty" json:"size,omitempty" yaml:"size,omitempty"`
}

// CEProjectAreas stores usage and value flag names.
type CEProjectAreas struct {
	// Usages stores Area Usage Flags labels from CEProject UI.
	// These labels map to usage bit positions used by `layers/layer@usage_flags`.
	Usages *CEProjectAreaUsageList `xml:"usages,omitempty" json:"usages,omitempty" yaml:"usages,omitempty"`
	// Values stores Area Value Flags labels from CEProject UI (tiers).
	// These labels map to value bit positions used by `layers/layer@value_flags`.
	Values *CEProjectAreaValueList `xml:"values,omitempty" json:"values,omitempty" yaml:"values,omitempty"`
}

// CEProjectAreaUsageList stores usage name entries.
type CEProjectAreaUsageList struct {
	// Items stores `<usage name="..."/>` values.
	Items []NamedRef `xml:"usage,omitempty" json:"usage,omitempty" yaml:"usage,omitempty"`
}

// CEProjectAreaValueList stores value name entries.
type CEProjectAreaValueList struct {
	// Items stores `<value name="..."/>` values.
	Items []NamedRef `xml:"value,omitempty" json:"value,omitempty" yaml:"value,omitempty"`
}

// CEProjectBrush stores GUI brush settings.
type CEProjectBrush struct {
	// Color is packed brush fill color.
	Color string `xml:"color,attr,omitempty" json:"color,omitempty" yaml:"color,omitempty"`
	// Outline toggles outline drawing.
	Outline string `xml:"outline,attr,omitempty" json:"outline,omitempty" yaml:"outline,omitempty"`
	// OutlineColor is packed outline color.
	OutlineColor string `xml:"outline_color,attr,omitempty" json:"outline_color,omitempty" yaml:"outline_color,omitempty"`
	// OutlineWidth is outline width.
	OutlineWidth string `xml:"outline_width,attr,omitempty" json:"outline_width,omitempty" yaml:"outline_width,omitempty"`
	// Size is brush radius/size.
	Size string `xml:"size,attr,omitempty" json:"size,omitempty" yaml:"size,omitempty"`
}

// CEProjectLayers stores layer definitions.
type CEProjectLayers struct {
	// Layers stores paint layer entries.
	Layers []CEProjectLayer `xml:"layer,omitempty" json:"layer,omitempty" yaml:"layer,omitempty"`
}

// CEProjectLayer stores one paint layer entry from `<layers>`.
type CEProjectLayer struct {
	// UsageFlags is bitmask of usage flags this layer writes.
	UsageFlags string `xml:"usage_flags,attr,omitempty" json:"usage_flags,omitempty" yaml:"usage_flags,omitempty"`
	// ValueFlags is bitmask of value flags this layer writes.
	ValueFlags string `xml:"value_flags,attr,omitempty" json:"value_flags,omitempty" yaml:"value_flags,omitempty"`
	// Color is packed layer display color.
	Color string `xml:"color,attr,omitempty" json:"color,omitempty" yaml:"color,omitempty"`
	// Visible toggles layer visibility in CEProject UI only.
	// It does not affect areaflags binary payload.
	Visible string `xml:"visible,attr,omitempty" json:"visible,omitempty" yaml:"visible,omitempty"`
	// Name is layer name, also used as mask file base name.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
}

// CEProjectTerritoryTypeList stores all territory types from config.
type CEProjectTerritoryTypeList struct {
	// Types stores individual `<territory-type>` entries.
	// CEProject save/export may rewrite this section from `territoryTypes/*.xml`.
	Types []CEProjectTerritoryType `xml:"territory-type,omitempty" json:"territory-type,omitempty" yaml:"territory-type,omitempty"`
}

// CEProjectTerritoryType stores one territory type and its territories.
type CEProjectTerritoryType struct {
	// Name is territory type name, usually used for output file name.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	// Territories stores territory blocks.
	Territories []CEProjectTerritory `xml:"territory,omitempty" json:"territory,omitempty" yaml:"territory,omitempty"`
}

// CEProjectTerritory stores one territory block from CEProject config.
type CEProjectTerritory struct {
	// Name is territory name in editor.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	// Visible toggles visibility in editor.
	Visible string `xml:"visible,attr,omitempty" json:"visible,omitempty" yaml:"visible,omitempty"`
	// Color is packed ARGB color value.
	Color string `xml:"color,attr,omitempty" json:"color,omitempty" yaml:"color,omitempty"`
	// Zones stores territory zones.
	Zones []CEProjectTerritoryZone `xml:"zone,omitempty" json:"zone,omitempty" yaml:"zone,omitempty"`
}

// CEProjectTerritoryZone stores one zone entry in CEProject config.
type CEProjectTerritoryZone struct {
	// Name is zone name.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	// SMin is minimum static spawn count.
	SMin string `xml:"smin,attr,omitempty" json:"smin,omitempty" yaml:"smin,omitempty"`
	// SMax is maximum static spawn count.
	SMax string `xml:"smax,attr,omitempty" json:"smax,omitempty" yaml:"smax,omitempty"`
	// DMin is minimum dynamic spawn count.
	DMin string `xml:"dmin,attr,omitempty" json:"dmin,omitempty" yaml:"dmin,omitempty"`
	// DMax is maximum dynamic spawn count.
	DMax string `xml:"dmax,attr,omitempty" json:"dmax,omitempty" yaml:"dmax,omitempty"`
	// X is zone center world X coordinate.
	X string `xml:"x,attr,omitempty" json:"x,omitempty" yaml:"x,omitempty"`
	// Z is zone center world Z coordinate.
	Z string `xml:"z,attr,omitempty" json:"z,omitempty" yaml:"z,omitempty"`
	// R is zone radius.
	R string `xml:"r,attr,omitempty" json:"r,omitempty" yaml:"r,omitempty"`
	// D is optional additional tool-only zone spacing value.
	D string `xml:"d,attr,omitempty" json:"d,omitempty" yaml:"d,omitempty"`
}
