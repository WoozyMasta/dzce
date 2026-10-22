// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomicReplacesAndPreservesPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte("old"), 0o640); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() before write error: %v", err)
	}

	if err = writeFileAtomic(path, []byte("new config")); err != nil {
		t.Fatalf("writeFileAtomic() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if string(data) != "new config" {
		t.Fatalf("file content = %q, want %q", data, "new config")
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() after write error: %v", err)
	}
	if after.Mode().Perm() != before.Mode().Perm() {
		t.Fatalf("file permissions = %04o, want %04o", after.Mode().Perm(), before.Mode().Perm())
	}

	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatalf("ReadDir() error: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != filepath.Base(path) {
		t.Fatalf("directory entries after write = %v, want only %q", entries, filepath.Base(path))
	}
}

func TestWriteFileAtomicCreatesWithPrivatePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new.json")
	if err := writeFileAtomic(path, []byte("new")); err != nil {
		t.Fatalf("writeFileAtomic() error: %v", err)
	}

	referencePath := filepath.Join(dir, "reference.json")
	if err := os.WriteFile(referencePath, []byte("reference"), 0o600); err != nil {
		t.Fatalf("WriteFile() reference error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	referenceInfo, err := os.Stat(referencePath)
	if err != nil {
		t.Fatalf("Stat() reference error: %v", err)
	}
	if info.Mode().Perm() != referenceInfo.Mode().Perm() {
		t.Fatalf("file permissions = %04o, want os.WriteFile(0600) behavior %04o", info.Mode().Perm(), referenceInfo.Mode().Perm())
	}
}
