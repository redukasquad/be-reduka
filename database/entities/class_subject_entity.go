package entities

import "gorm.io/gorm"

type ClassSubject struct {
	gorm.Model

	CourseID    uint   `json:"courseId" form:"courseId" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`

	// relations
	Course  Course        `json:"course,omitempty"`
	Lessons []ClassLesson `json:"lessons,omitempty" gorm:"foreignKey:SubjectID"`
}
