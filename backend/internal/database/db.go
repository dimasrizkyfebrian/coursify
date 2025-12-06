package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dimasrizkyfebrian/coursify/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Method ConnectDB
func ConnectDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
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
