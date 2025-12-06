package dto

import (
	"time"

	"github.com/dimasrizkyfebrian/coursify/internal/model"
)

type CreateCourseRequest struct {
	Title         string `json:"title" example:"Introduction to Go"`
	Description   string `json:"description" example:"A beginner's guide to Golang."`
	CoverImageURL string `json:"cover_image_url,omitempty" example:"/images/covers/cover-1.jpg"`
}

type AddMaterialRequest struct {
	Title       string `json:"title" example:"Chapter 1: Introduction"`
	ContentType string `json:"content_type" enums:"text,video,pdf"`
	TextContent string `json:"text_content,omitempty" example:"This is the lesson content."`
	VideoURL    string `json:"video_url,omitempty" example:"https://youtube.com/watch?v=..."`
}

type CourseResponse struct {
	ID             string    `json:"id"`
	InstructorID   string    `json:"instructor_id"`
	InstructorName string    `json:"instructor_name,omitempty"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	CoverImageURL  *string   `json:"cover_image_url,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CourseWithMaterialsResponse struct {
	CourseResponse
	Materials []model.LearningMaterial `json:"materials"`
}
