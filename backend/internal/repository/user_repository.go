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

// Method CreateUser
func (r *UserRepository) CreateUser(user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)

	query := `INSERT INTO users (full_name, email, password_hash, role)
	           VALUES ($1, $2, $3, $4)`

	_, err = r.DB.Exec(query, user.FullName, user.Email, user.PasswordHash, user.Role)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	return nil
}

// Method GetUserByEmail
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, full_name, email, password_hash, role, status FROM users WHERE email = $1`

	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role, &user.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Method GetUsersByStatus
func (r *UserRepository) GetUsersByStatus(status string) ([]model.User, error) {
	query := `SELECT id, full_name, email, role, status, created_at, updated_at FROM users WHERE status = $1 ORDER BY created_at ASC`

	rows, err := r.DB.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Method UpdateUserStatus
func (r *UserRepository) UpdateUserStatus(userID, status string) error {
	query := `UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.DB.Exec(query, status, userID)
	if err != nil {
		log.Printf("Error updating user status: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Method GetUsersByID
func (r *UserRepository) GetUsersByID(userID string) (*model.User, error) {
	var user model.User
	query := `SELECT id, full_name, email, role, status, created_at, updated_at FROM users WHERE id = $1`

	err := r.DB.QueryRow(query, userID).Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}