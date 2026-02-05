package entities

import "gorm.io/gorm"

type UniversityProgram struct {
	gorm.Model

	UniversityID uint    `json:"universityId" form:"universityId" binding:"required"`
	Name         string  `json:"name" form:"name" binding:"required"`
	PassingGrade float64 `json:"passingGrade" form:"passingGrade" binding:"required,min=0,max=100"`

	// relations
	University University  `json:"university,omitempty"`
	Targets    []UserTarget `json:"targets,omitempty"`
}