package schemas

import "time"

type CreateCourseRequest struct {
	Name        string `json:"name" example:"Introduction to Go Programming"`
	Description string `json:"description" example:"Learn the basics of Go programming language"`
}
type UpdateCourseRequest struct {
	Name        string `json:"name" example:"Introduction to Go Programming"`
	Description string `json:"description" example:"Learn the basics of Go programming language"`
}

type CourseResponse struct {
	ID          uint      `json:"id,omitempty" example:"1"`
	Name        string    `json:"name" example:"Introduction to Go Programming"`
	Description string    `json:"description" example:"Learn the basics of Go programming language"`
	CreatedAt   time.Time `json:"created_at,omitempty" example:"2020-01-01T12:00:00Z"`
}
