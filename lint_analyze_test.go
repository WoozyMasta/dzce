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

func TestAnalyzeLintContentJSONDecodeError(t *testing.T) {
	t.Parallel()

	diagnostics := AnalyzeLintContent("cfggameplay.json", []byte(`{"version":`))
	if !hasCode(diagnostics, CodeParseInvalidJSON) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeParseInvalidJSON))
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
	if hasCode(diagnostics, CodeValidateInvalidBool) {
		t.Fatalf("economy flag diagnostic must not be duplicated by generic bool code: %#v", diagnostics)
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
	if hasCode(diagnostics, CodeValidateInvalidIntRange) {
		t.Fatalf("quantity diagnostics must not be duplicated by generic range code: %#v", diagnostics)
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

func TestAnalyzeLintContentSpawnableChanceModes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		content         string
		wantDiagnostics int
	}{
		{
			name: "probability values do not need to sum to one",
			content: `<spawnabletypes><type name="T1">
			<attachments chance="0.8">
				<item name="A" chance="0.3"/><item name="B" chance="0.2"/>
			</attachments>
		</type></spawnabletypes>`,
		},
		{
			name: "percentage weights totaling one hundred",
			content: `<spawnabletypes><type name="AKM">
			<attachments chance="1">
				<item name="A" chance="40"/><item name="B" chance="30"/>
				<item name="C" chance="20"/><item name="D" chance="5"/>
				<item name="E" chance="5"/>
			</attachments>
		</type></spawnabletypes>`,
		},
		{
			name: "mixed scales totaling one hundred are diagnosed",
			content: `<spawnabletypes><type name="T1">
			<cargo chance="0.5">
				<item name="A" chance="99"/><item name="B" chance="0.5"/>
				<item name="C" chance="0.5"/>
			</cargo>
		</type></spawnabletypes>`,
			wantDiagnostics: 1,
		},
		{
			name: "percentage weights with a minority fractional value",
			content: `<spawnabletypes><type name="T1">
			<cargo chance="0.5">
				<item name="A" chance="40"/><item name="B" chance="30"/>
				<item name="C" chance="29.5"/><item name="D" chance="0.5"/>
			</cargo>
		</type></spawnabletypes>`,
		},
		{
			name: "percentage weights with a non-one-hundred total",
			content: `<spawnabletypes><type name="T1">
			<cargo chance="0.5">
				<item name="A" chance="40"/><item name="B" chance="50"/>
			</cargo>
		</type></spawnabletypes>`,
			wantDiagnostics: 1,
		},
		{
			name: "percentage weight above one hundred",
			content: `<spawnabletypes><type name="T1">
			<cargo chance="0.5"><item name="A" chance="101"/></cargo>
		</type></spawnabletypes>`,
			wantDiagnostics: 1,
		},
		{
			name: "incomplete percentage weights skip total check",
			content: `<spawnabletypes><type name="T1">
			<cargo chance="0.5">
				<item name="A" chance="40"/><item name="B"/>
			</cargo>
		</type></spawnabletypes>`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			diagnostics := AnalyzeLintContent(
				"cfgspawnabletypes.xml",
				[]byte(test.content),
			)
			if got := countCode(diagnostics, CodeSpawnableChanceRange); got != test.wantDiagnostics {
				t.Fatalf(
					"chance diagnostics = %d, want %d; diagnostics: %#v",
					got,
					test.wantDiagnostics,
					diagnostics,
				)
			}
		})
	}
}

func TestAnalyzeLintContentSpawnableItemChanceRejectsNonFiniteValues(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<spawnabletypes>
	<type name="T1">
		<attachments chance="1">
			<item name="A" chance="NaN"/>
		</attachments>
	</type>
</spawnabletypes>`,
	)

	diagnostics := AnalyzeLintContent("cfgspawnabletypes.xml", content)
	if !hasCode(diagnostics, CodeSpawnableChanceRange) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeSpawnableChanceRange))
	}
}

func TestAnalyzeLintContentNonFiniteFloatValues(t *testing.T) {
	t.Parallel()

	spawnable := AnalyzeLintContent(
		"cfgspawnabletypes.xml",
		[]byte(`<spawnabletypes><type name="T1"><damage min="NaN" max="1"/><cargo chance="NaN"/></type></spawnabletypes>`),
	)
	if !hasCode(spawnable, CodeSpawnableDamageRange) {
		t.Fatalf("expected non-finite damage diagnostic, got %#v", spawnable)
	}
	if !hasCode(spawnable, CodeSpawnableChanceRange) {
		t.Fatalf("expected non-finite chance diagnostic, got %#v", spawnable)
	}

	globals := AnalyzeLintContent(
		"db/globals.xml",
		[]byte(`<variables><var name="LootDamageMin" type="1" value="NaN"/></variables>`),
	)
	if !hasCode(globals, CodeGlobalsValueTypeMismatch) {
		t.Fatalf("expected non-finite globals diagnostic, got %#v", globals)
	}

	defaults := AnalyzeLintContent(
		"cfgeconomycore.xml",
		[]byte(`<economycore><defaults><default name="dyn_radius" value="NaN"/></defaults></economycore>`),
	)
	if !hasCode(defaults, CodeEconomyCoreDefaultOutOfRange) {
		t.Fatalf("expected non-finite default diagnostic, got %#v", defaults)
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

func TestAnalyzeLintContentEventsZeroSentinels(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<events>
	<event name="AnimalBear">
		<nominal>0</nominal>
		<min>5</min>
		<max>8</max>
	</event>
	<event name="StaticContainerLocked">
		<nominal>13</nominal>
		<min>11</min>
		<max>0</max>
	</event>
</events>`,
	)

	diagnostics := AnalyzeLintContent("db/events.xml", content)
	if hasCode(diagnostics, CodeEventsInvalidLimitWindow) {
		t.Fatalf("zero sentinel values should not trigger limit-window diagnostics: %#v", diagnostics)
	}
}

func TestAnalyzeLintContentTypesZeroNominal(t *testing.T) {
	t.Parallel()

	content := []byte(
		`<?xml version="1.0"?>
<types>
	<type name="ExampleType">
		<nominal>0</nominal>
		<min>1</min>
	</type>
</types>`,
	)

	diagnostics := AnalyzeLintContent("db/types.xml", content)
	if hasCode(diagnostics, CodeTypesMinGreaterThanNominal) {
		t.Fatalf("min > nominal should be ignored when nominal is zero: %#v", diagnostics)
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
	if !hasCode(eventsDiagnostics, CodeEventsFlagNonCanonical) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeEventsFlagNonCanonical),
		)
	}
	if !hasCode(eventsDiagnostics, CodeEventsUnknownPosition) ||
		!hasCode(eventsDiagnostics, CodeEventsUnknownLimit) {
		t.Fatalf(
			"expected specialized enum diagnostics in diagnostics",
		)
	}
	if hasCode(eventsDiagnostics, CodeValidateInvalidBool) ||
		hasCode(eventsDiagnostics, CodeValidateUnknownEnum) {
		t.Fatalf("specific event diagnostics must not be duplicated by generic codes: %#v", eventsDiagnostics)
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
	if !hasCode(typesDiagnostics, CodeTypesNominalNegative) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeTypesNominalNegative),
		)
	}
	if hasCode(typesDiagnostics, CodeValidateInvalidIntRange) {
		t.Fatalf("nominal diagnostic must not be duplicated by generic range code: %#v", typesDiagnostics)
	}

	minDiagnostics := AnalyzeLintContent(
		"db/types.xml",
		[]byte(`<types><type name="A"><min>-1</min></type></types>`),
	)
	if !hasCode(minDiagnostics, CodeValidateInvalidIntRange) {
		t.Fatalf("expected generic range code for negative min: %#v", minDiagnostics)
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

func TestAnalyzeLintContentSchemaConstraints(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		content  string
		wantCode lint.Code
	}{
		{
			name:     "gameplay numeric minimum",
			path:     "cfggameplay.json",
			content:  `{"PlayerData":{"MovementData":{"timeToSprint":0}}}`,
			wantCode: CodeValidateNumericRange,
		},
		{
			name:     "gameplay enum",
			path:     "cfggameplay.json",
			content:  `{"PlayerData":{"WeaponObstructionData":{"staticMode":3}}}`,
			wantCode: CodeValidateUnknownEnum,
		},
		{
			name:     "gameplay fixed array length",
			path:     "cfggameplay.json",
			content:  `{"WorldsData":{"environmentMinTemps":[1]}}`,
			wantCode: CodeValidateInvalidLength,
		},
		{
			name:     "underground trigger vector length",
			path:     "cfgundergroundtriggers.json",
			content:  `{"Triggers":[{"Position":[1,2]}]}`,
			wantCode: CodeValidateInvalidLength,
		},
		{
			name:     "underground accommodation range",
			path:     "cfgundergroundtriggers.json",
			content:  `{"Triggers":[{"EyeAccommodation":1.1}]}`,
			wantCode: CodeValidateNumericRange,
		},
		{
			name:     "effect area position length",
			path:     "cfgeffectarea.json",
			content:  `{"Areas":[{"Data":{"Pos":[1,2]}}]}`,
			wantCode: CodeValidateInvalidLength,
		},
		{
			name:     "bunker setting range",
			path:     "cfgbunkerbroadcast.json",
			content:  `{"bunkerSettings":{"numberOfActive":31}}`,
			wantCode: CodeValidateInvalidIntRange,
		},
		{
			name:     "object spawner position length",
			path:     "object-spawner.json",
			content:  `{"Objects":[{"pos":[1,2]}]}`,
			wantCode: CodeValidateInvalidLength,
		},
		{
			name:     "bunker broadcast time length",
			path:     "cfgbunkerbroadcast.json",
			content:  `{"schedulerSettings":{"periods":[{"startTime":"700","endTime":"1200"}]}}`,
			wantCode: CodeValidateInvalidLength,
		},
		{
			name:     "gear preset minimum weight",
			path:     "gameplay-gear-presets.json",
			content:  `[{"spawnWeight":0}]`,
			wantCode: CodeValidateInvalidIntRange,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			diagnostics := AnalyzeLintContent(test.path, []byte(test.content))
			if !hasCode(diagnostics, test.wantCode) {
				t.Fatalf("expected %s in diagnostics, got %#v", codeToken(test.wantCode), diagnostics)
			}
		})
	}
}

func TestAnalyzeLintContentEconomyCoreStringEnums(t *testing.T) {
	t.Parallel()

	diagnostics := AnalyzeLintContent(
		"cfgeconomycore.xml",
		[]byte(`<economycore><classes><rootclass name="ItemBase" reportMemoryLOD="sometimes" act="zombie"/></classes></economycore>`),
	)
	if got := countCode(diagnostics, CodeValidateUnknownEnum); got != 2 {
		t.Fatalf("economycore enum diagnostics = %d, want 2; diagnostics: %#v", got, diagnostics)
	}
}

func TestAnalyzeLintContentGearPresetBounds(t *testing.T) {
	t.Parallel()

	diagnostics := AnalyzeLintContent(
		"gameplay-gear-presets.json",
		[]byte(`[{
  "attachmentSlotItemSets": [{
    "discreteItemSets": [{
      "quickBarSlot": -2,
      "complexChildrenTypes": [{
        "quickBarSlot": -2
      }]
    }]
  }],
  "discreteUnsortedItemSets": [{
    "spawnWeight": 0
  }]
}]`),
	)
	if got := countCode(diagnostics, CodeValidateInvalidIntRange); got != 3 {
		t.Fatalf("gear preset bound diagnostics = %d, want 3; diagnostics: %#v", got, diagnostics)
	}
}

func TestAnalyzeLintContentObjectSpawnerOptionalYPR(t *testing.T) {
	t.Parallel()

	for _, ypr := range []string{"[]", "[1,2,3]"} {
		diagnostics := AnalyzeLintContent(
			"object-spawner.json",
			[]byte(`{"Objects":[{"pos":[1,2,3],"ypr":`+ypr+`}]}`),
		)
		if hasCode(diagnostics, CodeValidateInvalidLength) {
			t.Errorf("valid ypr %s produced length diagnostic: %#v", ypr, diagnostics)
		}
	}

	diagnostics := AnalyzeLintContent(
		"object-spawner.json",
		[]byte(`{"Objects":[{"pos":[1,2,3],"ypr":[1,2]}]}`),
	)
	if !hasCode(diagnostics, CodeValidateInvalidLength) {
		t.Fatalf("expected invalid optional YPR vector length: %#v", diagnostics)
	}
}

func TestAnalyzeLintContentSchemaConstraintsInXML(t *testing.T) {
	t.Parallel()

	messages := AnalyzeLintContent(
		"db/messages.xml",
		[]byte(`<messages>
	<message>
		<deadline>-1</deadline>
		<shutdown>2</shutdown>
		<delay>-1</delay>
		<repeat>-2</repeat>
		<onconnect>3</onconnect>
</message></messages>`),
	)
	if got := countCode(messages, CodeValidateInvalidIntRange); got != 3 {
		t.Fatalf("message range diagnostics = %d, want 3; diagnostics: %#v", got, messages)
	}
	if got := countCode(messages, CodeValidateUnknownEnum); got != 2 {
		t.Fatalf("message enum diagnostics = %d, want 2; diagnostics: %#v", got, messages)
	}

	weather := AnalyzeLintContent(
		"cfgweather.xml",
		[]byte(`<weather reset="2" enable="1"></weather>`),
	)
	if !hasCode(weather, CodeParseInvalidXML) {
		t.Fatalf("expected decoder to reject unsupported weather toggle: %#v", weather)
	}
}

func TestAnalyzeLintContentRandomPresetChanceConstraints(t *testing.T) {
	t.Parallel()

	invalid := AnalyzeLintContent(
		"cfgrandompresets.xml",
		[]byte(`<randompresets><cargo name="p" chance="1.1"><item name="A" chance="-0.1"/></cargo></randompresets>`),
	)
	if got := countCode(invalid, CodeSpawnableChanceRange); got != 2 {
		t.Fatalf("random preset chance diagnostics = %d, want 2; diagnostics: %#v", got, invalid)
	}

	validWeight := AnalyzeLintContent(
		"cfgrandompresets.xml",
		[]byte(`<randompresets><cargo name="p"><item name="A" chance="1.2"/></cargo></randompresets>`),
	)
	if hasCode(validWeight, CodeSpawnableChanceRange) {
		t.Fatalf("unbounded candidate chance was rejected: %#v", validWeight)
	}

	nested := AnalyzeLintContent(
		"cfgspawnabletypes.xml",
		[]byte(`<spawnabletypes>
	<type name="T"><cargo chance="0.5"><item name="A" chance="0.5"><cargo chance="1.1"><item name="B" chance="-0.1"/></cargo></item></cargo></type>
</spawnabletypes>`),
	)
	if got := countCode(nested, CodeSpawnableChanceRange); got != 2 {
		t.Fatalf("nested chance diagnostics = %d, want 2; diagnostics: %#v", got, nested)
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

// countCode returns the number of diagnostics with one public code token.
func countCode(diagnostics []lint.Diagnostic, code lint.Code) int {
	needle := codeToken(code)
	count := 0
	for index := range diagnostics {
		if diagnostics[index].Code == needle {
			count++
		}
	}

	return count
}

// codeToken renders public code token for dzce diagnostics.
func codeToken(code lint.Code) string {
	return lint.ApplyCodePrefix("DZCE", code)
}
