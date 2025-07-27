# Name-Based Routing Transport

This transport implementation automatically converts proto RPC method calls to REST API endpoints based on naming patterns. It uses method and service names to determine the appropriate HTTP method and URL structure.

## Concept

The transport analyzes the RPC method name in the format `ServiceName_MethodName` and converts it to a REST API call by:

1. **Service Name Extraction**: Extracts the service name and converts it to lowercase
2. **Method Name Analysis**: Determines HTTP method based on the RPC method name prefix
3. **URL Construction**: Builds REST endpoints using the service and method names
4. **Parameter Handling**: Routes parameters to query params (GET) or request body (POST/PUT/DELETE)

## Naming Pattern to HTTP Method Mapping

| RPC Method Prefix | HTTP Method | Description |
|------------------|-------------|-------------|
| `Get*` | GET | Retrieve operations |
| `Update*` | PUT | Update operations |
| `Delete*` | DELETE | Delete operations |
| `Create*`, `Start*`, `Describe*`, etc. | POST | All other operations |

## Examples

### Example 1: User Service

**Proto RPC Method**: `UserService_GetUser`

```
Input: { id: "123" }
↓
HTTP GET /user/getuser?id=123
```

**Proto RPC Method**: `UserService_CreateUser`

```
Input: { name: "John", email: "john@example.com" }
↓
HTTP POST /user/createuser
Body: {"name": "John", "email": "john@example.com"}
```

**Proto RPC Method**: `UserService_UpdateUser`

```
Input: { id: "123", name: "Jane" }
↓
HTTP PUT /user/updateuser
Body: {"id": "123", "name": "Jane"}
```

**Proto RPC Method**: `UserService_DeleteUser`

```
Input: { id: "123" }
↓
HTTP DELETE /user/deleteuser
Body: {"id": "123"}
```

### Example 2: Order Service

**Proto RPC Method**: `OrderService_GetOrdersByStatus`

```
Input: { status: "pending", limit: 10 }
↓
HTTP GET /order/getordersbystatus?status=pending&limit=10
```

**Proto RPC Method**: `OrderService_CreateOrder`

```
Input: { 
  customerId: "456", 
  items: [{"productId": "789", "quantity": 2}] 
}
↓
HTTP POST /order/createorder
Body: {
  "customerId": "456",
  "items": [{"productId": "789", "quantity": 2}]
}
```

### Example 3: Workflow Service

**Proto RPC Method**: `WorkflowService_StartWorkflow`

```
Input: { workflowId: "wf-001", parameters: {...} }
↓
HTTP POST /workflow/startworkflow
Body: {"workflowId": "wf-001", "parameters": {...}}
```

**Proto RPC Method**: `WorkflowService_DescribeWorkflow`

```
Input: { workflowId: "wf-001" }
↓
HTTP POST /workflow/describeworkflow
Body: {"workflowId": "wf-001"}
```

## Key Features

1. **Automatic HTTP Method Detection**: Determines appropriate HTTP methods based on RPC method name prefixes
2. **Service Name Normalization**: Converts service names to lowercase and removes "Service" suffix
3. **Smart Parameter Routing**: 
   - GET requests: Parameters converted to query strings
   - Other methods: Parameters sent in JSON request body
4. **Flexible Input Handling**: Supports structs with JSON tags for proper field mapping
5. **Type-Safe Output**: Creates properly typed response objects

## Usage

```go
// Create transport
transport := NewHTTPTransport("http://localhost:8080")

// Example call - this would convert to HTTP GET /user/getuser?id=123
result, err := transport.Send(
    ctx, 
    "UserService_GetUser", 
    &GetUserRequest{ID: "123"}, 
    (*GetUserResponse)(nil),
)
```
