// HTTP Transport implementation for Go
// This transport handles method routing based on method names and service names

package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// HTTPTransport implements the Transport interface for HTTP communication
type HTTPTransport struct {
	baseURL string
	client  *http.Client
}

// NewHTTPTransport creates a new HTTP transport instance
func NewHTTPTransport(baseURL string) *HTTPTransport {
	return &HTTPTransport{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Send implements the Transport interface
func (t *HTTPTransport) Send(ctx context.Context, methodName string, inputData interface{}, outputType interface{}) (interface{}, error) {
	// Parse service name and method from methodName (format: ServiceName_MethodName)
	parts := strings.Split(methodName, "_")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid method name format: %s", methodName)
	}

	serviceName := strings.ToLower(parts[0])
	// Remove "Service" suffix if present
	if strings.HasSuffix(serviceName, "service") {
		serviceName = serviceName[:len(serviceName)-7]
	}

	methodNamePart := parts[1]

	// Determine HTTP method and endpoint
	var httpMethod string
	var endpoint string
	var body io.Reader

	// Convert method name to lowercase for URL
	methodPath := strings.ToLower(methodNamePart)

	// Determine HTTP method based on method name
	if strings.HasPrefix(methodNamePart, "Get") {
		httpMethod = "GET"
		endpoint = fmt.Sprintf("/%s/%s", serviceName, methodPath)

		// For GET requests, convert input to query parameters
		if inputData != nil {
			queryParams, err := t.structToQueryParams(inputData)
			if err != nil {
				return nil, fmt.Errorf("failed to convert input to query params: %w", err)
			}
			if queryParams != "" {
				endpoint += "?" + queryParams
			}
		}
	} else if strings.HasPrefix(methodNamePart, "Update") {
		httpMethod = "PUT"
		endpoint = fmt.Sprintf("/%s/%s", serviceName, methodPath)

		if inputData != nil {
			jsonData, err := json.Marshal(inputData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request: %w", err)
			}
			body = bytes.NewReader(jsonData)
		}
	} else if strings.HasPrefix(methodNamePart, "Delete") {
		httpMethod = "DELETE"
		endpoint = fmt.Sprintf("/%s/%s", serviceName, methodPath)

		if inputData != nil {
			jsonData, err := json.Marshal(inputData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request: %w", err)
			}
			body = bytes.NewReader(jsonData)
		}
	} else {
		// Default to POST for everything else (Create, Start, Describe, etc.)
		httpMethod = "POST"
		endpoint = fmt.Sprintf("/%s/%s", serviceName, methodPath)

		if inputData != nil {
			jsonData, err := json.Marshal(inputData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request: %w", err)
			}
			body = bytes.NewReader(jsonData)
		}
	}

	// Build full URL
	fullURL := t.baseURL + endpoint

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, httpMethod, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type for requests with body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read response body
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Create result instance based on output type
	result := t.createOutputInstance(outputType)
	if result == nil {
		return nil, fmt.Errorf("failed to create output instance")
	}

	// Unmarshal response
	if err := json.Unmarshal(respData, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// structToQueryParams converts a struct to URL query parameters
func (t *HTTPTransport) structToQueryParams(data interface{}) (string, error) {
	values := url.Values{}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("input must be a struct")
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get JSON tag for field name
		jsonTag := fieldType.Tag.Get("json")
		fieldName := fieldType.Name
		if jsonTag != "" && jsonTag != "-" {
			// Use the part before comma (ignore omitempty, etc.)
			if commaIdx := strings.Index(jsonTag, ","); commaIdx >= 0 {
				fieldName = jsonTag[:commaIdx]
			} else {
				fieldName = jsonTag
			}
		}

		// Convert field value to string
		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				values.Add(fieldName, field.String())
			}
		case reflect.Int, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				values.Add(fieldName, strconv.FormatInt(field.Int(), 10))
			}
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			if field.Uint() != 0 {
				values.Add(fieldName, strconv.FormatUint(field.Uint(), 10))
			}
		case reflect.Bool:
			values.Add(fieldName, strconv.FormatBool(field.Bool()))
		case reflect.Float32, reflect.Float64:
			if field.Float() != 0 {
				values.Add(fieldName, strconv.FormatFloat(field.Float(), 'f', -1, 64))
			}
		case reflect.Slice:
			// Handle slices by adding multiple values with the same key
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elem.Kind() == reflect.String {
					values.Add(fieldName, elem.String())
				} else {
					values.Add(fieldName, fmt.Sprintf("%v", elem.Interface()))
				}
			}
		default:
			// For complex types, skip or handle as needed
			if !field.IsZero() {
				values.Add(fieldName, fmt.Sprintf("%v", field.Interface()))
			}
		}
	}

	return values.Encode(), nil
}

// createOutputInstance creates an instance of the output type
func (t *HTTPTransport) createOutputInstance(outputType interface{}) interface{} {
	if outputType == nil {
		return nil
	}

	// outputType should be a pointer to the desired type (e.g., (*UserResponse)(nil))
	v := reflect.ValueOf(outputType)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil
	}

	// Get the type that the pointer points to
	elemType := v.Type().Elem()

	// Create a new instance of that type
	newValue := reflect.New(elemType)
	return newValue.Interface()
}
