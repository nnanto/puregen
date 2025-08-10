# Metadata Support in puregen

puregen now supports the `puregen:metadata` directive on messages, enums, and fields, in addition to the existing support for service methods.

## Overview

The `puregen:metadata` directive allows you to attach custom metadata to various protobuf elements:

- **Service Methods** (existing): Metadata is stored in `<ServiceName>Methods.METHOD_METADATA`
- **Messages** (new): Metadata is stored in `<MessageName>Metadata`
- **Enums** (new): Metadata is stored in `<EnumName>Metadata`
- **Fields** (new): Metadata is stored in `<MessageName><FieldName>Metadata`

## Usage

### Message Metadata

```protobuf
// Message with metadata for database table mapping
// puregen:metadata: {"table": "users", "cache": "true", "partition_key": "tenant_id"}
message User {
    string id = 1;
    string name = 2;
}
```

Generated code includes:
- **Python**: `UserMetadata = {"table": "users", "cache": "true", "partition_key": "tenant_id"}`
- **Go**: `var UserMetadata = map[string]string{"table": "users", "cache": "true", "partition_key": "tenant_id"}`
- **Java**: `UserMetadata.METADATA` contains the key-value pairs

### Enum Metadata

```protobuf
// Enum with validation and UI metadata
// puregen:metadata: {"validation": "required", "ui_type": "dropdown", "default": "PENDING"}
enum Status {
    UNKNOWN = 0;
    PENDING = 1;
    APPROVED = 2;
    REJECTED = 3;
}
```

Generated code includes:
- **Python**: `StatusMetadata = {"validation": "required", "ui_type": "dropdown", "default": "PENDING"}`
- **Go**: `var StatusMetadata = map[string]string{"validation": "required", "ui_type": "dropdown", "default": "PENDING"}`
- **Java**: `StatusMetadata.METADATA` contains the key-value pairs

### Field Metadata

```protobuf
message User {
    // Field with database and validation metadata
    // puregen:metadata: {"db_column": "user_id", "index": "primary", "validation": "uuid"}
    string id = 1;
    
    // Field with validation constraints
    // puregen:metadata: {"validation": "required", "min_length": "2", "max_length": "100"}
    string name = 2;
    
    // Field with UI metadata
    // puregen:metadata: {"ui_widget": "email_input", "placeholder": "Enter email"}
    string email = 3;
}
```

Generated code includes:
- **Python**: `UserIdMetadata`, `UserNameMetadata`, `UserEmailMetadata`
- **Go**: `var UserIdMetadata`, `var UserNameMetadata`, `var UserEmailMetadata`
- **Java**: `UserIdMetadata.METADATA`, `UserNameMetadata.METADATA`, `UserEmailMetadata.METADATA`

## Common Use Cases

### Database Mapping
```protobuf
// puregen:metadata: {"table": "orders", "schema": "commerce"}
message Order {
    // puregen:metadata: {"column": "order_id", "type": "uuid", "primary_key": "true"}
    string id = 1;
    
    // puregen:metadata: {"column": "customer_id", "foreign_key": "customers.id"}
    string customer_id = 2;
}
```

### Validation Rules
```protobuf
message CreateUserRequest {
    // puregen:metadata: {"required": "true", "pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"}
    string email = 1;
    
    // puregen:metadata: {"required": "true", "min_length": "8", "pattern": "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)"}
    string password = 2;
}
```

### UI Configuration
```protobuf
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
```protobuf
// puregen:metadata: {"endpoint": "/api/v1/users", "rate_limit": "100/hour"}
message User {
    // puregen:metadata: {"searchable": "true", "filterable": "true"}
    string name = 1;
    
    // puregen:metadata: {"private": "true", "admin_only": "true"}
    string internal_notes = 2;
}
```

## Language-Specific Access

### Python
```python
# Access message metadata
user_table = UserMetadata.get("table")  # "users"

# Access field metadata  
id_validation = UserIdMetadata.get("validation")  # "uuid"

# Access enum metadata
status_default = StatusMetadata.get("default")  # "PENDING"
```

### Go
```go
// Access message metadata
userTable := UserMetadata["table"]  // "users"

// Access field metadata
idValidation := UserIdMetadata["validation"]  // "uuid"

// Access enum metadata
statusDefault := StatusMetadata["default"]  // "PENDING"
```

### Java
```java
// Access message metadata
String userTable = UserMetadata.METADATA.get("table");  // "users"

// Access field metadata
String idValidation = UserIdMetadata.METADATA.get("validation");  // "uuid"

// Access enum metadata
String statusDefault = StatusMetadata.METADATA.get("default");  // "PENDING"
```

## Format Requirements

- Metadata must be valid JSON
- Keys and values must be strings
- Use proper JSON escaping for special characters
- Metadata must be placed in comments immediately before the element

### Valid Examples
```protobuf
// puregen:metadata: {"key": "value"}
// puregen:metadata:{"key":"value","key2":"value2"}
// puregen:metadata: {"escaped": "value with \"quotes\""}
```

### Invalid Examples
```protobuf
// puregen:metadata: {key: "value"}  // Keys must be quoted
// puregen:metadata: {"key": value}  // Values must be quoted
// puregen:metadata: {"key": 123}    // Only string values supported
```
