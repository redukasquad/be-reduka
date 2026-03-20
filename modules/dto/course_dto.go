package dto

import (
	"time"

	"github.com/redukasquad/be-reduka/database/entities"
)

type CourseResponse struct {
	ID                uint                   `json:"id"`
	NameCourse        string                 `json:"nameCourse"`
	Description       string                 `json:"description,omitempty"`
	StartDate         time.Time              `json:"startDate"`
	EndDate           time.Time              `json:"endDate"`
	IsFree            bool                   `json:"isFree"`
	WhatsAppGroupLink string                 `json:"whatsAppGroupLink,omitempty"`
	Program           *ProgramBriefResponse  `json:"program,omitempty"`
	Creator           *CreatorResponse       `json:"creator,omitempty"`
	Classes           []SubjectBriefResponse `json:"classes,omitempty"`
	CreatedAt         time.Time              `json:"createdAt"`
	Image             string                 `json:"image,omitempty"`
}

func ToCourseResponse(course entities.Course) CourseResponse {
	response := CourseResponse{
		ID:                course.ID,
		NameCourse:        course.NameCourse,
		Description:       course.Description,
		StartDate:         course.StartDate,
		EndDate:           course.EndDate,
		IsFree:            course.IsFree,
		Image:             course.Image,
		WhatsAppGroupLink: course.WhatsappGroupLink,
		CreatedAt:         course.CreatedAt,
	}

	// Map Program jika ada (ID != 0)
	if course.Program.ID != 0 {
		program := ToProgramBriefResponse(course.Program)
		response.Program = &program
	}

	// Map Creator jika ada
	if course.Creator.ID != 0 {
		creator := ToCreatorResponse(course.Creator)
		response.Creator = &creator
	}

	// Map Classes
	for _, class := range course.Classes {
		response.Classes = append(response.Classes, ToSubjectBriefResponse(class))
	}

	return response
}
