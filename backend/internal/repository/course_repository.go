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

// CreateCourse method
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

// GetCourseByInstructorId method
func (r *CourseRepository) GetCoursesByInstructorID(instructorID string) ([]model.Course, error) {
    query := `SELECT id, instructor_id, title, description, cover_image_url, created_at, updated_at
               FROM courses WHERE instructor_id = $1 ORDER BY created_at DESC`

    rows, err := r.DB.Query(query, instructorID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var courses []model.Course
    for rows.Next() {
        var course model.Course
        if err := rows.Scan(
            &course.ID,
            &course.InstructorID,
            &course.Title,
            &course.Description,
            &course.CoverImageURL,
            &course.CreatedAt,
            &course.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        courses = append(courses, course)
    }

    return courses, nil
}