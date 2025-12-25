package entities

import "gorm.io/gorm"

type CourseRegistration struct {
	gorm.Model

	UserID   uint `json:"userId" form:"userId" binding:"required"`
	CourseID uint `json:"courseId" form:"courseId" binding:"required"`

	Status string `json:"status" form:"status" binding:"required,oneof=pending approved rejected" gorm:"type:enum('pending','approved','rejected');default:'pending'"`

	// relations
	User   User   `json:"user,omitempty"`
	Course Course `json:"course,omitempty"`

	Answers []RegistrationAnswer `json:"answers,omitempty" gorm:"foreignKey:RegistrationID"`
}
