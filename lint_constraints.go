// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"fmt"
	"slices"
	"unicode/utf8"

	"github.com/woozymasta/lintkit/lint"
)

// analyzeMessagesFile validates numeric and enum constraints for message entries.
func analyzeMessagesFile(path string, file *MessagesFile) []lint.Diagnostic {
	if file == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, len(file.Messages))
	for index := range file.Messages {
		message := &file.Messages[index]
		fieldPrefix := fmt.Sprintf("message[%d]", index)
		diagnostics = append(diagnostics, checkIntRange(
			path, fieldPrefix+".deadline", message.Deadline, 0, 0, false,
		)...)
		diagnostics = append(diagnostics, checkIntEnum(
			path, fieldPrefix+".shutdown", message.Shutdown, 0, 1,
		)...)
		diagnostics = append(diagnostics, checkIntRange(
			path, fieldPrefix+".delay", message.Delay, 0, 0, false,
		)...)
		diagnostics = append(diagnostics, checkIntRange(
			path, fieldPrefix+".repeat", message.Repeat, 0, 0, false,
		)...)
		diagnostics = append(diagnostics, checkIntEnum(
			path, fieldPrefix+".onconnect", message.OnConnect, 0, 1,
		)...)
	}

	return diagnostics
}

// analyzeGameplayFile validates declared numeric, enum, and fixed-array constraints.
func analyzeGameplayFile(path string, file *GameplayFile) []lint.Diagnostic {
	if file == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, 16)
	if player := file.PlayerData; player != nil {
		if movement := player.MovementData; movement != nil {
			for _, field := range []struct {
				name  string
				value *float64
			}{
				{name: "timeToStrafeJog", value: movement.TimeToStrafeJog},
				{name: "rotationSpeedJog", value: movement.RotationSpeedJog},
				{name: "timeToSprint", value: movement.TimeToSprint},
				{name: "timeToStrafeSprint", value: movement.TimeToStrafeSprint},
				{name: "rotationSpeedSprint", value: movement.RotationSpeedSprint},
			} {
				diagnostics = append(diagnostics, checkFloatRange(
					path,
					"PlayerData.MovementData."+field.name,
					field.value,
					0.01,
					0,
					false,
				)...)
			}
		}

		if obstruction := player.WeaponObstructionData; obstruction != nil {
			diagnostics = append(diagnostics, checkIntEnum(
				path, "PlayerData.WeaponObstructionData.staticMode", obstruction.StaticMode, 0, 1, 2,
			)...)
			diagnostics = append(diagnostics, checkIntEnum(
				path, "PlayerData.WeaponObstructionData.dynamicMode", obstruction.DynamicMode, 0, 1, 2,
			)...)
		}
	}

	if worlds := file.WorldsData; worlds != nil {
		for _, field := range []struct {
			name  string
			value []float64
			want  int
		}{
			{name: "WorldsData.environmentMinTemps", value: worlds.EnvironmentMinTemps, want: 12},
			{name: "WorldsData.environmentMaxTemps", value: worlds.EnvironmentMaxTemps, want: 12},
			{name: "WorldsData.wetnessWeightModifiers", value: worlds.WetnessWeightModifiers, want: 5},
		} {
			diagnostics = append(diagnostics, checkArrayLength(path, field.name, field.value, field.want)...)
		}
	}

	if ui := file.UIData; ui != nil && ui.HitIndicationData != nil {
		hit := ui.HitIndicationData

		diagnostics = append(diagnostics, checkIntEnum(
			path,
			"UIData.HitIndicationData.hitDirectionBehaviour",
			hit.HitDirectionBehaviour,
			0, 1, 2,
		)...)

		diagnostics = append(diagnostics, checkIntEnum(
			path,
			"UIData.HitIndicationData.hitDirectionStyle",
			hit.HitDirectionStyle,
			0, 1, 2,
		)...)

		diagnostics = append(diagnostics, checkFloatRange(
			path,
			"UIData.HitIndicationData.hitDirectionBreakPointRelative",
			hit.HitDirectionBreakPointRelative,
			0,
			1,
			true,
		)...)
	}

	if file.SandstormData != nil {
		diagnostics = append(diagnostics, checkFloatEnum(
			path,
			"SandstormData.sandstormFrequency",
			file.SandstormData.SandstormFrequency,
			0,
			1,
			2,
		)...)
	}

	return diagnostics
}

// analyzeBunkerBroadcastFile validates bunker settings and proxy coordinates.
func analyzeBunkerBroadcastFile(path string, file *BunkerBroadcastFile) []lint.Diagnostic {
	if file == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, len(file.Proxies)+6)
	if settings := file.BunkerSettings; settings != nil {
		diagnostics = append(diagnostics, checkIntRange(
			path, "bunkerSettings.numberOfActive", settings.NumberOfActive, 1, 30, true,
		)...)
		diagnostics = append(diagnostics, checkIntRange(
			path, "bunkerSettings.lockCodeLength", settings.LockCodeLength, 3, 6, true,
		)...)
		diagnostics = append(diagnostics, checkIntEnum(
			path, "bunkerSettings.lockCodeMode", settings.LockCodeMode, 0, 1,
		)...)
	}

	if settings := file.SchedulerSettings; settings != nil {
		diagnostics = append(diagnostics, checkIntRange(
			path, "schedulerSettings.batchLifetimeLength", settings.BatchLifetimeLength, 1, 1000, true,
		)...)
		diagnostics = append(diagnostics, checkIntRange(
			path, "schedulerSettings.reportInterval", settings.ReportInterval, 1, 120, true,
		)...)

		for index := range settings.Periods {
			period := &settings.Periods[index]
			prefix := fmt.Sprintf("schedulerSettings.periods[%d]", index)
			diagnostics = append(diagnostics,
				checkStringLength(path, prefix+".startTime", period.StartTime, 4, 4)...,
			)
			diagnostics = append(diagnostics,
				checkStringLength(path, prefix+".endTime", period.EndTime, 4, 4)...,
			)
		}
	}

	for index := range file.Proxies {
		diagnostics = append(diagnostics, checkArrayLength(
			path,
			fmt.Sprintf("proxies[%d].position", index),
			file.Proxies[index].Position,
			3,
		)...)
	}

	return diagnostics
}

// analyzeEconomyCoreEnums validates string enums declared by economycore models.
func analyzeEconomyCoreEnums(path string, file *EconomyCoreFile) []lint.Diagnostic {
	if file == nil || file.Classes == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, len(file.Classes.RootClasses)*2)
	for index := range file.Classes.RootClasses {
		class := &file.Classes.RootClasses[index]
		prefix := fmt.Sprintf("classes.rootclass[%d]", index)
		diagnostics = append(diagnostics,
			checkStringEnum(path, prefix+".reportMemoryLOD", class.ReportMemoryLOD, "yes", "no")...,
		)
		if class.Act != "" {
			diagnostics = append(diagnostics,
				checkStringEnum(path, prefix+".act", class.Act, "none", "character", "car")...,
			)
		}
	}

	return diagnostics
}

// analyzeGameplayGearPresetsFile validates weight and quickbar bounds in gear presets.
func analyzeGameplayGearPresetsFile(path string, presets *GameplayGearPresetsFile) []lint.Diagnostic {
	if presets == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, len(*presets))
	for presetIndex := range *presets {
		preset := &(*presets)[presetIndex]
		prefix := fmt.Sprintf("presets[%d]", presetIndex)
		diagnostics = append(diagnostics, checkIntRange(
			path, prefix+".spawnWeight", preset.SpawnWeight, 1, 0, false,
		)...)
		for slotIndex := range preset.AttachmentSlotItemSets {
			slot := &preset.AttachmentSlotItemSets[slotIndex]

			for itemIndex := range slot.DiscreteItemSets {
				item := &slot.DiscreteItemSets[itemIndex]
				itemPrefix := fmt.Sprintf("%s.attachmentSlotItemSets[%d].discreteItemSets[%d]", prefix, slotIndex, itemIndex)
				diagnostics = append(diagnostics, checkGearItemBounds(path, itemPrefix, item.SpawnWeight, item.QuickBarSlot)...)
				for childIndex := range item.ComplexChildrenTypes {
					child := &item.ComplexChildrenTypes[childIndex]
					childPrefix := fmt.Sprintf("%s.complexChildrenTypes[%d]", itemPrefix, childIndex)
					diagnostics = append(diagnostics, checkGearItemBounds(path, childPrefix, nil, child.QuickBarSlot)...)
				}
			}
		}

		for setIndex := range preset.DiscreteUnsortedItemSets {
			set := &preset.DiscreteUnsortedItemSets[setIndex]
			setPrefix := fmt.Sprintf("%s.discreteUnsortedItemSets[%d]", prefix, setIndex)
			diagnostics = append(diagnostics, checkGearItemBounds(path, setPrefix, set.SpawnWeight, nil)...)

			for childIndex := range set.ComplexChildrenTypes {
				child := &set.ComplexChildrenTypes[childIndex]
				childPrefix := fmt.Sprintf("%s.complexChildrenTypes[%d]", setPrefix, childIndex)
				diagnostics = append(diagnostics, checkGearItemBounds(path, childPrefix, nil, child.QuickBarSlot)...)
			}
		}
	}

	return diagnostics
}

// checkGearItemBounds validates constraints shared by weighted and nested gear entries.
func checkGearItemBounds(path string, prefix string, weight *int, quickBarSlot *int) []lint.Diagnostic {
	diagnostics := make([]lint.Diagnostic, 0, 2)
	diagnostics = append(diagnostics, checkIntRange(path, prefix+".spawnWeight", weight, 1, 0, false)...)
	diagnostics = append(diagnostics, checkIntRange(path, prefix+".quickBarSlot", quickBarSlot, -1, 0, false)...)

	return diagnostics
}

// analyzeEffectAreaFile validates fixed XYZ positions in effect-area entries.
func analyzeEffectAreaFile(path string, file *EffectAreaFile) []lint.Diagnostic {
	if file == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, len(file.Areas))
	for index := range file.Areas {
		area := &file.Areas[index]
		if area.Data == nil {
			continue
		}

		diagnostics = append(diagnostics, checkArrayLength(
			path,
			fmt.Sprintf("Areas[%d].Data.Pos", index),
			area.Data.Pos,
			3,
		)...)
	}

	return diagnostics
}

// analyzeUndergroundTriggersFile validates trigger ranges, enum flags, and XYZ vectors.
func analyzeUndergroundTriggersFile(path string, file *UndergroundTriggersFile) []lint.Diagnostic {
	if file == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, len(file.Triggers)*4)
	for triggerIndex := range file.Triggers {
		trigger := &file.Triggers[triggerIndex]
		prefix := fmt.Sprintf("Triggers[%d]", triggerIndex)

		diagnostics = append(diagnostics, checkFloatRange(
			path,
			prefix+".EyeAccommodation",
			trigger.EyeAccommodation,
			0,
			1,
			true,
		)...)

		diagnostics = append(diagnostics, checkIntEnum(
			path,
			prefix+".UseLinePointFade",
			trigger.UseLinePointFade,
			0, 1,
		)...)

		diagnostics = append(diagnostics, checkArrayLength(
			path,
			prefix+".Position",
			trigger.Position, 3,
		)...)

		diagnostics = append(diagnostics, checkArrayLength(
			path,
			prefix+".Orientation",
			trigger.Orientation, 3,
		)...)

		diagnostics = append(diagnostics, checkArrayLength(
			path,
			prefix+".Size",
			trigger.Size, 3,
		)...)

		for breadcrumbIndex := range trigger.Breadcrumbs {
			breadcrumb := &trigger.Breadcrumbs[breadcrumbIndex]
			breadcrumbPrefix := fmt.Sprintf("%s.Breadcrumbs[%d]", prefix, breadcrumbIndex)

			diagnostics = append(diagnostics, checkFloatRange(
				path,
				breadcrumbPrefix+".EyeAccommodation",
				breadcrumb.EyeAccommodation,
				0,
				1,
				true,
			)...)

			diagnostics = append(diagnostics, checkIntEnum(
				path,
				breadcrumbPrefix+".UseRaycast",
				breadcrumb.UseRaycast,
				0,
				1,
			)...)

			diagnostics = append(diagnostics, checkArrayLength(
				path,
				breadcrumbPrefix+".Position",
				breadcrumb.Position,
				3,
			)...)
		}
	}

	return diagnostics
}

// analyzeObjectSpawnerFile validates object position and optional rotation vector lengths.
func analyzeObjectSpawnerFile(path string, file *ObjectSpawnerFile) []lint.Diagnostic {
	if file == nil {
		return nil
	}

	diagnostics := make([]lint.Diagnostic, 0, len(file.Objects))
	for index := range file.Objects {
		object := &file.Objects[index]
		diagnostics = append(diagnostics, checkArrayLength(
			path,
			fmt.Sprintf("Objects[%d].pos", index),
			object.Pos,
			3,
		)...)

		diagnostics = append(diagnostics, checkOptionalArrayLength(
			path,
			fmt.Sprintf("Objects[%d].ypr", index),
			object.YPR,
			3,
		)...)
	}

	return diagnostics
}

// checkOptionalArrayLength allows an empty vector or one with exactly the expected length.
func checkOptionalArrayLength(path string, field string, values []float64, expected int) []lint.Diagnostic {
	if len(values) == 0 || len(values) == expected {
		return nil
	}

	return []lint.Diagnostic{newDiagnostic(
		CodeValidateInvalidLength,
		path,
		fmt.Sprintf("%s has %d values, want 0 or exactly %d", field, len(values), expected),
	)}
}

// checkRandomPresetItemChances validates chance constraints throughout a preset tree.
func checkRandomPresetItemChances(
	path string,
	group string,
	presetName string,
	items []SpawnableItem,
) []lint.Diagnostic {
	diagnostics := make([]lint.Diagnostic, 0, len(items))
	for index := range items {
		item := &items[index]
		if item.Chance != nil && (!isFiniteFloat64(*item.Chance) || *item.Chance < 0) {
			diagnostics = append(diagnostics, newDiagnostic(
				CodeSpawnableChanceRange,
				path,
				fmt.Sprintf(
					"randompresets %s preset %q item %q chance=%v outside finite non-negative range",
					group,
					presetName,
					item.Name,
					*item.Chance,
				),
			))
		}

		for cargoIndex := range item.Cargo {
			cargo := &item.Cargo[cargoIndex]
			if cargo.Chance != nil {
				diagnostics = append(diagnostics, checkSpawnableChance(
					path,
					fmt.Sprintf("randompresets %s preset %q item %q", group, presetName, item.Name),
					"nested cargo",
					*cargo.Chance,
				)...)
			}
			diagnostics = append(diagnostics, checkRandomPresetItemChances(
				path,
				group,
				presetName,
				cargo.Items,
			)...)
		}

		for attachmentIndex := range item.Attachments {
			attachment := &item.Attachments[attachmentIndex]
			if attachment.Chance != nil {
				diagnostics = append(diagnostics, checkSpawnableChance(
					path,
					fmt.Sprintf("randompresets %s preset %q item %q", group, presetName, item.Name),
					"nested attachment",
					*attachment.Chance,
				)...)
			}
			diagnostics = append(diagnostics, checkRandomPresetItemChances(
				path,
				group,
				presetName,
				attachment.Items,
			)...)
		}
	}

	return diagnostics
}

// checkNestedSpawnableChances validates chance constraints below one spawnable group.
func checkNestedSpawnableChances(
	path string,
	typeName string,
	items []SpawnableItem,
) []lint.Diagnostic {
	diagnostics := make([]lint.Diagnostic, 0, len(items))
	for index := range items {
		item := &items[index]
		for cargoIndex := range item.Cargo {
			cargo := &item.Cargo[cargoIndex]
			if cargo.Chance != nil {
				diagnostics = append(diagnostics, checkSpawnableChance(
					path,
					typeName,
					"nested cargo",
					*cargo.Chance,
				)...)
			}

			diagnostics = append(diagnostics,
				checkSpawnableItemChances(path, typeName, "nested cargo", cargo.Items)...,
			)
			diagnostics = append(diagnostics,
				checkNestedSpawnableChances(path, typeName, cargo.Items)...,
			)
		}

		for attachmentIndex := range item.Attachments {
			attachment := &item.Attachments[attachmentIndex]
			if attachment.Chance != nil {
				diagnostics = append(diagnostics, checkSpawnableChance(
					path,
					typeName,
					"nested attachments",
					*attachment.Chance,
				)...)
			}

			diagnostics = append(diagnostics,
				checkSpawnableItemChances(path, typeName, "nested attachments", attachment.Items)...,
			)
			diagnostics = append(diagnostics,
				checkNestedSpawnableChances(path, typeName, attachment.Items)...,
			)
		}
	}

	return diagnostics
}

// checkIntRange reports an integer below its minimum or above an optional maximum.
func checkIntRange(
	path string,
	field string,
	value *int,
	minimum int,
	maximum int,
	hasMaximum bool,
) []lint.Diagnostic {
	if value == nil || *value >= minimum && (!hasMaximum || *value <= maximum) {
		return nil
	}

	constraint := ""
	if hasMaximum {
		constraint = fmt.Sprintf("%d..%d", minimum, maximum)
	} else {
		constraint = fmt.Sprintf("at least %d", minimum)
	}

	return []lint.Diagnostic{newDiagnostic(
		CodeValidateInvalidIntRange,
		path,
		fmt.Sprintf("%s=%d is outside allowed range %s", field, *value, constraint),
	)}
}

// checkFloatRange reports a floating-point value outside optional inclusive bounds.
func checkFloatRange(
	path string,
	field string,
	value *float64,
	minimum float64,
	maximum float64,
	hasMaximum bool,
) []lint.Diagnostic {
	if value == nil || isFiniteFloat64(*value) &&
		*value >= minimum && (!hasMaximum || *value <= maximum) {
		return nil
	}

	constraint := ""
	if hasMaximum {
		constraint = fmt.Sprintf("%v..%v", minimum, maximum)
	} else {
		constraint = fmt.Sprintf("at least %v", minimum)
	}

	return []lint.Diagnostic{newDiagnostic(
		CodeValidateNumericRange,
		path,
		fmt.Sprintf("%s=%v is outside allowed range %s", field, *value, constraint),
	)}
}

// checkIntEnum reports an integer outside the declared enum values.
func checkIntEnum(path string, field string, value *int, allowed ...int) []lint.Diagnostic {
	if value == nil {
		return nil
	}

	if slices.Contains(allowed, *value) {
		return nil
	}

	return []lint.Diagnostic{newDiagnostic(
		CodeValidateUnknownEnum,
		path,
		fmt.Sprintf("%s=%d is not one of %v", field, *value, allowed),
	)}
}

// checkFloatEnum reports a floating-point value outside the declared enum values.
func checkFloatEnum(path string, field string, value *float64, allowed ...float64) []lint.Diagnostic {
	if value == nil {
		return nil
	}

	if slices.Contains(allowed, *value) {
		return nil
	}

	return []lint.Diagnostic{newDiagnostic(
		CodeValidateUnknownEnum,
		path,
		fmt.Sprintf("%s=%v is not one of %v", field, *value, allowed),
	)}
}

// checkStringEnum reports a string outside the declared enum values.
func checkStringEnum(path string, field string, value string, allowed ...string) []lint.Diagnostic {
	if value == "" || slices.Contains(allowed, value) {
		return nil
	}

	return []lint.Diagnostic{newDiagnostic(
		CodeValidateUnknownEnum,
		path,
		fmt.Sprintf("%s=%q is not one of %v", field, value, allowed),
	)}
}

// checkStringLength reports a present string outside inclusive rune-count bounds.
func checkStringLength(path string, field string, value string, minimum int, maximum int) []lint.Diagnostic {
	if value == "" {
		return nil
	}

	length := utf8.RuneCountInString(value)
	if length >= minimum && length <= maximum {
		return nil
	}

	return []lint.Diagnostic{newDiagnostic(
		CodeValidateInvalidLength,
		path,
		fmt.Sprintf("%s has %d characters, want %d..%d", field, length, minimum, maximum),
	)}
}

// checkArrayLength reports a present array whose length differs from its schema size.
func checkArrayLength(path string, field string, values []float64, expected int) []lint.Diagnostic {
	if values == nil || len(values) == expected {
		return nil
	}

	return []lint.Diagnostic{newDiagnostic(
		CodeValidateInvalidLength,
		path,
		fmt.Sprintf("%s has %d values, want exactly %d", field, len(values), expected),
	)}
}
