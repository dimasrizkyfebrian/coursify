package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Method ConnectDB
func ConnectDB() *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	// Connection Database
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}

	// Verification with ping
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB: ", err)
	}

	log.Println("Database connected successfully")
	return db
}