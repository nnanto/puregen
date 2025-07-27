# Go Client Example

```go
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    
    "github.com/nnanto/puregen/examples/proto/gen/go"
)

// HTTP Transport implementation
type HTTPTransport struct {
    baseURL string
    client  *http.Client
}

func NewHTTPTransport(baseURL string) *HTTPTransport {
    return &HTTPTransport{
        baseURL: baseURL,
        client:  &http.Client{},
    }
}

func (t *HTTPTransport) Send(ctx context.Context, methodName string, inputData interface{}, outputType interface{}) (interface{}, error) {
    // Determine endpoint based on method name
    var endpoint string
    var method string
    
    switch methodName {
    case proto.UserService_CreateUser:
        endpoint = "/users"
        method = "POST"
    case proto.UserService_GetUser:
        // For GET requests, we need to extract the ID from the input
        if req, ok := inputData.(*proto.GetUserRequest); ok {
            endpoint = fmt.Sprintf("/users/%d", req.Id)
            method = "GET"
        }
    default:
        return nil, fmt.Errorf("unknown method: %s", methodName)
    }
    
    url := t.baseURL + endpoint
    
    var body io.Reader
    if method == "POST" {
        jsonData, err := json.Marshal(inputData)
        if err != nil {
            return nil, err
        }
        body = bytes.NewReader(jsonData)
    }
    
    req, err := http.NewRequestWithContext(ctx, method, url, body)
    if err != nil {
        return nil, err
    }
    
    if method == "POST" {
        req.Header.Set("Content-Type", "application/json")
    }
    
    resp, err := t.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
    }
    
    respData, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    // Create the output type based on the nil pointer type
    switch outputType.(type) {
    case (*proto.CreateUserResponse):
        result := &proto.CreateUserResponse{}
        err = json.Unmarshal(respData, result)
        return result, err
    case (*proto.GetUserResponse):
        result := &proto.GetUserResponse{}
        err = json.Unmarshal(respData, result)
        return result, err
    default:
        return nil, fmt.Errorf("unknown output type")
    }
}

func main() {
    // Create HTTP transport
    transport := NewHTTPTransport("http://localhost:8080")
    
    // Create client
    client := proto.NewUserServiceClient(transport)
    
    ctx := context.Background()
    
    // Create a user
    profile := &proto.UserProfile{
        Bio: "Software Engineer",
        AvatarUrl: "https://example.com/avatar.jpg",
        CreatedAt: 1640995200,
    }
    
    createReq := &proto.CreateUserRequest{
        Name: "Jane Doe",
        Email: "jane@example.com",
        Profile: profile,
    }
    
    createResp, err := client.CreateUser(ctx, createReq)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created user: %+v\n", createResp.User)
    
    // Get the user
    getReq := &proto.GetUserRequest{
        Id: createResp.User.Id,
    }
    
    getResp, err := client.GetUser(ctx, getReq)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Retrieved user: %+v\n", getResp.User)
}
```
