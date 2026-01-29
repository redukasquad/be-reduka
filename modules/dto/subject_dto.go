package dto

import "github.com/redukasquad/be-reduka/database/entities"

type SubjectBriefResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func ToSubjectBriefResponse(subject entities.ClassSubject) SubjectBriefResponse {
	return SubjectBriefResponse{
		ID:          subject.ID,
		Name:        subject.Name,
		Description: subject.Description,
	}
}