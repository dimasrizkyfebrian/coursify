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
	"github.com/dimasrizkyfebrian/coursify/internal/handler/middleware"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// setupTestApp initializes the entire application for testing
func setupTestApp() (*chi.Mux, *sql.DB, func()) {
	// Load environment variables from .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file for test: %v", err)
	}

	os.Setenv("DB_NAME", "coursify_test")

	db := database.ConnectDB()

	// Create the router and all its dependencies
	r := chi.NewRouter()
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// --- Public Route ---
	r.Post("/api/login", userHandler.Login)
	r.Post("/api/register", userHandler.Register)

	// --- Protected Admin Route ---
	r.Group(func(r chi.Router) {
        r.Use(middleware.AuthMiddleware)
        r.Use(middleware.AdminOnly)

        r.Put("/api/admin/users/{id}/approve", userHandler.ApproveUser)
    })

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

	// Clean the users table before each test
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("Failed to clean users table before test: %v", err)
	}

	// Preparation data test
	// Create test users directly in the test database
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Insert an active test user into the database
	_, err = db.Exec("INSERT INTO users (full_name, email, password_hash, role, status) VALUES ($1, $2, $3, $4, $5)",
		"Active User", "active@example.com", string(hashedPassword), "student", "active")
	if err != nil {
		t.Fatalf("Failed to insert active user for test: %v", err)
	}

	// Insert a pending test user into the database
	_, err = db.Exec("INSERT INTO users (full_name, email, password_hash, role, status) VALUES ($1, $2, $3, $4, $5)",
		"Pending User", "pending@example.com", string(hashedPassword), "student", "pending")
	if err != nil {
		t.Fatalf("Failed to insert pending user for test: %v", err)
	}

	// Define Test Cases
	testCases := []struct {
		name           string
		email          string
		password       string
		expectedStatus int
	}{
		{
			name:           "successful login",
			email:          "active@example.com",
			password:       "password123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "wrong password",
			email:          "active@example.com",
			password:       "wrongpassword",
			expectedStatus: http.StatusUnauthorized, // 401
		},
		{
			name:           "email not found",
			email:          "notfound@example.com",
			password:       "password123",
			expectedStatus: http.StatusUnauthorized, // 401
		},
		{
			name:           "account not active",
			email:          "pending@example.com",
			password:       "password123",
			expectedStatus: http.StatusForbidden, // 403
		},
	}

	// Run Test All Scenario
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			credentials := map[string]string{
				"email":    tc.email,
				"password": tc.password,
			}
			body, _ := json.Marshal(credentials)

			resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Check result (Assert)
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("expected status %v; got %v", tc.expectedStatus, resp.Status)
			}

			// Specifically for a successful test, check the token
			if tc.expectedStatus == http.StatusOK {
				var responseBody map[string]string
				json.NewDecoder(resp.Body).Decode(&responseBody)
				if _, ok := responseBody["token"]; !ok {
					t.Errorf("expected response body to contain a token")
				}
			}
		})
	}
}

func TestRegisterIntegration(t *testing.T) {
	// Setup application
	router, db, teardown := setupTestApp()
	defer teardown()

	// Create a test server that uses an application router
	server := httptest.NewServer(router)
	defer server.Close()

	// Clean the users table before each test
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("Failed to clean users table before test: %v", err)
	}

	// Run Test Scenario: Registration Successful
	t.Run("successful registration", func(t *testing.T) {
		// Create a request body with new user data
		newUser := map[string]string{
			"full_name": "New Register Test",
			"email":     "register@example.com",
			"password":  "password123",
			"role":      "student",
		}
		body, _ := json.Marshal(newUser)

		// Send a request to the test server
		resp, err := http.Post(server.URL+"/api/register", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Check the Results (Assert)
		// Check Status Code
		if resp.StatusCode != http.StatusCreated { // Should be 201 Created
			t.Errorf("expected status %v; got %v", http.StatusCreated, resp.Status)
		}

		// Check whether the user was actually created in the database
		var user model.User
		err = db.QueryRow("SELECT id, full_name, email, role, status FROM users WHERE email = $1", "register@example.com").
			Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Status)
		
		if err != nil {
			t.Fatalf("Failed to find user in database after registration: %v", err)
		}

		// Check if the new user's status is 'pending'
		if user.Status != "pending" {
			t.Errorf("expected user status to be 'pending'; got '%s'", user.Status)
		}
	})
}

func TestApproveUserIntegration(t *testing.T) {
	// Setup Application
	router, db, teardown := setupTestApp()
	defer teardown()
	server := httptest.NewServer(router)
	defer server.Close()

	// Clean the users table before each test
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("Failed to clean users table: %v", err)
	}

	// Data test preparation
	adminUser := model.User{FullName: "Admin Test", Email: "admin@test.com", Role: "admin", Status: "active"}
	instructorUser := model.User{FullName: "Instructor Test", Email: "instructor@test.com", Role: "instructor", Status: "active"}
	pendingUser := model.User{FullName: "Pending Test", Email: "pending@test.com", Role: "student", Status: "pending"}

	// Insert test users into the database
	for _, u := range []*model.User{&adminUser, &instructorUser, &pendingUser} {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		err := db.QueryRow("INSERT INTO users (full_name, email, password_hash, role, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			u.FullName, u.Email, string(hashedPassword), u.Role, u.Status).Scan(&u.ID)
		if err != nil {
			t.Fatalf("Failed to insert user %s: %v", u.Email, err)
		}
	}

	// Helper function for login and obtaining the original token
	getToken := func(email, password string) string {
		credentials := map[string]string{"email": email, "password": password}
		body, _ := json.Marshal(credentials)
		resp, _ := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(body))
		var tokenResp map[string]string
		json.NewDecoder(resp.Body).Decode(&tokenResp)
		return tokenResp["token"]
	}

	// Test Cases
	t.Run("fails when non-admin tries to approve", func(t *testing.T) {
		token := getToken("instructor@test.com", "password123")

		req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/admin/users/"+pendingUser.ID+"/approve", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden; got %v", resp.Status)
		}
	})

	t.Run("successfully approves user when admin", func(t *testing.T) {
		token := getToken("admin@test.com", "password123")

		req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/admin/users/"+pendingUser.ID+"/approve", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200 OK; got %v", resp.Status)
		}

		// Verify changes in the database
		var updatedStatus string
		err = db.QueryRow("SELECT status FROM users WHERE id = $1", pendingUser.ID).Scan(&updatedStatus)
		if err != nil {
			t.Fatalf("Failed to query updated user: %v", err)
		}

		if updatedStatus != "active" {
			t.Errorf("expected status to be 'active'; got '%s'", updatedStatus)
		}
	})
}