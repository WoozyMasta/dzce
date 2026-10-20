<!-- Automatically generated file, do not modify! -->

# Lint Rules Registry

This document contains the current registry of lint rules.

Total rules: 34.

## dzce

DayZ CE

> Lint rules for DayZ central economy configuration files.

Rule groups for `dzce`:

* [crossref](#crossref)
* [merge](#merge)
* [parse](#parse)
* [semantic](#semantic)
* [validate](#validate)

### crossref

> Cross-file reference diagnostics for merged CE trees.

Codes:
[DZCE4001](#dzce4001),
[DZCE4002](#dzce4002),

#### `DZCE4001`

Missing type reference

> A merged CE event child references a type missing in final merged
> `types.xml`. Add the missing type definition or fix the reference.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.crossref.missing-type-reference` |
| Scope | `crossref` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `DZCE4002`

Missing event reference

> A merged `events.xml` event.secondary references an event name missing in
> final merged `events.xml`.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.crossref.missing-event-reference` |
| Scope | `crossref` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

### merge

> Include graph and merge diagnostics for economycore trees.

Codes:
[DZCE4003](#dzce4003),
[DZCE4004](#dzce4004),
[DZCE4005](#dzce4005),

#### `DZCE4003`

Include cycle detected

> `cfgeconomycore.xml` include graph has a recursive cycle. Break the cycle to
> get deterministic CE merge order.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.merge.include-cycle-detected` |
| Scope | `merge` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `DZCE4004`

Include file not found

> `cfgeconomycore.xml` references an include file that cannot be found at
> resolved path. Fix CE include folder/name/type or restore the file.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.merge.include-file-not-found` |
| Scope | `merge` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `DZCE4005`

Duplicate type override across includes

> The same CE type name is defined in multiple included `types.xml` files.
> This can be intentional override behavior, but verify final include
> priority/order.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.merge.duplicate-type-override-across-includes` |
| Scope | `merge` |
| Severity | `notice` |
| Enabled | `true` (implicit) |

### parse

> Parse diagnostics for CE input payloads.

Codes:
[DZCE1001](#dzce1001),
[DZCE1002](#dzce1002),

#### `DZCE1001`

Xml decode failed

> CE XML file is malformed and cannot be parsed. Check broken tags, invalid
> attributes, and XML syntax.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.parse.xml-decode-failed` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `DZCE1002`

Unsupported xml root

> XML is valid, but root tag is not a supported CE file model for this check
> set. Use a supported CE root or exclude this file from CE lint input.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.parse.unsupported-xml-root` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

### semantic

> Semantic validation diagnostics for decoded CE models.

Codes:
[DZCE3001](#dzce3001),
[DZCE3002](#dzce3002),
[DZCE3003](#dzce3003),
[DZCE3004](#dzce3004),
[DZCE3005](#dzce3005),
[DZCE3102](#dzce3102),
[DZCE3103](#dzce3103),
[DZCE3104](#dzce3104),
[DZCE3202](#dzce3202),
[DZCE3203](#dzce3203),
[DZCE3204](#dzce3204),
[DZCE3301](#dzce3301),
[DZCE3302](#dzce3302),
[DZCE3303](#dzce3303),
[DZCE3304](#dzce3304),
[DZCE3305](#dzce3305),
[DZCE3601](#dzce3601),
[DZCE3604](#dzce3604),
[DZCE3605](#dzce3605),
[DZCE3606](#dzce3606),
[DZCE3702](#dzce3702),
[DZCE3703](#dzce3703),

#### `DZCE3001`

Duplicate type name

> `types.xml` contains multiple `<type>` entries with the same name. Keep one
> canonical CE type definition, or verify override order is intentional.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.duplicate-type-name` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3002`

Type nominal is negative

> `types.xml` has `<nominal>` below zero. In CE this value is usually expected
> to be 0 or greater.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.type-nominal-is-negative` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3003`

Type min is greater than nominal

> `types.xml` has `<min>` larger than `<nominal>`. This can be intentional in
> some setups, but often indicates inconsistent balancing values.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.type-min-is-greater-than-nominal` |
| Scope | `semantic` |
| Severity | `notice` |
| Enabled | `true` (implicit) |

#### `DZCE3004`

Event limit window looks inconsistent

> `events.xml` has potentially inconsistent numeric window values (for example
> min > nominal or max < min). Review this event configuration.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.event-limit-window-looks-inconsistent` |
| Scope | `semantic` |
| Severity | `notice` |
| Enabled | `true` (implicit) |

#### `DZCE3005`

Duplicate spawnable child entry

> `cfgspawnabletypes.xml` contains duplicate child item names in the same
> parent block. Remove duplicates unless intentional weighted duplication is
> required.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.duplicate-spawnable-child-entry` |
| Scope | `semantic` |
| Severity | `notice` |
| Enabled | `true` (implicit) |

#### `DZCE3102`

Globals var type tag is invalid

> `globals.xml` var@type uses an unsupported CE type tag. Allowed values are
> 0, 1, 2.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.globals-var-type-tag-is-invalid` |
| Scope | `semantic` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `DZCE3103`

Globals value type mismatch

> `globals.xml` var@value does not match declared var@type. Fix CE value
> format so it matches the selected type tag.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.globals-value-type-mismatch` |
| Scope | `semantic` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `DZCE3104`

Globals value is out of range

> `globals.xml` value is outside recommended CE range for this variable.
> Review gameplay impact and adjust if not intentional.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.globals-value-is-out-of-range` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3202`

Duplicate economycore default name

> `cfgeconomycore.xml` `<defaults>` contains duplicate default@name keys. Keep
> one value per key to avoid ambiguous CE runtime defaults.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.duplicate-economycore-default-name` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3203`

Economycore bool default is invalid

> `cfgeconomycore.xml` bool-like default key uses non-bool value. Use
> true/false or 0/1.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.economycore-bool-default-is-invalid` |
| Scope | `semantic` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `DZCE3204`

Economycore default is out of range

> A numeric default in `cfgeconomycore.xml` is outside expected CE limits.
> This may lead to unstable CE behavior.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.economycore-default-is-out-of-range` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3301`

Economy section flags are invalid

> `economy.xml` section is missing required CE flags or uses invalid values.
> Each section should define init/load/respawn/save as 0 or 1.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.economy-section-flags-are-invalid` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3302`

Type flags block is incomplete

> `types.xml` `<flags>` block does not define all commonly paired attributes.
> This can cause implicit inheritance/merge side effects.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.type-flags-block-is-incomplete` |
| Scope | `semantic` |
| Severity | `notice` |
| Enabled | `true` (implicit) |

#### `DZCE3303`

Invalid type quantity range

> `types.xml` quantity range (quantmin/quantmax) is invalid for CE. Allowed
> values are -1 or 0..100, and quantmin must not exceed quantmax.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.invalid-type-quantity-range` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3304`

Invalid spawnable damage range

> `cfgspawnabletypes.xml` damage range is invalid for CE spawn rules. Use
> values in 0..1 and ensure min is not greater than max.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.invalid-spawnable-damage-range` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3305`

Invalid spawnable chance range

> `cfgspawnabletypes.xml` chance values are inconsistent for CE spawn rules.
> Use one mode consistently: normalized 0..1 or percent 0..100.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.invalid-spawnable-chance-range` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3601`

Duplicate event name

> `events.xml` contains duplicate `<event>` names. Keep CE event names unique
> to avoid ambiguous merge and spawn behavior.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.duplicate-event-name` |
| Scope | `semantic` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `DZCE3604`

Event flags are not 0/1

> `events.xml` active/flag attributes use non-canonical CE values. Use 0 or 1.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.event-flags-are-not-0-1` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3605`

Unsupported event position

> `events.xml` position uses an unsupported CE value. Supported values: fixed,
> player, uniform.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.unsupported-event-position` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3606`

Unsupported event limit

> `events.xml` limit uses an unsupported CE value. Supported values: child,
> custom, mixed, parent.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.unsupported-event-limit` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE3702`

Duplicate random preset name

> `cfgrandompresets.xml` contains duplicate preset names in one section. Keep
> preset names unique inside cargo and attachments blocks.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.duplicate-random-preset-name` |
| Scope | `semantic` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `DZCE3703`

Random preset has no items

> `cfgrandompresets.xml` preset is declared without `<item>` entries and has
> no effect during generation.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.semantic.random-preset-has-no-items` |
| Scope | `semantic` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

### validate

> Validation diagnostics for CE XML shape and values.

Codes:
[DZCE2001](#dzce2001),
[DZCE2002](#dzce2002),
[DZCE2003](#dzce2003),
[DZCE2004](#dzce2004),
[DZCE2005](#dzce2005),

#### `DZCE2001`

Missing required attribute

> Required XML attribute is missing. Add the attribute so CE can interpret
> this element deterministically.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.validate.missing-required-attribute` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `DZCE2002`

Empty required attribute

> Required XML attribute is present but empty. Provide a non-empty value.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.validate.empty-required-attribute` |
| Scope | `validate` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `DZCE2003`

Invalid bool value

> Boolean-like field uses a non-canonical value. Use 0 or 1 for CE XML boolean
> fields.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.validate.invalid-bool-value` |
| Scope | `validate` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE2004`

Invalid integer range

> Integer value is outside expected range for this field.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.validate.invalid-integer-range` |
| Scope | `validate` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `DZCE2005`

Unknown enum value

> Enum field contains unsupported token for this CE context.

| Field | Value |
| --- | --- |
| Rule ID | `dzce.validate.unknown-enum-value` |
| Scope | `validate` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

---

> Generated with
> [lintkit](https://github.com/woozymasta/lintkit)
> version `dev`
> commit `unknown`

<!-- Automatically generated file, do not modify! -->
