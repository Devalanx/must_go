package must_go

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

// TestExampleBasic demonstrates basic usage of the must_go package
func TestExampleBasic(t *testing.T) {
	// Basic error handling
	err := fmt.Errorf("some error")
	// Note: Must(err) would panic here, but we're just demonstrating the usage
	_ = err
}

// TestExampleWithMessage demonstrates custom error messages
func TestExampleWithMessage(t *testing.T) {
	err := fmt.Errorf("database connection failed")
	// Note: MustWithMessage(err, "Failed to connect to database") would panic here
	_ = err
	// This would panic with: "Failed to connect to database: database connection failed"
}

// TestExampleHTTP demonstrates HTTP error handling
func TestExampleHTTP(t *testing.T) {
	err := fmt.Errorf("user not found")
	// Note: MustHTTP(err, http.StatusNotFound, "User not found") would panic here
	_ = err
	// This would panic with HTTPError{StatusCode: 404, Message: "User not found"}
}

// TestExampleHTTPDefault demonstrates automatic HTTP error detection
func TestExampleHTTPDefault(t *testing.T) {
	err := fmt.Errorf("user not found")
	// Note: MustHTTPWithDefault(err) would panic here
	_ = err
	// This would panic with HTTPError{StatusCode: 404, Message: "Resource not found"}
}

// TestExampleMiddleware demonstrates how to use the recovery middleware
func TestExampleMiddleware(t *testing.T) {
	// Create a simple HTTP server with recovery middleware
	mux := http.NewServeMux()
	
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		// This would panic, but the middleware would recover and return proper HTTP response
		err := fmt.Errorf("user not found")
		// Note: MustHTTP(err, http.StatusNotFound, "User not found") would panic here
		_ = err
	})
	
	// Wrap the mux with recovery middleware
	handler := RecoveryMiddleware(mux)
	
	// Start server
	// http.ListenAndServe(":8080", handler)
	_ = handler // Use handler to avoid unused variable warning
}

// TestExampleHelperFunctions demonstrates the helper functions
func TestExampleHelperFunctions(t *testing.T) {
	// Parse a string to int
	str := "123"
	num := MustParse(strconv.Atoi(str))
	_ = num // Use num to avoid unused variable warning
	
	// Parse with custom error message
	num2, err2 := strconv.Atoi("abc")
	// Note: MustParseWithMessage(num2, err2, "Invalid number format") would panic here
	_ = num2
	_ = err2
	
	// Parse with HTTP error
	num3, err3 := strconv.Atoi("abc")
	// Note: MustParseHTTP(num3, err3, http.StatusBadRequest, "Invalid number") would panic here
	_ = num3
	_ = err3
}

// TestExampleJSON demonstrates JSON parsing with must_go
func TestExampleJSON(t *testing.T) {
	jsonData := `{"name": "John", "age": 30}`
	
	var user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	
	// Parse JSON with automatic error handling
	err := json.Unmarshal([]byte(jsonData), &user)
	// Note: MustHTTP(err, http.StatusBadRequest, "Invalid JSON format") would panic here
	_ = err
	
	// Or use the helper function
	err2 := json.Unmarshal([]byte(jsonData), &user)
	// Note: MustParseHTTPDefault(&user, err2) would panic here
	_ = err2
	_ = user
}

// TestExampleCustomRecovery demonstrates custom panic handling
func TestExampleCustomRecovery(t *testing.T) {
	customHandler := func(w http.ResponseWriter, r *http.Request, err interface{}) {
		// Custom panic handling logic
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Custom error handling",
		})
	}
	
	mux := http.NewServeMux()
	handler := CustomRecoveryMiddleware(customHandler)(mux)
	_ = handler
}

// TestExampleCommonHTTPErrors demonstrates the common HTTP error helpers
func TestExampleCommonHTTPErrors(t *testing.T) {
	// Not found error
	err := fmt.Errorf("user not found")
	// Note: MustNotFound(err) would panic here
	_ = err
	
	// Bad request error
	validationErr := fmt.Errorf("invalid input")
	// Note: MustBadRequest(validationErr) would panic here
	_ = validationErr
	
	// Unauthorized error
	authErr := fmt.Errorf("invalid credentials")
	// Note: MustUnauthorized(authErr) would panic here
	_ = authErr
	
	// Forbidden error
	permissionErr := fmt.Errorf("insufficient permissions")
	// Note: MustForbidden(permissionErr) would panic here
	_ = permissionErr
	
	// Conflict error
	conflictErr := fmt.Errorf("resource already exists")
	// Note: MustConflict(conflictErr) would panic here
	_ = conflictErr
	
	// Validation error
	validationErr2 := fmt.Errorf("validation failed")
	// Note: MustValidation(validationErr2) would panic here
	_ = validationErr2
	
	// Internal server error
	internalErr := fmt.Errorf("database error")
	// Note: MustInternal(internalErr) would panic here
	_ = internalErr
} 