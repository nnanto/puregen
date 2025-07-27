# HTTP Transport implementation for Python
# This transport handles method routing based on method names and service names

import requests
from typing import Dict, Any, Type, TypeVar
from abc import ABC, abstractmethod
import re

T = TypeVar("T")


class Transport(ABC):
    """Abstract transport interface for client communication"""

    @abstractmethod
    def send(
        self,
        ctx: Dict[str, Any],
        method_name: str,
        input_data: Any,
        output_type: Type[T],
    ) -> T:
        """Send request and return response"""
        pass


class HTTPTransport(Transport):
    """HTTP Transport implementation that routes methods based on naming conventions"""

    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip("/")
        self.session = requests.Session()

    def send(
        self,
        ctx: Dict[str, Any],
        method_name: str,
        input_data: Any,
        output_type: Type[T],
    ) -> T:
        """Send HTTP request based on method name and routing rules"""

        # Parse service name and method from method_name (format: ServiceName_MethodName)
        parts = method_name.split("_", 1)
        if len(parts) < 2:
            raise ValueError(f"Invalid method name format: {method_name}")

        service_name = parts[0].lower()
        # Remove "Service" suffix if present
        if service_name.endswith("service"):
            service_name = service_name[:-7]

        method_name_part = parts[1]

        # Convert method name to snake_case for URL
        method_path = self._camel_to_snake_case(method_name_part)

        # Build base endpoint
        endpoint = f"/{service_name}/{method_path}"
        url = self.base_url + endpoint

        # Prepare headers from context
        headers = {}
        if ctx:
            for key, value in ctx.items():
                if key.startswith("header."):
                    header_name = key[7:]  # Remove 'header.' prefix
                    headers[header_name] = str(value)

        # Prepare data for request
        json_data = None
        params = None

        if input_data is not None:
            if method_name_part.startswith("Get"):
                # For GET requests, convert input to query parameters
                params = self._object_to_params(input_data)
            else:
                # For other methods, convert input to JSON
                json_data = self._object_to_dict(input_data)

        # Make HTTP request based on method name
        try:
            if method_name_part.startswith("Get"):
                response = self.session.get(url, params=params, headers=headers)
            elif method_name_part.startswith("Update"):
                response = self.session.put(url, json=json_data, headers=headers)
            elif method_name_part.startswith("Delete"):
                response = self.session.delete(url, json=json_data, headers=headers)
            else:
                # Default to POST for everything else (Create, Start, Describe, etc.)
                response = self.session.post(url, json=json_data, headers=headers)

            # Check for HTTP errors
            response.raise_for_status()

            # Parse JSON response
            response_data = response.json()

            # Convert to output type
            if hasattr(output_type, "from_dict"):
                return output_type.from_dict(response_data)
            else:
                # Fallback: try to instantiate with the response data
                try:
                    return output_type(**response_data)
                except Exception:
                    # If that fails, return the raw data
                    return response_data

        except requests.RequestException as e:
            raise RuntimeError(f"HTTP request failed: {e}")

    def _object_to_params(self, obj: Any) -> Dict[str, Any]:
        """Convert an object to query parameters for GET requests"""
        if obj is None:
            return {}

        params = {}

        if hasattr(obj, "to_dict"):
            obj_dict = obj.to_dict()
        elif hasattr(obj, "__dict__"):
            obj_dict = obj.__dict__
        elif isinstance(obj, dict):
            obj_dict = obj
        else:
            return {}

        for key, value in obj_dict.items():
            if value is None:
                continue

            if isinstance(value, str):
                if value:  # Non-empty string
                    params[key] = value
            elif isinstance(value, (int, float)):
                if value != 0:  # Non-zero number
                    params[key] = value
            elif isinstance(value, bool):
                params[key] = value
            elif isinstance(value, list):
                # Handle lists - requests will handle multiple values automatically
                if value:  # Non-empty list
                    params[key] = value
            else:
                # For complex objects, convert to string
                params[key] = str(value)

        return params

    def _object_to_dict(self, obj: Any) -> Dict[str, Any]:
        """Convert an object to dictionary for JSON serialization"""
        if obj is None:
            return {}

        if hasattr(obj, "to_dict"):
            return obj.to_dict()
        elif hasattr(obj, "__dict__"):
            # Convert object attributes to dict, handling nested objects
            result = {}
            for key, value in obj.__dict__.items():
                if not key.startswith("_"):  # Skip private attributes
                    if hasattr(value, "to_dict"):
                        result[key] = value.to_dict()
                    elif hasattr(value, "__dict__"):
                        result[key] = self._object_to_dict(value)
                    elif isinstance(value, list):
                        result[key] = [
                            self._object_to_dict(item)
                            if hasattr(item, "__dict__") or hasattr(item, "to_dict")
                            else item
                            for item in value
                        ]
                    else:
                        result[key] = value
            return result
        elif isinstance(obj, dict):
            return obj
        else:
            return {"value": obj}

    def _camel_to_snake_case(self, camel_case: str) -> str:
        """Convert camelCase to snake_case for URL paths"""
        # Insert underscore before uppercase letters (except the first character)
        snake_case = re.sub("([a-z0-9])([A-Z])", r"\1_\2", camel_case)
        return snake_case.lower()
