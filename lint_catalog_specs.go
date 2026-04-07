// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "github.com/woozymasta/lintkit/lint"

// withDescription attaches optional documentation text to one catalog spec.
func withDescription(spec lint.CodeSpec, description string) lint.CodeSpec {
	spec.Description = description
	return spec
}

// diagnosticCatalog stores stable dzce diagnostics metadata table.
var diagnosticCatalog = []lint.CodeSpec{
	withDescription(
		lint.ErrorCodeSpec(
			CodeParseInvalidXML,
			StageParse,
			"xml decode failed",
		),
		"CE XML file is malformed and cannot be parsed. Check broken tags, invalid attributes, and XML syntax.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParseUnknownRoot,
			StageParse,
			"unsupported xml root",
		),
		"XML is valid, but root tag is not a supported CE file model for this check set. Use a supported CE root or exclude this file from CE lint input.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateMissingRequiredAttr,
			StageValidate,
			"missing required attribute",
		),
		"Required XML attribute is missing. Add the attribute so CE can interpret this element deterministically.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateEmptyRequiredAttr,
			StageValidate,
			"empty required attribute",
		),
		"Required XML attribute is present but empty. Provide a non-empty value.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeValidateInvalidBool,
			StageValidate,
			"invalid bool value",
		),
		"Boolean-like field uses a non-canonical value. Use 0 or 1 for CE XML boolean fields.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeValidateInvalidIntRange,
			StageValidate,
			"invalid integer range",
		),
		"Integer value is outside expected range for this field.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeValidateUnknownEnum,
			StageValidate,
			"unknown enum value",
		),
		"Enum field contains unsupported token for this CE context.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeTypesDuplicateName,
			StageSemantic,
			"duplicate type name",
		),
		"`types.xml` contains multiple `<type>` entries with the same name. Keep one canonical CE type definition, or verify override order is intentional.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeTypesNominalNegative,
			StageSemantic,
			"type nominal is negative",
		),
		"`types.xml` has `<nominal>` below zero. In CE this value is usually expected to be 0 or greater.",
	),
	withDescription(
		lint.NoticeCodeSpec(
			CodeTypesMinGreaterThanNominal,
			StageSemantic,
			"type min is greater than nominal",
		),
		"`types.xml` has `<min>` larger than `<nominal>`. This can be intentional in some setups, but often indicates inconsistent balancing values.",
	),
	withDescription(
		lint.NoticeCodeSpec(
			CodeEventsInvalidLimitWindow,
			StageSemantic,
			"event limit window looks inconsistent",
		),
		"`events.xml` has potentially inconsistent numeric window values (for example min > nominal or max < min). Review this event configuration.",
	),
	withDescription(
		lint.NoticeCodeSpec(
			CodeSpawnableDuplicateChild,
			StageSemantic,
			"duplicate spawnable child entry",
		),
		"`cfgspawnabletypes.xml` contains duplicate child item names in the same parent block. Remove duplicates unless intentional weighted duplication is required.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeGlobalsInvalidTypeTag,
			StageSemantic,
			"globals var type tag is invalid",
		),
		"`globals.xml` var@type uses an unsupported CE type tag. Allowed values are 0, 1, 2.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeGlobalsValueTypeMismatch,
			StageSemantic,
			"globals value type mismatch",
		),
		"`globals.xml` var@value does not match declared var@type. Fix CE value format so it matches the selected type tag.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeGlobalsOutOfRange,
			StageSemantic,
			"globals value is out of range",
		),
		"`globals.xml` value is outside recommended CE range for this variable. Review gameplay impact and adjust if not intentional.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeEconomyCoreDuplicateDefaultName,
			StageSemantic,
			"duplicate economycore default name",
		),
		"`cfgeconomycore.xml` `<defaults>` contains duplicate default@name keys. Keep one value per key to avoid ambiguous CE runtime defaults.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeEconomyCoreDefaultInvalidBool,
			StageSemantic,
			"economycore bool default is invalid",
		),
		"`cfgeconomycore.xml` bool-like default key uses non-bool value. Use true/false or 0/1.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeEconomyCoreDefaultOutOfRange,
			StageSemantic,
			"economycore default is out of range",
		),
		"A numeric default in `cfgeconomycore.xml` is outside expected CE limits. This may lead to unstable CE behavior.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeEconomyIncompleteSection,
			StageSemantic,
			"economy section flags are invalid",
		),
		"`economy.xml` section is missing required CE flags or uses invalid values. Each section should define init/load/respawn/save as 0 or 1.",
	),
	withDescription(
		lint.NoticeCodeSpec(
			CodeTypesFlagsIncomplete,
			StageSemantic,
			"type flags block is incomplete",
		),
		"`types.xml` `<flags>` block does not define all commonly paired attributes. This can cause implicit inheritance/merge side effects.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeTypesQuantityRange,
			StageSemantic,
			"invalid type quantity range",
		),
		"`types.xml` quantity range (quantmin/quantmax) is invalid for CE. Allowed values are -1 or 0..100, and quantmin must not exceed quantmax.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeSpawnableDamageRange,
			StageSemantic,
			"invalid spawnable damage range",
		),
		"`cfgspawnabletypes.xml` damage range is invalid for CE spawn rules. Use values in 0..1 and ensure min is not greater than max.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeSpawnableChanceRange,
			StageSemantic,
			"invalid spawnable chance range",
		),
		"`cfgspawnabletypes.xml` chance values are inconsistent for CE spawn rules. Use one mode consistently: normalized 0..1 or percent 0..100.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeEventsDuplicateName,
			StageSemantic,
			"duplicate event name",
		),
		"`events.xml` contains duplicate `<event>` names. Keep CE event names unique to avoid ambiguous merge and spawn behavior.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeEventsFlagNonCanonical,
			StageSemantic,
			"event flags are not 0/1",
		),
		"`events.xml` active/flag attributes use non-canonical CE values. Use 0 or 1.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeEventsUnknownPosition,
			StageSemantic,
			"unsupported event position",
		),
		"`events.xml` position uses an unsupported CE value. Supported values: fixed, player, uniform.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeEventsUnknownLimit,
			StageSemantic,
			"unsupported event limit",
		),
		"`events.xml` limit uses an unsupported CE value. Supported values: child, custom, mixed, parent.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeRandomPresetsDuplicateName,
			StageSemantic,
			"duplicate random preset name",
		),
		"`cfgrandompresets.xml` contains duplicate preset names in one section. Keep preset names unique inside cargo and attachments blocks.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeRandomPresetsEmptyItems,
			StageSemantic,
			"random preset has no items",
		),
		"`cfgrandompresets.xml` preset is declared without `<item>` entries and has no effect during generation.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeCrossRefMissingType,
			StageCrossref,
			"missing type reference",
		),
		"A merged CE event child references a type missing in final merged `types.xml`. Add the missing type definition or fix the reference.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeCrossRefMissingEvent,
			StageCrossref,
			"missing event reference",
		),
		"A merged `events.xml` event.secondary references an event name missing in final merged `events.xml`.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeMergeIncludeCycle,
			StageMerge,
			"include cycle detected",
		),
		"`cfgeconomycore.xml` include graph has a recursive cycle. Break the cycle to get deterministic CE merge order.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeMergeMissingIncludeFile,
			StageMerge,
			"include file not found",
		),
		"`cfgeconomycore.xml` references an include file that cannot be found at resolved path. Fix CE include folder/name/type or restore the file.",
	),
	withDescription(
		lint.NoticeCodeSpec(
			CodeMergeDuplicateTypeOverride,
			StageMerge,
			"duplicate type override across includes",
		),
		"The same CE type name is defined in multiple included `types.xml` files. This can be intentional override behavior, but verify final include priority/order.",
	),
}
