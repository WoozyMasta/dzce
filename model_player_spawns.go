// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// PlayerSpawnProfile defines spawn selection and candidate generation for one player category:
// fresh characters, same-map server hops, or map travel.
type PlayerSpawnProfile struct {
	// SpawnParams scores candidate points by their distance from players, infected, buildings, and trigger volumes at spawn time.
	SpawnParams *PlayerSpawnParams `json:"spawn_params,omitempty" yaml:"spawn_params,omitempty" xml:"spawn_params,omitempty"`
	// GeneratorParams defines how candidate points are generated around the configured centers and which terrain or geometry is acceptable.
	GeneratorParams *PlayerGeneratorParams `json:"generator_params,omitempty" yaml:"generator_params,omitempty" xml:"generator_params,omitempty"`
	// GroupParams enables rotating groups of spawn centers and sets their shared lifetime and reuse counter.
	GroupParams *PlayerGroupParams `json:"group_params,omitempty" yaml:"group_params,omitempty" xml:"group_params,omitempty"`
	// GeneratorPosBubbles lists the map positions around which the generator builds candidate spawn grids.
	// Positions can be grouped for group-based selection.
	GeneratorPosBubbles *PlayerGeneratorBubbles `json:"generator_posbubbles,omitempty" yaml:"generator_posbubbles,omitempty" xml:"generator_posbubbles,omitempty"`
}

// PlayerSpawnParams provides runtime distance bands used to rate candidate points.
// Distances below a minimum invalidate a point; values between the minimum and maximum are more favorable than values beyond the maximum.
type PlayerSpawnParams struct {
	// MinDistInfected is the minimum allowed distance from infected, in meters.
	MinDistInfected string `json:"min_dist_infected,omitempty" yaml:"min_dist_infected,omitempty" xml:"min_dist_infected,omitempty"`
	// MaxDistInfected is the distance from infected beyond which extra distance no longer improves the candidate's rating, in meters.
	MaxDistInfected string `json:"max_dist_infected,omitempty" yaml:"max_dist_infected,omitempty" xml:"max_dist_infected,omitempty"`
	// MinDistPlayer is the minimum allowed distance from other players, in meters.
	MinDistPlayer string `json:"min_dist_player,omitempty" yaml:"min_dist_player,omitempty" xml:"min_dist_player,omitempty"`
	// MaxDistPlayer is the distance from other players beyond which extra distance no longer improves the candidate's rating, in meters.
	MaxDistPlayer string `json:"max_dist_player,omitempty" yaml:"max_dist_player,omitempty" xml:"max_dist_player,omitempty"`
	// MinDistStatic is the minimum allowed distance from buildings and other static geometry, in meters.
	MinDistStatic string `json:"min_dist_static,omitempty" yaml:"min_dist_static,omitempty" xml:"min_dist_static,omitempty"`
	// MaxDistStatic is the upper end of the favorable distance band from static geometry, in meters.
	MaxDistStatic string `json:"max_dist_static,omitempty" yaml:"max_dist_static,omitempty" xml:"max_dist_static,omitempty"`
	// MinDistTrigger is the minimum allowed distance from trigger volumes, in meters.
	MinDistTrigger string `json:"min_dist_trigger,omitempty" yaml:"min_dist_trigger,omitempty" xml:"min_dist_trigger,omitempty"`
	// MaxDistTrigger is the upper end of the favorable distance band from trigger volumes, in meters.
	MaxDistTrigger string `json:"max_dist_trigger,omitempty" yaml:"max_dist_trigger,omitempty" xml:"max_dist_trigger,omitempty"`
}

// PlayerGeneratorParams defines the rectangular sampling grid
// around each generator position and filters candidate points by water, nearby geometry, and terrain slope.
type PlayerGeneratorParams struct {
	// AllowInWater permits candidate points in water when set to a true value.
	AllowInWater string `json:"allow_in_water,omitempty" yaml:"allow_in_water,omitempty" xml:"allow_in_water,omitempty"`
	// GridDensity sets the density of candidate samples across the generated rectangle.
	GridDensity string `json:"grid_density,omitempty" yaml:"grid_density,omitempty" xml:"grid_density,omitempty"`
	// GridWidth sets the generated rectangle's total width, in meters.
	GridWidth string `json:"grid_width,omitempty" yaml:"grid_width,omitempty" xml:"grid_width,omitempty"`
	// GridHeight sets the generated rectangle's total height, in meters.
	GridHeight string `json:"grid_height,omitempty" yaml:"grid_height,omitempty" xml:"grid_height,omitempty"`
	// MinDistStatic is the minimum distance from buildings required for a candidate point to be valid, in meters.
	MinDistStatic string `json:"min_dist_static,omitempty" yaml:"min_dist_static,omitempty" xml:"min_dist_static,omitempty"`
	// MaxDistStatic is the upper distance from buildings used to rate a valid candidate point, in meters.
	MaxDistStatic string `json:"max_dist_static,omitempty" yaml:"max_dist_static,omitempty" xml:"max_dist_static,omitempty"`
	// MinSteepness is the minimum terrain slope accepted for a candidate.
	MinSteepness string `json:"min_steepness,omitempty" yaml:"min_steepness,omitempty" xml:"min_steepness,omitempty"`
	// MaxSteepness is the maximum terrain slope accepted for a candidate.
	MaxSteepness string `json:"max_steepness,omitempty" yaml:"max_steepness,omitempty" xml:"max_steepness,omitempty"`
}

// PlayerGroupParams controls rotation through named groups of spawn centers.
type PlayerGroupParams struct {
	// EnableGroups enables group-based point selection. When disabled, points are treated as a regular candidate list.
	EnableGroups string `json:"enablegroups,omitempty" yaml:"enablegroups,omitempty" xml:"enablegroups,omitempty"`
	// GroupsAsRegular allows points inside groups to be used as regular spawn candidates instead of being selected only through group rotation.
	GroupsAsRegular string `json:"groups_as_regular,omitempty" yaml:"groups_as_regular,omitempty" xml:"groups_as_regular,omitempty"`
	// Lifetime is how long a selected group remains active before CE rotates to another group. A value of -1 disables the timer.
	Lifetime string `json:"lifetime,omitempty" yaml:"lifetime,omitempty" xml:"lifetime,omitempty"`
	// Counter is how many player spawns in the active group reset its lifetime. A value of -1 disables the counter.
	Counter string `json:"counter,omitempty" yaml:"counter,omitempty" xml:"counter,omitempty"`
}

// PlayerGeneratorBubbles contains the map positions used as spawn-grid centers, optionally organized into named groups.
type PlayerGeneratorBubbles struct {
	// Groups are named sets of centers used by group-based spawn selection.
	Groups []PlayerSpawnGroup `json:"group,omitempty" yaml:"group,omitempty" xml:"group,omitempty"`
}

// PlayerSpawnGroup defines one named spawn region and the centers used to generate candidate points in that region.
type PlayerSpawnGroup struct {
	// Name identifies the region for selection, configuration, and diagnostics.
	Name string `json:"name" yaml:"name" xml:"name,attr"`
	// Lifetime overrides the shared group lifetime, in seconds. A value of -1 disables the group timer.
	Lifetime string `json:"lifetime,omitempty" yaml:"lifetime,omitempty" xml:"lifetime,attr,omitempty"`
	// Pos are centers around which the generator creates candidate spawn grids.
	Pos []PlayerSpawnPos `json:"pos,omitempty" yaml:"pos,omitempty" xml:"pos,omitempty"`
}

// PlayerSpawnPos stores one horizontal map coordinate used as a spawn-grid center.
type PlayerSpawnPos struct {
	// X is the map's horizontal X coordinate.
	X string `json:"x,omitempty" yaml:"x,omitempty" xml:"x,attr,omitempty"`
	// Z is the map's horizontal Z coordinate.
	Z string `json:"z,omitempty" yaml:"z,omitempty" xml:"z,attr,omitempty"`
}
