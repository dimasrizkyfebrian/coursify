package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	db := database.ConnectDB() // Database connection

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Allow port 5173
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)
	courseRepo := repository.NewCourseRepository(db)
	courseHandler := handler.NewCourseHandler(courseRepo)


	// --- Public Routes ---
	r.With(middleware.RateLimitMiddleware).Post("/api/register", userHandler.Register)
	r.Post("/api/login", userHandler.Login)

	// --- Protected Admin Routes ---
	r.Group(func(r chi.Router) {
	r.Use(middleware.AuthMiddleware)
	r.Use(middleware.AdminOnly)

	r.Get("/api/admin/users/stats", userHandler.GetUserStats)
	r.Get("/api/admin/users/pending", userHandler.GetPendingUsers)
	r.Get("/api/admin/users/pending/count", userHandler.GetPendingUserCount)
	r.Get("/api/admin/users/all", userHandler.GetAllUsers)
	r.Get("/api/admin/users/{id}", userHandler.GetUserByIDForAdmin)
	r.Put("/api/admin/users/{id}/approve", userHandler.ApproveUser)
	r.Put("/api/admin/users/{id}/reject", userHandler.RejectUser)
	r.Put("/api/admin/users/{id}", userHandler.UpdateUser)
	r.Delete("/api/admin/users/{id}", userHandler.DeleteUser)
	})

	// --- Protected Instructor Routes ---
	r.Group(func(r chi.Router) {
    r.Use(middleware.AuthMiddleware)
    r.Use(middleware.InstructorOnly)

	r.Get("/api/instructor/courses", courseHandler.GetMyCourses)
    r.Post("/api/instructor/courses", courseHandler.CreateCourse)
	})
	
	// --- Protected General Routes ---
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/api/profile", userHandler.GetProfile)
	})

	port := ":8080"
	log.Printf("Server is starting on port %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}