// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package schemas

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// TestEmbeddedSchemaKeywordsAreModeled prevents generated schema keywords
// from silently falling outside the parser and public field contract.
func TestEmbeddedSchemaKeywordsAreModeled(t *testing.T) {
	documentKeywords := schemaKeywordFields(reflect.TypeOf(schemaDocument{}))
	definitionKeywords := schemaKeywordFields(reflect.TypeOf(schemaDefinition{}))
	fieldKeywords := schemaKeywordFields(reflect.TypeOf(FieldSchema{}))
	structuralKeywords := map[string]struct{}{
		"$defs":       {},
		"$id":         {},
		"$schema":     {},
		"allOf":       {},
		"anyOf":       {},
		"oneOf":       {},
		"prefixItems": {},
		"properties":  {},
		"required":    {},
	}

	for keyword := range definitionKeywords {
		if _, structural := structuralKeywords[keyword]; structural {
			continue
		}

		fieldKeyword := keyword
		if keyword == "$ref" {
			fieldKeyword = "reference"
		}
		if _, modeled := fieldKeywords[fieldKeyword]; !modeled {
			t.Errorf("schema keyword %q is parsed but missing from FieldSchema", keyword)
		}
	}

	for _, name := range Names() {
		data, err := Read(name)
		if err != nil {
			t.Fatalf("Read(%q) error = %v", name, err)
		}

		var document map[string]json.RawMessage
		if err := json.Unmarshal(data, &document); err != nil {
			t.Fatalf("decode schema %q: %v", name, err)
		}
		for keyword := range document {
			if _, modeled := documentKeywords[keyword]; !modeled {
				t.Errorf("schema %q uses unmodeled document keyword %q", name, keyword)
			}
		}

		if definitions, ok := document["$defs"]; ok {
			var nodes map[string]json.RawMessage
			if err := json.Unmarshal(definitions, &nodes); err != nil {
				t.Fatalf("decode definitions in %q: %v", name, err)
			}
			for definition, node := range nodes {
				checkSchemaDefinition(t, name+"/$defs/"+definition, node, definitionKeywords)
			}
		}
	}
}

// schemaKeywordFields returns JSON property names declared by a struct.
func schemaKeywordFields(structType reflect.Type) map[string]struct{} {
	keywords := make(map[string]struct{}, structType.NumField())
	for index := range structType.NumField() {
		field := structType.Field(index)
		keyword := strings.Split(field.Tag.Get("json"), ",")[0]
		if keyword == "" || keyword == "-" {
			continue
		}

		keywords[keyword] = struct{}{}
	}

	return keywords
}

// checkSchemaDefinition validates one schema node and recursively checks child schemas.
func checkSchemaDefinition(
	t *testing.T,
	path string,
	data json.RawMessage,
	allowedKeywords map[string]struct{},
) {
	t.Helper()
	var node map[string]json.RawMessage
	if err := json.Unmarshal(data, &node); err != nil {
		t.Fatalf("decode schema node %s: %v", path, err)
	}

	for keyword, value := range node {
		if _, modeled := allowedKeywords[keyword]; !modeled {
			t.Errorf("schema node %s uses unmodeled keyword %q", path, keyword)
			continue
		}
		if keyword == "allOf" || keyword == "anyOf" || keyword == "oneOf" || keyword == "prefixItems" {
			t.Errorf("schema node %s uses composition keyword %q, which FieldSchema cannot represent", path, keyword)
		}

		switch keyword {
		case "$defs", "properties":
			var children map[string]json.RawMessage
			if err := json.Unmarshal(value, &children); err != nil {
				t.Fatalf("decode %s at %s: %v", keyword, path, err)
			}
			for name, child := range children {
				checkSchemaDefinition(t, path+"/"+keyword+"/"+name, child, allowedKeywords)
			}

		case "additionalProperties", "items":
			if strings.HasPrefix(strings.TrimSpace(string(value)), "{") {
				if keyword == "additionalProperties" {
					t.Errorf("schema node %s uses schema-valued additionalProperties, which FieldSchema cannot represent", path)
				}
				checkSchemaDefinition(t, path+"/"+keyword, value, allowedKeywords)
			}

		case "allOf", "anyOf", "oneOf", "prefixItems":
			var children []json.RawMessage
			if err := json.Unmarshal(value, &children); err != nil {
				t.Fatalf("decode %s at %s: %v", keyword, path, err)
			}
			for index, child := range children {
				checkSchemaDefinition(t, path+"/"+keyword+"/"+strconv.Itoa(index), child, allowedKeywords)
			}
		}
	}
}
