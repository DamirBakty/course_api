package schemas

type LessonRequest struct {
	Name        string `json:"name" example:"Introduction to Go"`
	Description string `json:"description" example:"Learn the basics of Go programming language"`
	Content     string `json:"content" example:"This lesson covers the basic concepts of the chapter."`
	Order       int    `json:"order" example:"1"`
	CreatedBy   *uint  `json:"created_by,omitempty"`
}

type LessonResponse struct {
	ID          uint   `json:"id,omitempty" example:"1"`
	Name        string `json:"name" example:"Introduction to Go"`
	Description string `json:"description" example:"Learn the basics of Go programming language"`
	Content     string `json:"content" example:"This lesson covers the basic concepts of the chapter."`
	Order       int    `json:"order" example:"1"`
	CreatedBy   *uint  `json:"created_by,omitempty"`
	CreatedAt   string `json:"created_at,omitempty" example:"2020-01-01T12:00:00Z"`
}
