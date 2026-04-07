// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// TypeDef is a single `<type>` entry in types.xml.
type TypeDef struct {
	// Name is game config class name.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
	// Nominal is target amount in economy.
	Nominal *int `xml:"nominal,omitempty" json:"nominal,omitempty" yaml:"nominal,omitempty"`
	// Lifetime is cleanup lifetime in seconds.
	Lifetime *int `xml:"lifetime,omitempty" json:"lifetime,omitempty" yaml:"lifetime,omitempty"`
	// Restock is respawn delay in seconds.
	Restock *int `xml:"restock,omitempty" json:"restock,omitempty" yaml:"restock,omitempty"`
	// Min is minimum amount target.
	Min *int `xml:"min,omitempty" json:"min,omitempty" yaml:"min,omitempty"`
	// QuantityMin is item minimum quantity percent.
	QuantityMin *int `xml:"quantmin,omitempty" json:"quantmin,omitempty" yaml:"quantmin,omitempty"`
	// QuantityMax is item maximum quantity percent.
	QuantityMax *int `xml:"quantmax,omitempty" json:"quantmax,omitempty" yaml:"quantmax,omitempty"`
	// Cost is respawn cost weight.
	Cost *int `xml:"cost,omitempty" json:"cost,omitempty" yaml:"cost,omitempty"`
	// Flags controls count/crafted/deloot semantics.
	Flags *TypeFlags `xml:"flags,omitempty" json:"flags,omitempty" yaml:"flags,omitempty"`
	// Categories are `<category name="..."/>` entries.
	Categories []NamedRef `xml:"category,omitempty" json:"category,omitempty" yaml:"category,omitempty"`
	// Usages are `<usage name="..."/>` entries.
	Usages []NamedRef `xml:"usage,omitempty" json:"usage,omitempty" yaml:"usage,omitempty"`
	// Tags are `<tag name="..."/>` entries.
	Tags []NamedRef `xml:"tag,omitempty" json:"tag,omitempty" yaml:"tag,omitempty"`
	// Values are `<value name="..."/>` entries.
	Values []NamedRef `xml:"value,omitempty" json:"value,omitempty" yaml:"value,omitempty"`
}

// TypeFlags stores attributes of `<flags />` in types.xml.
type TypeFlags struct {
	// CountInCargo includes items currently inside containers.
	CountInCargo *int `xml:"count_in_cargo,attr,omitempty" json:"count_in_cargo,omitempty" yaml:"count_in_cargo,omitempty"`
	// CountInHoarder includes buried stashes.
	CountInHoarder *int `xml:"count_in_hoarder,attr,omitempty" json:"count_in_hoarder,omitempty" yaml:"count_in_hoarder,omitempty"`
	// CountInMap includes world-spawned items.
	CountInMap *int `xml:"count_in_map,attr,omitempty" json:"count_in_map,omitempty" yaml:"count_in_map,omitempty"`
	// CountInPlayer includes items on connected players.
	CountInPlayer *int `xml:"count_in_player,attr,omitempty" json:"count_in_player,omitempty" yaml:"count_in_player,omitempty"`
	// Crafted marks items that can be crafted instead of spawned.
	Crafted *int `xml:"crafted,attr,omitempty" json:"crafted,omitempty" yaml:"crafted,omitempty"`
	// Deloot controls delayed loot cleanup behavior.
	Deloot *int `xml:"deloot,attr,omitempty" json:"deloot,omitempty" yaml:"deloot,omitempty"`
}

// EventDef is a single `<event>` entry in events.xml.
type EventDef struct {
	// Nominal is target active event amount.
	Nominal *int `xml:"nominal,omitempty" json:"nominal,omitempty" yaml:"nominal,omitempty"`
	// Min is lower event amount bound.
	Min *int `xml:"min,omitempty" json:"min,omitempty" yaml:"min,omitempty"`
	// Max is upper event amount bound.
	Max *int `xml:"max,omitempty" json:"max,omitempty" yaml:"max,omitempty"`
	// Lifetime is event lifetime in seconds.
	Lifetime *int `xml:"lifetime,omitempty" json:"lifetime,omitempty" yaml:"lifetime,omitempty"`
	// Restock is event respawn delay in seconds.
	Restock *int `xml:"restock,omitempty" json:"restock,omitempty" yaml:"restock,omitempty"`
	// SafeRadius is protected radius around active events.
	SafeRadius *int `xml:"saferadius,omitempty" json:"saferadius,omitempty" yaml:"saferadius,omitempty"`
	// DistanceRadius is minimum distance for new spawns.
	DistanceRadius *int `xml:"distanceradius,omitempty" json:"distanceradius,omitempty" yaml:"distanceradius,omitempty"`
	// CleanupRadius is cleanup scan radius.
	CleanupRadius *int `xml:"cleanupradius,omitempty" json:"cleanupradius,omitempty" yaml:"cleanupradius,omitempty"`
	// Secondary links this event to a secondary event name.
	Secondary *string `xml:"secondary,omitempty" json:"secondary,omitempty" yaml:"secondary,omitempty"`
	// Flags stores event runtime behavior toggles.
	Flags *EventFlags `xml:"flags,omitempty" json:"flags,omitempty" yaml:"flags,omitempty"`
	// Position selects position generation mode.
	Position *string `xml:"position,omitempty" json:"position,omitempty" yaml:"position,omitempty"`
	// Limit selects CE limit pool mode.
	Limit *string `xml:"limit,omitempty" json:"limit,omitempty" yaml:"limit,omitempty"`
	// Active enables or disables event processing.
	Active *int `xml:"active,omitempty" json:"active,omitempty" yaml:"active,omitempty"`
	// Children defines objects spawned by the event.
	Children *EventChildren `xml:"children,omitempty" json:"children,omitempty" yaml:"children,omitempty"`
	// Name is unique event identifier.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
}

// EventFlags stores attributes of `<flags />` in events.xml.
type EventFlags struct {
	// Deletable allows CE to remove this event instance.
	Deletable *int `xml:"deletable,attr,omitempty" json:"deletable,omitempty" yaml:"deletable,omitempty"`
	// InitRandom randomizes initial spawn state.
	InitRandom *int `xml:"init_random,attr,omitempty" json:"init_random,omitempty" yaml:"init_random,omitempty"`
	// RemoveDamaged removes damaged spawned entities.
	RemoveDamaged *int `xml:"remove_damaged,attr,omitempty" json:"remove_damaged,omitempty" yaml:"remove_damaged,omitempty"`
}

// EventChildren wraps event children entries.
type EventChildren struct {
	// Children is list of objects that belong to the event.
	Children []EventChild `xml:"child,omitempty" json:"child,omitempty" yaml:"child,omitempty"`
}

// EventChild is a `<child />` entry in events.xml.
type EventChild struct {
	// LootMax limits max loot spawned in this child.
	LootMax *int `xml:"lootmax,attr,omitempty" json:"lootmax,omitempty" yaml:"lootmax,omitempty"`
	// LootMin limits min loot spawned in this child.
	LootMin *int `xml:"lootmin,attr,omitempty" json:"lootmin,omitempty" yaml:"lootmin,omitempty"`
	// Max is maximum allowed child object amount.
	Max *int `xml:"max,attr,omitempty" json:"max,omitempty" yaml:"max,omitempty"`
	// Min is minimum value used by CE spawner logic for this child.
	// For some spawners (for example ambient), wiki describes it as weight.
	Min *int `xml:"min,attr,omitempty" json:"min,omitempty" yaml:"min,omitempty"`
	// Type is spawned entity class name.
	Type string `xml:"type,attr" json:"type" yaml:"type"`
}

// EconomySection is a `<dynamic/>` or similar section in economy.xml.
type EconomySection struct {
	// Init enables first-time initialization of this CE section.
	Init *int `xml:"init,attr,omitempty" json:"init,omitempty" yaml:"init,omitempty"`
	// Load enables loading section state from storage.
	Load *int `xml:"load,attr,omitempty" json:"load,omitempty" yaml:"load,omitempty"`
	// Respawn enables runtime respawn processing.
	Respawn *int `xml:"respawn,attr,omitempty" json:"respawn,omitempty" yaml:"respawn,omitempty"`
	// Save enables writing state to storage.
	Save *int `xml:"save,attr,omitempty" json:"save,omitempty" yaml:"save,omitempty"`
}

// GlobalVar is a single `<var />` entry in globals.xml.
type GlobalVar struct {
	// Name is CE variable key.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
	// Value is raw variable payload as text.
	Value string `xml:"value,attr" json:"value" yaml:"value"`
	// Type defines Value interpretation in globals.xml:
	// 0=int, 1=float, 2=string.
	Type VariableType `xml:"type,attr" json:"type" yaml:"type"`
}

// MessageDef is one `<message>` entry in messages.xml.
type MessageDef struct {
	// Deadline is countdown until server shutdown message trigger (minutes).
	Deadline *int `xml:"deadline,omitempty" json:"deadline,omitempty" yaml:"deadline,omitempty"`
	// Shutdown controls whether server shutdown is triggered.
	Shutdown *int `xml:"shutdown,omitempty" json:"shutdown,omitempty" yaml:"shutdown,omitempty"`
	// Delay is delay before first message display (minutes).
	Delay *int `xml:"delay,omitempty" json:"delay,omitempty" yaml:"delay,omitempty"`
	// Repeat is repeat period for message display (minutes).
	Repeat *int `xml:"repeat,omitempty" json:"repeat,omitempty" yaml:"repeat,omitempty"`
	// OnConnect controls whether message is shown after player connection.
	// In wiki examples this is represented as 0/1 integer.
	OnConnect *int `xml:"onconnect,omitempty" json:"onconnect,omitempty" yaml:"onconnect,omitempty"`
	// Text is localized message template text.
	Text string `xml:"text,omitempty" json:"text,omitempty" yaml:"text,omitempty"`
}

// SpawnableTypeDef is a single `<type>` entry in cfgspawnabletypes.xml.
type SpawnableTypeDef struct {
	// Name is base item class to which rules apply.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
	// Damage bounds spawned item health range.
	Damage *SpawnableMinMax `xml:"damage,omitempty" json:"damage,omitempty" yaml:"damage,omitempty"`
	// Hoarder enables use in hoarder/stash spawn logic.
	Hoarder *EmptyElement `xml:"hoarder,omitempty" json:"hoarder,omitempty" yaml:"hoarder,omitempty"`
	// Unique marks spawn list as unique-only.
	Unique *EmptyElement `xml:"unique,omitempty" json:"unique,omitempty" yaml:"unique,omitempty"`
	// Tags maps item to preset tag filters.
	Tags []NamedRef `xml:"tag,omitempty" json:"tag,omitempty" yaml:"tag,omitempty"`
	// Cargo defines nested cargo item generation.
	Cargo []SpawnableCargo `xml:"cargo,omitempty" json:"cargo,omitempty" yaml:"cargo,omitempty"`
	// Attachments defines attachment generation.
	Attachments []SpawnableAttachment `xml:"attachments,omitempty" json:"attachments,omitempty" yaml:"attachments,omitempty"`
}

// SpawnableMinMax stores `min/max` attributes for damage and similar nodes.
type SpawnableMinMax struct {
	// Min is minimum allowed value.
	Min *float64 `xml:"min,attr,omitempty" json:"min,omitempty" yaml:"min,omitempty"`
	// Max is maximum allowed value.
	Max *float64 `xml:"max,attr,omitempty" json:"max,omitempty" yaml:"max,omitempty"`
}

// SpawnableCargo is `<cargo ...>` in cfgspawnabletypes.xml.
type SpawnableCargo struct {
	// Preset links cargo item list to named random preset.
	Preset string `xml:"preset,attr,omitempty" json:"preset,omitempty" yaml:"preset,omitempty"`
	// Chance is probability for this cargo branch.
	Chance *float64 `xml:"chance,attr,omitempty" json:"chance,omitempty" yaml:"chance,omitempty"`
	// Items is explicit cargo item candidates.
	Items []SpawnableItem `xml:"item,omitempty" json:"item,omitempty" yaml:"item,omitempty"`
}

// SpawnableAttachment is `<attachments ...>` in cfgspawnabletypes.xml.
type SpawnableAttachment struct {
	// Chance is probability for this attachment branch.
	Chance *float64 `xml:"chance,attr,omitempty" json:"chance,omitempty" yaml:"chance,omitempty"`
	// Items is explicit attachment item candidates.
	Items []SpawnableItem `xml:"item,omitempty" json:"item,omitempty" yaml:"item,omitempty"`
}

// SpawnableItem is `<item .../>` in spawnable types and random presets.
type SpawnableItem struct {
	// Chance is per-item probability inside parent list.
	Chance *float64 `xml:"chance,attr,omitempty" json:"chance,omitempty" yaml:"chance,omitempty"`
	// Name is class name of spawned item.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
}

// RandomPreset is a named preset for cargo or attachments.
type RandomPreset struct {
	// Name is preset identifier used by `preset="..."`.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
	// Chance is probability of choosing this preset.
	Chance *float64 `xml:"chance,attr,omitempty" json:"chance,omitempty" yaml:"chance,omitempty"`
	// Items is item list inside the preset.
	Items []SpawnableItem `xml:"item,omitempty" json:"item,omitempty" yaml:"item,omitempty"`
}
