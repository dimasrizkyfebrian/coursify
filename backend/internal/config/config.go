package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	JWTSecretKey       string
	Port               string
	BaseURL            string
	CorsAllowedOrigins []string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		// Try loading from parent directory (root) for local dev
		err = godotenv.Load("../.env")
		if err != nil {
			log.Println("Warning: .env file not found in current or parent directory, using environment variables from runtime")
		}
	}

	return &Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "coursify_db"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecretKey:       getEnv("JWT_SECRET_KEY", "secret"),
		Port:               getEnv("PORT", "8080"),
		BaseURL:            getEnv("BASE_URL", "http://localhost:8080"),
		CorsAllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"), ","),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
