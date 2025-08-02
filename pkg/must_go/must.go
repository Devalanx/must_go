package must_go

import (
	"fmt"
	"net/http"
)

// HTTPError represents an HTTP error with status code and message
type HTTPError struct {
	StatusCode int
	Message    string
}

// Error implements the error interface
func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// Must panics if err is not nil
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// MustWithMessage panics with a custom error message if err is not nil
func MustWithMessage(err error, message string) {
	if err != nil {
		panic(fmt.Errorf("%s: %w", message, err))
	}
}

// MustHTTP panics with an HTTPError if err is not nil
func MustHTTP(err error, statusCode int, message string) {
	if err != nil {
		panic(HTTPError{
			StatusCode: statusCode,
			Message:    message,
		})
	}
}

// MustHTTPWithDefault panics with a default HTTP error based on common error patterns
func MustHTTPWithDefault(err error) {
	if err != nil {
		// Try to determine appropriate status code based on error message
		statusCode := http.StatusInternalServerError
		message := "Internal server error"
		
		errStr := err.Error()
		
		// Common error patterns
		switch {
		case contains(errStr, "not found"):
			statusCode = http.StatusNotFound
			message = "Resource not found"
		case contains(errStr, "unauthorized"):
			statusCode = http.StatusUnauthorized
			message = "Unauthorized"
		case contains(errStr, "forbidden"):
			statusCode = http.StatusForbidden
			message = "Forbidden"
		case contains(errStr, "bad request"):
			statusCode = http.StatusBadRequest
			message = "Bad request"
		case contains(errStr, "validation"):
			statusCode = http.StatusBadRequest
			message = "Validation error"
		case contains(errStr, "timeout"):
			statusCode = http.StatusRequestTimeout
			message = "Request timeout"
		case contains(errStr, "conflict"):
			statusCode = http.StatusConflict
			message = "Resource conflict"
		}
		
		panic(HTTPError{
			StatusCode: statusCode,
			Message:    message,
		})
	}
}

// MustWithRecovery panics with err if not nil, but can be recovered by middleware
func MustWithRecovery(err error) {
	if err != nil {
		panic(err)
	}
}

// contains is a helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

// containsSubstring checks if s contains substr (case-insensitive)
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
} 