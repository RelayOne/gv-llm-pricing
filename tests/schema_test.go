package tests_test

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

const schemaPath = "../v1/schema.json"
const schemaURL = "https://raw.githubusercontent.com/RelayOne/gv-llm-pricing/main/v1/schema.json"

// loadSchema compiles ../v1/schema.json relative to the test file's working directory.
func loadSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	raw, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(schemaURL, bytes.NewReader(raw)); err != nil {
		t.Fatalf("AddResource: %v", err)
	}
	schema, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	return schema
}

// mustMarshalAndValidate marshals obj to JSON, unmarshals it into `any`, and
// validates. Returns any validation error.
func mustMarshalAndValidate(t *testing.T, schema *jsonschema.Schema, obj any) error {
	t.Helper()
	raw, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return schema.Validate(v)
}

func validModel() map[string]any {
	return map[string]any{
		"provider":                     "openai",
		"prompt_usd_per_1k_tokens":     0.01,
		"completion_usd_per_1k_tokens": 0.03,
	}
}

func validDoc() map[string]any {
	return map[string]any{
		"schema_version": "v1",
		"generated_at":   "2026-04-18T12:00:00Z",
		"models": map[string]any{
			"gpt-4o-mini": validModel(),
		},
	}
}

func TestSchema_ValidFixture(t *testing.T) {
	schema := loadSchema(t)
	if err := mustMarshalAndValidate(t, schema, validDoc()); err != nil {
		t.Errorf("expected valid document to pass, got error: %v", err)
	}
}

func TestSchema_MissingSchemaVersion(t *testing.T) {
	schema := loadSchema(t)
	doc := validDoc()
	delete(doc, "schema_version")
	err := mustMarshalAndValidate(t, schema, doc)
	if err == nil {
		t.Fatalf("expected validation error for missing schema_version, got nil")
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "schema_version") {
		t.Errorf("error should mention schema_version, got: %v", err)
	}
	if !strings.Contains(msg, "required") {
		t.Errorf("error should mention required, got: %v", err)
	}
}

func TestSchema_NegativePrice(t *testing.T) {
	schema := loadSchema(t)
	doc := validDoc()
	bad := validModel()
	bad["prompt_usd_per_1k_tokens"] = -0.01
	doc["models"] = map[string]any{"gpt-4o-mini": bad}
	err := mustMarshalAndValidate(t, schema, doc)
	if err == nil {
		t.Fatalf("expected validation error for negative price, got nil")
	}
	msg := err.Error()
	lower := strings.ToLower(msg)
	if !strings.Contains(lower, "prompt_usd_per_1k_tokens") {
		t.Errorf("error should mention prompt_usd_per_1k_tokens, got: %v", err)
	}
	if !strings.Contains(lower, "minimum") && !strings.Contains(msg, ">= 0") {
		t.Errorf("error should mention minimum or >= 0, got: %v", err)
	}
}

func TestSchema_ExtraTopLevelKey(t *testing.T) {
	schema := loadSchema(t)
	doc := validDoc()
	doc["extra_key"] = "x"
	err := mustMarshalAndValidate(t, schema, doc)
	if err == nil {
		t.Fatalf("expected validation error for extra top-level key, got nil")
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "additionalproperties") && !strings.Contains(msg, "extra_key") {
		t.Errorf("error should mention additionalProperties or extra_key, got: %v", err)
	}
}

func TestSchema_UnknownProvider(t *testing.T) {
	schema := loadSchema(t)
	doc := validDoc()
	bad := validModel()
	bad["provider"] = "ibm"
	doc["models"] = map[string]any{"some-model": bad}
	err := mustMarshalAndValidate(t, schema, doc)
	if err == nil {
		t.Fatalf("expected validation error for unknown provider, got nil")
	}
	msg := err.Error()
	lower := strings.ToLower(msg)
	if !strings.Contains(lower, "provider") {
		t.Errorf("error should mention provider, got: %v", err)
	}
	validProviders := []string{"openai", "anthropic", "google", "meta", "mistral", "cohere", "other"}
	hasEnum := strings.Contains(lower, "enum")
	hasProvider := false
	for _, p := range validProviders {
		if strings.Contains(msg, p) {
			hasProvider = true
			break
		}
	}
	if !hasEnum && !hasProvider {
		t.Errorf("error should mention enum or a valid provider, got: %v", err)
	}
}
