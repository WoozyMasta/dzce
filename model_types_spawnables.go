// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// TypeDef defines how Central Economy counts and replenishes one game class.
// The nominal/min targets apply to the instances included by Flags;
// category, usage, tag, and value entries restrict the locations or contexts where the class can spawn.
type TypeDef struct {
	// Nominal is the target number of counted instances the economy maintains for this class across the world.
	Nominal *int `json:"nominal,omitempty" yaml:"nominal,omitempty" xml:"nominal,omitempty" jsonschema:"example=10"`
	// Lifetime is how long, in seconds, an eligible spawned instance can remain before economy cleanup removes it.
	Lifetime *int `json:"lifetime,omitempty" yaml:"lifetime,omitempty" xml:"lifetime,omitempty" jsonschema:"example=14400"`
	// Restock is the delay, in seconds, used by the economy before it can replenish this class after its count falls below target.
	Restock *int `json:"restock,omitempty" yaml:"restock,omitempty" xml:"restock,omitempty" jsonschema:"example=1800"`
	// Min is the lower count threshold for this class. When the counted amount falls below it, the economy prioritizes bringing the amount back up.
	Min *int `json:"min,omitempty" yaml:"min,omitempty" xml:"min,omitempty" jsonschema:"example=6"`
	// QuantityMin is the minimum fill percentage for items with a variable quantity, such as magazines, food, or liquid containers.
	QuantityMin *int `json:"quantmin,omitempty" yaml:"quantmin,omitempty" xml:"quantmin,omitempty" jsonschema:"example=-1"`
	// QuantityMax is the maximum fill percentage for items with a variable quantity, such as magazines, food, or liquid containers.
	QuantityMax *int `json:"quantmax,omitempty" yaml:"quantmax,omitempty" xml:"quantmax,omitempty" jsonschema:"example=-1"`
	// Cost is the economy cost assigned to spawning this class and is used when the CE schedules eligible types for replenishment.
	Cost *int `json:"cost,omitempty" yaml:"cost,omitempty" xml:"cost,omitempty" jsonschema:"example=100"`
	// Flags determine which existing instances count toward the economy target and whether this entry is treated as crafted or dynamic-event loot.
	Flags *TypeFlags `json:"flags,omitempty" yaml:"flags,omitempty" xml:"flags,omitempty"`
	// Name is the game configuration class name whose economy rules are set by this entry.
	Name string `json:"name" yaml:"name" xml:"name,attr" jsonschema:"example=ACOGOptic"`
	// Categories classify the item and can constrain which loot locations accept it. Names must be defined in cfglimitsdefinition.xml.
	Categories []NamedRef `json:"category,omitempty" yaml:"category,omitempty" xml:"category,omitempty"`
	// Usages identify the location contexts in which the item may spawn, such as military or town. Names must be defined in cfglimitsdefinition.xml.
	Usages []NamedRef `json:"usage,omitempty" yaml:"usage,omitempty" xml:"usage,omitempty"`
	// Tags provide additional location filters, commonly selecting placement within a building. Names must be defined in cfglimitsdefinition.xml.
	Tags []NamedRef `json:"tag,omitempty" yaml:"tag,omitempty" xml:"tag,omitempty"`
	// Values classify the item's tier or special status and can constrain its spawn locations. Names must be defined in cfglimitsdefinition.xml.
	Values []NamedRef `json:"value,omitempty" yaml:"value,omitempty" xml:"value,omitempty"`
}

// TypeFlags controls which locations contribute existing instances to the type's economy count and marks special spawn handling.
// In CE override files, specifying flags replaces the original flag set, so include every flag that should remain enabled.
type TypeFlags struct {
	// CountInCargo includes instances stored inside another item's cargo when the CE compares the current count with Nominal and Min.
	CountInCargo *int `json:"count_in_cargo,omitempty" yaml:"count_in_cargo,omitempty" xml:"count_in_cargo,attr,omitempty" jsonschema:"enum=0,enum=1,example=0"`
	// CountInHoarder includes instances kept in hoarder storage, such as tents, barrels, and buried stashes.
	CountInHoarder *int `json:"count_in_hoarder,omitempty" yaml:"count_in_hoarder,omitempty" xml:"count_in_hoarder,attr,omitempty" jsonschema:"enum=0,enum=1,example=0"`
	// CountInMap includes instances placed in the world, including loot at ordinary spawn locations.
	CountInMap *int `json:"count_in_map,omitempty" yaml:"count_in_map,omitempty" xml:"count_in_map,attr,omitempty" jsonschema:"enum=0,enum=1,example=1"`
	// CountInPlayer includes instances in player inventories when calculating the economy count.
	CountInPlayer *int `json:"count_in_player,omitempty" yaml:"count_in_player,omitempty" xml:"count_in_player,attr,omitempty" jsonschema:"enum=0,enum=1,example=0"`
	// Crafted marks a type as crafted for CE handling; it does not define a crafting recipe.
	Crafted *int `json:"crafted,omitempty" yaml:"crafted,omitempty" xml:"crafted,attr,omitempty" jsonschema:"enum=0,enum=1,example=0"`
	// Deloot marks a type for dynamic-event loot handling rather than ordinary building loot placement.
	Deloot *int `json:"deloot,omitempty" yaml:"deloot,omitempty" xml:"deloot,attr,omitempty" jsonschema:"enum=0,enum=1,example=0"`
}

// EventDef configures one dynamic event and the rules CE uses to place, populate, retain, and replenish its instances.
type EventDef struct {
	// Nominal is the target number of active instances CE tries to maintain for this event.
	Nominal *int `json:"nominal,omitempty" yaml:"nominal,omitempty" xml:"nominal,omitempty" jsonschema:"example=0"`
	// Min is the lower event count bound used by the selected Limit mode.
	Min *int `json:"min,omitempty" yaml:"min,omitempty" xml:"min,omitempty" jsonschema:"example=0"`
	// Max is the upper event count bound used by the selected Limit mode.
	Max *int `json:"max,omitempty" yaml:"max,omitempty" xml:"max,omitempty" jsonschema:"example=25"`
	// Lifetime is how long an event instance may remain active before it expires, in seconds.
	// Entity lifetime may also be controlled by its type entry.
	Lifetime *int `json:"lifetime,omitempty" yaml:"lifetime,omitempty" xml:"lifetime,omitempty" jsonschema:"example=33"`
	// Restock is the delay, in seconds, before CE may create a replacement event instance.
	Restock *int `json:"restock,omitempty" yaml:"restock,omitempty" xml:"restock,omitempty" jsonschema:"example=25"`
	// SafeRadius is the minimum clearance used by CE's event placement rules.
	SafeRadius *int `json:"saferadius,omitempty" yaml:"saferadius,omitempty" xml:"saferadius,omitempty" jsonschema:"example=0"`
	// DistanceRadius controls the distance condition for placing event instances.
	// For ambient zone events, it is the player distance at which a zone becomes eligible to spawn its entities.
	DistanceRadius *int `json:"distanceradius,omitempty" yaml:"distanceradius,omitempty" xml:"distanceradius,omitempty" jsonschema:"example=80"`
	// CleanupRadius controls when an event zone may clean up its spawned entities.
	// Ambient zones remain active while a player is within this radius.
	CleanupRadius *int `json:"cleanupradius,omitempty" yaml:"cleanupradius,omitempty" xml:"cleanupradius,omitempty" jsonschema:"example=120"`
	// Secondary names an additional event CE may spawn for each instance, such as the infected event paired with a crash site or convoy.
	Secondary *string `json:"secondary,omitempty" yaml:"secondary,omitempty" xml:"secondary,omitempty" jsonschema:"example=InfectedArmyDesert"`
	// Flags control event cleanup and initial-spawn behavior.
	Flags *EventFlags `json:"flags,omitempty" yaml:"flags,omitempty" xml:"flags,omitempty"`
	// Position selects how CE chooses event locations, for example from fixed mission positions or from a generated distribution.
	Position *string `json:"position,omitempty" yaml:"position,omitempty" xml:"position,omitempty" jsonschema:"example=fixed"`
	// Limit selects how Min and Max apply to the child entries and their combined event count.
	Limit *string `json:"limit,omitempty" yaml:"limit,omitempty" xml:"limit,omitempty" jsonschema:"example=mixed"`
	// Active enables or disables spawning for this event (`1` or `0`).
	Active *int `json:"active,omitempty" yaml:"active,omitempty" xml:"active,omitempty" jsonschema:"enum=0,enum=1,example=1"`
	// Children define the entity classes or objects created by this event and any per-child count or loot limits.
	Children *EventChildren `json:"children,omitempty" yaml:"children,omitempty" xml:"children,omitempty"`
	// Name is the event identifier referenced by position and group data.
	Name string `json:"name" yaml:"name" xml:"name,attr" jsonschema:"example=AmbientFox"`
}

// EventFlags controls whether CE can delete an event instance, randomizes its initial spawning, and removes damaged spawned entities.
type EventFlags struct {
	// Deletable allows CE to remove this event instance during cleanup.
	Deletable *int `json:"deletable,omitempty" yaml:"deletable,omitempty" xml:"deletable,attr,omitempty" jsonschema:"enum=0,enum=1,example=0"`
	// InitRandom randomizes the event's initial spawn state.
	InitRandom *int `json:"init_random,omitempty" yaml:"init_random,omitempty" xml:"init_random,attr,omitempty" jsonschema:"enum=0,enum=1,example=0"`
	// RemoveDamaged allows CE to remove damaged or ruined entities created by this event so they can be replaced.
	RemoveDamaged *int `json:"remove_damaged,omitempty" yaml:"remove_damaged,omitempty" xml:"remove_damaged,attr,omitempty" jsonschema:"enum=0,enum=1,example=0"`
}

// EventChildren groups the possible entity types or objects created by an event.
type EventChildren struct {
	// Children are the spawnable entries and their per-type limits or weights.
	Children []EventChild `json:"child,omitempty" yaml:"child,omitempty" xml:"child,omitempty"`
}

// EventChild defines one entity type or object that an event can spawn.
type EventChild struct {
	// LootMax caps the amount of loot generated around this child, where the event type uses child loot limits.
	LootMax *int `json:"lootmax,omitempty" yaml:"lootmax,omitempty" xml:"lootmax,attr,omitempty" jsonschema:"example=0"`
	// LootMin sets the lower loot amount generated around this child, where the event type uses child loot limits.
	LootMin *int `json:"lootmin,omitempty" yaml:"lootmin,omitempty" xml:"lootmin,attr,omitempty" jsonschema:"example=0"`
	// Max is the upper number of instances for this child type when the event's limit mode applies bounds per child.
	Max *int `json:"max,omitempty" yaml:"max,omitempty" xml:"max,attr,omitempty" jsonschema:"example=0"`
	// Min is the lower number of instances for this child type when the event's limit mode applies bounds per child.
	// Ambient spawners use this value as a relative weight when selecting among child types.
	Min *int `json:"min,omitempty" yaml:"min,omitempty" xml:"min,attr,omitempty" jsonschema:"example=100"`
	// Type is the game class name created for this child entry.
	Type string `json:"type" yaml:"type" xml:"type,attr" jsonschema:"example=Animal_VulpesVulpes"`
}

// EconomySection selects which lifecycle operations Central Economy performs for one subsystem:
// initialization, loading and saving persistent state, and runtime respawning.
type EconomySection struct {
	// Init allows CE to initialize this subsystem when the mission starts.
	Init *int `json:"init,omitempty" yaml:"init,omitempty" xml:"init,attr,omitempty" jsonschema:"enum=0,enum=1,example=1"`
	// Load allows CE to restore this subsystem's saved state from server storage.
	Load *int `json:"load,omitempty" yaml:"load,omitempty" xml:"load,attr,omitempty" jsonschema:"enum=0,enum=1,example=1"`
	// Respawn allows CE to replenish entities managed by this subsystem during runtime.
	Respawn *int `json:"respawn,omitempty" yaml:"respawn,omitempty" xml:"respawn,attr,omitempty" jsonschema:"enum=0,enum=1,example=1"`
	// Save allows CE to write this subsystem's state to server storage.
	Save *int `json:"save,omitempty" yaml:"save,omitempty" xml:"save,attr,omitempty" jsonschema:"enum=0,enum=1,example=1"`
}

// GlobalVar overrides one named Central Economy variable.
// Variable names are defined by the game;
// missions should preserve the declared type of existing variables when overriding them.
type GlobalVar struct {
	// Name is the built-in CE variable identifier, such as `AnimalMaxCount` or `CleanupLifetimeDefault`.
	Name string `json:"name" yaml:"name" xml:"name,attr" jsonschema:"example=AnimalMaxCount"`
	// Value is the variable's configured value encoded as text and interpreted according to Type.
	Value string `json:"value" yaml:"value" xml:"value,attr" jsonschema:"example=200"`
	// Type tells CE how to parse Value: 0 for integer, 1 for floating-point, and 2 for string.
	Type VariableType `json:"type" yaml:"type" xml:"type,attr" jsonschema:"enum=0,enum=1,enum=2,example=0"`
}

// MessageDef is one `<message>` entry in messages.xml.
type MessageDef struct {
	// Deadline is the countdown duration in minutes.
	// A positive value enables countdown messages and determines when the countdown reaches zero.
	Deadline *int `json:"deadline,omitempty" yaml:"deadline,omitempty" xml:"deadline,omitempty" jsonschema:"minimum=0"`
	// Shutdown set to 1 stops the server when the deadline countdown reaches zero.
	// It has no effect when Deadline is zero or omitted.
	Shutdown *int `json:"shutdown,omitempty" yaml:"shutdown,omitempty" xml:"shutdown,omitempty" jsonschema:"enum=0,enum=1"`
	// Delay is the wait in minutes before sending an OnConnect message.
	// Zero sends it immediately after the player connects.
	Delay *int `json:"delay,omitempty" yaml:"delay,omitempty" xml:"delay,omitempty" jsonschema:"minimum=0"`
	// Repeat is the interval in minutes between message deliveries.
	// A positive value enables recurring delivery; zero disables it.
	Repeat *int `json:"repeat,omitempty" yaml:"repeat,omitempty" xml:"repeat,omitempty" jsonschema:"minimum=0"`
	// OnConnect set to 1 sends the message once after each player connects; 0 disables connection-triggered delivery.
	// Delay sets the wait before sending.
	OnConnect *int `json:"onconnect,omitempty" yaml:"onconnect,omitempty" xml:"onconnect,omitempty" jsonschema:"enum=0,enum=1"`
	// Text is localized message template text.
	Text string `json:"text,omitempty" yaml:"text,omitempty" xml:"text,omitempty"`
}

// SpawnableTypeDef defines randomized cargo and attachment contents for one item class.
// The class must also be part of the economy's spawnable types for its spawn rules to be used.
// In a CE override, specifying cargo or attachments replaces the corresponding original list for that class.
type SpawnableTypeDef struct {
	// Damage sets the minimum and maximum health for the spawned item.
	// In an override file, this setting affects only new or redefined types.
	Damage *SpawnableMinMax `json:"damage,omitempty" yaml:"damage,omitempty" xml:"damage,omitempty"`
	// Hoarder marks this class for hoarder-related spawn handling.
	Hoarder *EmptyElement `json:"hoarder,omitempty" yaml:"hoarder,omitempty" xml:"hoarder,omitempty"`
	// Unique marks this type's spawn rules as unique. In an override file, this flag replaces the original type's flag set.
	Unique *EmptyElement `json:"unique,omitempty" yaml:"unique,omitempty" xml:"unique,omitempty"`
	// Name is the game class whose cargo and attachment spawn rules are defined.
	Name string `json:"name" yaml:"name" xml:"name,attr" jsonschema:"example=Barrel_Blue"`
	// Tags attach named filters used to select compatible spawn presets.
	Tags []NamedRef `json:"tag,omitempty" yaml:"tag,omitempty" xml:"tag,omitempty"`
	// Cargo defines the items or named preset that can be generated inside this item's cargo.
	// Defining cargo in an override replaces the original cargo rules for this type.
	Cargo []SpawnableCargo `json:"cargo,omitempty" yaml:"cargo,omitempty" xml:"cargo,omitempty"`
	// Attachments defines the items or named preset that can be attached to this item.
	// Defining attachments in an override replaces the original attachment rules for this type.
	Attachments []SpawnableAttachment `json:"attachments,omitempty" yaml:"attachments,omitempty" xml:"attachments,omitempty"`
}

// SpawnableMinMax stores a lower and upper bound used to randomize a numeric property when an item is created.
type SpawnableMinMax struct {
	// Min is the inclusive lower bound for the generated property.
	Min *float64 `json:"min,omitempty" yaml:"min,omitempty" xml:"min,attr,omitempty" jsonschema:"example=0"`
	// Max is the inclusive upper bound for the generated property.
	Max *float64 `json:"max,omitempty" yaml:"max,omitempty" xml:"max,attr,omitempty" jsonschema:"example=0.32"`
}

// SpawnableCargo defines one cargo-generation rule for an item.
// It can either reference a named preset or list item candidates directly.
type SpawnableCargo struct {
	// Chance is the probability, from 0 to 1, that this cargo rule is selected.
	Chance *float64 `json:"chance,omitempty" yaml:"chance,omitempty" xml:"chance,attr,omitempty" jsonschema:"minimum=0,maximum=1,example=0.6"`
	// Preset names a cargo preset from cfgrandompresets.xml.
	// When set, the preset provides the candidate items for this cargo rule.
	Preset string `json:"preset,omitempty" yaml:"preset,omitempty" xml:"preset,attr,omitempty" jsonschema:"example=mixArmy"`
	// Items are item candidates listed directly instead of in a named preset.
	Items []SpawnableItem `json:"item,omitempty" yaml:"item,omitempty" xml:"item,omitempty"`
}

// SpawnableAttachment defines one attachment-generation rule for an item.
// It can reference a named preset or list attachment candidates directly.
type SpawnableAttachment struct {
	// Chance is the probability, from 0 to 1, that this attachment rule is selected.
	Chance *float64 `json:"chance,omitempty" yaml:"chance,omitempty" xml:"chance,attr,omitempty" jsonschema:"minimum=0,maximum=1,example=0.2"`
	// Preset names an attachment preset from cfgrandompresets.xml.
	// When set, the preset provides the candidate items for this attachment rule.
	Preset string `json:"preset,omitempty" yaml:"preset,omitempty" xml:"preset,attr,omitempty" jsonschema:"example=headtorches"`
	// Items are attachment candidates listed directly instead of in a named preset.
	Items []SpawnableItem `json:"item,omitempty" yaml:"item,omitempty" xml:"item,omitempty"`
}

// SpawnableItem is one item candidate in cargo, attachment, or random-preset rules.
// Its nested rules can define contents and attachments recursively.
type SpawnableItem struct {
	// Damage sets the health range for this item when it is spawned.
	Damage *SpawnableMinMax `json:"damage,omitempty" yaml:"damage,omitempty" xml:"damage,omitempty"`
	// Chance is a probability from 0 to 1, or a percentage weight
	// when the explicit chance values in its parent group use the 0 to 100 scale.
	Chance *float64 `json:"chance,omitempty" yaml:"chance,omitempty" xml:"chance,attr,omitempty" jsonschema:"minimum=0,example=0.11"`
	// Name is the game class name of the candidate item.
	Name string `json:"name" yaml:"name" xml:"name,attr" jsonschema:"example=TunaCan"`
	// Cargo defines items that may be generated inside this item.
	Cargo []SpawnableCargo `json:"cargo,omitempty" yaml:"cargo,omitempty" xml:"cargo,omitempty"`
	// Attachments defines items that may be attached to this item.
	Attachments []SpawnableAttachment `json:"attachments,omitempty" yaml:"attachments,omitempty" xml:"attachments,omitempty"`
}

// RandomPreset is a reusable item list selected by name from cargo or attachment rules in cfgspawnabletypes.xml.
type RandomPreset struct {
	// Chance is the probability, from 0 to 1, that this preset is selected when it is one of the candidates.
	Chance *float64 `json:"chance,omitempty" yaml:"chance,omitempty" xml:"chance,attr,omitempty" jsonschema:"minimum=0,maximum=1,example=0.15"`
	// Name is the identifier referenced by a spawn rule's `preset` attribute.
	Name string `json:"name" yaml:"name" xml:"name,attr" jsonschema:"example=foodHermit"`
	// Items are the candidates provided by this preset.
	Items []SpawnableItem `json:"item,omitempty" yaml:"item,omitempty" xml:"item,omitempty"`
}
