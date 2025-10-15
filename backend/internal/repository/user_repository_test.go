package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
)

func TestGetUserByID(t *testing.T) {
	// Setup Mock Database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a repository instance with a fake database
	repo := NewUserRepository(db)

	// Define Data & Expectations
	// Expected user data will be returned by the database
	expectedUser := &model.User{
		ID:        "abc-123",
		FullName:  "Test User",
		Email:     "test@example.com",
		Role:      "student",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Expected SQL query will be executed by function
	// Use regexp.QuoteMeta to "escape" special characters in SQL
	expectedSQL := regexp.QuoteMeta(`SELECT id, full_name, email, role, status, created_at, updated_at FROM users WHERE id = $1`)

	// Define the data rows that will be 'returned' by the fake database
	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "role", "status", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.FullName, expectedUser.Email, expectedUser.Role, expectedUser.Status, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	// Set Expectations on the Mock
	mock.ExpectQuery(expectedSQL).WithArgs(expectedUser.ID).WillReturnRows(rows)

	// Run the Function to be tested
	user, err := repo.GetUserByID(expectedUser.ID)

	// Check the Result (Assert)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if user == nil {
		t.Errorf("expected user to be returned, but got nil")
	}
	if user.Email != expectedUser.Email {
		t.Errorf("expected email %s, but got %s", expectedUser.Email, user.Email)
	}

	// Ensure all expectations are met
	// The test will fail if the fake database does not accept the expected query
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUser(t *testing.T) {
	// Setup Mock Database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	// Define input data
	newUser := &model.User{
		FullName: "New Test User",
		Email:    "newuser@example.com",
		Password: "password123",
		Role:     "student",
	}

	// Query SQL that is expected to be executed
	expectedSQL := regexp.QuoteMeta(`INSERT INTO users (full_name, email, password_hash, role) VALUES ($1, $2, $3, $4)`)

	// Set expectations in the mock
	mock.ExpectExec(expectedSQL).
		WithArgs(newUser.FullName, newUser.Email, sqlmock.AnyArg(), newUser.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Run function to be tested
	err = repo.CreateUser(newUser)

	// Check the result (Assert)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Ensure all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUserStatus(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	// Define input data
	userID := "test-user-id"
	newStatus := "active"

	// Query sql that is expected to be executed
	expectedSQL := regexp.QuoteMeta(`UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2`)

	// Set expectation in mock
	mock.ExpectExec(expectedSQL).
		WithArgs(newStatus, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Run function to be tested
	err = repo.UpdateUserStatus(userID, newStatus)

	// Check the result (Assert)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Ensure all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}