package registrations

import (
	"time"

	"github.com/redukasquad/be-reduka/database/entities"
)

// ==========================================
// REGISTRATION DTOs
// ==========================================

// RegistrationResponse is the response DTO for a registration
type RegistrationResponse struct {
	ID              uint                 `json:"id"`
	TryOutID        uint                 `json:"tryOutId"`
	TryOut          *TryOutBriefResponse `json:"tryOut,omitempty"`
	User            *UserBriefResponse   `json:"user,omitempty"`
	PaymentProofURL string               `json:"paymentProofUrl,omitempty"`
	PaymentStatus   string               `json:"paymentStatus"`
	RejectionReason string               `json:"rejectionReason,omitempty"`
	ApprovedBy      *UserBriefResponse   `json:"approvedBy,omitempty"`
	ApprovedAt      *time.Time           `json:"approvedAt,omitempty"`
	RegisteredAt    time.Time            `json:"registeredAt"`
	HasAttempt      bool                 `json:"hasAttempt"`
}

// TryOutBriefResponse is a minimal try out info
type TryOutBriefResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	ImageURL          string    `json:"imageUrl,omitempty"`
	IsFree            bool      `json:"isFree"`
	Price             float64   `json:"price,omitempty"`
	RegistrationStart time.Time `json:"registrationStart"`
	RegistrationEnd   time.Time `json:"registrationEnd"`
}

// UserBriefResponse is a minimal user info
type UserBriefResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// PendingPaymentResponse for admin view
type PendingPaymentResponse struct {
	ID              uint                 `json:"id"`
	TryOut          *TryOutBriefResponse `json:"tryOut,omitempty"`
	User            *UserBriefResponse   `json:"user,omitempty"`
	PaymentProofURL string               `json:"paymentProofUrl"`
	RegisteredAt    time.Time            `json:"registeredAt"`
}

// UploadPaymentProofInput is the input for uploading payment proof
type UploadPaymentProofInput struct {
	PaymentProofURL string `json:"paymentProofUrl" binding:"required"`
}

// ApprovePaymentInput is the input for approving/rejecting payment
type ApprovePaymentInput struct {
	RejectionReason string `json:"rejectionReason"` // Only used when rejecting
}

// ==========================================
// Helper Functions
// ==========================================

func ToRegistrationResponse(r entities.TryOutRegistration) RegistrationResponse {
	response := RegistrationResponse{
		ID:              r.ID,
		TryOutID:        r.TryOutPackageID,
		PaymentProofURL: r.PaymentProofURL,
		PaymentStatus:   string(r.PaymentStatus),
		RejectionReason: r.RejectionReason,
		ApprovedAt:      r.ApprovedAt,
		RegisteredAt:    r.RegisteredAt,
		HasAttempt:      r.Attempt != nil,
	}

	if r.TryOutPackage.ID != 0 {
		tryOut := ToTryOutBriefResponse(r.TryOutPackage)
		response.TryOut = &tryOut
	}

	if r.User.ID != 0 {
		user := ToUserBriefResponse(r.User)
		response.User = &user
	}

	if r.ApprovedBy != nil && r.ApprovedBy.ID != 0 {
		approvedBy := ToUserBriefResponse(*r.ApprovedBy)
		response.ApprovedBy = &approvedBy
	}

	return response
}

func ToPendingPaymentResponse(r entities.TryOutRegistration) PendingPaymentResponse {
	response := PendingPaymentResponse{
		ID:              r.ID,
		PaymentProofURL: r.PaymentProofURL,
		RegisteredAt:    r.RegisteredAt,
	}

	if r.TryOutPackage.ID != 0 {
		tryOut := ToTryOutBriefResponse(r.TryOutPackage)
		response.TryOut = &tryOut
	}

	if r.User.ID != 0 {
		user := ToUserBriefResponse(r.User)
		response.User = &user
	}

	return response
}

func ToTryOutBriefResponse(t entities.TryOut) TryOutBriefResponse {
	return TryOutBriefResponse{
		ID:                t.ID,
		Name:              t.Name,
		ImageURL:          t.ImageURL,
		IsFree:            t.IsFree,
		Price:             t.Price,
		RegistrationStart: t.RegistrationStart,
		RegistrationEnd:   t.RegistrationEnd,
	}
}

func ToUserBriefResponse(u entities.User) UserBriefResponse {
	return UserBriefResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}
