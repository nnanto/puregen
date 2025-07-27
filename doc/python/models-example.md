# Python Models Example

```python
from example.v1.user import User, UserProfile
import json

def main():
    # Create a new user
    user = User(
        id=1,
        name="John Doe",
        email="john@example.com",
        is_active=True
    )
    
    # Create profile
    profile = UserProfile(
        bio="Software Engineer",
        avatar_url="https://example.com/avatar.jpg",
        created_at=1640995200
    )
    user.profile = profile
    
    # Validate
    if not user.validate():
        raise ValueError("Validation failed")
    
    # Convert to JSON
    json_str = user.to_json()
    print(f"User JSON: {json_str}")
    
    # Create from JSON
    new_user = User.from_json(json_str)
    print(f"Restored user: {new_user.name}")

if __name__ == "__main__":
    main()
```
