package schemas

type ChapterRequest struct {
	Name        string `json:"name" example:"Chapter 1: Getting Started"`
	Description string `json:"description" example:"Introduction to the course material"`
	Order       int    `json:"order" example:"1"`
	CreatedBy   *uint  `json:"created_by,omitempty"`
}

type ChapterResponse struct {
	ID          uint   `json:"id,omitempty" example:"1"`
	Name        string `json:"name" example:"Chapter 1: Getting Started"`
	Description string `json:"description" example:"Introduction to the course material"`
	CreatedBy   *uint  `json:"created_by,omitempty"`
	CreatedAt   string `json:"created_at,omitempty" example:"2020-01-01T12:00:00Z"`
	UpdatedAt   string `json:"updated_at,omitempty" example:"2020-01-01T12:00:00Z"`
}

type ChapterResponseWithLessonsCount struct {
	ID           uint   `json:"id,omitempty" example:"1"`
	Name         string `json:"name" example:"Chapter 1: Getting Started"`
	Description  string `json:"description" example:"Introduction to the course material"`
	CreatedBy    *uint  `json:"created_by,omitempty"`
	CreatedAt    string `json:"created_at,omitempty" example:"2020-01-01T12:00:00Z"`
	UpdatedAt    string `json:"updated_at,omitempty" example:"2020-01-01T12:00:00Z"`
	LessonsCount int    `json:"lessons_count" example:"1"`
}
