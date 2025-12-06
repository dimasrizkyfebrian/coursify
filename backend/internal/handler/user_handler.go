package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dimasrizkyfebrian/coursify/internal/dto"
	"github.com/dimasrizkyfebrian/coursify/internal/handler/middleware"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/dimasrizkyfebrian/coursify/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
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

	if err := h.Service.Register(&user); err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully, waiting for admin approval"})
}

// @Summary      Log in a user
// @Description  Authenticates a user and returns a JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body dto.LoginRequest true "User credentials"
// @Success      200  {object}  map[string]string "{"token": "JWT_TOKEN"}"
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.Service.Login(credentials.Email, credentials.Password)
	if err != nil {
		if err.Error() == "invalid email or password" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == "account is not active" {
			http.Error(w, "Account is not active, please wait for admin approval", http.StatusForbidden)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// @Summary      Get user profile
// @Description  Retrieves the profile information for the currently logged-in user.
// @Tags         Users
// @Produce      json
// @Success      200  {object}  dto.UserResponse
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

	user, err := h.Service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Could not fetch user profile", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := h.Service.MapUserToResponse(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary      Get pending users (Admin only)
// @Description  Retrieves a list of users with 'pending' status.
// @Tags         Admin
// @Produce      json
// @Success      200  {array}  dto.UserResponse
// @Failure      403  {object}  map[string]string
// @Router       /admin/users/pending [get]
// @Security     BearerAuth
func (h *UserHandler) GetPendingUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetPendingUsers()
	if err != nil {
		http.Error(w, "Could not fetch users", http.StatusInternalServerError)
		return
	}

	response := h.Service.MapUsersToResponse(users)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

	err := h.Service.ApproveUser(userID)
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

	err := h.Service.RejectUser(userID)
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
// @Success      200  {object}  dto.UserResponse
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/users/{id} [get]
// @Security     BearerAuth
func (h *UserHandler) GetUserByIDForAdmin(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	user, err := h.Service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Could not fetch user profile", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := h.Service.MapUserToResponse(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
	count, err := h.Service.GetPendingUserCount()
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
// @Success      200  {array}   dto.UserResponse
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/users/all [get]
// @Security     BearerAuth
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		http.Error(w, "Could not fetch users", http.StatusInternalServerError)
		return
	}

	response := h.Service.MapUsersToResponse(users)

	w.Header().Set("Content_Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary      Update a user (Admin only)
// @Description  Updates a user's full_name, email, or role.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Param        user body      dto.UpdateUserRequest true "User data to update"
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

	err := h.Service.UpdateUser(&userUpdates)
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

	err := h.Service.DeleteUser(userID)
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
	stats, err := h.Service.GetUserStats()
	if err != nil {
		http.Error(w, "Could not fetch user statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

// @Summary      Upload or update user avatar
// @Description  Uploads a new avatar image (jpg, png, max 2MB) for the logged-in user. Converts to WebP, deletes old avatar.
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

	// Get the file from form data
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

	// --- Decode Uploaded Image ---
	img, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Error decoding image for user %s: %v", userID, err)
		http.Error(w, "Failed to decode image. Ensure it is a valid JPG or PNG.", http.StatusBadRequest)
		return
	}
	log.Printf("Decoded image format: %s for user %s", format, userID)
	if format != "jpeg" && format != "png" {
		http.Error(w, "Invalid file format after decoding. Only JPG and PNG are allowed.", http.StatusBadRequest)
		return
	}
	// --- End Image Decoding ---

	// --- Encode Image as WebP using kolesa-team/go-webp ---
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 80)
	if err != nil {
		log.Printf("Error creating webp encoder options: %v", err)
		http.Error(w, "Failed to configure image conversion", http.StatusInternalServerError)
		return
	}

	// Create a buffer to store the encoding results
	var webpBuffer bytes.Buffer

	// Perform encoding from 'img' to 'webpBuffer' with the specified options
	if err := webp.Encode(&webpBuffer, img, options); err != nil {
		log.Printf("Error encoding image to WebP for user %s: %v", userID, err)
		http.Error(w, "Failed to convert image to WebP", http.StatusInternalServerError)
		return
	}
	// --- End WebP Encoding ---

	// Generate unique filename and paths
	fileName := fmt.Sprintf("%s-%d.webp", userID, time.Now().UnixNano())
	uploadDir := filepath.Join("uploads", "avatars")
	filePath := filepath.Join(uploadDir, fileName)
	fileURL := "/" + strings.ReplaceAll(filePath, "\\", "/") // URL path for DB

	// Ensure directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Error creating avatar directory %s: %v", uploadDir, err)
		http.Error(w, "Could not create upload directory", http.StatusInternalServerError)
		return
	}

	// --- Delete Old Avatar File ---
	currentUser, err := h.Service.GetUserByID(userID)
	if err != nil || currentUser == nil {
		log.Printf("Warning: Error fetching user %s to delete old avatar: %v", userID, err)
	} else if currentUser.AvatarURL.Valid && currentUser.AvatarURL.String != "" {
		oldFilePath := filepath.Join(".", strings.TrimPrefix(currentUser.AvatarURL.String, "/"))
		oldFilePath = filepath.FromSlash(oldFilePath)
		errRemove := os.Remove(oldFilePath)
		if errRemove != nil && !os.IsNotExist(errRemove) {
			log.Printf("Warning: Failed to delete old avatar file %s for user %s: %v", oldFilePath, userID, errRemove)
		} else if errRemove == nil {
			log.Printf("Successfully deleted old avatar file %s for user %s", oldFilePath, userID)
		}
	}
	// --- End Delete Old Avatar ---

	// --- Save the NEW WebP file ---
	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file %s: %v", filePath, err)
		http.Error(w, "Could not save the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy WebP data from buffer to file
	if _, err := io.Copy(dst, &webpBuffer); err != nil {
		log.Printf("Error writing WebP data to file %s: %v", filePath, err)
		http.Error(w, "Could not write WebP file", http.StatusInternalServerError)
		os.Remove(filePath) // Cleanup file if copying fails
		return
	}
	log.Printf("Successfully saved NEW WebP avatar to: %s", filePath)
	// --- End Save New File ---

	// Update database with the new file URL
	err = h.Service.UpdateUserAvatarURL(userID, fileURL)
	if err != nil {
		os.Remove(filePath) // Clean up file if DB update fails
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
