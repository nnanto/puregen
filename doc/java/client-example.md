# Java Client Example

```java
import com.example.proto.v1.*;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.net.http.*;
import java.net.URI;
import java.util.*;

public class HTTPTransport implements Transport {
    private final String baseUrl;
    private final HttpClient client;
    private final ObjectMapper mapper;
    
    public HTTPTransport(String baseUrl) {
        this.baseUrl = baseUrl;
        this.client = HttpClient.newHttpClient();
        this.mapper = new ObjectMapper();
    }
    
    @Override
    public <T> T send(Map<String, Object> ctx, String methodName, Object inputData, Class<T> responseClass) throws Exception {
        String endpoint;
        String method;
        
        switch (methodName) {
            case UserServiceMethods.UserService_CreateUser:
                endpoint = "/users";
                method = "POST";
                break;
            case UserServiceMethods.UserService_GetUser:
                GetUserRequest req = (GetUserRequest) inputData;
                endpoint = "/users/" + req.getId();
                method = "GET";
                break;
            default:
                throw new IllegalArgumentException("Unknown method: " + methodName);
        }
        
        String url = baseUrl + endpoint;
        HttpRequest.Builder requestBuilder = HttpRequest.newBuilder().uri(URI.create(url));
        
        if ("POST".equals(method)) {
            String json = mapper.writeValueAsString(inputData);
            requestBuilder.POST(HttpRequest.BodyPublishers.ofString(json))
                         .header("Content-Type", "application/json");
        } else {
            requestBuilder.GET();
        }
        
        HttpRequest request = requestBuilder.build();
        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        
        if (response.statusCode() != 200) {
            throw new RuntimeException("HTTP error: " + response.statusCode());
        }
        
        return mapper.readValue(response.body(), responseClass);
    }
}

public class UserClientExample {
    public static void main(String[] args) throws Exception {
        // Create HTTP transport
        HTTPTransport transport = new HTTPTransport("http://localhost:8080");
        
        // Create client
        UserServiceClient client = new UserServiceClient(transport);
        
        Map<String, Object> ctx = new HashMap<>();
        
        // Create a user
        UserProfile profile = new UserProfile.Builder()
            .setBio("Software Engineer")
            .setAvatarUrl("https://example.com/avatar.jpg")
            .setCreatedAt(1640995200L)
            .build();
        
        CreateUserRequest createReq = new CreateUserRequest.Builder()
            .setName("Jane Doe")
            .setEmail("jane@example.com")
            .setProfile(profile)
            .build();
        
        CreateUserResponse createResp = client.createUser(ctx, createReq);
        System.out.println("Created user: " + createResp.getUser().getName());
        
        // Get the user
        GetUserRequest getReq = new GetUserRequest.Builder()
            .setId(createResp.getUser().getId())
            .build();
        
        GetUserResponse getResp = client.getUser(ctx, getReq);
        System.out.println("Retrieved user: " + getResp.getUser().getName());
    }
}
```
