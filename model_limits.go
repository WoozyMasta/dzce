// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// LimitsCategoryList stores category entries.
type LimitsCategoryList struct {
	// Categories holds `<category name="..."/>` entries.
	Categories []NamedRef `xml:"category,omitempty" json:"category,omitempty" yaml:"category,omitempty"`
}

// LimitsTagList stores tag entries.
type LimitsTagList struct {
	// Tags holds `<tag name="..."/>` entries.
	Tags []NamedRef `xml:"tag,omitempty" json:"tag,omitempty" yaml:"tag,omitempty"`
}

// LimitsUsageList stores usage entries.
type LimitsUsageList struct {
	// Usages holds `<usage name="..."/>` entries.
	Usages []NamedRef `xml:"usage,omitempty" json:"usage,omitempty" yaml:"usage,omitempty"`
}

// LimitsValueList stores value entries.
type LimitsValueList struct {
	// Values holds `<value name="..."/>` entries.
	Values []NamedRef `xml:"value,omitempty" json:"value,omitempty" yaml:"value,omitempty"`
}

// LimitsUserBindingList stores user aliases for usage/value flags.
type LimitsUserBindingList struct {
	// Users is list of named aliases.
	Users []LimitsUserBinding `xml:"user,omitempty" json:"user,omitempty" yaml:"user,omitempty"`
}

// LimitsUserBinding stores one alias with usage/value entries.
type LimitsUserBinding struct {
	// Name is alias identifier.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
	// Usages is usage flags included in alias.
	Usages []NamedRef `xml:"usage,omitempty" json:"usage,omitempty" yaml:"usage,omitempty"`
	// Values is value flags included in alias.
	Values []NamedRef `xml:"value,omitempty" json:"value,omitempty" yaml:"value,omitempty"`
}
