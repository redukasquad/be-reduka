package lessons

import "time"

type CreateLessonInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	LessonOrder int        `json:"lessonOrder" binding:"required,min=1"`
	StartTime   *time.Time `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
}

type UpdateLessonInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	LessonOrder *int       `json:"lessonOrder"`
	StartTime   *time.Time `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
}

type LessonResponse struct {
	ID            uint               `json:"id"`
	SubjectID     uint               `json:"subjectId"`
	SubjectName   string             `json:"subjectName,omitempty"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	LessonOrder   int                `json:"lessonOrder"`
	StartTime     *time.Time         `json:"startTime,omitempty"`
	EndTime       *time.Time         `json:"endTime,omitempty"`
	Resources     []ResourceResponse `json:"resources,omitempty"`
	ResourceCount int                `json:"resourceCount,omitempty"`
}

type ResourceResponse struct {
	ID    uint   `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
}
