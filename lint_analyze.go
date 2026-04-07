// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/woozymasta/lintkit/lint"
)

const (
	// globalsVarLootDamageMin stores globals.xml key name.
	globalsVarLootDamageMin = "lootdamagemin"

	// globalsVarLootDamageMax stores globals.xml key name.
	globalsVarLootDamageMax = "lootdamagemax"

	// globalsVarTimeLogin stores globals.xml key name.
	globalsVarTimeLogin = "timelogin"

	// globalsVarTimeLogout stores globals.xml key name.
	globalsVarTimeLogout = "timelogout"

	// globalsVarInitialSpawn stores globals.xml key name.
	globalsVarInitialSpawn = "initialspawn"

	// globalsVarRestartSpawn stores globals.xml key name.
	globalsVarRestartSpawn = "restartspawn"

	// globalsVarAnimalMaxCount stores globals.xml key name.
	globalsVarAnimalMaxCount = "animalmaxcount"

	// globalsVarZombieMaxCount stores globals.xml key name.
	globalsVarZombieMaxCount = "zombiemaxcount"

	// globalsVarRespawnLimit stores globals.xml key name.
	globalsVarRespawnLimit = "respawnlimit"

	// globalsVarRespawnTypes stores globals.xml key name.
	globalsVarRespawnTypes = "respawntypes"

	// economyCoreDefaultBackupPeriod stores defaults key.
	economyCoreDefaultBackupPeriod = "backup_period"

	// economyCoreDefaultBackupCount stores defaults key.
	economyCoreDefaultBackupCount = "backup_count"

	// economyCoreDefaultWorldSegments stores defaults key.
	economyCoreDefaultWorldSegments = "world_segments"

	// economyCoreDefaultDynRadius stores defaults key.
	economyCoreDefaultDynRadius = "dyn_radius"

	// economyCoreDefaultDynSMin stores defaults key.
	economyCoreDefaultDynSMin = "dyn_smin"

	// economyCoreDefaultDynSMax stores defaults key.
	economyCoreDefaultDynSMax = "dyn_smax"

	// economyCoreDefaultDynDMin stores defaults key.
	economyCoreDefaultDynDMin = "dyn_dmin"

	// economyCoreDefaultDynDMax stores defaults key.
	economyCoreDefaultDynDMax = "dyn_dmax"
)

// globalsRangeSpec stores baseline range constraints for one globals key.
type globalsRangeSpec struct {
	// min stores minimum accepted value.
	min float64

	// max stores maximum accepted value.
	max float64

	// hasMin reports whether minimum bound is enabled.
	hasMin bool

	// hasMax reports whether maximum bound is enabled.
	hasMax bool
}

var (
	// globalsRangeBaseline stores known range checks by lower-cased key.
	globalsRangeBaseline = map[string]globalsRangeSpec{
		globalsVarLootDamageMin: {
			min:    0,
			max:    1,
			hasMin: true,
			hasMax: true,
		},
		globalsVarLootDamageMax: {
			min:    0,
			max:    1,
			hasMin: true,
			hasMax: true,
		},
		globalsVarTimeLogin: {
			min:    0,
			max:    65536,
			hasMin: true,
			hasMax: true,
		},
		globalsVarTimeLogout: {
			min:    0,
			max:    65536,
			hasMin: true,
			hasMax: true,
		},
		globalsVarInitialSpawn: {
			min:    0,
			max:    100,
			hasMin: true,
			hasMax: true,
		},
		globalsVarRestartSpawn: {
			min:    0,
			max:    100,
			hasMin: true,
			hasMax: true,
		},
		globalsVarAnimalMaxCount: {
			min:    0,
			hasMin: true,
		},
		globalsVarZombieMaxCount: {
			min:    0,
			hasMin: true,
		},
		globalsVarRespawnLimit: {
			min:    0,
			hasMin: true,
		},
		globalsVarRespawnTypes: {
			min:    0,
			hasMin: true,
		},
	}

	// eventsPositionValues stores supported events.xml position values.
	eventsPositionValues = map[string]struct{}{
		"fixed":   {},
		"player":  {},
		"uniform": {},
	}

	// eventsLimitValues stores supported events.xml limit values.
	eventsLimitValues = map[string]struct{}{
		"child":  {},
		"custom": {},
		"mixed":  {},
		"parent": {},
	}
)

// AnalyzeLintContent runs first-pass CE lint checks for one file payload.
func AnalyzeLintContent(path string, data []byte) []lint.Diagnostic {
	normalizedPath := strings.TrimSpace(path)
	kind := DetectKind(normalizedPath)
	if kind == KindUnknown {
		kind = detectKindByContent(normalizedPath, data)
	}

	if kind == KindUnknown {
		return analyzeUnknownKindXML(normalizedPath, data)
	}

	value, err := Decode(kind, data)
	if err != nil {
		if !isXMLKind(kind) {
			return nil
		}

		return []lint.Diagnostic{
			newDiagnostic(
				CodeParseInvalidXML,
				normalizedPath,
				fmt.Sprintf("failed to decode XML kind %s: %v", kind, err),
			),
		}
	}

	switch typed := value.(type) {
	case *TypesFile:
		return analyzeTypesFile(normalizedPath, typed)
	case *EventsFile:
		return analyzeEventsFile(normalizedPath, typed)
	case *GlobalsFile:
		return analyzeGlobalsFile(normalizedPath, typed)
	case *EconomyCoreFile:
		return analyzeEconomyCoreFile(normalizedPath, typed)
	case *EconomyFile:
		return analyzeEconomyFile(normalizedPath, typed)
	case *SpawnableTypesFile:
		return analyzeSpawnableTypesFile(normalizedPath, typed)
	default:
		return nil
	}
}

// analyzeUnknownKindXML tries to classify unknown XML payload parse issues.
func analyzeUnknownKindXML(path string, data []byte) []lint.Diagnostic {
	if strings.ToLower(strings.TrimSpace(fileExt(path))) != ".xml" {
		return nil
	}

	var root struct {
		// XMLName stores root tag name.
		XMLName xml.Name
	}

	if err := xml.Unmarshal(data, &root); err != nil {
		return []lint.Diagnostic{
			newDiagnostic(
				CodeParseInvalidXML,
				path,
				fmt.Sprintf("failed to parse XML payload: %v", err),
			),
		}
	}

	return []lint.Diagnostic{
		newDiagnostic(
			CodeParseUnknownRoot,
			path,
			fmt.Sprintf("unsupported XML root <%s>", root.XMLName.Local),
		),
	}
}

// analyzeTypesFile runs types.xml semantic checks.
func analyzeTypesFile(path string, file *TypesFile) []lint.Diagnostic {
	if file == nil {
		return nil
	}

	seen := make(map[string]struct{}, len(file.Types))
	diagnostics := make([]lint.Diagnostic, 0, 16)

	for index := range file.Types {
		item := file.Types[index]
		nameKey := strings.ToLower(strings.TrimSpace(item.Name))
		if nameKey != "" {
			if _, exists := seen[nameKey]; exists {
				diagnostics = append(diagnostics, newDiagnostic(
					CodeTypesDuplicateName,
					path,
					fmt.Sprintf("types.xml has duplicate type name %q", item.Name),
				))
			} else {
				seen[nameKey] = struct{}{}
			}
		}

		if item.Nominal != nil && *item.Nominal < 0 {
			diagnostics = append(diagnostics, newDiagnostic(
				CodeTypesNominalNegative,
				path,
				fmt.Sprintf("type %q has negative nominal=%d", item.Name, *item.Nominal),
			))
		}

		diagnostics = append(
			diagnostics,
			checkTypeQuantity(path, item.Name, item.QuantityMin, item.QuantityMax)...,
		)
	}

	return diagnostics
}

// checkTypeQuantity validates quantmin/quantmax constraints.
func checkTypeQuantity(
	path string,
	typeName string,
	quantityMin *int,
	quantityMax *int,
) []lint.Diagnostic {
	diagnostics := make([]lint.Diagnostic, 0, 2)

	if quantityMin != nil && !isValidQuantityValue(*quantityMin) {
		diagnostics = append(diagnostics, newDiagnostic(
			CodeTypesQuantityRange,
			path,
			fmt.Sprintf("type %q has invalid quantmin=%d", typeName, *quantityMin),
		))
	}

	if quantityMax != nil && !isValidQuantityValue(*quantityMax) {
		diagnostics = append(diagnostics, newDiagnostic(
			CodeTypesQuantityRange,
			path,
			fmt.Sprintf("type %q has invalid quantmax=%d", typeName, *quantityMax),
		))
	}

	if quantityMin != nil &&
		quantityMax != nil &&
		*quantityMin >= 0 &&
		*quantityMax >= 0 &&
		*quantityMin > *quantityMax {
		diagnostics = append(diagnostics, newDiagnostic(
			CodeTypesQuantityRange,
			path,
			fmt.Sprintf(
				"type %q has quantmin=%d greater than quantmax=%d",
				typeName,
				*quantityMin,
				*quantityMax,
			),
		))
	}

	return diagnostics
}

// analyzeEventsFile runs events.xml semantic checks.
func analyzeEventsFile(path string, file *EventsFile) []lint.Diagnostic {
	if file == nil {
		return nil
	}

	seen := make(map[string]struct{}, len(file.Events))
	diagnostics := make([]lint.Diagnostic, 0, 16)

	for index := range file.Events {
		item := file.Events[index]
		nameKey := strings.ToLower(strings.TrimSpace(item.Name))
		if nameKey != "" {
			if _, exists := seen[nameKey]; exists {
				diagnostics = append(diagnostics, newDiagnostic(
					CodeEventsDuplicateName,
					path,
					fmt.Sprintf("events.xml has duplicate event name %q", item.Name),
				))
			} else {
				seen[nameKey] = struct{}{}
			}
		}

		if item.Active != nil && !isCanonicalBool(*item.Active) {
			diagnostics = append(diagnostics, newDiagnostic(
				CodeEventsFlagNonCanonical,
				path,
				fmt.Sprintf("event %q has non-canonical active=%d", item.Name, *item.Active),
			))
		}

		if item.Flags != nil {
			diagnostics = append(
				diagnostics,
				checkEventFlag(path, item.Name, "deletable", item.Flags.Deletable)...,
			)
			diagnostics = append(
				diagnostics,
				checkEventFlag(path, item.Name, "init_random", item.Flags.InitRandom)...,
			)
			diagnostics = append(
				diagnostics,
				checkEventFlag(path, item.Name, "remove_damaged", item.Flags.RemoveDamaged)...,
			)
		}

		if item.Position != nil {
			value := strings.ToLower(strings.TrimSpace(*item.Position))
			if value != "" {
				if _, ok := eventsPositionValues[value]; !ok {
					diagnostics = append(diagnostics, newDiagnostic(
						CodeEventsUnknownPosition,
						path,
						fmt.Sprintf("event %q has unknown position=%q", item.Name, *item.Position),
					))
				}
			}
		}

		if item.Limit != nil {
			value := strings.ToLower(strings.TrimSpace(*item.Limit))
			if value != "" {
				if _, ok := eventsLimitValues[value]; !ok {
					diagnostics = append(diagnostics, newDiagnostic(
						CodeEventsUnknownLimit,
						path,
						fmt.Sprintf("event %q has unknown limit=%q", item.Name, *item.Limit),
					))
				}
			}
		}
	}

	return diagnostics
}

// checkEventFlag validates one event flags integer value.
func checkEventFlag(
	path string,
	eventName string,
	flagName string,
	value *int,
) []lint.Diagnostic {
	if value == nil || isCanonicalBool(*value) {
		return nil
	}

	return []lint.Diagnostic{
		newDiagnostic(
			CodeEventsFlagNonCanonical,
			path,
			fmt.Sprintf(
				"event %q has non-canonical flags.%s=%d",
				eventName,
				flagName,
				*value,
			),
		),
	}
}

// analyzeGlobalsFile runs globals.xml semantic checks.
func analyzeGlobalsFile(path string, file *GlobalsFile) []lint.Diagnostic {
	if file == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, 16)

	for index := range file.Vars {
		item := file.Vars[index]
		nameKey := strings.ToLower(strings.TrimSpace(item.Name))
		rawValue := strings.TrimSpace(item.Value)

		switch item.Type {
		case VariableTypeInt:
			parsed, err := strconv.Atoi(rawValue)
			if err != nil {
				diagnostics = append(diagnostics, newDiagnostic(
					CodeGlobalsValueTypeMismatch,
					path,
					fmt.Sprintf(
						"globals var %q has type=0 but value %q is not int",
						item.Name,
						item.Value,
					),
				))
				continue
			}

			diagnostics = append(
				diagnostics,
				checkGlobalsRange(path, item.Name, nameKey, float64(parsed))...,
			)

		case VariableTypeFloat:
			parsed, err := strconv.ParseFloat(rawValue, 64)
			if err != nil {
				diagnostics = append(diagnostics, newDiagnostic(
					CodeGlobalsValueTypeMismatch,
					path,
					fmt.Sprintf(
						"globals var %q has type=1 but value %q is not float",
						item.Name,
						item.Value,
					),
				))
				continue
			}

			diagnostics = append(
				diagnostics,
				checkGlobalsRange(path, item.Name, nameKey, parsed)...,
			)

		case VariableTypeString:
			continue

		default:
			diagnostics = append(diagnostics, newDiagnostic(
				CodeGlobalsInvalidTypeTag,
				path,
				fmt.Sprintf(
					"globals var %q has unsupported type=%d",
					item.Name,
					item.Type,
				),
			))
		}
	}

	return diagnostics
}

// analyzeSpawnableTypesFile runs cfgspawnabletypes.xml semantic checks.
func analyzeSpawnableTypesFile(
	path string,
	file *SpawnableTypesFile,
) []lint.Diagnostic {
	if file == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, 16)
	chanceValues := make([]float64, 0, 64)

	for index := range file.Types {
		item := file.Types[index]

		if item.Damage != nil {
			diagnostics = append(
				diagnostics,
				checkSpawnableMinMax(path, item.Name, "damage", item.Damage)...,
			)
		}

		for cargoIndex := range item.Cargo {
			if item.Cargo[cargoIndex].Chance != nil {
				chanceValues = append(chanceValues, *item.Cargo[cargoIndex].Chance)
			}

			for childIndex := range item.Cargo[cargoIndex].Items {
				if item.Cargo[cargoIndex].Items[childIndex].Chance != nil {
					chanceValues = append(
						chanceValues,
						*item.Cargo[cargoIndex].Items[childIndex].Chance,
					)
				}
			}
		}

		for attachmentIndex := range item.Attachments {
			if item.Attachments[attachmentIndex].Chance != nil {
				chanceValues = append(
					chanceValues,
					*item.Attachments[attachmentIndex].Chance,
				)
			}

			for childIndex := range item.Attachments[attachmentIndex].Items {
				if item.Attachments[attachmentIndex].Items[childIndex].Chance != nil {
					chanceValues = append(
						chanceValues,
						*item.Attachments[attachmentIndex].Items[childIndex].Chance,
					)
				}
			}
		}
	}

	diagnostics = append(diagnostics, checkSpawnableChance(path, chanceValues)...)

	return diagnostics
}

// checkSpawnableMinMax validates spawnable min/max in 0..1 interval.
func checkSpawnableMinMax(
	path string,
	typeName string,
	field string,
	value *SpawnableMinMax,
) []lint.Diagnostic {
	if value == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, 3)
	if value.Min != nil && (*value.Min < 0 || *value.Min > 1) {
		diagnostics = append(diagnostics, newDiagnostic(
			CodeSpawnableDamageRange,
			path,
			fmt.Sprintf(
				"type %q has %s.min=%v outside 0..1",
				typeName,
				field,
				*value.Min,
			),
		))
	}

	if value.Max != nil && (*value.Max < 0 || *value.Max > 1) {
		diagnostics = append(diagnostics, newDiagnostic(
			CodeSpawnableDamageRange,
			path,
			fmt.Sprintf(
				"type %q has %s.max=%v outside 0..1",
				typeName,
				field,
				*value.Max,
			),
		))
	}

	if value.Min != nil && value.Max != nil && *value.Min > *value.Max {
		diagnostics = append(diagnostics, newDiagnostic(
			CodeSpawnableDamageRange,
			path,
			fmt.Sprintf(
				"type %q has %s.min=%v greater than %s.max=%v",
				typeName,
				field,
				*value.Min,
				field,
				*value.Max,
			),
		))
	}

	return diagnostics
}

// checkSpawnableChance validates chance values in 0..1 or 0..100 mode.
func checkSpawnableChance(path string, values []float64) []lint.Diagnostic {
	if len(values) == 0 {
		return nil
	}

	percentMode := false
	for index := range values {
		if values[index] > 1 {
			percentMode = true
			break
		}
	}

	maxAllowed := 1.0
	if percentMode {
		maxAllowed = 100
	}

	diagnostics := make([]lint.Diagnostic, 0, 2)
	for index := range values {
		value := values[index]
		if value >= 0 && value <= maxAllowed {
			continue
		}

		diagnostics = append(diagnostics, newDiagnostic(
			CodeSpawnableChanceRange,
			path,
			fmt.Sprintf(
				"chance value %v is outside allowed range 0..%v",
				value,
				maxAllowed,
			),
		))
	}

	return diagnostics
}

// checkGlobalsRange validates one parsed globals value against baseline range.
func checkGlobalsRange(
	path string,
	name string,
	nameKey string,
	value float64,
) []lint.Diagnostic {
	spec, ok := globalsRangeBaseline[nameKey]
	if !ok {
		return nil
	}

	if spec.hasMin && value < spec.min {
		return []lint.Diagnostic{
			newDiagnostic(
				CodeGlobalsOutOfRange,
				path,
				fmt.Sprintf(
					"globals var %q value %v is less than %v",
					name,
					value,
					spec.min,
				),
			),
		}
	}

	if spec.hasMax && value > spec.max {
		return []lint.Diagnostic{
			newDiagnostic(
				CodeGlobalsOutOfRange,
				path,
				fmt.Sprintf(
					"globals var %q value %v is greater than %v",
					name,
					value,
					spec.max,
				),
			),
		}
	}

	return nil
}

// analyzeEconomyCoreFile checks duplicate default names in economycore.
func analyzeEconomyCoreFile(path string, file *EconomyCoreFile) []lint.Diagnostic {
	if file == nil || file.Defaults == nil || len(file.Defaults.Defaults) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(file.Defaults.Defaults))
	parsed := make(map[string]float64, len(file.Defaults.Defaults))
	diagnostics := make([]lint.Diagnostic, 0, 8)

	for index := range file.Defaults.Defaults {
		name := file.Defaults.Defaults[index].Name
		value := file.Defaults.Defaults[index].Value
		nameKey := strings.ToLower(strings.TrimSpace(name))
		if nameKey == "" {
			continue
		}

		if _, exists := seen[nameKey]; exists {
			diagnostics = append(diagnostics, newDiagnostic(
				CodeEconomyCoreDuplicateDefaultName,
				path,
				fmt.Sprintf(
					"cfgeconomycore duplicate default name %q",
					name,
				),
			))
			continue
		}

		seen[nameKey] = struct{}{}

		number, ok := parseEconomyCoreNumericDefault(nameKey, value)
		if !ok {
			continue
		}

		parsed[nameKey] = number
		diagnostics = append(
			diagnostics,
			checkEconomyCoreDefaultRange(path, nameKey, number)...,
		)
	}

	diagnostics = append(diagnostics, checkEconomyCoreDefaultWindow(path, parsed)...)

	return diagnostics
}

// parseEconomyCoreNumericDefault parses known numeric defaults as float64.
func parseEconomyCoreNumericDefault(nameKey string, value string) (float64, bool) {
	switch nameKey {
	case economyCoreDefaultBackupPeriod,
		economyCoreDefaultBackupCount,
		economyCoreDefaultWorldSegments,
		economyCoreDefaultDynRadius,
		economyCoreDefaultDynSMin,
		economyCoreDefaultDynSMax,
		economyCoreDefaultDynDMin,
		economyCoreDefaultDynDMax:
	default:
		return 0, false
	}

	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0, false
	}

	return parsed, true
}

// checkEconomyCoreDefaultRange validates one numeric default against baseline.
func checkEconomyCoreDefaultRange(
	path string,
	nameKey string,
	value float64,
) []lint.Diagnostic {
	switch nameKey {
	case economyCoreDefaultBackupPeriod:
		if value < 15 {
			return []lint.Diagnostic{
				newDiagnostic(
					CodeEconomyCoreDefaultOutOfRange,
					path,
					fmt.Sprintf("%s=%v must be >= 15", nameKey, value),
				),
			}
		}
	case economyCoreDefaultBackupCount, economyCoreDefaultWorldSegments:
		if value < 1 {
			return []lint.Diagnostic{
				newDiagnostic(
					CodeEconomyCoreDefaultOutOfRange,
					path,
					fmt.Sprintf("%s=%v must be >= 1", nameKey, value),
				),
			}
		}
	case economyCoreDefaultDynRadius:
		if value < 0 {
			return []lint.Diagnostic{
				newDiagnostic(
					CodeEconomyCoreDefaultOutOfRange,
					path,
					fmt.Sprintf("%s=%v must be >= 0", nameKey, value),
				),
			}
		}
	}

	return nil
}

// checkEconomyCoreDefaultWindow validates paired min/max defaults.
func checkEconomyCoreDefaultWindow(
	path string,
	parsed map[string]float64,
) []lint.Diagnostic {
	diagnostics := make([]lint.Diagnostic, 0, 2)
	if minValue, okMin := parsed[economyCoreDefaultDynSMin]; okMin {
		if maxValue, okMax := parsed[economyCoreDefaultDynSMax]; okMax && minValue > maxValue {
			diagnostics = append(diagnostics, newDiagnostic(
				CodeEconomyCoreDefaultOutOfRange,
				path,
				fmt.Sprintf(
					"%s=%v must be <= %s=%v",
					economyCoreDefaultDynSMin,
					minValue,
					economyCoreDefaultDynSMax,
					maxValue,
				),
			))
		}
	}

	if minValue, okMin := parsed[economyCoreDefaultDynDMin]; okMin {
		if maxValue, okMax := parsed[economyCoreDefaultDynDMax]; okMax && minValue > maxValue {
			diagnostics = append(diagnostics, newDiagnostic(
				CodeEconomyCoreDefaultOutOfRange,
				path,
				fmt.Sprintf(
					"%s=%v must be <= %s=%v",
					economyCoreDefaultDynDMin,
					minValue,
					economyCoreDefaultDynDMax,
					maxValue,
				),
			))
		}
	}

	return diagnostics
}

// analyzeEconomyFile checks required sections and canonical 0/1 section flags.
func analyzeEconomyFile(path string, file *EconomyFile) []lint.Diagnostic {
	if file == nil {
		return nil
	}

	sections := []struct {
		// name stores XML section tag.
		name string

		// item stores one economy section payload.
		item *EconomySection
	}{
		{name: "dynamic", item: file.Dynamic},
		{name: "animals", item: file.Animals},
		{name: "zombies", item: file.Zombies},
		{name: "vehicles", item: file.Vehicles},
		{name: "randoms", item: file.Randoms},
		{name: "custom", item: file.Custom},
		{name: "building", item: file.Building},
		{name: "player", item: file.Player},
	}

	diagnostics := make([]lint.Diagnostic, 0, 8)
	for index := range sections {
		section := sections[index]
		if section.item == nil {
			diagnostics = append(diagnostics, newDiagnostic(
				CodeEconomyIncompleteSection,
				path,
				fmt.Sprintf("economy section <%s> is missing", section.name),
			))
			continue
		}

		diagnostics = append(
			diagnostics,
			checkEconomySectionFlags(path, section.name, section.item)...,
		)
	}

	return diagnostics
}

// checkEconomySectionFlags validates canonical section attrs for one section.
func checkEconomySectionFlags(
	path string,
	sectionName string,
	section *EconomySection,
) []lint.Diagnostic {
	if section == nil {
		return nil
	}

	flags := []struct {
		// name stores XML attribute token.
		name string

		// value stores parsed numeric flag value.
		value *int
	}{
		{name: "init", value: section.Init},
		{name: "load", value: section.Load},
		{name: "respawn", value: section.Respawn},
		{name: "save", value: section.Save},
	}

	diagnostics := make([]lint.Diagnostic, 0, 4)
	for index := range flags {
		item := flags[index]
		if item.value == nil {
			diagnostics = append(diagnostics, newDiagnostic(
				CodeEconomyIncompleteSection,
				path,
				fmt.Sprintf(
					"economy section <%s> is missing @%s",
					sectionName,
					item.name,
				),
			))
			continue
		}

		if *item.value == 0 || *item.value == 1 {
			continue
		}

		diagnostics = append(diagnostics, newDiagnostic(
			CodeEconomyIncompleteSection,
			path,
			fmt.Sprintf(
				"economy section <%s> has non-canonical @%s=%d",
				sectionName,
				item.name,
				*item.value,
			),
		))
	}

	return diagnostics
}

// isXMLKind reports whether one CE kind is XML-based.
func isXMLKind(kind Kind) bool {
	switch kind {
	case KindTypes,
		KindEvents,
		KindEconomy,
		KindGlobals,
		KindMessages,
		KindSpawnableTypes,
		KindRandomPresets,
		KindEconomyCore,
		KindEnvironment,
		KindEventSpawns,
		KindEventGroups,
		KindPlayerSpawnPoints,
		KindWeather,
		KindLimitsDefinition,
		KindLimitsDefinitionUser,
		KindIgnoreList,
		KindTerritories,
		KindCEProjectConfig,
		KindMapGroupProto,
		KindMapClusterProto,
		KindMapGroupPos,
		KindMapGroupDirt,
		KindMapGroupCluster:
		return true
	default:
		return false
	}
}

// fileExt extracts file extension from path in lower-level helper style.
func fileExt(path string) string {
	for index := len(path) - 1; index >= 0; index-- {
		char := path[index]
		if char == '.' {
			return path[index:]
		}
		if char == '\\' || char == '/' {
			break
		}
	}

	return ""
}

// isValidQuantityValue reports whether quantity value is -1 sentinel or 0..100.
func isValidQuantityValue(value int) bool {
	return value == -1 || (value >= 0 && value <= 100)
}

// isCanonicalBool reports whether value is canonical numeric boolean.
func isCanonicalBool(value int) bool {
	return value == 0 || value == 1
}

// newDiagnostic builds one shared lint diagnostic for dzce code token.
func newDiagnostic(code lint.Code, path string, message string) lint.Diagnostic {
	return lint.Diagnostic{
		Code:    lint.ApplyCodePrefix("DZCE", code),
		Path:    path,
		Message: strings.TrimSpace(message),
	}
}
