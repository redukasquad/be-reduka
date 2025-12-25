package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username     string `json:"username" form:"username" binding:"required"`
	Email        string `json:"email" form:"email" binding:"required,email" gorm:"uniqueIndex;type:varchar(191)"`
	Password     string `json:"password" form:"password" binding:"required"`
	NoTelp       string `json:"noTelp" form:"noTelp"`
	JenisKelamin bool   `json:"jenisKelamin" form:"jenisKelamin"`

	Kelas string `json:"kelas" form:"kelas" binding:"omitempty,oneof='Kelas 10' 'Kelas 11' 'Kelas 12' 'Gapyer (Alumni)'" gorm:"type:enum('Kelas 10','Kelas 11','Kelas 12','Gapyer (Alumni)')"`

	Role string `json:"role" form:"role" binding:"omitempty,oneof=Students Tutor Admin" gorm:"type:enum('Students','Tutor','Admin')"`

	ProfileImage string `json:"profileImage" form:"profileImage"`
	IsVerified               bool       `json:"isVerified" gorm:"default:false"`
	VerificationCode         string     `json:"verificationCode,omitempty"`
	ResetPasswordToken       string     `json:"resetPasswordToken" gorm:"index"`
	ResetPasswordTokenExpiry *time.Time `json:"resetPasswordTokenExpiry"`

	// relations
	CourseRegistrations []CourseRegistration `json:"courseRegistrations,omitempty"`
	UserTargets         []UserTarget         `json:"userTargets,omitempty"`
}

