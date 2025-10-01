package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Import driver pgx secara anonim
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ConnectDB berfungsi untuk membuat dan mengembalikan koneksi ke database
func ConnectDB() *sql.DB {
	// Membuat Data Source Name (DSN) string dari environment variables
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	// Membuka koneksi ke database
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}

	// Memverifikasi koneksi ke database dengan Ping
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB: ", err)
	}

	log.Println("Database connected successfully")
	return db
}