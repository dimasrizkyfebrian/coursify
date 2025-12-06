package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimasrizkyfebrian/coursify/internal/config"
	"github.com/dimasrizkyfebrian/coursify/internal/database"
	"github.com/dimasrizkyfebrian/coursify/internal/handler"
	"github.com/dimasrizkyfebrian/coursify/internal/handler/middleware"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
	"github.com/dimasrizkyfebrian/coursify/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// setupTestApp initializes the router and database for testing
func setupTestApp() (*chi.Mux, *sql.DB, func()) {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		// If .env not found, try loading from root (for CI/CD or different execution context)
		if err := godotenv.Load("../.env"); err != nil {
			// Fallback to system environment variables if .env is missing
		}
	}

	cfg := config.LoadConfig()
	db := database.ConnectDB(cfg)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg)
	userHandler := handler.NewUserHandler(userService)

	r := chi.NewRouter()

	// Public Routes
	r.Post("/api/register", userHandler.Register)
	r.Post("/api/login", userHandler.Login)

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWTSecretKey))
		r.Get("/api/profile", userHandler.GetProfile)
	})

	// Admin Routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWTSecretKey))
		r.Use(middleware.AdminOnly)
		r.Get("/api/admin/users/pending", userHandler.GetPendingUsers)
		r.Put("/api/admin/users/{id}/approve", userHandler.ApproveUser)
		r.Put("/api/admin/users/{id}/reject", userHandler.RejectUser)
	})

	// Clear users table before tests
	db.Exec("DELETE FROM users")

	return r, db, func() { db.Close() }
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

func TestDeleteUserIntegration(t *testing.T) {
	// Setup Application
	router, db, teardown := setupTestApp()
	defer teardown()
	server := httptest.NewServer(router)
	defer server.Close()

	// Data test preparation
	getToken := func(email, password string) string {
		credentials := map[string]string{"email": email, "password": password}
		body, _ := json.Marshal(credentials)
		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(body))
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Helper function getToken failed for email %s", email)
		}
		var tokenResp map[string]string
		json.NewDecoder(resp.Body).Decode(&tokenResp)
		resp.Body.Close()
		return tokenResp["token"]
	}

	// Test case
	t.Run("fails when non-admin tries to delete", func(t *testing.T) {
		// Set up specific data for this test
		db.Exec("DELETE FROM users")
		instructorUser := model.User{FullName: "Instructor Test", Email: "instructor@test.com", Role: "instructor", Status: "active"}
		userToDelete := model.User{FullName: "To Delete", Email: "todelete@test.com", Role: "student", Status: "active"}
		for _, u := range []*model.User{&instructorUser, &userToDelete} {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
			db.QueryRow("INSERT INTO users (full_name, email, password_hash, role, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
				u.FullName, u.Email, string(hashedPassword), u.Role, u.Status).Scan(&u.ID)
		}

		token := getToken("instructor@test.com", "password123") // Get tokens as an instructor

		req, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/admin/users/"+userToDelete.ID, nil)
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

	t.Run("successfully deletes user when admin", func(t *testing.T) {
		// Set up specific data for this test
		db.Exec("DELETE FROM users")
		adminUser := model.User{FullName: "Admin Test", Email: "admin@test.com", Role: "admin", Status: "active"}
		userToDelete := model.User{FullName: "To Delete", Email: "todelete@test.com", Role: "student", Status: "active"}
		for _, u := range []*model.User{&adminUser, &userToDelete} {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
			db.QueryRow("INSERT INTO users (full_name, email, password_hash, role, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
				u.FullName, u.Email, string(hashedPassword), u.Role, u.Status).Scan(&u.ID)
		}

		token := getToken("admin@test.com", "password123") // Get tokens as an admin

		req, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/admin/users/"+userToDelete.ID, nil)
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

		// Verify that the user is truly removed from the database
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", userToDelete.ID).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query deleted user: %v", err)
		}

		if count != 0 {
			t.Errorf("expected user count to be 0; got %d", count)
		}
	})
}
