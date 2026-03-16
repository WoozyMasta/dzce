// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// EffectAreaFile is the root of `cfgeffectarea.json`.
type EffectAreaFile struct {
	// Areas is configured effect area definitions.
	Areas []EffectArea `json:"Areas,omitempty" yaml:"Areas,omitempty"`
	// SafePositions is list of safe teleport/spawn points.
	SafePositions [][]float64 `json:"SafePositions,omitempty" yaml:"SafePositions,omitempty"`
}

// EffectArea stores one effect area definition.
type EffectArea struct {
	// Data stores geometry and effect settings for area.
	Data *EffectAreaData `json:"Data,omitempty" yaml:"Data,omitempty"`
	// PlayerData stores optional player-side VFX settings.
	PlayerData *EffectAreaPlayerData `json:"PlayerData,omitempty" yaml:"PlayerData,omitempty"`
	// AreaName is unique area identifier.
	AreaName string `json:"AreaName,omitempty" yaml:"AreaName,omitempty"`
	// Type is script class implementing area behavior.
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
	// TriggerType is trigger class used by the area.
	TriggerType string `json:"TriggerType,omitempty" yaml:"TriggerType,omitempty"`
}

// EffectAreaData stores area geometry and effect settings.
type EffectAreaData struct {
	// Radius is main area radius.
	Radius *float64 `json:"Radius,omitempty" yaml:"Radius,omitempty"`
	// PosHeight is positive vertical half-height.
	PosHeight *float64 `json:"PosHeight,omitempty" yaml:"PosHeight,omitempty"`
	// NegHeight is negative vertical half-height.
	NegHeight *float64 `json:"NegHeight,omitempty" yaml:"NegHeight,omitempty"`
	// InnerPartDist is particle spacing/size parameter used by area fill logic.
	// On newer configs (1.28+), wiki describes it as `partSize`.
	InnerPartDist *float64 `json:"InnerPartDist,omitempty" yaml:"InnerPartDist,omitempty"`
	// OuterOffset extends visible particle radius beyond the area boundary.
	// On newer configs (1.28+), wiki describes it as `outwardsBleed`.
	OuterOffset *float64 `json:"OuterOffset,omitempty" yaml:"OuterOffset,omitempty"`
	// ParticleName is particle system resource path.
	ParticleName *string `json:"ParticleName,omitempty" yaml:"ParticleName,omitempty"`
	// EffectInterval is period between periodic effects.
	EffectInterval *float64 `json:"EffectInterval,omitempty" yaml:"EffectInterval,omitempty"`
	// EffectDuration is active duration of periodic effect.
	EffectDuration *float64 `json:"EffectDuration,omitempty" yaml:"EffectDuration,omitempty"`
	// Pos is area center position in world coordinates.
	Pos []float64 `json:"Pos,omitempty" yaml:"Pos,omitempty"`
}

// EffectAreaPlayerData stores player visual effect details.
type EffectAreaPlayerData struct {
	// AroundPartName is particle around-player effect path.
	AroundPartName string `json:"AroundPartName,omitempty" yaml:"AroundPartName,omitempty"`
	// TinyPartName is near-camera tiny particle effect path.
	TinyPartName string `json:"TinyPartName,omitempty" yaml:"TinyPartName,omitempty"`
	// PPERequesterType is post-processing profile class.
	PPERequesterType string `json:"PPERequesterType,omitempty" yaml:"PPERequesterType,omitempty"`
}
