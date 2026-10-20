// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// ObjectSpawnerFile is a mission object list loaded through `cfggameplay.json -> WorldsData.objectSpawnersArr`.
// Its entries are spawned into the world when the mission starts.
// Large object lists can affect both server and client performance.
type ObjectSpawnerFile struct {
	// Objects are the map objects to create with the configured class or model, transform, and persistence behavior.
	Objects []ObjectSpawnerEntry `json:"Objects,omitempty" yaml:"Objects,omitempty"`
}

// ObjectSpawnerEntry defines one object placement created by the mission object spawner.
type ObjectSpawnerEntry struct {
	// Scale is the size multiplier applied to the spawned object or model.
	Scale *float64 `json:"scale,omitempty" yaml:"scale,omitempty"`
	// EnableCEPersistency spawns the object without Central Economy persistence.
	// If a player takes it, persistence is enabled when the item supports it.
	EnableCEPersistency *bool `json:"enableCEPersistency,omitempty" yaml:"enableCEPersistency,omitempty"`
	// CustomString passes custom data to the spawned class's OnSpawnByObjectSpawner script callback.
	CustomString *string `json:"customString,omitempty" yaml:"customString,omitempty"`
	// Name is either a game config class name or a P3D model path.
	// P3D spawning is limited to DZ/plants, DZ/plants_bliss, DZ/plants_sakhal, DZ/rocks, DZ/rocks_bliss, and DZ/rocks_sakhal.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Pos is the object's world position as an X/Y/Z vector in meters; Y is the vertical coordinate.
	Pos []float64 `json:"pos,omitempty" yaml:"pos,omitempty" jsonschema:"minItems=3,maxItems=3"`
	// YPR is the object's rotation in yaw, pitch, and roll order, in degrees. An empty vector leaves rotation unspecified.
	YPR []float64 `json:"ypr,omitempty" yaml:"ypr,omitempty" jsonschema:"minItems=3,maxItems=3"`
}
