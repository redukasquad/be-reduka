package entities

import (
	"time"

	"gorm.io/gorm"
)

// TryOutPackage represents a Try Out package created by admin.
// Each package contains 7 subtests with questions created by tutors.
type TryOut struct {
	gorm.Model

	Name        string `json:"name" gorm:"size:100;not null"`
	Description string `json:"description" gorm:"type:text"`
	ImageURL    string `json:"imageUrl" gorm:"size:500"` // Avatar/thumbnail

	// Pricing
	IsFree       bool    `json:"isFree" gorm:"default:false"`
	Price        float64 `json:"price" gorm:"type:decimal(12,2)"`
	QrisImageURL string  `json:"qrisImageUrl" gorm:"size:500"`
	PaymentLink  string  `json:"paymentLink" gorm:"size:500"`

	// Registration period
	RegistrationStart time.Time `json:"registrationStart" gorm:"not null"`
	RegistrationEnd   time.Time `json:"registrationEnd" gorm:"not null"`

	IsPublished     bool `json:"isPublished" gorm:"default:false"`
	CreatedByUserID uint `json:"createdByUserId"`

	// Relations
	Creator          User                 `json:"creator,omitempty" gorm:"foreignKey:CreatedByUserID"`
	TutorPermissions []TutorPermission    `json:"tutorPermissions,omitempty" gorm:"foreignKey:TryOutPackageID"`
	Questions        []TryOutQuestion     `json:"questions,omitempty" gorm:"foreignKey:TryOutPackageID"`
	Registrations    []TryOutRegistration `json:"registrations,omitempty" gorm:"foreignKey:TryOutPackageID"`
}
