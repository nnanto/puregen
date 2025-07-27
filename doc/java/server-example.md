# Java Server Example

```java
import com.example.proto.v1.*;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.*;
import java.net.*;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;

public class UserServiceServer implements UserServiceService {
    private final Map<Integer, User> users = new ConcurrentHashMap<>();
    private int nextId = 1;
    
    @Override
    public CreateUserResponse createUser(Map<String, Object> ctx, CreateUserRequest request) throws Exception {
        User user = new User.Builder()
            .setId(nextId++)
            .setName(request.getName())
            .setEmail(request.getEmail())
            .setIsActive(true)
            .setProfile(request.getProfile())
            .build();
        
        users.put(user.getId(), user);
        
        return new CreateUserResponse.Builder()
            .setUser(user)
            .setMessage("User created successfully")
            .build();
    }
    
    @Override
    public GetUserResponse getUser(Map<String, Object> ctx, GetUserRequest request) throws Exception {
        User user = users.get(request.getId());
        if (user == null) {
            throw new RuntimeException("User not found");
        }
        
        return new GetUserResponse.Builder()
            .setUser(user)
            .build();
    }
    
    // Simple HTTP server implementation
    public static void main(String[] args) throws Exception {
        UserServiceServer service = new UserServiceServer();
        HttpServer server = HttpServer.create(new InetSocketAddress(8080), 0);
        
        server.createContext("/users", exchange -> {
            try {
                service.handleRequest(exchange);
            } catch (Exception e) {
                e.printStackTrace();
                exchange.sendResponseHeaders(500, 0);
                exchange.close();
            }
        });
        
        server.setExecutor(null);
        server.start();
        System.out.println("Server started on port 8080");
    }
    
    private void handleRequest(HttpExchange exchange) throws Exception {
        String method = exchange.getRequestMethod();
        String path = exchange.getRequestURI().getPath();
        
        if ("POST".equals(method) && "/users".equals(path)) {
            handleCreateUser(exchange);
        } else if ("GET".equals(method) && path.startsWith("/users/")) {
            handleGetUser(exchange);
        } else {
            exchange.sendResponseHeaders(404, 0);
            exchange.close();
        }
    }
    
    private void handleCreateUser(HttpExchange exchange) throws Exception {
        ObjectMapper mapper = new ObjectMapper();
        CreateUserRequest request = mapper.readValue(exchange.getRequestBody(), CreateUserRequest.class);
        
        CreateUserResponse response = createUser(new HashMap<>(), request);
        
        String responseJson = response.toJson();
        exchange.getResponseHeaders().set("Content-Type", "application/json");
        exchange.sendResponseHeaders(200, responseJson.length());
        
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(responseJson.getBytes());
        }
    }
    
    private void handleGetUser(HttpExchange exchange) throws Exception {
        String path = exchange.getRequestURI().getPath();
        int id = Integer.parseInt(path.substring(7)); // Remove "/users/"
        
        GetUserRequest request = new GetUserRequest.Builder().setId(id).build();
        GetUserResponse response = getUser(new HashMap<>(), request);
        
        String responseJson = response.toJson();
        exchange.getResponseHeaders().set("Content-Type", "application/json");
        exchange.sendResponseHeaders(200, responseJson.length());
        
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(responseJson.getBytes());
        }
    }
}
```
