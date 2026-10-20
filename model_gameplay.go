// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// GameplayFile is the root of `cfggameplay.json`, the mission-level file for server-side gameplay rules.
// The server reads it when `enableCfgGameplayFile` is enabled in `serverDZ.cfg`.
type GameplayFile struct {
	// Version is the game's internal version marker for this configuration format.
	// It is maintained by the game and is not a gameplay setting.
	Version *int `json:"version,omitempty" yaml:"version,omitempty" jsonschema:"example=123"`
	// GeneralData controls general rules such as damage to base structures and containers, and the availability of respawn actions.
	GeneralData *GameplayGeneralData `json:"GeneralData,omitempty" yaml:"GeneralData,omitempty"`
	// PlayerData controls player-side rules, including stamina, movement, shock recovery, drowning, weapon obstruction, and spawn gear presets.
	PlayerData *GameplayPlayerData `json:"PlayerData,omitempty" yaml:"PlayerData,omitempty"`
	// WorldsData controls world-level settings such as night lighting, temperature and wetness effects, object spawners, and restricted areas.
	WorldsData *GameplayWorldsData `json:"WorldsData,omitempty" yaml:"WorldsData,omitempty"`
	// BaseBuildingData controls checks applied when previewing and constructing base-building objects.
	// Disabling a check permits placements that thecorresponding rule would otherwise reject.
	BaseBuildingData *GameplayBaseBuildingData `json:"BaseBuildingData,omitempty" yaml:"BaseBuildingData,omitempty"`
	// UIData controls gameplay-related interface behavior, including the directional hit indicator and map presentation.
	UIData *GameplayUIData `json:"UIData,omitempty" yaml:"UIData,omitempty"`
	// MapData controls whether map and navigation features require their corresponding inventory items and which navigation details are shown.
	MapData *GameplayMapData `json:"MapData,omitempty" yaml:"MapData,omitempty"`
	// VehicleData controls vehicle behavior such as boat decay and whether attached car or bike wheels can detach or take damage.
	VehicleData *GameplayVehicleData `json:"VehicleData,omitempty" yaml:"VehicleData,omitempty"`
	// ExternalLockData configures lockout thresholds and timers that slow repeated failed attempts to open an external lock.
	ExternalLockData *GameplayExternalLockData `json:"ExternalLockData,omitempty" yaml:"ExternalLockData,omitempty"`
	// SandstormData scales how frequently sandstorms occur in worlds that use the sandstorm system.
	SandstormData *GameplaySandstormData `json:"SandstormData,omitempty" yaml:"SandstormData,omitempty"`
}

// GameplayGeneralData controls mission-wide damage and respawn behavior.
type GameplayGeneralData struct {
	// DisableBaseDamage prevents damage to structures built through the base-building system, such as fences and watchtowers.
	DisableBaseDamage *bool `json:"disableBaseDamage,omitempty" yaml:"disableBaseDamage,omitempty" jsonschema:"example=false"`
	// DisableContainerDamage prevents damage to storage containers such as tents and barrels.
	// It is separate from damage to base-building structures.
	DisableContainerDamage *bool `json:"disableContainerDamage,omitempty" yaml:"disableContainerDamage,omitempty" jsonschema:"example=false"`
	// DisableRespawnDialog removes the dialog that lets a player choose a respawn type after requesting respawn.
	DisableRespawnDialog *bool `json:"disableRespawnDialog,omitempty" yaml:"disableRespawnDialog,omitempty" jsonschema:"example=false"`
	// DisableRespawnInUnconsciousness hides the pause-menu Respawn action while the player is unconscious, preventing respawn from that state.
	DisableRespawnInUnconsciousness *bool `json:"disableRespawnInUnconsciousness,omitempty" yaml:"disableRespawnInUnconsciousness,omitempty" jsonschema:"example=false"`
}

// GameplayPlayerData groups rules that affect a player's controls, movement, survival, and initial equipment.
type GameplayPlayerData struct {
	// DisablePersonalLight turns off the local light that normally illuminates nearby objects around a player at night.
	DisablePersonalLight *bool `json:"disablePersonalLight,omitempty" yaml:"disablePersonalLight,omitempty" jsonschema:"example=false"`
	// StaminaData sets the stamina pool and the rates at which common actions consume it, including sprinting, swimming, and climbing.
	StaminaData *GameplayStaminaData `json:"StaminaData,omitempty" yaml:"StaminaData,omitempty"`
	// ShockHandlingData sets how quickly shock recovers while conscious or unconscious and whether ammunition can modify that recovery.
	ShockHandlingData *GameplayShockHandlingData `json:"ShockHandlingData,omitempty" yaml:"ShockHandlingData,omitempty"`
	// MovementData sets transition times and turning rates for jogging and sprinting, and whether stamina affects movement inertia.
	MovementData *GameplayMovementData `json:"MovementData,omitempty" yaml:"MovementData,omitempty"`
	// DrowningData sets the per-second stamina, health, and shock loss while a character is drowning.
	DrowningData *GameplayDrowningData `json:"DrowningData,omitempty" yaml:"DrowningData,omitempty"`
	// WeaponObstructionData selects how the game handles a weapon blocked by static scenery or by dynamic objects.
	WeaponObstructionData *GameplayWeaponObstructionData `json:"WeaponObstructionData,omitempty" yaml:"WeaponObstructionData,omitempty"`
	// SpawnGearPresetFiles lists mission-relative JSON files that define server spawn gear presets.
	// Files may be in subdirectories; when the list is non-empty,
	// these presets control character equipment at spawn and replace the usual StartingEquipSetup behavior.
	SpawnGearPresetFiles []string `json:"spawnGearPresetFiles,omitempty" yaml:"spawnGearPresetFiles,omitempty"`
}

// GameplayStaminaData controls the stamina pool and action-specific stamina costs.
// Modifiers scale the rate of consumption: increasing one makes the associated action consume stamina faster.
type GameplayStaminaData struct {
	// SprintStaminaModifierErc scales stamina consumption while sprinting upright.
	SprintStaminaModifierErc *float64 `json:"sprintStaminaModifierErc,omitempty" yaml:"sprintStaminaModifierErc,omitempty" jsonschema:"example=1"`
	// SprintStaminaModifierCro scales stamina consumption while sprinting crouched.
	SprintStaminaModifierCro *float64 `json:"sprintStaminaModifierCro,omitempty" yaml:"sprintStaminaModifierCro,omitempty" jsonschema:"example=1"`
	// StaminaWeightLimitThreshold is the amount of stamina points, divided by 1000, excluded from the load-based reduction of maximum stamina.
	StaminaWeightLimitThreshold *float64 `json:"staminaWeightLimitThreshold,omitempty" yaml:"staminaWeightLimitThreshold,omitempty" jsonschema:"example=6000"`
	// StaminaMax sets the maximum stamina pool. Setting it to zero may produce unexpected results.
	StaminaMax *float64 `json:"staminaMax,omitempty" yaml:"staminaMax,omitempty" jsonschema:"example=100"`
	// StaminaKgToStaminaPercentPenalty scales the reduction of maximum stamina calculated from the character's carried load.
	StaminaKgToStaminaPercentPenalty *float64 `json:"staminaKgToStaminaPercentPenalty,omitempty" yaml:"staminaKgToStaminaPercentPenalty,omitempty" jsonschema:"example=1.75"`
	// StaminaMinCap sets the minimum size of the stamina cap. Setting it to zero may produce unexpected results.
	StaminaMinCap *float64 `json:"staminaMinCap,omitempty" yaml:"staminaMinCap,omitempty" jsonschema:"example=5"`
	// SprintSwimmingStaminaModifier scales stamina consumption during fast swimming.
	SprintSwimmingStaminaModifier *float64 `json:"sprintSwimmingStaminaModifier,omitempty" yaml:"sprintSwimmingStaminaModifier,omitempty" jsonschema:"example=1"`
	// SprintLadderStaminaModifier scales stamina consumption while climbing a ladder quickly.
	SprintLadderStaminaModifier *float64 `json:"sprintLadderStaminaModifier,omitempty" yaml:"sprintLadderStaminaModifier,omitempty" jsonschema:"example=1"`
	// MeleeStaminaModifier scales stamina consumption for heavy melee attacks and evasive moves.
	MeleeStaminaModifier *float64 `json:"meleeStaminaModifier,omitempty" yaml:"meleeStaminaModifier,omitempty" jsonschema:"example=1"`
	// ObstacleTraversalStaminaModifier scales stamina consumption when jumping, climbing, or vaulting over obstacles.
	ObstacleTraversalStaminaModifier *float64 `json:"obstacleTraversalStaminaModifier,omitempty" yaml:"obstacleTraversalStaminaModifier,omitempty" jsonschema:"example=1"`
	// HoldBreathStaminaModifier scales stamina consumption while holding breath.
	HoldBreathStaminaModifier *float64 `json:"holdBreathStaminaModifier,omitempty" yaml:"holdBreathStaminaModifier,omitempty" jsonschema:"example=1"`
}

// GameplayShockHandlingData controls how quickly the player's shock value recovers in each consciousness state.
type GameplayShockHandlingData struct {
	// ShockRefillSpeedConscious sets shock recovery per second while the player is conscious.
	ShockRefillSpeedConscious *float64 `json:"shockRefillSpeedConscious,omitempty" yaml:"shockRefillSpeedConscious,omitempty" jsonschema:"example=5"`
	// ShockRefillSpeedUnconscious sets shock recovery per second while the player is unconscious.
	ShockRefillSpeedUnconscious *float64 `json:"shockRefillSpeedUnconscious,omitempty" yaml:"shockRefillSpeedUnconscious,omitempty" jsonschema:"example=1"`
	// AllowRefillSpeedModifier allows ammunition settings to modify shock recovery speed, which can change how quickly a shot player wakes up.
	AllowRefillSpeedModifier *bool `json:"allowRefillSpeedModifier,omitempty" yaml:"allowRefillSpeedModifier,omitempty" jsonschema:"example=true"`
}

// GameplayMovementData controls how quickly jogging and sprinting movement transitions blend and how quickly the character can rotate.
// Transition times and rotation speeds have a documented minimum of 0.01.
type GameplayMovementData struct {
	// TimeToStrafeJog sets the blend time for diagonal movement while jogging.
	TimeToStrafeJog *float64 `json:"timeToStrafeJog,omitempty" yaml:"timeToStrafeJog,omitempty" jsonschema:"minimum=0.01,example=0.1"`
	// RotationSpeedJog sets the character's rotation speed while jogging.
	RotationSpeedJog *float64 `json:"rotationSpeedJog,omitempty" yaml:"rotationSpeedJog,omitempty" jsonschema:"minimum=0.01,example=0.15"`
	// TimeToSprint sets the transition time from jogging to sprinting.
	TimeToSprint *float64 `json:"timeToSprint,omitempty" yaml:"timeToSprint,omitempty" jsonschema:"minimum=0.01,example=0.45"`
	// TimeToStrafeSprint sets the blend time for diagonal movement while sprinting.
	TimeToStrafeSprint *float64 `json:"timeToStrafeSprint,omitempty" yaml:"timeToStrafeSprint,omitempty" jsonschema:"minimum=0.01,example=0.3"`
	// RotationSpeedSprint sets the character's rotation speed while sprinting.
	RotationSpeedSprint *float64 `json:"rotationSpeedSprint,omitempty" yaml:"rotationSpeedSprint,omitempty" jsonschema:"minimum=0.01,example=0.15"`
	// AllowStaminaAffectInertia lets the character's current stamina affect movement inertia.
	AllowStaminaAffectInertia *bool `json:"allowStaminaAffectInertia,omitempty" yaml:"allowStaminaAffectInertia,omitempty" jsonschema:"example=true"`
}

// GameplayDrowningData sets how quickly player resources are depleted while drowning. Each speed is expressed as loss per second.
type GameplayDrowningData struct {
	// StaminaDepletionSpeed sets stamina points lost per second while drowning.
	StaminaDepletionSpeed *float64 `json:"staminaDepletionSpeed,omitempty" yaml:"staminaDepletionSpeed,omitempty" jsonschema:"example=10"`
	// HealthDepletionSpeed sets health lost per second while drowning.
	HealthDepletionSpeed *float64 `json:"healthDepletionSpeed,omitempty" yaml:"healthDepletionSpeed,omitempty" jsonschema:"example=10"`
	// ShockDepletionSpeed sets shock lost per second while drowning.
	ShockDepletionSpeed *float64 `json:"shockDepletionSpeed,omitempty" yaml:"shockDepletionSpeed,omitempty" jsonschema:"example=10"`
}

// GameplayWeaponObstructionData selects the response when an obstacle blocks a weapon.
// Static objects are scenery; dynamic objects include entities that can move.
// Mode 1 lets the weapon lift to clear the obstacle, while mode 2 keeps it obstructed.
type GameplayWeaponObstructionData struct {
	// StaticMode selects obstruction handling for static objects:
	//  - 0 does not obstruct but still lifts the weapon,
	//  - 1 obstructs before lifting it,
	//  - 2 always obstructs without lifting.
	StaticMode *int `json:"staticMode,omitempty" yaml:"staticMode,omitempty" jsonschema:"enum=0,enum=1,enum=2,example=1"`
	// DynamicMode selects obstruction handling for dynamic objects:
	//  - 0 does not obstruct but still lifts the weapon,
	//  - 1 obstructs before lifting it,
	//  - 2 always obstructs without lifting.
	DynamicMode *int `json:"dynamicMode,omitempty" yaml:"dynamicMode,omitempty" jsonschema:"enum=0,enum=1,enum=2,example=1"`
}

// GameplayWorldsData controls environment behavior and additional world configuration files loaded by the mission.
type GameplayWorldsData struct {
	// BunkerBroadcastEnabled controls Morse-code broadcasts associated with discovering the bunker.
	BunkerBroadcastEnabled *bool `json:"bunkerBroadcastEnabled,omitempty" yaml:"bunkerBroadcastEnabled,omitempty" jsonschema:"example=true"`
	// LightingConfig selects the night lighting preset: 0 is brighter and 1 is darker. World-specific settings may define additional values.
	LightingConfig *int `json:"lightingConfig,omitempty" yaml:"lightingConfig,omitempty" jsonschema:"example=0"`
	// ObjectSpawnersArr lists mission-relative JSON files containing object spawner data. An empty or omitted list loads no custom spawner files.
	ObjectSpawnersArr []string `json:"objectSpawnersArr,omitempty" yaml:"objectSpawnersArr,omitempty"`
	// EnvironmentMinTemps sets the minimum ambient temperature for each month, ordered January through December.
	// The game expects exactly 12 values.
	EnvironmentMinTemps []float64 `json:"environmentMinTemps,omitempty" yaml:"environmentMinTemps,omitempty" jsonschema:"minItems=12,maxItems=12"`
	// EnvironmentMaxTemps sets the maximum ambient temperature for each month, ordered January through December.
	// The game expects exactly 12 values.
	EnvironmentMaxTemps []float64 `json:"environmentMaxTemps,omitempty" yaml:"environmentMaxTemps,omitempty" jsonschema:"minItems=12,maxItems=12"`
	// WetnessWeightModifiers scales item weight at the five wetness levels, in order: dry, damp, wet, soaked, and drenched.
	WetnessWeightModifiers []float64 `json:"wetnessWeightModifiers,omitempty" yaml:"wetnessWeightModifiers,omitempty" jsonschema:"minItems=5,maxItems=5"`
	// PlayerRestrictedAreaFiles lists mission-relative JSON files that define areas with player restrictions.
	// A non-empty list enables loading them.
	PlayerRestrictedAreaFiles []string `json:"playerRestrictedAreaFiles,omitempty" yaml:"playerRestrictedAreaFiles,omitempty"`
}

// GameplayBaseBuildingData separates checks performed while placing an object hologram from checks performed when the construction action is completed.
type GameplayBaseBuildingData struct {
	// HologramData controls preview-time placement checks. These checks decide whether the hologram is considered a valid placement before construction.
	HologramData *GameplayHologramData `json:"HologramData,omitempty" yaml:"HologramData,omitempty"`
	// ConstructionData controls checks and item handling during construction, after the placement preview has been accepted.
	ConstructionData *GameplayConstructionData `json:"ConstructionData,omitempty" yaml:"ConstructionData,omitempty"`
}

// GameplayHologramData controls validation performed while positioning a base-building hologram.
// Each Disable field bypasses one named placement check, allowing placement in situations that check would reject.
type GameplayHologramData struct {
	// DisableIsCollidingBBoxCheck permits placement when the hologram's bounding box overlaps objects in the world.
	DisableIsCollidingBBoxCheck *bool `json:"disableIsCollidingBBoxCheck,omitempty" yaml:"disableIsCollidingBBoxCheck,omitempty" jsonschema:"example=false"`
	// DisableIsCollidingPlayerCheck permits placement when the hologram overlaps a player.
	DisableIsCollidingPlayerCheck *bool `json:"disableIsCollidingPlayerCheck,omitempty" yaml:"disableIsCollidingPlayerCheck,omitempty" jsonschema:"example=false"`
	// DisableIsClippingRoofCheck permits placement that would clip through a roof.
	DisableIsClippingRoofCheck *bool `json:"disableIsClippingRoofCheck,omitempty" yaml:"disableIsClippingRoofCheck,omitempty" jsonschema:"example=false"`
	// DisableIsBaseViableCheck permits placement on dynamic objects or other surfaces that the game would otherwise reject as unsuitable.
	DisableIsBaseViableCheck *bool `json:"disableIsBaseViableCheck,omitempty" yaml:"disableIsBaseViableCheck,omitempty" jsonschema:"example=false"`
	// DisableIsCollidingGPlotCheck permits garden plots on surface types that are normally incompatible.
	DisableIsCollidingGPlotCheck *bool `json:"disableIsCollidingGPlotCheck,omitempty" yaml:"disableIsCollidingGPlotCheck,omitempty" jsonschema:"example=false"`
	// DisableIsCollidingAngleCheck permits placement even when its roll, pitch, or yaw exceeds the normal angle limits.
	DisableIsCollidingAngleCheck *bool `json:"disableIsCollidingAngleCheck,omitempty" yaml:"disableIsCollidingAngleCheck,omitempty" jsonschema:"example=false"`
	// DisableIsPlacementPermittedCheck bypasses the general permission check that can reject a placement after the basic placement checks.
	DisableIsPlacementPermittedCheck *bool `json:"disableIsPlacementPermittedCheck,omitempty" yaml:"disableIsPlacementPermittedCheck,omitempty" jsonschema:"example=false"`
	// DisableHeightPlacementCheck permits placement when there is not enough vertical clearance for the object.
	DisableHeightPlacementCheck *bool `json:"disableHeightPlacementCheck,omitempty" yaml:"disableHeightPlacementCheck,omitempty" jsonschema:"example=false"`
	// DisableIsUnderwaterCheck permits placement below the water surface.
	DisableIsUnderwaterCheck *bool `json:"disableIsUnderwaterCheck,omitempty" yaml:"disableIsUnderwaterCheck,omitempty" jsonschema:"example=false"`
	// DisableIsInTerrainCheck permits placement when the object intersects the terrain.
	DisableIsInTerrainCheck *bool `json:"disableIsInTerrainCheck,omitempty" yaml:"disableIsInTerrainCheck,omitempty" jsonschema:"example=false"`
	// DisableColdAreaBuildingCheck is the legacy spelling found in some official mission files.
	// The current gameplay script uses DisableColdAreaPlacementCheck for this rule.
	DisableColdAreaBuildingCheck *bool `json:"disableColdAreaBuildingCheck,omitempty" yaml:"disableColdAreaBuildingCheck,omitempty" jsonschema:"example=true"`
	// DisableColdAreaPlacementCheck permits garden-plot placement on frozen ground.
	// This is the key exposed by the current gameplay script.
	DisableColdAreaPlacementCheck *bool `json:"disableColdAreaPlacementCheck,omitempty" yaml:"disableColdAreaPlacementCheck,omitempty" jsonschema:"example=true"`
	// DisallowedTypesInUnderground lists item types that cannot be placed in underground areas.
	// The restriction also applies to inherited types.
	DisallowedTypesInUnderground []string `json:"disallowedTypesInUnderground,omitempty" yaml:"disallowedTypesInUnderground,omitempty"`
}

// GameplayConstructionData controls validation and item handling while a player completes base construction.
// These settings are distinct from the checks used to position the hologram.
type GameplayConstructionData struct {
	// DisablePerformRoofCheck permits construction that clips through a roof.
	DisablePerformRoofCheck *bool `json:"disablePerformRoofCheck,omitempty" yaml:"disablePerformRoofCheck,omitempty" jsonschema:"example=false"`
	// DisableIsCollidingCheck permits construction when the object collides with objects in the world.
	DisableIsCollidingCheck *bool `json:"disableIsCollidingCheck,omitempty" yaml:"disableIsCollidingCheck,omitempty" jsonschema:"example=false"`
	// DisableDistanceCheck bypasses the normal distance restriction applied to the player during construction.
	DisableDistanceCheck *bool `json:"disableDistanceCheck,omitempty" yaml:"disableDistanceCheck,omitempty" jsonschema:"example=false"`
	// DisableItemFalling prevents construction items dropped during the construction process from falling.
	DisableItemFalling *bool `json:"disableItemFalling,omitempty" yaml:"disableItemFalling,omitempty" jsonschema:"example=false"`
	// DisableDropAttachmentsOnFall keeps attached items on a construction object when it falls.
	DisableDropAttachmentsOnFall *bool `json:"disableDropAttachmentsOnFall,omitempty" yaml:"disableDropAttachmentsOnFall,omitempty" jsonschema:"example=false"`
	// DisableRefundMaterialsOnFall prevents construction materials from being refunded when the construction object falls.
	DisableRefundMaterialsOnFall *bool `json:"disableRefundMaterialsOnFall,omitempty" yaml:"disableRefundMaterialsOnFall,omitempty" jsonschema:"example=false"`
}

// GameplayUIData controls interface features that affect gameplay feedback or navigation.
type GameplayUIData struct {
	// Use3DMap enables the 3D map and disables the default 2D map overlay.
	Use3DMap *bool `json:"use3DMap,omitempty" yaml:"use3DMap,omitempty" jsonschema:"example=false"`
	// HitIndicationData configures the direction marker shown when the player takes damage, including its appearance and timing.
	HitIndicationData *GameplayHitIndicationData `json:"HitIndicationData,omitempty" yaml:"HitIndicationData,omitempty"`
}

// GameplayHitIndicationData controls the directional damage indicator.
// The game uses this section only when HitDirectionOverrideEnabled is true.
type GameplayHitIndicationData struct {
	// HitDirectionOverrideEnabled tells the game to use the values in this section.
	// Unspecified fields are treated as zero when the section is used.
	HitDirectionOverrideEnabled *bool `json:"hitDirectionOverrideEnabled,omitempty" yaml:"hitDirectionOverrideEnabled,omitempty" jsonschema:"example=false"`
	// HitDirectionBehaviour selects indicator motion:
	//   - 0 disables it,
	//   - 1 keeps it static,
	//   - 2 uses the moving dynamic behavior.
	HitDirectionBehaviour *int `json:"hitDirectionBehaviour,omitempty" yaml:"hitDirectionBehaviour,omitempty" jsonschema:"enum=0,enum=1,enum=2,example=1"`
	// HitDirectionStyle selects the indicator graphic: 0 splash, 1 spike, or 2 arrow.
	HitDirectionStyle *int `json:"hitDirectionStyle,omitempty" yaml:"hitDirectionStyle,omitempty" jsonschema:"enum=0,enum=1,enum=2,example=0"`
	// HitDirectionIndicatorColorStr sets the indicator color as a hexadecimal ARGB string, for example `0xFFBB0A1E`.
	HitDirectionIndicatorColorStr *string `json:"hitDirectionIndicatorColorStr,omitempty" yaml:"hitDirectionIndicatorColorStr,omitempty" jsonschema:"example=0xffbb0a1e"`
	// HitDirectionMaxDuration sets the upper bound for how long the indicator is displayed.
	// Actual duration varies with hit severity and is between 60% and 100% of this value.
	HitDirectionMaxDuration *float64 `json:"hitDirectionMaxDuration,omitempty" yaml:"hitDirectionMaxDuration,omitempty" jsonschema:"example=2"`
	// HitDirectionBreakPointRelative is the fraction of the actual display duration after which the indicator starts fading.
	// 0 fades immediately; 1 leaves no time for fading.
	HitDirectionBreakPointRelative *float64 `json:"hitDirectionBreakPointRelative,omitempty" yaml:"hitDirectionBreakPointRelative,omitempty" jsonschema:"minimum=0,maximum=1,example=0.2"`
	// HitDirectionScatter is the maximum random angular offset applied to the indicated hit direction. The offset is randomized in both directions.
	HitDirectionScatter *float64 `json:"hitDirectionScatter,omitempty" yaml:"hitDirectionScatter,omitempty" jsonschema:"example=10"`
	// HitIndicationPostProcessEnabled controls the legacy red-flash postprocess effect that appears when the player is hit.
	HitIndicationPostProcessEnabled *bool `json:"hitIndicationPostProcessEnabled,omitempty" yaml:"hitIndicationPostProcessEnabled,omitempty" jsonschema:"example=true"`
}

// GameplayMapData controls map access and the navigation information shown on the map interface.
type GameplayMapData struct {
	// IgnoreMapOwnership allows the player to open the map without carrying a map item.
	IgnoreMapOwnership *bool `json:"ignoreMapOwnership,omitempty" yaml:"ignoreMapOwnership,omitempty" jsonschema:"example=false"`
	// IgnoreNavItemsOwnership allows map navigation helpers without requiring a compass or GPS receiver.
	IgnoreNavItemsOwnership *bool `json:"ignoreNavItemsOwnership,omitempty" yaml:"ignoreNavItemsOwnership,omitempty" jsonschema:"example=false"`
	// DisplayPlayerPosition shows the player's red position marker and heading on the map.
	DisplayPlayerPosition *bool `json:"displayPlayerPosition,omitempty" yaml:"displayPlayerPosition,omitempty" jsonschema:"example=false"`
	// DisplayNavInfo controls whether compass and GPS information appears in the map legend, including when the player carries those items.
	DisplayNavInfo *bool `json:"displayNavInfo,omitempty" yaml:"displayNavInfo,omitempty" jsonschema:"example=true"`
}

// GameplayVehicleData controls vehicle decay and wheel behavior.
// Wheel fields apply to attached wheels; they do not determine whether a vehicle accepts a particular wheel type.
type GameplayVehicleData struct {
	// BoatDecayMultiplier scales how quickly boats decay.
	// A larger multiplier makes decay faster; 1 leaves the normal rate unchanged.
	BoatDecayMultiplier *float64 `json:"boatDecayMultiplier,omitempty" yaml:"boatDecayMultiplier,omitempty" jsonschema:"example=1"`
	// CanDetachAttachedCarWheels allows attached car wheels to detach.
	CanDetachAttachedCarWheels *bool `json:"canDetachAttachedCarWheels,omitempty" yaml:"canDetachAttachedCarWheels,omitempty" jsonschema:"example=true"`
	// CanDamageAttachedCarWheels allows damage to reduce the condition of wheels attached to cars.
	CanDamageAttachedCarWheels *bool `json:"canDamageAttachedCarWheels,omitempty" yaml:"canDamageAttachedCarWheels,omitempty" jsonschema:"example=true"`
	// CanDetachAttachedBikeWheels allows attached motorcycle wheels to detach.
	CanDetachAttachedBikeWheels *bool `json:"canDetachAttachedBikeWheels,omitempty" yaml:"canDetachAttachedBikeWheels,omitempty" jsonschema:"example=true"`
	// CanDamageAttachedBikeWheels allows damage to reduce the condition of wheels attached to motorcycles.
	CanDamageAttachedBikeWheels *bool `json:"canDamageAttachedBikeWheels,omitempty" yaml:"canDamageAttachedBikeWheels,omitempty" jsonschema:"example=true"`
}

// GameplayExternalLockData controls temporary lockouts after failed attempts to open an external lock.
// Thresholds count failed attempts; durations are expressed in seconds.
type GameplayExternalLockData struct {
	// ProtectionCountStageOne sets the number of failed attempts before the first protection stage is applied.
	ProtectionCountStageOne *int `json:"protectionCountStageOne,omitempty" yaml:"protectionCountStageOne,omitempty" jsonschema:"example=3"`
	// ProtectionCountStageTwo sets the number of failed attempts before the second protection stage is applied.
	ProtectionCountStageTwo *int `json:"protectionCountStageTwo,omitempty" yaml:"protectionCountStageTwo,omitempty" jsonschema:"example=1"`
	// ProtectionTime sets the duration of the protection lockout in seconds.
	ProtectionTime *float64 `json:"protectionTime,omitempty" yaml:"protectionTime,omitempty" jsonschema:"example=30"`
	// ProtectionResetTime sets how long the failed-attempt counter remains active before it is reset, in seconds.
	ProtectionResetTime *float64 `json:"protectionResetTime,omitempty" yaml:"protectionResetTime,omitempty" jsonschema:"example=3600"`
}

// GameplaySandstormData controls the occurrence rate of sandstorms.
type GameplaySandstormData struct {
	// SandstormFrequency sets how often sandstorms occur:
	//  - 0 disables them,
	//  - 1 uses the normal rate,
	//  - 2 makes them occur twice as often.
	SandstormFrequency *float64 `json:"sandstormFrequency,omitempty" yaml:"sandstormFrequency,omitempty" jsonschema:"enum=0,enum=1,enum=2,example=1"`
}
