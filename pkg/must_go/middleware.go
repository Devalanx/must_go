package must_go

import (
	"encoding/json"
	"log"
	"net/http"
)

// RecoveryMiddleware recovers from panics and returns appropriate HTTP responses
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				handlePanic(w, r, err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddlewareFunc is a function-based version of RecoveryMiddleware
func RecoveryMiddlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				handlePanic(w, r, err)
			}
		}()
		next(w, r)
	}
}

// handlePanic processes the panic and returns appropriate HTTP response
func handlePanic(w http.ResponseWriter, r *http.Request, err interface{}) {
	log.Printf("Panic recovered: %v", err)

	// Set default values
	statusCode := http.StatusInternalServerError
	message := "Internal server error"

	// Check if it's our custom HTTPError
	if httpErr, ok := err.(HTTPError); ok {
		statusCode = httpErr.StatusCode
		message = httpErr.Message
	} else if errStr, ok := err.(string); ok {
		// Handle string panics
		message = errStr
	} else if errObj, ok := err.(error); ok {
		// Handle error panics
		message = errObj.Error()
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Create error response
	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
			"status":  statusCode,
		},
	}

	// Encode and send response
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Failed to encode error response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// CustomRecoveryMiddleware allows custom panic handling
func CustomRecoveryMiddleware(panicHandler func(http.ResponseWriter, *http.Request, interface{})) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					panicHandler(w, r, err)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// SimpleRecoveryMiddleware provides a simple recovery that logs and returns 500
func SimpleRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
} 