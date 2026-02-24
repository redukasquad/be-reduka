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

type TryOutRegistration struct {
	gorm.Model

	UserID          uint `json:"userId" gorm:"uniqueIndex:idx_user_tryout;not null"`
	TryOutPackageID uint `json:"tryOutPackageId" gorm:"uniqueIndex:idx_user_tryout;not null"`

	PaymentProofURL string        `json:"paymentProofUrl" gorm:"size:500"`
	PaymentStatus   PaymentStatus `json:"paymentStatus" gorm:"size:20;default:'pending'"`
	RejectionReason string        `json:"rejectionReason" gorm:"type:text"`

	ApprovedByUserID *uint      `json:"approvedByUserId"`
	ApprovedAt       *time.Time `json:"approvedAt"`
	RegisteredAt     time.Time  `json:"registeredAt" gorm:"autoCreateTime"`

	// ✅ RELATIONS (FIXED)
	User          User           `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	TryOutPackage TryOut         `json:"tryOutPackage,omitempty" gorm:"foreignKey:TryOutPackageID;references:ID"`
	ApprovedBy    *User          `json:"approvedBy,omitempty" gorm:"foreignKey:ApprovedByUserID;references:ID"`
	Attempt       *TryOutAttempt `json:"attempt,omitempty" gorm:"foreignKey:RegistrationID;references:ID"`
}