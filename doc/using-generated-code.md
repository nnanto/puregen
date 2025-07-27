# Using the Generated Code

The generated code can be used in three different ways:

1. **As Models Only** - Simple data structures with validation and serialization
2. **As Server Implementation** - Service implementations with HTTP endpoints
3. **As Client with Custom Transport** - Client libraries with pluggable transport

## Method Metadata Support

The generator supports extracting metadata from service method comments. This is useful for HTTP routing, authorization, and other transport-specific configurations.

### Defining Metadata in Proto Files

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

### Accessing Metadata in Generated Code

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

### Example Use Cases

1. **HTTP Routing:** Use `method` and `path` metadata for automatic route registration
2. **Authentication:** Use `auth` metadata to determine if authentication is required
3. **Authorization:** Use `role` metadata for role-based access control
4. **Rate Limiting:** Add custom metadata for rate limiting configurations
5. **Documentation:** Include API versioning or documentation URLs

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
