// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// GameplayGearPresetsFile is one gear preset JSON payload loaded via
// `cfggameplay.json -> PlayerData.spawnGearPresetFiles`.
//
// The root JSON value is an array of presets.
type GameplayGearPresetsFile []GameplayGearPreset

// GameplayGearPreset stores one player spawn gear preset.
type GameplayGearPreset struct {
	// Name is preset identifier.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// SpawnWeight is weighted random selection priority.
	SpawnWeight *int `json:"spawnWeight,omitempty" yaml:"spawnWeight,omitempty"`
	// CharacterTypes is list of allowed survivor class names.
	CharacterTypes []string `json:"characterTypes,omitempty" yaml:"characterTypes,omitempty"`
	// AttachmentSlotItemSets configures spawn variants per character slot.
	AttachmentSlotItemSets []GameplayGearAttachmentSlotSet `json:"attachmentSlotItemSets,omitempty" yaml:"attachmentSlotItemSets,omitempty"`
	// DiscreteUnsortedItemSets configures random cargo additions.
	DiscreteUnsortedItemSets []GameplayGearUnsortedSet `json:"discreteUnsortedItemSets,omitempty" yaml:"discreteUnsortedItemSets,omitempty"`
}

// GameplayGearAttachmentSlotSet configures one attachment slot.
type GameplayGearAttachmentSlotSet struct {
	// SlotName is CfgSlots slot id (`Body`, `Back`, `shoulderL`, and so on).
	SlotName string `json:"slotName,omitempty" yaml:"slotName,omitempty"`
	// DiscreteItemSets is weighted item variant list for this slot.
	DiscreteItemSets []GameplayGearDiscreteItemSet `json:"discreteItemSets,omitempty" yaml:"discreteItemSets,omitempty"`
}

// GameplayGearDiscreteItemSet is one weighted item variant entry.
type GameplayGearDiscreteItemSet struct {
	// SpawnWeight is weighted random selection priority.
	SpawnWeight *int `json:"spawnWeight,omitempty" yaml:"spawnWeight,omitempty"`
	// Attributes defines item health and quantity ranges.
	Attributes *GameplayGearItemAttributes `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	// QuickBarSlot sets assigned quickbar slot (`-1` means no assignment).
	QuickBarSlot *int `json:"quickBarSlot,omitempty" yaml:"quickBarSlot,omitempty"`
	// SimpleChildrenUseDefaultAttributes toggles default-attribute inheritance.
	SimpleChildrenUseDefaultAttributes *bool `json:"simpleChildrenUseDefaultAttributes,omitempty" yaml:"simpleChildrenUseDefaultAttributes,omitempty"`
	// ItemType is spawned item class name.
	ItemType string `json:"itemType,omitempty" yaml:"itemType,omitempty"`
	// ComplexChildrenTypes are nested children with their own settings.
	ComplexChildrenTypes []GameplayGearComplexChild `json:"complexChildrenTypes,omitempty" yaml:"complexChildrenTypes,omitempty"`
	// SimpleChildrenTypes are simple nested class names.
	SimpleChildrenTypes []string `json:"simpleChildrenTypes,omitempty" yaml:"simpleChildrenTypes,omitempty"`
}

// GameplayGearUnsortedSet configures weighted random cargo additions.
type GameplayGearUnsortedSet struct {
	// SpawnWeight is weighted random selection priority.
	SpawnWeight *int `json:"spawnWeight,omitempty" yaml:"spawnWeight,omitempty"`
	// Attributes defines default health/quantity ranges.
	Attributes *GameplayGearItemAttributes `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	// SimpleChildrenUseDefaultAttributes toggles default-attribute inheritance.
	SimpleChildrenUseDefaultAttributes *bool `json:"simpleChildrenUseDefaultAttributes,omitempty" yaml:"simpleChildrenUseDefaultAttributes,omitempty"`
	// Name is set identifier.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// ComplexChildrenTypes are nested children with their own settings.
	ComplexChildrenTypes []GameplayGearComplexChild `json:"complexChildrenTypes,omitempty" yaml:"complexChildrenTypes,omitempty"`
	// SimpleChildrenTypes are simple nested class names.
	SimpleChildrenTypes []string `json:"simpleChildrenTypes,omitempty" yaml:"simpleChildrenTypes,omitempty"`
}

// GameplayGearComplexChild is one nested child definition.
type GameplayGearComplexChild struct {
	// ItemType is spawned item class name.
	ItemType string `json:"itemType,omitempty" yaml:"itemType,omitempty"`
	// Attributes defines item health and quantity ranges.
	Attributes *GameplayGearItemAttributes `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	// QuickBarSlot sets assigned quickbar slot (`-1` means no assignment).
	QuickBarSlot *int `json:"quickBarSlot,omitempty" yaml:"quickBarSlot,omitempty"`
	// SimpleChildrenUseDefaultAttributes toggles default-attribute inheritance.
	SimpleChildrenUseDefaultAttributes *bool `json:"simpleChildrenUseDefaultAttributes,omitempty" yaml:"simpleChildrenUseDefaultAttributes,omitempty"`
	// SimpleChildrenTypes are simple nested class names.
	SimpleChildrenTypes []string `json:"simpleChildrenTypes,omitempty" yaml:"simpleChildrenTypes,omitempty"`
}

// GameplayGearItemAttributes stores common min/max item attributes.
type GameplayGearItemAttributes struct {
	// HealthMin is minimum health multiplier.
	HealthMin *float64 `json:"healthMin,omitempty" yaml:"healthMin,omitempty"`
	// HealthMax is maximum health multiplier.
	HealthMax *float64 `json:"healthMax,omitempty" yaml:"healthMax,omitempty"`
	// QuantityMin is minimum quantity multiplier.
	QuantityMin *float64 `json:"quantityMin,omitempty" yaml:"quantityMin,omitempty"`
	// QuantityMax is maximum quantity multiplier.
	QuantityMax *float64 `json:"quantityMax,omitempty" yaml:"quantityMax,omitempty"`
}
