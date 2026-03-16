// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// GameplayFile is the root of `cfggameplay.json`.
type GameplayFile struct {
	// Version stores schema/config version marker.
	Version *int `json:"version,omitempty" yaml:"version,omitempty"`
	// GeneralData stores base gameplay toggles.
	GeneralData *GameplayGeneralData `json:"GeneralData,omitempty" yaml:"GeneralData,omitempty"`
	// PlayerData stores player movement and stamina settings.
	PlayerData *GameplayPlayerData `json:"PlayerData,omitempty" yaml:"PlayerData,omitempty"`
	// WorldsData stores world environment and object-spawner settings.
	WorldsData *GameplayWorldsData `json:"WorldsData,omitempty" yaml:"WorldsData,omitempty"`
	// BaseBuildingData stores hologram and construction restrictions.
	BaseBuildingData *GameplayBaseBuildingData `json:"BaseBuildingData,omitempty" yaml:"BaseBuildingData,omitempty"`
	// UIData stores HUD/UI-specific gameplay tuning.
	UIData *GameplayUIData `json:"UIData,omitempty" yaml:"UIData,omitempty"`
	// MapData stores map and navigation ownership behavior.
	MapData *GameplayMapData `json:"MapData,omitempty" yaml:"MapData,omitempty"`
	// VehicleData stores vehicle-specific modifiers.
	VehicleData *GameplayVehicleData `json:"VehicleData,omitempty" yaml:"VehicleData,omitempty"`
}

// GameplayGeneralData stores generic gameplay toggles.
type GameplayGeneralData struct {
	// DisableBaseDamage disables damage to base-building structures.
	DisableBaseDamage *bool `json:"disableBaseDamage,omitempty" yaml:"disableBaseDamage,omitempty"`
	// DisableContainerDamage disables damage to containers (tents/barrels).
	DisableContainerDamage *bool `json:"disableContainerDamage,omitempty" yaml:"disableContainerDamage,omitempty"`
	// DisableRespawnDialog disables respawn dialog UI.
	DisableRespawnDialog *bool `json:"disableRespawnDialog,omitempty" yaml:"disableRespawnDialog,omitempty"`
	// DisableRespawnInUnconsciousness disables respawn while unconscious.
	DisableRespawnInUnconsciousness *bool `json:"disableRespawnInUnconsciousness,omitempty" yaml:"disableRespawnInUnconsciousness,omitempty"`
}

// GameplayPlayerData stores player movement and stamina settings.
type GameplayPlayerData struct {
	// DisablePersonalLight disables personal night light.
	DisablePersonalLight *bool `json:"disablePersonalLight,omitempty" yaml:"disablePersonalLight,omitempty"`
	// StaminaData configures stamina consumption and limits.
	StaminaData *GameplayStaminaData `json:"StaminaData,omitempty" yaml:"StaminaData,omitempty"`
	// ShockHandlingData configures shock refill behavior.
	ShockHandlingData *GameplayShockHandlingData `json:"ShockHandlingData,omitempty" yaml:"ShockHandlingData,omitempty"`
	// MovementData configures inertia and turning values.
	MovementData *GameplayMovementData `json:"MovementData,omitempty" yaml:"MovementData,omitempty"`
	// DrowningData configures drowning depletion rates.
	DrowningData *GameplayDrowningData `json:"DrowningData,omitempty" yaml:"DrowningData,omitempty"`
	// WeaponObstructionData configures obstruction behavior modes.
	WeaponObstructionData *GameplayWeaponObstructionData `json:"WeaponObstructionData,omitempty" yaml:"WeaponObstructionData,omitempty"`
	// SpawnGearPresetFiles lists player gear preset JSON files loaded from
	// mission folder/subfolders. File names are arbitrary.
	SpawnGearPresetFiles []string `json:"spawnGearPresetFiles,omitempty" yaml:"spawnGearPresetFiles,omitempty"`
}

// GameplayStaminaData stores stamina settings.
type GameplayStaminaData struct {
	// SprintStaminaModifierErc modifies sprint drain while erect.
	SprintStaminaModifierErc *float64 `json:"sprintStaminaModifierErc,omitempty" yaml:"sprintStaminaModifierErc,omitempty"`
	// SprintStaminaModifierCro modifies sprint drain while crouched.
	SprintStaminaModifierCro *float64 `json:"sprintStaminaModifierCro,omitempty" yaml:"sprintStaminaModifierCro,omitempty"`
	// StaminaWeightLimitThreshold sets weight threshold before penalties.
	StaminaWeightLimitThreshold *float64 `json:"staminaWeightLimitThreshold,omitempty" yaml:"staminaWeightLimitThreshold,omitempty"`
	// StaminaMax sets max stamina points.
	StaminaMax *float64 `json:"staminaMax,omitempty" yaml:"staminaMax,omitempty"`
	// StaminaKgToStaminaPercentPenalty controls kg-to-penalty factor.
	StaminaKgToStaminaPercentPenalty *float64 `json:"staminaKgToStaminaPercentPenalty,omitempty" yaml:"staminaKgToStaminaPercentPenalty,omitempty"`
	// StaminaMinCap sets minimum stamina cap.
	StaminaMinCap *float64 `json:"staminaMinCap,omitempty" yaml:"staminaMinCap,omitempty"`
	// SprintSwimmingStaminaModifier modifies swimming sprint drain.
	SprintSwimmingStaminaModifier *float64 `json:"sprintSwimmingStaminaModifier,omitempty" yaml:"sprintSwimmingStaminaModifier,omitempty"`
	// SprintLadderStaminaModifier modifies ladder sprint drain.
	SprintLadderStaminaModifier *float64 `json:"sprintLadderStaminaModifier,omitempty" yaml:"sprintLadderStaminaModifier,omitempty"`
	// MeleeStaminaModifier modifies heavy melee stamina drain.
	MeleeStaminaModifier *float64 `json:"meleeStaminaModifier,omitempty" yaml:"meleeStaminaModifier,omitempty"`
	// ObstacleTraversalStaminaModifier modifies jump/vault drain.
	ObstacleTraversalStaminaModifier *float64 `json:"obstacleTraversalStaminaModifier,omitempty" yaml:"obstacleTraversalStaminaModifier,omitempty"`
	// HoldBreathStaminaModifier modifies hold-breath stamina drain.
	HoldBreathStaminaModifier *float64 `json:"holdBreathStaminaModifier,omitempty" yaml:"holdBreathStaminaModifier,omitempty"`
}

// GameplayShockHandlingData stores shock refill settings.
type GameplayShockHandlingData struct {
	// ShockRefillSpeedConscious is conscious refill speed per second.
	ShockRefillSpeedConscious *float64 `json:"shockRefillSpeedConscious,omitempty" yaml:"shockRefillSpeedConscious,omitempty"`
	// ShockRefillSpeedUnconscious is unconscious refill speed per second.
	ShockRefillSpeedUnconscious *float64 `json:"shockRefillSpeedUnconscious,omitempty" yaml:"shockRefillSpeedUnconscious,omitempty"`
	// AllowRefillSpeedModifier enables ammo-type shock modifiers.
	AllowRefillSpeedModifier *bool `json:"allowRefillSpeedModifier,omitempty" yaml:"allowRefillSpeedModifier,omitempty"`
}

// GameplayMovementData stores movement inertia settings.
type GameplayMovementData struct {
	// TimeToStrafeJog is blend time for jogging strafe movement.
	TimeToStrafeJog *float64 `json:"timeToStrafeJog,omitempty" yaml:"timeToStrafeJog,omitempty"`
	// RotationSpeedJog is yaw rotation speed while jogging.
	RotationSpeedJog *float64 `json:"rotationSpeedJog,omitempty" yaml:"rotationSpeedJog,omitempty"`
	// TimeToSprint is blend time from jog to sprint.
	TimeToSprint *float64 `json:"timeToSprint,omitempty" yaml:"timeToSprint,omitempty"`
	// TimeToStrafeSprint is blend time for sprint strafing.
	TimeToStrafeSprint *float64 `json:"timeToStrafeSprint,omitempty" yaml:"timeToStrafeSprint,omitempty"`
	// RotationSpeedSprint is yaw rotation speed while sprinting.
	RotationSpeedSprint *float64 `json:"rotationSpeedSprint,omitempty" yaml:"rotationSpeedSprint,omitempty"`
	// AllowStaminaAffectInertia enables stamina-inertia coupling.
	AllowStaminaAffectInertia *bool `json:"allowStaminaAffectInertia,omitempty" yaml:"allowStaminaAffectInertia,omitempty"`
}

// GameplayDrowningData stores drowning depletion rates.
type GameplayDrowningData struct {
	// StaminaDepletionSpeed sets stamina loss per second in drowning.
	StaminaDepletionSpeed *float64 `json:"staminaDepletionSpeed,omitempty" yaml:"staminaDepletionSpeed,omitempty"`
	// HealthDepletionSpeed sets health loss per second in drowning.
	HealthDepletionSpeed *float64 `json:"healthDepletionSpeed,omitempty" yaml:"healthDepletionSpeed,omitempty"`
	// ShockDepletionSpeed sets shock loss per second in drowning.
	ShockDepletionSpeed *float64 `json:"shockDepletionSpeed,omitempty" yaml:"shockDepletionSpeed,omitempty"`
}

// GameplayWeaponObstructionData stores obstruction mode settings.
type GameplayWeaponObstructionData struct {
	// StaticMode controls obstruction against static objects:
	// 0=off, 1=obstruct+lift, 2=always obstruct.
	StaticMode *int `json:"staticMode,omitempty" yaml:"staticMode,omitempty"`
	// DynamicMode controls obstruction against dynamic objects:
	// 0=off, 1=obstruct+lift, 2=always obstruct.
	DynamicMode *int `json:"dynamicMode,omitempty" yaml:"dynamicMode,omitempty"`
}

// GameplayWorldsData stores world-level gameplay settings.
type GameplayWorldsData struct {
	// LightingConfig selects night lighting mode (0=bright, 1=dark).
	LightingConfig *int `json:"lightingConfig,omitempty" yaml:"lightingConfig,omitempty"`
	// ObjectSpawnersArr lists object spawner config files.
	ObjectSpawnersArr []string `json:"objectSpawnersArr,omitempty" yaml:"objectSpawnersArr,omitempty"`
	// EnvironmentMinTemps stores monthly minimum temperatures (12 values).
	EnvironmentMinTemps []float64 `json:"environmentMinTemps,omitempty" yaml:"environmentMinTemps,omitempty"`
	// EnvironmentMaxTemps stores monthly maximum temperatures (12 values).
	EnvironmentMaxTemps []float64 `json:"environmentMaxTemps,omitempty" yaml:"environmentMaxTemps,omitempty"`
	// WetnessWeightModifiers stores wetness-to-weight multipliers.
	WetnessWeightModifiers []float64 `json:"wetnessWeightModifiers,omitempty" yaml:"wetnessWeightModifiers,omitempty"`
	// PlayerRestrictedAreaFiles lists player restricted area config files.
	PlayerRestrictedAreaFiles []string `json:"playerRestrictedAreaFiles,omitempty" yaml:"playerRestrictedAreaFiles,omitempty"`
}

// GameplayBaseBuildingData stores base building restriction settings.
type GameplayBaseBuildingData struct {
	// HologramData configures placement checks for holograms.
	HologramData *GameplayHologramData `json:"HologramData,omitempty" yaml:"HologramData,omitempty"`
	// ConstructionData configures checks for construction phase.
	ConstructionData *GameplayConstructionData `json:"ConstructionData,omitempty" yaml:"ConstructionData,omitempty"`
}

// GameplayHologramData stores hologram placement check toggles.
type GameplayHologramData struct {
	// DisableIsCollidingBBoxCheck allows placement despite bbox overlap.
	DisableIsCollidingBBoxCheck *bool `json:"disableIsCollidingBBoxCheck,omitempty" yaml:"disableIsCollidingBBoxCheck,omitempty"`
	// DisableIsCollidingPlayerCheck allows placement despite player overlap.
	DisableIsCollidingPlayerCheck *bool `json:"disableIsCollidingPlayerCheck,omitempty" yaml:"disableIsCollidingPlayerCheck,omitempty"`
	// DisableIsClippingRoofCheck allows clipping through roofs.
	DisableIsClippingRoofCheck *bool `json:"disableIsClippingRoofCheck,omitempty" yaml:"disableIsClippingRoofCheck,omitempty"`
	// DisableIsBaseViableCheck allows otherwise invalid base surfaces.
	DisableIsBaseViableCheck *bool `json:"disableIsBaseViableCheck,omitempty" yaml:"disableIsBaseViableCheck,omitempty"`
	// DisableIsCollidingGPlotCheck allows garden plot collisions.
	DisableIsCollidingGPlotCheck *bool `json:"disableIsCollidingGPlotCheck,omitempty" yaml:"disableIsCollidingGPlotCheck,omitempty"`
	// DisableIsCollidingAngleCheck allows placement beyond angle limits.
	DisableIsCollidingAngleCheck *bool `json:"disableIsCollidingAngleCheck,omitempty" yaml:"disableIsCollidingAngleCheck,omitempty"`
	// DisableIsPlacementPermittedCheck allows placement when denied by checks.
	DisableIsPlacementPermittedCheck *bool `json:"disableIsPlacementPermittedCheck,omitempty" yaml:"disableIsPlacementPermittedCheck,omitempty"`
	// DisableHeightPlacementCheck allows placement with limited height space.
	DisableHeightPlacementCheck *bool `json:"disableHeightPlacementCheck,omitempty" yaml:"disableHeightPlacementCheck,omitempty"`
	// DisableIsUnderwaterCheck allows underwater placement.
	DisableIsUnderwaterCheck *bool `json:"disableIsUnderwaterCheck,omitempty" yaml:"disableIsUnderwaterCheck,omitempty"`
	// DisableIsInTerrainCheck allows placement intersecting terrain.
	DisableIsInTerrainCheck *bool `json:"disableIsInTerrainCheck,omitempty" yaml:"disableIsInTerrainCheck,omitempty"`
	// DisableColdAreaBuildingCheck is old key name used by some worlds.
	DisableColdAreaBuildingCheck *bool `json:"disableColdAreaBuildingCheck,omitempty" yaml:"disableColdAreaBuildingCheck,omitempty"`
	// DisableColdAreaPlacementCheck is new key name used by some worlds.
	DisableColdAreaPlacementCheck *bool `json:"disableColdAreaPlacementCheck,omitempty" yaml:"disableColdAreaPlacementCheck,omitempty"`
	// DisallowedTypesInUnderground blocks listed kits in underground.
	DisallowedTypesInUnderground []string `json:"disallowedTypesInUnderground,omitempty" yaml:"disallowedTypesInUnderground,omitempty"`
}

// GameplayConstructionData stores construction check toggles.
type GameplayConstructionData struct {
	// DisablePerformRoofCheck allows construction clipping through roof.
	DisablePerformRoofCheck *bool `json:"disablePerformRoofCheck,omitempty" yaml:"disablePerformRoofCheck,omitempty"`
	// DisableIsCollidingCheck allows construction despite collisions.
	DisableIsCollidingCheck *bool `json:"disableIsCollidingCheck,omitempty" yaml:"disableIsCollidingCheck,omitempty"`
	// DisableDistanceCheck disables range restriction checks.
	DisableDistanceCheck *bool `json:"disableDistanceCheck,omitempty" yaml:"disableDistanceCheck,omitempty"`
}

// GameplayUIData stores UI-specific settings.
type GameplayUIData struct {
	// Use3DMap enables only 3D map behavior.
	Use3DMap *bool `json:"use3DMap,omitempty" yaml:"use3DMap,omitempty"`
	// HitIndicationData configures hit indicator rendering.
	HitIndicationData *GameplayHitIndicationData `json:"HitIndicationData,omitempty" yaml:"HitIndicationData,omitempty"`
}

// GameplayHitIndicationData stores hit indicator settings.
type GameplayHitIndicationData struct {
	// HitDirectionOverrideEnabled enables this section override.
	HitDirectionOverrideEnabled *bool `json:"hitDirectionOverrideEnabled,omitempty" yaml:"hitDirectionOverrideEnabled,omitempty"`
	// HitDirectionBehaviour selects behavior mode:
	// 0=disabled, 1=static, 2=dynamic.
	HitDirectionBehaviour *int `json:"hitDirectionBehaviour,omitempty" yaml:"hitDirectionBehaviour,omitempty"`
	// HitDirectionStyle selects style preset:
	// 0=splash, 1=spike, 2=arrow.
	HitDirectionStyle *int `json:"hitDirectionStyle,omitempty" yaml:"hitDirectionStyle,omitempty"`
	// HitDirectionIndicatorColorStr stores ARGB color string.
	HitDirectionIndicatorColorStr *string `json:"hitDirectionIndicatorColorStr,omitempty" yaml:"hitDirectionIndicatorColorStr,omitempty"`
	// HitDirectionMaxDuration sets maximum display duration.
	HitDirectionMaxDuration *float64 `json:"hitDirectionMaxDuration,omitempty" yaml:"hitDirectionMaxDuration,omitempty"`
	// HitDirectionBreakPointRelative sets fade-start fraction.
	HitDirectionBreakPointRelative *float64 `json:"hitDirectionBreakPointRelative,omitempty" yaml:"hitDirectionBreakPointRelative,omitempty"`
	// HitDirectionScatter sets randomized directional inaccuracy.
	HitDirectionScatter *float64 `json:"hitDirectionScatter,omitempty" yaml:"hitDirectionScatter,omitempty"`
	// HitIndicationPostProcessEnabled toggles postprocess effect.
	HitIndicationPostProcessEnabled *bool `json:"hitIndicationPostProcessEnabled,omitempty" yaml:"hitIndicationPostProcessEnabled,omitempty"`
}

// GameplayMapData stores map ownership/navigation settings.
type GameplayMapData struct {
	// IgnoreMapOwnership allows opening map without owning map item.
	IgnoreMapOwnership *bool `json:"ignoreMapOwnership,omitempty" yaml:"ignoreMapOwnership,omitempty"`
	// IgnoreNavItemsOwnership bypasses compass/GPS ownership checks.
	IgnoreNavItemsOwnership *bool `json:"ignoreNavItemsOwnership,omitempty" yaml:"ignoreNavItemsOwnership,omitempty"`
	// DisplayPlayerPosition toggles red player marker on map.
	DisplayPlayerPosition *bool `json:"displayPlayerPosition,omitempty" yaml:"displayPlayerPosition,omitempty"`
	// DisplayNavInfo toggles nav helper display on map.
	DisplayNavInfo *bool `json:"displayNavInfo,omitempty" yaml:"displayNavInfo,omitempty"`
}

// GameplayVehicleData stores vehicle-related settings.
type GameplayVehicleData struct {
	// BoatDecayMultiplier multiplies decay speed for boats.
	BoatDecayMultiplier *float64 `json:"boatDecayMultiplier,omitempty" yaml:"boatDecayMultiplier,omitempty"`
}
