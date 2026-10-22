<!-- markdownlint-disable MD024 -->
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog][],
and this project adheres to [Semantic Versioning][].

<!--
## Unreleased

### Added
### Changed
### Removed
-->

## Unreleased

### Added

* Generated JSON Schemas and Markdown references for supported CE models,
  with a typed catalog, exact field lookup,
  and ranked search with weighted aliases for common and technical terms.
* Models current official mission JSON and XML schemas and validates
  `areaflags.map` files in integration tests.
* Validates declared numeric, enum, and length constraints in modeled
  CE configuration files.

### Changed

* Saves individual CE files atomically and preserves existing file permissions.
* Validates spawnable candidate chances as probabilities
  or percentage weights within each cargo or attachment group.

### Fixed

* Rejects overflowing `areaflags.map` section sizes before allocating layers.

## [0.2.0][] - 2026-04-07

### Added

* Linting support via `lintkit` for CE configuration.

### Changed

* File kind detection now uses `bimime` as primary source
  (extension/path and content-aware detection).

[0.2.0]: https://github.com/WoozyMasta/dzce/compare/v0.1.0...v0.2.0

## [0.1.0][] - 2026-03-17

### Added

* First public release

[0.1.0]: https://github.com/WoozyMasta/dzce/tree/v0.1.0

<!--links-->
[Keep a Changelog]: https://keepachangelog.com/en/1.1.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
