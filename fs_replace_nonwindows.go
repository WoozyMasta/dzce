// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

//go:build !windows

package dzce

import "os"

// replaceFileAtomic atomically renames a same-directory file over its target.
func replaceFileAtomic(tempPath, targetPath string) error {
	return os.Rename(tempPath, targetPath)
}
