// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package schemas

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"unicode"
)

// fieldLiteralMatchScore matches aliases without fuzzy subsequence matching.
func fieldLiteralMatchScore(query string, field Field) (int, bool) {
	const (
		nameMatchScore        = 800_000
		pathMatchScore        = 500_000
		descriptionMatchScore = 300_000
	)

	bestScore := 0
	matched := false
	for _, candidate := range []struct {
		text      string
		baseScore int
	}{
		{
			text:      field.Name,
			baseScore: nameMatchScore,
		},
		{
			text:      field.Path,
			baseScore: pathMatchScore,
		},
		{
			text:      field.Description,
			baseScore: descriptionMatchScore,
		},
	} {
		candidateLower := strings.ToLower(candidate.text)
		index := strings.Index(candidateLower, query)
		if index < 0 {
			continue
		}

		score := candidate.baseScore - index*100 - len(candidateLower) + len(query)
		if !matched || score > bestScore {
			bestScore = score
			matched = true
		}
	}

	return bestScore, matched
}

// loadSchemaIndex parses and caches the embedded schemas on first use.
func loadSchemaIndex() error {
	schemaIndex.once.Do(func() {
		names := Names()
		schemaIndex.entries = make([]schemaEntry, 0, len(names))
		schemaIndex.byName = make(map[string]schemaEntry, len(names))

		for _, name := range names {
			data, err := Read(name)
			if err != nil {
				schemaIndex.err = err
				return
			}

			var document schemaDocument
			if err := json.Unmarshal(data, &document); err != nil {
				schemaIndex.err = fmt.Errorf("decode schema %q: %w", name, err)
				return
			}

			rootName, ok := strings.CutPrefix(document.Ref, "#/$defs/")
			if !ok {
				schemaIndex.err = fmt.Errorf("schema %q has unsupported root reference %q", name, document.Ref)
				return
			}

			root, ok := document.Definitions[rootName]
			if !ok {
				schemaIndex.err = fmt.Errorf("schema %q has no root definition %q", name, rootName)
				return
			}

			entry := schemaEntry{
				SchemaInfo: SchemaInfo{
					Name:        name,
					Description: root.Description,
				},
				document: document,
				rootName: rootName,
			}
			schemaIndex.entries = append(schemaIndex.entries, entry)
			schemaIndex.byName[name] = entry
		}
	})

	return schemaIndex.err
}

// schemaByName returns a cached schema entry or a filesystem-compatible error.
func schemaByName(name string) (schemaEntry, error) {
	if err := loadSchemaIndex(); err != nil {
		return schemaEntry{}, err
	}

	entry, ok := schemaIndex.byName[name]
	if !ok {
		return schemaEntry{}, fmt.Errorf("read schema %q: %w", name, fs.ErrNotExist)
	}

	return entry, nil
}

// fieldsForSchema collects fields according to the requested nesting policy.
func fieldsForSchema(name string, includeNested bool) ([]Field, error) {
	entry, err := schemaByName(name)
	if err != nil {
		return nil, err
	}

	return fieldsFromEntry(entry, includeNested)
}

// fieldsFromEntry collects root properties and optionally follows nested schemas.
func fieldsFromEntry(entry schemaEntry, includeNested bool) ([]Field, error) {
	root, ok := entry.document.Definitions[entry.rootName]
	if !ok {
		return nil, fmt.Errorf("schema %q has no root definition %q", entry.Name, entry.rootName)
	}
	if !includeNested && schemaTypeName(root, entry.document.Definitions) != "object" {
		return []Field{}, nil
	}

	visitedDefinitions := map[string]struct{}{entry.rootName: {}}
	seenPaths := make(map[string]struct{})
	fields := make([]Field, 0, len(root.Properties))
	collectSchemaFields(
		entry,
		root,
		"$",
		visitedDefinitions,
		seenPaths,
		includeNested,
		&fields,
	)

	return fields, nil
}

// collectSchemaFields walks one schema node while stopping recursive references.
func collectSchemaFields(
	entry schemaEntry,
	node schemaDefinition,
	path string,
	visitedDefinitions map[string]struct{},
	seenPaths map[string]struct{},
	includeNested bool,
	fields *[]Field,
) {
	if reference, ok := strings.CutPrefix(node.Ref, "#/$defs/"); ok {
		if _, visited := visitedDefinitions[reference]; !visited {
			if definition, exists := entry.document.Definitions[reference]; exists {
				nextVisited := cloneDefinitionSet(visitedDefinitions)
				nextVisited[reference] = struct{}{}
				collectSchemaFields(
					entry,
					definition,
					path,
					nextVisited,
					seenPaths,
					includeNested,
					fields,
				)
			}
		}
	}

	fieldNames := make([]string, 0, len(node.Properties))
	for name := range node.Properties {
		fieldNames = append(fieldNames, name)
	}
	sort.Strings(fieldNames)

	required := make(map[string]struct{}, len(node.Required))
	for _, name := range node.Required {
		required[name] = struct{}{}
	}

	for _, name := range fieldNames {
		property := node.Properties[name]
		fieldPath := schemaFieldPath(path, name)
		if _, exists := seenPaths[fieldPath]; !exists {
			_, isRequired := required[name]
			*fields = append(*fields, Field{
				Schema:      entry.Name,
				Name:        name,
				Path:        fieldPath,
				Required:    isRequired,
				FieldSchema: fieldSchemaFromNode(property.schemaDefinition, entry.document.Definitions, nil),
			})
			seenPaths[fieldPath] = struct{}{}
		}

		if includeNested {
			collectSchemaFields(
				entry,
				property.schemaDefinition,
				fieldPath,
				visitedDefinitions,
				seenPaths,
				true,
				fields,
			)
		}
	}

	if !includeNested {
		return
	}

	if node.Items != nil {
		collectSchemaFields(
			entry,
			node.Items.schemaDefinition,
			schemaArrayPath(path),
			visitedDefinitions,
			seenPaths,
			true,
			fields,
		)
	}
	for _, item := range node.PrefixItems {
		collectSchemaFields(
			entry,
			item.schemaDefinition,
			schemaArrayPath(path),
			visitedDefinitions,
			seenPaths,
			true,
			fields,
		)
	}

	for _, alternative := range node.AnyOf {
		collectSchemaFields(entry, alternative.schemaDefinition, path, visitedDefinitions, seenPaths, true, fields)
	}
	for _, alternative := range node.OneOf {
		collectSchemaFields(entry, alternative.schemaDefinition, path, visitedDefinitions, seenPaths, true, fields)
	}
	for _, alternative := range node.AllOf {
		collectSchemaFields(entry, alternative.schemaDefinition, path, visitedDefinitions, seenPaths, true, fields)
	}

	if len(node.AdditionalProperties) > 0 && node.AdditionalProperties[0] == '{' {
		var additional schemaDefinition
		if json.Unmarshal(node.AdditionalProperties, &additional) == nil {
			collectSchemaFields(
				entry,
				additional,
				path+".*",
				visitedDefinitions,
				seenPaths,
				true,
				fields,
			)
		}
	}
}

// fieldSchemaFromNode resolves one schema node into the public field contract.
func fieldSchemaFromNode(
	node schemaDefinition,
	definitions map[string]schemaDefinition,
	visitedDefinitions map[string]struct{},
) FieldSchema {
	if visitedDefinitions == nil {
		visitedDefinitions = make(map[string]struct{})
	}

	reference, hasReference := strings.CutPrefix(node.Ref, "#/$defs/")
	_, referenceVisited := visitedDefinitions[reference]
	if hasReference && !referenceVisited {
		visitedDefinitions = cloneDefinitionSet(visitedDefinitions)
		visitedDefinitions[reference] = struct{}{}
	}

	fieldSchema := FieldSchema{
		Description: node.Description,
		Type:        schemaTypeName(node, definitions),
		Enum:        cloneSchemaValues(node.Enum),
		Default:     cloneSchemaValue(node.Default),
		Examples:    cloneSchemaValues(node.Examples),
		Minimum:     cloneFloat64(node.Minimum),
		Maximum:     cloneFloat64(node.Maximum),
		MinLength:   cloneUint64(node.MinLength),
		MaxLength:   cloneUint64(node.MaxLength),
		MinItems:    cloneUint64(node.MinItems),
		MaxItems:    cloneUint64(node.MaxItems),
	}
	if reference, ok := strings.CutPrefix(node.Ref, "#/$defs/"); ok {
		fieldSchema.Reference = reference
	}
	if node.Items != nil {
		itemSchema := fieldSchemaFromNode(node.Items.schemaDefinition, definitions, visitedDefinitions)
		fieldSchema.Items = &itemSchema
	}
	if len(node.AdditionalProperties) > 0 {
		var additionalProperties bool
		if json.Unmarshal(node.AdditionalProperties, &additionalProperties) == nil {
			fieldSchema.AdditionalProperties = &additionalProperties
		}
	}

	if hasReference && !referenceVisited {
		if definition, exists := definitions[reference]; exists {
			fieldSchema = mergeFieldSchema(
				fieldSchema,
				fieldSchemaFromNode(definition, definitions, visitedDefinitions),
			)
		}
	}

	return fieldSchema
}

// mergeFieldSchema gives inline property metadata precedence over its reference.
func mergeFieldSchema(local FieldSchema, referenced FieldSchema) FieldSchema {
	if local.Description != "" {
		referenced.Description = local.Description
	}
	if local.Type != "" {
		referenced.Type = local.Type
	}
	if local.Reference != "" {
		referenced.Reference = local.Reference
	}
	if local.Enum != nil {
		referenced.Enum = local.Enum
	}
	if local.Default != nil {
		referenced.Default = local.Default
	}
	if local.Examples != nil {
		referenced.Examples = local.Examples
	}
	if local.Minimum != nil {
		referenced.Minimum = local.Minimum
	}
	if local.Maximum != nil {
		referenced.Maximum = local.Maximum
	}
	if local.MinLength != nil {
		referenced.MinLength = local.MinLength
	}
	if local.MaxLength != nil {
		referenced.MaxLength = local.MaxLength
	}
	if local.MinItems != nil {
		referenced.MinItems = local.MinItems
	}
	if local.MaxItems != nil {
		referenced.MaxItems = local.MaxItems
	}
	if local.AdditionalProperties != nil {
		referenced.AdditionalProperties = local.AdditionalProperties
	}
	if local.Items != nil {
		referenced.Items = local.Items
	}

	return referenced
}

// cloneSchemaValues copies schema values so callers cannot mutate cached data.
func cloneSchemaValues(values []any) []any {
	if values == nil {
		return nil
	}

	cloned := make([]any, len(values))
	for index, value := range values {
		cloned[index] = cloneSchemaValue(value)
	}

	return cloned
}

// cloneSchemaValue copies JSON values with nested maps and slices.
func cloneSchemaValue(value any) any {
	switch value := value.(type) {
	case []any:
		return cloneSchemaValues(value)

	case map[string]any:
		cloned := make(map[string]any, len(value))
		for key, nested := range value {
			cloned[key] = cloneSchemaValue(nested)
		}
		return cloned

	default:
		return value
	}
}

// cloneFloat64 returns an independent copy of an optional numeric bound.
func cloneFloat64(value *float64) *float64 {
	if value == nil {
		return nil
	}

	cloned := *value
	return &cloned
}

// cloneUint64 returns an independent copy of an optional integer bound.
func cloneUint64(value *uint64) *uint64 {
	if value == nil {
		return nil
	}

	cloned := *value
	return &cloned
}

// cloneDefinitionSet gives a referenced branch its own recursion stack.
func cloneDefinitionSet(source map[string]struct{}) map[string]struct{} {
	clone := make(map[string]struct{}, len(source)+1)
	for name := range source {
		clone[name] = struct{}{}
	}

	return clone
}

// schemaFieldPath appends a property name to the readable schema path.
func schemaFieldPath(parent string, name string) string {
	name = strings.ReplaceAll(name, `\`, `\\`)
	name = strings.ReplaceAll(name, ".", `\.`)

	return parent + "." + name
}

// schemaArrayPath appends the marker used for an array item's schema path.
func schemaArrayPath(parent string) string {
	return parent + "[]"
}

// schemaTypeName resolves a schema node's declared type or local reference type.
func schemaTypeName(node schemaDefinition, definitions map[string]schemaDefinition) string {
	var typeName string
	if json.Unmarshal(node.Type, &typeName) == nil && typeName != "" {
		return typeName
	}

	var typeNames []string
	if json.Unmarshal(node.Type, &typeNames) == nil && len(typeNames) > 0 {
		return strings.Join(typeNames, "|")
	}

	if reference, ok := strings.CutPrefix(node.Ref, "#/$defs/"); ok {
		return definitions[reference].typeNameWithoutReference()
	}

	return ""
}

// typeNameWithoutReference returns the explicitly declared type of a definition.
func (definition schemaDefinition) typeNameWithoutReference() string {
	var typeName string
	if json.Unmarshal(definition.Type, &typeName) == nil {
		return typeName
	}

	var typeNames []string
	if json.Unmarshal(definition.Type, &typeNames) == nil {
		return strings.Join(typeNames, "|")
	}

	return ""
}

// searchOptions validates optional settings and supplies their zero-value defaults.
func searchOptions(options []SearchOptions) (SearchOptions, error) {
	if len(options) > 1 {
		return SearchOptions{}, fmt.Errorf("%w: provide at most one options value", ErrInvalidSearchOptions)
	}
	if len(options) == 0 {
		return SearchOptions{}, nil
	}

	return options[0], nil
}

// fieldMatchScore prioritizes field names, then paths, then descriptions.
func fieldMatchScore(query string, field Field) (int, bool) {
	bestScore := 0
	matched := false

	nameScore, nameMatched := fieldNameMatchScore(query, field.Name)
	if nameMatched {
		bestScore = nameScore
		matched = true
	}

	pathScore, pathMatched := substringMatchScore(query, field.Path, 500_000)
	if pathMatched && (!matched || pathScore > bestScore) {
		bestScore = pathScore
		matched = true
	}

	descriptionScore, descriptionMatched := substringMatchScore(query, field.Description, 300_000)
	if descriptionMatched && (!matched || descriptionScore > bestScore) {
		bestScore = descriptionScore
		matched = true
	}

	return bestScore, matched
}

// fieldNameMatchScore ranks exact, prefix, substring, and subsequence matches.
func fieldNameMatchScore(query string, name string) (int, bool) {
	const (
		exactMatchScore    = 1_000_000
		prefixMatchScore   = 900_000
		containsMatchScore = 800_000
		fuzzyMatchScore    = 700_000
	)

	nameLower := strings.ToLower(name)
	switch {
	case nameLower == query:
		return exactMatchScore, true

	case strings.HasPrefix(nameLower, query):
		return prefixMatchScore - len(nameLower) + len(query), true

	case strings.Contains(nameLower, query):
		return containsMatchScore - strings.Index(nameLower, query)*100 - len(nameLower) + len(query), true
	}

	if len([]rune(query)) < 2 {
		return 0, false
	}

	score, matched := subsequenceMatchScore(query, name)
	if !matched {
		return 0, false
	}

	return fuzzyMatchScore + score, true
}

// substringMatchScore rewards direct matches while preferring earlier positions.
func substringMatchScore(query string, candidate string, baseScore int) (int, bool) {
	candidateLower := strings.ToLower(candidate)
	index := strings.Index(candidateLower, query)
	if index < 0 {
		return 0, false
	}

	return baseScore - index*100 - len(candidateLower) + len(query), true
}

// subsequenceMatchScore scores ordered character matches in a field name.
func subsequenceMatchScore(query string, candidate string) (int, bool) {
	queryRunes := []rune(query)
	candidateRunes := []rune(candidate)
	if len(queryRunes) == 0 || len(queryRunes) > len(candidateRunes) {
		return 0, false
	}

	queryIndex := 0
	previousMatch := -1
	score := 0
	for index, candidateRune := range candidateRunes {
		if unicode.ToLower(candidateRune) != queryRunes[queryIndex] {
			continue
		}

		score += 10
		if index == 0 {
			score += 20
		} else if isSearchWordBoundary(candidateRunes[index-1], candidateRune) {
			score += 12
		}

		if previousMatch >= 0 {
			gap := index - previousMatch - 1
			if gap == 0 {
				score += 12
			} else {
				score -= min(gap*2, 20)
			}
		}

		previousMatch = index
		queryIndex++
		if queryIndex == len(queryRunes) {
			score -= len(candidateRunes) - index - 1
			return score, true
		}
	}

	return 0, false
}

// isSearchWordBoundary identifies separator and camel-case boundaries.
func isSearchWordBoundary(previous rune, current rune) bool {
	if previous == '_' || previous == '-' || previous == '.' || previous == '/' || unicode.IsSpace(previous) {
		return true
	}

	return unicode.IsLower(previous) && unicode.IsUpper(current)
}
