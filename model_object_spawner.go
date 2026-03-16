// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// ObjectSpawnerFile is one object spawner JSON payload loaded via
// `cfggameplay.json -> WorldsData.objectSpawnersArr`.
type ObjectSpawnerFile struct {
	// Objects is map object spawn list.
	Objects []ObjectSpawnerEntry `json:"Objects,omitempty" yaml:"Objects,omitempty"`
}

// ObjectSpawnerEntry stores one object spawn entry.
type ObjectSpawnerEntry struct {
	// Scale is model scale multiplier.
	Scale *float64 `json:"scale,omitempty" yaml:"scale,omitempty"`
	// EnableCEPersistency toggles CE persistence integration for item.
	EnableCEPersistency *bool `json:"enableCEPersistency,omitempty" yaml:"enableCEPersistency,omitempty"`
	// CustomString is custom user payload for script-side handling.
	CustomString *string `json:"customString,omitempty" yaml:"customString,omitempty"`
	// Name is config class name or supported p3d model path.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Pos is world position vector (x, y, z).
	Pos []float64 `json:"pos,omitempty" yaml:"pos,omitempty"`
	// YPR is orientation vector (yaw, pitch, roll).
	YPR []float64 `json:"ypr,omitempty" yaml:"ypr,omitempty"`
}
