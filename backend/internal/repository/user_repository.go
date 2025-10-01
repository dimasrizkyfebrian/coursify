package repository

import (
	"database/sql"
	"log"

	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser akan melakukan hash password dan menyimpan user baru ke DB
func (r *UserRepository) CreateUser(user *model.User) error {
	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)

	// Query untuk memasukkan user baru
	query := `INSERT INTO users (full_name, email, password_hash, role)
	           VALUES ($1, $2, $3, $4)`

	_, err = r.DB.Exec(query, user.FullName, user.Email, user.PasswordHash, user.Role)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	return nil
}