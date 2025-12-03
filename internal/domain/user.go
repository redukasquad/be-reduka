package domain

import (
	"encoding/json"
)

type User struct {
	UserID            int             `gorm:"primaryKey;column:user_id" json:"user_id"`
	Username          string          `gorm:"type:varchar(100);not null" json:"username"`
	Email             string          `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password          string          `gorm:"type:varchar(255);not null" json:"-"`
	NoTelp            string          `gorm:"type:varchar(20)" json:"no_telp"`
	JenisKelamin      string          `gorm:"type:varchar(20)" json:"jenis_kelamin"`
	Kelas             string          `gorm:"type:enum('Kelas 12', 'Gapyear')" json:"kelas"`
	Role              string          `gorm:"type:enum('Students', 'Tutor', 'Admin')" json:"role"`
	ProfileImage      json.RawMessage `gorm:"type:json" json:"profile_image"`
	IsVerified        bool            `gorm:"default:false" json:"is_verified"`
	VerificationToken string          `gorm:"type:varchar(255)" json:"-"`
}

type UserRepository interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
	FindByVerificationToken(token string) (*User, error)
	UpdateUser(user *User) error
}