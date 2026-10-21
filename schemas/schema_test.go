// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package schemas_test

import (
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"testing"

	"github.com/woozymasta/dzce/schemas"
)

func ExampleRead() {
	data, err := schemas.Read("economy")
	if err != nil {
		panic(err)
	}

	fmt.Println(len(data) > 0)
	// Output: true
}

func ExampleLookup() {
	field, err := schemas.Lookup(
		"types",
		"lifetime",
		schemas.SearchOptions{IncludeNested: true},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s (%s): %s\n", field.Path, field.Type, field.Description)
	fmt.Printf("Example: %.0f\n", field.Examples[0])
	// Output:
	// $.type[].lifetime (integer): Lifetime is how long, in seconds, an eligible spawned instance can remain before economy cleanup removes it.
	// Example: 14400
}

func ExampleSearch() {
	results, err := schemas.Search(
		"lft",
		schemas.SearchOptions{
			Schema:        "types",
			IncludeNested: true,
		},
	)
	if err != nil {
		panic(err)
	}
	if len(results) == 0 {
		panic("no matching fields")
	}

	fmt.Printf("%s %s %s\n", results[0].Schema, results[0].Name, results[0].Path)
	// Output: types lifetime $.type[].lifetime
}

func TestNames(t *testing.T) {
	t.Parallel()

	names := schemas.Names()
	if len(names) == 0 {
		t.Fatal("Names() returned no schemas")
	}
	if !sort.StringsAreSorted(names) {
		t.Fatalf("Names() is not sorted: %v", names)
	}
	if !contains(names, "types") || !contains(names, "gameplay") {
		t.Fatalf("Names() is missing expected schemas: %v", names)
	}
	for i := 1; i < len(names); i++ {
		if names[i-1] == names[i] {
			t.Fatalf("Names() contains duplicate %q", names[i])
		}
	}
}

func TestCatalog(t *testing.T) {
	t.Parallel()

	catalog, err := schemas.Catalog()
	if err != nil {
		t.Fatalf("Catalog() error = %v", err)
	}
	if len(catalog) != len(schemas.Names()) {
		t.Fatalf("Catalog() returned %d schemas, want %d", len(catalog), len(schemas.Names()))
	}

	for _, schema := range catalog {
		if schema.Name == "types" {
			if schema.Description == "" {
				t.Fatal("Catalog(types) has no description")
			}
			return
		}
	}

	t.Fatal("Catalog() has no types schema")
}

func TestRead(t *testing.T) {
	t.Parallel()

	data, err := schemas.Read("types")
	if err != nil {
		t.Fatalf("Read(types) error = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("Read(types) returned an empty schema")
	}
	if _, err := fs.Stat(schemas.FS(), "types.schema.json"); err != nil {
		t.Fatalf("FS() cannot read the types schema: %v", err)
	}

	if _, err := schemas.Read("missing"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Read(missing) error = %v, want fs.ErrNotExist", err)
	}
}

func TestFields(t *testing.T) {
	t.Parallel()

	fields, err := schemas.Fields("economy")
	if err != nil {
		t.Fatalf("Fields(economy) error = %v", err)
	}

	var economyField *schemas.Field
	for i := range fields {
		if fields[i].Name == "economy" {
			economyField = &fields[i]
			break
		}
	}
	if economyField == nil {
		t.Fatal("Fields(economy) has no economy property")
	}
	if !economyField.Required {
		t.Fatal("economy field is not marked required")
	}
	if economyField.Type != "object" || economyField.Reference != "xml_Name" {
		t.Fatalf("economy field type/reference = %q/%q, want object/xml_Name", economyField.Type, economyField.Reference)
	}
	if economyField.Description == "" {
		t.Fatal("economy field has no description")
	}

	fields, err = schemas.Fields("gameplay-gear-presets")
	if err != nil {
		t.Fatalf("Fields(gameplay-gear-presets) error = %v", err)
	}
	if len(fields) != 0 {
		t.Fatalf("Fields(gameplay-gear-presets) = %v, want no root object fields", fields)
	}
}

func TestLookup(t *testing.T) {
	t.Parallel()

	if _, err := schemas.Lookup("types", "lifetime"); !errors.Is(err, schemas.ErrFieldNotFound) {
		t.Fatalf("Lookup() without nested fields error = %v, want ErrFieldNotFound", err)
	}

	field, err := schemas.Lookup(
		"types",
		"lifetime",
		schemas.SearchOptions{IncludeNested: true},
	)
	if err != nil {
		t.Fatalf("Lookup(types, lifetime, IncludeNested) error = %v", err)
	}
	if field.Path != "$.type[].lifetime" {
		t.Fatalf("Lookup() path = %q, want $.type[].lifetime", field.Path)
	}
	if field.Type != "integer" || field.Description == "" {
		t.Fatalf("Lookup() returned incomplete field metadata: %+v", field)
	}
	if len(field.Examples) != 1 || field.Examples[0] != float64(14400) {
		t.Fatalf("Lookup() examples = %v, want [14400]", field.Examples)
	}

	if _, err := schemas.Lookup(
		"spawnable-types",
		"item",
		schemas.SearchOptions{IncludeNested: true},
	); !errors.Is(err, schemas.ErrAmbiguousField) {
		t.Fatalf("Lookup() duplicate nested key error = %v, want ErrAmbiguousField", err)
	}
}

func TestLookupExposesTypedSchemaKeywords(t *testing.T) {
	t.Parallel()

	options := schemas.SearchOptions{IncludeNested: true}
	length, err := schemas.Lookup("bunker-broadcast", "lockCodeLength", options)
	if err != nil {
		t.Fatalf("Lookup(lockCodeLength) error = %v", err)
	}
	if length.Minimum == nil || *length.Minimum != 3 || length.Maximum == nil || *length.Maximum != 6 {
		t.Fatalf("Lookup(lockCodeLength) bounds = %v..%v, want 3..6", length.Minimum, length.Maximum)
	}
	*length.Minimum = 99

	mode, err := schemas.Lookup("bunker-broadcast", "lockCodeMode", options)
	if err != nil {
		t.Fatalf("Lookup(lockCodeMode) error = %v", err)
	}
	if len(mode.Enum) != 2 || mode.Enum[0] != float64(0) || mode.Enum[1] != float64(1) {
		t.Fatalf("Lookup(lockCodeMode) enum = %v, want [0 1]", mode.Enum)
	}
	mode.Enum[0] = float64(99)

	length, err = schemas.Lookup("bunker-broadcast", "lockCodeLength", options)
	if err != nil {
		t.Fatalf("second Lookup(lockCodeLength) error = %v", err)
	}
	if length.Minimum == nil || *length.Minimum != 3 {
		t.Fatalf("mutating a returned field changed cached data: minimum = %v", length.Minimum)
	}

	proxyPositions, err := schemas.Lookup("bunker-broadcast", "proxies", options)
	if err != nil {
		t.Fatalf("Lookup(proxies) error = %v", err)
	}
	if proxyPositions.Items == nil || proxyPositions.Items.Reference != "BunkerBroadcastProxy" {
		t.Fatalf("Lookup(proxies) item schema = %+v, want BunkerBroadcastProxy", proxyPositions.Items)
	}

	act, err := schemas.Lookup("economy-core", "act", options)
	if err != nil {
		t.Fatalf("Lookup(act) error = %v", err)
	}
	if act.Default != "none" {
		t.Fatalf("Lookup(act) default = %v, want none", act.Default)
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()

	options := schemas.SearchOptions{Schema: "types", IncludeNested: true}
	results, err := schemas.Search("lft", options)
	if err != nil {
		t.Fatalf("Search(lft) error = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("Search(lft) returned no results")
	}
	if results[0].Name != "lifetime" {
		t.Fatalf("Search(lft) first result = %q, want lifetime", results[0].Name)
	}
	results, err = schemas.Search("LFT", options)
	if err != nil {
		t.Fatalf("Search(LFT) error = %v", err)
	}
	if results[0].Name != "lifetime" {
		t.Fatalf("Search(LFT) first result = %q, want lifetime", results[0].Name)
	}

	results, err = schemas.Search("lft", schemas.SearchOptions{Schema: "types"})
	if err != nil {
		t.Fatalf("Search(lft) without nested fields error = %v", err)
	}
	for _, field := range results {
		if field.Name == "lifetime" {
			t.Fatalf("Search(lft) without nested fields returned nested field: %+v", field)
		}
	}

	results, err = schemas.Search("no-such-field-9f14", options)
	if err != nil {
		t.Fatalf("Search(no-such-field) error = %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("Search(no-such-field) returned %d results, want 0", len(results))
	}

	results, err = schemas.Search("zombie", schemas.SearchOptions{IncludeNested: true})
	if err != nil {
		t.Fatalf("Search(zombie) error = %v", err)
	}
	if len(results) == 0 || results[0].Path != "$.zombies" {
		var firstPath string
		if len(results) > 0 {
			firstPath = results[0].Path
		}
		t.Fatalf("Search(zombie) first path = %q, want $.zombies", firstPath)
	}
	for _, field := range results {
		if field.Path == "$.animals" {
			t.Fatal("Search(zombie) returned unrelated $.animals field")
		}
	}

	results, err = schemas.Search("min", schemas.SearchOptions{IncludeNested: true})
	if err != nil {
		t.Fatalf("Search(min) error = %v", err)
	}
	if len(results) == 0 || results[0].Name != "min" {
		var firstName string
		if len(results) > 0 {
			firstName = results[0].Name
		}
		t.Fatalf("Search(min) first name = %q, want exact field name min", firstName)
	}

	results, err = schemas.Search("cleanup removes", schemas.SearchOptions{Schema: "types", IncludeNested: true})
	if err != nil {
		t.Fatalf("Search(cleanup removes) error = %v", err)
	}
	if len(results) == 0 || results[0].Path != "$.type[].lifetime" {
		var firstPath string
		if len(results) > 0 {
			firstPath = results[0].Path
		}
		t.Fatalf("Search(cleanup removes) first path = %q, want $.type[].lifetime", firstPath)
	}

	if _, err := schemas.Search("type", schemas.SearchOptions{Schema: "missing"}); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Search() missing schema error = %v, want fs.ErrNotExist", err)
	}
}

func TestSearchAliases(t *testing.T) {
	t.Parallel()

	options := schemas.SearchOptions{IncludeNested: true}
	tests := []struct {
		query      string
		wantSchema string
		wantPath   string
	}{
		{query: "undead", wantSchema: "economy", wantPath: "$.zombies"},
		{query: "character", wantSchema: "economy", wantPath: "$.animals"},
		{query: "bike", wantSchema: "economy", wantPath: "$.vehicles"},
		{query: "cloud cover", wantSchema: "weather", wantPath: "$.overcast"},
		{query: "dungeon", wantSchema: "underground-triggers", wantPath: "$.Triggers[].UseLinePointFade"},
		{query: "dungeon", wantSchema: "bunker-broadcast", wantPath: "$.proxies"},
		{query: "cooldown", wantSchema: "types", wantPath: "$.type[].restock"},
		{query: "ui", wantSchema: "effect-area", wantPath: "$.Areas[].PlayerData.PPERequesterType"},
	}

	for _, test := range tests {
		t.Run(test.query+"/"+test.wantSchema+test.wantPath, func(t *testing.T) {
			t.Parallel()

			results, err := schemas.Search(test.query, options)
			if err != nil {
				t.Fatalf("Search(%q) error = %v", test.query, err)
			}

			for _, field := range results {
				if field.Schema == test.wantSchema && field.Path == test.wantPath {
					return
				}
			}

			t.Fatalf("Search(%q) did not return %s:%s", test.query, test.wantSchema, test.wantPath)
		})
	}
}

func TestSearchDirectMatchesRankAboveAliases(t *testing.T) {
	t.Parallel()

	results, err := schemas.Search(
		"weather",
		schemas.SearchOptions{Schema: "weather", IncludeNested: true},
	)
	if err != nil {
		t.Fatalf("Search(weather) error = %v", err)
	}

	weatherIndex := -1
	overcastIndex := -1
	for index, field := range results {
		switch field.Path {
		case "$.weather":
			weatherIndex = index
		case "$.overcast":
			overcastIndex = index
		}
	}

	if weatherIndex < 0 || overcastIndex < 0 {
		t.Fatalf("Search(weather) returned no direct/alias pair: weather=%d overcast=%d", weatherIndex, overcastIndex)
	}
	if weatherIndex >= overcastIndex {
		t.Fatalf("direct match index = %d, alias match index = %d; direct match must rank first", weatherIndex, overcastIndex)
	}
}

// contains reports whether values contains want.
func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}

	return false
}
