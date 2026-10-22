// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReplaceExistingFileWindows covers Windows replacement
// of an existing destination through the Windows replacement path.
func TestReplaceExistingFileWindows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte("before"), 0o600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	if err := writeFileAtomic(path, []byte("after")); err != nil {
		t.Fatalf("writeFileAtomic() replacement error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if string(data) != "after" {
		t.Fatalf("file content = %q, want %q", data, "after")
	}
}
