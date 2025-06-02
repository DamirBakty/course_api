package schemas

type AttachmentResponse struct {
	ID        uint   `json:"id,omitempty" example:"1"`
	Name      string `json:"name" example:"lecture_slides.pdf"`
	URL       string `json:"url" example:"https://storage.example.com/files/lecture_slides.pdf"`
	LessonID  uint   `json:"lesson_id,omitempty" example:"1"`
	CreatedAt string `json:"created_at,omitempty" example:"2020-01-01T12:00:00Z"`
}

type UploadResponse struct {
	ID       uint   `json:"id,omitempty" example:"1"`
	Name     string `json:"name" example:"lecture_slides.pdf"`
	URL      string `json:"url" example:"https://storage.example.com/files/lecture_slides.pdf"`
	LessonID uint   `json:"lesson_id,omitempty" example:"1"`
}