package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dimasrizkyfebrian/coursify/internal/handler/middleware"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
	"github.com/go-chi/chi/v5"
)

type CourseHandler struct {
    Repo *repository.CourseRepository
}

func NewCourseHandler(repo *repository.CourseRepository) *CourseHandler {
    return &CourseHandler{Repo: repo}
}

type createCourseRequest struct {
	Title       string `json:"title" example:"Introduction to Go"`
	Description string `json:"description" example:"A beginner's guide to Golang."`
}

type addMaterialRequest struct {
	Title       string `json:"title" example:"Chapter 1: Introduction"`
	ContentType string `json:"content_type" enums:"text,video,pdf"`
	TextContent string `json:"text_content,omitempty" example:"This is the lesson content."`
	VideoURL    string `json:"video_url,omitempty" example:"https://youtube.com/watch?v=..."`
}

type courseWithMaterials struct {
    model.Course
    Materials []model.LearningMaterial `json:"materials"`
}

// @Summary      Create a new course (Instructor only)
// @Description  Creates a new course for the logged-in instructor.
// @Tags         Instructor
// @Accept       json
// @Produce      json
// @Param        course body createCourseRequest true "Course Information"
// @Success      201  {object}  model.Course
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /instructor/courses [post]
// @Security     BearerAuth
// CreateCourse handles requests to create new courses
func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
    // Get the instructor ID from the JWT context
    instructorID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Could not retrieve instructor ID from context", http.StatusInternalServerError)
        return
    }

	// Parse the request body into a Course struct
    var course model.Course
    if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate the course fields
    if strings.TrimSpace(course.Title) == "" || strings.TrimSpace(course.Description) == "" {
        http.Error(w, "Title and description cannot be empty", http.StatusBadRequest)
        return
    }

    // Set the instructor ID from the token, not from the request body
    course.InstructorID = instructorID

	// Create the course in the database
    if err := h.Repo.CreateCourse(&course); err != nil {
        http.Error(w, "Failed to create course", http.StatusInternalServerError)
        return
    }

	// Respond with the created course
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(course)
}

// @Summary      Get my courses (Instructor only)
// @Description  Retrieves a list of all courses created by the logged-in instructor.
// @Tags         Instructor
// @Produce      json
// @Success      200  {array}   model.Course
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /instructor/courses [get]
// @Security     BearerAuth
// GetMyCourses handles requests to retrieve courses owned by the logged-in instructor
func (h *CourseHandler) GetMyCourses(w http.ResponseWriter, r *http.Request) {
    // Get the instructor ID from the JWT context
    instructorID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Could not retrieve instructor ID from context", http.StatusInternalServerError)
        return
    }

	// Fetch courses owned by the instructor from the database
    courses, err := h.Repo.GetCoursesByInstructorID(instructorID)
    if err != nil {
        http.Error(w, "Failed to fetch courses", http.StatusInternalServerError)
        return
    }

	// Respond with the list of courses
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(courses)
}

// @Summary      Get my course details (Instructor only)
// @Description  Retrieves the details of a specific course owned by the logged-in instructor.
// @Tags         Instructor
// @Produce      json
// @Param        id   path      string  true  "Course ID"
// @Success      200  {object}  model.Course
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /instructor/courses/{id} [get]
// @Security     BearerAuth
// GetMyCourseDetails handles requests to retrieve details of a specific course
func (h *CourseHandler) GetMyCourseDetails(w http.ResponseWriter, r *http.Request) {
    // Get the instructor ID from the JWT context
    instructorID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Could not retrieve instructor ID from context", http.StatusInternalServerError)
        return
    }

    // Get the course ID from the URL parameter
    courseID := chi.URLParam(r, "id")

    // Get course from repo
    course, err := h.Repo.GetCourseByID(courseID)
    if err != nil || course == nil {
        http.Error(w, "Course not found", http.StatusNotFound)
        return
    }

    // Check if the course belongs to the instructor
    if course.InstructorID != instructorID {
        http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
        return
    }

    // Respond with the course details
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(course)
}

// @Summary      Update a course (Instructor only)
// @Description  Updates the title and description of a course owned by the logged-in instructor.
// @Tags         Instructor
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Course ID"
// @Param        course body      createCourseRequest true "Updated Course Information"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /instructor/courses/{id} [put]
// @Security     BearerAuth
// UpdateCourse handles request to edit courses
func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
    // Get instructor ID from context JWT
    instructorID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Could not retrieve instructor ID", http.StatusInternalServerError)
        return
    }

    // Get course id from url parameter
    courseID := chi.URLParam(r, "id")

    // Check if the course exists and belongs to the instructor
    existingCourse, err := h.Repo.GetCourseByID(courseID)
    if err != nil || existingCourse == nil {
        http.Error(w, "Course not found", http.StatusNotFound)
        return
    }
    if existingCourse.InstructorID != instructorID {
        http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
        return
    }

    // Parse the request body into a Course struct
    var courseUpdates model.Course
    if err := json.NewDecoder(r.Body).Decode(&courseUpdates); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate the course fields
    courseUpdates.ID = courseID

    // Update the course in the database
    if err := h.Repo.UpdateCourse(&courseUpdates); err != nil {
        http.Error(w, "Failed to update course", http.StatusInternalServerError)
        return
    }

    // Respond with success message
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Course updated successfully"})
}

// @Summary      Add material to a course (Instructor only)
// @Description  Adds a new learning material to a specific course.
// @Tags         Instructor - Materials
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Course ID"
// @Param        material body      addMaterialRequest true "Material Information"
// @Success      201      {object}  model.LearningMaterial
// @Failure      400      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /instructor/courses/{id}/materials [post]
// @Security     BearerAuth
// AddMaterialToCourse handles request to add material to a course
func (h *CourseHandler) AddMaterialToCourse(w http.ResponseWriter, r *http.Request) {
    instructorID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Could not retrieve instructor ID", http.StatusInternalServerError)
        return
    }

    courseID := chi.URLParam(r, "id")

    // Verify course ownership before adding material
    existingCourse, err := h.Repo.GetCourseByID(courseID)
    if err != nil || existingCourse == nil {
        http.Error(w, "Course not found", http.StatusNotFound)
        return
    }
    if existingCourse.InstructorID != instructorID {
        http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
        return
    }

    var material model.LearningMaterial
    if err := json.NewDecoder(r.Body).Decode(&material); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Simple input validation
    if strings.TrimSpace(material.Title) == "" || strings.TrimSpace(material.ContentType) == "" {
        http.Error(w, "Title and content_type are required", http.StatusBadRequest)
        return
    }

    material.CourseID = courseID // Set course id from URL

    if err := h.Repo.AddMaterialToCourse(&material); err != nil {
        http.Error(w, "Failed to add material", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(material)
}

// @Summary      Get course materials (Instructor only)
// @Description  Retrieves all learning materials for a specific course.
// @Tags         Instructor - Materials
// @Produce      json
// @Param        id   path      string  true  "Course ID"
// @Success      200  {array}   model.LearningMaterial
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /instructor/courses/{id}/materials [get]
// @Security     BearerAuth
// GetMaterialsByCourseID handles request to retrieve materials of a course
func (h *CourseHandler) GetMaterialsByCourseID(w http.ResponseWriter, r *http.Request) {
    instructorID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Could not retrieve instructor ID", http.StatusInternalServerError)
        return
    }

    courseId := chi.URLParam(r, "id")

    // Verify course ownership before retrieving materials
    existingCourse, err := h.Repo.GetCourseByID(courseId)
    if err != nil || existingCourse == nil {
        http.Error(w, "Course not found", http.StatusNotFound)
        return
    }
    if existingCourse.InstructorID != instructorID {
        http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
        return
    }

    // Take material from repository
    materials, err := h.Repo.GetMaterialsByCourseID(courseId)
    if err != nil {
        http.Error(w, "Failed to fetch materials", http.StatusInternalServerError)
        return
    }

    // Respond with the list of materials
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(materials)
}

// @Summary      Update course material (Instructor only)
// @Description  Updates a specific learning material within a course.
// @Tags         Instructor - Materials
// @Accept       json
// @Produce      json
// @Param        id         path      string  true  "Course ID"
// @Param        materialId path      string  true  "Material ID"
// @Param        material   body      addMaterialRequest true "Updated Material Information"
// @Success      200        {object}  map[string]string
// @Failure      400        {object}  map[string]string
// @Failure      403        {object}  map[string]string
// @Failure      404        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /instructor/courses/{id}/materials/{materialId} [put]
// @Security     BearerAuth
// UpdateMaterial handles requests to edit course materials
func (h *CourseHandler) UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	instructorID, _ := r.Context().Value(middleware.UserIDKey).(string)
	courseID := chi.URLParam(r, "id")
	materialID := chi.URLParam(r, "materialId")

	// Verify course ownership
	existingCourse, err := h.Repo.GetCourseByID(courseID)
	if err != nil || existingCourse == nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}
	if existingCourse.InstructorID != instructorID {
		http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
		return
	}

	// Decode data update from body
	var materialUpdates model.LearningMaterial
	if err := json.NewDecoder(r.Body).Decode(&materialUpdates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set ID from URL
	materialUpdates.ID = materialID
	materialUpdates.CourseID = courseID

	// Call repository to update
	if err := h.Repo.UpdateMaterial(&materialUpdates); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Material not found in this course", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update material", http.StatusInternalServerError)
		return
	}

    // Respond with success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Material updated successfully"})
}

// @Summary      Delete course material (Instructor only)
// @Description  Deletes a specific learning material from a course.
// @Tags         Instructor - Materials
// @Produce      json
// @Param        id         path      string  true  "Course ID"
// @Param        materialId path      string  true  "Material ID"
// @Success      200        {object}  map[string]string
// @Failure      403        {object}  map[string]string
// @Failure      404        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /instructor/courses/{id}/materials/{materialId} [delete]
// @Security     BearerAuth
// DeleteMaterial handles requests to delete course materials
func (h *CourseHandler) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
    // Get instructor ID from context JWT
	instructorID, _ := r.Context().Value(middleware.UserIDKey).(string)
	courseID := chi.URLParam(r, "id")
	materialID := chi.URLParam(r, "materialId")

	// Verify course ownership
	existingCourse, err := h.Repo.GetCourseByID(courseID)
	if err != nil || existingCourse == nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}
	if existingCourse.InstructorID != instructorID {
		http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
		return
	}

    // Call repository to delete
	if err := h.Repo.DeleteMaterial(courseID, materialID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Material not found in this course", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete material", http.StatusInternalServerError)
		return
	}

    // Respond with success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Material deleted successfully"})
}

// @Summary      Get public course catalog
// @Description  Retrieves a list of all available courses for anyone to see.
// @Tags         Public
// @Produce      json
// @Success      200  {array}   model.Course
// @Failure      500  {object}  map[string]string
// @Router       /courses [get]
// GetAllCoursesPublic handles requests to retrieve public course catalog
func (h *CourseHandler) GetAllCoursesPublic(w http.ResponseWriter, r *http.Request) {
    courses, err := h.Repo.GetAllCourses()
    if err != nil {
        http.Error(w, "Could not fetch courses", http.StatusInternalServerError)
        return
    }

    // Respond with the list of courses
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(courses)
}

// @Summary      Enroll in a course (Student only)
// @Description  Enrolls the currently logged-in student into a specific course.
// @Tags         Student
// @Produce      json
// @Param        id   path      string  true  "Course ID"
// @Success      201  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string "Course not found (due to foreign key constraint)"
// @Failure      409  {object}  map[string]string "Student is already enrolled in this course"
// @Failure      500  {object}  map[string]string
// @Router       /courses/{id}/enroll [post]
// @Security     BearerAuth
// EnrollInCourse handles requests to enroll students in courses
func (h *CourseHandler) EnrollInCourse(w http.ResponseWriter, r *http.Request) {
    // Get student ID from the JWT context
    studentID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Could not retrieve student ID from context", http.StatusInternalServerError)
        return
    }

    // Get course id from url parameter
    courseID := chi.URLParam(r, "id")

    // Call repository to register students
    err := h.Repo.EnrollStudent(studentID, courseID)
    if err != nil {
        // Check if the error is caused by duplication (unique constraint violation)
        // Code '23505' is the standard PostgreSQL error code for this.
        if strings.Contains(err.Error(), "23505") {
            http.Error(w, "You are already enrolled in this course", http.StatusConflict) // 409 Conflict
            return
        }
        http.Error(w, "Failed to enroll in course", http.StatusInternalServerError)
        return
    }

    // Respond with success message
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Successfully enrolled in the course"})
}

// @Summary      Get my enrolled courses (Student only)
// @Description  Retrieves a list of all courses the logged-in student is enrolled in.
// @Tags         Student
// @Produce      json
// @Success      200  {array}   model.Course
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /student/my-courses [get]
// @Security     BearerAuth
// GetMyEnrolledCourses handles requests to retrieve enrolled courses of a student
func (h *CourseHandler) GetMyEnrolledCourses(w http.ResponseWriter, r *http.Request) {
    // Get the student ID from the JWT context
    studentID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Could not retrieve student ID from context", http.StatusInternalServerError)
        return
    }

    // Fetch enrolled courses from the repository
    courses, err := h.Repo.GetEnrolledCoursesByStudentID(studentID)
    if err != nil {
        http.Error(w, "Failed to fetch enrolled courses", http.StatusInternalServerError)
        return
    }

    // Respond with the list of courses
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(courses)
}

// @Summary      Get enrolled course details (Student only)
// @Description  Retrieves details and all materials for a specific course the student is enrolled in.
// @Tags         Student
// @Produce      json
// @Param        id   path      string  true  "Course ID"
// @Success      200  {object}  handler.courseWithMaterials
// @Failure      403  {object}  map[string]string "Returned if the student is not enrolled in the course"
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /student/courses/{id} [get]
// @Security     BearerAuth
// GetEnrolledCourseDetails handles requests to retrieve enrolled course details
func (h *CourseHandler) GetEnrolledCourseDetails(w http.ResponseWriter, r *http.Request) {
    // Get the student ID from the JWT context
    studentID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Could not retrieve student ID", http.StatusInternalServerError)
        return
    }

    courseID := chi.URLParam(r, "id")

    // Verify enrollment
    isEnrolled, err := h.Repo.IsStudentEnrolled(studentID, courseID)
    if err != nil {
        http.Error(w, "Failed to verify enrollment", http.StatusInternalServerError)
        return
    }
    if !isEnrolled {
        http.Error(w, "Forbidden: You are not enrolled in this course", http.StatusForbidden)
        return
    }

    // Get course details and materials
    course, err := h.Repo.GetCourseByID(courseID)
    if err != nil || course == nil {
        http.Error(w, "Course not found", http.StatusNotFound)
        return
    }

    materials, err := h.Repo.GetMaterialsByCourseID(courseID)
    if err != nil {
        http.Error(w, "Failed to fetch materials", http.StatusInternalServerError)
        return
    }

    // Combine into one response
    response := courseWithMaterials{
        Course:    *course,
        Materials: materials,
    }

    // Respond with the course details and materials
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}