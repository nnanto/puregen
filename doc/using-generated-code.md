# Using the Generated Code

The generated code can be used in three different ways:

1. **As Models Only** - Simple data structures with validation and serialization
2. **As Server Implementation** - Service implementations with HTTP endpoints
3. **As Client with Custom Transport** - Client libraries with pluggable transport

## Code Generation Options

The generator supports several options:

- `language` - Target language (go, java, python, or all)
- `common_namespace` - Namespace for common/shared classes like Transport interface (e.g., 'shared', 'common.transport')

When `common_namespace` is specified, Transport interfaces/classes are generated in a global namespace and imported by clients. When not specified, each client generates its own local Transport interface.

#### Usage
# Common namespace for all languages
protoc --plugin=./build/protoc-gen-puregen \
       --puregen_out=./output \
       --puregen_opt=language=all,common_namespace=shared \
       examples/proto/user.proto

# Multi-level common namespace
protoc --plugin=./build/protoc-gen-puregen \
       --puregen_out=./output \
       --puregen_opt=language=python,common_namespace=common.transport \
       examples/proto/user.proto

# Without common namespace (local transport interfaces)
protoc --plugin=./build/protoc-gen-puregen \
       --puregen_out=./output \
       --puregen_opt=language=all \
       examples/proto/user.proto

## Language-Specific Examples

### As Models Only

- **[Go Models Example](golang/models-example.md)** - Creating and working with generated Go structs
- **[Java Models Example](java/models-example.md)** - Using generated Java classes with builder pattern  
- **[Python Models Example](python/models-example.md)** - Working with generated Python dataclasses

### As Server Implementation

- **[Go Server Example](golang/server-example.md)** - HTTP server implementation with generated service interface
- **[Java Server Example](java/server-example.md)** - Java HTTP server using generated service classes
- **[Python Server Example](python/server-example.md)** - Flask-based server with generated service interface

### As Client with Custom Transport

- **[Go Client Example](golang/client-example.md)** - HTTP client with custom transport implementation
- **[Java Client Example](java/client-example.md)** - Java HTTP client with transport abstraction
- **[Python Client Example](python/client-example.md)** - Python client with requests-based transport
