package dto

import (
	"time"
)

type UserResponse struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name" example:"John Doe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Role     string `json:"role" enums:"admin,instructor,student"`
}

type UserStatsResponse struct {
	TotalUsers   int `json:"total_users"`
	ActiveUsers  int `json:"active_users"`
	PendingUsers int `json:"pending_users"`
}
