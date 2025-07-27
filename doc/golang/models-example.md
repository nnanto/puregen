# Go Models Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    
    "github.com/nnanto/puregen/examples/proto/gen/go"
)

func main() {
    // Create a new user
    user := proto.NewUser()
    user.Id = 1
    user.Name = "John Doe"
    user.Email = "john@example.com"
    user.IsActive = true
    
    // Create profile
    profile := proto.NewUserProfile()
    profile.Bio = "Software Engineer"
    profile.AvatarUrl = "https://example.com/avatar.jpg"
    profile.CreatedAt = 1640995200 // Unix timestamp
    user.Profile = profile
    
    // Validate
    if err := user.Validate(); err != nil {
        log.Fatal("Validation failed:", err)
    }
    
    // Convert to JSON
    jsonData, err := user.ToJSON()
    if err != nil {
        log.Fatal("JSON marshaling failed:", err)
    }
    
    fmt.Printf("User JSON: %s\n", string(jsonData))
    
    // Create from JSON
    newUser := proto.NewUser()
    if err := newUser.FromJSON(jsonData); err != nil {
        log.Fatal("JSON unmarshaling failed:", err)
    }
    
    fmt.Printf("Restored user: %+v\n", newUser)
}
```
