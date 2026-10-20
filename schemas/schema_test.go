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

// contains reports whether values contains want.
func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}

	return false
}
