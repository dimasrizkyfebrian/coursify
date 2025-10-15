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