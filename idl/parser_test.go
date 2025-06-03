package idl

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := NewParser()

	yamlContent := `
name: "TestService"
version: "1.0.0"
messages:
  TestMessage:
    fields:
      id:
        type: "string"
        required: true
services:
  TestService:
    methods:
      TestMethod:
        input: "TestMessage"
        output: "TestMessage"
`

	schema, err := parser.Parse([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to parse valid YAML: %v", err)
	}

	if schema.Name != "TestService" {
		t.Errorf("Expected name 'TestService', got '%s'", schema.Name)
	}

	if len(schema.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(schema.Messages))
	}

	if len(schema.Services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(schema.Services))
	}
}

func TestParser_ValidateRequired(t *testing.T) {
	parser := NewParser()

	// Test missing name
	yamlContent := `
version: "1.0.0"
messages: {}
services: {}
`

	_, err := parser.Parse([]byte(yamlContent))
	if err == nil {
		t.Error("Expected validation error for missing name")
	}
}
