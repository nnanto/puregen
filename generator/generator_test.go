package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/nnanto/puregen/idl"
)

func TestNew(t *testing.T) {
	g := New()
	if g == nil {
		t.Fatal("New() returned nil")
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name         string
		template     string
		schema       *idl.Schema
		expectError  bool
		expectedFile string
	}{
		{
			name: "valid template with metadata",
			template: `{{define "metadata"}}{"extension": "go"}{{end}}
package {{.Name | lower}}

{{range .Messages}}
type {{.Name}} struct {
{{- range $name, $field := .Fields}}
    {{$name}} {{$field.Type}}
{{- end}}
}
{{end}}`,
			schema: &idl.Schema{
				Name: "User",
				Messages: []idl.Message{
					{
						Name: "User",
						Fields: map[string]idl.Field{
							"ID":   {Type: "int", Required: true},
							"Name": {Type: "string", Required: true},
						},
					},
				},
			},
			expectError:  false,
			expectedFile: "user.go",
		},
		{
			name: "template without metadata",
			template: `package {{.Name | lower}}

type {{.Name}} struct {}`,
			schema: &idl.Schema{
				Name: "Product",
			},
			expectError: true,
		},
		{
			name: "template with empty extension",
			template: `{{define "metadata"}}{"extension": ""}{{end}}
package test`,
			schema:      &idl.Schema{Name: "Test"},
			expectError: true,
		},
		{
			name: "complex schema with multiple messages and services",
			template: `{{define "metadata"}}{"extension": "go"}{{end}}
package {{.Name | lower}}

{{range .Messages}}
type {{.Name}} struct {
{{- range $name, $field := .Fields}}
    {{$name}} {{$field.Type}}{{if $field.Required}} // required{{end}}
{{- end}}
}
{{end}}

{{range .Services}}
type {{.Name}}Service interface {
{{- range .Methods}}
    {{.Name}}({{.Input}}) {{.Output}}
{{- end}}
}
{{end}}`,
			schema: &idl.Schema{
				Name: "UserService",
				Messages: []idl.Message{
					{
						Name: "User",
						Fields: map[string]idl.Field{
							"ID":    {Type: "int64", Required: true},
							"Email": {Type: "string", Required: true},
							"Name":  {Type: "string", Required: false},
							"Tags":  {Type: "[]string", Repeated: true},
						},
					},
					{
						Name: "CreateUserRequest",
						Fields: map[string]idl.Field{
							"Email": {Type: "string", Required: true},
							"Name":  {Type: "string", Required: true},
						},
					},
				},
				Services: []idl.Service{
					{
						Name: "User",
						Methods: []idl.Method{
							{Name: "CreateUser", Input: "CreateUserRequest", Output: "User"},
							{Name: "GetUser", Input: "int64", Output: "User"},
						},
					},
				},
			},
			expectError:  false,
			expectedFile: "userservice.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			tmpDir := t.TempDir()

			reader := strings.NewReader(tt.template)
			err := g.Generate(tt.schema, reader, tmpDir)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check if file was created
			expectedPath := filepath.Join(tmpDir, tt.expectedFile)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Fatalf("expected file %s was not created", tt.expectedFile)
			}

			// Check file content
			content, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("failed to read generated file: %v", err)
			}

			contentStr := string(content)
			if !strings.Contains(contentStr, "package user") {
				t.Errorf("generated content doesn't contain expected package declaration")
			}
		})
	}
}

func TestGetTemplateWithMetadata(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		expectedExt string
		expectError bool
	}{
		{
			name:        "valid template with metadata",
			template:    `{{define "metadata"}}{"extension": "go"}{{end}}Hello {{.Name}}`,
			expectedExt: "go",
			expectError: false,
		},
		{
			name:        "template without metadata",
			template:    `Hello {{.Name}}`,
			expectedExt: "",
			expectError: true,
		},
		{
			name:        "template with invalid JSON metadata",
			template:    `{{define "metadata"}}invalid json{{end}}Hello`,
			expectedExt: "",
			expectError: true,
		},
		{
			name:        "template with empty metadata",
			template:    `{{define "metadata"}}{{end}}Hello`,
			expectedExt: "",
			expectError: true,
		},
		{
			name:        "invalid template syntax",
			template:    `{{invalid syntax`,
			expectedExt: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			reader := strings.NewReader(tt.template)

			tmpl, metadata, err := g.getTemplateWithMetadata(reader)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tmpl == nil {
				t.Fatal("template is nil")
			}

			if metadata.Extension != tt.expectedExt {
				t.Errorf("expected extension %q, got %q", tt.expectedExt, metadata.Extension)
			}
		})
	}
}

func TestExtractMetadataFromTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		expectedExt string
		expectError bool
	}{
		{
			name:        "valid metadata",
			template:    `{{define "metadata"}}{"extension": "py"}{{end}}`,
			expectedExt: "py",
			expectError: false,
		},
		{
			name:        "no metadata template",
			template:    `Hello World`,
			expectedExt: "",
			expectError: true,
		},
		{
			name:        "invalid JSON in metadata",
			template:    `{{define "metadata"}}not json{{end}}`,
			expectedExt: "",
			expectError: true,
		},
		{
			name:        "empty metadata",
			template:    `{{define "metadata"}}{{end}}`,
			expectedExt: "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()

			tmpl, err := template.New("test").Funcs(TemplateFuncs()).Parse(tt.template)
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			metadata, err := g.extractMetadataFromTemplate(tmpl)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if metadata.Extension != tt.expectedExt {
				t.Errorf("expected extension %q, got %q", tt.expectedExt, metadata.Extension)
			}
		})
	}
}

func TestGetOutputFilename(t *testing.T) {
	tests := []struct {
		name      string
		extension string
		schema    *idl.Schema
		expected  string
	}{
		{
			name:      "simple case",
			extension: "go",
			schema:    &idl.Schema{Name: "User"},
			expected:  "user.go",
		},
		{
			name:      "uppercase schema name",
			extension: "py",
			schema:    &idl.Schema{Name: "PRODUCT"},
			expected:  "product.py",
		},
		{
			name:      "mixed case schema name",
			extension: "js",
			schema:    &idl.Schema{Name: "UserProfile"},
			expected:  "userprofile.js",
		},
		{
			name:      "empty extension",
			extension: "",
			schema:    &idl.Schema{Name: "Test"},
			expected:  "test.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			result := g.getOutputFilename(tt.extension, "", tt.schema)

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGenerateFileCreationError(t *testing.T) {
	g := New()
	schema := &idl.Schema{Name: "Test"}
	template := `{{define "metadata"}}{"extension": "go"}{{end}}package test`
	reader := strings.NewReader(template)

	// Try to create file in non-existent directory without permission
	invalidDir := "/invalid/path/that/does/not/exist"
	err := g.Generate(schema, reader, invalidDir)

	if err == nil {
		t.Fatal("expected error for invalid output directory")
	}

	if !strings.Contains(err.Error(), "failed to create output directory") {
		t.Errorf("expected directory creation error, got: %v", err)
	}
}
