// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"errors"
	"fmt"
	"io"
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

// writeFileAtomic writes a complete file beside its destination before
// replacing it, preserving the destination's permission bits when it exists.
func writeFileAtomic(path string, data []byte) error {
	cleanPath, err := resolveAtomicPath(path)
	if err != nil {
		return fmt.Errorf("resolve destination: %w", err)
	}

	mode := os.FileMode(0o600)
	info, err := os.Stat(cleanPath)
	if err == nil {
		mode = info.Mode().Perm()
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat destination: %w", err)
	}

	temp, err := os.CreateTemp(filepath.Dir(cleanPath), ".dzce-tmp-*")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}

	tempPath := temp.Name()
	closed := false
	defer func() {
		if !closed {
			_ = temp.Close()
		}

		_ = os.Remove(tempPath)
	}()

	n, err := temp.Write(data)
	if err != nil {
		return fmt.Errorf("write temporary file: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("write temporary file: %w", io.ErrShortWrite)
	}

	if err = temp.Sync(); err != nil {
		return fmt.Errorf("sync temporary file: %w", err)
	}
	if err = temp.Chmod(mode); err != nil {
		return fmt.Errorf("set temporary file permissions: %w", err)
	}
	if err = temp.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}

	closed = true
	if err = replaceFileAtomic(tempPath, cleanPath); err != nil {
		return fmt.Errorf("replace destination: %w", err)
	}

	return nil
}

// resolveAtomicPath follows symlinks so replacing the file keeps os.WriteFile's
// behavior of writing to the link target, including when that target is absent.
func resolveAtomicPath(path string) (string, error) {
	cleanPath := filepath.Clean(path)

	for range 255 {
		info, err := os.Lstat(cleanPath)
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			linkTarget, readlinkErr := os.Readlink(cleanPath)
			if readlinkErr != nil {
				return "", readlinkErr
			}

			if !filepath.IsAbs(linkTarget) {
				linkTarget = filepath.Join(filepath.Dir(cleanPath), linkTarget)
			}

			cleanPath = filepath.Clean(linkTarget)
			continue
		}

		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		dir, err := filepath.EvalSymlinks(filepath.Dir(cleanPath))
		if err != nil {
			return "", err
		}

		return filepath.Join(dir, filepath.Base(cleanPath)), nil
	}

	return "", fmt.Errorf("too many symlinks in path %q", path)
}
