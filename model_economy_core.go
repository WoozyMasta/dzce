// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// EconomyCoreClasses holds root classes used by economy.
type EconomyCoreClasses struct {
	// RootClasses are class roots used by CE logic.
	RootClasses []EconomyCoreRootClass `xml:"rootclass,omitempty" json:"rootclass,omitempty" yaml:"rootclass,omitempty"`
}

// EconomyCoreRootClass is a single `<rootclass />` in economy core.
type EconomyCoreRootClass struct {
	// Name is root class name from config inheritance.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
	// ReportMemoryLOD enables/disables missing Memory LOD warnings (`yes`/`no`).
	// Wiki default is `yes`.
	ReportMemoryLOD string `xml:"reportMemoryLOD,attr,omitempty" json:"reportMemoryLOD,omitempty" yaml:"reportMemoryLOD,omitempty"`
	// Act selects CE actor class category: `none` (or omitted), `character`,
	// or `car`. Omitted/`none` is treated as regular loot entity.
	Act string `xml:"act,attr,omitempty" json:"act,omitempty" yaml:"act,omitempty"`
}

// EconomyCoreDefaults holds `<default />` values.
type EconomyCoreDefaults struct {
	// Defaults is the list of CE default key/value pairs.
	Defaults []EconomyCoreDefault `xml:"default,omitempty" json:"default,omitempty" yaml:"default,omitempty"`
}

// EconomyCoreDefault is a single `<default />` in economy core.
type EconomyCoreDefault struct {
	// Name is default variable key.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
	// Value is default variable value.
	Value string `xml:"value,attr" json:"value" yaml:"value"`
}

// EconomyCoreCE describes mission-level CE override folder.
type EconomyCoreCE struct {
	// Folder is mod/mission folder with CE overrides, resolved relative to
	// current economycore file location.
	Folder string `xml:"folder,attr" json:"folder" yaml:"folder"`
	// Files are include entries loaded from this override folder.
	Files []EconomyCoreCEFile `xml:"file,omitempty" json:"file,omitempty" yaml:"file,omitempty"`
}
