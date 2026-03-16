// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MergeSource describes one included file consumed during merge.
type MergeSource struct {
	// Kind is CE kind used for decoding this source.
	Kind Kind
	// Path is absolute path to source file.
	Path string
}

// MergeOptions configures economycore include processing behavior.
type MergeOptions struct {
	// RelaxedIncludeTypes allows include `type` values outside strict wiki set.
	//
	// When enabled, include kind resolution fallback order is:
	// 1. strict wiki mapping (`types`, `spawnabletypes`, `globals`, `economy`,
	//    `events`, `messages`, plus `economycore`)
	// 2. direct kind names known to registry (`undergroundtriggers`, `gameplay`,
	//    `effectarea`, and so on)
	// 3. kind detection from include file name/path (`DetectKind`)
	RelaxedIncludeTypes bool
}

// MergedConfig contains merged CE payload grouped by kind.
type MergedConfig struct {
	// Values stores merged payload by CE kind.
	Values map[Kind]any
	// Sources stores consumed include sources in load order.
	Sources []MergeSource
}

// Get returns merged payload for requested kind.
func (config *MergedConfig) Get(kind Kind) (any, bool) {
	if config == nil || config.Values == nil {
		return nil, false
	}

	value, ok := config.Values[kind]
	return value, ok
}

// LoadMergedEconomyCore loads economycore tree and merges included CE files.
//
// Input path can be any economycore XML file name. Include entries are resolved
// relative to each currently processed economycore file and its `<ce folder=...>`
// value. Include type is resolved only from `<file type="...">` and must
// match wiki-compatible include set (plus recursive `economycore`).
func LoadMergedEconomyCore(path string) (*MergedConfig, error) {
	return LoadMergedEconomyCoreWithOptions(path, MergeOptions{})
}

// LoadMergedEconomyCoreWithOptions loads economycore tree with custom options.
func LoadMergedEconomyCoreWithOptions(
	path string,
	options MergeOptions,
) (*MergedConfig, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve economycore path %q: %w", path, err)
	}

	result := &MergedConfig{
		Values: make(map[Kind]any),
	}

	if err = loadEconomyCoreRecursive(
		absolutePath,
		result,
		map[string]struct{}{},
		options,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// loadEconomyCoreRecursive walks one economycore file and nested includes.
func loadEconomyCoreRecursive(
	path string,
	result *MergedConfig,
	stack map[string]struct{},
	options MergeOptions,
) error {
	cleanPath := filepath.Clean(path)

	if _, exists := stack[cleanPath]; exists {
		return fmt.Errorf("%w: %s", ErrEconomyCoreCycle, cleanPath)
	}

	stack[cleanPath] = struct{}{}
	defer delete(stack, cleanPath)

	raw, err := os.ReadFile(cleanPath)
	if err != nil {
		return fmt.Errorf("read economycore %q: %w", cleanPath, err)
	}

	core, err := DecodeEconomyCore(raw)
	if err != nil {
		return fmt.Errorf("decode economycore %q: %w", cleanPath, err)
	}

	baseDir := filepath.Dir(cleanPath)

	for _, includeGroup := range core.CE {
		folder := strings.TrimSpace(includeGroup.Folder)

		for _, include := range includeGroup.Files {
			includePath := filepath.Join(
				baseDir,
				folder,
				strings.TrimSpace(include.Name),
			)
			includePath = filepath.Clean(includePath)

			kind := resolveIncludeKind(include.Type, includePath, options)
			if kind == KindUnknown {
				return fmt.Errorf(
					"%w: type=%q file=%q in %q",
					ErrUnknownFileKind,
					include.Type,
					include.Name,
					cleanPath,
				)
			}

			if kind == KindEconomyCore {
				if err = loadEconomyCoreRecursive(
					includePath,
					result,
					stack,
					options,
				); err != nil {
					return err
				}

				continue
			}

			if err = loadAndMergeInclude(includePath, kind, result); err != nil {
				return err
			}
		}
	}

	return nil
}

// resolveIncludeKind resolves include kind according to selected mode.
func resolveIncludeKind(
	includeType string,
	includePath string,
	options MergeOptions,
) Kind {
	kind := KindFromEconomyCoreType(includeType)
	if kind != KindUnknown {
		return kind
	}

	if !options.RelaxedIncludeTypes {
		return KindUnknown
	}

	if kind = parseRegisteredKind(includeType); kind != KindUnknown {
		return kind
	}

	return DetectKind(includePath)
}

// parseRegisteredKind resolves raw kind string if registry has such kind.
func parseRegisteredKind(value string) Kind {
	kind := Kind(strings.ToLower(strings.TrimSpace(value)))
	if kind == KindUnknown {
		return KindUnknown
	}

	if _, ok := DefaultRegistry().Get(kind); !ok {
		return KindUnknown
	}

	return kind
}

// loadAndMergeInclude reads one include file and merges it into result.
func loadAndMergeInclude(path string, kind Kind, result *MergedConfig) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read include %q: %w", path, err)
	}

	value, err := Decode(kind, raw)
	if err != nil {
		return fmt.Errorf("decode include %q as %q: %w", path, kind, err)
	}

	merged, err := mergeKindValue(kind, result.Values[kind], value)
	if err != nil {
		return fmt.Errorf("merge include %q as %q: %w", path, kind, err)
	}

	result.Values[kind] = merged
	result.Sources = append(result.Sources, MergeSource{
		Kind: kind,
		Path: path,
	})

	return nil
}
