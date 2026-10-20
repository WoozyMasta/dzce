// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "encoding/json"

// UndergroundTriggersFile defines the trigger volumes used to manage
// underground lighting transitions and ambient audio for the mission.
type UndergroundTriggersFile struct {
	// Triggers are the world-space volumes and transition points placed on the map.
	// A trigger with breadcrumbs performs a gradual transition;
	// without them, its eye-accommodation target determines whether it is an outer or inner area.
	Triggers []UndergroundTrigger `json:"Triggers,omitempty" yaml:"Triggers,omitempty"`
}

// UndergroundTrigger defines one box volume that controls underground
// visibility and, optionally, ambient sound or breadcrumb transitions.
type UndergroundTrigger struct {
	// Comment is a human-readable label for identifying the trigger in mission data and diagnostic tools.
	Comment string `json:"Comment,omitempty" yaml:"Comment,omitempty"`
	// CustomSpawn marks a trigger whose creation is associated with a parent map object rather than an independently placed trigger carrier.
	CustomSpawn json.RawMessage `json:"CustomSpawn,omitempty" yaml:"CustomSpawn,omitempty"`
	// ParentNetworkID is the network identifier of the parent entity used to associate a custom-spawn trigger with that object.
	ParentNetworkID []int64 `json:"ParentNetworkId,omitempty" yaml:"ParentNetworkId,omitempty"`
	// EyeAccommodation is the target eye-adaptation level inside the volume: 0 represents darkness and 1 represents normal visibility.
	EyeAccommodation *float64 `json:"EyeAccommodation,omitempty" yaml:"EyeAccommodation,omitempty"`
	// InterpolationSpeed controls how quickly eye accommodation approaches the trigger's target value,
	// especially for outer and inner triggers that do not use breadcrumb points.
	InterpolationSpeed *float64 `json:"InterpolationSpeed,omitempty" yaml:"InterpolationSpeed,omitempty"`
	// UseLinePointFade enables a linear fade along the configured breadcrumb path, which is useful for tunnels and corridors.
	UseLinePointFade *int `json:"UseLinePointFade,omitempty" yaml:"UseLinePointFade,omitempty"`
	// AmbientSoundType selects an engine-provided ambient sound controller.
	AmbientSoundType *string `json:"AmbientSoundType,omitempty" yaml:"AmbientSoundType,omitempty"`
	// AmbientSoundSet selects a sound set for ambient audio played in the trigger volume.
	AmbientSoundSet *string `json:"AmbientSoundSet,omitempty" yaml:"AmbientSoundSet,omitempty"`
	// Position is the center of the trigger volume in world X/Y/Z coordinates.
	Position []float64 `json:"Position,omitempty" yaml:"Position,omitempty"`
	// Orientation is the trigger rotation as yaw, pitch, and roll in degrees.
	Orientation []float64 `json:"Orientation,omitempty" yaml:"Orientation,omitempty"`
	// Size is the trigger volume's X/Y/Z dimensions in meters.
	Size []float64 `json:"Size,omitempty" yaml:"Size,omitempty"`
	// Breadcrumbs are points inside a transitional volume that shape the gradual eye-accommodation change as a player moves through it.
	Breadcrumbs []UndergroundBreadcrumb `json:"Breadcrumbs,omitempty" yaml:"Breadcrumbs,omitempty"`
}

// UndergroundBreadcrumb defines a point that contributes to the visibility transition
// inside its parent underground trigger.
type UndergroundBreadcrumb struct {
	// Comment is an optional label for identifying the point in mission data.
	Comment string `json:"Comment,omitempty" yaml:"Comment,omitempty"`
	// ExternalValueController makes the point's contribution depend on a script-controlled value, such as the state of a door.
	ExternalValueController *UndergroundExternalValueController `json:"ExternalValueController,omitempty" yaml:"ExternalValueController,omitempty"`
	// EyeAccommodation is the visibility target contributed by this point.
	EyeAccommodation *float64 `json:"EyeAccommodation,omitempty" yaml:"EyeAccommodation,omitempty"`
	// UseRaycast requires a successful line-of-sight check from the player to this point before it contributes to the transition.
	UseRaycast *int `json:"UseRaycast,omitempty" yaml:"UseRaycast,omitempty"`
	// Radius is the point's influence radius in meters. A value of -1 uses the engine's default radius.
	Radius *float64 `json:"Radius,omitempty" yaml:"Radius,omitempty"`
	// LightLerp controls world-light interpolation for line-point fade mode.
	LightLerp *float64 `json:"LightLerp,omitempty" yaml:"LightLerp,omitempty"`
	// Position is the point's world X/Y/Z coordinate.
	Position []float64 `json:"Position,omitempty" yaml:"Position,omitempty"`
}

// UndergroundExternalValueController selects a script-side condition that controls whether a breadcrumb contributes to the transition.
type UndergroundExternalValueController struct {
	// Type identifies the controller implementation that reads the external condition.
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
	// Params are arguments interpreted by the selected controller, such as the name of a door whose state controls the point.
	Params []string `json:"Params,omitempty" yaml:"Params,omitempty"`
}
