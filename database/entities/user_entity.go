package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username                 string     `json:"username" binding:"required"`
	Email                    string     `json:"email" binding:"required,email" gorm:"uniqueIndex;type:varchar(191)"`
	Password                 string     `json:"password" binding:"required"`
	NoTelp                   string     `json:"no_telp" binding:"required"`
	JenisKelamin             bool       `json:"jenis_kelamin" binding:"required"`
	Kelas                    string     `json:"kelas" binding:"required,oneof='Kelas 10' 'Kelas 11' 'Kelas 12' 'Gapyer (Alumni)'" gorm:"type:enum('Kelas 10','Kelas 11','Kelas 12','Gapyer (Alumni)')"`
	Role                     string     `json:"role" binding:"required,oneof=Students Tutor Admin" gorm:"type:enum('Students','Tutor','Admin')"`
	ProfileImage             string     `json:"profile_image"`
	IsVerified               bool       `json:"is_verified" gorm:"default:false"`
	VerificationCode         string     `json:"verification_code" gorm:"omitempty"`
	ResetPasswordToken       string     `json:"reset_password_token" gorm:"index"`
	ResetPasswordTokenExpiry *time.Time `json:"reset_password_token_expiry"`
}
