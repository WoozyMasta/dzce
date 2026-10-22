// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

/*
Package dzce provides typed models and codecs for
DayZ Central Economy (CE) mission configuration files,
including XML, JSON, CEProject layer masks, and areaflags.map data.

[DetectKind] identifies formats from file paths,
while [Decode] and [Encode] convert payloads through the built-in registry.
[LoadFile] and [SaveFile] provide path-based access;
[LoadFile] also detects CEProject XML from its root element when the filename is not recognized.
Use [LoadFileAs] when a format uses arbitrary filenames.
[SaveFile] can infer the kind for CEProject configurations and area-flags maps from their value types.

[AnalyzeLintContent] and [ValidateContent] decode a CE file
and run its existing structural and semantic checks.
Their [AnalysisResult] reports status, the count of error-severity diagnostics, and all findings.
[ValidateFile] provides the same analysis for a file path
and returns an error only when the file cannot be read.
JSON decoding does not validate against generated JSON Schemas or reject unknown fields.

[CEPipeline.Forward] builds areaflags.map and territory XML files
from a CEProject configuration and its source layers.
[CEPipeline.Restore] recreates the configuration and layer masks from compiled map data.
The lower-level [BuildAreaFlagsFromCEProjectLayerDir]
and [EncodeAreaFlagsMap] functions are available for custom workflows.
The package includes a TGA mask codec and supports custom codecs
through [RegisterMaskImageCodec] and [UnregisterMaskImageCodec].

[LoadMergedEconomyCore] resolves and merges an economycore include tree.
Strictmode accepts the wiki-compatible include types and recursive economycore includes;
[LoadMergedEconomyCoreWithOptions] enables additional kind resolution through [MergeOptions].
Use [MergedConfig.Get] to retrieve a merged document by kind.
Merge results contain included files only and do not preload vanilla or base CE datasets.

Runnable usage examples are in the package example tests.
*/
package dzce
