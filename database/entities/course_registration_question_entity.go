package entities

import "gorm.io/gorm"

type RegistrationQuestion struct {
	gorm.Model
	CourseID uint

	QuestionText  string
	QuestionType  string 
	QuestionOrder int

	Course  Course
	Answers []RegistrationAnswer `gorm:"foreignKey:QuestionID"` 
}
