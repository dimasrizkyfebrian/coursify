package handler

import (
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