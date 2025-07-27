# Go Server Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "strconv"
    
    "github.com/nnanto/puregen/examples/proto/gen/go"
)

// Implement the service interface
type UserServiceImpl struct {
    users map[int32]*proto.User
    nextID int32
}

func NewUserService() *UserServiceImpl {
    return &UserServiceImpl{
        users: make(map[int32]*proto.User),
        nextID: 1,
    }
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
    user := proto.NewUser()
    user.Id = s.nextID
    user.Name = req.Name
    user.Email = req.Email
    user.IsActive = true
    user.Profile = req.Profile
    
    s.users[user.Id] = user
    s.nextID++
    
    return &proto.CreateUserResponse{
        User: user,
        Message: "User created successfully",
    }, nil
}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
    user, exists := s.users[req.Id]
    if !exists {
        return nil, fmt.Errorf("user not found")
    }
    
    return &proto.GetUserResponse{
        User: user,
    }, nil
}

// HTTP handler for JSON-based API
func (s *UserServiceImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    switch r.URL.Path {
    case "/users":
        if r.Method == "POST" {
            s.handleCreateUser(w, r)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    default:
        if r.Method == "GET" && len(r.URL.Path) > 7 { // "/users/"
            s.handleGetUser(w, r)
        } else {
            http.NotFound(w, r)
        }
    }
}

func (s *UserServiceImpl) handleCreateUser(w http.ResponseWriter, r *http.Request) {
    var req proto.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    resp, err := s.CreateUser(r.Context(), &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(resp)
}

func (s *UserServiceImpl) handleGetUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Path[7:] // Remove "/users/"
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    req := &proto.GetUserRequest{Id: int32(id)}
    resp, err := s.GetUser(r.Context(), req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    json.NewEncoder(w).Encode(resp)
}

func main() {
    service := NewUserService()
    
    fmt.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", service))
}
```
