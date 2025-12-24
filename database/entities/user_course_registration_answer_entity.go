package entities

import "gorm.io/gorm"

type RegistrationAnswer struct {
	gorm.Model
	RegistrationID uint
	QuestionID     uint

	AnswerText string

	Registration CourseRegistration `gorm:"foreignKey:RegistrationID"`
	Question     RegistrationQuestion `gorm:"foreignKey:QuestionID"`
}
