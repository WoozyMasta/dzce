// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "github.com/woozymasta/lintkit/lint"

const (
	// LintModule is stable lint module namespace for dzce rules.
	LintModule = "dzce"
)

const (
	// StageParse marks XML/JSON parsing diagnostics.
	StageParse lint.Stage = "parse"

	// StageSemantic marks CE semantic validation diagnostics.
	StageSemantic lint.Stage = "semantic"

	// StageCrossref marks cross-file reference diagnostics.
	StageCrossref lint.Stage = "crossref"

	// StageMerge marks include graph and merge diagnostics.
	StageMerge lint.Stage = "merge"
)

const (
	// CodeParseInvalidXML reports malformed XML payload.
	CodeParseInvalidXML lint.Code = 1001

	// CodeParseUnknownRoot reports unsupported XML root tag.
	CodeParseUnknownRoot lint.Code = 1002
)

const (
	// CodeGlobalsInvalidTypeTag reports unsupported globals type tag.
	CodeGlobalsInvalidTypeTag lint.Code = 3102

	// CodeGlobalsValueTypeMismatch reports globals value/type mismatch.
	CodeGlobalsValueTypeMismatch lint.Code = 3103

	// CodeGlobalsOutOfRange reports globals value outside known range.
	CodeGlobalsOutOfRange lint.Code = 3104

	// CodeEconomyCoreDuplicateDefaultName reports duplicate default name.
	CodeEconomyCoreDuplicateDefaultName lint.Code = 3202

	// CodeEconomyCoreDefaultOutOfRange reports invalid numeric default value.
	CodeEconomyCoreDefaultOutOfRange lint.Code = 3204

	// CodeEconomyIncompleteSection reports incomplete economy section flags.
	CodeEconomyIncompleteSection lint.Code = 3301

	// CodeTypesDuplicateName reports duplicate type names in one types.xml.
	CodeTypesDuplicateName lint.Code = 3001

	// CodeTypesNominalNegative reports negative nominal in types.xml.
	CodeTypesNominalNegative lint.Code = 3002

	// CodeTypesQuantityRange reports invalid quantmin/quantmax values.
	CodeTypesQuantityRange lint.Code = 3303

	// CodeTypesQuantRange is deprecated alias of CodeTypesQuantityRange.
	CodeTypesQuantRange lint.Code = CodeTypesQuantityRange

	// CodeSpawnableDamageRange reports invalid damage min/max range.
	CodeSpawnableDamageRange lint.Code = 3304

	// CodeSpawnableChanceRange reports invalid chance values.
	CodeSpawnableChanceRange lint.Code = 3305

	// CodeEventsDuplicateName reports duplicate event names in one file.
	CodeEventsDuplicateName lint.Code = 3601

	// CodeEventsFlagNonCanonical reports non-canonical event flags values.
	CodeEventsFlagNonCanonical lint.Code = 3604

	// CodeEventsUnknownPosition reports unsupported event position token.
	CodeEventsUnknownPosition lint.Code = 3605

	// CodeEventsUnknownLimit reports unsupported event limit token.
	CodeEventsUnknownLimit lint.Code = 3606

	// CodeCrossRefMissingType reports unresolved type reference.
	CodeCrossRefMissingType lint.Code = 4001

	// CodeCrossRefMissingEvent reports unresolved event reference.
	CodeCrossRefMissingEvent lint.Code = 4002

	// CodeMergeIncludeCycle reports economycore include cycle.
	CodeMergeIncludeCycle lint.Code = 4003

	// CodeMergeMissingIncludeFile reports missing include target file.
	CodeMergeMissingIncludeFile lint.Code = 4004

	// CodeMergeDuplicateTypeOverride reports duplicate type overrides.
	CodeMergeDuplicateTypeOverride lint.Code = 4005
)
