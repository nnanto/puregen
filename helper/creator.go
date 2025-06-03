package helper

import (
	"fmt"
	"os"

	"github.com/nnanto/puregen/idl"
	"gopkg.in/yaml.v3"
)

// Creator handles interactive creation and editing of YAML IDL files
type Creator struct {
	filePath   string
	schema     *idl.SchemaYAML
	fileExists bool
}

// NewCreator creates a new Creator instance
func NewCreator(filePath string) *Creator {
	return &Creator{
		filePath: filePath,
		schema: &idl.SchemaYAML{
			Messages: make(map[string]idl.Message),
			Services: make(map[string]idl.ServiceYAML),
		},
	}
}

// LoadOrCreate loads an existing file or creates a new one
func (c *Creator) LoadOrCreate() error {
	if _, err := os.Stat(c.filePath); os.IsNotExist(err) {
		// File doesn't exist, create new schema with defaults
		c.schema.Name = "MySchema"
		c.schema.Version = "1.0.0"
		c.schema.Package = "main"
		c.fileExists = false
		return nil
	} else if err != nil {
		return err
	}

	// File exists, load it
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, c.schema); err != nil {
		return err
	}

	// Initialize maps if they're nil
	if c.schema.Messages == nil {
		c.schema.Messages = make(map[string]idl.Message)
	}
	if c.schema.Services == nil {
		c.schema.Services = make(map[string]idl.ServiceYAML)
	}

	c.fileExists = true
	return nil
}

// FileExists returns true if the file was loaded from disk
func (c *Creator) FileExists() bool {
	return c.fileExists
}

// Save writes the current schema to the file
func (c *Creator) Save() error {
	data, err := yaml.Marshal(c.schema)
	if err != nil {
		return err
	}

	return os.WriteFile(c.filePath, data, 0644)
}

// CreateMessage creates a new message
func (c *Creator) CreateMessage(name, description string) idl.Message {
	return idl.Message{
		Name:        name,
		Description: description,
		Fields:      make(map[string]idl.Field),
		Metadata:    make(map[string]string),
	}
}

// CreateField creates a new field
func (c *Creator) CreateField(fieldType, description string, required, repeated bool) idl.Field {
	return idl.Field{
		Type:        fieldType,
		Description: description,
		Required:    required,
		Repeated:    repeated,
		Metadata:    make(map[string]string),
	}
}

// CreateService creates a new service
func (c *Creator) CreateService(description string) idl.ServiceYAML {
	return idl.ServiceYAML{
		Description: description,
		Methods:     make(map[string]idl.Method),
		Metadata:    make(map[string]string),
	}
}

// CreateMethod creates a new method
func (c *Creator) CreateMethod(name, description, input, output string, streaming bool) idl.Method {
	return idl.Method{
		Name:        name,
		Description: description,
		Input:       input,
		Output:      output,
		Streaming:   streaming,
		Metadata:    make(map[string]string),
	}
}

// AddMessage adds a message to the schema
func (c *Creator) AddMessage(name string, message idl.Message) error {
	if name == "" {
		return fmt.Errorf("message name cannot be empty")
	}
	c.schema.Messages[name] = message
	return nil
}

// AddService adds a service to the schema
func (c *Creator) AddService(name string, service idl.ServiceYAML) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	c.schema.Services[name] = service
	return nil
}

// IsValidType checks if a type is valid (primitive or defined message)
func (c *Creator) IsValidType(typeName string) bool {
	// Check primitive types
	primitives := []string{
		"string", "int32", "int64", "uint32", "uint64",
		"float32", "float64", "bool", "bytes",
	}

	for _, primitive := range primitives {
		if typeName == primitive {
			return true
		}
	}

	// Check if it's a defined message
	_, exists := c.schema.Messages[typeName]
	return exists
}

// ShowSchema displays the current schema
func (c *Creator) ShowSchema() {
	fmt.Printf("\n=== Schema: %s ===\n", c.schema.Name)
	fmt.Printf("Version: %s\n", c.schema.Version)
	fmt.Printf("Package: %s\n", c.schema.Package)

	if len(c.schema.Messages) > 0 {
		fmt.Println("\nMessages:")
		for name, message := range c.schema.Messages {
			fmt.Printf("  %s: %s\n", name, message.Description)
			for fieldName, field := range message.Fields {
				required := ""
				if field.Required {
					required = " (required)"
				}
				repeated := ""
				if field.Repeated {
					repeated = " (repeated)"
				}
				fmt.Printf("    - %s: %s%s%s\n", fieldName, field.Type, required, repeated)
			}
		}
	}

	if len(c.schema.Services) > 0 {
		fmt.Println("\nServices:")
		for name, service := range c.schema.Services {
			fmt.Printf("  %s: %s\n", name, service.Description)
			for methodName, method := range service.Methods {
				streaming := ""
				if method.Streaming {
					streaming = " (streaming)"
				}
				fmt.Printf("    - %s: %s -> %s%s\n", methodName, method.Input, method.Output, streaming)
			}
		}
	}

	if len(c.schema.Messages) == 0 && len(c.schema.Services) == 0 {
		fmt.Println("\nNo messages or services defined yet.")
	}
}

// GetSchema returns the current schema
func (c *Creator) GetSchema() *idl.SchemaYAML {
	return c.schema
}
