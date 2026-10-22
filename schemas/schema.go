// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package schemas

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"sync"
)

const schemaSuffix = ".schema.json"

// embeddedFiles stores the schema documents under the data directory.
//
//go:embed data/*.schema.json
var embeddedFiles embed.FS

// files exposes embedded schema documents at the package filesystem root.
var files = func() fs.FS {
	sub, err := fs.Sub(embeddedFiles, "data")
	if err != nil {
		panic(err)
	}

	return sub
}()

var (
	// ErrFieldNotFound reports that an exact schema field key was not found.
	ErrFieldNotFound = errors.New("schema field not found")
	// ErrAmbiguousField reports that an exact key matched multiple fields.
	ErrAmbiguousField = errors.New("schema field is ambiguous")
	// ErrInvalidSearchOptions reports conflicting or repeated search options.
	ErrInvalidSearchOptions = errors.New("invalid schema search options")
)

// schemaIndex caches parsed embedded schema metadata and definitions.
var schemaIndex struct {
	once    sync.Once              // once ensures embedded schema documents are parsed only once.
	entries []schemaEntry          // entries contains schemas in catalog order.
	byName  map[string]schemaEntry // byName provides constant-time schema lookup.
	err     error                  // err stores the first schema loading failure.
}

// SchemaInfo describes one embedded JSON Schema document.
type SchemaInfo struct {
	// Name is the schema name accepted by Read and Lookup.
	Name string `json:"name" yaml:"name"`
	// Description describes the root configuration type represented by the schema.
	Description string `json:"description" yaml:"description"`
}

// SearchOptions limits schema searches and controls nested field traversal.
type SearchOptions struct {
	// Schema limits a search to one schema. An empty value searches all schemas.
	Schema string `json:"schema,omitempty" yaml:"schema,omitempty"`
	// IncludeNested includes fields below the root object. It defaults to false.
	IncludeNested bool `json:"includeNested,omitempty" yaml:"includeNested,omitempty"`
}

// Field describes one property in an embedded JSON Schema document.
type Field struct {
	// Schema is the name of the schema containing this field.
	Schema string `json:"schema" yaml:"schema"`
	// Name is the property name as written in the JSON document.
	Name string `json:"name" yaml:"name"`
	// Path is the field's schema path; [] marks an array item and * a map value.
	Path string `json:"path" yaml:"path"`
	// FieldSchema contains the property's modeled JSON Schema keywords.
	FieldSchema
	// Required reports whether the property is required by its parent object.
	Required bool `json:"required" yaml:"required"`
}

// FieldSchema contains the modeled metadata and constraints for a property.
type FieldSchema struct {
	// Default is the schema's default value, when declared.
	Default any `json:"default,omitempty" yaml:"default,omitempty"`
	// Minimum is the inclusive lower bound for numeric values.
	Minimum *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	// Maximum is the inclusive upper bound for numeric values.
	Maximum *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	// MinLength is the minimum allowed string length.
	MinLength *uint64 `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	// MaxLength is the maximum allowed string length.
	MaxLength *uint64 `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	// MinItems is the minimum number of items allowed in an array.
	MinItems *uint64 `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	// MaxItems is the maximum number of items allowed in an array.
	MaxItems *uint64 `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	// AdditionalProperties declares whether an object accepts undeclared keys.
	AdditionalProperties *bool `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	// Items describes the schema applied to each array element.
	Items *FieldSchema `json:"items,omitempty" yaml:"items,omitempty"`
	// Description explains the property's purpose.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Type is the JSON Schema type, such as object, array, or string.
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Reference is the referenced definition name, when the property uses $ref.
	Reference string `json:"reference,omitempty" yaml:"reference,omitempty"`
	// Enum lists the allowed values when the schema declares an enum.
	Enum []any `json:"enum,omitempty" yaml:"enum,omitempty"`
	// Examples contains representative values declared by the schema.
	Examples []any `json:"examples,omitempty" yaml:"examples,omitempty"`
}

// schemaDocument contains document metadata and reusable type definitions.
type schemaDocument struct {
	// ID identifies this schema document.
	ID string `json:"$id"`
	// Schema is the JSON Schema dialect used by this document.
	Schema string `json:"$schema"`
	// Definitions contains named schema definitions from $defs.
	Definitions map[string]schemaDefinition `json:"$defs"`
	// Ref points to the definition representing the document root.
	Ref string `json:"$ref"`
}

// schemaDefinition contains JSON Schema metadata used for field lookup.
type schemaDefinition struct {
	// Description explains the configuration type or field.
	Description string `json:"description"`
	// Type retains the JSON Schema type value, including unions.
	Type json.RawMessage `json:"type"`
	// Ref points to a reusable definition when this node is a reference.
	Ref string `json:"$ref"`
	// Enum contains the values accepted by an enum constraint.
	Enum []any `json:"enum"`
	// Default is the schema's default value.
	Default any `json:"default"`
	// Examples contains representative values for this schema.
	Examples []any `json:"examples"`
	// Minimum is the inclusive lower bound for numeric values.
	Minimum *float64 `json:"minimum"`
	// Maximum is the inclusive upper bound for numeric values.
	Maximum *float64 `json:"maximum"`
	// MinLength is the minimum allowed string length.
	MinLength *uint64 `json:"minLength"`
	// MaxLength is the maximum allowed string length.
	MaxLength *uint64 `json:"maxLength"`
	// MinItems is the minimum number of items allowed in an array.
	MinItems *uint64 `json:"minItems"`
	// MaxItems is the maximum number of items allowed in an array.
	MaxItems *uint64 `json:"maxItems"`
	// Properties contains object fields keyed by their JSON names.
	Properties map[string]schemaProperty `json:"properties"`
	// Required lists properties required by this object.
	Required []string `json:"required"`
	// Items describes array elements.
	Items *schemaProperty `json:"items"`
	// PrefixItems describes schemas for tuple positions.
	PrefixItems []schemaProperty `json:"prefixItems"`
	// AnyOf contains alternative schemas accepted by anyOf.
	AnyOf []schemaProperty `json:"anyOf"`
	// OneOf contains alternative schemas accepted by oneOf.
	OneOf []schemaProperty `json:"oneOf"`
	// AllOf contains schemas whose constraints are combined.
	AllOf []schemaProperty `json:"allOf"`
	// AdditionalProperties retains the schema or boolean policy for extra keys.
	AdditionalProperties json.RawMessage `json:"additionalProperties"`
}

// schemaEntry holds one parsed schema document and its catalog metadata.
type schemaEntry struct {
	SchemaInfo                // document is the parsed source schema.
	document   schemaDocument // rootName identifies the root definition in document.
	rootName   string
}

// schemaProperty contains one field's JSON Schema definition.
type schemaProperty struct {
	schemaDefinition // schemaDefinition contains the decoded schema keywords.
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

// Catalog returns the available schemas and their root type descriptions.
func Catalog() ([]SchemaInfo, error) {
	if err := loadSchemaIndex(); err != nil {
		return nil, err
	}

	result := make([]SchemaInfo, len(schemaIndex.entries))
	for index, entry := range schemaIndex.entries {
		result[index] = entry.SchemaInfo
	}

	return result, nil
}

// Read returns the raw JSON Schema document identified by name.
// The name is one of the values returned by Names.
func Read(name string) ([]byte, error) {
	data, err := fs.ReadFile(files, name+schemaSuffix)
	if err != nil {
		return nil, fmt.Errorf("read schema %q: %w", name, err)
	}

	return data, nil
}

// Fields returns the documented properties of a schema's root object.
// It returns an empty slice when the root schema is not an object.
func Fields(name string) ([]Field, error) {
	entry, err := schemaByName(name)
	if err != nil {
		return nil, err
	}

	return fieldsFromEntry(entry, false)
}

// Lookup returns the field whose key exactly matches key in schemaName.
// Nested fields are excluded unless IncludeNested is enabled.
// If nested lookup finds multiple fields with the same key, Lookup returns ErrAmbiguousField.
func Lookup(
	schemaName string,
	key string,
	options ...SearchOptions,
) (Field, error) {
	settings, err := searchOptions(options)
	if err != nil {
		return Field{}, err
	}
	if settings.Schema != "" && settings.Schema != schemaName {
		return Field{}, fmt.Errorf("%w: schema option conflicts with lookup schema", ErrInvalidSearchOptions)
	}

	fields, err := fieldsForSchema(schemaName, settings.IncludeNested)
	if err != nil {
		return Field{}, err
	}

	var found Field
	for _, field := range fields {
		if field.Name != key {
			continue
		}

		if found.Name != "" {
			return Field{}, fmt.Errorf(
				"%w: %q occurs at %s and %s",
				ErrAmbiguousField,
				key,
				found.Path,
				field.Path,
			)
		}

		found = field
	}

	if found.Name == "" {
		return Field{}, fmt.Errorf("%w: schema=%q key=%q", ErrFieldNotFound, schemaName, key)
	}

	return found, nil
}

// Search returns fields matching query, ordered from strongest to weakest match.
// Matching is case-insensitive. Literal matches rank above weighted aliases,
// followed by ordered character subsequences in field names.
// An empty result is not an error. Nested fields are excluded unless IncludeNested is enabled.
func Search(query string, options ...SearchOptions) ([]Field, error) {
	settings, err := searchOptions(options)
	if err != nil {
		return nil, err
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return []Field{}, nil
	}
	query = strings.ToLower(query)
	if err := loadSchemaIndex(); err != nil {
		return nil, err
	}

	entries := schemaIndex.entries
	if settings.Schema != "" {
		entry, ok := schemaIndex.byName[settings.Schema]
		if !ok {
			return nil, fmt.Errorf("read schema %q: %w", settings.Schema, fs.ErrNotExist)
		}

		entries = []schemaEntry{entry}
	}

	type scoredField struct {
		field Field // field is the matching schema property.
		score int   // score determines this result's position within its match tier.
		tier  int   // tier ranks literal, alias, then subsequence matches.
	}

	scored := make([]scoredField, 0)
	expansions := searchExpansionsFor(query)
	for _, entry := range entries {
		fields, fieldsErr := fieldsFromEntry(entry, settings.IncludeNested)
		if fieldsErr != nil {
			return nil, fieldsErr
		}

		for _, field := range fields {
			score, literalMatched := fieldLiteralMatchScore(query, field)
			if literalMatched {
				score, _ = fieldMatchScore(query, field)
				scored = append(scored, scoredField{field: field, score: score, tier: 2})
				continue
			}

			for _, expansion := range expansions {
				aliasScore, aliasMatched := fieldLiteralMatchScore(expansion.term, field)
				if !aliasMatched {
					continue
				}

				weightedScore := int(float64(aliasScore) * expansion.weight)
				if weightedScore > score {
					score = weightedScore
				}
			}
			if score > 0 {
				scored = append(scored, scoredField{field: field, score: score, tier: 1})
				continue
			}

			if score, matched := fieldMatchScore(query, field); matched {
				scored = append(scored, scoredField{field: field, score: score, tier: 0})
			}
		}
	}

	sort.Slice(scored, func(left, right int) bool {
		if scored[left].tier != scored[right].tier {
			return scored[left].tier > scored[right].tier
		}
		if scored[left].score != scored[right].score {
			return scored[left].score > scored[right].score
		}
		if scored[left].field.Schema != scored[right].field.Schema {
			return scored[left].field.Schema < scored[right].field.Schema
		}

		return scored[left].field.Path < scored[right].field.Path
	})

	result := make([]Field, len(scored))
	for index, item := range scored {
		result[index] = item.field
	}

	return result, nil
}
