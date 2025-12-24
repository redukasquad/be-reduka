package entities

import "gorm.io/gorm"

type ClassSubject struct {
	gorm.Model
	CourseID uint
	Name     string
	Description string

	Course  Course
	Lessons []ClassLesson `gorm:"foreignKey:SubjectID"` 
}
