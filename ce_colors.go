// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"hash/fnv"
	"strconv"
)

// syntheticLayerColor returns deterministic pseudo-random ARGB color for layer.
func syntheticLayerColor(name string, usageMask uint32, valueMask uint32) string {
	seed := hashSeed32(name) ^ usageMask ^ (valueMask << 1)
	return syntheticColorFromSeed(seed)
}

// syntheticTerritoryColor returns deterministic pseudo-random ARGB color.
func syntheticTerritoryColor(typeName string, territoryIndex int) string {
	indexKey := strconv.Itoa(territoryIndex + 1)
	seed := hashSeed32(typeName + ":" + indexKey)
	return syntheticColorFromSeed(seed)
}

// syntheticColorFromSeed builds opaque ARGB color from seed.
func syntheticColorFromSeed(seed uint32) string {
	// Keep colors away from pure black/white for better editor visibility.
	r := 40 + (((seed >> 0) & 0xFF) % 176)
	g := 40 + (((seed >> 8) & 0xFF) % 176)
	b := 40 + (((seed >> 16) & 0xFF) % 176)
	color := uint32(0xFF000000) | (r << 16) | (g << 8) | b

	return strconv.FormatUint(uint64(color), 10)
}

// hashSeed32 computes deterministic seed from text.
func hashSeed32(text string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(text))
	return h.Sum32()
}
