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

// GetUserByEmail mengambil satu user dari DB berdasarkan email
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, full_name, email, password_hash, role, status FROM users WHERE email = $1`

	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role, &user.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			// Ini bukan error fatal, hanya user tidak ditemukan
			return nil, nil
		}
		// Error lain saat query
		return nil, err
	}

	return &user, nil
}