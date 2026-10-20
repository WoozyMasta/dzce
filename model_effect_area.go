// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// EffectAreaFile configures static world effect areas and the global safe destinations used by contaminated-area recovery logic.
// Static areas are created when the mission starts; dynamic contaminated areas are managed by Central Economy instead.
type EffectAreaFile struct {
	// Areas define the effect volumes, their script behavior, and optional visual effects shown to players.
	Areas []EffectArea `json:"Areas,omitempty" yaml:"Areas,omitempty"`
	// SafePositions are global X/Z destination candidates used when a player connects inside a contaminated area.
	// The game selects the nearest point in 2D and derives its Y coordinate from the terrain.
	SafePositions [][]float64 `json:"SafePositions,omitempty" yaml:"SafePositions,omitempty"`
}

// EffectArea binds an area name and geometry to the script classes that implement its trigger and gameplay effects.
type EffectArea struct {
	// Data defines the area's position, dimensions, and effect-specific settings.
	Data *EffectAreaData `json:"Data,omitempty" yaml:"Data,omitempty"`
	// PlayerData defines optional particles and post-processing applied to players affected by this area.
	PlayerData *EffectAreaPlayerData `json:"PlayerData,omitempty" yaml:"PlayerData,omitempty"`
	// AreaName identifies the area in mission configuration and diagnostics.
	AreaName string `json:"AreaName,omitempty" yaml:"AreaName,omitempty"`
	// Type selects the script class that implements the area's gameplay behavior, such as a static contaminated area.
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
	// TriggerType selects the script trigger class that detects players inside the area.
	TriggerType string `json:"TriggerType,omitempty" yaml:"TriggerType,omitempty"`
}

// EffectAreaData defines the area's cylindrical bounds and the particle or timed-effect parameters consumed by its script class.
// Particle emitters are limited to 1,000 per area; large emitter counts can affect performance.
type EffectAreaData struct {
	// Radius is the horizontal radius of the area in meters.
	// The game skips the area when this value is zero or negative.
	Radius *float64 `json:"Radius,omitempty" yaml:"Radius,omitempty"`
	// PosHeight is the distance in meters that the area extends above Pos.
	PosHeight *float64 `json:"PosHeight,omitempty" yaml:"PosHeight,omitempty"`
	// NegHeight is the distance in meters that the area extends below Pos.
	NegHeight *float64 `json:"NegHeight,omitempty" yaml:"NegHeight,omitempty"`
	// InnerPartDist controls particle size when filling the area with emitters.
	// The 1.28 configuration describes this setting as partSize.
	InnerPartDist *float64 `json:"InnerPartDist,omitempty" yaml:"InnerPartDist,omitempty"`
	// OuterOffset extends visible particles beyond the configured area radius, in meters.
	OuterOffset *float64 `json:"OuterOffset,omitempty" yaml:"OuterOffset,omitempty"`
	// ParticleName is the resource path for particles emitted by the area.
	ParticleName *string `json:"ParticleName,omitempty" yaml:"ParticleName,omitempty"`
	// EffectInterval is the delay, in seconds, between recurring effects for area types that implement timed effects.
	EffectInterval *float64 `json:"EffectInterval,omitempty" yaml:"EffectInterval,omitempty"`
	// EffectDuration is how long, in seconds, a recurring effect remains active for area types that implement timed effects.
	EffectDuration *float64 `json:"EffectDuration,omitempty" yaml:"EffectDuration,omitempty"`
	// Pos is the area's center or anchor in world X/Y/Z coordinates. It does not need to be placed on the terrain surface.
	Pos []float64 `json:"Pos,omitempty" yaml:"Pos,omitempty"`
}

// EffectAreaPlayerData selects the local visual effects applied to a player while the area behavior considers them affected.
type EffectAreaPlayerData struct {
	// AroundPartName is the particle resource rendered around the player.
	AroundPartName string `json:"AroundPartName,omitempty" yaml:"AroundPartName,omitempty"`
	// TinyPartName is the particle resource rendered close to the camera.
	TinyPartName string `json:"TinyPartName,omitempty" yaml:"TinyPartName,omitempty"`
	// PPERequesterType is the post-processing requester class used for the player's screen effect.
	PPERequesterType string `json:"PPERequesterType,omitempty" yaml:"PPERequesterType,omitempty"`
}
