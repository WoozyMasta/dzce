// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "errors"

var (
	// ErrUnsupportedKind reports an unsupported CE file kind.
	ErrUnsupportedKind = errors.New("unsupported ce kind")

	// ErrUnsupportedValue reports unsupported value type for selected kind.
	ErrUnsupportedValue = errors.New("unsupported ce value")

	// ErrUnknownFileKind reports an unknown CE file kind by path.
	ErrUnknownFileKind = errors.New("unknown ce file kind")

	// ErrEconomyCoreCycle reports recursive economycore include cycle.
	ErrEconomyCoreCycle = errors.New("economycore include cycle")

	// ErrMergeConflict reports incompatible CE include merge values.
	ErrMergeConflict = errors.New("ce merge conflict")

	// ErrUnsupportedImageFormat reports unsupported mask image format.
	ErrUnsupportedImageFormat = errors.New("unsupported image format")

	// ErrNilAreaFlagsFile reports a missing `*AreaFlagsMapFile` input.
	ErrNilAreaFlagsFile = errors.New("nil area flags file")

	// ErrNilMaskImage reports a missing `*MaskImage` input.
	ErrNilMaskImage = errors.New("nil mask image")

	// ErrNilCEProjectConfig reports a missing `*CEProjectConfigFile` input.
	ErrNilCEProjectConfig = errors.New("nil CEProject config")

	// ErrMissingCEProjectLayers reports absent CEProject `<layers>` section.
	ErrMissingCEProjectLayers = errors.New("missing CEProject layers section")

	// ErrEmptyLayerName reports a layer entry with empty `name`.
	ErrEmptyLayerName = errors.New("layer with empty name")

	// ErrEmptyTerritoryDirPath reports an empty territory directory path.
	ErrEmptyTerritoryDirPath = errors.New("empty territory directory path")
)
