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
		suffix    string
		schema    *idl.Schema
		expected  string
	}{
		{
			name:      "simple case",
			extension: "go",
			suffix:    "",
			schema:    &idl.Schema{Name: "User"},
			expected:  "user.go",
		},
		{
			name:      "with suffix",
			extension: "go",
			suffix:    "_client",
			schema:    &idl.Schema{Name: "User"},
			expected:  "user_client.go",
		},
		{
			name:      "uppercase schema name with suffix",
			extension: "py",
			suffix:    "_service",
			schema:    &idl.Schema{Name: "PRODUCT"},
			expected:  "product_service.py",
		},
		{
			name:      "mixed case schema name",
			extension: "js",
			suffix:    "",
			schema:    &idl.Schema{Name: "UserProfile"},
			expected:  "userprofile.js",
		},
		{
			name:      "empty extension with suffix",
			extension: "",
			suffix:    "_api",
			schema:    &idl.Schema{Name: "Test"},
			expected:  "test_api.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			result := g.getOutputFilename(tt.extension, tt.suffix, tt.schema)

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGenerateWithOutputFileSuffix(t *testing.T) {
	g := New()
	tmpDir := t.TempDir()

	template := `{{define "metadata"}}{"extension": "go", "outputFileSuffix": "_client"}{{end}}
package {{.Name | lower}}

type {{.Name}}Client struct {
    endpoint string
}
`
	schema := &idl.Schema{
		Name: "User",
		Messages: []idl.Message{
			{Name: "User", Fields: map[string]idl.Field{"ID": {Type: "int"}}},
		},
	}

	reader := strings.NewReader(template)
	err := g.Generate(schema, reader, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if file was created with suffix
	expectedPath := filepath.Join(tmpDir, "user_client.go")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("expected file user_client.go was not created")
	}

	// Check file content
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "type UserClient struct") {
		t.Errorf("generated content doesn't contain expected client struct")
	}
}

func TestGeneratorMetadata(t *testing.T) {
	g := New()
	tmpDir := t.TempDir()

	template := `{{define "metadata"}}{"extension": "go"}{{end}}
// Generated file: {{.GeneratorMetadata.OutputFilePath}}
package {{.Name | lower}}

type {{.Name}} struct {}
`
	schema := &idl.Schema{
		Name: "User",
		Messages: []idl.Message{
			{Name: "User", Fields: map[string]idl.Field{"ID": {Type: "int"}}},
		},
	}

	reader := strings.NewReader(template)
	err := g.Generate(schema, reader, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if file was created
	expectedPath := filepath.Join(tmpDir, "user.go")
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	contentStr := string(content)
	// Check that GeneratorMetadata.OutputFilePath was used in template
	if !strings.Contains(contentStr, "// Generated file: "+expectedPath) {
		t.Errorf("generated content doesn't contain expected generator metadata comment")
	}
}

func TestApplyTypeMapping(t *testing.T) {
	tests := []struct {
		name        string
		schema      *idl.Schema
		typeMapping map[string]string
		expected    map[string]string // field name -> expected type
	}{
		{
			name: "no type mapping",
			schema: &idl.Schema{
				Name: "User",
				Messages: []idl.Message{
					{
						Name: "User",
						Fields: map[string]idl.Field{
							"ID":   {Type: "int"},
							"Name": {Type: "string"},
						},
					},
				},
			},
			typeMapping: nil,
			expected:    map[string]string{"ID": "int", "Name": "string"},
		},
		{
			name: "with type mapping",
			schema: &idl.Schema{
				Name: "User",
				Messages: []idl.Message{
					{
						Name: "User",
						Fields: map[string]idl.Field{
							"ID":    {Type: "int"},
							"Score": {Type: "float"},
							"Name":  {Type: "string"},
						},
					},
				},
			},
			typeMapping: map[string]string{
				"int":   "int64",
				"float": "float64",
			},
			expected: map[string]string{"ID": "int64", "Score": "float64", "Name": "string"},
		},
		{
			name: "service method type mapping",
			schema: &idl.Schema{
				Name: "UserService",
				Services: []idl.Service{
					{
						Name: "User",
						Methods: []idl.Method{
							{Name: "GetUser", Input: "int", Output: "User"},
							{Name: "UpdateUser", Input: "User", Output: "bool"},
						},
					},
				},
			},
			typeMapping: map[string]string{
				"int":  "int64",
				"bool": "boolean",
			},
			expected: map[string]string{"GetUser.Input": "int64", "UpdateUser.Output": "boolean"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			result := g.applyTypeMapping(tt.schema, tt.typeMapping)

			// Check message field types
			if len(result.Messages) > 0 {
				for fieldName, expectedType := range tt.expected {
					if strings.Contains(fieldName, ".") {
						continue // Skip service method checks in this loop
					}
					actualType := result.Messages[0].Fields[fieldName].Type
					if actualType != expectedType {
						t.Errorf("field %s: expected type %q, got %q", fieldName, expectedType, actualType)
					}
				}
			}

			// Check service method types
			if len(result.Services) > 0 {
				for key, expectedType := range tt.expected {
					if !strings.Contains(key, ".") {
						continue // Skip message field checks in this loop
					}
					parts := strings.Split(key, ".")
					methodName, fieldType := parts[0], parts[1]

					for _, method := range result.Services[0].Methods {
						if method.Name == methodName {
							var actualType string
							if fieldType == "Input" {
								actualType = method.Input
							} else if fieldType == "Output" {
								actualType = method.Output
							}
							if actualType != expectedType {
								t.Errorf("method %s.%s: expected type %q, got %q", methodName, fieldType, expectedType, actualType)
							}
							break
						}
					}
				}
			}
		})
	}
}

func TestMapType(t *testing.T) {
	tests := []struct {
		name         string
		originalType string
		typeMapping  map[string]string
		expected     string
	}{
		{
			name:         "type exists in mapping",
			originalType: "int",
			typeMapping:  map[string]string{"int": "int64", "string": "varchar"},
			expected:     "int64",
		},
		{
			name:         "type not in mapping",
			originalType: "bool",
			typeMapping:  map[string]string{"int": "int64", "string": "varchar"},
			expected:     "bool",
		},
		{
			name:         "empty mapping",
			originalType: "string",
			typeMapping:  map[string]string{},
			expected:     "string",
		},
		{
			name:         "nil mapping",
			originalType: "float",
			typeMapping:  nil,
			expected:     "float",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			result := g.mapType(tt.originalType, tt.typeMapping)

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGenerateWithTypeMapping(t *testing.T) {
	g := New()
	tmpDir := t.TempDir()

	template := `{{define "metadata"}}{"extension": "go", "typeMapping": {"int": "int64", "string": "varchar"}}{{end}}
package {{.Name | lower}}

{{range .Messages}}
type {{.Name}} struct {
{{- range $name, $field := .Fields}}
    {{$name}} {{$field.Type}}
{{- end}}
}
{{end}}
`
	schema := &idl.Schema{
		Name: "User",
		Messages: []idl.Message{
			{
				Name: "User",
				Fields: map[string]idl.Field{
					"ID":     {Type: "int"},
					"Name":   {Type: "string"},
					"Active": {Type: "bool"},
				},
			},
		},
	}

	reader := strings.NewReader(template)
	err := g.Generate(schema, reader, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check file content
	expectedPath := filepath.Join(tmpDir, "user.go")
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	contentStr := string(content)
	// Check that type mapping was applied
	if !strings.Contains(contentStr, "ID int64") {
		t.Errorf("type mapping for int -> int64 was not applied")
	}
	if !strings.Contains(contentStr, "Name varchar") {
		t.Errorf("type mapping for string -> varchar was not applied")
	}
	if !strings.Contains(contentStr, "Active bool") {
		t.Errorf("unmapped type bool should remain unchanged")
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
