# PureGen - IDL Parser and Code Generator

PureGen is a code generation tool that transforms YAML into code for any target language. Define your data structures and services once, then generate consistent, type-safe code across multiple languages and platforms.

## Why PureGen?

- **Language Agnostic**: Generate code for any language using Go templates
- **Type Safety**: Ensures consistent types across your entire system
- **DRY Principle**: Define once, generate everywhere
- **Flexible Templates**: Write custom templates for your specific needs
- **Service Contracts**: Define RPC services alongside your data structures
- **Field Validation**: Built-in support for required fields, arrays, and custom types
- **Interactive Creation**: Create IDL files interactively with the creator command
- **Validation**: Validate IDL files for type consistency

## Installation

### Released Binaries (Recommended)

Download the latest release for your platform:

**Linux (x64)**
```bash
curl -L https://github.com/nnanto/puregen/releases/download/latest/puregen-linux-amd64.tar.gz | tar -xz
sudo mv puregen-linux-amd64 /usr/local/bin/puregen
```

**macOS (Apple Silicon)**
```bash
curl -L https://github.com/nnanto/puregen/releases/download/latest/puregen-darwin-arm64.tar.gz | tar -xz
sudo mv puregen-darwin-arm64 /usr/local/bin/puregen
```

**macOS (Intel)**
```bash
curl -L https://github.com/nnanto/puregen/releases/download/latest/puregen-darwin-amd64.tar.gz | tar -xz
sudo mv puregen-darwin-amd64 /usr/local/bin/puregen
```

**Windows**
Download `puregen-windows-amd64.zip` from the [releases page](https://github.com/nnanto/puregen/releases/latest), extract it, and add the executable to your PATH.

### From Source

```bash
git clone https://github.com/nnanto/puregen.git
cd puregen
go mod tidy
make build
```

## Quick Start

1. **Create an IDL file** (`user_service.yaml`):
```yaml
name: "UserService"
version: "1.0"

messages:
  User:
    description: "User entity with profile information"
    fields:
      id:
        type: "string"
        required: true
        description: "Unique user identifier"
      name:
        type: "string"
        required: true
      email:
        type: "string"
        required: true
      tags:
        type: "string"
        repeated: true
        description: "User tags for categorization"
      active:
        type: "bool"
        required: true

  GetUserRequest:
    fields:
      id:
        type: "string"
        required: true

services:
  UserService:
    description: "User management service"
    methods:
      GetUser:
        description: "Retrieve a user by ID"
        input: "GetUserRequest"
        output: "User"
      ListUsers:
        description: "Stream all users"
        input: "void"
        output: "User"
        streaming: true
```

2. **Generate code using a template**:
```bash
puregen generate --input user_service.yaml --templates templates/typescript.tmpl --output ./generated
```

3. **Result**: PureGen generates type-safe code in your target language!

## Usage

### Command Line

```bash
# Basic code generation
puregen generate --input <idl-file> --templates <template-file> --output <output-directory>

# Short flags
puregen generate -i <idl-file> -t <template-file> -o <output-directory>

# Multiple templates (comma-separated)
puregen generate --input <idl-file> --templates <template1,template2,template3> --output <output-directory>

# Interactive IDL creation
puregen creator --output-file <output-yaml-file>

# Validate IDL file
puregen validate --input <idl-file>

# Examples
puregen generate -i service.yaml -t templates/go.tmpl -o ./gen
puregen generate -i api.yaml -t templates/typescript.tmpl -o ./src/types
puregen creator -o my_service.yaml
puregen validate -i my_service.yaml

# Generate multiple languages at once
puregen generate -i user_service.yaml -t templates/go.tmpl,templates/typescript.tmpl,templates/python.tmpl -o ./generated

# Check version
puregen version
```

### Available Commands

- `puregen generate` - Generate code from IDL files using templates
- `puregen creator` - Interactively create a new IDL file
- `puregen validate` - Validate an IDL file for type consistency and potential issues
- `puregen version` - Show version information
- `puregen help` - Show help information

### Creator Mode

Use the interactive creator to build IDL files step by step:

```bash
puregen creator --output-file my_service.yaml
```

The creator will guide you through:
- Setting service name and version
- Defining messages with typed fields
- Creating service methods
- Setting up field validation rules

### Validation

Validate your IDL files to catch potential issues:

```bash
puregen validate --input my_service.yaml
```

The validator checks for:
- **Type consistency**: Ensures all field types are either primitive types or defined messages
- **Missing references**: Warns about custom types that aren't defined in the Messages section
- **Syntax errors**: Basic YAML and structure validation

Example validation output:
```
Warning: Field 'profile' in message 'User' has type 'UserProfile' which is not primitive and not defined in Messages
Warning: Field 'settings' in message 'User' has type '[]AppSettings' which is not primitive and not defined in Messages
```

### Available Templates

PureGen comes with sample templates :
- `templates/go.tmpl` - Go structs and interfaces
- `templates/typescript.tmpl` - TypeScript interfaces
- `templates/python.tmpl` - Python dataclasses

### Generated Code Example

From the IDL above, TypeScript generation produces:

```typescript
// Generated by PureGen - do not modify

export interface User {
  /** Unique user identifier */
  id: string;
  name: string;
  email: string;
  /** User tags for categorization */
  tags?: string[];
  active: boolean;
}

export interface GetUserRequest {
  id: string;
}

/** User management service */
export interface UserService {
  /** Retrieve a user by ID */
  getUser(input: GetUserRequest): Promise<User>;
  /** Stream all users */
  listUsers(): AsyncIterable<User>;
}
```


## IDL Format

### Messages

Define your data structures with typed fields:

```yaml
messages:
  User:
    description: "User entity"
    fields:
      id:
        type: "string"
        required: true
        description: "Primary key"
      tags:
        type: "string"
        repeated: true
      metadata:
        type: "UserMetadata"  # Custom type reference
```

### Services

Define RPC services with methods:

```yaml
services:
  UserService:
    description: "User management service"
    methods:
      GetUser:
        description: "Fetch user by ID"
        input: "GetUserRequest"
        output: "User"
      StreamUsers:
        description: "Stream all users"
        input: "ListUsersRequest"
        output: "User"
        streaming: true
      DeleteUser:
        input: "DeleteUserRequest"
        output: "void"  # No return value
```

## Supported Field Types

- **Primitives**: `string`, `int`, `int32`, `int64`, `float`, `float32`, `float64`, `bool`, `byte`, `bytes`
- **Custom Types**: Reference other messages by name
- **Arrays**: Use `repeated: true` for array/slice types or prefix with `[]`
- **Optional**: Fields without `required: true` are optional

## Writing Custom Templates

Create your own templates for any language or framework. See [Writing Templates Guide](./writing_template.md) for detailed documentation.

Example template snippet:
```go
{{range .Messages}}
class {{.Name | pascal}} {
{{range $name, $field := .Fields}}
  {{$name}}: {{$field.Type}}{{if not $field.Required}}?{{end}};
{{end}}
}
{{end}}
```

## Testing

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for your changes
4. Submit a pull request

## License

MIT License - see LICENSE file for details.
