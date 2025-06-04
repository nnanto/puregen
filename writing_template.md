# Writing Templates for PureGen

This guide explains how to write Go templates for PureGen code generation, including the available data structures and helper functions.

## Template Structure

Every template must include a `metadata` template that defines the output file extension and optional type mappings:

```go
{{define "metadata"}}
{
  "extension": "go",
  "outputFileSuffix": "_generated",
  "typeMapping": {
    "string": "string",
    "int": "int64",
    "bool": "bool"
  }
}
{{end}}
```

**Fields:**
- `extension` (required): File extension for generated code (e.g., "go", "py", "java")
- `outputFileSuffix` (optional): Suffix to add to filename before extension
  - If schema name is "UserService" and extension is "go":
    - Without suffix: `userservice.go`
    - With suffix "_generated": `userservice_generated.go`
    - With suffix ".models": `userservice.models.go`
- `typeMapping` (optional): Map of IDL types to target language types

## Output File Naming

The `outputFileSuffix` field controls how generated files are named:

### Examples:
- Schema: "UserAPI", Extension: "go"
  - No suffix: `userapi.go`  
  - Suffix "_client": `userapi_client.go`
  - Suffix "_models": `userapi_models.go`
  - Suffix ".pb": `userapi.pb.go`

### Usage in Templates:
```go
{{define "metadata"}}
{
  "extension": "go",
  "outputFileSuffix": "_generated"
}
{{end}}
```

## Schema Data Structure

The main data passed to your template is a `Schema` struct with the following structure:

```go
type Schema struct {
    Name              string                 // Schema name
    Messages          []Message              // List of message definitions
    Services          []Service              // List of service definitions
    AdditionalContext map[string]interface{} // Additional JSON context passed via --additional-context-json
}

type Message struct {
    Name        string            // Message name
    Description string            // Optional description
    Fields      map[string]Field  // Field definitions
    Metadata    map[string]string // Optional metadata
}

type Field struct {
    Type        string            // Field type
    Description string            // Optional description
    Required    bool              // Whether field is required
    Repeated    bool              // Whether field is an array/slice
    Metadata    map[string]string // Optional metadata
}

type Service struct {
    Name        string            // Service name
    Description string            // Optional description
    Methods     []Method          // List of methods
    Metadata    map[string]string // Optional metadata
}

type Method struct {
    Name        string            // Method name
    Description string            // Optional description
    Input       string            // Input type
    Output      string            // Output type
    Streaming   bool              // Whether method uses streaming
    Metadata    map[string]string // Optional metadata
}
```

## Accessing Schema Data

### Basic Schema Information
```go
// Schema name
{{.Name}}

// Number of messages
{{len .Messages}}

// Number of services
{{len .Services}}

// Access additional context passed via --additional-context-json
{{.AdditionalContext}}
```

### Using Additional Context

The `AdditionalContext` field contains any JSON data passed via the `--additional-context-json` flag. This allows you to pass custom configuration, metadata, or other data to your templates.

#### Command Line Usage
```bash
puregen generate --schema schema.yaml --template template.go --additional-context-json '{"package":"myapi","version":"1.0.0","author":"John Doe"}'
```

#### Accessing in Templates
```go
// Access string values
Package: {{index .AdditionalContext "package"}}
Version: {{index .AdditionalContext "version"}}
Author: {{index .AdditionalContext "author"}}

// Type assertion for complex types
{{$config := index .AdditionalContext "config"}}
{{if $config}}
{{range $key, $value := $config}}
Config {{$key}}: {{$value}}
{{end}}
{{end}}

// Conditional logic based on additional context
{{if index .AdditionalContext "enableLogging"}}
// Generate logging code
{{end}}
```

#### Complex JSON Example
```bash
puregen generate --schema api.yaml --template go.tmpl --additional-context-json '{
  "package": "userapi",
  "version": "2.0.0",
  "features": {
    "authentication": true,
    "logging": true,
    "metrics": false
  },
  "endpoints": ["v1", "v2"],
  "author": {
    "name": "API Team",
    "email": "api-team@company.com"
  }
}'
```

```go
// In template
{{$features := index .AdditionalContext "features"}}
{{if index $features "authentication"}}
// Generate authentication middleware
{{end}}

{{$author := index .AdditionalContext "author"}}
// Author: {{index $author "name"}} <{{index $author "email"}}>

{{$endpoints := index .AdditionalContext "endpoints"}}
{{range $endpoints}}
// Endpoint: {{.}}
{{end}}
```

### Iterating Over Messages
```go
{{range .Messages}}
Message: {{.Name}}
{{if .Description}}Description: {{.Description}}{{end}}
Fields:
{{range $name, $field := .Fields}}
  - {{$name}}: {{$field.Type}}{{if $field.Required}} (required){{end}}
{{end}}
{{end}}
```

### Iterating Over Services
```go
{{range .Services}}
Service: {{.Name}}
{{range .Methods}}
  Method: {{.Name}}
  Input: {{.Input}} -> Output: {{.Output}}
{{end}}
{{end}}
```

## Available Template Functions

### String Case Conversion

#### `title` - Title Case
Converts string to title case (first letter of each word capitalized).
```go
{{title "hello world"}} // "Hello World"
{{.Name | title}}       // Apply to schema name
```

#### `capitalize` - Capitalize First Letter
Converts the first character to uppercase.
```go
{{capitalize "hello"}} // "Hello"
{{$field.Name | capitalize}}
```

#### `upper` - Uppercase
Converts entire string to uppercase.
```go
{{upper "hello"}} // "HELLO"
{{.Name | upper}}
```

#### `lower` - Lowercase
Converts entire string to lowercase.
```go
{{lower "HELLO"}} // "hello"
{{.Name | lower}}
```

#### `snake` - Snake Case
Converts string to snake_case.
```go
{{snake "HelloWorld"}} // "hello_world"
{{snake "XMLParser"}}  // "xml_parser"
```

#### `camel` - Camel Case
Converts string to camelCase.
```go
{{camel "hello_world"}} // "helloWorld"
{{camel "XML-parser"}}  // "xmlParser"
```

#### `pascal` - Pascal Case
Converts string to PascalCase.
```go
{{pascal "hello_world"}} // "HelloWorld"
{{pascal "xml-parser"}}  // "XmlParser"
```

#### `kebab` - Kebab Case
Converts string to kebab-case.
```go
{{kebab "HelloWorld"}} // "hello-world"
{{kebab "XMLParser"}}  // "xml-parser"
```

### String Manipulation

#### `split` - Split String
Splits string by delimiter.
```go
{{split "a,b,c" ","}} // ["a", "b", "c"]
{{range split .Name "_"}}{{.}}{{end}}
```

#### `trim` - Trim Whitespace
Removes leading and trailing whitespace.
```go
{{trim "  hello  "}} // "hello"
```

#### `join` - Join Strings
Joins slice of strings with delimiter.
```go
{{join (split "a,b,c" ",") "-"}} // "a-b-c"
```

#### `replace` - Replace All
Replaces all occurrences of old with new.
```go
{{replace "hello world" "world" "Go"}} // "hello Go"
```

### String Testing

#### `contains` - Contains Substring
Checks if string contains substring.
```go
{{if contains .Type "string"}}// Handle string type{{end}}
```

#### `hasPrefix` - Has Prefix
Checks if string starts with prefix.
```go
{{if hasPrefix .Name "User"}}// Handle User types{{end}}
```

#### `hasSuffix` - Has Suffix
Checks if string ends with suffix.
```go
{{if hasSuffix .Type "[]"}}// Handle array types{{end}}
```

### Comparison

#### `eq` - Equal
Checks if two values are equal.
```go
{{if eq .Type "string"}}// String type{{end}}
```

#### `ne` - Not Equal
Checks if two values are not equal.
```go
{{if ne .Type "string"}}// Not a string type{{end}}
```

### Map Access

#### `index` - Safe Map Access
Safely gets value from map with string key.
```go
{{index .Metadata "author"}} // Gets "author" from metadata map
```

## Complete Example Template

Here's a complete example that generates a Go struct with additional context:

```go
{{define "metadata"}}
{
  "extension": "go",
  "outputFileSuffix": "_models",
  "typeMapping": {
    "string": "string",
    "int": "int64",
    "bool": "bool",
    "double": "float64"
  }
}
{{end}}

package {{index .AdditionalContext "package" | default (.Name | lower)}}

// Generated code - do not modify
// Version: {{index .AdditionalContext "version" | default "unknown"}}
{{$author := index .AdditionalContext "author"}}{{if $author}}// Author: {{$author}}{{end}}

{{range .Messages}}
// {{.Name}}{{if .Description}} - {{.Description}}{{end}}
type {{.Name | pascal}} struct {
{{range $name, $field := .Fields}}
    {{$name | pascal}} {{if $field.Repeated}}[]{{end}}{{$field.Type}}{{if not $field.Required}} `json:"{{$name}},omitempty"`{{else}} `json:"{{$name}}"`{{end}}{{if $field.Description}} // {{$field.Description}}{{end}}
{{end}}
}
{{end}}

{{range .Services}}
// {{.Name}}{{if .Description}} - {{.Description}}{{end}}
type {{.Name | pascal}}Service interface {
{{range .Methods}}
    {{.Name | pascal}}({{if ne .Input "void"}}input {{.Input}}{{end}}) {{if ne .Output "void"}}({{.Output}}, error){{else}}error{{end}}{{if .Description}} // {{.Description}}{{end}}
{{end}}
}
{{end}}

{{$features := index .AdditionalContext "features"}}
{{if $features}}
// Features configuration
const (
{{if index $features "authentication"}}    AuthenticationEnabled = true{{end}}
{{if index $features "logging"}}    LoggingEnabled = true{{end}}
{{if index $features "metrics"}}    MetricsEnabled = true{{end}}
)
{{end}}
```

## Best Practices

1. **Always define metadata**: Include the `metadata` template with at least the `extension` field.

2. **Use type mapping**: Define type mappings to convert IDL types to target language types.

3. **Handle empty cases**: Check for empty descriptions, optional fields, etc.

4. **Use appropriate case conversion**: Match the target language's naming conventions.

5. **Add comments**: Include generated code warnings and documentation.

6. **Validate types**: Use conditional logic to handle different field types appropriately.

7. **Leverage additional context**: Use `--additional-context-json` to make templates more flexible and reusable.

8. **Provide defaults**: Use the `default` function when accessing additional context to handle missing values gracefully.

## Common Patterns

### Conditional Generation
```go
{{if .Description}}// {{.Description}}{{end}}
{{if $field.Required}}*{{end}}{{$field.Type}}
```

### Using Additional Context with Defaults
```go
// Package name from context or schema name as fallback
package {{index .AdditionalContext "package" | default (.Name | lower)}}

// Version with default
// Version: {{index .AdditionalContext "version" | default "1.0.0"}}

// Feature flags
{{$features := index .AdditionalContext "features"}}
{{if and $features (index $features "enableDebug")}}
// Debug mode enabled
{{end}}
```

### Complex Type Handling
```go
{{if eq $field.Type "string"}}
    // String-specific logic
{{else if eq $field.Type "int"}}
    // Integer-specific logic
{{else}}
    // Default handling
{{end}}
```

### Nested Iterations
```go
{{range .Services}}
Service: {{.Name}}
{{range .Methods}}
    Method: {{.Name}}
{{end}}
{{end}}
```

### Working with JSON Arrays in Additional Context
```go
{{$endpoints := index .AdditionalContext "endpoints"}}
{{if $endpoints}}
// Available endpoints:
{{range $endpoints}}
// - {{.}}
{{end}}
{{end}}
```
