package entities

import "gorm.io/gorm"

type UserTarget struct {
	gorm.Model

	UserID              uint `json:"userId" form:"userId" binding:"required"`
	UniversityProgramID uint `json:"universityProgramId" form:"universityProgramId" binding:"required"`
	Priority            int  `json:"priority" form:"priority" binding:"required,min=1"`

	// relations
	User              User              `json:"user,omitempty"`
	UniversityProgram UniversityProgram `json:"universityProgram,omitempty"`
}
