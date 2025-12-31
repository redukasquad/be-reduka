package courses

import "time"

type CreateCourseInput struct {
	ProgramID         uint      `json:"programId" binding:"required"`
	NameCourse        string    `json:"nameCourse" binding:"required"`
	Description       string    `json:"description"`
	StartDate         time.Time `json:"startDate" binding:"required"`
	EndDate           time.Time `json:"endDate" binding:"required"`
	IsFree            bool      `json:"isFree"`
	WhatsappGroupLink string    `json:"whatsappGroupLink"`
}

type UpdateCourseInput struct {
	ProgramID         *uint      `json:"programId"`
	NameCourse        *string    `json:"nameCourse"`
	Description       *string    `json:"description"`
	StartDate         *time.Time `json:"startDate"`
	EndDate           *time.Time `json:"endDate"`
	IsFree            *bool      `json:"isFree"`
	WhatsappGroupLink *string    `json:"whatsappGroupLink"`
}

type CourseResponse struct {
	ID                uint      `json:"id"`
	ProgramID         uint      `json:"programId"`
	NameCourse        string    `json:"nameCourse"`
	Description       string    `json:"description"`
	StartDate         time.Time `json:"startDate"`
	EndDate           time.Time `json:"endDate"`
	IsFree            bool      `json:"isFree"`
	WhatsappGroupLink string    `json:"whatsappGroupLink,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
}

type CourseListResponse struct {
	Courses    []CourseResponse `json:"courses"`
	TotalCount int64            `json:"totalCount"`
}
