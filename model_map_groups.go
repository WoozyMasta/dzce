// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// MapPrototypeDefaults stores `<defaults>` block values.
type MapPrototypeDefaults struct {
	// Defaults is list of default entries.
	Defaults []MapPrototypeDefault `xml:"default,omitempty" json:"default,omitempty" yaml:"default,omitempty"`
}

// MapPrototypeDefault stores one `<default />` entry.
type MapPrototypeDefault struct {
	// Name is default type identifier.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	// LootMax is default loot limit for group/container.
	LootMax string `xml:"lootmax,attr,omitempty" json:"lootmax,omitempty" yaml:"lootmax,omitempty"`
	// Enabled toggles selected default feature.
	Enabled string `xml:"enabled,attr,omitempty" json:"enabled,omitempty" yaml:"enabled,omitempty"`
	// DE is dynamic event trajectory name.
	DE string `xml:"de,attr,omitempty" json:"de,omitempty" yaml:"de,omitempty"`
	// Width is cluster matrix width.
	Width string `xml:"width,attr,omitempty" json:"width,omitempty" yaml:"width,omitempty"`
	// Height is cluster matrix height.
	Height string `xml:"height,attr,omitempty" json:"height,omitempty" yaml:"height,omitempty"`
}

// MapPrototypeExports stores `<clusters><export .../></clusters>` mappings.
type MapPrototypeExports struct {
	// Exports is list of name-to-shape mappings.
	Exports []MapPrototypeExport `xml:"export,omitempty" json:"export,omitempty" yaml:"export,omitempty"`
}

// MapPrototypeExport stores one cluster export mapping.
type MapPrototypeExport struct {
	// Name is logical cluster group name.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	// Shape is p3d model path.
	Shape string `xml:"shape,attr,omitempty" json:"shape,omitempty" yaml:"shape,omitempty"`
}

// MapPrototypeEntry stores one `<group>` or `<cluster>` entry.
type MapPrototypeEntry struct {
	// Name is prototype group name.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	// LootMax limits amount for this group.
	LootMax string `xml:"lootmax,attr,omitempty" json:"lootmax,omitempty" yaml:"lootmax,omitempty"`
	// MaxInstances limits cluster instances across map.
	MaxInstances string `xml:"maxinstances,attr,omitempty" json:"maxinstances,omitempty" yaml:"maxinstances,omitempty"`

	// DELinks references dynamic trajectory presets.
	DELinks []MapPrototypeDE `xml:"de,omitempty" json:"de,omitempty" yaml:"de,omitempty"`
	// Usages stores usage limiter values.
	Usages []NamedRef `xml:"usage,omitempty" json:"usage,omitempty" yaml:"usage,omitempty"`
	// Categories stores category limiter values.
	Categories []NamedRef `xml:"category,omitempty" json:"category,omitempty" yaml:"category,omitempty"`
	// Tags stores tag limiter values.
	Tags []NamedRef `xml:"tag,omitempty" json:"tag,omitempty" yaml:"tag,omitempty"`
	// Values stores value limiter values.
	Values []NamedRef `xml:"value,omitempty" json:"value,omitempty" yaml:"value,omitempty"`
	// Containers stores nested loot containers.
	Containers []MapPrototypeContainer `xml:"container,omitempty" json:"container,omitempty" yaml:"container,omitempty"`
}

// MapPrototypeDE stores one `<de name="..."/>` link.
type MapPrototypeDE struct {
	// Name is trajectory preset identifier.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
}

// MapPrototypeContainer stores one prototype container block.
type MapPrototypeContainer struct {
	// Name is container identifier.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	// LootMax limits spawned items in this container.
	LootMax string `xml:"lootmax,attr,omitempty" json:"lootmax,omitempty" yaml:"lootmax,omitempty"`

	// Usages stores usage limiter values.
	Usages []NamedRef `xml:"usage,omitempty" json:"usage,omitempty" yaml:"usage,omitempty"`
	// Categories stores category limiter values.
	Categories []NamedRef `xml:"category,omitempty" json:"category,omitempty" yaml:"category,omitempty"`
	// Tags stores tag limiter values.
	Tags []NamedRef `xml:"tag,omitempty" json:"tag,omitempty" yaml:"tag,omitempty"`
	// Values stores value limiter values.
	Values []NamedRef `xml:"value,omitempty" json:"value,omitempty" yaml:"value,omitempty"`
	// Points stores loot spawn points in container space.
	Points []MapPrototypePoint `xml:"point,omitempty" json:"point,omitempty" yaml:"point,omitempty"`
}

// MapPrototypePoint stores one loot spawn point in prototype space.
type MapPrototypePoint struct {
	// Pos is x/y/z local position string.
	Pos string `xml:"pos,attr,omitempty" json:"pos,omitempty" yaml:"pos,omitempty"`
	// Range is point radius.
	Range string `xml:"range,attr,omitempty" json:"range,omitempty" yaml:"range,omitempty"`
	// Height is point height allowance.
	Height string `xml:"height,attr,omitempty" json:"height,omitempty" yaml:"height,omitempty"`
	// Flags is CE point flags bitmask.
	Flags string `xml:"flags,attr,omitempty" json:"flags,omitempty" yaml:"flags,omitempty"`
}

// MapGroupInstance stores one `<group .../>` in map group files.
type MapGroupInstance struct {
	// Name is prototype/group name.
	Name string `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	// Pos is world position vector as text.
	Pos string `xml:"pos,attr,omitempty" json:"pos,omitempty" yaml:"pos,omitempty"`
	// A is heading angle in degrees.
	A string `xml:"a,attr,omitempty" json:"a,omitempty" yaml:"a,omitempty"`
	// RPY is roll/pitch/yaw vector as text.
	RPY string `xml:"rpy,attr,omitempty" json:"rpy,omitempty" yaml:"rpy,omitempty"`
}
