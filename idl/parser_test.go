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

func TestParser_ValidateMessageFields(t *testing.T) {
	parser := NewParser()

	// Test missing field type
	yamlContent := `
name: "TestService"
version: "1.0.0"
messages:
  TestMessage:
    fields:
      id:
        required: true
      name:
        type: "string"
services: {}
`

	_, err := parser.Parse([]byte(yamlContent))
	if err == nil {
		t.Error("Expected validation error for missing field type")
	}

	// Verify the error message contains the field information
	expectedErrMsg := "field type is required for TestMessage.id"
	if err != nil && err.Error() != "validation failed: "+expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestParser_ValidateServiceMethods(t *testing.T) {
	parser := NewParser()

	// Test method with both input and output empty
	yamlContent := `
name: "TestService"
version: "1.0.0"
messages: {}
services:
  TestService:
    methods:
      EmptyMethod:
        input: ""
        output: ""
`

	_, err := parser.Parse([]byte(yamlContent))
	if err == nil {
		t.Error("Expected validation error for method with empty input and output")
	}

	// Verify the error message contains the method information
	expectedErrMsg := "both input and output types cannot be empty for TestService.EmptyMethod"
	if err != nil && err.Error() != "validation failed: "+expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestParser_ValidateServiceMethodsWithOneType(t *testing.T) {
	parser := NewParser()

	// Test method with only input (should be valid)
	yamlContent := `
name: "TestService"
version: "1.0.0"
messages:
  TestMessage:
    fields:
      id:
        type: "string"
services:
  TestService:
    methods:
      InputOnlyMethod:
        input: "TestMessage"
        output: ""
`

	_, err := parser.Parse([]byte(yamlContent))
	if err != nil {
		t.Errorf("Expected no validation error for method with only input, got: %v", err)
	}

	// Test method with only output (should be valid)
	yamlContent2 := `
name: "TestService"
version: "1.0.0"
messages:
  TestMessage:
    fields:
      id:
        type: "string"
services:
  TestService:
    methods:
      OutputOnlyMethod:
        input: ""
        output: "TestMessage"
`

	_, err2 := parser.Parse([]byte(yamlContent2))
	if err2 != nil {
		t.Errorf("Expected no validation error for method with only output, got: %v", err2)
	}
}

func TestParser_ValidateComplexSchema(t *testing.T) {
	parser := NewParser()

	// Test valid complex schema with multiple messages and services
	yamlContent := `
name: "ComplexService"
version: "2.0.0"
package: "com.example"
messages:
  User:
    description: "User entity"
    fields:
      id:
        type: "string"
        required: true
        description: "User ID"
      name:
        type: "string"
        required: true
      email:
        type: "string"
        required: false
  UserList:
    fields:
      users:
        type: "User"
        repeated: true
services:
  UserService:
    description: "User management service"
    methods:
      GetUser:
        description: "Get user by ID"
        input: "string"
        output: "User"
      ListUsers:
        input: ""
        output: "UserList"
      CreateUser:
        input: "User"
        output: "User"
  NotificationService:
    methods:
      SendNotification:
        input: "string"
        output: ""
`

	schema, err := parser.Parse([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to parse valid complex YAML: %v", err)
	}

	// Validate the parsed schema structure
	if schema.Name != "ComplexService" {
		t.Errorf("Expected name 'ComplexService', got '%s'", schema.Name)
	}

	if schema.Version != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got '%s'", schema.Version)
	}

	if schema.Package != "com.example" {
		t.Errorf("Expected package 'com.example', got '%s'", schema.Package)
	}

	if len(schema.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(schema.Messages))
	}

	if len(schema.Services) != 2 {
		t.Errorf("Expected 2 services, got %d", len(schema.Services))
	}

	// Validate specific message fields
	userMessage := findMessage(schema.Messages, "User")
	if userMessage == nil {
		t.Error("Expected to find User message")
	} else {
		if len(userMessage.Fields) != 3 {
			t.Errorf("Expected User message to have 3 fields, got %d", len(userMessage.Fields))
		}
	}

	// Validate specific service methods
	userService := findService(schema.Services, "UserService")
	if userService == nil {
		t.Error("Expected to find UserService")
	} else {
		if len(userService.Methods) != 3 {
			t.Errorf("Expected UserService to have 3 methods, got %d", len(userService.Methods))
		}
	}
}

// Helper functions for test validation
func findMessage(messages []Message, name string) *Message {
	for _, msg := range messages {
		if msg.Name == name {
			return &msg
		}
	}
	return nil
}

func findService(services []Service, name string) *Service {
	for _, svc := range services {
		if svc.Name == name {
			return &svc
		}
	}
	return nil
}
