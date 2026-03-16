// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// readDirNamesSorted reads direct child names and returns deterministic order.
func readDirNamesSorted(path string) ([]string, error) {
	handle, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = handle.Close()
	}()

	names, err := handle.Readdirnames(-1)
	if err != nil {
		return nil, fmt.Errorf("readdir names: %w", err)
	}

	sort.Strings(names)

	return names, nil
}

// writeFile600 writes file payload with fixed user-only permissions.
func writeFile600(path string, data []byte) error {
	cleanPath := filepath.Clean(path)

	// #nosec G304,G703 -- paths are explicit API inputs by package design.
	return os.WriteFile(cleanPath, data, 0o600)
}
