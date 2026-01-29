package subjects

type CreateSubjectInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateSubjectInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type SubjectResponse struct {
	ID          uint   `json:"id"`
	CourseID    uint   `json:"courseId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LessonCount int    `json:"lessonCount,omitempty"`
}
