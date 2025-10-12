// backend/cmd/seeder/main.go
package main

import (
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

func main() {
	// Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.ConnectDB()
	userRepo := repository.NewUserRepository(db)

	// Create a new random source based on the current time
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Role options in a slice
	roles := []string{"student", "instructor"}

	fmt.Println("Seeding users...")

	// Create 10 users
	for i := 0; i < 10; i++ {
		// Choose a random role from the roles slice
		randomRole := roles[r.Intn(len(roles))]

		user := model.User{
			FullName: faker.Name(),
			Email:    faker.Email(),
			Password: "password123", // Default password
			Role:     randomRole,
		}
		err := userRepo.CreateUser(&user)
		if err != nil {
			log.Printf("Could not create user %s: %v\n", user.Email, err)
		} else {
			fmt.Printf("Created %s: %s\n", user.Role, user.Email)
		}
	}

	fmt.Println("Seeding complete!")
}