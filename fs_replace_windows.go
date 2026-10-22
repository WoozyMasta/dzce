// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// replaceFileProc invokes the Windows API for replacing an existing file.
var replaceFileProc = windows.NewLazySystemDLL("kernel32.dll").NewProc("ReplaceFileW")

// replaceFileAtomic replaces an existing file with ReplaceFileW
// and moves a new file into place when the destination does not exist.
func replaceFileAtomic(tempPath, targetPath string) error {
	tempName, err := windows.UTF16PtrFromString(tempPath)
	if err != nil {
		return fmt.Errorf("encode temporary path: %w", err)
	}

	targetName, err := windows.UTF16PtrFromString(targetPath)
	if err != nil {
		return fmt.Errorf("encode destination path: %w", err)
	}

	result, _, callErr := replaceFileProc.Call(
		uintptr(unsafe.Pointer(targetName)),
		uintptr(unsafe.Pointer(tempName)),
		0,
		0,
		0,
		0,
	)
	if result != 0 {
		return nil
	}

	if callErr == windows.ERROR_FILE_NOT_FOUND {
		return windows.MoveFileEx(tempName, targetName, windows.MOVEFILE_WRITE_THROUGH)
	}

	if callErr == windows.ERROR_SUCCESS {
		return errors.New("ReplaceFileW failed without an error code")
	}

	return fmt.Errorf("ReplaceFileW: %w", callErr)
}
