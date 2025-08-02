# must_go

A Go package that implements the "must" pattern for error handling with HTTP-specific error recovery. This package allows you to panic on errors and recover them gracefully in HTTP handlers with proper status codes and error messages.

## Features

- **Basic Error Handling**: Simple panic on error with `Must(err)`
- **Custom Error Messages**: Add context to errors with `MustWithMessage(err, message)`
- **HTTP-Specific Errors**: Panic with HTTP status codes and messages
- **Automatic Error Detection**: Automatically determine HTTP status codes based on error messages
- **Recovery Middleware**: Recover from panics and return proper HTTP responses
- **Helper Functions**: Common HTTP error scenarios (404, 400, 401, etc.)
- **Generic Parsing**: Type-safe parsing with automatic error handling

## Installation

```bash
go get github.com/Devalanx/must_go
```

## Quick Start

```go
package main

import (
    "net/http"
    "github.com/yourusername/must_go"
)

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        // This will panic, but middleware will recover and return 404
        err := fmt.Errorf("user not found")
        must_go.MustHTTP(err, http.StatusNotFound, "User not found")
    })
    
    // Wrap with recovery middleware
    handler := must_go.RecoveryMiddleware(mux)
    http.ListenAndServe(":8080", handler)
}
```

## Core Functions

### Basic Error Handling

```go
// Panic on any error
must_go.Must(err)

// Panic with custom message
must_go.MustWithMessage(err, "Custom error message")

// Panic with HTTP error
must_go.MustHTTP(err, http.StatusNotFound, "Resource not found")

// Panic with automatic HTTP error detection
must_go.MustHTTPWithDefault(err)
```

### HTTP Error Helpers

```go
// Common HTTP error scenarios
must_go.MustNotFound(err)           // 404
must_go.MustBadRequest(err)         // 400
must_go.MustUnauthorized(err)       // 401
must_go.MustForbidden(err)          // 403
must_go.MustConflict(err)           // 409
must_go.MustValidation(err)         // 400
must_go.MustInternal(err)           // 500
must_go.MustTimeout(err)            // 408
must_go.MustServiceUnavailable(err) // 503
must_go.MustTooManyRequests(err)    // 429
must_go.MustUnprocessableEntity(err) // 422
```

### Generic Parsing

```go
// Parse with automatic error handling
num := must_go.MustParse(strconv.Atoi("123"))

// Parse with custom message
num := must_go.MustParseWithMessage(strconv.Atoi("abc"), "Invalid number")

// Parse with HTTP error
num := must_go.MustParseHTTP(strconv.Atoi("abc"), http.StatusBadRequest, "Invalid number")

// Parse with automatic HTTP error detection
num := must_go.MustParseHTTPDefault(strconv.Atoi("abc"))
```

## Middleware

### Recovery Middleware

```go
// Standard recovery middleware
handler := must_go.RecoveryMiddleware(mux)

// Function-based middleware
handler := must_go.RecoveryMiddlewareFunc(func(w http.ResponseWriter, r *http.Request) {
    // Your handler code
})

// Custom recovery middleware
customHandler := func(w http.ResponseWriter, r *http.Request, err interface{}) {
    // Custom panic handling
}
handler := must_go.CustomRecoveryMiddleware(customHandler)(mux)

// Simple recovery (logs and returns 500)
handler := must_go.SimpleRecoveryMiddleware(mux)
```

## Error Response Format

When a panic is recovered, the middleware returns a JSON response:

```json
{
  "error": {
    "message": "User not found",
    "status": 404
  }
}
```

## Automatic Error Detection

The `MustHTTPWithDefault` function automatically detects common error patterns:

- "not found" → 404
- "unauthorized" → 401
- "forbidden" → 403
- "bad request" → 400
- "validation" → 400
- "timeout" → 408
- "conflict" → 409

## Examples

### Database Operations

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    userID := must_go.MustParse(strconv.Atoi(r.URL.Query().Get("id")))
    
    user, err := db.GetUser(userID)
    must_go.MustNotFound(err) // Panics with 404 if user not found
    
    json.NewEncoder(w).Encode(user)
}
```

### JSON Parsing

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    must_go.MustHTTP(err, http.StatusBadRequest, "Invalid JSON format")
    
    // Process user...
}
```

### Validation

```go
func validateUser(user User) {
    if user.Email == "" {
        must_go.MustValidation(fmt.Errorf("email is required"))
    }
    
    if user.Age < 0 {
        must_go.MustValidation(fmt.Errorf("age must be positive"))
    }
}
```

### Custom Error Handling

```go
func customErrorHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
    w.Header().Set("Content-Type", "application/json")
    
    switch e := err.(type) {
    case must_go.HTTPError:
        w.WriteHeader(e.StatusCode)
        json.NewEncoder(w).Encode(map[string]string{
            "error": e.Message,
        })
    default:
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Internal server error",
        })
    }
}

handler := must_go.CustomRecoveryMiddleware(customErrorHandler)(mux)
```

## Best Practices

1. **Use specific error helpers** when you know the exact HTTP status code needed
2. **Use automatic detection** when you want the package to guess based on error messages
3. **Always wrap HTTP handlers** with recovery middleware
4. **Use generic parsing** for type-safe operations with automatic error handling
5. **Provide meaningful error messages** for better user experience

## Error Types

The package handles various error types:

- `error` interface
- `string` panics
- Custom `HTTPError` struct
- Any other interface{} panic

## Contributing

Feel free to submit issues and pull requests to improve the package.

## License

MIT License 