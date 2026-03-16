// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// UndergroundTriggersFile is the root of `cfgundergroundtriggers.json`.
type UndergroundTriggersFile struct {
	// Triggers is underground trigger list for map.
	Triggers []UndergroundTrigger `json:"Triggers,omitempty" yaml:"Triggers,omitempty"`
}

// UndergroundTrigger stores one underground trigger area definition.
type UndergroundTrigger struct {
	// EyeAccommodation adjusts scene eye adaptation in trigger.
	EyeAccommodation *float64 `json:"EyeAccommodation,omitempty" yaml:"EyeAccommodation,omitempty"`
	// InterpolationSpeed controls eye-accommodation transition speed.
	// Wiki examples use this for non-breadcrumb (outer/inner) triggers.
	InterpolationSpeed *float64 `json:"InterpolationSpeed,omitempty" yaml:"InterpolationSpeed,omitempty"`
	// UseLinePointFade enables line-point fade transitions.
	UseLinePointFade *int `json:"UseLinePointFade,omitempty" yaml:"UseLinePointFade,omitempty"`
	// AmbientSoundType chooses ambient sound preset.
	AmbientSoundType *string `json:"AmbientSoundType,omitempty" yaml:"AmbientSoundType,omitempty"`
	// Position is trigger center in world coordinates.
	Position []float64 `json:"Position,omitempty" yaml:"Position,omitempty"`
	// Orientation is trigger orientation in Euler angles.
	Orientation []float64 `json:"Orientation,omitempty" yaml:"Orientation,omitempty"`
	// Size is box extents for trigger volume.
	Size []float64 `json:"Size,omitempty" yaml:"Size,omitempty"`
	// Breadcrumbs defines transitional points for gradual darkening/brightening.
	// Presence of breadcrumbs marks this trigger as transitional in wiki docs.
	Breadcrumbs []UndergroundBreadcrumb `json:"Breadcrumbs,omitempty" yaml:"Breadcrumbs,omitempty"`
}

// UndergroundBreadcrumb stores one underground trigger breadcrumb.
type UndergroundBreadcrumb struct {
	// EyeAccommodation overrides adaptation at breadcrumb point.
	EyeAccommodation *float64 `json:"EyeAccommodation,omitempty" yaml:"EyeAccommodation,omitempty"`
	// UseRaycast enables line-of-sight check from player to breadcrumb.
	// When enabled, breadcrumb contributes only if raycast succeeds.
	UseRaycast *int `json:"UseRaycast,omitempty" yaml:"UseRaycast,omitempty"`
	// Radius is local influence radius.
	Radius *float64 `json:"Radius,omitempty" yaml:"Radius,omitempty"`
	// LightLerp is lighting interpolation factor.
	LightLerp *float64 `json:"LightLerp,omitempty" yaml:"LightLerp,omitempty"`
	// Position is breadcrumb position in world coordinates.
	Position []float64 `json:"Position,omitempty" yaml:"Position,omitempty"`
}
