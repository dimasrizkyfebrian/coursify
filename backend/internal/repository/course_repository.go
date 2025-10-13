package repository

import (
	"database/sql"
	"log"

	"github.com/dimasrizkyfebrian/coursify/internal/model"
)

type CourseRepository struct {
    DB *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
    return &CourseRepository{DB: db}
}

// CreateCourse Method
func (r *CourseRepository) CreateCourse(course *model.Course) error {
    query := `INSERT INTO courses (title, description, instructor_id) 
               VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

    err := r.DB.QueryRow(query, course.Title, course.Description, course.InstructorID).Scan(&course.ID, &course.CreatedAt, &course.UpdatedAt)
    if err != nil {
        log.Printf("Error creating course: %v", err)
        return err
    }

    return nil
}