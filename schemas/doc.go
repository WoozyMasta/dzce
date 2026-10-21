// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

/*
Package schemas provides access to the embedded JSON Schemas
for DayZ Central Economy configuration models.

[Catalog] lists schemas and their descriptions.
[Lookup] finds an exact field key,
while [Search] ranks case-insensitive matches and weighted query aliases.
Literal matches rank above alias matches, followed by subsequence matches.
Both APIs include only root fields by default;
[SearchOptions] enables nested fields when needed.
[Field] exposes modeled JSON Schema constraints and examples as typed values.
*/
package schemas
