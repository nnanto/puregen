# Java Models Example

```java
import com.example.proto.v1.*;
import com.fasterxml.jackson.databind.ObjectMapper;

public class ModelsExample {
    public static void main(String[] args) throws Exception {
        // Create a new user using builder pattern
        User user = new User.Builder()
            .setId(1)
            .setName("John Doe")
            .setEmail("john@example.com")
            .setIsActive(true)
            .build();
        
        // Create profile
        UserProfile profile = new UserProfile.Builder()
            .setBio("Software Engineer")
            .setAvatarUrl("https://example.com/avatar.jpg")
            .setCreatedAt(1640995200L)
            .build();
        
        user.setProfile(profile);
        
        // Validate
        if (!user.validate()) {
            throw new RuntimeException("Validation failed");
        }
        
        // Convert to JSON
        String json = user.toJson();
        System.out.println("User JSON: " + json);
        
        // Create from JSON
        User newUser = User.fromJson(json);
        System.out.println("Restored user: " + newUser.getName());
    }
}
```
