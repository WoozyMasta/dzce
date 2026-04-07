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

	// StageValidate marks XML shape and primitive validation diagnostics.
	StageValidate lint.Stage = "validate"

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
	// CodeValidateMissingRequiredAttr reports missing required XML attributes.
	CodeValidateMissingRequiredAttr lint.Code = 2001

	// CodeValidateEmptyRequiredAttr reports empty required XML attributes.
	CodeValidateEmptyRequiredAttr lint.Code = 2002

	// CodeValidateInvalidBool reports invalid boolean-like values.
	CodeValidateInvalidBool lint.Code = 2003

	// CodeValidateInvalidIntRange reports invalid integer ranges.
	CodeValidateInvalidIntRange lint.Code = 2004

	// CodeValidateUnknownEnum reports unsupported enum values.
	CodeValidateUnknownEnum lint.Code = 2005
)

const (
	// CodeTypesDuplicateName reports duplicate type names in one types.xml.
	CodeTypesDuplicateName lint.Code = 3001

	// CodeTypesNominalNegative reports negative nominal in types.xml.
	CodeTypesNominalNegative lint.Code = 3002

	// CodeTypesMinGreaterThanNominal reports min larger than nominal.
	CodeTypesMinGreaterThanNominal lint.Code = 3003

	// CodeEventsInvalidLimitWindow reports potentially invalid min/max windows.
	CodeEventsInvalidLimitWindow lint.Code = 3004

	// CodeSpawnableDuplicateChild reports duplicate child entries.
	CodeSpawnableDuplicateChild lint.Code = 3005

	// CodeGlobalsInvalidTypeTag reports unsupported globals type tag.
	CodeGlobalsInvalidTypeTag lint.Code = 3102

	// CodeGlobalsValueTypeMismatch reports globals value/type mismatch.
	CodeGlobalsValueTypeMismatch lint.Code = 3103

	// CodeGlobalsOutOfRange reports globals value outside known range.
	CodeGlobalsOutOfRange lint.Code = 3104

	// CodeEconomyCoreDuplicateDefaultName reports duplicate default name.
	CodeEconomyCoreDuplicateDefaultName lint.Code = 3202

	// CodeEconomyCoreDefaultInvalidBool reports invalid bool default value.
	CodeEconomyCoreDefaultInvalidBool lint.Code = 3203

	// CodeEconomyCoreDefaultOutOfRange reports invalid numeric default value.
	CodeEconomyCoreDefaultOutOfRange lint.Code = 3204

	// CodeEconomyIncompleteSection reports incomplete economy section flags.
	CodeEconomyIncompleteSection lint.Code = 3301

	// CodeTypesFlagsIncomplete reports incomplete `flags` block in types.xml.
	CodeTypesFlagsIncomplete lint.Code = 3302

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

	// CodeRandomPresetsDuplicateName reports duplicate preset name in file.
	CodeRandomPresetsDuplicateName lint.Code = 3702

	// CodeRandomPresetsEmptyItems reports preset without items.
	CodeRandomPresetsEmptyItems lint.Code = 3703

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
