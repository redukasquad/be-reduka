package tryouts

import (
	"time"
)

// ==========================================
// TRY OUT DTOs
// ==========================================

// TryOutResponse is the response DTO for Try Out
type TryOutResponse struct {
	ID                uint                  `json:"id"`
	Name              string                `json:"name"`
	Description       string                `json:"description,omitempty"`
	ImageURL          string                `json:"imageUrl,omitempty"`
	IsFree            bool                  `json:"isFree"`
	Price             float64               `json:"price,omitempty"`
	QrisImageURL      string                `json:"qrisImageUrl,omitempty"`
	PaymentLink       string                `json:"paymentLink,omitempty"`
	RegistrationStart time.Time             `json:"registrationStart"`
	RegistrationEnd   time.Time             `json:"registrationEnd"`
	IsPublished       bool                  `json:"isPublished"`
	Creator           *CreatorBriefResponse `json:"creator,omitempty"`
	CreatedAt         time.Time             `json:"createdAt"`
}

// TryOutBriefResponse is a minimal response for list views
type TryOutBriefResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	ImageURL          string    `json:"imageUrl,omitempty"`
	IsFree            bool      `json:"isFree"`
	Price             float64   `json:"price,omitempty"`
	RegistrationStart time.Time `json:"registrationStart"`
	RegistrationEnd   time.Time `json:"registrationEnd"`
	IsPublished       bool      `json:"isPublished"`
}

// CreatorBriefResponse is a minimal creator info
type CreatorBriefResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// CreateTryOutInput is the input for creating a new Try Out
type CreateTryOutInput struct {
	Name              string    `json:"name" binding:"required"`
	Description       string    `json:"description"`
	ImageURL          string    `json:"imageUrl"`
	IsFree            bool      `json:"isFree"`
	Price             float64   `json:"price"`
	QrisImageURL      string    `json:"qrisImageUrl"`
	PaymentLink       string    `json:"paymentLink"`
	RegistrationStart time.Time `json:"registrationStart" binding:"required"`
	RegistrationEnd   time.Time `json:"registrationEnd" binding:"required"`
	IsPublished       bool      `json:"isPublished"`
}

// UpdateTryOutInput is the input for updating a Try Out
type UpdateTryOutInput struct {
	Name              *string    `json:"name"`
	Description       *string    `json:"description"`
	ImageURL          *string    `json:"imageUrl"`
	IsFree            *bool      `json:"isFree"`
	Price             *float64   `json:"price"`
	QrisImageURL      *string    `json:"qrisImageUrl"`
	PaymentLink       *string    `json:"paymentLink"`
	RegistrationStart *time.Time `json:"registrationStart"`
	RegistrationEnd   *time.Time `json:"registrationEnd"`
	IsPublished       *bool      `json:"isPublished"`
}

// ==========================================
// TUTOR PERMISSION DTOs
// ==========================================

// TutorPermissionResponse is the response for tutor permission
type TutorPermissionResponse struct {
	ID        uint                  `json:"id"`
	TryOutID  uint                  `json:"tryOutId"`
	User      *CreatorBriefResponse `json:"user"`
	GrantedBy *CreatorBriefResponse `json:"grantedBy,omitempty"`
	GrantedAt time.Time             `json:"grantedAt"`
}

// GrantTutorPermissionInput is the input for granting tutor permission
type GrantTutorPermissionInput struct {
	UserID uint `json:"userId" binding:"required"`
}
