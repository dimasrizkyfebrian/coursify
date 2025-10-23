package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dimasrizkyfebrian/coursify/internal/handler/middleware"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

// @Summary      Register a new user
// @Description  Creates a new user account with a 'pending' status.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body model.User true "User registration info"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Repo.CreateUser(&user); err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully, waiting for admin approval"})
}

type loginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// @Summary      Log in a user
// @Description  Authenticates a user and returns a JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body loginRequest true "User credentials"
// @Success      200  {object}  map[string]string "{"token": "JWT_TOKEN"}"
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetUserByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if user.Status != "active" {
		http.Error(w, "Account is not active, please wait for admin approval", http.StatusForbidden)
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// @Summary      Get user profile
// @Description  Retrieves the profile information for the currently logged-in user.
// @Tags         Users
// @Produce      json
// @Success      200  {object}  model.User
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /profile [get]
// @Security     BearerAuth
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	user, err := h.Repo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Could not fetch user profile", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// @Summary      Get pending users (Admin only)
// @Description  Retrieves a list of users with 'pending' status.
// @Tags         Admin
// @Produce      json
// @Success      200  {array}  model.User
// @Failure      403  {object}  map[string]string
// @Router       /admin/users/pending [get]
// @Security     BearerAuth
func (h *UserHandler) GetPendingUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Repo.GetUsersByStatus("pending")
	if err != nil {
		http.Error(w, "Could not fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// @Summary      Approve a user (Admin only)
// @Description  Changes a user's status from 'pending' to 'active'.
// @Tags         Admin
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/users/{id}/approve [put]
// @Security     BearerAuth
func (h *UserHandler) ApproveUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	err := h.Repo.UpdateUserStatus(userID, "active")
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to approve user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User approved successfully"})
}

// @Summary      Reject a user (Admin only)
// @Description  Changes a user's status from 'pending' to 'rejected'.
// @Tags         Admin
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/users/{id}/reject [put]
// @Security     BearerAuth
func (h *UserHandler) RejectUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	
	err := h.Repo.UpdateUserStatus(userID, "rejected")
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to reject user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User rejected successfully"})
}

// @Summary      Get a single user's details (Admin only)
// @Description  Retrieves the full details of a single user by their ID.
// @Tags         Admin
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  model.User
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/users/{id} [get]
// @Security     BearerAuth
func (h *UserHandler) GetUserByIDForAdmin(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	user, err := h.Repo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Could not fetch user profile", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// @Summary      Get pending user count (Admin only)
// @Description  Retrieves the number of users with 'pending' status.
// @Tags         Admin
// @Produce      json
// @Success      200  {object}  map[string]int "{"count": 5}"
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/users/pending/count [get]
// @Security     BearerAuth
func (h *UserHandler) GetPendingUserCount(w http.ResponseWriter, r *http.Request) {
    count, err := h.Repo.GetPendingUserCount()
    if err != nil {
        http.Error(w, "Could not get pending user count", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]int{"count": count})
}

// @Summary      Get all users (Admin only)
// @Description  Retrieves a list of all users regardless of their status.
// @Tags         Admin
// @Produce      json
// @Success      200  {array}   model.User
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/users/all [get]
// @Security     BearerAuth
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Repo.GetAllUsers()
	if err != nil {
		http.Error(w, "Could not fetch users", http.StatusInternalServerError)
	return
	}

	w.Header().Set("Content_Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

type updateUserRequest struct {
	FullName string `json:"full_name" example:"John Doe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Role     string `json:"role" enums:"admin,instructor,student"`
}

// @Summary      Update a user (Admin only)
// @Description  Updates a user's full_name, email, or role.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Param        user body      updateUserRequest true "User data to update"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/users/{id} [put]
// @Security     BearerAuth
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	var userUpdates model.User
	if err := json.NewDecoder(r.Body).Decode(&userUpdates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userUpdates.ID = userID

	err := h.Repo.UpdateUser(&userUpdates)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

// @Summary      Delete a user (Admin only)
// @Description  Permanently deletes a user account.
// @Tags         Admin
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/users/{id} [delete]
// @Security     BearerAuth
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	err := h.Repo.DeleteUser(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

// @Summary      Get user statistics (Admin only)
// @Description  Retrieves key statistics like total, active, and pending users.
// @Tags         Admin
// @Produce      json
// @Success      200  {object}  map[string]int "{"total_users": 10, "active_users": 5, "pending_users": 2}"
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/users/stats [get]
// @Security     BearerAuth
func (h *UserHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
    stats, err := h.Repo.GetUserStats()
    if err != nil {
        http.Error(w, "Could not fetch user statistics", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(stats)
}

// @Summary      Upload or update user avatar
// @Description  Uploads a new avatar image (jpg, png) for the logged-in user. Converts to WebP.
// @Tags         Profile
// @Accept       multipart/form-data
// @Produce      json
// @Param        avatar formData file true "Avatar image file (jpg, png, max 2MB)"
// @Success      200 {object} map[string]string "{"message": "Avatar updated successfully", "url": "/uploads/avatars/..."}"
// @Failure      400 {object} map[string]string "e.g., No file, file too large, invalid file type"
// @Failure      500 {object} map[string]string
// @Router       /profile/avatar [put]
// @Security     BearerAuth
// UploadAvatar handles avatar image uploads
func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	// Get User ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	// Parse Multipart Form (e.g., max 2MB)
	maxUploadSize := int64(2 * 1024 * 1024) // 2 MB
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, fmt.Sprintf("File is too large. Max size is %dMB.", maxUploadSize/(1024*1024)), http.StatusBadRequest)
		return
	}

	// Get the file from form data (key should be 'avatar')
	file, _, err := r.FormFile("avatar")
	if err != nil {
		if err == http.ErrMissingFile {
			http.Error(w, "No file uploaded. Please use 'avatar' as the key.", http.StatusBadRequest)
		} else {
			http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	// Validate File Type (MIME Type Check is more reliable)
	// Read the first 512 bytes to determine the actual file type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(w, "Error reading file for type validation", http.StatusInternalServerError)
		return
	}
	// Reset the read pointer back to the beginning of the file
	file.Seek(0, 0)

	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" {
		http.Error(w, "Invalid file type. Only JPG and PNG are allowed.", http.StatusBadRequest)
		return
	}

	// --- TODO: Image Conversion (using an imaging library) ---
	// - Decode the uploaded image (PNG or JPG) from 'file'
	// - Encode the image as WebP into a buffer or temporary file
	// - For now, we'll just pretend this happened and continue with file saving logic

	// Generate unique filename (e.g., user_id-timestamp.webp)
	// Use timestamp for uniqueness, extension is always .webp
	fileName := fmt.Sprintf("%s-%d.webp", userID, time.Now().UnixNano())
	uploadDir := filepath.Join("uploads", "avatars") // Subdirectory for avatars
	filePath := filepath.Join(uploadDir, fileName)
	fileURL := "/" + strings.ReplaceAll(filePath, "\\", "/") // Ensure URL uses forward slashes

	// Ensure directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Error creating avatar directory %s: %v", uploadDir, err)
		http.Error(w, "Could not create upload directory", http.StatusInternalServerError)
		return
	}

	// --- TODO: Delete Old Avatar File ---
	// - Get the current avatar_url from the database for the user
	// - If it exists, construct the old file path on the server
	// - Use os.Remove() to delete the old file before saving the new one

	// Save the (converted) file
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file %s: %v", filePath, err)
		http.Error(w, "Could not save the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// In a real scenario, you would copy the *converted WebP data* here
	// For now, we copy the original file just to test the flow
	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Error copying file content to %s: %v", filePath, err)
		http.Error(w, "Could not copy file content", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully saved (placeholder) avatar to: %s", filePath) // Add log

	// Update database with the new file URL
	err = h.Repo.UpdateUserAvatarURL(userID, fileURL)
	if err != nil {
		// Attempt to clean up the newly saved file if DB update fails
		os.Remove(filePath)
		log.Printf("Error updating avatar URL in DB for user %s: %v", userID, err)
		http.Error(w, "Could not update avatar information", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Avatar updated successfully",
		"url":     fileURL,
	})
}