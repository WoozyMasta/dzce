// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// BunkerBroadcastFile configures the bunker locations, broadcast schedule, access-code options, and Morse playback.
// The game uses this file when `bunkerBroadcastEnabled` is enabled in `cfggameplay.json`;
// otherwise the broadcast system does not load it.
type BunkerBroadcastFile struct {
	// BunkerSettings sets how many bunkers participate in a broadcast batch and how their access codes are generated.
	BunkerSettings *BunkerBroadcastBunkerSettings `json:"bunkerSettings,omitempty" yaml:"bunkerSettings,omitempty"`
	// SchedulerSettings controls batch lifetime, report cadence, and the in-game time windows when broadcasts may run.
	SchedulerSettings *BunkerBroadcastSchedulerSettings `json:"schedulerSettings,omitempty" yaml:"schedulerSettings,omitempty"`
	// PlaybackSettings controls Morse tone lengths, pauses, and repetitions.
	PlaybackSettings *BunkerBroadcastPlaybackSettings `json:"playbackSettings,omitempty" yaml:"playbackSettings,omitempty"`
	// Proxies are world-space positions associated with bunker broadcasts; each position contains X, Y, and Z coordinates.
	Proxies []BunkerBroadcastProxy `json:"proxies,omitempty" yaml:"proxies,omitempty"`
}

// BunkerBroadcastBunkerSettings controls bunker selection and access-code generation for each broadcast batch.
type BunkerBroadcastBunkerSettings struct {
	// NumberOfActive is how many bunkers receive fresh codes in each batch. The game accepts values from 1 through 30.
	NumberOfActive *int `json:"numberOfActive,omitempty" yaml:"numberOfActive,omitempty"`
	// LockCodeLength is the number of digits in generated access codes. The game accepts values from 3 through 6.
	LockCodeLength *int `json:"lockCodeLength,omitempty" yaml:"lockCodeLength,omitempty"`
	// LockCodeMode selects the access-code generation mode. The game accepts only 0 or 1.
	LockCodeMode *int `json:"lockCodeMode,omitempty" yaml:"lockCodeMode,omitempty"`
}

// BunkerBroadcastSchedulerSettings configures when broadcasts are scheduled
// and how long each generated batch remains active.
type BunkerBroadcastSchedulerSettings struct {
	// BatchLifetimeLength is the batch lifetime in in-game minutes. Valid values are 1 through 1000.
	BatchLifetimeLength *int `json:"batchLifetimeLength,omitempty" yaml:"batchLifetimeLength,omitempty"`
	// ReportInterval is the scheduler report interval in in-game minutes. Valid values are 1 through 120.
	ReportInterval *int `json:"reportInterval,omitempty" yaml:"reportInterval,omitempty"`
	// Periods are in-game clock windows during which broadcasts may play.
	// The start and end use HHMM strings, and a window may continue past midnight.
	Periods []BunkerBroadcastPeriod `json:"periods,omitempty" yaml:"periods,omitempty"`
}

// BunkerBroadcastPeriod defines one allowed window on the in-game clock.
type BunkerBroadcastPeriod struct {
	// StartTime is the window start as four digits in HHMM form, such as `0700`.
	StartTime string `json:"startTime,omitempty" yaml:"startTime,omitempty"`
	// EndTime is the window end as four digits in HHMM form, such as `1930`.
	EndTime string `json:"endTime,omitempty"   yaml:"endTime,omitempty"`
}

// BunkerBroadcastPlaybackSettings controls the timing and repetition of Morse messages sent by bunker broadcasts.
// Duration values are in seconds.
type BunkerBroadcastPlaybackSettings struct {
	// DotDuration is the sound duration of one Morse dot, in seconds.
	DotDuration *float64 `json:"dotDuration,omitempty" yaml:"dotDuration,omitempty"`
	// DashDuration is the sound duration of one Morse dash, in seconds.
	DashDuration *float64 `json:"dashDuration,omitempty" yaml:"dashDuration,omitempty"`
	// SymbolGapDuration is the pause between symbols within a Morse character, in seconds.
	SymbolGapDuration *float64 `json:"symbolGapDuration,omitempty" yaml:"symbolGapDuration,omitempty"`
	// WordGapDuration is the pause between Morse words, in seconds.
	WordGapDuration *float64 `json:"wordGapDuration,omitempty" yaml:"wordGapDuration,omitempty"`
	// PauseBetweenBunkers is the pause before playback moves to the next bunker, in seconds.
	PauseBetweenBunkers *float64 `json:"pauseBetweenBunkers,omitempty" yaml:"pauseBetweenBunkers,omitempty"`
	// MessageRepetitions is the number of additional times a message is repeated.
	MessageRepetitions *int `json:"messageRepetitions,omitempty" yaml:"messageRepetitions,omitempty"`
	// BunkerRepetitions is the number of additional times the current bunker's Morse broadcast is played
	// before playback moves to the next bunker.
	BunkerRepetitions *int `json:"bunkerRepetitions,omitempty" yaml:"bunkerRepetitions,omitempty"`
}

// BunkerBroadcastProxy stores a world-space position where the game creates a proxy object for the bunker broadcast system.
type BunkerBroadcastProxy struct {
	// Position is the proxy's X, Y, and Z world coordinate.
	// The position should correspond to a bunker location that can be activated by the broadcast.
	Position []float64 `json:"position,omitempty" yaml:"position,omitempty"`
}

// ExternalLocksFile configures brute-force protection for external locks.
type ExternalLocksFile struct {
	// LocksSettings contains gameplay limits applied when players try external locks.
	LocksSettings *GameplayExternalLockData `json:"locksSettings,omitempty" yaml:"locksSettings,omitempty"`
}

// PlayerRestrictedAreaFile defines a restricted player area and safe positions for that area.
// The file is loaded through `cfggameplay.json -> WorldsData.playerRestrictedAreaFiles`.
type PlayerRestrictedAreaFile struct {
	// AreaName is an optional label used to identify the area in game diagnostics.
	AreaName string `json:"areaName,omitempty" yaml:"areaName,omitempty"`
	// PRABoxes defines rotated 3D boxes. Each box contains three XYZ vectors in order: size, rotation, and world position.
	PRABoxes [][][]float64 `json:"PRABoxes,omitempty" yaml:"PRABoxes,omitempty"`
	// SafePositions3D lists explicit XYZ destinations used when selecting a safe position for this area.
	SafePositions3D [][]float64 `json:"safePositions3D,omitempty" yaml:"safePositions3D,omitempty"`
}
