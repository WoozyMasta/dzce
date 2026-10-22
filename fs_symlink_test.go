// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomicFollowsSymlinks(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.json")
	link := filepath.Join(dir, "link.json")
	if err := os.WriteFile(target, []byte("before"), 0o600); err != nil {
		t.Fatalf("WriteFile() target error: %v", err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("Symlink() unavailable: %v", err)
	}

	if err := writeFileAtomic(link, []byte("after")); err != nil {
		t.Fatalf("writeFileAtomic() through symlink error: %v", err)
	}

	if info, err := os.Lstat(link); err != nil {
		t.Fatalf("Lstat() link error: %v", err)
	} else if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("link was replaced, mode = %v", info.Mode())
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile() target error: %v", err)
	}
	if string(data) != "after" {
		t.Fatalf("target content = %q, want %q", data, "after")
	}
}

func TestWriteFileAtomicCreatesDanglingSymlinkTarget(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "new-target.json")
	link := filepath.Join(dir, "link.json")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("Symlink() unavailable: %v", err)
	}

	if err := writeFileAtomic(link, []byte("created")); err != nil {
		t.Fatalf("writeFileAtomic() through dangling symlink error: %v", err)
	}

	if info, err := os.Lstat(link); err != nil {
		t.Fatalf("Lstat() link error: %v", err)
	} else if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("link was replaced, mode = %v", info.Mode())
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile() target error: %v", err)
	}
	if string(data) != "created" {
		t.Fatalf("target content = %q, want %q", data, "created")
	}
}
