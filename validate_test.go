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
		},
		{
			name:         "binary",
			path:         "areaflags.map",
			data:         []byte("invalid"),
			wantDiagCode: CodeParseInvalidContent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidateContent(test.path, test.data)
			if result.Status != AnalysisStatusInvalid {
				t.Fatalf("ValidateContent() status = %q, want %q", result.Status, AnalysisStatusInvalid)
			}
			if result.ErrorCount == 0 {
				t.Fatal("ValidateContent() ErrorCount = 0, want at least one error")
			}
			if !hasCode(result.Diagnostics, test.wantDiagCode) {
				t.Fatalf("diagnostics do not contain code %s", codeToken(test.wantDiagCode))
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

	result := ValidateContent("db/globals.xml", content)
	if result.Status != AnalysisStatusInvalid {
		t.Fatalf("ValidateContent() status = %q, want %q", result.Status, AnalysisStatusInvalid)
	}
	if !hasCode(result.Diagnostics, CodeGlobalsInvalidTypeTag) {
		t.Fatalf("diagnostics do not contain %s", codeToken(CodeGlobalsInvalidTypeTag))
	}
}

func TestValidateContentUnknownKind(t *testing.T) {
	result := ValidateContent("unknown.bin", []byte{0, 1, 2})
	if result.Status != AnalysisStatusUnsupported {
		t.Fatalf("ValidateContent() status = %q, want %q", result.Status, AnalysisStatusUnsupported)
	}
	if result.ErrorCount != 0 || len(result.Diagnostics) != 0 {
		t.Fatalf("ValidateContent() result = %+v, want unsupported without findings", result)
	}
}

func TestValidateContentUnknownXMLRootIsUnsupported(t *testing.T) {
	result := ValidateContent("unknown.xml", []byte(`<unsupported/>`))
	if result.Status != AnalysisStatusUnsupported {
		t.Fatalf("ValidateContent() status = %q, want %q", result.Status, AnalysisStatusUnsupported)
	}
	if !hasCode(result.Diagnostics, CodeParseUnknownRoot) {
		t.Fatalf("diagnostics do not contain %s", codeToken(CodeParseUnknownRoot))
	}
}

func TestAnalyzeLintContentReturnsAnalysisResult(t *testing.T) {
	result := AnalyzeLintContent("cfggameplay.json", []byte(`{"version":`))
	if result.Status != AnalysisStatusInvalid || result.ErrorCount != 1 {
		t.Fatalf("AnalyzeLintContent() = %+v, want invalid with one error", result)
	}
	if !hasCode(result.Diagnostics, CodeParseInvalidJSON) {
		t.Fatalf("diagnostics do not contain %s", codeToken(CodeParseInvalidJSON))
	}
}

func TestValidateContentCountsErrorSeverity(t *testing.T) {
	result := makeAnalysisResult([]lint.Diagnostic{
		{Severity: lint.SeverityWarning},
		{Severity: lint.SeverityInfo},
		{Severity: lint.SeverityNotice},
		{Severity: lint.SeverityError},
	})
	if result.Status != AnalysisStatusInvalid || result.ErrorCount != 1 {
		t.Fatalf("makeAnalysisResult() = %+v, want invalid with one error", result)
	}

	result = makeAnalysisResult([]lint.Diagnostic{
		{Severity: lint.SeverityWarning},
		{Severity: lint.SeverityInfo},
		{Severity: lint.SeverityNotice},
	})
	if result.Status != AnalysisStatusOK || result.ErrorCount != 0 {
		t.Fatalf("makeAnalysisResult() = %+v, want ok with no errors", result)
	}
}

func TestValidateFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "types.xml")
	if err := os.WriteFile(path, []byte(`<types/>`), 0o600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	result, err := ValidateFile(path)
	if err != nil {
		t.Fatalf("ValidateFile() error: %v", err)
	}
	if result.Status != AnalysisStatusOK || result.ErrorCount != 0 {
		t.Fatalf("ValidateFile() result = %+v, want ok", result)
	}

	if _, err = ValidateFile(filepath.Join(t.TempDir(), "missing.xml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ValidateFile() missing file error = %v, want os.ErrNotExist", err)
	}
}
