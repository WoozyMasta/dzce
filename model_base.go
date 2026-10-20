// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"fmt"
	"strings"
)

const (
	// VariableTypeInt is integer `type="0"` in globals.xml.
	VariableTypeInt VariableType = 0
	// VariableTypeFloat is float `type="1"` in globals.xml.
	VariableTypeFloat VariableType = 1
	// VariableTypeString is string `type="2"` in globals.xml.
	VariableTypeString VariableType = 2
)

// VariableType is the numeric value stored in a globals.xml variable's type attribute:
// 0 for integer, 1 for float, and 2 for string values.
type VariableType int

// WeatherToggle represents the boolean-like reset and enable attributes in cfgweather.xml.
// Text decoding accepts 0/1, true/false, yes/no, and on/off; an empty value decodes as false.
type WeatherToggle uint8

// UnmarshalText decodes a trimmed, case-insensitive WeatherToggle value.
// Values other than the supported boolean forms return an error.
func (toggle *WeatherToggle) UnmarshalText(text []byte) error {
	value := strings.ToLower(strings.TrimSpace(string(text)))

	switch value {
	case "1", "true", "yes", "on":
		*toggle = 1

	case "0", "false", "no", "off", "":
		*toggle = 0

	default:
		return fmt.Errorf("unsupported weather toggle %q", value)
	}

	return nil
}

// MarshalText encodes zero as 0 and every non-zero value as 1.
func (toggle WeatherToggle) MarshalText() ([]byte, error) {
	if toggle == 0 {
		return []byte("0"), nil
	}

	return []byte("1"), nil
}

// EconomyCoreCEFile registers one mission file in an economycore.xml <ce> block.
// The parent block supplies the folder;
// this entry supplies the file name and the CE document type to load.
type EconomyCoreCEFile struct {
	// Name is the file name relative to the folder on the parent <ce> element.
	Name string `json:"name" yaml:"name" xml:"name,attr"`
	// Type selects how Central Economy interprets the registered file.
	// The game documents types, spawnabletypes, globals, economy, events, and messages;
	// dzce also accepts economycore when traversing nested includes.
	Type string `json:"type" yaml:"type" xml:"type,attr"`
}

// EnvironmentTerritoryFile selects territory-zone data for a territory profile.
// A file entry can identify the data directly by path
// or refer to a reusable file declaration by its usable identifier.
type EnvironmentTerritoryFile struct {
	// Path is the direct path to a territory data file.
	Path string `json:"path,omitempty" yaml:"path,omitempty" xml:"path,attr,omitempty"`
	// Usable references a reusable territory file declaration by its identifier.
	Usable string `json:"usable,omitempty" yaml:"usable,omitempty" xml:"usable,attr,omitempty"`
}

// NamedRef stores the name attribute used by CE reference elements
// such as usage, category, tag, value, type, and group.
// The parent element determines what the name refers to.
type NamedRef struct {
	// Name is the identifier referenced by the containing CE element.
	Name string `json:"name" yaml:"name" xml:"name,attr"`
}

// EmptyElement represents a presence-only XML flag such as <unique/> or <hoarder/>.
// A non-nil value means the flag element is present; it has no attributes or text content.
type EmptyElement struct{}
