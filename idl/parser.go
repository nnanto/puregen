package idl

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type GeneratorMetadata struct {
	OutputFilePath string `json:"outputFilePath"`
}

// Schema represents the root IDL schema
type Schema struct {
	Name              string                 `yaml:"name"`
	Version           string                 `yaml:"version"`
	Package           string                 `yaml:"package"`
	Messages          []Message              `yaml:"-"`
	Services          []Service              `yaml:"-"`
	GeneratorMetadata *GeneratorMetadata     `yaml:"-"`
	AdditionalContext map[string]interface{} `yaml:"-"`
}

// Service represents a service with RPC methods
type Service struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description,omitempty"`
	Methods     []Method          `yaml:"-"`
	Metadata    map[string]string `yaml:"metadata,omitempty"`
}

// ToSchema converts SchemaYAML to Schema by converting maps to slices
func (sy *SchemaYAML) ToSchema() *Schema {
	schema := &Schema{
		Name:     sy.Name,
		Version:  sy.Version,
		Package:  sy.Package,
		Messages: make([]Message, 0, len(sy.Messages)),
		Services: make([]Service, 0, len(sy.Services)),
	}

	// Convert messages map to slice
	for name, message := range sy.Messages {
		message.Name = name
		schema.Messages = append(schema.Messages, message)
	}

	// Convert services map to slice
	for name, serviceYAML := range sy.Services {
		service := Service{
			Name:        name,
			Description: serviceYAML.Description,
			Methods:     make([]Method, 0, len(serviceYAML.Methods)),
			Metadata:    serviceYAML.Metadata,
		}

		// Convert methods map to slice
		for methodName, method := range serviceYAML.Methods {
			method.Name = methodName
			service.Methods = append(service.Methods, method)
		}

		schema.Services = append(schema.Services, service)
	}

	return schema
}

// Parser handles IDL file parsing
type Parser struct{}

// NewParser creates a new IDL parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile parses an IDL file from disk
func (p *Parser) ParseFile(filename string) (*Schema, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return p.Parse(data)
}

// Parse parses IDL content from bytes
func (p *Parser) Parse(data []byte) (*Schema, error) {
	var schemaYAML SchemaYAML

	if err := yaml.Unmarshal(data, &schemaYAML); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	schema := schemaYAML.ToSchema()

	if err := p.validate(schema); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return schema, nil
}

// validate performs basic validation on the parsed schema
func (p *Parser) validate(schema *Schema) error {
	if schema.Name == "" {
		return fmt.Errorf("schema name is required")
	}

	// Validate message field types
	for _, msg := range schema.Messages {
		for fieldName, field := range msg.Fields {
			if field.Type == "" {
				return fmt.Errorf("field type is required for %s.%s", msg.Name, fieldName)
			}
		}
	}

	// Validate service method input/output types
	for _, service := range schema.Services {
		for _, method := range service.Methods {
			if method.Input == "" {
				return fmt.Errorf("input type is required for %s.%s", service.Name, method.Name)
			}
			if method.Output == "" {
				return fmt.Errorf("output type is required for %s.%s", service.Name, method.Name)
			}
		}
	}

	return nil
}
