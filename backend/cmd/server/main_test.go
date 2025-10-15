package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dimasrizkyfebrian/coursify/internal/database"
	"github.com/dimasrizkyfebrian/coursify/internal/handler"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// setupTestApp initializes the entire application for testing
func setupTestApp() (*chi.Mux, *sql.DB, func()) {
	// Muat environment variables untuk tes
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file for test: %v", err)
	}

	os.Setenv("DB_NAME", "coursify_test")

	db := database.ConnectDB()

	// Create the router and all its dependencies
	r := chi.NewRouter()
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// Register the route that will be tested
	r.Post("/api/login", userHandler.Login)

	// Return the router and teardown function to clean the DB
	teardown := func() {
		db.Exec("DELETE FROM users") // Delete all user data after the test is completed
		db.Close()
	}

	return r, db, teardown
}


func TestLoginIntegration(t *testing.T) {
	// Setup application
	router, db, teardown := setupTestApp()
	defer teardown() // Make sure the teardown is executed at the end

	// Create a test server that uses an application router
	server := httptest.NewServer(router)
	defer server.Close()

	// Preparation data test
	// Create test users directly in the test database
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Insert a test user into the database
	_, err := db.Exec("INSERT INTO users (full_name, email, password_hash, role, status) VALUES ($1, $2, $3, $4, $5)", 
        "Test User", "test@example.com", string(hashedPassword), "student", "active")
    if err != nil {
        t.Fatalf("Failed to insert test user: %v", err)
    }

	// Run Test Scenario: Successful Login
	t.Run("successful login", func(t *testing.T) {
		// Create body request
		credentials := map[string]string{
			"email":    "test@example.com",
			"password": password,
		}
		body, _ := json.Marshal(credentials)

		// Send a request to the test server
		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Check Results (Assert)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}

		var responseBody map[string]string
		json.NewDecoder(resp.Body).Decode(&responseBody)

		if _, ok := responseBody["token"]; !ok {
			t.Errorf("expected response body to contain a token")
		}
	})
}