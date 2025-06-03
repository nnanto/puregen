package main

import (
	"flag"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainWithVersionFlag(t *testing.T) {
	// Phase 2: When BE_CRASHER=1, run the actual main() function
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen", "-version"}
		main() // This will call os.Exit() or return normally
		return
	}

	// Phase 1: Spawn a subprocess with BE_CRASHER=1 to test main()
	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithVersionFlag")
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

func TestMainWithVersionShortFlag(t *testing.T) {
	// Test -v flag
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen", "-v"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithVersionShortFlag")
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

func TestMainWithInsufficientArgs(t *testing.T) {
	// Test with no arguments
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithInsufficientArgs")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("Expected non-zero exit code")
	}

	output := string(out)
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage message, got: %s", output)
	}
}

func TestMainWithOneArg(t *testing.T) {
	// Test with only one argument
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen", "test.yaml"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithOneArg")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("Expected non-zero exit code")
	}

	output := string(out)
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage message, got: %s", output)
	}
}

func TestMainWithNonExistentYamlFile(t *testing.T) {
	// Test with non-existent YAML file
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"puregen", "nonexistent.yaml", "template.tmpl"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
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
		os.Args = []string{"puregen", yamlFile.Name(), "nonexistent.tmpl"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
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
		os.Args = []string{"puregen", yamlFile.Name(), "nonexistent1.tmpl,nonexistent2.tmpl"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
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
