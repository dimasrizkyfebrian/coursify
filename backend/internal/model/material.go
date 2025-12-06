package model

import "time"

type LearningMaterial struct {
    ID           string    `json:"id"`
    CourseID     string    `json:"course_id"`
    Title        string    `json:"title"`
    ContentType  string    `json:"content_type"` // 'text', 'video', 'pdf'
    TextContent  string    `json:"text_content,omitempty"`
    VideoURL     string    `json:"video_url,omitempty"`
    FileURL      string    `json:"file_url,omitempty"`
    Position     int       `json:"position"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}