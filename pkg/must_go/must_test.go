package must_go

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMust(t *testing.T) {
	// Test with nil error (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Must() panicked with nil error: %v", r)
		}
	}()
	Must(nil)
}

func TestMustWithError(t *testing.T) {
	// Test with non-nil error (should panic)
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Must() should have panicked with non-nil error")
		}
		if r.(error).Error() != "test error" {
			t.Errorf("Expected panic with 'test error', got: %v", r)
		}
	}()
	Must(fmt.Errorf("test error"))
}

func TestMustWithMessage(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("MustWithMessage() should have panicked")
		}
		err := r.(error)
		if err.Error() != "Custom message: test error" {
			t.Errorf("Expected 'Custom message: test error', got: %v", err)
		}
	}()
	MustWithMessage(fmt.Errorf("test error"), "Custom message")
}

func TestMustHTTP(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("MustHTTP() should have panicked")
		}
		httpErr, ok := r.(HTTPError)
		if !ok {
			t.Errorf("Expected HTTPError, got: %T", r)
		}
		if httpErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got: %d", httpErr.StatusCode)
		}
		if httpErr.Message != "Not found" {
			t.Errorf("Expected message 'Not found', got: %s", httpErr.Message)
		}
	}()
	MustHTTP(fmt.Errorf("test error"), http.StatusNotFound, "Not found")
}

func TestMustHTTPWithDefault(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "not found",
			err:        fmt.Errorf("user not found"),
			wantStatus: http.StatusNotFound,
			wantMsg:    "Resource not found",
		},
		{
			name:       "unauthorized",
			err:        fmt.Errorf("unauthorized access"),
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "Unauthorized",
		},
		{
			name:       "forbidden",
			err:        fmt.Errorf("forbidden action"),
			wantStatus: http.StatusForbidden,
			wantMsg:    "Forbidden",
		},
		{
			name:       "bad request",
			err:        fmt.Errorf("bad request data"),
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Bad request",
		},
		{
			name:       "validation",
			err:        fmt.Errorf("validation failed"),
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Validation error",
		},
		{
			name:       "timeout",
			err:        fmt.Errorf("request timeout"),
			wantStatus: http.StatusRequestTimeout,
			wantMsg:    "Request timeout",
		},
		{
			name:       "conflict",
			err:        fmt.Errorf("resource conflict"),
			wantStatus: http.StatusConflict,
			wantMsg:    "Resource conflict",
		},
		{
			name:       "unknown error",
			err:        fmt.Errorf("unknown error"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Error("MustHTTPWithDefault() should have panicked")
				}
				httpErr, ok := r.(HTTPError)
				if !ok {
					t.Errorf("Expected HTTPError, got: %T", r)
				}
				if httpErr.StatusCode != tt.wantStatus {
					t.Errorf("Expected status %d, got: %d", tt.wantStatus, httpErr.StatusCode)
				}
				if httpErr.Message != tt.wantMsg {
					t.Errorf("Expected message '%s', got: '%s'", tt.wantMsg, httpErr.Message)
				}
			}()
			MustHTTPWithDefault(tt.err)
		})
	}
}

func TestHTTPError_Error(t *testing.T) {
	err := HTTPError{
		StatusCode: http.StatusNotFound,
		Message:    "Not found",
	}
	expected := "HTTP 404: Not found"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	// Create a handler that panics
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MustHTTP(fmt.Errorf("test error"), http.StatusNotFound, "Test error")
	})

	// Wrap with recovery middleware
	wrappedHandler := RecoveryMiddleware(handler)

	// Create test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Serve the request
	wrappedHandler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got: %d", w.Code)
	}

	// Check content type
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got: %s", w.Header().Get("Content-Type"))
	}

	// Check response body contains error
	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body")
	}
}

func TestMustParse(t *testing.T) {
	// Test successful parsing
	result := MustParse(42, nil)
	if result != 42 {
		t.Errorf("Expected 42, got: %v", result)
	}

	// Test parsing with error
	defer func() {
		r := recover()
		if r == nil {
			t.Error("MustParse() should have panicked with error")
		}
	}()
	MustParse(0, fmt.Errorf("parse error"))
}

func TestMustParseHTTP(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("MustParseHTTP() should have panicked")
		}
		httpErr, ok := r.(HTTPError)
		if !ok {
			t.Errorf("Expected HTTPError, got: %T", r)
		}
		if httpErr.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got: %d", httpErr.StatusCode)
		}
	}()
	MustParseHTTP(0, fmt.Errorf("parse error"), http.StatusBadRequest, "Invalid input")
}

func TestCommonHTTPHelpers(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(error)
		wantStatus int
		wantMsg    string
	}{
		{"MustNotFound", MustNotFound, http.StatusNotFound, "Resource not found"},
		{"MustBadRequest", MustBadRequest, http.StatusBadRequest, "Bad request"},
		{"MustUnauthorized", MustUnauthorized, http.StatusUnauthorized, "Unauthorized"},
		{"MustForbidden", MustForbidden, http.StatusForbidden, "Forbidden"},
		{"MustConflict", MustConflict, http.StatusConflict, "Resource conflict"},
		{"MustValidation", MustValidation, http.StatusBadRequest, "Validation error"},
		{"MustInternal", MustInternal, http.StatusInternalServerError, "Internal server error"},
		{"MustTimeout", MustTimeout, http.StatusRequestTimeout, "Request timeout"},
		{"MustServiceUnavailable", MustServiceUnavailable, http.StatusServiceUnavailable, "Service unavailable"},
		{"MustTooManyRequests", MustTooManyRequests, http.StatusTooManyRequests, "Too many requests"},
		{"MustUnprocessableEntity", MustUnprocessableEntity, http.StatusUnprocessableEntity, "Unprocessable entity"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Error("Function should have panicked")
				}
				httpErr, ok := r.(HTTPError)
				if !ok {
					t.Errorf("Expected HTTPError, got: %T", r)
				}
				if httpErr.StatusCode != tt.wantStatus {
					t.Errorf("Expected status %d, got: %d", tt.wantStatus, httpErr.StatusCode)
				}
				if httpErr.Message != tt.wantMsg {
					t.Errorf("Expected message '%s', got: '%s'", tt.wantMsg, httpErr.Message)
				}
			}()
			tt.fn(fmt.Errorf("test error"))
		})
	}
} 