package dzce

import (
	"testing"

	"github.com/woozymasta/lintkit/lint"
)

func TestAnalyzeLintContentUnknownRoot(t *testing.T) {
	t.Parallel()

	content := []byte(`<?xml version="1.0"?><config></config>`)
	diagnostics := AnalyzeLintContent("unknown.xml", content)
	if !hasCode(diagnostics, CodeParseUnknownRoot) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeParseUnknownRoot))
	}
}

func TestAnalyzeLintContentGlobalsTypeAndRange(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<variables>
	<var name="TimeLogin" type="0" value="-1"/>
	<var name="RespawnLimit" type="8" value="10"/>
</variables>`,
	)

	diagnostics := AnalyzeLintContent("db/globals.xml", content)
	if !hasCode(diagnostics, CodeGlobalsOutOfRange) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeGlobalsOutOfRange))
	}
	if !hasCode(diagnostics, CodeGlobalsInvalidTypeTag) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeGlobalsInvalidTypeTag),
		)
	}
}

func TestAnalyzeLintContentEconomyCoreDuplicateDefaults(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<economycore>
	<defaults>
		<default name="backup_period" value="60"/>
		<default name="backup_period" value="120"/>
	</defaults>
</economycore>`,
	)

	diagnostics := AnalyzeLintContent("cfgeconomycore.xml", content)
	if !hasCode(diagnostics, CodeEconomyCoreDuplicateDefaultName) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeEconomyCoreDuplicateDefaultName),
		)
	}
}

func TestAnalyzeLintContentEconomyCoreRange(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<economycore>
	<defaults>
		<default name="backup_period" value="5"/>
		<default name="dyn_smin" value="10"/>
		<default name="dyn_smax" value="2"/>
	</defaults>
</economycore>`,
	)

	diagnostics := AnalyzeLintContent("cfgeconomycore.xml", content)
	if !hasCode(diagnostics, CodeEconomyCoreDefaultOutOfRange) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeEconomyCoreDefaultOutOfRange),
		)
	}
}

func TestAnalyzeLintContentEconomySectionFlags(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<economy>
	<dynamic init="1" load="1" respawn="1" save="1"/>
	<animals init="1" load="1" respawn="1" save="1"/>
	<zombies init="1" load="1" respawn="2" save="1"/>
	<vehicles init="1" load="1" respawn="1" save="1"/>
	<randoms init="1" load="1" respawn="1" save="1"/>
	<custom init="1" load="1" respawn="1" save="1"/>
	<building init="1" load="1" respawn="1" save="1"/>
</economy>`,
	)

	diagnostics := AnalyzeLintContent("db/economy.xml", content)
	if !hasCode(diagnostics, CodeEconomyIncompleteSection) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeEconomyIncompleteSection),
		)
	}
}

func TestAnalyzeLintContentTypesRules(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<types>
	<type name="BandageDressing">
		<nominal>-1</nominal>
		<quantmin>101</quantmin>
		<quantmax>1</quantmax>
	</type>
	<type name="BandageDressing">
		<nominal>1</nominal>
		<quantmin>50</quantmin>
		<quantmax>40</quantmax>
	</type>
</types>`,
	)

	diagnostics := AnalyzeLintContent("db/types.xml", content)
	if !hasCode(diagnostics, CodeTypesDuplicateName) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeTypesDuplicateName))
	}
	if !hasCode(diagnostics, CodeTypesNominalNegative) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeTypesNominalNegative))
	}
	if !hasCode(diagnostics, CodeTypesQuantityRange) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeTypesQuantityRange))
	}
}

func TestAnalyzeLintContentEventsRules(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<events>
	<event name="StaticHeliCrash" active="2">
		<position>invalid</position>
		<limit>invalid</limit>
		<flags deletable="1" init_random="2" remove_damaged="0"/>
	</event>
	<event name="StaticHeliCrash"/>
</events>`,
	)

	diagnostics := AnalyzeLintContent("db/events.xml", content)
	if !hasCode(diagnostics, CodeEventsDuplicateName) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeEventsDuplicateName))
	}
	if !hasCode(diagnostics, CodeEventsFlagNonCanonical) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeEventsFlagNonCanonical))
	}
	if !hasCode(diagnostics, CodeEventsUnknownPosition) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeEventsUnknownPosition))
	}
	if !hasCode(diagnostics, CodeEventsUnknownLimit) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeEventsUnknownLimit))
	}
}

func TestAnalyzeLintContentSpawnableRules(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<spawnabletypes>
	<type name="AmmoBox_556x45_20Rnd">
		<damage min="-0.1" max="1.2"/>
		<cargo chance="120">
			<item name="Rag" chance="-1"/>
		</cargo>
	</type>
</spawnabletypes>`,
	)

	diagnostics := AnalyzeLintContent("cfgspawnabletypes.xml", content)
	if !hasCode(diagnostics, CodeSpawnableDamageRange) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeSpawnableDamageRange))
	}
	if !hasCode(diagnostics, CodeSpawnableChanceRange) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeSpawnableChanceRange))
	}
}

func TestAnalyzeLintContentValidateRequiredAttrs(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<types>
	<type>
		<nominal>1</nominal>
	</type>
</types>`,
	)

	diagnostics := AnalyzeLintContent("db/types.xml", content)
	if !hasCode(diagnostics, CodeValidateMissingRequiredAttr) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeValidateMissingRequiredAttr),
		)
	}
}

func TestAnalyzeLintContentNewSemanticRules(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<types>
	<type name="ExampleType">
		<nominal>1</nominal>
		<min>2</min>
		<flags count_in_cargo="1"/>
	</type>
</types>`,
	)

	diagnostics := AnalyzeLintContent("db/types.xml", content)
	if !hasCode(diagnostics, CodeTypesMinGreaterThanNominal) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeTypesMinGreaterThanNominal),
		)
	}
	if !hasCode(diagnostics, CodeTypesFlagsIncomplete) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeTypesFlagsIncomplete),
		)
	}
}

func TestAnalyzeLintContentEconomyCoreUnknownAndBool(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<economycore>
	<defaults>
		<default name="log_ce_startup" value="maybe"/>
	</defaults>
</economycore>`,
	)

	diagnostics := AnalyzeLintContent("cfgeconomycore.xml", content)
	if !hasCode(diagnostics, CodeEconomyCoreDefaultInvalidBool) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeEconomyCoreDefaultInvalidBool),
		)
	}
}

func TestAnalyzeLintContentEventsLimitWindow(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<events>
	<event name="E1">
		<nominal>1</nominal>
		<min>2</min>
		<max>1</max>
	</event>
</events>`,
	)

	diagnostics := AnalyzeLintContent("db/events.xml", content)
	if !hasCode(diagnostics, CodeEventsInvalidLimitWindow) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeEventsInvalidLimitWindow),
		)
	}
}

func TestAnalyzeLintContentSpawnableDuplicateChild(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<spawnabletypes>
	<type name="T1">
		<cargo chance="1">
			<item name="Rag" chance="1"/>
			<item name="Rag" chance="1"/>
		</cargo>
	</type>
</spawnabletypes>`,
	)

	diagnostics := AnalyzeLintContent("cfgspawnabletypes.xml", content)
	if !hasCode(diagnostics, CodeSpawnableDuplicateChild) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeSpawnableDuplicateChild),
		)
	}
}

func TestAnalyzeLintContentValidateCodes(t *testing.T) {
	t.Parallel()

	eventsContent := []byte(
		`<?xml version="1.0"?>
<events>
	<event name="" active="2">
		<position>bad</position>
		<limit>bad</limit>
		<flags deletable="2"/>
	</event>
</events>`,
	)
	eventsDiagnostics := AnalyzeLintContent("db/events.xml", eventsContent)
	if !hasCode(eventsDiagnostics, CodeValidateEmptyRequiredAttr) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeValidateEmptyRequiredAttr),
		)
	}
	if !hasCode(eventsDiagnostics, CodeValidateInvalidBool) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeValidateInvalidBool),
		)
	}
	if !hasCode(eventsDiagnostics, CodeValidateUnknownEnum) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeValidateUnknownEnum),
		)
	}

	typesContent := []byte(
		`<?xml version="1.0"?>
<types>
	<type name="A">
		<nominal>-1</nominal>
	</type>
</types>`,
	)
	typesDiagnostics := AnalyzeLintContent("db/types.xml", typesContent)
	if !hasCode(typesDiagnostics, CodeValidateInvalidIntRange) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeValidateInvalidIntRange),
		)
	}
}

func TestAnalyzeLintContentRandomPresetsRules(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<randompresets>
	<cargo name="p1">
		<item name="Rag" chance="1"/>
	</cargo>
	<cargo name="p1"/>
	<attachments name="a1"/>
</randompresets>`,
	)

	diagnostics := AnalyzeLintContent("cfgrandompresets.xml", content)
	if !hasCode(diagnostics, CodeRandomPresetsDuplicateName) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeRandomPresetsDuplicateName),
		)
	}
	if !hasCode(diagnostics, CodeRandomPresetsEmptyItems) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeRandomPresetsEmptyItems),
		)
	}
}

// hasCode checks whether diagnostic list contains one numeric code.
func hasCode(diagnostics []lint.Diagnostic, code lint.Code) bool {
	needle := codeToken(code)
	for index := range diagnostics {
		if diagnostics[index].Code == needle {
			return true
		}
	}

	return false
}

// codeToken renders public code token for dzce diagnostics.
func codeToken(code lint.Code) string {
	return lint.ApplyCodePrefix("DZCE", code)
}
