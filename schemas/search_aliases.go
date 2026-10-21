// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package schemas

// searchAliases maps user-facing terms to one-step search expansions.
var searchAliases = []searchAlias{
	{query: "infected", target: "zombies", weight: 0.95},
	{query: "zombie", target: "infected", weight: 0.95},
	{query: "zombies", target: "infected", weight: 0.95},
	{query: "undead", target: "infected", weight: 0.65},
	{query: "survivor", target: "player", weight: 0.90},
	{query: "character", target: "player", weight: 0.50},
	{query: "character", target: "infected", weight: 0.45},
	{query: "character", target: "animal", weight: 0.45},
	{query: "man", target: "player", weight: 0.40},
	{query: "wildlife", target: "animal", weight: 0.90},
	{query: "animal", target: "wildlife", weight: 0.55},
	{query: "motorbike", target: "vehicle", weight: 0.65},
	{query: "motorcycle", target: "vehicle", weight: 0.65},
	{query: "bike", target: "vehicle", weight: 0.40},
	{query: "boat", target: "vehicle", weight: 0.55},
	{query: "watercraft", target: "vehicle", weight: 0.50},
	{query: "weather", target: "overcast", weight: 0.55},
	{query: "overcast", target: "weather", weight: 0.45},
	{query: "cloud cover", target: "overcast", weight: 0.80},
	{query: "cloudiness", target: "overcast", weight: 0.80},
	{query: "coordinates", target: "position", weight: 0.85},
	{query: "coords", target: "position", weight: 0.85},
	{query: "position", target: "coordinates", weight: 0.55},
	{query: "location", target: "position", weight: 0.65},
	{query: "location", target: "territory", weight: 0.40},
	{query: "location", target: "zone", weight: 0.40},
	{query: "place", target: "position", weight: 0.60},
	{query: "place", target: "territory", weight: 0.35},
	{query: "place", target: "zone", weight: 0.35},
	{query: "area", target: "zone", weight: 0.55},
	{query: "zone", target: "area", weight: 0.45},
	{query: "region", target: "territory", weight: 0.55},
	{query: "territory", target: "region", weight: 0.45},
	{query: "bunker", target: "tunnel", weight: 0.30},
	{query: "underground", target: "bunker", weight: 0.35},
	{query: "underground", target: "tunnel", weight: 0.40},
	{query: "dungeon", target: "bunker", weight: 0.35},
	{query: "dungeon", target: "tunnel", weight: 0.40},
	{query: "cave", target: "bunker", weight: 0.35},
	{query: "cave", target: "tunnel", weight: 0.40},
	{query: "tunnel", target: "bunker", weight: 0.35},
	{query: "cooldown", target: "restock", weight: 0.65},
	{query: "expiration", target: "lifetime", weight: 0.65},
	{query: "decay", target: "lifetime", weight: 0.50},
	{query: "target count", target: "nominal", weight: 0.70},
	{query: "population target", target: "nominal", weight: 0.65},
	{query: "buried stash", target: "hoarder", weight: 0.80},
	{query: "stash", target: "hoarder", weight: 0.60},
	{query: "storage", target: "hoarder", weight: 0.50},
	{query: "event loot", target: "deloot", weight: 0.75},
	{query: "dynamic loot", target: "deloot", weight: 0.70},
	{query: "chance", target: "weight", weight: 0.45},
	{query: "probability", target: "chance", weight: 0.35},
	{query: "lightning", target: "storm", weight: 0.65},
	{query: "line of sight", target: "raycast", weight: 0.75},
	{query: "route", target: "breadcrumbs", weight: 0.35},
	{query: "smooth transition", target: "interpolation", weight: 0.45},
	{query: "darkness", target: "eye accommodation", weight: 0.45},
	{query: "dim", target: "eye accommodation", weight: 0.35},
	{query: "screen effect", target: "post-processing", weight: 0.55},
	{query: "hud", target: "screen effect", weight: 0.30},
	{query: "ui", target: "screen effect", weight: 0.30},
	{query: "overlay", target: "screen effect", weight: 0.30},
}

// searchAlias describes a directed, weighted query expansion.
type searchAlias struct {
	query  string  // query is the user term that activates the expansion.
	target string  // target is the literal term searched in field metadata.
	weight float64 // weight scales the alias match score between zero and one.
}

// searchExpansion is one target term and its ranking multiplier.
type searchExpansion struct {
	term   string  // term is the literal search target.
	weight float64 // weight scales the target's match score.
}

// searchExpansionsFor returns direct aliases without following alias chains.
func searchExpansionsFor(query string) []searchExpansion {
	var expansions []searchExpansion
	for _, alias := range searchAliases {
		if alias.query != query {
			continue
		}

		expansions = append(expansions, searchExpansion{
			term:   alias.target,
			weight: alias.weight,
		})
	}

	return expansions
}
