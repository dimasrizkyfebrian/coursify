// backend/cmd/seeder/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"

	"github.com/dimasrizkyfebrian/coursify/internal/database"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
)

// Helper function to create and activate a user
func createAndActivateUser(repo *repository.UserRepository, user *model.User) error {
	// Check if user already exists
	existingUser, err := repo.GetUserByEmail(user.Email)
	if err != nil && err != sql.ErrNoRows { // Handle potential DB errors, ignore ErrNoRows
		log.Printf("Error checking for user %s: %v\n", user.Email, err)
		return err // Return the error to stop processing this user
	}
	if existingUser != nil {
		fmt.Printf("User %s already exists, skipping.\n", user.Email)
		// Optionally update status just in case it wasn't active
		if existingUser.Status != "active" {
			errApprove := repo.UpdateUserStatus(existingUser.ID, "active")
			if errApprove != nil {
				log.Printf("Could not activate existing user %s: %v\n", user.Email, errApprove)
			} else {
				fmt.Printf("Activated existing user %s\n", user.Email)
			}
		}
		return nil // Successfully skipped or updated
	}

	// Create the user
	err = repo.CreateUser(user)
	if err != nil {
		log.Printf("Could not create user %s: %v\n", user.Email, err)
		return err // Return the error
	}
	fmt.Printf("Created default %s: %s\n", user.Role, user.Email)

	// Activate the newly created user
	if user.ID == "" {
		log.Printf("Warning: User ID not populated after creation for %s. Cannot activate.", user.Email)
		return fmt.Errorf("user ID not populated after creation for %s", user.Email) // Indicate failure
	}
	errApprove := repo.UpdateUserStatus(user.ID, "active")
	if errApprove != nil {
		log.Printf("Could not activate newly created user %s: %v\n", user.Email, errApprove)
		// Even if activation fails, we consider creation successful for seeding purposes
	} else {
		fmt.Printf("Activated newly created user %s\n", user.Email)
	}
	return nil // User created (and activation attempted)
}

func main() {
	// Load env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or error loading, relying on environment variables.")
	}

	db := database.ConnectDB()
	userRepo := repository.NewUserRepository(db)

	fmt.Println("Seeding default users...")

	// --- 1. Define Default Users ---
	defaultUsers := []model.User{
		{
			FullName: "Admin User",
			Email:    "admin@example.com",
			Password: "password123",
			Role:     "admin",
		},
		{
			FullName: "Instructor User",
			Email:    "instructor@example.com",
			Password: "password123",
			Role:     "instructor",
		},
		{
			FullName: "Student User",
			Email:    "student@example.com",
			Password: "password123",
			Role:     "student",
		},
	}

	// --- 2. Create and Activate Default Users ---
	for i := range defaultUsers {
		// Pass the address of the user in the slice
		err := createAndActivateUser(userRepo, &defaultUsers[i])
		if err != nil {
			// Log the error but continue with other default users
			log.Printf("Failed to process default user %s: %v", defaultUsers[i].Email, err)
		}
	}

	// --- 3. Create Random Users ---
	fmt.Println("\nSeeding random users...")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	roles := []string{"student", "instructor"}

	for i := 0; i < 7; i++ { // Create 7 random users
		randomRole := roles[r.Intn(len(roles))]
		user := model.User{
			FullName: faker.Name(),
			Email:    faker.Email(),
			Password: "password123",
			Role:     randomRole,
		}
		// Don't check for existing random users, just create
		err := userRepo.CreateUser(&user)
		if err != nil {
			log.Printf("Could not create random user %s: %v\n", user.Email, err)
		} else {
			fmt.Printf("Created random %s: %s (Status: pending)\n", user.Role, user.Email)
		}
	}

	fmt.Println("\nSeeding complete!")
}