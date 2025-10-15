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

// GetCourseByID method
func (r *CourseRepository) GetCourseByID(courseID string) (*model.Course, error) {
    var course model.Course
    query := `SELECT id, instructor_id, title, description, cover_image_url, created_at, updated_at
               FROM courses WHERE id = $1`

    err := r.DB.QueryRow(query, courseID).Scan(
        &course.ID, &course.InstructorID, &course.Title, &course.Description,
        &course.CoverImageURL, &course.CreatedAt, &course.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &course, nil
}

// UpdateCourse method
func (r *CourseRepository) UpdateCourse(course *model.Course) error {
    query := `UPDATE courses SET title = $1, description = $2, updated_at = NOW() WHERE id = $3`

    _, err := r.DB.Exec(query, course.Title, course.Description, course.ID)
    if err != nil {
        log.Printf("Error updating course: %v", err)
        return err
    }
    return nil
}

// AddMaterialToCourse method
func (r *CourseRepository) AddMaterialToCourse(material *model.LearningMaterial) error {
    // Get the last position to determine the position of new material
    var lastPosition int
    posQuery := `SELECT COALESCE(MAX(position), 0) FROM learning_materials WHERE course_id = $1`
    err := r.DB.QueryRow(posQuery, material.CourseID).Scan(&lastPosition)
    if err != nil {
        log.Printf("Error getting last material position: %v", err)
        return err
    }
    material.Position = lastPosition + 1

    // Insert new material
    insertQuery := `
        INSERT INTO learning_materials (course_id, title, content_type, text_content, position)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `
    err = r.DB.QueryRow(
        insertQuery,
        material.CourseID,
        material.Title,
        material.ContentType,
        material.TextContent,
        material.Position,
    ).Scan(&material.ID, &material.CreatedAt, &material.UpdatedAt)

    if err != nil {
        log.Printf("Error adding material to course: %v", err)
        return err
    }

    return nil
}

// GetMaterialsByCourseID method
func (r *CourseRepository) GetMaterialsByCourseID(courseID string) ([]model.LearningMaterial, error) {
    query := `
        SELECT id, course_id, title, content_type, text_content, video_url, file_url, position, created_at, updated_at
        FROM learning_materials 
        WHERE course_id = $1 
        ORDER BY position ASC
    `
    rows, err := r.DB.Query(query, courseID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var materials []model.LearningMaterial
    for rows.Next() {
        var material model.LearningMaterial
        // Use sql.NullString for fields that can be NULL
        var textContent, videoURL, fileURL sql.NullString

        if err := rows.Scan(
            &material.ID,
            &material.CourseID,
            &material.Title,
            &material.ContentType,
            &textContent,
            &videoURL,
            &fileURL,
            &material.Position,
            &material.CreatedAt,
            &material.UpdatedAt,
        ); err != nil {
            return nil, err
        }

        // Conversion from sql.NullString to a regular string
        if textContent.Valid {
            material.TextContent = textContent.String
        }
        if videoURL.Valid {
            material.VideoURL = videoURL.String
        }
        if fileURL.Valid {
            material.FileURL = fileURL.String
        }

        materials = append(materials, material)
    }
    return materials, nil
}

// UpdateMaterial method
func (r *CourseRepository) UpdateMaterial(material *model.LearningMaterial) error {
	query := `
		UPDATE learning_materials 
		SET title = $1, text_content = $2, video_url = $3, updated_at = NOW()
		WHERE id = $4 AND course_id = $5
	`

    // Execute the update query
	result, err := r.DB.Exec(query, material.Title, material.TextContent, material.VideoURL, material.ID, material.CourseID)
	if err != nil {
		log.Printf("Error updating material: %v", err)
		return err
	}

    // Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Indicates that the material was not found or does not match
	}

	return nil
}

// DeleteMaterial method
func (r *CourseRepository) DeleteMaterial(courseID, materialID string) error {
	query := `DELETE FROM learning_materials WHERE id = $1 AND course_id = $2`

    // Execute the delete query
	result, err := r.DB.Exec(query, materialID, courseID)
	if err != nil {
		return err
	}

    // Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Indicates that the material was not found or does not match
	}

	return nil
}

// GetAllCourses method
func (r *CourseRepository) GetAllCourses() ([]model.Course, error) {
    query := `SELECT id, instructor_id, title, description, cover_image_url, created_at, updated_at
               FROM courses ORDER BY created_at DESC`

    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var courses []model.Course
    for rows.Next() {
        var course model.Course
        if err := rows.Scan(
            &course.ID, &course.InstructorID, &course.Title, &course.Description,
            &course.CoverImageURL, &course.CreatedAt, &course.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        courses = append(courses, course)
    }

    return courses, nil
}

// EnrollStudent method
func (r *CourseRepository) EnrollStudent(studentID, courseID string) error {
    query := `INSERT INTO enrollments (user_id, course_id) VALUES ($1, $2)`

    // Execute the insert query
    _, err := r.DB.Exec(query, studentID, courseID)
    if err != nil {
        log.Printf("Error enrolling student: %v", err)
        return err
    }

    return nil
}

// GetEnrolledCoursesByStudentID method
func (r *CourseRepository) GetEnrolledCoursesByStudentID(studentID string) ([]model.Course, error) {
    query := `
        SELECT c.id, c.instructor_id, c.title, c.description, c.cover_image_url, c.created_at, c.updated_at
        FROM courses c
        JOIN enrollments e ON c.id = e.course_id
        WHERE e.user_id = $1
        ORDER BY e.enrollment_date DESC
    `

    rows, err := r.DB.Query(query, studentID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var courses []model.Course
    for rows.Next() {
        var course model.Course
        if err := rows.Scan(
            &course.ID, &course.InstructorID, &course.Title, &course.Description,
            &course.CoverImageURL, &course.CreatedAt, &course.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        courses = append(courses, course)
    }

    return courses, nil
}

// IsStudentEnrolled method
func (r *CourseRepository) IsStudentEnrolled(studentID, courseID string) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM enrollments WHERE user_id = $1 AND course_id = $2)`

    err := r.DB.QueryRow(query, studentID, courseID).Scan(&exists)
    if err != nil {
        return false, err
    }
    return exists, nil
}

// UpdateCourseCoverImage method
func (r *CourseRepository) UpdateCourseCoverImage(courseID, imageURL string) error {
    query := `UPDATE courses SET cover_image_url = $1, updated_at = NOW() WHERE id = $2`

    // Execute the update query
    result, err := r.DB.Exec(query, imageURL, courseID)
    if err != nil {
        return err
    }

    // Check if any rows were affected
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return sql.ErrNoRows // Indicates that the course was not found or does not match
    }

    return nil
}

// AddFileMaterialToCourse method
func (r *CourseRepository) AddFileMaterialToCourse(material *model.LearningMaterial) error {
    // Get the latest position for new material
    var lastPosition int
    posQuery := `SELECT COALESCE(MAX(position), 0) FROM learning_materials WHERE course_id = $1`
    err := r.DB.QueryRow(posQuery, material.CourseID).Scan(&lastPosition)
    if err != nil {
        log.Printf("Error getting last material position: %v", err)
        return err
    }
    material.Position = lastPosition + 1

    // Add new material with the type 'pdf' and file_url
    insertQuery := `
        INSERT INTO learning_materials (course_id, title, content_type, file_url, position)
        VALUES ($1, $2, 'pdf', $3, $4)
        RETURNING id, created_at, updated_at
    `
    err = r.DB.QueryRow(
        insertQuery,
        material.CourseID,
        material.Title,
        material.FileURL,
        material.Position,
    ).Scan(&material.ID, &material.CreatedAt, &material.UpdatedAt)

    if err != nil {
        log.Printf("Error adding file material to course: %v", err)
        return err
    }

    return nil
}