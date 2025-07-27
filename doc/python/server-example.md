# Python Server Example

```python
from example.v1.user import (
    User, UserProfile, CreateUserRequest, CreateUserResponse,
    GetUserRequest, GetUserResponse, UserServiceService
)
from typing import Dict, Any
from flask import Flask, request, jsonify
import traceback

class UserServiceImpl(UserServiceService):
    def __init__(self):
        self.users: Dict[int, User] = {}
        self.next_id = 1
    
    def create_user(self, ctx: Dict[str, Any], request: CreateUserRequest) -> CreateUserResponse:
        user = User(
            id=self.next_id,
            name=request.name,
            email=request.email,
            is_active=True,
            profile=request.profile
        )
        
        self.users[user.id] = user
        self.next_id += 1
        
        return CreateUserResponse(
            user=user,
            message="User created successfully"
        )
    
    def get_user(self, ctx: Dict[str, Any], request: GetUserRequest) -> GetUserResponse:
        user = self.users.get(request.id)
        if user is None:
            raise ValueError("User not found")
        
        return GetUserResponse(user=user)

def create_app():
    app = Flask(__name__)
    service = UserServiceImpl()
    
    @app.route('/users', methods=['POST'])
    def create_user():
        try:
            request_data = request.get_json()
            user_request = CreateUserRequest.from_dict(request_data)
            
            response = service.create_user({}, user_request)
            
            return jsonify(response.to_dict())
        except Exception as e:
            return jsonify({'error': str(e)}), 400
    
    @app.route('/users/<int:user_id>', methods=['GET'])
    def get_user(user_id):
        try:
            user_request = GetUserRequest(id=user_id)
            response = service.get_user({}, user_request)
            
            return jsonify(response.to_dict())
        except ValueError as e:
            return jsonify({'error': str(e)}), 404
        except Exception as e:
            return jsonify({'error': str(e)}), 400
    
    return app

def main():
    app = create_app()
    print("Server starting on port 8080")
    app.run(host='localhost', port=8080, debug=True)

if __name__ == "__main__":
    main()
```
