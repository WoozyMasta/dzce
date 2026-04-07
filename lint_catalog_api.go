// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "github.com/woozymasta/lintkit/lint"

var diagnosticCodeCatalogHandle = lint.NewCodeCatalogHandle(
	lint.CodeCatalogConfig{
		Module:            LintModule,
		CodePrefix:        "DZCE",
		ModuleName:        "DayZ CE",
		ModuleDescription: "Lint rules for DayZ central economy configuration files.",
		ScopeDescriptions: map[lint.Stage]string{
			StageParse:    "Parse diagnostics for CE input payloads.",
			StageSemantic: "Semantic validation diagnostics for decoded CE models.",
			StageCrossref: "Cross-file reference diagnostics for merged CE trees.",
			StageMerge:    "Include graph and merge diagnostics for economycore trees.",
		},
	},
	diagnosticCatalog,
)

// getDiagnosticCodeCatalog returns lazy-initialized diagnostics catalog.
func getDiagnosticCodeCatalog() (lint.CodeCatalog, error) {
	return diagnosticCodeCatalogHandle.Catalog()
}

// DiagnosticRuleSpec converts one diagnostic spec into lint rule metadata.
func DiagnosticRuleSpec(spec lint.CodeSpec) (lint.RuleSpec, error) {
	return diagnosticCodeCatalogHandle.RuleSpec(spec)
}

// LintRuleID returns lint rule ID mapped from stable dzce diagnostic code.
func LintRuleID(code lint.Code) string {
	return diagnosticCodeCatalogHandle.RuleIDOrUnknown(code)
}

// DiagnosticCatalog returns stable diagnostics metadata list.
func DiagnosticCatalog() []lint.CodeSpec {
	return diagnosticCodeCatalogHandle.CodeSpecs()
}

// DiagnosticByCode returns diagnostic metadata for one code token.
func DiagnosticByCode(code lint.Code) (lint.CodeSpec, bool) {
	return diagnosticCodeCatalogHandle.ByCode(code)
}

// LintRuleSpecs returns deterministic lint rule specs from diagnostics catalog.
func LintRuleSpecs() []lint.RuleSpec {
	return diagnosticCodeCatalogHandle.RuleSpecs()
}
