package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dimasrizkyfebrian/coursify/internal/handler/middleware"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

// Register adalah method untuk menangani request registrasi
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	// Decode JSON body dari request ke dalam struct User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Panggil repository untuk membuat user
	if err := h.Repo.CreateUser(&user); err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	// Kirim response sukses
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully, waiting for admin approval"})
}

// Login adalah method untuk menangani request login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 1. Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Ambil user berdasarkan email dari repository
	user, err := h.Repo.GetUserByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 3. Bandingkan password yang di-request dengan hash di DB
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 4. Cek apakah status user sudah 'active'
	if user.Status != "active" {
		http.Error(w, "Account is not active, please wait for admin approval", http.StatusForbidden)
		return
	}

	// 5. Jika semua benar, buat JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token berlaku 3 hari
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// 6. Kirim token sebagai response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// GetProfile adalah handler untuk mengambil data profil pengguna (contoh endpoint terproteksi)
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Ambil user_id dari context yang sudah ditambahkan oleh middleware
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	// (Untuk sekarang kita hanya kembalikan ID-nya. Nanti kita bisa query ke DB)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to your profile!",
		"user_id": userID,
	})
}