package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Devalanx/must_go/pkg/must_go"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Mock database
var users = map[int]User{
	1: {ID: 1, Name: "John Doe", Email: "john@example.com"},
	2: {ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from query parameter
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		must_go.MustBadRequest(fmt.Errorf("user id is required"))
	}

	// Parse user ID to integer
	userID := must_go.MustParse(strconv.Atoi(userIDStr))

	// Get user from database
	user, exists := users[userID]
	if !exists {
		must_go.MustNotFound(fmt.Errorf("user with id %d not found", userID))
	}

	// Return user as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	must_go.MustHTTP(err, http.StatusBadRequest, "Invalid JSON format")

	// Validate user data
	if user.Name == "" {
		must_go.MustValidation(fmt.Errorf("name is required"))
	}
	if user.Email == "" {
		must_go.MustValidation(fmt.Errorf("email is required"))
	}

	// Check if user already exists
	for _, existingUser := range users {
		if existingUser.Email == user.Email {
			must_go.MustConflict(fmt.Errorf("user with email %s already exists", user.Email))
		}
	}

	// Generate new ID
	user.ID = len(users) + 1
	users[user.ID] = user

	// Return created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from query parameter
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		must_go.MustBadRequest(fmt.Errorf("user id is required"))
	}

	userID := must_go.MustParse(strconv.Atoi(userIDStr))

	// Check if user exists
	_, exists := users[userID]
	if !exists {
		must_go.MustNotFound(fmt.Errorf("user with id %d not found", userID))
	}

	// Parse JSON request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	must_go.MustHTTP(err, http.StatusBadRequest, "Invalid JSON format")

	// Validate user data
	if user.Name == "" {
		must_go.MustValidation(fmt.Errorf("name is required"))
	}
	if user.Email == "" {
		must_go.MustValidation(fmt.Errorf("email is required"))
	}

	// Update user
	user.ID = userID
	users[userID] = user

	// Return updated user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from query parameter
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		must_go.MustBadRequest(fmt.Errorf("user id is required"))
	}

	userID := must_go.MustParse(strconv.Atoi(userIDStr))

	// Check if user exists
	_, exists := users[userID]
	if !exists {
		must_go.MustNotFound(fmt.Errorf("user with id %d not found", userID))
	}

	// Delete user
	delete(users, userID)

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("GET /users", getUserHandler)
	mux.HandleFunc("POST /users", createUserHandler)
	mux.HandleFunc("PUT /users", updateUserHandler)
	mux.HandleFunc("DELETE /users", deleteUserHandler)

	// Wrap with recovery middleware
	handler := must_go.RecoveryMiddleware(mux)

	log.Println("Server starting on :8080")
	log.Println("Try these endpoints:")
	log.Println("  GET  /users?id=1")
	log.Println("  POST /users (with JSON body)")
	log.Println("  PUT  /users?id=1 (with JSON body)")
	log.Println("  DELETE /users?id=1")

	log.Fatal(http.ListenAndServe(":8080", handler))
} 