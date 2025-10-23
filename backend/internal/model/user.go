package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID           string            `json:"id"`
	FullName     string            `json:"full_name"`
	Email        string            `json:"email"`
	Password     string            `json:"-"`
	PasswordHash string            `json:"-"`
	Role         string            `json:"role"`
	Status       string            `json:"status"`
	AvatarURL    sql.NullString    `json:"avatar_url,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}