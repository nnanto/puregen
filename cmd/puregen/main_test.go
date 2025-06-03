package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainWithVersionCommand(t *testing.T) {
	// Phase 2: When BE_CRASHER=1, run the actual main() function
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen", "version"}
		main() // This will call os.Exit() or return normally
		return
	}

	// Phase 1: Spawn a subprocess with BE_CRASHER=1 to test main()
	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithVersionCommand")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	out, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Expected successful execution, got error: %v", err)
	}

	output := string(out)
	if !strings.Contains(output, "puregen version") {
		t.Errorf("Expected version output, got: %s", output)
	}
}

func TestMainWithGenerateNoFlags(t *testing.T) {
	// Test generate command with no flags
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen", "generate"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithGenerateNoFlags")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("Expected non-zero exit code")
	}

	output := string(out)
	if !strings.Contains(output, "required flag") {
		t.Errorf("Expected required flag error, got: %s", output)
	}
}

func TestMainWithGenerateMissingInput(t *testing.T) {
	// Test generate command with missing input flag
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen", "generate", "--templates", "test.tmpl"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithGenerateMissingInput")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("Expected non-zero exit code")
	}

	output := string(out)
	if !strings.Contains(output, "required flag") && !strings.Contains(output, "input") {
		t.Errorf("Expected input flag error, got: %s", output)
	}
}

func TestMainWithGenerateMissingTemplates(t *testing.T) {
	// Test generate command with missing templates flag
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen", "generate", "--input", "test.yaml"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithGenerateMissingTemplates")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("Expected non-zero exit code")
	}

	output := string(out)
	if !strings.Contains(output, "required flag") && !strings.Contains(output, "templates") {
		t.Errorf("Expected templates flag error, got: %s", output)
	}
}

func TestMainWithNonExistentYamlFile(t *testing.T) {
	// Test with non-existent YAML file
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen", "generate", "--input", "nonexistent.yaml", "--templates", "template.tmpl"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithNonExistentYamlFile")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("Expected non-zero exit code")
	}

	output := string(out)
	if !strings.Contains(output, "Error parsing IDL file") {
		t.Errorf("Expected IDL parsing error, got: %s", output)
	}
}

func TestMainWithNonExistentTemplateFile(t *testing.T) {
	// Create a temporary YAML file
	yamlFile, err := os.CreateTemp("", "test*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(yamlFile.Name())

	// Write minimal valid YAML content
	_, err = yamlFile.WriteString("services:\n  test:\n    methods:\n      - name: test\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	yamlFile.Close()

	// Test with non-existent template file
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen", "generate", "--input", yamlFile.Name(), "--templates", "nonexistent.tmpl"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithNonExistentTemplateFile")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("Expected non-zero exit code")
	}

	output := string(out)
	if !strings.Contains(output, "Error parsing") {
		t.Errorf("Expected template file error, got: %s", output)
	}
}

func TestMainWithMultipleTemplates(t *testing.T) {
	// Create temporary YAML file
	yamlFile, err := os.CreateTemp("", "test*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(yamlFile.Name())

	_, err = yamlFile.WriteString("services:\n  test:\n    methods:\n      - name: test\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	yamlFile.Close()

	// Test with multiple template files (some non-existent)
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen", "generate", "--input", yamlFile.Name(), "--templates", "nonexistent1.tmpl,nonexistent2.tmpl"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithMultipleTemplates")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("Expected non-zero exit code")
	}

	output := string(out)
	if !strings.Contains(output, "Error parsing") {
		t.Errorf("Expected template file error, got: %s", output)
	}
}
