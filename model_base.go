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

// VariableType is a globals.xml variable value type id.
type VariableType int

// WeatherToggle stores cfgweather boolean-like attribute values.
// Wiki allows 0/1, true/false, yes/no for reset/enable attributes.
type WeatherToggle uint8

// UnmarshalText decodes 0/1/true/false/yes/no toggle formats.
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

// MarshalText encodes weather toggle as 0/1.
func (toggle WeatherToggle) MarshalText() ([]byte, error) {
	if toggle == 0 {
		return []byte("0"), nil
	}

	return []byte("1"), nil
}

// EconomyCoreCEFile is a file mapping in `<ce>`.
type EconomyCoreCEFile struct {
	// Name is CE file name inside folder.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
	// Type is CE include category from wiki set:
	// types, spawnabletypes, globals, economy, events, messages.
	// economycore is additionally accepted for recursive include traversal.
	Type string `xml:"type,attr" json:"type" yaml:"type"`
}

// EnvironmentTerritoryFile references file by path or usable id.
type EnvironmentTerritoryFile struct {
	// Path is direct path to territory file.
	Path string `xml:"path,attr,omitempty" json:"path,omitempty" yaml:"path,omitempty"`
	// Usable is symbolic alias for previously declared file.
	Usable string `xml:"usable,attr,omitempty" json:"usable,omitempty" yaml:"usable,omitempty"`
}

// NamedRef stores `<x name="..."/>` style value.
type NamedRef struct {
	// Name is limiter or group name in CE config.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
}

// EmptyElement marks enabled empty XML element flags.
type EmptyElement struct{}
