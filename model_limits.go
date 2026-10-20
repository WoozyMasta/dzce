// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// LimitsCategoryList declares the category names that item types may use in their Central Economy filters.
type LimitsCategoryList struct {
	// Categories are category identifiers accepted by the economy's item definitions.
	Categories []NamedRef `json:"category,omitempty" yaml:"category,omitempty" xml:"category,omitempty"`
}

// LimitsTagList declares the tag names that item types may use in their Central Economy filters.
type LimitsTagList struct {
	// Tags are tag identifiers accepted by the economy's item definitions.
	Tags []NamedRef `json:"tag,omitempty" yaml:"tag,omitempty" xml:"tag,omitempty"`
}

// LimitsUsageList declares location-purpose flags used to restrict where an item type can spawn.
type LimitsUsageList struct {
	// Usages are identifiers assigned to map areas and referenced by item spawn filters.
	Usages []NamedRef `json:"usage,omitempty" yaml:"usage,omitempty" xml:"usage,omitempty"`
}

// LimitsValueList declares loot-tier flags used to restrict item spawning by map region or economy value.
type LimitsValueList struct {
	// Values are tier or value identifiers referenced by item spawn filters.
	Values []NamedRef `json:"value,omitempty" yaml:"value,omitempty" xml:"value,omitempty"`
}

// LimitsUserBindingList defines shorthand names that expand to several usage or value flags in Central Economy item definitions.
type LimitsUserBindingList struct {
	// Users are the named flag groups available to item definitions.
	Users []LimitsUserBinding `json:"user,omitempty" yaml:"user,omitempty" xml:"user,omitempty"`
}

// LimitsUserBinding maps one shorthand name to the usage or value flags it represents.
type LimitsUserBinding struct {
	// Name is the shorthand identifier referenced from an item definition.
	Name string `json:"name" yaml:"name" xml:"name,attr"`
	// Usages are the location-purpose flags included in this shorthand group.
	Usages []NamedRef `json:"usage,omitempty" yaml:"usage,omitempty" xml:"usage,omitempty"`
	// Values are the loot-tier flags included in this shorthand group.
	Values []NamedRef `json:"value,omitempty" yaml:"value,omitempty" xml:"value,omitempty"`
}
