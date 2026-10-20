// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "encoding/xml"

// CEProjectConfigFile is a CEProject workspace file (`mapname.xml`).
// It stores editor settings, area-flag layer metadata, and territory data.
type CEProjectConfigFile struct {
	// Selected stores editor selection state and does not describe map data.
	Selected *CEProjectSelected `json:"selected,omitempty" yaml:"selected,omitempty" xml:"selected,omitempty"`
	// Global defines the map preview and the raster and world dimensions used by the project.
	Global *CEProjectGlobal `json:"global,omitempty" yaml:"global,omitempty" xml:"global,omitempty"`
	// Areas declares the usage and value labels available in the editor.
	// Their order maps labels to bit positions in layer masks.
	Areas *CEProjectAreas `json:"areas,omitempty" yaml:"areas,omitempty" xml:"areas,omitempty"`
	// Brush stores editor painting defaults; it is not an area-flags layer.
	Brush *CEProjectBrush `json:"brush,omitempty" yaml:"brush,omitempty" xml:"brush,omitempty"`
	// Layers maps named source masks to usage/value masks and editor display properties.
	Layers *CEProjectLayers `json:"layers,omitempty" yaml:"layers,omitempty" xml:"layers,omitempty"`
	// TerritoryTypeList embeds territory definitions in the project file.
	// CEProject may rewrite this section from the separate territory XML files when saving or exporting the project.
	TerritoryTypeList *CEProjectTerritoryTypeList `json:"territory-type-list,omitempty" yaml:"territory-type-list,omitempty" xml:"territory-type-list,omitempty"`
	// XMLName identifies the CEProject document root.
	XMLName xml.Name `json:"zg-config" yaml:"zg-config" xml:"zg-config"`
}

// CEProjectSelected stores the editor's current row selection.
type CEProjectSelected struct {
	// Row is the selected row value persisted by the editor.
	Row string `json:"row,omitempty" yaml:"row,omitempty" xml:"row,attr,omitempty"`
}

// CEProjectGlobal stores map preview metadata and project dimensions.
type CEProjectGlobal struct {
	// Background selects and tints the map image shown in the editor.
	Background *CEProjectBackground `json:"background,omitempty" yaml:"background,omitempty" xml:"background,omitempty"`
	// Layer.Size is the width and height of the square area-flags raster.
	Layer *CEProjectDimension `json:"layer,omitempty" yaml:"layer,omitempty" xml:"layer,omitempty"`
	// World.Size is the world-map extent used to relate raster cells to map coordinates.
	World *CEProjectDimension `json:"world,omitempty" yaml:"world,omitempty" xml:"world,omitempty"`
}

// CEProjectBackground stores the editor map preview image and its tint.
type CEProjectBackground struct {
	// File is the map image path opened for the editor preview.
	File string `json:"file,omitempty" yaml:"file,omitempty" xml:"file,attr,omitempty"`
	// RGBA is the packed color tint applied to the preview image.
	RGBA string `json:"rgba,omitempty" yaml:"rgba,omitempty" xml:"rgba,attr,omitempty"`
}

// CEProjectDimension stores a dimension value from a CEProject `size` attribute. The XML format represents it as text.
type CEProjectDimension struct {
	// Size is the raster or world dimension, depending on its parent element.
	Size string `json:"size,omitempty" yaml:"size,omitempty" xml:"size,attr,omitempty"`
}

// CEProjectAreas stores the ordered flag labels displayed in the editor.
type CEProjectAreas struct {
	// Usages defines labels for area usage bits.
	// The list order corresponds to bit positions in `layers/layer@usage_flags`.
	Usages *CEProjectAreaUsageList `json:"usages,omitempty" yaml:"usages,omitempty" xml:"usages,omitempty"`
	// Values defines labels for area value bits, commonly used for loot tiers.
	// The list order corresponds to bit positions in `layers/layer@value_flags`.
	Values *CEProjectAreaValueList `json:"values,omitempty" yaml:"values,omitempty" xml:"values,omitempty"`
}

// CEProjectAreaUsageList stores usage flag labels in bit-position order.
type CEProjectAreaUsageList struct {
	// Items are the names shown by CEProject for usage flags.
	Items []NamedRef `json:"usage,omitempty" yaml:"usage,omitempty" xml:"usage,omitempty"`
}

// CEProjectAreaValueList stores value flag labels in bit-position order.
type CEProjectAreaValueList struct {
	// Items are the names shown by CEProject for value flags.
	Items []NamedRef `json:"value,omitempty" yaml:"value,omitempty" xml:"value,omitempty"`
}

// CEProjectBrush stores the editor's current painting tool settings.
type CEProjectBrush struct {
	// Color is the packed fill color used to preview the brush.
	Color string `json:"color,omitempty" yaml:"color,omitempty" xml:"color,attr,omitempty"`
	// Outline controls whether the brush preview draws an outline.
	Outline string `json:"outline,omitempty" yaml:"outline,omitempty" xml:"outline,attr,omitempty"`
	// OutlineColor is the packed color used for the brush preview outline.
	OutlineColor string `json:"outline_color,omitempty" yaml:"outline_color,omitempty" xml:"outline_color,attr,omitempty"`
	// OutlineWidth is the brush preview outline width.
	OutlineWidth string `json:"outline_width,omitempty" yaml:"outline_width,omitempty" xml:"outline_width,attr,omitempty"`
	// Size is the brush size used when painting source masks.
	Size string `json:"size,omitempty" yaml:"size,omitempty" xml:"size,attr,omitempty"`
}

// CEProjectLayers stores named source-mask layers used to paint area flags.
type CEProjectLayers struct {
	// Layers associate source mask files with area-flag bitmasks.
	Layers []CEProjectLayer `json:"layer,omitempty" yaml:"layer,omitempty" xml:"layer,omitempty"`
}

// CEProjectLayer describes one source mask and how its painted cells contribute to the compiled area-flags map.
type CEProjectLayer struct {
	// UsageFlags is the decimal bitmask of usage flags written by this layer.
	UsageFlags string `json:"usage_flags,omitempty" yaml:"usage_flags,omitempty" xml:"usage_flags,attr,omitempty"`
	// ValueFlags is the decimal bitmask of value flags written by this layer.
	ValueFlags string `json:"value_flags,omitempty" yaml:"value_flags,omitempty" xml:"value_flags,attr,omitempty"`
	// Color is the packed color used to display this layer in the editor.
	Color string `json:"color,omitempty" yaml:"color,omitempty" xml:"color,attr,omitempty"`
	// Visible controls editor visibility only; it does not change the compiled area-flags data.
	Visible string `json:"visible,omitempty" yaml:"visible,omitempty" xml:"visible,attr,omitempty"`
	// Name identifies the layer and provides the base name for its source mask.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty"`
}

// CEProjectTerritoryTypeList stores territory definitions embedded in the editor project.
// These definitions are separate from the area-flags raster.
type CEProjectTerritoryTypeList struct {
	// Types stores individual `<territory-type>` entries.
	// CEProject save/export may rewrite this section from `territoryTypes/*.xml`.
	Types []CEProjectTerritoryType `json:"territory-type,omitempty" yaml:"territory-type,omitempty" xml:"territory-type,omitempty"`
}

// CEProjectTerritoryType groups territories under the type name used for a separate territory definition file.
type CEProjectTerritoryType struct {
	// Name is the territory type identifier and commonly matches the output file name, such as `sheep_goat_territories`.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty"`
	// Territories are the named regions and their spawn zones for this type.
	Territories []CEProjectTerritory `json:"territory,omitempty" yaml:"territory,omitempty" xml:"territory,omitempty"`
}

// CEProjectTerritory stores one named territory and its editor display state.
type CEProjectTerritory struct {
	// Name identifies this territory in the editor and exported XML.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty"`
	// Visible controls whether the territory is drawn in the editor.
	Visible string `json:"visible,omitempty" yaml:"visible,omitempty" xml:"visible,attr,omitempty"`
	// Color is the packed display color used for the territory.
	Color string `json:"color,omitempty" yaml:"color,omitempty" xml:"color,attr,omitempty"`
	// Zones define the named spawn areas and their static/dynamic population limits.
	Zones []CEProjectTerritoryZone `json:"zone,omitempty" yaml:"zone,omitempty" xml:"zone,omitempty"`
}

// CEProjectTerritoryZone stores one circular spawn area in a territory.
type CEProjectTerritoryZone struct {
	// Name identifies the zone's role, such as `Graze`, `Rest`, or `Water`.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty"`
	// SMin is the minimum number of static members maintained in this zone.
	SMin string `json:"smin,omitempty" yaml:"smin,omitempty" xml:"smin,attr,omitempty"`
	// SMax is the maximum number of static members allowed in this zone.
	SMax string `json:"smax,omitempty" yaml:"smax,omitempty" xml:"smax,attr,omitempty"`
	// DMin is the minimum number of dynamic members maintained in this zone.
	DMin string `json:"dmin,omitempty" yaml:"dmin,omitempty" xml:"dmin,attr,omitempty"`
	// DMax is the maximum number of dynamic members allowed in this zone.
	DMax string `json:"dmax,omitempty" yaml:"dmax,omitempty" xml:"dmax,attr,omitempty"`
	// X is the zone center's world X coordinate.
	X string `json:"x,omitempty" yaml:"x,omitempty" xml:"x,attr,omitempty"`
	// Z is the zone center's world Z coordinate.
	Z string `json:"z,omitempty" yaml:"z,omitempty" xml:"z,attr,omitempty"`
	// R is the zone radius in world units.
	R string `json:"r,omitempty" yaml:"r,omitempty" xml:"r,attr,omitempty"`
	// D is an optional CEProject zone attribute.
	// Its runtime meaning is not established by the available format references.
	D string `json:"d,omitempty" yaml:"d,omitempty" xml:"d,attr,omitempty"`
}
