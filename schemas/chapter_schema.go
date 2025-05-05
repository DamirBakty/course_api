package schemas

type ChapterRequest struct {
	Name        string `json:"name" example:"Chapter 1: Getting Started"`
	Description string `json:"description" example:"Introduction to the course material"`
	Order       int    `json:"order" example:"1"`
}

type ChapterResponse struct {
	ID          uint   `json:"id,omitempty" example:"1"`
	Name        string `json:"name" example:"Chapter 1: Getting Started"`
	Description string `json:"description" example:"Introduction to the course material"`
	CreatedAt   string `json:"created_at,omitempty" example:"2020-01-01T12:00:00Z"`
	UpdatedAt   string `json:"updated_at,omitempty" example:"2020-01-01T12:00:00Z"`
}
