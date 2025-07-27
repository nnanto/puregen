# puregen - Protobuf Code Generator

puregen is a protobuf plugin that generates **simple, dependency-minimal code** for Go, Java, and Python from `.proto` files. The generated code focuses on simplicity and uses built-in language features rather than heavy dependencies.

puregen is ideal for projects that need simple, readable generated code without heavy protobuf runtime dependencies, with the flexibility to use any transport mechanism (HTTP, gRPC, message queues, etc.).

## Comparison with Standard Generators

| Feature | puregen | protoc-gen-go | protoc-gen-java |
|---------|-----------|---------------|-----------------|
| Dependencies | Minimal | protobuf runtime | protobuf runtime |
| Code size | Small | Large | Large |
| JSON support | Built-in | Requires jsonpb | Requires additional libs |
| Transport abstraction | Pluggable | gRPC only | gRPC only |
| Customization | Easy | Complex | Complex |
| Learning curve | Low | Medium | Medium |


## Features

- **Multi-language support**: Generate code for Go, Java, and Python
- **Minimal dependencies**: Uses only built-in libraries and standard patterns
- **Simple data structures**: Generated classes/structs are easy to understand and modify
- **JSON serialization**: Built-in JSON marshaling/unmarshaling support
- **Service interfaces**: Clean interface definitions for RPC services
- **Method metadata support**: Extract metadata from service method comments (*//metadata:{...}*) for HTTP routing, authorization, etc. [See details](#method-metadata-support)
- **Client generation**: Ready-to-use clients with pluggable transport. [See details](#using-the-generated-code)

## Installation

### Prerequisites

- Protocol Buffers compiler (`protoc`). Install it from [the official site](https://protobuf.dev/installation/).

### Build the Plugin

#### Pre-built Binaries (Recommended)

You can download pre-built binaries for your platform from the [releases page](https://github.com/nnanto/puregen/releases).

```bash
curl -L https://github.com/nnanto/puregen/releases/download/latest/protoc-gen-puregen-linux-amd64.tar.gz | tar -xz
          sudo mv protoc-gen-puregen-* /usr/local/bin/protoc-gen-puregen

# Or for macOS
curl -L https://github.com/nnanto/puregen/releases/download/latest/protoc-gen-puregen-darwin-amd64.tar.gz | tar -xz
sudo mv protoc-gen-puregen-* /usr/local/bin/protoc-gen-puregen

# Run the generator
protoc --puregen_out=./examples/generated --puregen_opt=language=python examples/proto/*.proto
```

#### From Source

```bash
# Clone the repository
git clone https://github.com/nnanto/puregen
cd puregen

# Build the plugin
go build -o protoc-gen-puregen ./cmd/protoc-gen-puregen

# Make it available in your PATH (optional)
sudo mv protoc-gen-puregen /usr/local/bin/
```

## Usage

### Basic Generation

Generate code for all supported languages:

```bash
protoc --puregen_out=./generated --puregen_opt=language=all user.proto
```

Generate code for a specific language:

```bash
# Go only
protoc --puregen_out=./generated --puregen_opt=language=go user.proto

# Java only
protoc --puregen_out=./generated --puregen_opt=language=java user.proto

# Python only
protoc --puregen_out=./generated --puregen_opt=language=python user.proto
```

### Example Proto File

```protobuf
syntax = "proto3";

package example.v1;

option go_package = "github.com/nnanto/puregen/examples/proto/gen/go";
option java_package = "com.example.proto.v1";

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
  bool is_active = 4;
  UserProfile profile = 5;
}

message UserProfile {
  string bio = 1;
  string avatar_url = 2;
  int64 created_at = 3;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
  UserProfile profile = 3;
}

message CreateUserResponse {
  User user = 1;
  string message = 2;
}

message GetUserRequest {
  int32 id = 1;
}

message GetUserResponse {
  User user = 1;
}

service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}
```

## Using the Generated Code

### As Models Only
Simple data structures with validation and serialization
- **[Go Models Example](doc/golang/models-example.md)** - Creating and working with generated Go structs
- **[Java Models Example](doc/java/models-example.md)** - Using generated Java classes with builder pattern
- **[Python Models Example](doc/python/models-example.md)** - Working with generated Python dataclasses

### As Server Implementation
Service implementations with HTTP endpoints

- **[Go Server Example](doc/golang/server-example.md)** - HTTP server implementation with generated service interface
- **[Java Server Example](doc/java/server-example.md)** - Java HTTP server using generated service classes
- **[Python Server Example](doc/python/server-example.md)** - Flask-based server with generated service interface

### As Client with Custom Transport
Client libraries with pluggable transport

- **[Go Client Example](doc/golang/client-example.md)** - HTTP client with custom transport implementation
- **[Java Client Example](doc/java/client-example.md)** - Java HTTP client with transport abstraction
- **[Python Client Example](doc/python/client-example.md)** - Python client with requests-based transport

You can define custom transports for different protocols (HTTP, gRPC, etc.) by implementing the `Transport` interface in each language.
Example: [Name-Based Routing Transport](examples/transport/name_based_routing_transport/README.md)


## Generated Code Features

### Method Metadata Support

The generator supports extracting metadata from service method comments. This is useful for HTTP routing, authorization, and other transport-specific configurations.

#### Defining Metadata in Proto Files

Add metadata to method comments using the `metadata:` prefix followed by a JSON object:

```protobuf
service UserService {
    // metadata:{"method":"GET", "path":"/users/{id}", "auth":"required"}
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    
    // metadata:{"method":"POST", "path":"/users", "auth":"required", "role":"admin"}
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
    
    // metadata:{"method":"DELETE", "path":"/users/{id}", "auth":"required", "role":"admin"}
    rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
}
```

#### Accessing Metadata in Generated Code

**Go:**
```go
// Access metadata using method constants
metadata := MethodMetadata[UserService_GetUser]
httpMethod := metadata["method"]  // "GET"
path := metadata["path"]          // "/users/{id}"
auth := metadata["auth"]          // "required"
```

**Java:**
```java
// Access metadata through the Methods class
Map<String, String> metadata = UserServiceMethods.METHOD_METADATA.get(UserServiceMethods.UserService_GetUser);
String httpMethod = metadata.get("method");  // "GET"
String path = metadata.get("path");          // "/users/{id}"
String auth = metadata.get("auth");          // "required"
```

**Python:**
```python
# Access metadata through the Methods class
metadata = UserServiceMethods.METHOD_METADATA[UserServiceMethods.UserService_GetUser]
http_method = metadata["method"]  # "GET"
path = metadata["path"]           # "/users/{id}"
auth = metadata["auth"]           # "required"
```

#### Example Use Cases

1. **HTTP Routing:** Use `method` and `path` metadata for automatic route registration
2. **Authentication:** Use `auth` metadata to determine if authentication is required
3. **Authorization:** Use `role` metadata for role-based access control
4. **Rate Limiting:** Add custom metadata for rate limiting configurations
5. **Documentation:** Include API versioning or documentation URLs

### Go

- Struct definitions with JSON tags
- Constructor functions (`NewMessageName()`)
- Validation methods
- JSON serialization (`ToJSON()`, `FromJSON()`)
- Service interfaces with default implementations
- Clients with pluggable Transport interface

### Java

- POJO classes with Jackson annotations
- Builder pattern support
- Getters and setters
- JSON serialization methods
- Service interfaces with default implementations
- Clients with generic Transport interface

### Python

- Dataclasses with type hints
- JSON serialization support
- Validation methods
- Service abstract base classes
- Clients with abstract Transport base class

## Testing the Plugin

Test with the provided example:

```bash
# Generate code for the example proto file
cd examples
protoc --puregen_out=./generated --puregen_opt=language=all proto/user.proto

# Check the generated files
ls -la generated/
```

## Development

### Project Structure

```
├── cmd/protoc-gen-puregen/   # Main plugin entry point
├── internal/generator/         # Code generation logic
│   ├── go.go                  # Go code generator
│   ├── java.go                # Java code generator
│   └── python.go              # Python code generator
├── examples/                   # Example proto files and usage
└── README.md
```

### Adding New Language Support

1. Create a new generator file in `internal/generator/`
2. Implement the `GenerateXXXFile` function
3. Add the language option to the main plugin

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

