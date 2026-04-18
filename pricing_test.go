package pricing_test

import (
	"reflect"
	"testing"

	"github.com/RelayOne/gv-llm-pricing"
)

func TestLoad_Succeeds(t *testing.T) {
	table, err := pricing.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if table.SchemaVersion != "v1" {
		t.Errorf("SchemaVersion = %q, want %q", table.SchemaVersion, "v1")
	}
	if got, want := len(table.Models), 19; got != want {
		t.Errorf("len(Models) = %d, want %d", got, want)
	}
}

func TestGet_UnknownModel(t *testing.T) {
	table, err := pricing.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	rec, ok := table.Get("does-not-exist")
	if ok {
		t.Errorf("Get(unknown) ok = true, want false")
	}
	if !reflect.DeepEqual(rec, pricing.Record{}) {
		t.Errorf("Get(unknown) record = %+v, want zero value", rec)
	}
}

func TestGet_GPT4o_NonZero(t *testing.T) {
	table, err := pricing.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	rec, ok := table.Get("gpt-4o")
	if !ok {
		t.Fatalf(`Get("gpt-4o") ok = false, want true`)
	}
	if rec.PromptUSDPer1KTokens <= 0 {
		t.Errorf("PromptUSDPer1KTokens = %v, want > 0", rec.PromptUSDPer1KTokens)
	}
}

func TestLoad_Idempotent(t *testing.T) {
	first, err := pricing.Load()
	if err != nil {
		t.Fatalf("Load() #1 error: %v", err)
	}
	second, err := pricing.Load()
	if err != nil {
		t.Fatalf("Load() #2 error: %v", err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("Load() returned differing tables")
	}
	first.Models["gpt-4o"] = pricing.Record{}
	third, err := pricing.Load()
	if err != nil {
		t.Fatalf("Load() #3 error: %v", err)
	}
	rec, ok := third.Get("gpt-4o")
	if !ok || rec.PromptUSDPer1KTokens <= 0 {
		t.Errorf("after mutating first, third.Get(gpt-4o) = (%+v, %v); want non-zero, ok", rec, ok)
	}
}
