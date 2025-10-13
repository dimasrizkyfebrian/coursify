package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dimasrizkyfebrian/coursify/internal/handler/middleware"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
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