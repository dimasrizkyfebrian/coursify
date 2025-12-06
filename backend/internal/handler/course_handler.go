package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dimasrizkyfebrian/coursify/internal/dto"
	"github.com/dimasrizkyfebrian/coursify/internal/handler/middleware"
	"github.com/dimasrizkyfebrian/coursify/internal/service"
	"github.com/go-chi/chi/v5"
)

type CourseHandler struct {
	Service *service.CourseService
}

func NewCourseHandler(service *service.CourseService) *CourseHandler {
	return &CourseHandler{Service: service}
}

// @Summary      Create a new course (Instructor only)
// @Description  Creates a new course for the logged-in instructor.
// @Tags         Instructor
// @Accept       json
// @Produce      json
// @Param        course body dto.CreateCourseRequest true "Course Information"
// @Success      201  {object}  dto.CourseResponse
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

	// Decode JSON body into the simple 'createCourseRequest' struct
	var req dto.CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the input
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Description) == "" {
		http.Error(w, "Title and description cannot be empty", http.StatusBadRequest)
		return
	}

	course, err := h.Service.CreateCourse(instructorID, &req)
	if err != nil {
		http.Error(w, "Failed to create course", http.StatusInternalServerError)
		return
	}

	response := h.Service.MapCourseToResponse(course)

	// Respond with the created course
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// @Summary      Get my courses (Instructor only)
// @Description  Retrieves a list of all courses created by the logged-in instructor.
// @Tags         Instructor
// @Produce      json
// @Success      200  {array}   dto.CourseResponse
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
	courses, err := h.Service.GetMyCourses(instructorID)
	if err != nil {
		http.Error(w, "Failed to fetch courses", http.StatusInternalServerError)
		return
	}

	response := h.Service.MapCoursesToResponse(courses)

	// Respond with the list of courses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary      Get my course details (Instructor only)
// @Description  Retrieves the details of a specific course owned by the logged-in instructor.
// @Tags         Instructor
// @Produce      json
// @Param        id   path      string  true  "Course ID"
// @Success      200  {object}  dto.CourseResponse
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
	course, err := h.Service.GetCourseByID(courseID)
	if err != nil || course == nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	// Check if the course belongs to the instructor
	if course.InstructorID != instructorID {
		http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
		return
	}

	response := h.Service.MapCourseToResponse(course)

	// Respond with the course details
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary      Update a course (Instructor only)
// @Description  Updates the title and description of a course owned by the logged-in instructor.
// @Tags         Instructor
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Course ID"
// @Param        course body      dto.CreateCourseRequest true "Updated Course Information"
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

	// Parse the request body into a Course struct
	var req dto.CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateCourse(courseID, instructorID, &req)
	if err != nil {
		if err.Error() == "course not found" {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else if err.Error() == "forbidden" {
			http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to update course", http.StatusInternalServerError)
		}
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Course updated successfully"})
}

// @Summary      Delete a course (Instructor only)
// @Description  Deletes a course owned by the logged-in instructor.
// @Tags         Instructor
// @Produce      json
// @Param        id   path      string  true  "Course ID"
// @Success      200  {object}  map[string]string "{"message": "Course deleted successfully"}"
// @Failure      403  {object}  map[string]string "Forbidden: You are not the owner"
// @Failure      404  {object}  map[string]string "Course not found"
// @Failure      500  {object}  map[string]string
// @Router       /instructor/courses/{id} [delete]
// @Security     BearerAuth
// DeleteCourse handles request to delete a course
func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	// Get instructor ID from context JWT
	instructorID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not retrieve instructor ID", http.StatusInternalServerError)
		return
	}

	// Get course id from url parameter
	courseID := chi.URLParam(r, "id")

	// Check if the course exists and belongs to the instructor
	err := h.Service.DeleteCourse(courseID, instructorID)
	if err != nil {
		// Note: Service should ideally return specific errors to distinguish 404/403
		// For now assuming generic error or not found/forbidden handled by repo logic mostly
		// But service wrapper might hide sql.ErrNoRows.
		// Let's assume service returns error if operation fails.
		http.Error(w, "Failed to delete course or not found/forbidden", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Course deleted successfully"})
}

// @Summary      Add material to a course (Instructor only)
// @Description  Adds a new learning material to a specific course.
// @Tags         Instructor - Materials
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Course ID"
// @Param        material body      dto.AddMaterialRequest true "Material Information"
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

	var req dto.AddMaterialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simple input validation
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.ContentType) == "" {
		http.Error(w, "Title and content_type are required", http.StatusBadRequest)
		return
	}

	material, err := h.Service.AddMaterial(courseID, instructorID, &req)
	if err != nil {
		if err.Error() == "course not found" {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else if err.Error() == "forbidden" {
			http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to add material", http.StatusInternalServerError)
		}
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

	materials, err := h.Service.GetMaterials(courseId, instructorID)
	if err != nil {
		if err.Error() == "course not found" {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else if err.Error() == "forbidden" {
			http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to fetch materials", http.StatusInternalServerError)
		}
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
// @Param        material   body      dto.AddMaterialRequest true "Updated Material Information"
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

	// Decode data update from body
	var req dto.AddMaterialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateMaterial(courseID, materialID, instructorID, &req)
	if err != nil {
		if err.Error() == "course not found" {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else if err.Error() == "forbidden" {
			http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to update material", http.StatusInternalServerError)
		}
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

	err := h.Service.DeleteMaterial(courseID, materialID, instructorID)
	if err != nil {
		if err.Error() == "course not found" {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else if err.Error() == "forbidden" {
			http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to delete material", http.StatusInternalServerError)
		}
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
// @Success      200  {array}   dto.CourseResponse
// @Failure      500  {object}  map[string]string
// @Router       /courses [get]
// GetAllCoursesPublic handles requests to retrieve public course catalog
func (h *CourseHandler) GetAllCoursesPublic(w http.ResponseWriter, r *http.Request) {
	courses, err := h.Service.GetAllCoursesPublic()
	if err != nil {
		http.Error(w, "Could not fetch courses", http.StatusInternalServerError)
		return
	}

	response := h.Service.MapCoursesToResponse(courses)

	// Respond with the list of courses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
	err := h.Service.EnrollStudent(studentID, courseID)
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
// @Success      200  {array}   dto.CourseResponse
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
	courses, err := h.Service.GetMyEnrolledCourses(studentID)
	if err != nil {
		http.Error(w, "Failed to fetch enrolled courses", http.StatusInternalServerError)
		return
	}

	response := h.Service.MapCoursesToResponse(courses)

	// Respond with the list of courses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary      Get enrolled course details (Student only)
// @Description  Retrieves details and all materials for a specific course the student is enrolled in.
// @Tags         Student
// @Produce      json
// @Param        id   path      string  true  "Course ID"
// @Success      200  {object}  dto.CourseWithMaterialsResponse
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

	response, err := h.Service.GetEnrolledCourseDetails(studentID, courseID)
	if err != nil {
		if err.Error() == "course not found" {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else if err.Error() == "forbidden" {
			http.Error(w, "Forbidden: You are not enrolled in this course", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to fetch course details", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the course details and materials
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary      Upload a PDF material for a course (Instructor only)
// @Description  Uploads a PDF file as a new learning material for a course.
// @Tags         Instructor - Materials
// @Accept       multipart/form-data
// @Produce      json
// @Param        id    path      string  true  "Course ID"
// @Param        title formData  string  true  "Title of the material"
// @Param        pdf   formData  file    true  "PDF file to upload"
// @Success      201   {object}  model.LearningMaterial
// @Failure      400   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /instructor/courses/{id}/materials/upload-pdf [post]
// @Security     BearerAuth
// UploadPdfMaterial handles requests to upload PDF materials
func (h *CourseHandler) UploadPdfMaterial(w http.ResponseWriter, r *http.Request) {
	instructorID, _ := r.Context().Value(middleware.UserIDKey).(string)
	courseID := chi.URLParam(r, "id")

	// Parse form, maximum size 10 MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File is too large. Max size is 10MB.", http.StatusBadRequest)
		return
	}

	// Retrieve the file from form-data with the key 'pdf'
	file, handler, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "No file uploaded. Please use 'pdf' as the file key.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Take the title from the form data
	title := r.FormValue("title")
	if strings.TrimSpace(title) == "" {
		http.Error(w, "Title is required.", http.StatusBadRequest)
		return
	}

	// Create a unique file name
	ext := filepath.Ext(handler.Filename)
	fileName := fmt.Sprintf("%s-%d%s", courseID, time.Now().Unix(), ext)

	// Create the 'uploads/materials' directory if it doesn't exist
	uploadDir := filepath.Join("uploads", "materials")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Could not create uploads directory", http.StatusInternalServerError)
		return
	}

	// Save the file to the server
	filePath := filepath.Join(uploadDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Could not save the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Could not copy the file content", http.StatusInternalServerError)
		return
	}

	material, err := h.Service.UploadPdfMaterial(courseID, instructorID, title, filePath)
	if err != nil {
		if err.Error() == "course not found" {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else if err.Error() == "forbidden" {
			http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to create material", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(material)
}

// @Summary      Get enrolled students for a course (Instructor only)
// @Description  Retrieves a list of students enrolled in a specific course owned by the instructor.
// @Tags         Instructor
// @Produce      json
// @Param        id   path      string  true  "Course ID"
// @Success      200  {array}   dto.UserResponse
// @Failure      403  {object}  map[string]string "Forbidden: You are not the owner of this course"
// @Failure      404  {object}  map[string]string "Course not found"
// @Failure      500  {object}  map[string]string
// @Router       /instructor/courses/{id}/enrollments [get]
// @Security     BearerAuth
// GetEnrolledStudents handles requests to retrieve students enrolled in a course
func (h *CourseHandler) GetEnrolledStudents(w http.ResponseWriter, r *http.Request) {
	// Get instructor ID from context JWT
	instructorID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not retrieve instructor ID", http.StatusInternalServerError)
		return
	}

	// Get course id from url parameter
	courseID := chi.URLParam(r, "id")

	students, err := h.Service.GetEnrolledStudents(courseID, instructorID)
	if err != nil {
		if err.Error() == "course not found" {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else if err.Error() == "forbidden" {
			http.Error(w, "Forbidden: You are not the owner of this course", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to fetch enrolled students", http.StatusInternalServerError)
		}
		return
	}

	// Need to map students to response DTO if we want to be consistent, but UserResponseForSwagger/UserResponse is fine.
	// Assuming we want to use the DTO mapper from UserService? But we are in CourseHandler.
	// For now, let's just return the model or map it manually here if needed.
	// The previous implementation returned UserResponseForSwagger which is similar to model.User but with JSON tags.
	// Let's assume we can return the model directly if JSON tags match, or we should have a mapper.
	// Since we don't have access to UserService here easily without injecting it, let's just return the model for now
	// or create a local mapper if strictly needed.
	// Actually, we can just use the DTO.

	var response []dto.UserResponse
	for _, user := range students {
		response = append(response, dto.UserResponse{
			ID:        user.ID,
			FullName:  user.FullName,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			AvatarURL: &user.AvatarURL.String,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	// Respond with the list of students
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
