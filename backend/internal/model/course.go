package model

import (
	"database/sql"
	"time"
)

type Course struct {
    ID              string            `json:"id"`
    InstructorID    string            `json:"instructor_id"`
    Title           string            `json:"title"`
    Description     string            `json:"description"`
    CoverImageURL   sql.NullString    `json:"cover_image_url,omitzero"`
    CreatedAt       time.Time         `json:"created_at"`
    UpdatedAt       time.Time         `json:"updated_at"`
}