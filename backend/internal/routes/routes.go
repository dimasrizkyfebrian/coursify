package routes

import (
	"github.com/dimasrizkyfebrian/coursify/internal/handler"
	"github.com/dimasrizkyfebrian/coursify/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Router struct {
	UserHandler   *handler.UserHandler
	CourseHandler *handler.CourseHandler
	JWTSecret     string
}

func NewRouter(userHandler *handler.UserHandler, courseHandler *handler.CourseHandler, jwtSecret string) *Router {
	return &Router{
		UserHandler:   userHandler,
		CourseHandler: courseHandler,
		JWTSecret:     jwtSecret,
	}
}

func (router *Router) RegisterRoutes(r chi.Router) {
	// --- Swagger Documentation ---
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Point to doc.json
	))

	// --- API v1 ---
	r.Route("/api/v1", func(r chi.Router) {

		// --- Auth Routes ---
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", router.UserHandler.Register)
			r.Post("/auth/login", router.UserHandler.Login)
		})

		// --- Public Course Routes ---
		r.Group(func(r chi.Router) {
			r.Get("/courses", router.CourseHandler.GetAllCoursesPublic)
		})

		// --- Protected Routes (Require Login) ---
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(router.JWTSecret))

			// --- User/Profile Routes ---
			r.Group(func(r chi.Router) {
				r.Get("/users/profile", router.UserHandler.GetProfile)
				r.Put("/users/profile/avatar", router.UserHandler.UploadAvatar)
			})

			// --- Student Routes ---
			r.Group(func(r chi.Router) {
				r.Post("/student/courses/{id}/enroll", router.CourseHandler.EnrollInCourse)
				r.Get("/student/my-courses", router.CourseHandler.GetMyEnrolledCourses)
				r.Get("/student/courses/{id}", router.CourseHandler.GetEnrolledCourseDetails)
			})

			// --- Instructor Routes ---
			r.Group(func(r chi.Router) {
				r.Use(middleware.RoleMiddleware("instructor"))
				r.Post("/instructor/courses", router.CourseHandler.CreateCourse)
				r.Get("/instructor/courses", router.CourseHandler.GetMyCourses)
				r.Get("/instructor/courses/{id}", router.CourseHandler.GetMyCourseDetails)
				r.Put("/instructor/courses/{id}", router.CourseHandler.UpdateCourse)
				r.Delete("/instructor/courses/{id}", router.CourseHandler.DeleteCourse)
				r.Post("/instructor/courses/{id}/materials", router.CourseHandler.AddMaterialToCourse)
				r.Get("/instructor/courses/{id}/materials", router.CourseHandler.GetMaterialsByCourseID)
				r.Put("/instructor/courses/{id}/materials/{materialId}", router.CourseHandler.UpdateMaterial)
				r.Delete("/instructor/courses/{id}/materials/{materialId}", router.CourseHandler.DeleteMaterial)
				r.Post("/instructor/courses/{id}/materials/upload-pdf", router.CourseHandler.UploadPdfMaterial)
				r.Get("/instructor/courses/{id}/enrollments", router.CourseHandler.GetEnrolledStudents)
			})

			// --- Admin Routes ---
			r.Group(func(r chi.Router) {
				r.Use(middleware.RoleMiddleware("admin"))
				r.Get("/admin/users/pending", router.UserHandler.GetPendingUsers)
				r.Put("/admin/users/{id}/approve", router.UserHandler.ApproveUser)
				r.Put("/admin/users/{id}/reject", router.UserHandler.RejectUser)
				r.Get("/admin/users/{id}", router.UserHandler.GetUserByIDForAdmin)
				r.Get("/admin/users/pending/count", router.UserHandler.GetPendingUserCount)
				r.Get("/admin/users/all", router.UserHandler.GetAllUsers)
				r.Put("/admin/users/{id}", router.UserHandler.UpdateUser)
				r.Delete("/admin/users/{id}", router.UserHandler.DeleteUser)
				r.Get("/admin/users/stats", router.UserHandler.GetUserStats)
			})
		})
	})
}
