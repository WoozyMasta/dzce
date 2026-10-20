//go:build integration

// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const missionDirEnv = "DZCE_EXTERNAL_MPMISSION_DIR"

// TestDecodeRealWorldFixtures validates codec compatibility
// against external CE directories listed in DZCE_EXTERNAL_MPMISSION_DIR.
func TestDecodeRealWorldFixtures(t *testing.T) {
	roots := missionRoots(t)
	if len(roots) == 0 {
		t.Skip(missionDirEnv + " is not configured")
	}

	kinds := []struct {
		kind Kind
		rel  string
	}{
		{kind: KindTypes, rel: filepath.Join("db", "types.xml")},
		{kind: KindEvents, rel: filepath.Join("db", "events.xml")},
		{kind: KindEconomy, rel: filepath.Join("db", "economy.xml")},
		{kind: KindGlobals, rel: filepath.Join("db", "globals.xml")},
		{kind: KindSpawnableTypes, rel: "cfgspawnabletypes.xml"},
		{kind: KindRandomPresets, rel: "cfgrandompresets.xml"},
		{kind: KindEconomyCore, rel: "cfgeconomycore.xml"},
		{kind: KindEnvironment, rel: "cfgenvironment.xml"},
		{kind: KindEventSpawns, rel: "cfgeventspawns.xml"},
		{kind: KindEventGroups, rel: "cfgeventgroups.xml"},
		{kind: KindPlayerSpawnPoints, rel: "cfgplayerspawnpoints.xml"},
		{kind: KindWeather, rel: "cfgweather.xml"},
		{kind: KindLimitsDefinition, rel: "cfglimitsdefinition.xml"},
		{kind: KindLimitsDefinitionUser, rel: "cfglimitsdefinitionuser.xml"},
		{kind: KindIgnoreList, rel: "cfgignorelist.xml"},
		{kind: KindTerritories, rel: filepath.Join("env", "hare_territories.xml")},
		{kind: KindUndergroundTriggers, rel: "cfgundergroundtriggers.json"},
		{kind: KindEffectArea, rel: "cfgeffectarea.json"},
		{kind: KindGameplay, rel: "cfggameplay.json"},
		{kind: KindMapGroupProto, rel: "mapgroupproto.xml"},
		{kind: KindMapClusterProto, rel: "mapclusterproto.xml"},
		{kind: KindMapGroupPos, rel: "mapgrouppos.xml"},
		{kind: KindMapGroupDirt, rel: "mapgroupdirt.xml"},
		{kind: KindMapGroupCluster, rel: "mapgroupcluster.xml"},
	}

	roots := strings.Split(raw, string(os.PathListSeparator))

	for _, root := range roots {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}

		for _, target := range kinds {
			target := target
			path := filepath.Join(root, target.rel)

			t.Run(path, func(t *testing.T) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("ReadFile(%s) error: %v", path, err)
				}

				value, err := Decode(target.kind, data)
				if err != nil {
					t.Fatalf("Decode(%s) error: %v", path, err)
				}

				if _, err = Encode(target.kind, value); err != nil {
					t.Fatalf("Encode(%s) error: %v", path, err)
				}
			})
		}
	}
}

// TestRealMissionJSONModelCoverage detects unmodeled JSON files
// and fields in installed mission directories under DZCE_EXTERNAL_MPMISSION_DIR.
func TestRealMissionJSONModelCoverage(t *testing.T) {
	roots := missionRoots(t)
	if len(roots) == 0 {
		t.Skip(missionDirEnv + " is not configured")
	}

	for _, root := range roots {
		root := root
		t.Run(filepath.Base(root), func(t *testing.T) {
			references := missionJSONReferences(t, root)

			err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if entry.IsDir() || !strings.EqualFold(filepath.Ext(path), ".json") {
					return nil
				}

				kind := DetectKind(path)
				if kind == KindUnknown {
					kind = references[missionPathKey(root, path)]
				}
				if kind == KindUnknown {
					t.Errorf("unmodeled JSON file: %s", path)
					return nil
				}

				data, readErr := os.ReadFile(path)
				if readErr != nil {
					t.Errorf("read %s: %v", path, readErr)
					return nil
				}

				value, decodeErr := Decode(kind, data)
				if decodeErr != nil {
					t.Errorf("decode %s as %s: %v", path, kind, decodeErr)
					return nil
				}

				if _, encodeErr := Encode(kind, value); encodeErr != nil {
					t.Errorf("encode %s as %s: %v", path, kind, encodeErr)
				}

				unknown, inspectErr := unknownJSONFields(data, reflect.TypeOf(value))
				if inspectErr != nil {
					t.Errorf("inspect %s: %v", path, inspectErr)
					return nil
				}
				for _, field := range unknown {
					t.Errorf("unmodeled JSON field in %s: %s", path, field)
				}

				return nil
			})
			if err != nil {
				t.Fatalf("walk mission %s: %v", root, err)
			}
		})
	}
}

// TestRealMissionXMLModelCoverage detects unmodeled XML elements
// and attributes in installed mission directories under DZCE_EXTERNAL_MPMISSION_DIR.
func TestRealMissionXMLModelCoverage(t *testing.T) {
	roots := missionRoots(t)
	if len(roots) == 0 {
		t.Skip(missionDirEnv + " is not configured")
	}

	for _, root := range roots {
		root := root
		t.Run(filepath.Base(root), func(t *testing.T) {
			err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if entry.IsDir() || !strings.EqualFold(filepath.Ext(path), ".xml") {
					return nil
				}
				if strings.EqualFold(filepath.Base(path), "build.xml") {
					return nil
				}

				data, readErr := os.ReadFile(path)
				if readErr != nil {
					t.Errorf("read %s: %v", path, readErr)
					return nil
				}

				kind := DetectKind(path)
				if kind == KindUnknown {
					kind = detectKindByContent(path, data)
				}
				if kind == KindUnknown {
					t.Errorf("unmodeled XML file: %s", path)
					return nil
				}

				value, decodeErr := Decode(kind, data)
				if decodeErr != nil {
					t.Errorf("decode %s as %s: %v", path, kind, decodeErr)
					return nil
				}

				if _, encodeErr := Encode(kind, value); encodeErr != nil {
					t.Errorf("encode %s as %s: %v", path, kind, encodeErr)
				}

				unknown, inspectErr := unknownXMLFields(data, reflect.TypeOf(value))
				if inspectErr != nil {
					t.Errorf("inspect %s: %v", path, inspectErr)
					return nil
				}
				for _, field := range unknown {
					t.Errorf("unmodeled XML field in %s: %s", path, field)
				}

				return nil
			})
			if err != nil {
				t.Fatalf("walk mission %s: %v", root, err)
			}
		})
	}
}

// TestRealMissionAreaFlagsModelCoverage validates areaflags.map in each mission.
func TestRealMissionAreaFlagsModelCoverage(t *testing.T) {
	roots := missionRoots(t)
	if len(roots) == 0 {
		t.Skip(missionDirEnv + " is not configured")
	}

	for _, root := range roots {
		root := root
		t.Run(filepath.Base(root), func(t *testing.T) {
			path := filepath.Join(root, "areaflags.map")
			if kind := DetectKind(path); kind != KindAreaFlagsMap {
				t.Fatalf("DetectKind(%s) = %s, want %s", path, kind, KindAreaFlagsMap)
			}

			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}

			decoded, err := Decode(KindAreaFlagsMap, data)
			if err != nil {
				t.Fatalf("decode %s: %v", path, err)
			}

			encoded, err := Encode(KindAreaFlagsMap, decoded)
			if err != nil {
				t.Fatalf("encode %s: %v", path, err)
			}

			roundTripped, err := Decode(KindAreaFlagsMap, encoded)
			if err != nil {
				t.Fatalf("decode encoded %s: %v", path, err)
			}
			if !reflect.DeepEqual(decoded, roundTripped) {
				t.Fatalf("round-trip changed decoded data in %s", path)
			}
		})
	}
}

// missionRoots finds mission directories under the configured mpmissions root.
func missionRoots(t *testing.T) []string {
	t.Helper()

	raw := strings.TrimSpace(os.Getenv(missionDirEnv))
	if raw == "" {
		var err error
		raw, err = missionDirFromDotEnv()
		if err != nil {
			t.Fatalf("read .env: %v", err)
		}
	}
	if raw == "" {
		return nil
	}

	root := filepath.Clean(raw)
	if fileExists(filepath.Join(root, "cfggameplay.json")) {
		return []string{root}
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read mpmissions directory %s: %v", root, err)
	}

	roots := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		mission := filepath.Join(root, entry.Name())
		if fileExists(filepath.Join(mission, "cfggameplay.json")) {
			roots = append(roots, mission)
		}
	}
	sort.Strings(roots)
	if len(roots) == 0 {
		t.Fatalf("no mission directories with cfggameplay.json found under %s", root)
	}

	return roots
}

// missionDirFromDotEnv reads the mission root only when the environment omits it.
func missionDirFromDotEnv() (string, error) {
	data, err := os.ReadFile(".env")
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) != missionDirEnv {
			continue
		}
		value = strings.TrimSpace(value)
		if len(value) >= 2 &&
			(value[0] == '\'' && value[len(value)-1] == '\'' ||
				value[0] == '"' && value[len(value)-1] == '"') {
			value = value[1 : len(value)-1]
		}

		return strings.TrimSpace(value), nil
	}

	return "", nil
}

// fileExists reports whether path names an existing regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// missionJSONReferences maps gameplay-listed JSON files to their model kinds.
func missionJSONReferences(t *testing.T, root string) map[string]Kind {
	t.Helper()

	path := filepath.Join(root, "cfggameplay.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read gameplay config %s: %v", path, err)
	}

	gameplay, err := DecodeGameplay(data)
	if err != nil {
		t.Fatalf("decode gameplay config %s: %v", path, err)
	}

	references := make(map[string]Kind)
	if gameplay.PlayerData != nil {
		for _, name := range gameplay.PlayerData.SpawnGearPresetFiles {
			addMissionJSONReference(t, root, references, name, KindGameplayGearPresets)
		}
	}
	if gameplay.WorldsData != nil {
		for _, name := range gameplay.WorldsData.ObjectSpawnersArr {
			addMissionJSONReference(t, root, references, name, KindObjectSpawner)
		}
		for _, name := range gameplay.WorldsData.PlayerRestrictedAreaFiles {
			addMissionJSONReference(t, root, references, name, KindPlayerRestrictedArea)
		}
	}

	return references
}

// addMissionJSONReference adds one safe path referenced by cfggameplay.json.
func addMissionJSONReference(
	t *testing.T,
	root string,
	references map[string]Kind,
	name string,
	kind Kind,
) {
	t.Helper()

	if !filepath.IsLocal(name) {
		t.Errorf("mission JSON reference escapes mission root: %q", name)
		return
	}

	path := filepath.Join(root, name)
	references[missionPathKey(root, path)] = kind
}

// missionPathKey returns a case-insensitive normalized relative path key.
func missionPathKey(root string, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return strings.ToLower(filepath.Clean(path))
	}

	return strings.ToLower(filepath.Clean(rel))
}

// unknownJSONFields reports object keys absent from the decoded Go model type.
func unknownJSONFields(data []byte, model reflect.Type) ([]string, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("decode JSON for model inspection: %w", err)
	}

	unknown := make([]string, 0)
	collectUnknownJSONFields(value, model, "$", &unknown)
	unknown = uniqueSorted(unknown)

	return unknown, nil
}

// missionXMLNode stores XML element names, attributes, and child elements.
type missionXMLNode struct {
	name     string
	attrs    []xml.Attr
	children []*missionXMLNode
}

// unknownXMLFields reports XML elements and attributes absent from a model.
func unknownXMLFields(data []byte, model reflect.Type) ([]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var root *missionXMLNode
	stack := make([]*missionXMLNode, 0)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse XML for model inspection: %w", err)
		}

		switch element := token.(type) {
		case xml.StartElement:
			node := &missionXMLNode{name: element.Name.Local, attrs: element.Attr}
			if len(stack) == 0 {
				root = node
			} else {
				parent := stack[len(stack)-1]
				parent.children = append(parent.children, node)
			}
			stack = append(stack, node)

		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}

	if root == nil {
		return nil, fmt.Errorf("XML document has no root element")
	}

	unknown := make([]string, 0)
	collectUnknownXMLFields(root, model, "$", &unknown)
	unknown = uniqueSorted(unknown)

	return unknown, nil
}

// uniqueSorted returns stable unique strings for readable schema failures.
func uniqueSorted(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	unique := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}

		seen[value] = struct{}{}
		unique = append(unique, value)
	}
	sort.Strings(unique)

	return unique
}

// collectUnknownXMLFields walks XML nodes using model XML tags.
func collectUnknownXMLFields(node *missionXMLNode, model reflect.Type, path string, unknown *[]string) {
	for model.Kind() == reflect.Pointer {
		model = model.Elem()
	}
	if model.Kind() == reflect.Slice || model.Kind() == reflect.Array {
		model = model.Elem()
	}
	if model.Kind() != reflect.Struct {
		return
	}

	elements := make(map[string]reflect.Type)
	attributes := make(map[string]struct{})
	allowAnyElement := false
	allowAnyAttribute := false
	for index := 0; index < model.NumField(); index++ {
		field := model.Field(index)
		if !field.IsExported() || field.Name == "XMLName" {
			continue
		}

		parts := strings.Split(field.Tag.Get("xml"), ",")
		name := parts[0]
		if name == "-" {
			continue
		}

		isAttribute := false
		isAny := false
		isText := false
		for _, option := range parts[1:] {
			switch option {
			case "attr":
				isAttribute = true
			case "any":
				isAny = true
			case "chardata", "innerxml", "comment":
				isText = true
			}
		}

		if isAny {
			if isAttribute {
				allowAnyAttribute = true
			} else {
				allowAnyElement = true
			}
			continue
		}

		if isText {
			continue
		}

		if name == "" {
			name = field.Name
		}
		if isAttribute {
			attributes[name] = struct{}{}
			continue
		}

		elements[name] = field.Type
	}

	for _, attribute := range node.attrs {
		if attribute.Name.Space == "xmlns" || attribute.Name.Local == "xmlns" {
			continue
		}
		if _, ok := attributes[attribute.Name.Local]; !ok && !allowAnyAttribute {
			*unknown = append(*unknown, path+".@"+attribute.Name.Local)
		}
	}

	for _, child := range node.children {
		fieldType, ok := elements[child.name]
		childPath := path + "." + child.name
		if !ok {
			if !allowAnyElement {
				*unknown = append(*unknown, childPath)
			}
			continue
		}

		collectUnknownXMLFields(child, fieldType, childPath, unknown)
	}
}

// collectUnknownJSONFields walks JSON objects using model JSON tags.
func collectUnknownJSONFields(value any, model reflect.Type, path string, unknown *[]string) {
	for model.Kind() == reflect.Pointer {
		model = model.Elem()
	}

	switch typed := value.(type) {
	case map[string]any:
		if model.Kind() != reflect.Struct {
			return
		}

		fields := jsonFields(model)
		for name, child := range typed {
			fieldType, ok := fields[name]
			childPath := path + "." + name
			if !ok {
				*unknown = append(*unknown, childPath)
				continue
			}

			collectUnknownJSONFields(child, fieldType, childPath, unknown)
		}

	case []any:
		if model.Kind() != reflect.Slice && model.Kind() != reflect.Array {
			return
		}

		model = model.Elem()
		for index, child := range typed {
			collectUnknownJSONFields(child, model, fmt.Sprintf("%s[%d]", path, index), unknown)
		}
	}
}

// jsonFields indexes exported struct fields by their JSON property names.
func jsonFields(model reflect.Type) map[string]reflect.Type {
	fields := make(map[string]reflect.Type, model.NumField())
	for index := 0; index < model.NumField(); index++ {
		field := model.Field(index)
		if !field.IsExported() {
			continue
		}

		tag := strings.Split(field.Tag.Get("json"), ",")[0]
		if tag == "-" {
			continue
		}
		if tag == "" {
			tag = field.Name
		}
		fields[tag] = field.Type
	}

	return fields
}
