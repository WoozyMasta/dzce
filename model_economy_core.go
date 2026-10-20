// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// EconomyCoreClasses lists entity-class roots that Central Economy tracks.
// A root includes its derived game classes, so this section determines which entities CE recognizes and how it handles them.
type EconomyCoreClasses struct {
	// RootClasses are class inheritance roots registered with CE.
	RootClasses []EconomyCoreRootClass `json:"rootclass,omitempty" yaml:"rootclass,omitempty" xml:"rootclass,omitempty"`
}

// EconomyCoreRootClass registers one game class hierarchy with Central Economy
// and identifies whether its instances are loot, characters, or cars.
type EconomyCoreRootClass struct {
	// Name is the root class name from the game's configuration inheritance.
	Name string `json:"name" yaml:"name" xml:"name,attr" jsonschema:"example=DefaultWeapon"`
	// ReportMemoryLOD controls warnings for missing Memory LODs on instances of this class (`yes` or `no`). The documented default is `yes`.
	ReportMemoryLOD string `json:"reportMemoryLOD,omitempty" yaml:"reportMemoryLOD,omitempty" xml:"reportMemoryLOD,attr,omitempty" jsonschema:"enum=yes,enum=no,default=yes,example=no"`
	// Act tells CE how to initialize this class: `character` for players, infected, and animals;
	// `car` for vehicles; or omitted/`none` for ordinary economy items.
	Act string `json:"act,omitempty" yaml:"act,omitempty" xml:"act,attr,omitempty" jsonschema:"enum=none,enum=character,enum=car,default=none,example=character"`
}

// EconomyCoreDefaults overrides built-in Central Economy default variables.
type EconomyCoreDefaults struct {
	// Defaults are named CE default values applied by the mission.
	Defaults []EconomyCoreDefault `json:"default,omitempty" yaml:"default,omitempty" xml:"default,omitempty"`
}

// EconomyCoreDefault assigns a value to one Central Economy default variable.
type EconomyCoreDefault struct {
	// Name is the CE variable name being overridden.
	Name string `json:"name" yaml:"name" xml:"name,attr" jsonschema:"example=dyn_radius"`
	// Value is the configured value encoded as text. Its expected format is determined by the named variable.
	Value string `json:"value" yaml:"value" xml:"value,attr" jsonschema:"example=30"`
}

// EconomyCoreCE registers a folder of mission-provided Central Economy files.
// The game uses these entries to append to or override selected vanilla CE data without copying the full source files into the mission.
type EconomyCoreCE struct {
	// Folder is the mission-relative directory containing the registered CE files.
	Folder string `json:"folder" yaml:"folder" xml:"folder,attr"`
	// Files lists the source file and CE data type to append or override.
	Files []EconomyCoreCEFile `json:"file,omitempty" yaml:"file,omitempty" xml:"file,omitempty"`
}
