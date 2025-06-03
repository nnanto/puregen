package generator

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nnanto/puregen/idl"
)

type Generator struct {
}

type TemplateMetadata struct {
	Extension        string            `json:"extension"`
	OutputFileSuffix string            `json:"outputFileSuffix"`
	TypeMapping      map[string]string `json:"typeMapping"`
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(schema *idl.Schema, templateReader io.Reader, outputDir string) error {
	tmpl, metadata, err := g.getTemplateWithMetadata(templateReader)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	extension := metadata.Extension

	// Apply type mapping if present
	transformedSchema := g.applyTypeMapping(schema, metadata.TypeMapping)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputFile := filepath.Join(outputDir, g.getOutputFilename(extension, metadata.OutputFileSuffix, schema))
	transformedSchema.GeneratorMetadata = &idl.GeneratorMetadata{
		OutputFilePath: outputFile,
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, transformedSchema); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("Generated %s code in %s\n", extension, outputFile)
	return nil
}

func (g *Generator) getTemplateWithMetadata(reader io.Reader) (*template.Template, *TemplateMetadata, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read template from reader: %w", err)
	}

	// Parse template once with all functions
	tmpl, err := template.New("main").Funcs(TemplateFuncs()).Parse(string(content))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Extract metadata from the already parsed template
	metadata, err := g.extractMetadataFromTemplate(tmpl)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract template metadata: %w", err)
	}

	if metadata.Extension == "" {
		return nil, nil, fmt.Errorf("no extension found in template metadata")
	}

	return tmpl, metadata, nil
}

func (g *Generator) extractMetadataFromTemplate(tmpl *template.Template) (*TemplateMetadata, error) {
	// Execute the metadata template to get JSON
	var buf strings.Builder
	if err := tmpl.ExecuteTemplate(&buf, "metadata", nil); err != nil {
		return nil, fmt.Errorf("failed to execute metadata template: %w", err)
	}

	jsonContent := strings.TrimSpace(buf.String())
	if jsonContent == "" {
		return &TemplateMetadata{}, nil
	}

	var metadata TemplateMetadata
	if err := json.Unmarshal([]byte(jsonContent), &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata JSON: %w", err)
	}

	return &metadata, nil
}

func (g *Generator) getOutputFilename(extension, suffix string, schema *idl.Schema) string {
	baseName := strings.ToLower(schema.Name)
	fileName := baseName + "." + extension
	if suffix != "" {
		fileName = baseName + suffix + "." + extension
	}
	return fileName
}

func (g *Generator) applyTypeMapping(schema *idl.Schema, typeMapping map[string]string) *idl.Schema {
	if len(typeMapping) == 0 {
		return schema
	}

	// Create a deep copy of the schema
	transformedSchema := &idl.Schema{
		Name:     schema.Name,
		Messages: make([]idl.Message, len(schema.Messages)),
		Services: make([]idl.Service, len(schema.Services)),
	}

	// Transform message field types
	for i, message := range schema.Messages {
		transformedSchema.Messages[i] = idl.Message{
			Name:        message.Name,
			Description: message.Description,
			Fields:      make(map[string]idl.Field),
		}

		for fieldName, field := range message.Fields {
			transformedType := g.mapType(field.Type, typeMapping)
			transformedSchema.Messages[i].Fields[fieldName] = idl.Field{
				Type:        transformedType,
				Description: field.Description,
			}
		}
	}

	// Transform service method types
	for i, service := range schema.Services {
		transformedSchema.Services[i] = idl.Service{
			Name:        service.Name,
			Description: service.Description,
			Methods:     make([]idl.Method, len(service.Methods)),
		}

		for j, method := range service.Methods {
			transformedSchema.Services[i].Methods[j] = idl.Method{
				Name:        method.Name,
				Description: method.Description,
				Input:       g.mapType(method.Input, typeMapping),
				Output:      g.mapType(method.Output, typeMapping),
			}
		}
	}

	return transformedSchema
}

func (g *Generator) mapType(originalType string, typeMapping map[string]string) string {
	if mappedType, exists := typeMapping[originalType]; exists {
		return mappedType
	}
	return originalType
}
