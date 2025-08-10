# Puregen Directives

Puregen supports several directives that allow you to customize code generation behavior. Directives are specified in protobuf comments using JSON syntax.

> **Migration Note**: This document consolidates and supersedes the previous separate documentation files (`value-directive.md` and `metadata-support.md`). All directive documentation is now centralized here for easier maintenance and reference.

## Table of Contents

1. [Directive Syntax](#directive-syntax)
2. [Available Directives](#available-directives)
   - [`puregen:generate` - Code Generation Control](#1-puregengenerate---code-generation-control)
   - [`puregen:metadata` - Metadata Attachment](#2-puregenmetadata---metadata-attachment)
3. [Use Cases](#use-cases)
4. [Best Practices](#best-practices)
5. [Syntax Notes](#syntax-notes)
6. [Error Handling](#error-handling)

## Directive Syntax

```proto
// puregen:directive_name: {"key": "value", "key2": "value2"}
```

Directives can be applied to:
- **Service Methods**: For endpoint routing and behavior
- **Messages**: For database mapping, caching, and other configurations
- **Enums**: For generation type and validation rules
- **Fields**: For default values, validation, and UI configuration

## Available Directives

### 1. `puregen:generate` - Code Generation Control

Controls how code is generated for specific elements.

#### Enum Generation Type

Controls whether enums are generated as string constants (default) or integer enums.

```proto
// Generate as integer enum instead of string constants
// puregen:generate: {"enumType": "int"}
enum Status {
  STATUS_UNKNOWN = 0;
  STATUS_ACTIVE = 1;
  STATUS_INACTIVE = 2;
}

// Default behavior - generated as string constants
enum Priority {
  PRIORITY_LOW = 0;
  PRIORITY_MEDIUM = 1;
  PRIORITY_HIGH = 2;
}
```

**String Constants (Default):**
- **Go**: `const Priority_PRIORITY_LOW = "PRIORITY_LOW"`
- **Java**: `public static final String PRIORITY_LOW = "PRIORITY_LOW"`
- **Python**: `PRIORITY_LOW = "PRIORITY_LOW"`

**Integer Enums (with directive):**
- **Go**: `type Status int32; const Status_STATUS_UNKNOWN Status = 0`
- **Java**: `public enum Status { STATUS_UNKNOWN(0) }`
- **Python**: `class Status(IntEnum): STATUS_UNKNOWN = 0`

#### Field Default Values

Sets default values for primitive type fields when objects are created.

```proto
message UserProfile {
    // String field with default value
    // puregen:generate: {"value": "Anonymous"}
    string display_name = 1;
    
    // Integer field with default value
    // puregen:generate: {"value": "18"}
    int32 age = 2;
    
    // Boolean field with default value
    // puregen:generate: {"value": "true"}
    bool is_active = 3;
    
    // Float field with default value
    // puregen:generate: {"value": "0.0"}
    float score = 4;
}
```

**Supported Field Types:**
- **string**: Text values (enclosed in quotes in directive)
- **bool**: Boolean values (`"true"` or `"false"`)
- **int32, int64**: Integer values
- **uint32, uint64**: Unsigned integer values
- **sint32, sint64**: Signed integer values
- **fixed32, fixed64**: Fixed-size integer values
- **sfixed32, sfixed64**: Signed fixed-size integer values
- **float, double**: Floating-point values

**Generated Constructors:**

**Go:**
```go
func NewUserProfile() *UserProfile {
    return &UserProfile{
        DisplayName: "Anonymous",
        Age:         18,
        IsActive:    true,
        Score:       0.0,
    }
}
```

**Java:**
```java
public UserProfile() {
    this.displayName = "Anonymous";
    this.age = 18;
    this.isActive = true;
    this.score = 0.0f;
}
```

**Python:**
```python
@dataclass
class UserProfile:
    display_name: str = "Anonymous"
    age: int = 18
    is_active: bool = True
    score: float = 0.0
```

**Important Notes:**
- Only works with primitive types (not message types, enums, or repeated fields)
- Values are only applied when using generated constructors
- Invalid values fall back to language defaults

### 2. `puregen:metadata` - Metadata Attachment

Attaches custom metadata to protobuf elements for use in generated code.

#### Service Method Metadata

```proto
service UserService {
  // HTTP routing metadata
  // puregen:metadata: {"method":"POST", "path":"/users", "auth": "required"}
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);

  // Caching and rate limiting metadata
  // puregen:metadata: {"method":"GET", "path":"/users/{id}", "cache": "true", "cache_ttl": "300"}
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}
```

#### Message Metadata

```proto
// Database table mapping
// puregen:metadata: {"table": "users", "cache": "true", "partition_key": "tenant_id"}
message User {
    string id = 1;
    string name = 2;
}

// API configuration
// puregen:metadata: {"endpoint": "/api/v1/orders", "rate_limit": "100/hour"}
message Order {
    string order_id = 1;
    string customer_id = 2;
}
```

#### Enum Metadata

```proto
// Validation and UI metadata
// puregen:metadata: {"validation": "required", "ui_type": "dropdown", "default": "PENDING"}
enum TaskStatus {
    PENDING = 0;
    IN_PROGRESS = 1;
    COMPLETED = 2;
}
```

#### Field Metadata

```proto
message Task {
    // Database column mapping
    // puregen:metadata: {"db_column": "task_id", "index": "primary", "validation": "uuid"}
    string id = 1;
    
    // Validation rules
    // puregen:metadata: {"validation": "required", "min_length": "2", "max_length": "100"}
    string title = 2;
    
    // UI configuration
    // puregen:metadata: {"ui_widget": "email_input", "placeholder": "Enter email"}
    string email = 3;
}
```

**Generated Metadata Access:**

**Python:**
```python
# Message metadata
table_name = UserMetadata["table"]  # "users"

# Field metadata
id_column = TaskFieldMetadata[Task_Id_FIELD]["db_column"]  # "task_id"

# Method metadata
endpoint = UserServiceMethods[UserService_CreateUser_METHOD]["path"]  # "/users"
```

**Go:**
```go
// Message metadata
tableName := UserMetadata["table"]  // "users"

// Field metadata
idColumn := TaskFieldMetadata[Task_Id_FIELD]["db_column"]  // "task_id"

// Method metadata
endpoint := UserServiceMethods[UserService_CreateUser_METHOD]["path"]  // "/users"
```

**Java:**
```java
// Message metadata
String tableName = UserMetadata.METADATA.get("table");  // "users"

// Field metadata
String idColumn = TaskFieldMetadata.METADATA.get(Task.Id_FIELD).get("db_column");  // "task_id"

// Method metadata
String endpoint = UserServiceMethods.METADATA.get(UserServiceMethods.CreateUser_METHOD).get("path");  // "/users"
```

## Use Cases

### Database Mapping
```proto
// puregen:metadata: {"table": "orders", "schema": "commerce"}
message Order {
    // puregen:metadata: {"column": "order_id", "type": "uuid", "primary_key": "true"}
    string id = 1;
    
    // puregen:metadata: {"column": "customer_id", "foreign_key": "customers.id"}
    string customer_id = 2;
}
```

### Validation Rules
```proto
message UserRegistration {
    // puregen:metadata: {"required": "true", "pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"}
    string email = 1;
    
    // puregen:metadata: {"required": "true", "min_length": "8", "pattern": "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)"}
    string password = 2;
}
```

### UI Configuration
```proto
// puregen:metadata: {"form_title": "User Profile", "section": "account"}
message UserProfile {
    // puregen:metadata: {"widget": "text", "label": "Full Name", "required": "true"}
    string name = 1;
    
    // puregen:metadata: {"widget": "date", "label": "Date of Birth", "max_date": "today"}
    string birth_date = 2;
    
    // puregen:metadata: {"widget": "select", "options": "countries", "label": "Country"}
    string country = 3;
}
```

### API Configuration
```proto
// puregen:metadata: {"endpoint": "/api/v1/users", "rate_limit": "100/hour"}
message UserService {
    // puregen:metadata: {"searchable": "true", "filterable": "true"}
    string name = 1;
    
    // puregen:metadata: {"private": "true", "admin_only": "true"}
    string internal_notes = 2;
}
```

## Best Practices

1. **Use descriptive keys**: Choose meaningful metadata keys that clearly indicate their purpose
2. **Consistent naming**: Use consistent naming conventions across your project
3. **Validate metadata**: Ensure metadata values are valid for your use case
4. **Document custom metadata**: Document any custom metadata keys used in your project
5. **Keep it simple**: Avoid overly complex nested structures in metadata

## Syntax Notes

- JSON syntax is required: `{"key": "value"}`
- String values must be quoted: `{"name": "value"}`
- Multiple key-value pairs: `{"key1": "value1", "key2": "value2"}`
- Whitespace is flexible: `{"key":"value"}` or `{"key": "value"}` both work
- Comments can have multiple directive lines if needed

## Error Handling

- Invalid JSON syntax will be ignored
- Unsupported directive names will be ignored
- Invalid values for supported directives fall back to defaults
- Multiple directives on the same element will be merged (later ones override earlier ones for same keys)
