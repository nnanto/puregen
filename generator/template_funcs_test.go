package generator

import (
	"strings"
	"testing"
	"text/template"
)

func TestTemplateFuncs(t *testing.T) {
	funcs := TemplateFuncs()

	// Test that we have the expected functions
	expectedFuncs := []string{"lower", "upper", "title"}
	for _, funcName := range expectedFuncs {
		if _, exists := funcs[funcName]; !exists {
			t.Errorf("expected function %q not found in template functions", funcName)
		}
	}
}

func TestTemplateFuncsIntegration(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     interface{}
		expected string
	}{
		{
			name:     "lower function",
			template: `{{.Name | lower}}`,
			data:     struct{ Name string }{Name: "HELLO"},
			expected: "hello",
		},
		{
			name:     "upper function",
			template: `{{.Name | upper}}`,
			data:     struct{ Name string }{Name: "hello"},
			expected: "HELLO",
		},
		{
			name:     "title function",
			template: `{{.Name | title}}`,
			data:     struct{ Name string }{Name: "hello world"},
			expected: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New("test").Funcs(TemplateFuncs()).Parse(tt.template)
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			var buf strings.Builder
			err = tmpl.Execute(&buf, tt.data)
			if err != nil {
				t.Fatalf("failed to execute template: %v", err)
			}

			result := buf.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
