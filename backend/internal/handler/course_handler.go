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

    var course model.Course
    if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Set the instructor ID from the token, not from the request body
    course.InstructorID = instructorID

    if err := h.Repo.CreateCourse(&course); err != nil {
        http.Error(w, "Failed to create course", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(course)
}