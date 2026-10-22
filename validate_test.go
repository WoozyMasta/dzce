// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/woozymasta/lintkit/lint"
)

func TestValidateContentReportsDecodeFailures(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		data         []byte
		wantDiagCode lint.Code
		wantErr      error
	}{
		{
			name:         "json",
			path:         "cfggameplay.json",
			data:         []byte(`{"version":`),
			wantDiagCode: CodeParseInvalidJSON,
		},
		{
			name:         "json null root",
			path:         "cfggameplay.json",
			data:         []byte(`null`),
			wantDiagCode: CodeParseInvalidJSON,
			wantErr:      ErrInvalidJSONRoot,
		},
		{
			name:         "xml",
			path:         "db/types.xml",
			data:         []byte(`<types>`),
			wantDiagCode: CodeParseInvalidXML,
		},
		{
			name:         "xml multiple roots",
			path:         "db/types.xml",
			data:         []byte(`<types/><types/>`),
			wantDiagCode: CodeParseInvalidXML,
			wantErr:      ErrXMLMultipleRootElements,
		},
		{
			name: "binary",
			path: "areaflags.map",
			data: []byte("invalid"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			diagnostics, err := ValidateContent(test.path, test.data)
			if err == nil {
				t.Fatal("ValidateContent() error = nil, want decode error")
			}

			if test.wantDiagCode != 0 && !hasCode(diagnostics, test.wantDiagCode) {
				t.Fatalf("diagnostics do not contain code %s", codeToken(test.wantDiagCode))
			}
			if test.wantErr != nil && !errors.Is(err, test.wantErr) {
				t.Fatalf("ValidateContent() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestValidateContentReturnsExistingFindingsAfterDecode(t *testing.T) {
	content := []byte(
		`<?xml version="1.0"?><variables>` +
			`<var name="RespawnLimit" type="8" value="10"/>` +
			`</variables>`,
	)

	diagnostics, err := ValidateContent("db/globals.xml", content)
	if err != nil {
		t.Fatalf("ValidateContent() error: %v", err)
	}
	if !hasCode(diagnostics, CodeGlobalsInvalidTypeTag) {
		t.Fatalf("diagnostics do not contain %s", codeToken(CodeGlobalsInvalidTypeTag))
	}
}

func TestValidateContentUnknownKind(t *testing.T) {
	diagnostics, err := ValidateContent("unknown.bin", []byte{0, 1, 2})
	if !errors.Is(err, ErrUnknownFileKind) {
		t.Fatalf("ValidateContent() error = %v, want ErrUnknownFileKind", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %v, want none for unknown binary kind", diagnostics)
	}
}

func TestValidateFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "types.xml")
	if err := os.WriteFile(path, []byte(`<types/>`), 0o600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	diagnostics, err := ValidateFile(path)
	if err != nil {
		t.Fatalf("ValidateFile() error: %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("ValidateFile() diagnostics = %v, want none", diagnostics)
	}

	if _, err = ValidateFile(filepath.Join(t.TempDir(), "missing.xml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ValidateFile() missing file error = %v, want os.ErrNotExist", err)
	}
}
