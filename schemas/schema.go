// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

// Package schemas embeds the generated DayZ configuration JSON schemas.
package schemas

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

const schemaSuffix = ".schema.json"

//go:embed *.schema.json
var files embed.FS

// Field describes one property on a schema's root object.
type Field struct {
	// Name is the property name as written in the JSON document.
	Name string
	// Description is the property's schema description, if present.
	Description string
	// Type is the JSON Schema type, such as object, array, or string.
	Type string
	// Reference is the referenced definition name, when the property uses $ref.
	Reference string
	// Required reports whether the property is listed in the root definition.
	Required bool
}

// FS returns a read-only filesystem containing the embedded schema files.
func FS() fs.FS {
	return files
}

// Names returns the embedded schema names without the .schema.json suffix.
func Names() []string {
	paths, _ := fs.Glob(files, "*"+schemaSuffix)
	names := make([]string, 0, len(paths))

	for _, path := range paths {
		names = append(names, strings.TrimSuffix(path, schemaSuffix))
	}

	sort.Strings(names)

	return names
}

// Read returns the raw JSON Schema document identified by name. The name is
// one of the values returned by Names.
func Read(name string) ([]byte, error) {
	data, err := fs.ReadFile(files, name+schemaSuffix)
	if err != nil {
		return nil, fmt.Errorf("read schema %q: %w", name, err)
	}

	return data, nil
}

// Fields returns the documented properties of a schema's root object. It
// returns an empty slice when the root schema is not an object.
func Fields(name string) ([]Field, error) {
	data, err := Read(name)
	if err != nil {
		return nil, err
	}

	var document schemaDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("decode schema %q: %w", name, err)
	}

	rootName := strings.TrimPrefix(document.Ref, "#/$defs/")
	root, ok := document.Definitions[rootName]
	if !ok || root.Type != "object" {
		return []Field{}, nil
	}

	required := make(map[string]struct{}, len(root.Required))
	for _, fieldName := range root.Required {
		required[fieldName] = struct{}{}
	}

	fieldNames := make([]string, 0, len(root.Properties))
	for fieldName := range root.Properties {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	fields := make([]Field, 0, len(fieldNames))
	for _, fieldName := range fieldNames {
		property := root.Properties[fieldName]
		reference, hasReference := strings.CutPrefix(property.Ref, "#/$defs/")
		if !hasReference {
			reference = ""
		}
		fieldType := property.Type
		if fieldType == "" && reference != "" {
			fieldType = document.Definitions[reference].Type
		}
		_, isRequired := required[fieldName]
		fields = append(fields, Field{
			Name:        fieldName,
			Description: property.Description,
			Type:        fieldType,
			Reference:   reference,
			Required:    isRequired,
		})
	}

	return fields, nil
}

// schemaDocument contains the root reference and reusable type definitions.
type schemaDocument struct {
	Definitions map[string]schemaDefinition `json:"$defs"`
	Ref         string                      `json:"$ref"`
}

// schemaDefinition contains the schema metadata needed to list root fields.
type schemaDefinition struct {
	Type       string                    `json:"type"`
	Properties map[string]schemaProperty `json:"properties"`
	Required   []string                  `json:"required"`
}

// schemaProperty contains the field description and its JSON Schema type.
type schemaProperty struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Ref         string `json:"$ref"`
}
