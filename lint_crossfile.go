// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/woozymasta/lintkit/lint"
)

// AnalyzeLintEconomyCoreTree runs cross-file CE lint checks for economycore tree.
func AnalyzeLintEconomyCoreTree(path string) []lint.Diagnostic {
	rootPath := strings.TrimSpace(path)
	if rootPath == "" {
		return nil
	}

	collector := newEconomyCoreIncludeCollector(rootPath)
	collector.collect(rootPath)

	diagnostics := make([]lint.Diagnostic, 0, 32)
	diagnostics = append(diagnostics, collector.diagnostics...)
	diagnostics = append(diagnostics, analyzeTypeOverrideIncludes(collector.includes)...)

	merged, err := LoadMergedEconomyCore(rootPath)
	if err != nil {
		diagnostics = append(diagnostics, mapMergedLoadError(rootPath, err)...)
		return diagnostics
	}

	diagnostics = append(diagnostics, analyzeMergedCrossRefs(rootPath, merged)...)

	return diagnostics
}

// includeRef stores one parsed include entry.
type includeRef struct {
	// kind stores resolved include kind.
	kind Kind

	// path stores normalized include file path.
	path string
}

// economyCoreIncludeCollector traverses include graph and stores diagnostics.
type economyCoreIncludeCollector struct {
	// visited stores already visited economycore files.
	visited map[string]struct{}

	// stack stores current recursion stack for cycle detection.
	stack map[string]struct{}
	// rootPath stores normalized root economycore path.
	rootPath string

	// diagnostics stores collected include traversal diagnostics.
	diagnostics []lint.Diagnostic

	// includes stores non-economycore include entries.
	includes []includeRef
}

// newEconomyCoreIncludeCollector constructs include traversal collector.
func newEconomyCoreIncludeCollector(rootPath string) economyCoreIncludeCollector {
	return economyCoreIncludeCollector{
		rootPath:    filepath.Clean(rootPath),
		diagnostics: make([]lint.Diagnostic, 0, 8),
		includes:    make([]includeRef, 0, 64),
		visited:     make(map[string]struct{}),
		stack:       make(map[string]struct{}),
	}
}

// collect recursively walks one economycore file and nested includes.
func (collector *economyCoreIncludeCollector) collect(path string) {
	cleanPath := filepath.Clean(path)
	if _, exists := collector.stack[cleanPath]; exists {
		collector.diagnostics = append(collector.diagnostics, newDiagnostic(
			CodeMergeIncludeCycle,
			cleanPath,
			fmt.Sprintf("include cycle detected at %q", cleanPath),
		))
		return
	}

	if _, visited := collector.visited[cleanPath]; visited {
		return
	}

	raw, err := os.ReadFile(cleanPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			collector.diagnostics = append(collector.diagnostics, newDiagnostic(
				CodeMergeMissingIncludeFile,
				cleanPath,
				fmt.Sprintf("include file not found: %q", cleanPath),
			))
		}
		return
	}

	core, err := DecodeEconomyCore(raw)
	if err != nil {
		return
	}

	collector.visited[cleanPath] = struct{}{}
	collector.stack[cleanPath] = struct{}{}
	defer delete(collector.stack, cleanPath)

	baseDir := filepath.Dir(cleanPath)
	for index := range core.CE {
		folder := strings.TrimSpace(core.CE[index].Folder)
		for includeIndex := range core.CE[index].Files {
			include := core.CE[index].Files[includeIndex]
			includePath := filepath.Clean(
				filepath.Join(baseDir, folder, strings.TrimSpace(include.Name)),
			)
			kind := resolveIncludeKind(include.Type, includePath, MergeOptions{})

			if _, err = os.Stat(includePath); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					collector.diagnostics = append(collector.diagnostics, newDiagnostic(
						CodeMergeMissingIncludeFile,
						cleanPath,
						fmt.Sprintf(
							"missing include file %q referenced from %q",
							includePath,
							cleanPath,
						),
					))
				}
				continue
			}

			if kind == KindEconomyCore {
				collector.collect(includePath)
				continue
			}

			if kind == KindUnknown {
				continue
			}

			collector.includes = append(collector.includes, includeRef{
				kind: kind,
				path: includePath,
			})
		}
	}
}

// analyzeTypeOverrideIncludes reports duplicate type names across include files.
func analyzeTypeOverrideIncludes(includes []includeRef) []lint.Diagnostic {
	type typeSource struct {
		// name stores lower-cased type name.
		name string

		// source stores include file path.
		source string
	}

	typeSources := make([]typeSource, 0, 256)
	for index := range includes {
		if includes[index].kind != KindTypes {
			continue
		}

		raw, err := os.ReadFile(includes[index].path)
		if err != nil {
			continue
		}

		parsed, err := DecodeTypes(raw)
		if err != nil {
			continue
		}

		for typeIndex := range parsed.Types {
			nameKey := strings.ToLower(strings.TrimSpace(parsed.Types[typeIndex].Name))
			if nameKey == "" {
				continue
			}

			typeSources = append(typeSources, typeSource{
				name:   nameKey,
				source: includes[index].path,
			})
		}
	}

	byName := make(map[string]map[string]struct{}, len(typeSources))
	for index := range typeSources {
		sources, ok := byName[typeSources[index].name]
		if !ok {
			sources = make(map[string]struct{})
			byName[typeSources[index].name] = sources
		}
		sources[typeSources[index].source] = struct{}{}
	}

	diagnostics := make([]lint.Diagnostic, 0, 8)
	for name, sources := range byName {
		if len(sources) <= 1 {
			continue
		}

		diagnostics = append(diagnostics, newDiagnostic(
			CodeMergeDuplicateTypeOverride,
			"",
			fmt.Sprintf(
				"type %q appears in %d include files and uses override merge semantics",
				name,
				len(sources),
			),
		))
	}

	return diagnostics
}

// mapMergedLoadError converts merged loader failures to lint diagnostics.
func mapMergedLoadError(path string, err error) []lint.Diagnostic {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrEconomyCoreCycle):
		return []lint.Diagnostic{
			newDiagnostic(
				CodeMergeIncludeCycle,
				path,
				err.Error(),
			),
		}
	case errors.Is(err, os.ErrNotExist):
		return []lint.Diagnostic{
			newDiagnostic(
				CodeMergeMissingIncludeFile,
				path,
				err.Error(),
			),
		}
	default:
		return nil
	}
}

// analyzeMergedCrossRefs runs reference checks on merged CE payload.
func analyzeMergedCrossRefs(path string, merged *MergedConfig) []lint.Diagnostic {
	if merged == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, 16)

	typesValue, okTypes := merged.Get(KindTypes)
	eventsValue, okEvents := merged.Get(KindEvents)
	if !okTypes || !okEvents {
		return diagnostics
	}

	typesFile, ok := castValue[TypesFile](typesValue)
	if !ok || typesFile == nil {
		return diagnostics
	}

	eventsFile, ok := castValue[EventsFile](eventsValue)
	if !ok || eventsFile == nil {
		return diagnostics
	}

	typeSet := make(map[string]struct{}, len(typesFile.Types))
	for index := range typesFile.Types {
		nameKey := strings.ToLower(strings.TrimSpace(typesFile.Types[index].Name))
		if nameKey == "" {
			continue
		}
		typeSet[nameKey] = struct{}{}
	}

	eventSet := make(map[string]struct{}, len(eventsFile.Events))
	for index := range eventsFile.Events {
		nameKey := strings.ToLower(strings.TrimSpace(eventsFile.Events[index].Name))
		if nameKey == "" {
			continue
		}
		eventSet[nameKey] = struct{}{}
	}

	for eventIndex := range eventsFile.Events {
		event := eventsFile.Events[eventIndex]
		if event.Secondary != nil {
			secondary := strings.ToLower(strings.TrimSpace(*event.Secondary))
			if secondary != "" {
				if _, ok = eventSet[secondary]; !ok {
					diagnostics = append(diagnostics, newDiagnostic(
						CodeCrossRefMissingEvent,
						path,
						fmt.Sprintf(
							"event %q references unknown secondary event %q",
							event.Name,
							*event.Secondary,
						),
					))
				}
			}
		}

		if event.Children == nil {
			continue
		}

		for childIndex := range event.Children.Children {
			childType := strings.ToLower(strings.TrimSpace(event.Children.Children[childIndex].Type))
			if childType == "" {
				continue
			}
			if _, ok = typeSet[childType]; ok {
				continue
			}

			diagnostics = append(diagnostics, newDiagnostic(
				CodeCrossRefMissingType,
				path,
				fmt.Sprintf(
					"event %q child references unknown type %q",
					event.Name,
					event.Children.Children[childIndex].Type,
				),
			))
		}
	}

	return diagnostics
}
