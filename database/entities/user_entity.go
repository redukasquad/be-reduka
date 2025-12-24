package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username                 string     `binding:"required"`
	Email                    string     `binding:"required,email" gorm:"uniqueIndex;type:varchar(191)"`
	Password                 string     `binding:"required"`
	NoTelp                   string     
	JenisKelamin             bool       
	Kelas                    string     `binding:"required,oneof='Kelas 10' 'Kelas 11' 'Kelas 12' 'Gapyer (Alumni)'" gorm:"type:enum('Kelas 10','Kelas 11','Kelas 12','Gapyer (Alumni)')"`
	Role                     string     `binding:"required,oneof=Students Tutor Admin" gorm:"type:enum('Students','Tutor','Admin')"`
	ProfileImage             string   
	IsVerified               bool       `gorm:"default:false"`
	VerificationCode         string     `gorm:"omitempty"`
	ResetPasswordToken       string     `gorm:"index"`
	ResetPasswordTokenExpiry *time.Time 

	// relations
	CourseRegistrations []CourseRegistration 
	UserTargets         []UserTarget
}
