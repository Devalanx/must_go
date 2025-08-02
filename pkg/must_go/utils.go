package must_go

import (
	"net/http"
)

// Common HTTP error helpers

// MustNotFound panics with 404 if err is not nil
func MustNotFound(err error) {
	MustHTTP(err, http.StatusNotFound, "Resource not found")
}

// MustBadRequest panics with 400 if err is not nil
func MustBadRequest(err error) {
	MustHTTP(err, http.StatusBadRequest, "Bad request")
}

// MustUnauthorized panics with 401 if err is not nil
func MustUnauthorized(err error) {
	MustHTTP(err, http.StatusUnauthorized, "Unauthorized")
}

// MustForbidden panics with 403 if err is not nil
func MustForbidden(err error) {
	MustHTTP(err, http.StatusForbidden, "Forbidden")
}

// MustConflict panics with 409 if err is not nil
func MustConflict(err error) {
	MustHTTP(err, http.StatusConflict, "Resource conflict")
}

// MustValidation panics with 400 validation error if err is not nil
func MustValidation(err error) {
	MustHTTP(err, http.StatusBadRequest, "Validation error")
}

// MustInternal panics with 500 if err is not nil
func MustInternal(err error) {
	MustHTTP(err, http.StatusInternalServerError, "Internal server error")
}

// MustTimeout panics with 408 if err is not nil
func MustTimeout(err error) {
	MustHTTP(err, http.StatusRequestTimeout, "Request timeout")
}

// MustServiceUnavailable panics with 503 if err is not nil
func MustServiceUnavailable(err error) {
	MustHTTP(err, http.StatusServiceUnavailable, "Service unavailable")
}

// MustTooManyRequests panics with 429 if err is not nil
func MustTooManyRequests(err error) {
	MustHTTP(err, http.StatusTooManyRequests, "Too many requests")
}

// MustUnprocessableEntity panics with 422 if err is not nil
func MustUnprocessableEntity(err error) {
	MustHTTP(err, http.StatusUnprocessableEntity, "Unprocessable entity")
}

// Helper functions for common scenarios

// MustParse panics if parsing fails
func MustParse[T any](result T, err error) T {
	Must(err)
	return result
}

// MustParseWithMessage panics with custom message if parsing fails
func MustParseWithMessage[T any](result T, err error, message string) T {
	MustWithMessage(err, message)
	return result
}

// MustParseHTTP panics with HTTP error if parsing fails
func MustParseHTTP[T any](result T, err error, statusCode int, message string) T {
	MustHTTP(err, statusCode, message)
	return result
}

// MustParseHTTPDefault panics with default HTTP error if parsing fails
func MustParseHTTPDefault[T any](result T, err error) T {
	MustHTTPWithDefault(err)
	return result
} 