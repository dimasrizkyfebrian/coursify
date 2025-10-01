package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/dimasrizkyfebrian/coursify/internal/database"
	"github.com/dimasrizkyfebrian/coursify/internal/handler"
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
	r.Use(middleware.Logger) // Middleware untuk logging request

	// Membuat instance dari repository dan handler
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// Mendefinisikan route untuk API
	r.Post("/api/register", userHandler.Register)

	port := ":8080"
	log.Printf("Server is starting on port %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}