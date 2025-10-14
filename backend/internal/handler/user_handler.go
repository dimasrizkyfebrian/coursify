package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
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