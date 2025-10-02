package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/dimasrizkyfebrian/coursify/internal/database"
	"github.com/dimasrizkyfebrian/coursify/internal/handler"
	"github.com/dimasrizkyfebrian/coursify/internal/handler/middleware"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.ConnectDB()

	// Inisialisasi router chi
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger) // Middleware untuk logging request

	// Membuat instance dari repository dan handler
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// --- Public Routes ---
	r.Post("/api/register", userHandler.Register)
	r.Post("/api/login", userHandler.Login)

	// --- Protected Routes ---
	r.Group(func(r chi.Router) {
		// Terapkan middleware AuthMiddleware
		r.Use(middleware.AuthMiddleware)

		// Semua route di dalam grup ini akan terproteksi
		r.Get("/api/profile", userHandler.GetProfile)
	})

	port := ":8080"
	log.Printf("Server is starting on port %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}