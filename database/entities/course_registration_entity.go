package entities

import "gorm.io/gorm"


type CourseRegistration struct {
	gorm.Model
	UserID   uint
	CourseID uint

	Status string

	User   User
	Course Course

	Answers []RegistrationAnswer `gorm:"foreignKey:RegistrationID"`
}
