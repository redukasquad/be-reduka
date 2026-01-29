package dto

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"time"
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
	Subjects          []SubjectBriefResponse `json:"subjects,omitempty"`
	CreatedAt         time.Time              `json:"createdAt"`
} 

func ToCourseResponse(course entities.Course) CourseResponse {
	response := CourseResponse{
			ID:                course.ID,
			NameCourse:        course.NameCourse,
			Description:       course.Description,
			StartDate:         course.StartDate,
			EndDate:           course.EndDate,
			IsFree:            course.IsFree,
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
	
	// Map Subjects
	for _, subject := range course.Subjects {
			response.Subjects = append(response.Subjects, ToSubjectBriefResponse(subject))
	}
	
	return response
}