package entities

import "gorm.io/gorm"

type RegistrationAnswer struct {
	gorm.Model

	RegistrationID uint   `json:"registrationId" form:"registrationId" binding:"required"`
	QuestionID     uint   `json:"questionId" form:"questionId" binding:"required"`
	AnswerText     string `json:"answerText" form:"answerText" binding:"required"`

	// relations
	Registration CourseRegistration   `json:"registration,omitempty" gorm:"foreignKey:RegistrationID"`
	Question     RegistrationQuestion `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
}
