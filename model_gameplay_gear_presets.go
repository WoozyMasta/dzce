// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// GameplayGearPresetsFile defines character equipment presets loaded through `cfggameplay.json -> PlayerData.spawnGearPresetFiles`.
// The game selects a preset by weight during character creation and uses it to override the mission's `StartingEquipSetup` behavior.
// If the mission overrides `CreateCharacter` and chooses a character class there, `characterTypes` in these presets is ignored.
//
// The root JSON value is an array of presets.
type GameplayGearPresetsFile []GameplayGearPreset

// GameplayGearPreset defines the character classes, equipment slots, and cargo variants available for one starting loadout.
type GameplayGearPreset struct {
	// Name labels the preset for identification in configuration and logs.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// SpawnWeight is the preset's relative selection weight.
	// It must be at least 1; higher values make the preset more likely to be selected.
	SpawnWeight *int `json:"spawnWeight,omitempty" yaml:"spawnWeight,omitempty"`
	// CharacterTypes lists survivor classes eligible for this preset.
	// When omitted or empty, the character type selected in the creation menu is used.
	CharacterTypes []string `json:"characterTypes,omitempty" yaml:"characterTypes,omitempty"`
	// AttachmentSlotItemSets defines weighted equipment choices for character attachment slots such as body, back, and shoulders.
	AttachmentSlotItemSets []GameplayGearAttachmentSlotSet `json:"attachmentSlotItemSets,omitempty" yaml:"attachmentSlotItemSets,omitempty"`
	// DiscreteUnsortedItemSets defines weighted cargo loadouts added to the character independently of a specific attachment slot.
	DiscreteUnsortedItemSets []GameplayGearUnsortedSet `json:"discreteUnsortedItemSets,omitempty" yaml:"discreteUnsortedItemSets,omitempty"`
}

// GameplayGearAttachmentSlotSet defines the candidate items for one character attachment slot.
type GameplayGearAttachmentSlotSet struct {
	// SlotName is the inventory slot name from CfgSlots, such as `Body`, `Back`, or `shoulderL`.
	SlotName string `json:"slotName,omitempty" yaml:"slotName,omitempty"`
	// DiscreteItemSets are the weighted alternative items that may occupy this slot.
	DiscreteItemSets []GameplayGearDiscreteItemSet `json:"discreteItemSets,omitempty" yaml:"discreteItemSets,omitempty"`
}

// GameplayGearDiscreteItemSet defines one weighted item choice and any cargo or attachments created inside it.
type GameplayGearDiscreteItemSet struct {
	// SpawnWeight is this item's relative selection weight among alternatives for the slot or set. It must be at least 1.
	SpawnWeight *int `json:"spawnWeight,omitempty" yaml:"spawnWeight,omitempty"`
	// Attributes supplies the health and quantity ranges rolled for this item.
	Attributes *GameplayGearItemAttributes `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	// QuickBarSlot assigns the item to a quickbar slot; -1 leaves it unassigned.
	QuickBarSlot *int `json:"quickBarSlot,omitempty" yaml:"quickBarSlot,omitempty"`
	// SimpleChildrenUseDefaultAttributes makes simple child items inherit the attributes from this item set instead of using configuration defaults.
	SimpleChildrenUseDefaultAttributes *bool `json:"simpleChildrenUseDefaultAttributes,omitempty" yaml:"simpleChildrenUseDefaultAttributes,omitempty"`
	// ItemType is the game class name of the item created in this slot.
	ItemType string `json:"itemType,omitempty" yaml:"itemType,omitempty"`
	// ComplexChildrenTypes are items created inside the selected item, each with its own attributes, quickbar assignment, and optional simple cargo.
	ComplexChildrenTypes []GameplayGearComplexChild `json:"complexChildrenTypes,omitempty" yaml:"complexChildrenTypes,omitempty"`
	// SimpleChildrenTypes are game class names of child items created inside the selected item without per-child settings.
	SimpleChildrenTypes []string `json:"simpleChildrenTypes,omitempty" yaml:"simpleChildrenTypes,omitempty"`
}

// GameplayGearUnsortedSet defines one weighted cargo loadout that can be selected independently of a character attachment slot.
type GameplayGearUnsortedSet struct {
	// SpawnWeight is this cargo loadout's relative selection weight among alternatives. It must be at least 1.
	SpawnWeight *int `json:"spawnWeight,omitempty" yaml:"spawnWeight,omitempty"`
	// Attributes provides default health and quantity ranges for nested items that inherit attributes from this set.
	Attributes *GameplayGearItemAttributes `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	// SimpleChildrenUseDefaultAttributes makes simple child items inherit the attributes from this cargo set instead of using configuration defaults.
	SimpleChildrenUseDefaultAttributes *bool `json:"simpleChildrenUseDefaultAttributes,omitempty" yaml:"simpleChildrenUseDefaultAttributes,omitempty"`
	// Name identifies the cargo loadout in configuration and diagnostics.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// ComplexChildrenTypes are cargo items with their own attributes, quickbar assignment, and optional simple child items.
	ComplexChildrenTypes []GameplayGearComplexChild `json:"complexChildrenTypes,omitempty" yaml:"complexChildrenTypes,omitempty"`
	// SimpleChildrenTypes are game class names of cargo items created without per-item settings.
	SimpleChildrenTypes []string `json:"simpleChildrenTypes,omitempty" yaml:"simpleChildrenTypes,omitempty"`
}

// GameplayGearComplexChild defines an item created inside a selected parent item or cargo loadout, with optional nested simple cargo.
type GameplayGearComplexChild struct {
	// ItemType is the game class name of the nested item to create.
	ItemType string `json:"itemType,omitempty" yaml:"itemType,omitempty"`
	// Attributes supplies the health and quantity ranges rolled for this item.
	Attributes *GameplayGearItemAttributes `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	// QuickBarSlot assigns the item to a quickbar slot; -1 leaves it unassigned.
	QuickBarSlot *int `json:"quickBarSlot,omitempty" yaml:"quickBarSlot,omitempty"`
	// SimpleChildrenUseDefaultAttributes makes simple child items inherit the attributes from this item instead of using configuration defaults.
	SimpleChildrenUseDefaultAttributes *bool `json:"simpleChildrenUseDefaultAttributes,omitempty" yaml:"simpleChildrenUseDefaultAttributes,omitempty"`
	// SimpleChildrenTypes are game class names of items created inside this item without per-child settings.
	SimpleChildrenTypes []string `json:"simpleChildrenTypes,omitempty" yaml:"simpleChildrenTypes,omitempty"`
}

// GameplayGearItemAttributes defines the inclusive numeric ranges from which item health and quantity are selected during spawning.
type GameplayGearItemAttributes struct {
	// HealthMin is the lower bound for the spawned item's health value.
	HealthMin *float64 `json:"healthMin,omitempty" yaml:"healthMin,omitempty"`
	// HealthMax is the upper bound for the spawned item's health value.
	HealthMax *float64 `json:"healthMax,omitempty" yaml:"healthMax,omitempty"`
	// QuantityMin is the lower bound for the spawned item's quantity value.
	QuantityMin *float64 `json:"quantityMin,omitempty" yaml:"quantityMin,omitempty"`
	// QuantityMax is the upper bound for the spawned item's quantity value.
	QuantityMax *float64 `json:"quantityMax,omitempty" yaml:"quantityMax,omitempty"`
}
