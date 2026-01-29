package resources

type CreateResourceInput struct {
	Type  string `json:"type" binding:"required,oneof=video document link zoom recording"`
	Title string `json:"title" binding:"required"`
	URL   string `json:"url" binding:"required,url"`
}

type UpdateResourceInput struct {
	Type  *string `json:"type"`
	Title *string `json:"title"`
	URL   *string `json:"url"`
}

type ResourceResponse struct {
	ID            uint   `json:"id"`
	ClassLessonID uint   `json:"classLessonId"`
	LessonTitle   string `json:"lessonTitle,omitempty"`
	Type          string `json:"type"`
	Title         string `json:"title"`
	URL           string `json:"url"`
}
