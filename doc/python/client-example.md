# Python Client Example

```python
from example.v1.user import (
    User, UserProfile, CreateUserRequest, CreateUserResponse,
    GetUserRequest, GetUserResponse, UserServiceClient, Transport,
    UserServiceMethods
)
from typing import Dict, Any
import requests
import json

class HTTPTransport(Transport):
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.session = requests.Session()
    
    def send(self, ctx: Dict[str, Any], method_name: str, input_data: Any, output_type: type) -> Any:
        if method_name == UserServiceMethods.UserService_CreateUser:
            endpoint = "/users"
            method = "POST"
        elif method_name == UserServiceMethods.UserService_GetUser:
            endpoint = f"/users/{input_data.id}"
            method = "GET"
        else:
            raise ValueError(f"Unknown method: {method_name}")
        
        url = self.base_url + endpoint
        
        if method == "POST":
            data = input_data.to_dict() if hasattr(input_data, 'to_dict') else input_data
            response = self.session.post(url, json=data)
        else:
            response = self.session.get(url)
        
        response.raise_for_status()
        
        response_data = response.json()
        return output_type.from_dict(response_data)

def main():
    # Create HTTP transport
    transport = HTTPTransport("http://localhost:8080")
    
    # Create client
    client = UserServiceClient(transport)
    
    ctx = {}
    
    # Create a user
    profile = UserProfile(
        bio="Software Engineer",
        avatar_url="https://example.com/avatar.jpg",
        created_at=1640995200
    )
    
    create_req = CreateUserRequest(
        name="Jane Doe",
        email="jane@example.com",
        profile=profile
    )
    
    create_resp = client.create_user(ctx, create_req)
    print(f"Created user: {create_resp.user.name}")
    
    # Get the user
    get_req = GetUserRequest(id=create_resp.user.id)
    get_resp = client.get_user(ctx, get_req)
    print(f"Retrieved user: {get_resp.user.name}")

if __name__ == "__main__":
    main()
```
