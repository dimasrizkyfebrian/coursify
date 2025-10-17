package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/google/uuid"
)

func TestCreateCourse(t *testing.T) {
	// Setup Application
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCourseRepository(db)

	// Define the input data and expectations
	newCourse := &model.Course{
		Title:        "Test Course",
		Description:  "A description for the test course.",
		InstructorID: "instructor-123",
	}

	// Data that is expected to be returned by RETURNING
	expectedID := uuid.New().String()
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	// SQL query that is expected to be executed
	expectedSQL := regexp.QuoteMeta(`INSERT INTO courses (title, description, instructor_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`)

	// Set expectations in the Mock
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt)

	mock.ExpectQuery(expectedSQL).
		WithArgs(newCourse.Title, newCourse.Description, newCourse.InstructorID).
		WillReturnRows(rows)

	// Run the function that will be tested
	err = repo.CreateCourse(newCourse)

	// Check the Result (Assert)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Make sure that the returned ID and timestamp are filled into the struct
	if newCourse.ID != expectedID {
		t.Errorf("expected course ID to be '%s', but got '%s'", expectedID, newCourse.ID)
	}

	// Make sure all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}