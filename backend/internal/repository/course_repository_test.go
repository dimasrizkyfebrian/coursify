package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/google/uuid"
)

func TestCreateCourse(t *testing.T) {
	// Setup mock database
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

func TestGetCoursesByInstructorID(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCourseRepository(db)

	// Define the input data and expectations
	instructorID := "instructor-123"

	// Course data that is expected will be returned by the database
	expectedCourses := []model.Course{
		{ID: "course-1", InstructorID: instructorID, Title: "Course One", Description: "Desc One"},
		{ID: "course-2", InstructorID: instructorID, Title: "Course Two", Description: "Desc Two"},
	}

	// SQL query that is expected to be executed
	expectedSQL := regexp.QuoteMeta(`SELECT id, instructor_id, title, description, cover_image_url, created_at, updated_at FROM courses WHERE instructor_id = $1 ORDER BY created_at DESC`)

	// Prepare the row of data that will be 'returned' by the fake database
	rows := sqlmock.NewRows([]string{"id", "instructor_id", "title", "description", "cover_image_url", "created_at", "updated_at"}).
		AddRow(expectedCourses[0].ID, expectedCourses[0].InstructorID, expectedCourses[0].Title, expectedCourses[0].Description, sql.NullString{}, time.Now(), time.Now()).
		AddRow(expectedCourses[1].ID, expectedCourses[1].InstructorID, expectedCourses[1].Title, expectedCourses[1].Description, sql.NullString{}, time.Now(), time.Now())

	// Set expectations in the Mock
	mock.ExpectQuery(expectedSQL).WithArgs(instructorID).WillReturnRows(rows)

	// Run the function that will be tested
	courses, err := repo.GetCoursesByInstructorID(instructorID)

	// Check the result (Assert)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(courses) != 2 {
		t.Errorf("expected 2 courses, but got %d", len(courses))
	}
	if courses[0].Title != expectedCourses[0].Title {
		t.Errorf("expected first course title to be '%s', but got '%s'", expectedCourses[0].Title, courses[0].Title)
	}

	// Ensure all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}