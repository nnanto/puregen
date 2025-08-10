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

## Enum Generation

Puregen supports two different enum generation styles that can be controlled using puregen directives in your proto file comments:

1. **String Constants (Default)** - Enums are generated as string constants with validation functions
2. **Integer Enums** - Traditional integer-based enums with type safety

### Enum Directives

Use the `puregen:generate:` directive in enum comments to control generation:

```proto
// puregen:generate:{"enumType": "int"}
// This enum will be generated as integer constants
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

message ExampleMessage {
  Status status = 1;      // Field type: Status (Go), int (Python), Status (Java)
  Priority priority = 2;  // Field type: string (Go), str (Python), String (Java)
}
```

### Generated Output Examples

#### String Constants (Default)

**Go:**
```go
// Priority enum values as string constants
const (
    Priority_PRIORITY_LOW    = "PRIORITY_LOW"
    Priority_PRIORITY_MEDIUM = "PRIORITY_MEDIUM"
    Priority_PRIORITY_HIGH   = "PRIORITY_HIGH"
)

var PriorityValues = []string{
    Priority_PRIORITY_LOW,
    Priority_PRIORITY_MEDIUM,
    Priority_PRIORITY_HIGH,
}

func IsValidPriority(value string) bool {
    for _, v := range PriorityValues {
        if v == value {
            return true
        }
    }
    return false
}

type ExampleMessage struct {
    Priority string `json:"priority"`
}
```

**Python:**
```python
class Priority:
    """Priority enum values as string constants"""
    PRIORITY_LOW = "PRIORITY_LOW"
    PRIORITY_MEDIUM = "PRIORITY_MEDIUM"
    PRIORITY_HIGH = "PRIORITY_HIGH"

    VALUES = [PRIORITY_LOW, PRIORITY_MEDIUM, PRIORITY_HIGH]

    @classmethod
    def is_valid(cls, value: str) -> bool:
        return value in cls.VALUES

@dataclass
class ExampleMessage:
    priority: str = ""
```

**Java:**
```java
public final class Priority {
    private Priority() {} // Prevent instantiation
    
    public static final String PRIORITY_LOW = "PRIORITY_LOW";
    public static final String PRIORITY_MEDIUM = "PRIORITY_MEDIUM";
    public static final String PRIORITY_HIGH = "PRIORITY_HIGH";
    
    public static final String[] VALUES = {
        PRIORITY_LOW, PRIORITY_MEDIUM, PRIORITY_HIGH
    };
    
    public static boolean isValid(String value) {
        for (String v : VALUES) {
            if (v.equals(value)) return true;
        }
        return false;
    }
}

public class ExampleMessage {
    private String priority;
}
```

#### Integer Enums (with puregen:generate directive)

**Go:**
```go
type Status int32

const (
    Status_STATUS_UNKNOWN  Status = 0
    Status_STATUS_ACTIVE          = 1
    Status_STATUS_INACTIVE        = 2
)

func (x Status) String() string {
    if name, ok := Status_name[int32(x)]; ok {
        return name
    }
    return fmt.Sprintf("Status(%d)", x)
}

func ParseStatus(s string) (Status, error) {
    if value, ok := Status_value[s]; ok {
        return Status(value), nil
    }
    return 0, fmt.Errorf("invalid Status value: %s", s)
}

type ExampleMessage struct {
    Status Status `json:"status"`
}
```

**Python:**
```python
from enum import IntEnum

class Status(IntEnum):
    """Status enum values as integers"""
    STATUS_UNKNOWN = 0
    STATUS_ACTIVE = 1
    STATUS_INACTIVE = 2

    @classmethod
    def is_valid(cls, value: int) -> bool:
        return value in [item.value for item in cls]

@dataclass
class ExampleMessage:
    status: int = 0
```

**Java:**
```java
public enum Status {
    STATUS_UNKNOWN(0),
    STATUS_ACTIVE(1),
    STATUS_INACTIVE(2);
    
    private final int value;
    
    Status(int value) {
        this.value = value;
    }
    
    public int getValue() {
        return value;
    }
    
    public static Status fromValue(int value) {
        for (Status e : values()) {
            if (e.value == value) return e;
        }
        throw new IllegalArgumentException("Invalid Status value: " + value);
    }
}

public class ExampleMessage {
    private Status status;
}
```

### Usage Examples

**String Constants:**
```go
// Go
msg := &ExampleMessage{
    Priority: Priority_PRIORITY_HIGH,
}
if IsValidPriority(msg.Priority) {
    // Handle valid priority
}
```

```python
# Python
msg = ExampleMessage(priority=Priority.PRIORITY_HIGH)
if Priority.is_valid(msg.priority):
    # Handle valid priority
```

```java
// Java
ExampleMessage msg = new ExampleMessage();
msg.setPriority(Priority.PRIORITY_HIGH);
if (Priority.isValid(msg.getPriority())) {
    // Handle valid priority
}
```

**Integer Enums:**
```go
// Go
msg := &ExampleMessage{
    Status: Status_STATUS_ACTIVE,
}
if msg.Status.IsValid() {
    fmt.Println(msg.Status.String()) // Prints: "STATUS_ACTIVE"
}
```

```python
# Python
msg = ExampleMessage(status=Status.STATUS_ACTIVE)
if Status.is_valid(msg.status):
    print(Status(msg.status).name)  # Prints: "STATUS_ACTIVE"
```

```java
// Java
ExampleMessage msg = new ExampleMessage();
msg.setStatus(Status.STATUS_ACTIVE);
System.out.println(msg.getStatus().name()); // Prints: "STATUS_ACTIVE"
```

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
