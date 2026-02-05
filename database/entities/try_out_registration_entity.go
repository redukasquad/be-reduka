package entities

import (
	"time"

	"gorm.io/gorm"
)

// PaymentStatus represents the status of a Try Out registration payment.
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusApproved PaymentStatus = "approved"
	PaymentStatusRejected PaymentStatus = "rejected"
)

// TryOutRegistration represents a user's registration to a Try Out package.
// For paid packages, includes payment proof and approval workflow.
type TryOutRegistration struct {
	gorm.Model

	UserID          uint `json:"userId" gorm:"uniqueIndex:idx_user_tryout;not null"`
	TryOutPackageID uint `json:"tryOutPackageId" gorm:"uniqueIndex:idx_user_tryout;not null"`

	// Payment
	PaymentProofURL string        `json:"paymentProofUrl" gorm:"size:500"` // Bukti bayar (ImageKit)
	PaymentStatus   PaymentStatus `json:"paymentStatus" gorm:"size:20;default:'pending'"`
	RejectionReason string        `json:"rejectionReason" gorm:"type:text"`

	ApprovedByUserID *uint      `json:"approvedByUserId"`
	ApprovedAt       *time.Time `json:"approvedAt"`
	RegisteredAt     time.Time  `json:"registeredAt" gorm:"autoCreateTime"`

	// Relations
	User          User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	TryOutPackage TryOut  `json:"tryOutPackage,omitempty" gorm:"foreignKey:TryOutPackageID"`
	ApprovedBy    *User          `json:"approvedBy,omitempty" gorm:"foreignKey:ApprovedByUserID"`
	Attempt       *TryOutAttempt `json:"attempt,omitempty" gorm:"foreignKey:RegistrationID"`
}
