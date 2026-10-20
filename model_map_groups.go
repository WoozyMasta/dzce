// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// MapPrototypeDefaults stores values inherited by map-group or cluster prototypes when their own entries omit those properties.
type MapPrototypeDefaults struct {
	// Defaults are shared defaults applied to matching prototype entries.
	Defaults []MapPrototypeDefault `json:"default,omitempty" yaml:"default,omitempty" xml:"default,omitempty"`
}

// MapPrototypeDefault sets one shared property for a prototype category.
type MapPrototypeDefault struct {
	// Name identifies the prototype category or class this default applies to.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty" jsonschema:"example=group"`
	// LootMax is the default maximum amount of loot used for matching groups or containers.
	LootMax string `json:"lootmax,omitempty" yaml:"lootmax,omitempty" xml:"lootmax,attr,omitempty" jsonschema:"example=6"`
	// Enabled turns the selected default behavior on or off.
	Enabled string `json:"enabled,omitempty" yaml:"enabled,omitempty" xml:"enabled,attr,omitempty" jsonschema:"example=yes"`
	// DE names the dynamic-event trajectory associated with the prototype.
	DE string `json:"de,omitempty" yaml:"de,omitempty" xml:"de,attr,omitempty" jsonschema:"example=TrajectoryApple"`
	// Width sets the width of a generated cluster matrix.
	Width string `json:"width,omitempty" yaml:"width,omitempty" xml:"width,attr,omitempty" jsonschema:"example=12"`
	// Height sets the height of a generated cluster matrix.
	Height string `json:"height,omitempty" yaml:"height,omitempty" xml:"height,attr,omitempty" jsonschema:"example=12"`
}

// MapPrototypeExports maps cluster model shapes to logical cluster group names used by CE.
type MapPrototypeExports struct {
	// Exports are the model-to-cluster mappings exported by the terrain.
	Exports []MapPrototypeExport `json:"export,omitempty" yaml:"export,omitempty" xml:"export,omitempty"`
}

// MapPrototypeExport associates a model shape with the logical name used for its cluster prototype.
type MapPrototypeExport struct {
	// Name is the logical cluster group identifier referenced by prototype entries.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty" jsonschema:"example=PathD10"`
	// Shape is the P3D model path whose terrain placements use this mapping.
	Shape string `json:"shape,omitempty" yaml:"shape,omitempty" xml:"shape,attr,omitempty" jsonschema:"example=dz\\structures\\Roads\\Parts\\path_dirt_10.p3d"`
}

// MapPrototypeEntry defines one building group or cluster prototype:
// its loot limits, spawn filters, containers, and optional dynamic-event dispatches.
type MapPrototypeEntry struct {
	// Name is the prototype identifier used by map instances or cluster data.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty" jsonschema:"example=CaninaBush1"`
	// LootMax caps the total loot spawned in this group.
	LootMax string `json:"lootmax,omitempty" yaml:"lootmax,omitempty" xml:"lootmax,attr,omitempty" jsonschema:"example=2"`
	// MaxInstances caps how many instances of this cluster prototype can be used across the map.
	MaxInstances string `json:"maxinstances,omitempty" yaml:"maxinstances,omitempty" xml:"maxinstances,attr,omitempty" jsonschema:"example=150"`
	// DELinks associate this prototype with dynamic-event trajectory presets.
	DELinks []MapPrototypeDE `json:"de,omitempty" yaml:"de,omitempty" xml:"de,omitempty"`
	// Usages restrict which usage filters can select this group.
	Usages []NamedRef `json:"usage,omitempty" yaml:"usage,omitempty" xml:"usage,omitempty"`
	// Categories restrict which loot categories can use this group.
	Categories []NamedRef `json:"category,omitempty" yaml:"category,omitempty" xml:"category,omitempty"`
	// Tags restrict which loot tags can use this group.
	Tags []NamedRef `json:"tag,omitempty" yaml:"tag,omitempty" xml:"tag,omitempty"`
	// Values restrict which loot value tiers can use this group.
	Values []NamedRef `json:"value,omitempty" yaml:"value,omitempty" xml:"value,omitempty"`
	// Containers define lootable spaces and spawn points within the group.
	Containers []MapPrototypeContainer `json:"container,omitempty" yaml:"container,omitempty" xml:"container,omitempty"`
	// Dispatches define optional objects created through a dynamic event at this prototype.
	Dispatches []MapPrototypeDispatch `json:"dispatch,omitempty" yaml:"dispatch,omitempty" xml:"dispatch,omitempty"`
}

// MapPrototypeDispatch defines the chance and proxy objects for one dynamic event dispatch from a prototype.
type MapPrototypeDispatch struct {
	// DEChance is the probability that CE dispatches this dynamic event.
	DEChance string `json:"dechance,omitempty" yaml:"dechance,omitempty" xml:"dechance,attr,omitempty" jsonschema:"example=0.000000"`
	// Proxies are the objects created by the dispatch, with positions and rotations relative to the prototype.
	Proxies []MapPrototypeProxy `json:"proxy,omitempty" yaml:"proxy,omitempty" xml:"proxy,omitempty"`
}

// MapPrototypeProxy defines one object created by a dynamic-event dispatch.
type MapPrototypeProxy struct {
	// Type is the game class name of the dispatched object.
	Type string `json:"type,omitempty" yaml:"type,omitempty" xml:"type,attr,omitempty" jsonschema:"example=HatchbackWheel"`
	// Pos is the object's position relative to the prototype.
	Pos string `json:"pos,omitempty" yaml:"pos,omitempty" xml:"pos,attr,omitempty" jsonschema:"example=1.071 -0.736000 -0.852"`
	// RPY is the object's rotation relative to the prototype, in roll, pitch, yaw order.
	RPY string `json:"rpy,omitempty" yaml:"rpy,omitempty" xml:"rpy,attr,omitempty" jsonschema:"example=90.0 0.0 0.0"`
	// DEChance is the probability that this proxy is created by the dispatch.
	DEChance string `json:"dechance,omitempty" yaml:"dechance,omitempty" xml:"dechance,attr,omitempty" jsonschema:"example=0.000000"`
}

// MapPrototypeDE links a prototype to a named dynamic-event trajectory.
type MapPrototypeDE struct {
	// Name is the trajectory identifier used by the dynamic-event system.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty" jsonschema:"example=TrajectoryCanina"`
}

// MapPrototypeContainer defines one lootable part of a building prototype and the spawn points available within it.
type MapPrototypeContainer struct {
	// Name is the container identifier used by the building's game model.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty" jsonschema:"example=branch"`
	// LootMax caps how many items CE may spawn in this container. Each spawn point can hold at most one item.
	LootMax string `json:"lootmax,omitempty" yaml:"lootmax,omitempty" xml:"lootmax,attr,omitempty" jsonschema:"example=1"`
	// Usages restrict which item usage filters can select this container.
	Usages []NamedRef `json:"usage,omitempty" yaml:"usage,omitempty" xml:"usage,omitempty"`
	// Categories restrict which item categories can select this container.
	Categories []NamedRef `json:"category,omitempty" yaml:"category,omitempty" xml:"category,omitempty"`
	// Tags restrict which item tags can select this container.
	Tags []NamedRef `json:"tag,omitempty" yaml:"tag,omitempty" xml:"tag,omitempty"`
	// Values restrict which item value tiers can select this container.
	Values []NamedRef `json:"value,omitempty" yaml:"value,omitempty" xml:"value,omitempty"`
	// Points are the individual positions where CE can place one loot item.
	Points []MapPrototypePoint `json:"point,omitempty" yaml:"point,omitempty" xml:"point,omitempty"`
}

// MapPrototypePoint defines one candidate loot placement in a prototype's local model coordinates.
type MapPrototypePoint struct {
	// Pos is the local X/Y/Z position relative to the model origin, not a map position.
	Pos string `json:"pos,omitempty" yaml:"pos,omitempty" xml:"pos,attr,omitempty" jsonschema:"example=-1.86 -2.40 0.12"`
	// Range is the placement radius around the point.
	Range string `json:"range,omitempty" yaml:"range,omitempty" xml:"range,attr,omitempty" jsonschema:"example=0.55"`
	// Height is the vertical clearance available for an item at this point.
	Height string `json:"height,omitempty" yaml:"height,omitempty" xml:"height,attr,omitempty" jsonschema:"example=1.0"`
	// Flags is the bitmask of CE placement behaviors enabled for this point.
	Flags string `json:"flags,omitempty" yaml:"flags,omitempty" xml:"flags,attr,omitempty" jsonschema:"example=32"`
}

// MapGroupInstance places one prototype at a concrete map location in a map group file.
type MapGroupInstance struct {
	// Name is the prototype or group identifier resolved by CE.
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty" jsonschema:"example=PearTree2"`
	// Pos is the instance's world position as an X/Y/Z vector.
	Pos string `json:"pos,omitempty" yaml:"pos,omitempty" xml:"pos,attr,omitempty" jsonschema:"example=40.327084 10.632996 1613.091797"`
	// A is the instance's heading angle in degrees when the placement uses a single yaw value.
	A string `json:"a,omitempty" yaml:"a,omitempty" xml:"a,attr,omitempty" jsonschema:"example=1.999835"`
	// RPY is the instance's roll, pitch, and yaw rotation when the placement uses all three angles.
	RPY string `json:"rpy,omitempty" yaml:"rpy,omitempty" xml:"rpy,attr,omitempty" jsonschema:"example=-0.000000 0.000000 -70.013718"`
}
