// Package pricing exposes the embedded v1/pricing.json table as typed Go
// structs. Consumers call Load to get a freshly-parsed, mutation-safe copy
// of the pricing document, then use Table.Get to look up individual models.
package pricing

import (
	_ "embed"
	"encoding/json"
)

//go:embed v1/pricing.json
var rawPricingJSON []byte

// Record is one model's pricing entry.
type Record struct {
	Provider                   string  `json:"provider"`
	PromptUSDPer1KTokens       float64 `json:"prompt_usd_per_1k_tokens"`
	CompletionUSDPer1KTokens   float64 `json:"completion_usd_per_1k_tokens"`
	CachedPromptUSDPer1KTokens float64 `json:"cached_prompt_usd_per_1k_tokens,omitempty"`
	ContextWindowTokens        int     `json:"context_window_tokens,omitempty"`
	Deprecated                 bool    `json:"deprecated,omitempty"`
	DeprecationNotes           string  `json:"deprecation_notes,omitempty"`
}

// Table is the full pricing document.
type Table struct {
	SchemaVersion string            `json:"schema_version"`
	GeneratedAt   string            `json:"generated_at"`
	SourceNotes   string            `json:"source_notes,omitempty"`
	Models        map[string]Record `json:"models"`
}

// Load parses the embedded pricing.json. Callers may call Load repeatedly;
// each call returns a freshly parsed copy (callers are free to mutate).
func Load() (Table, error) {
	var t Table
	if err := json.Unmarshal(rawPricingJSON, &t); err != nil {
		return Table{}, err
	}
	return t, nil
}

// Get looks up model in the table. Returns (zero, false) on miss.
func (t Table) Get(model string) (Record, bool) {
	r, ok := t.Models[model]
	return r, ok
}
