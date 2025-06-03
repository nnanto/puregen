package idl

// SchemaYAML is used for parsing the original YAML format with maps
type SchemaYAML struct {
	Name     string                 `yaml:"name"`
	Version  string                 `yaml:"version"`
	Package  string                 `yaml:"package"`
	Messages map[string]Message     `yaml:"messages"`
	Services map[string]ServiceYAML `yaml:"services"`
}

// Message represents a data structure definition
type Message struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description,omitempty"`
	Fields      map[string]Field  `yaml:"fields"`
	Metadata    map[string]string `yaml:"metadata,omitempty"`
}

// Field represents a field in a message
type Field struct {
	Type        string            `yaml:"type"`
	Description string            `yaml:"description,omitempty"`
	Required    bool              `yaml:"required,omitempty"`
	Repeated    bool              `yaml:"repeated,omitempty"`
	Metadata    map[string]string `yaml:"metadata,omitempty"`
}

// ServiceYAML is used for parsing the original YAML format with method maps
type ServiceYAML struct {
	Description string            `yaml:"description,omitempty"`
	Methods     map[string]Method `yaml:"methods"`
	Metadata    map[string]string `yaml:"metadata,omitempty"`
}

// Method represents an RPC method
type Method struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description,omitempty"`
	Input       string            `yaml:"input"`
	Output      string            `yaml:"output"`
	Streaming   bool              `yaml:"streaming,omitempty"`
	Metadata    map[string]string `yaml:"metadata,omitempty"`
}
