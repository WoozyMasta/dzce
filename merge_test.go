// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestKindFromEconomyCoreType(t *testing.T) {
	tests := []struct {
		value string
		want  Kind
	}{
		{value: "types", want: KindTypes},
		{value: "events", want: KindEvents},
		{value: "economy", want: KindEconomy},
		{value: "globals", want: KindGlobals},
		{value: "messages", want: KindMessages},
		{value: "spawnabletypes", want: KindSpawnableTypes},
		{value: "cfgspawnabletypes", want: KindSpawnableTypes},
		{value: "economycore", want: KindEconomyCore},
		{value: "cfgeconomycore", want: KindEconomyCore},
		{value: "unknown", want: KindUnknown},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.value, func(t *testing.T) {
			if got := KindFromEconomyCoreType(tc.value); got != tc.want {
				t.Fatalf(
					"KindFromEconomyCoreType(%q) = %q, want %q",
					tc.value,
					got,
					tc.want,
				)
			}
		})
	}
}

func TestLoadMergedEconomyCore(t *testing.T) {
	path := filepath.Join("testdata", "merge", "main", "cfgeconomycore.xml")
	result, err := LoadMergedEconomyCore(path)
	if err != nil {
		t.Fatalf("LoadMergedEconomyCore(%s) error: %v", path, err)
	}

	if len(result.Sources) != 5 {
		t.Fatalf("sources count = %d, want 5", len(result.Sources))
	}

	rawTypes, ok := result.Get(KindTypes)
	if !ok {
		t.Fatal("merged types are missing")
	}

	typesDoc, ok := rawTypes.(*TypesFile)
	if !ok {
		t.Fatalf("types value type = %T, want *TypesFile", rawTypes)
	}

	if len(typesDoc.Types) != 2 {
		t.Fatalf("merged types count = %d, want 2", len(typesDoc.Types))
	}

	rawEvents, ok := result.Get(KindEvents)
	if !ok {
		t.Fatal("merged events are missing")
	}

	eventsDoc, ok := rawEvents.(*EventsFile)
	if !ok {
		t.Fatalf("events value type = %T, want *EventsFile", rawEvents)
	}

	if len(eventsDoc.Events) != 1 {
		t.Fatalf("merged events count = %d, want 1", len(eventsDoc.Events))
	}

	rawMessages, ok := result.Get(KindMessages)
	if !ok {
		t.Fatal("merged messages are missing")
	}

	messagesDoc, ok := rawMessages.(*MessagesFile)
	if !ok {
		t.Fatalf("messages value type = %T, want *MessagesFile", rawMessages)
	}

	if len(messagesDoc.Messages) != 2 {
		t.Fatalf("merged messages count = %d, want 2", len(messagesDoc.Messages))
	}
}

func TestLoadMergedEconomyCoreIncludeOrderPriority(t *testing.T) {
	root := t.TempDir()
	corePath := filepath.Join(root, "cfgeconomycore.xml")
	firstPath := filepath.Join(root, "first_types.xml")
	secondPath := filepath.Join(root, "second_types.xml")

	corePayload := `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<economycore>
    <ce folder=".">
        <file name="first_types.xml" type="types" />
        <file name="second_types.xml" type="types" />
    </ce>
</economycore>`
	firstPayload := `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<types>
    <type name="ItemA">
        <nominal>1</nominal>
        <lifetime>300</lifetime>
        <restock>60</restock>
    </type>
</types>`
	secondPayload := `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<types>
    <type name="ItemA">
        <nominal>9</nominal>
    </type>
</types>`

	if err := os.WriteFile(corePath, []byte(corePayload), 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", corePath, err)
	}

	if err := os.WriteFile(firstPath, []byte(firstPayload), 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", firstPath, err)
	}

	if err := os.WriteFile(secondPath, []byte(secondPayload), 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", secondPath, err)
	}

	result, err := LoadMergedEconomyCore(corePath)
	if err != nil {
		t.Fatalf("LoadMergedEconomyCore(%s) error: %v", corePath, err)
	}

	if len(result.Sources) != 2 {
		t.Fatalf("sources count = %d, want 2", len(result.Sources))
	}

	if result.Sources[0].Path != firstPath || result.Sources[1].Path != secondPath {
		t.Fatalf("sources order mismatch: %#v", result.Sources)
	}

	rawTypes, ok := result.Get(KindTypes)
	if !ok {
		t.Fatal("merged types are missing")
	}

	typesDoc, ok := rawTypes.(*TypesFile)
	if !ok {
		t.Fatalf("types value type = %T, want *TypesFile", rawTypes)
	}

	if len(typesDoc.Types) != 1 {
		t.Fatalf("merged types count = %d, want 1", len(typesDoc.Types))
	}

	if typesDoc.Types[0].Nominal == nil || *typesDoc.Types[0].Nominal != 9 {
		t.Fatalf("later include should override earlier nominal value")
	}

	if typesDoc.Types[0].Lifetime == nil || *typesDoc.Types[0].Lifetime != 300 {
		t.Fatalf("fields missing in later include should stay from earlier include")
	}
}

func TestLoadMergedEconomyCoreUnknownType(t *testing.T) {
	path := filepath.Join(
		"testdata",
		"merge",
		"negative_unknown_type",
		"cfgeconomycore.xml",
	)

	_, err := LoadMergedEconomyCore(path)
	if err == nil {
		t.Fatalf("LoadMergedEconomyCore(%s) expected error", path)
	}

	if !errors.Is(err, ErrUnknownFileKind) {
		t.Fatalf("error = %v, want ErrUnknownFileKind", err)
	}
}

func TestLoadMergedEconomyCoreUnknownTypeRelaxedByKind(t *testing.T) {
	root := t.TempDir()
	corePath := filepath.Join(root, "cfgeconomycore.xml")
	payloadPath := filepath.Join(root, "my_spawn_data.json")

	corePayload := `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<economycore>
    <ce folder=".">
        <file name="my_spawn_data.json" type="objectspawner" />
    </ce>
</economycore>`
	jsonPayload := `{
    "Objects": [
        {
            "name": "Land_Wall_Gate_FenR",
            "pos": [8406.5, 107.7, 12782.3],
            "ypr": [0.0, 0.0, 0.0]
        }
    ]
}`

	if err := os.WriteFile(corePath, []byte(corePayload), 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", corePath, err)
	}

	if err := os.WriteFile(payloadPath, []byte(jsonPayload), 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", payloadPath, err)
	}

	_, strictErr := LoadMergedEconomyCore(corePath)
	if strictErr == nil {
		t.Fatalf("LoadMergedEconomyCore(%s) expected strict unknown type error", corePath)
	}

	if !errors.Is(strictErr, ErrUnknownFileKind) {
		t.Fatalf("strict error = %v, want ErrUnknownFileKind", strictErr)
	}

	result, err := LoadMergedEconomyCoreWithOptions(corePath, MergeOptions{
		RelaxedIncludeTypes: true,
	})
	if err != nil {
		t.Fatalf("LoadMergedEconomyCoreWithOptions(%s) error: %v", corePath, err)
	}

	raw, ok := result.Get(KindObjectSpawner)
	if !ok {
		t.Fatal("merged object spawner payload is missing")
	}

	doc, ok := raw.(*ObjectSpawnerFile)
	if !ok {
		t.Fatalf("object spawner value type = %T, want *ObjectSpawnerFile", raw)
	}

	if len(doc.Objects) != 1 {
		t.Fatalf("object count = %d, want 1", len(doc.Objects))
	}
}

func TestLoadMergedEconomyCoreUnknownTypeRelaxedByFilename(t *testing.T) {
	root := t.TempDir()
	corePath := filepath.Join(root, "cfgeconomycore.xml")
	payloadPath := filepath.Join(root, "cfgundergroundtriggers.json")

	corePayload := `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<economycore>
    <ce folder=".">
        <file name="cfgundergroundtriggers.json" type="mystery" />
    </ce>
</economycore>`
	jsonPayload := `{
    "Triggers": [
        {
            "Position": [735.0, 533.7, 1229.1],
            "Orientation": [0.0, 0.0, 0.0],
            "Size": [15.0, 5.6, 10.8],
            "EyeAccommodation": 0.0
        }
    ]
}`

	if err := os.WriteFile(corePath, []byte(corePayload), 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", corePath, err)
	}

	if err := os.WriteFile(payloadPath, []byte(jsonPayload), 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", payloadPath, err)
	}

	result, err := LoadMergedEconomyCoreWithOptions(corePath, MergeOptions{
		RelaxedIncludeTypes: true,
	})
	if err != nil {
		t.Fatalf("LoadMergedEconomyCoreWithOptions(%s) error: %v", corePath, err)
	}

	raw, ok := result.Get(KindUndergroundTriggers)
	if !ok {
		t.Fatal("merged underground payload is missing")
	}

	doc, ok := raw.(*UndergroundTriggersFile)
	if !ok {
		t.Fatalf(
			"underground value type = %T, want *UndergroundTriggersFile",
			raw,
		)
	}

	if len(doc.Triggers) != 1 {
		t.Fatalf("trigger count = %d, want 1", len(doc.Triggers))
	}
}

func TestLoadMergedEconomyCoreCycle(t *testing.T) {
	root := t.TempDir()
	aPath := filepath.Join(root, "a.xml")
	bPath := filepath.Join(root, "b.xml")

	aPayload := `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<economycore>
    <ce folder=".">
        <file name="b.xml" type="economycore" />
    </ce>
</economycore>`
	bPayload := `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<economycore>
    <ce folder=".">
        <file name="a.xml" type="economycore" />
    </ce>
</economycore>`

	if err := os.WriteFile(aPath, []byte(aPayload), 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", aPath, err)
	}

	if err := os.WriteFile(bPath, []byte(bPayload), 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error: %v", bPath, err)
	}

	_, err := LoadMergedEconomyCore(aPath)
	if err == nil {
		t.Fatalf("LoadMergedEconomyCore(%s) expected cycle error", aPath)
	}

	if !errors.Is(err, ErrEconomyCoreCycle) {
		t.Fatalf("error = %v, want ErrEconomyCoreCycle", err)
	}
}

func TestMergeTypesValueOverridesByName(t *testing.T) {
	current := &TypesFile{
		Types: []TypeDef{
			{
				Name:    "AK101",
				Nominal: intPtr(2),
				Usages:  []NamedRef{{Name: "Military"}},
				Tags:    []NamedRef{{Name: "shelves"}},
			},
		},
	}
	incoming := &TypesFile{
		Types: []TypeDef{
			{
				Name:    "AK101",
				Nominal: intPtr(8),
				Usages:  []NamedRef{{Name: "Police"}},
			},
		},
	}

	mergedRaw, err := mergeTypesValue(current, incoming)
	if err != nil {
		t.Fatalf("mergeTypesValue() error: %v", err)
	}

	merged, ok := mergedRaw.(*TypesFile)
	if !ok {
		t.Fatalf("merged value type = %T, want *TypesFile", mergedRaw)
	}

	if len(merged.Types) != 1 {
		t.Fatalf("merged types count = %d, want 1", len(merged.Types))
	}

	if merged.Types[0].Nominal == nil || *merged.Types[0].Nominal != 8 {
		t.Fatalf("merged nominal mismatch")
	}

	if len(merged.Types[0].Usages) != 1 || merged.Types[0].Usages[0].Name != "Police" {
		t.Fatalf("merged usages mismatch")
	}

	if len(merged.Types[0].Tags) != 1 || merged.Types[0].Tags[0].Name != "shelves" {
		t.Fatalf("merged tags should remain unchanged when incoming tags are empty")
	}
}

func TestMergeSpawnableTypesValueOverridesByName(t *testing.T) {
	current := &SpawnableTypesFile{
		Types: []SpawnableTypeDef{
			{
				Name:    "FirstAidKit",
				Hoarder: &EmptyElement{},
				Tags:    []NamedRef{{Name: "medical"}},
				Cargo:   []SpawnableCargo{{Chance: floatPtr(1)}},
				Attachments: []SpawnableAttachment{
					{Items: []SpawnableItem{{Name: "BandageDressing"}}},
				},
			},
		},
	}
	incoming := &SpawnableTypesFile{
		Types: []SpawnableTypeDef{
			{
				Name:        "FirstAidKit",
				Unique:      &EmptyElement{},
				Cargo:       []SpawnableCargo{{Chance: floatPtr(0.5)}},
				Tags:        []NamedRef{{Name: "shelves"}},
				Damage:      &SpawnableMinMax{Min: floatPtr(0), Max: floatPtr(0.2)},
				Attachments: []SpawnableAttachment{},
			},
		},
	}

	mergedRaw, err := mergeSpawnableTypesValue(current, incoming)
	if err != nil {
		t.Fatalf("mergeSpawnableTypesValue() error: %v", err)
	}

	merged, ok := mergedRaw.(*SpawnableTypesFile)
	if !ok {
		t.Fatalf("merged value type = %T, want *SpawnableTypesFile", mergedRaw)
	}

	if len(merged.Types) != 1 {
		t.Fatalf("merged types count = %d, want 1", len(merged.Types))
	}

	out := merged.Types[0]
	if out.Hoarder != nil {
		t.Fatalf("hoarder flag should be reset when include defines unique/hoarder")
	}

	if out.Unique == nil {
		t.Fatalf("unique flag should be set")
	}

	if len(out.Tags) != 1 || out.Tags[0].Name != "shelves" {
		t.Fatalf("tags mismatch")
	}

	if len(out.Cargo) != 1 || out.Cargo[0].Chance == nil || *out.Cargo[0].Chance != 0.5 {
		t.Fatalf("cargo should be overridden")
	}
}

func TestMergeEventsValueMergesChildren(t *testing.T) {
	current := &EventsFile{
		Events: []EventDef{
			{
				Name: "Loot",
				Children: &EventChildren{
					Children: []EventChild{
						{Type: "TypeA", Min: intPtr(100), Max: intPtr(0)},
					},
				},
			},
		},
	}
	incoming := &EventsFile{
		Events: []EventDef{
			{
				Name: "Loot",
				Children: &EventChildren{
					Children: []EventChild{
						{Type: "TypeA", Min: intPtr(50)},
						{Type: "TypeB", Min: intPtr(25)},
					},
				},
			},
		},
	}

	mergedRaw, err := mergeEventsValue(current, incoming)
	if err != nil {
		t.Fatalf("mergeEventsValue() error: %v", err)
	}

	merged := mergedRaw.(*EventsFile)
	if len(merged.Events) != 1 {
		t.Fatalf("merged events count = %d, want 1", len(merged.Events))
	}

	children := merged.Events[0].Children.Children
	if len(children) != 2 {
		t.Fatalf("merged children count = %d, want 2", len(children))
	}

	if children[0].Min == nil || *children[0].Min != 50 {
		t.Fatalf("existing child should be modified")
	}

	if children[1].Type != "TypeB" {
		t.Fatalf("new child should be appended")
	}
}

func TestMergeGlobalsValueTypeConflict(t *testing.T) {
	current := &GlobalsFile{
		Vars: []GlobalVar{
			{Name: "AnimalMaxCount", Type: VariableTypeInt, Value: "200"},
		},
	}
	incoming := &GlobalsFile{
		Vars: []GlobalVar{
			{Name: "AnimalMaxCount", Type: VariableTypeFloat, Value: "200.0"},
		},
	}

	_, err := mergeGlobalsValue(current, incoming)
	if err == nil {
		t.Fatal("mergeGlobalsValue() expected conflict error")
	}

	if !errors.Is(err, ErrMergeConflict) {
		t.Fatalf("error = %v, want ErrMergeConflict", err)
	}
}

// intPtr returns pointer for test int literals.
func intPtr(value int) *int {
	return &value
}

// floatPtr returns pointer for test float literals.
func floatPtr(value float64) *float64 {
	return &value
}
