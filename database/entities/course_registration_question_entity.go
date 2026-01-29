package entities

import "gorm.io/gorm"

type RegistrationQuestion struct {
	gorm.Model

	CourseID uint `json:"courseId" form:"courseId" binding:"required"`

	QuestionText  string `json:"questionText" form:"questionText" binding:"required"`
	QuestionType  string `json:"questionType" form:"questionType" binding:"required,oneof=text textarea select radio checkbox"`
	QuestionOrder int    `json:"questionOrder" form:"questionOrder" binding:"required,min=1"`

	// relations
	Course  Course               `json:"course,omitempty"`
	Answers []RegistrationAnswer `json:"answers,omitempty" gorm:"foreignKey:QuestionID"`
}
