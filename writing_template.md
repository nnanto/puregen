# Writing Templates for PureGen

This guide explains how to write Go templates for PureGen code generation, including the available data structures and helper functions.

## Template Structure

Every template must include a `metadata` template that defines the output file extension and optional type mappings:

```go
{{define "metadata"}}
{
  "extension": "go",
  "typeMapping": {
    "string": "string",
    "int": "int64",
    "bool": "bool"
  }
}
{{end}}
```

## Schema Data Structure

The main data passed to your template is a `Schema` struct with the following structure:

```go
type Schema struct {
    Name     string     // Schema name
    Messages []Message  // List of message definitions
    Services []Service  // List of service definitions
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

Here's a complete example that generates a Go struct:

```go
{{define "metadata"}}
{
  "extension": "go",
  "typeMapping": {
    "string": "string",
    "int": "int64",
    "bool": "bool",
    "double": "float64"
  }
}
{{end}}

package {{.Name | lower}}

// Generated code - do not modify

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
```

## Best Practices

1. **Always define metadata**: Include the `metadata` template with at least the `extension` field.

2. **Use type mapping**: Define type mappings to convert IDL types to target language types.

3. **Handle empty cases**: Check for empty descriptions, optional fields, etc.

4. **Use appropriate case conversion**: Match the target language's naming conventions.

5. **Add comments**: Include generated code warnings and documentation.

6. **Validate types**: Use conditional logic to handle different field types appropriately.

## Common Patterns

### Conditional Generation
```go
{{if .Description}}// {{.Description}}{{end}}
{{if $field.Required}}*{{end}}{{$field.Type}}
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
