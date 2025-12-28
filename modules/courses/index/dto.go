package courses

import (
	"time"
)

// CreateCourseInput is the input for creating a new course
type CreateCourseInput struct {
	ProgramID         uint      `json:"programId" binding:"required"`
	NameCourse        string    `json:"nameCourse" binding:"required"`
	Description       string    `json:"description"`
	StartDate         time.Time `json:"startDate" binding:"required"`
	EndDate           time.Time `json:"endDate" binding:"required"`
	IsFree            bool      `json:"isFree"`
	WhatsappGroupLink string    `json:"whatsappGroupLink"`
}

// UpdateCourseInput is the input for updating a course
type UpdateCourseInput struct {
	ProgramID         *uint      `json:"programId"`
	NameCourse        *string    `json:"nameCourse"`
	Description       *string    `json:"description"`
	StartDate         *time.Time `json:"startDate"`
	EndDate           *time.Time `json:"endDate"`
	IsFree            *bool      `json:"isFree"`
	WhatsappGroupLink *string    `json:"whatsappGroupLink"`
}

// CourseResponse is the response format for a course (public view)
type CourseResponse struct {
	ID                uint      `json:"id"`
	ProgramID         uint      `json:"programId"`
	NameCourse        string    `json:"nameCourse"`
	Description       string    `json:"description"`
	StartDate         time.Time `json:"startDate"`
	EndDate           time.Time `json:"endDate"`
	IsFree            bool      `json:"isFree"`
	WhatsappGroupLink string    `json:"whatsappGroupLink,omitempty"` // Only shown to approved users or admin
	CreatedAt         time.Time `json:"createdAt"`
}

// CourseListResponse is the response for listing courses
type CourseListResponse struct {
	Courses    []CourseResponse `json:"courses"`
	TotalCount int64            `json:"totalCount"`
}
