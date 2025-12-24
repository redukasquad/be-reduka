package entities

import "gorm.io/gorm"

type Program struct {
	gorm.Model

	ProgramName  string `json:"programName" form:"programName" binding:"required" gorm:"type:varchar(100);uniqueIndex"`
	Description  string `json:"description" form:"description"`
	ImageProgram string `json:"imageProgram" form:"imageProgram"`

	// relations
	Courses []Course `json:"courses,omitempty"`
}
