# DayZ CE AreaFlags Format Notes

This document is a format-level reference for DayZ Central Economy area data.
It describes CE project `zg-config`, `areaflags.map`, and their relationship.

## Purpose

The goal is to explain:

* what data is authored in CE project files
* what data is actually encoded into `areaflags.map`
* why reverse reconstruction from `map` is not always identical to source

This is not a gameplay design guide and not a package API manual.

## CE Project Artifacts

Typical CE project contains:

* `mapname.xml` (`<zg-config>`)
* `layers/<layer-name>.<ext>` source masks
* `map/<mapname>.map` compiled areaflags binary
* `territoryTypes/*_territories.xml` territory data

`territoryTypes` data is related CE content, but it is not encoded into
`areaflags.map`.

## zg-config Parts Relevant to AreaFlags

Relevant subtree in `<zg-config>`:

* `global/layer/@size`: raster size (width and height)
* `global/world/@size`: world size metadata
* `areas/usages/usage[]`: usage flag names (UI labels)
* `areas/values/value[]`: value flag names (UI labels)
* `layers/layer[]`: layer definitions and flag masks

Important `layers/layer` attributes:

* `name`: logical layer name
* `usage_flags`: decimal usage bitmask
* `value_flags`: decimal value bitmask
* `color`, `visible`: editor/UI metadata

`visible` controls editor visibility only. It does not affect map encoding.

## Usage and Value Flags

`usage_flags` and `value_flags` are decimal bitmasks.

Examples:

* `usage_flags="32"` means usage bit #5 (`1 << 5`)
* `value_flags="16"` means value bit #4 (`1 << 4`)

At one pixel, a layer contributes bits only when its source pixel is ON.

## Per-Pixel Write Model

Per pixel, encoded map stores:

* `U`: usage bitfield
* `V`: value bitfield

For a layer with bitmasks `UL` and `VL`:

* if source pixel is OFF: no contribution
* if source pixel is ON:
  `U = U OR UL`
  `V = V OR VL`

Encoding is OR-accumulation across all layers.

## Why Reverse Is Not One-to-One

`areaflags.map` stores final bits, not provenance of source layers.

Because of this, `map -> layers` is generally not one-to-one:

* multiple source masks may write into the same bit
* different layers may use identical `(usage_flags, value_flags)` pairs
* `usage_flags=0` and `value_flags=0` layers write nothing to map

So reverse reconstruction can restore logical masks by current bit rules, but
not always original independent author masks.

## Special Case: Zero/Zero Layers

Layers with:

* `usage_flags="0"`
* `value_flags="0"`

are not encoded into `areaflags.map`.

They can exist as editor-only/source-only layers (for example point overlays).
Such layers require original source masks for full restoration.

## Binary Layout: areaflags.map

On-disk layout:

* `u32 layerWidth`
* `u32 layerHeight`
* `u32 worldWidth`
* `u32 worldHeight`
* `u32 usageBits` marker
* `usageSection` payload
* `u32 valueBits` marker
* `valueSection` payload

Section byte sizes:

* `usageSectionBytes = layerWidth * layerHeight * usageBits / 8`
* `valueSectionBytes = layerWidth * layerHeight * valueBits / 8`

For ChernarusPlus:

* `layer=4096`, `usageBits=32`, `valueBits=8`
* total size = `83886104` bytes

## Packed Section Rules

Each section uses 32-bit words and fixed-width slots.

Rules:

* `bitsPerPixel` is section bit depth (`usageBits` or `valueBits`)
* `bitsPerPixel` must be power of two and divide `32`
* `slotsPerWord = 32 / bitsPerPixel`
* one word stores `slotsPerWord` neighboring X pixels
* rows are packed in row-major order

## Coordinate Mapping

Source layer masks use top-left origin.
Packed map uses opposite Y orientation.

Mapping:

* `mapX = layerX`
* `mapY = layerHeight - 1 - layerY`

So X is direct, Y is flipped.

## territory-type-list Save Behavior

In CE project save/export flows, embedded `territory-type-list` in `mapname.xml`
may be rewritten from `territoryTypes/*.xml`.

This synchronization is independent from `areaflags.map` encoding.

## Practical Implications

For deterministic automation:

* keep layer bit assignments stable
* avoid accidental flag collisions when independent restoration is required
* treat `zero/zero` layers as source-dependent artifacts
* treat `territoryTypes` as separate data stream from areaflags binary
