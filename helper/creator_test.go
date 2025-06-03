package helper

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCreator(t *testing.T) {
	creator := NewCreator("test.yaml")

	if creator.filePath != "test.yaml" {
		t.Errorf("Expected filePath to be 'test.yaml', got %s", creator.filePath)
	}

	if creator.schema == nil {
		t.Error("Expected schema to be initialized")
	}

	if creator.schema.Messages == nil {
		t.Error("Expected Messages map to be initialized")
	}

	if creator.schema.Services == nil {
		t.Error("Expected Services map to be initialized")
	}

	if creator.fileExists {
		t.Error("Expected fileExists to be false for new creator")
	}
}

func TestLoadOrCreate_NewFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "new_schema.yaml")

	creator := NewCreator(filePath)
	err := creator.LoadOrCreate()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if creator.FileExists() {
		t.Error("Expected FileExists to return false for new file")
	}

	if creator.schema.Name != "MySchema" {
		t.Errorf("Expected default name 'MySchema', got %s", creator.schema.Name)
	}

	if creator.schema.Version != "1.0.0" {
		t.Errorf("Expected default version '1.0.0', got %s", creator.schema.Version)
	}

	if creator.schema.Package != "main" {
		t.Errorf("Expected default package 'main', got %s", creator.schema.Package)
	}
}

func TestLoadOrCreate_ExistingFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "existing_schema.yaml")

	// Create a test YAML file
	content := `name: TestSchema
version: 2.0.0
package: test
messages:
  User:
    name: User
    description: A user message
    fields:
      id:
        type: string
        required: true
services:
  UserService:
    description: User service
    methods: {}
`

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	creator := NewCreator(filePath)
	err = creator.LoadOrCreate()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !creator.FileExists() {
		t.Error("Expected FileExists to return true for existing file")
	}

	if creator.schema.Name != "TestSchema" {
		t.Errorf("Expected name 'TestSchema', got %s", creator.schema.Name)
	}

	if creator.schema.Version != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got %s", creator.schema.Version)
	}

	if creator.schema.Package != "test" {
		t.Errorf("Expected package 'test', got %s", creator.schema.Package)
	}

	if len(creator.schema.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(creator.schema.Messages))
	}

	if len(creator.schema.Services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(creator.schema.Services))
	}
}

func TestCreateMessage(t *testing.T) {
	creator := NewCreator("test.yaml")

	message := creator.CreateMessage("TestMessage", "A test message")

	if message.Name != "TestMessage" {
		t.Errorf("Expected name 'TestMessage', got %s", message.Name)
	}

	if message.Description != "A test message" {
		t.Errorf("Expected description 'A test message', got %s", message.Description)
	}

	if message.Fields == nil {
		t.Error("Expected Fields to be initialized")
	}

	if message.Metadata == nil {
		t.Error("Expected Metadata to be initialized")
	}
}

func TestCreateField(t *testing.T) {
	creator := NewCreator("test.yaml")

	field := creator.CreateField("string", "A test field", true, false)

	if field.Type != "string" {
		t.Errorf("Expected type 'string', got %s", field.Type)
	}

	if field.Description != "A test field" {
		t.Errorf("Expected description 'A test field', got %s", field.Description)
	}

	if !field.Required {
		t.Error("Expected field to be required")
	}

	if field.Repeated {
		t.Error("Expected field to not be repeated")
	}

	if field.Metadata == nil {
		t.Error("Expected Metadata to be initialized")
	}
}

func TestCreateService(t *testing.T) {
	creator := NewCreator("test.yaml")

	service := creator.CreateService("A test service")

	if service.Description != "A test service" {
		t.Errorf("Expected description 'A test service', got %s", service.Description)
	}

	if service.Methods == nil {
		t.Error("Expected Methods to be initialized")
	}

	if service.Metadata == nil {
		t.Error("Expected Metadata to be initialized")
	}
}

func TestCreateMethod(t *testing.T) {
	creator := NewCreator("test.yaml")

	method := creator.CreateMethod("TestMethod", "A test method", "InputType", "OutputType", true)

	if method.Name != "TestMethod" {
		t.Errorf("Expected name 'TestMethod', got %s", method.Name)
	}

	if method.Description != "A test method" {
		t.Errorf("Expected description 'A test method', got %s", method.Description)
	}

	if method.Input != "InputType" {
		t.Errorf("Expected input 'InputType', got %s", method.Input)
	}

	if method.Output != "OutputType" {
		t.Errorf("Expected output 'OutputType', got %s", method.Output)
	}

	if !method.Streaming {
		t.Error("Expected method to be streaming")
	}

	if method.Metadata == nil {
		t.Error("Expected Metadata to be initialized")
	}
}

func TestAddMessage(t *testing.T) {
	creator := NewCreator("test.yaml")
	creator.LoadOrCreate()

	message := creator.CreateMessage("TestMessage", "A test message")
	err := creator.AddMessage("TestMessage", message)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(creator.schema.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(creator.schema.Messages))
	}

	storedMessage, exists := creator.schema.Messages["TestMessage"]
	if !exists {
		t.Error("Expected message to be stored")
	}

	if storedMessage.Name != "TestMessage" {
		t.Errorf("Expected stored message name 'TestMessage', got %s", storedMessage.Name)
	}

	// Test empty name error
	err = creator.AddMessage("", message)
	if err == nil {
		t.Error("Expected error for empty message name")
	}
}

func TestAddService(t *testing.T) {
	creator := NewCreator("test.yaml")
	creator.LoadOrCreate()

	service := creator.CreateService("A test service")
	err := creator.AddService("TestService", service)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(creator.schema.Services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(creator.schema.Services))
	}

	storedService, exists := creator.schema.Services["TestService"]
	if !exists {
		t.Error("Expected service to be stored")
	}

	if storedService.Description != "A test service" {
		t.Errorf("Expected stored service description 'A test service', got %s", storedService.Description)
	}

	// Test empty name error
	err = creator.AddService("", service)
	if err == nil {
		t.Error("Expected error for empty service name")
	}
}

func TestIsValidType(t *testing.T) {
	creator := NewCreator("test.yaml")
	creator.LoadOrCreate()

	// Test primitive types
	primitives := []string{
		"string", "int32", "int64", "uint32", "uint64",
		"float32", "float64", "bool", "bytes",
	}

	for _, primitive := range primitives {
		if !creator.IsValidType(primitive) {
			t.Errorf("Expected %s to be a valid type", primitive)
		}
	}

	// Test invalid primitive
	if creator.IsValidType("invalid") {
		t.Error("Expected 'invalid' to be an invalid type")
	}

	// Add a message and test it's valid
	message := creator.CreateMessage("CustomMessage", "A custom message")
	creator.AddMessage("CustomMessage", message)

	if !creator.IsValidType("CustomMessage") {
		t.Error("Expected 'CustomMessage' to be a valid type")
	}

	// Test non-existent message
	if creator.IsValidType("NonExistentMessage") {
		t.Error("Expected 'NonExistentMessage' to be an invalid type")
	}
}

func TestSave(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "save_test.yaml")

	creator := NewCreator(filePath)
	creator.LoadOrCreate()

	// Add some content
	message := creator.CreateMessage("TestMessage", "A test message")
	field := creator.CreateField("string", "Test field", true, false)
	message.Fields["testField"] = field
	creator.AddMessage("TestMessage", message)

	err := creator.Save()
	if err != nil {
		t.Errorf("Expected no error saving, got %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Expected file to be created")
	}

	// Load the file and verify content
	newCreator := NewCreator(filePath)
	err = newCreator.LoadOrCreate()
	if err != nil {
		t.Errorf("Expected no error loading saved file, got %v", err)
	}

	if len(newCreator.schema.Messages) != 1 {
		t.Errorf("Expected 1 message in loaded file, got %d", len(newCreator.schema.Messages))
	}

	loadedMessage, exists := newCreator.schema.Messages["TestMessage"]
	if !exists {
		t.Error("Expected TestMessage to exist in loaded file")
	}

	if loadedMessage.Description != "A test message" {
		t.Errorf("Expected loaded message description 'A test message', got %s", loadedMessage.Description)
	}
}

func TestGetSchema(t *testing.T) {
	creator := NewCreator("test.yaml")

	schema := creator.GetSchema()

	if schema != creator.schema {
		t.Error("Expected GetSchema to return the same schema instance")
	}
}
