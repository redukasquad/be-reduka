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
	JenisKelamin *bool  `json:"jenisKelamin" form:"jenisKelamin"`

	Kelas *string `json:"kelas" form:"kelas" binding:"omitempty,oneof='Kelas 10' 'Kelas 11' 'Kelas 12' 'Gapyer (Alumni)'" gorm:"type:enum('Kelas 10','Kelas 11','Kelas 12','Gapyer (Alumni)')"`

	Role *string `json:"role" gorm:"type:enum('STUDENT','ADMIN','TUTOR');default:'STUDENT'"`

	// AuthProvider indicates how the user registered (PASSWORD or GOOGLE)
	AuthProvider string `json:"authProvider" gorm:"type:enum('PASSWORD','GOOGLE');default:'PASSWORD'"`

	ProfileImage             string     `json:"profileImage" form:"profileImage"`
	IsVerified               bool       `json:"isVerified" gorm:"default:false"`
	VerificationCode         *string    `json:"verificationCode,omitempty" gorm:"size:100"`
	ResetPasswordToken       string     `json:"resetPasswordToken" gorm:"index"`
	ResetPasswordTokenExpiry *time.Time `json:"resetPasswordTokenExpiry"`

	// relations
	CourseRegistrations []CourseRegistration `json:"courseRegistrations,omitempty"`
	UserTargets         []UserTarget         `json:"userTargets,omitempty"`
}
