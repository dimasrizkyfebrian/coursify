package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/dimasrizkyfebrian/coursify/internal/database"
)

func main() {
	// Memuat variabel dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

    // Memanggil fungsi koneksi database
    db := database.ConnectDB()
	log.Printf("Database instance created: %p", db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Coursify API!")
	})

	port := ":8080"
	log.Printf("Server is starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}