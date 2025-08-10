# Metadata Support in puregen

puregen now supports the `puregen:metadata` directive on messages, enums, and fields, in addi**Generated Access:**
```python
# Python
table_name = OrderMetadata["table"]  # "orders"
id_column = OrderFieldMetadata[Order_Id_FIELD]["column"]  # "order_id"
is_primary = OrderFieldMetadata[Order_Id_FIELD]["primary_key"]  # "true"
```

```go
// Go
tableName := OrderMetadata["table"]  // "orders"
idColumn := OrderFieldMetadata[Order_Id_FIELD]["column"]  // "order_id"  
isPrimary := OrderFieldMetadata[Order_Id_FIELD]["primary_key"]  // "true"
```isting support for service methods.

## Overview

The `puregen:metadata` directive allows you to attach custom metadata to various protobuf elements:

- **Service Methods** (existing): Metadata is stored in `MethodMetadata` map with method constants as keys
- **Messages** (new): Metadata is stored in `<MessageName>Metadata`
- **Enums** (new): Metadata is stored in `<EnumName>Metadata`
- **Fields** (new): Metadata is stored in `<MessageName>FieldMetadata` map with field constants as keys

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

**Field Constants and Metadata Map:**
- **Python**: 
  ```python
  User_Id_FIELD = "User_Id"
  User_Name_FIELD = "User_Name" 
  User_Email_FIELD = "User_Email"
  
  UserFieldMetadata = {
      User_Id_FIELD: {"db_column": "user_id", "index": "primary", "validation": "uuid"},
      User_Name_FIELD: {"validation": "required", "min_length": "2", "max_length": "100"},
      User_Email_FIELD: {"ui_widget": "email_input", "placeholder": "Enter email"}
  }
  ```

- **Go**: 
  ```go
  const (
      User_Id_FIELD    = "User_Id"
      User_Name_FIELD  = "User_Name"
      User_Email_FIELD = "User_Email"
  )
  
  var UserFieldMetadata = map[string]map[string]string{
      User_Id_FIELD: {"db_column": "user_id", "index": "primary", "validation": "uuid"},
      User_Name_FIELD: {"validation": "required", "min_length": "2", "max_length": "100"},
      User_Email_FIELD: {"ui_widget": "email_input", "placeholder": "Enter email"},
  }
  ```

- **Java**: 
  ```java
  public final class UserFieldMetadata {
      public static final String User_Id_FIELD = "User_Id";
      public static final String User_Name_FIELD = "User_Name";
      public static final String User_Email_FIELD = "User_Email";
      
      public static final Map<String, Map<String, String>> FIELD_METADATA = new HashMap<>();
      // ... populated with field metadata
  }
  ```

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

**Generated Access:**
```python
# Python
table_name = OrderMetadata["table"]  # "orders"
id_column = OrderFieldMetadata[Order_Id]["column"]  # "order_id"
is_primary = OrderFieldMetadata[Order_Id]["primary_key"]  # "true"
```

```go
// Go
tableName := OrderMetadata["table"]  // "orders"
idColumn := OrderFieldMetadata[Order_Id]["column"]  // "order_id"  
isPrimary := OrderFieldMetadata[Order_Id]["primary_key"]  // "true"
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

**Generated Access:**
```python
# Python - Validation example
def validate_user_request(request):
    for field_const, metadata in CreateUserRequestFieldMetadata.items():
        if metadata.get("required") == "true":
            # Check if field is present and validate
            pass
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

**Generated Access:**
```python
# Python - UI form generation example
def generate_form():
    form_title = UserProfileMetadata["form_title"]  # "User Profile"
    
    for field_const, metadata in UserProfileFieldMetadata.items():
        widget_type = metadata.get("widget", "text")
        label = metadata.get("label", field_const)
        required = metadata.get("required") == "true"
        # Generate UI element based on metadata
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

**Generated Access:**
```go
// Go - API configuration example
func configureAPI() {
    endpoint := UserMetadata["endpoint"]  // "/api/v1/users"
    rateLimit := UserMetadata["rate_limit"]  // "100/hour"
    
    // Configure searchable fields
    for fieldConst, metadata := range UserFieldMetadata {
        if metadata["searchable"] == "true" {
            enableSearchOnField(fieldConst)
        }
        if metadata["private"] == "true" {
            restrictFieldAccess(fieldConst)
        }
    }
}
```

## Language-Specific Access

### Field Metadata Structure

Field metadata follows the same pattern as service method metadata for consistency:

**Key Benefits:**
- **Centralized**: All field metadata for a message is in one unified map
- **Type Safety**: Field constants prevent typos in field names
- **Consistent**: Same pattern as service method metadata (`MethodMetadata`)
- **Iterable**: Easy to iterate over all field metadata for a message
- **Maintainable**: Single source of truth for field metadata per message

### Python
```python
# Access message metadata
user_table = UserMetadata.get("table")  # "users"

# Access field metadata using constants
id_validation = UserFieldMetadata[User_Id_FIELD].get("validation")  # "uuid"
name_min_length = UserFieldMetadata[User_Name_FIELD].get("min_length")  # "2"

# Access enum metadata
status_default = StatusMetadata.get("default")  # "PENDING"

# Iterate through all field metadata
for field_const, metadata in UserFieldMetadata.items():
    print(f"Field {field_const}: {metadata}")
```

### Go
```go
// Access message metadata
userTable := UserMetadata["table"]  // "users"

// Access field metadata using constants
idValidation := UserFieldMetadata[User_Id_FIELD]["validation"]  // "uuid"
nameMinLength := UserFieldMetadata[User_Name_FIELD]["min_length"]  // "2"

// Access enum metadata
statusDefault := StatusMetadata["default"]  // "PENDING"

// Iterate through all field metadata
for fieldConst, metadata := range UserFieldMetadata {
    fmt.Printf("Field %s: %v\n", fieldConst, metadata)
}
```

### Java
```java
// Access message metadata
String userTable = UserMetadata.METADATA.get("table");  // "users"

// Access field metadata using constants
String idValidation = UserFieldMetadata.FIELD_METADATA.get(UserFieldMetadata.User_Id_FIELD).get("validation");  // "uuid"
String nameMinLength = UserFieldMetadata.FIELD_METADATA.get(UserFieldMetadata.User_Name_FIELD).get("min_length");  // "2"

// Access enum metadata
String statusDefault = StatusMetadata.METADATA.get("default");  // "PENDING"

// Iterate through all field metadata
for (Map.Entry<String, Map<String, String>> entry : UserFieldMetadata.FIELD_METADATA.entrySet()) {
    System.out.println("Field " + entry.getKey() + ": " + entry.getValue());
}
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
