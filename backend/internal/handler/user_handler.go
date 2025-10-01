package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
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