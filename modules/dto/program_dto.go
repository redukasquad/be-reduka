package dto

import "github.com/redukasquad/be-reduka/database/entities"

type ProgramBriefResponse struct {
	ID           uint   `json:"id"`
	ProgramName  string `json:"programName"`
	ImageProgram string `json:"imageProgram,omitempty"`
}

type ProgramResponse struct {
	ID										uint							`json:"id"`
	ProgramName 					string						`json:"programName"`
	Description						string						`json:"description,omitempty"`
	ImageProgram					string						`json:"imageProgram,omitempty"`
	CourseCount 					int								`json:"courseCount"`
}

func ToProgramBriefResponse(program entities.Program) ProgramBriefResponse {
	return ProgramBriefResponse{
		ID:           program.ID,
		ProgramName:  program.ProgramName,
		ImageProgram: program.ImageProgram,
	}
}