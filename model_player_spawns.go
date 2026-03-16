// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// PlayerSpawnProfile stores one spawn profile block.
type PlayerSpawnProfile struct {
	// SpawnParams controls dynamic distance constraints.
	SpawnParams *PlayerSpawnParams `xml:"spawn_params,omitempty" json:"spawn_params,omitempty" yaml:"spawn_params,omitempty"`
	// GeneratorParams controls candidate grid generation.
	GeneratorParams *PlayerGeneratorParams `xml:"generator_params,omitempty" json:"generator_params,omitempty" yaml:"generator_params,omitempty"`
	// GroupParams controls group-based spawn reuse.
	GroupParams *PlayerGroupParams `xml:"group_params,omitempty" json:"group_params,omitempty" yaml:"group_params,omitempty"`
	// GeneratorPosBubbles stores manual generator centers.
	GeneratorPosBubbles *PlayerGeneratorBubbles `xml:"generator_posbubbles,omitempty" json:"generator_posbubbles,omitempty" yaml:"generator_posbubbles,omitempty"`
}

// PlayerSpawnParams defines dynamic spawn scoring distances.
type PlayerSpawnParams struct {
	// MinDistInfected is min distance from infected.
	MinDistInfected string `xml:"min_dist_infected,omitempty" json:"min_dist_infected,omitempty" yaml:"min_dist_infected,omitempty"`
	// MaxDistInfected is max distance from infected.
	MaxDistInfected string `xml:"max_dist_infected,omitempty" json:"max_dist_infected,omitempty" yaml:"max_dist_infected,omitempty"`
	// MinDistPlayer is min distance from other players.
	MinDistPlayer string `xml:"min_dist_player,omitempty" json:"min_dist_player,omitempty" yaml:"min_dist_player,omitempty"`
	// MaxDistPlayer is max distance from other players.
	MaxDistPlayer string `xml:"max_dist_player,omitempty" json:"max_dist_player,omitempty" yaml:"max_dist_player,omitempty"`
	// MinDistStatic is min distance from static obstacles.
	MinDistStatic string `xml:"min_dist_static,omitempty" json:"min_dist_static,omitempty" yaml:"min_dist_static,omitempty"`
	// MaxDistStatic is max distance from static obstacles.
	MaxDistStatic string `xml:"max_dist_static,omitempty" json:"max_dist_static,omitempty" yaml:"max_dist_static,omitempty"`
	// MinDistTrigger is min distance from trigger volumes.
	MinDistTrigger string `xml:"min_dist_trigger,omitempty" json:"min_dist_trigger,omitempty" yaml:"min_dist_trigger,omitempty"`
	// MaxDistTrigger is max distance from trigger volumes.
	MaxDistTrigger string `xml:"max_dist_trigger,omitempty" json:"max_dist_trigger,omitempty" yaml:"max_dist_trigger,omitempty"`
}

// PlayerGeneratorParams defines generator grid and slope parameters.
type PlayerGeneratorParams struct {
	// GridDensity is candidate sampling density.
	GridDensity string `xml:"grid_density,omitempty" json:"grid_density,omitempty" yaml:"grid_density,omitempty"`
	// GridWidth is sampling grid width in meters.
	GridWidth string `xml:"grid_width,omitempty" json:"grid_width,omitempty" yaml:"grid_width,omitempty"`
	// GridHeight is sampling grid height in meters.
	GridHeight string `xml:"grid_height,omitempty" json:"grid_height,omitempty" yaml:"grid_height,omitempty"`
	// MinDistStatic is min distance from static geometry.
	MinDistStatic string `xml:"min_dist_static,omitempty" json:"min_dist_static,omitempty" yaml:"min_dist_static,omitempty"`
	// MaxDistStatic is max distance from static geometry.
	MaxDistStatic string `xml:"max_dist_static,omitempty" json:"max_dist_static,omitempty" yaml:"max_dist_static,omitempty"`
	// MinSteepness is minimum terrain steepness.
	MinSteepness string `xml:"min_steepness,omitempty" json:"min_steepness,omitempty" yaml:"min_steepness,omitempty"`
	// MaxSteepness is maximum terrain steepness.
	MaxSteepness string `xml:"max_steepness,omitempty" json:"max_steepness,omitempty" yaml:"max_steepness,omitempty"`
}

// PlayerGroupParams controls spawn group selection behavior.
type PlayerGroupParams struct {
	// EnableGroups toggles group-based spawn system.
	EnableGroups string `xml:"enablegroups,omitempty" json:"enablegroups,omitempty" yaml:"enablegroups,omitempty"`
	// GroupsAsRegular treats grouped points as normal candidates.
	GroupsAsRegular string `xml:"groups_as_regular,omitempty" json:"groups_as_regular,omitempty" yaml:"groups_as_regular,omitempty"`
	// Lifetime is lifetime for generated group entries.
	Lifetime string `xml:"lifetime,omitempty" json:"lifetime,omitempty" yaml:"lifetime,omitempty"`
	// Counter controls generated group reuse count.
	Counter string `xml:"counter,omitempty" json:"counter,omitempty" yaml:"counter,omitempty"`
}

// PlayerGeneratorBubbles stores grid center groups for spawn generation.
type PlayerGeneratorBubbles struct {
	// Groups is list of named bubble centers.
	Groups []PlayerSpawnGroup `xml:"group,omitempty" json:"group,omitempty" yaml:"group,omitempty"`
}

// PlayerSpawnGroup stores one named set of spawn positions.
type PlayerSpawnGroup struct {
	// Name is spawn group identifier.
	Name string `xml:"name,attr" json:"name" yaml:"name"`
	// Pos are positions that belong to the group.
	Pos []PlayerSpawnPos `xml:"pos,omitempty" json:"pos,omitempty" yaml:"pos,omitempty"`
}

// PlayerSpawnPos stores one x/z position in spawn profile groups.
type PlayerSpawnPos struct {
	// X is world X coordinate.
	X string `xml:"x,attr,omitempty" json:"x,omitempty" yaml:"x,omitempty"`
	// Z is world Z coordinate.
	Z string `xml:"z,attr,omitempty" json:"z,omitempty" yaml:"z,omitempty"`
}
