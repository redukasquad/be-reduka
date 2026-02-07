package registrations

import (
	"errors"
	"time"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type registrationService struct {
	repo Repository
}

type Service interface {
	// User actions
	Register(tryOutID uint, userID uint, requestID string) (*RegistrationResponse, error)
	UploadPaymentProof(registrationID uint, input UploadPaymentProofInput, userID uint, requestID string) (*RegistrationResponse, error)
	GetMyRegistrations(userID uint, requestID string) ([]RegistrationResponse, error)
	GetRegistrationByID(id uint, requestID string) (*RegistrationResponse, error)

	// Admin actions
	GetPendingPayments(requestID string) ([]PendingPaymentResponse, error)
	GetRegistrationsByTryOut(tryOutID uint, requestID string) ([]RegistrationResponse, error)
	ApprovePayment(registrationID uint, adminUserID uint, requestID string) (*RegistrationResponse, error)
	RejectPayment(registrationID uint, input ApprovePaymentInput, adminUserID uint, requestID string) (*RegistrationResponse, error)
}

func NewService(repo Repository) Service {
	return &registrationService{repo: repo}
}

// ==========================================
// User Actions
// ==========================================

func (s *registrationService) Register(tryOutID uint, userID uint, requestID string) (*RegistrationResponse, error) {
	utils.LogInfo("registrations", "register", "User attempting to register for try out", requestID, userID, map[string]any{
		"try_out_id": tryOutID,
	})

	// Check if try out exists
	tryOut, err := s.repo.FindTryOutByID(tryOutID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("try out not found")
		}
		return nil, err
	}

	// Check registration period
	now := time.Now()
	if now.Before(tryOut.RegistrationStart) {
		return nil, errors.New("registration has not started yet")
	}
	if now.After(tryOut.RegistrationEnd) {
		return nil, errors.New("registration period has ended")
	}

	// Check if already registered
	_, err = s.repo.FindByUserAndTryOut(userID, tryOutID)
	if err == nil {
		return nil, errors.New("you are already registered for this try out")
	}

	// Determine initial payment status
	paymentStatus := entities.PaymentStatusPending
	if tryOut.IsFree {
		paymentStatus = entities.PaymentStatusApproved
	}

	registration := &entities.TryOutRegistration{
		UserID:          userID,
		TryOutPackageID: tryOutID,
		PaymentStatus:   paymentStatus,
		RegisteredAt:    now,
	}

	if err := s.repo.Create(registration); err != nil {
		utils.LogError("registrations", "register", "Failed to create registration: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	// Fetch with preload
	createdReg, err := s.repo.FindByID(registration.ID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("registrations", "register", "Registration created successfully", requestID, userID, map[string]any{
		"registration_id": registration.ID,
		"is_free":         tryOut.IsFree,
	})

	response := ToRegistrationResponse(createdReg)
	return &response, nil
}

func (s *registrationService) UploadPaymentProof(registrationID uint, input UploadPaymentProofInput, userID uint, requestID string) (*RegistrationResponse, error) {
	utils.LogInfo("registrations", "upload_payment_proof", "User uploading payment proof", requestID, userID, map[string]any{
		"registration_id": registrationID,
	})

	registration, err := s.repo.FindByID(registrationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("registration not found")
		}
		return nil, err
	}

	// Check ownership
	if registration.UserID != userID {
		return nil, errors.New("you can only upload payment proof for your own registration")
	}

	// Check if already approved
	if registration.PaymentStatus == entities.PaymentStatusApproved {
		return nil, errors.New("payment is already approved")
	}

	// Check if it's a free try out
	if registration.TryOutPackage.IsFree {
		return nil, errors.New("no payment required for free try out")
	}

	registration.PaymentProofURL = input.PaymentProofURL
	registration.PaymentStatus = entities.PaymentStatusPending

	if err := s.repo.Update(&registration); err != nil {
		utils.LogError("registrations", "upload_payment_proof", "Failed to update registration: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	// Fetch with preload
	updatedReg, err := s.repo.FindByID(registration.ID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("registrations", "upload_payment_proof", "Payment proof uploaded successfully", requestID, userID, map[string]any{
		"registration_id": registrationID,
	})

	response := ToRegistrationResponse(updatedReg)
	return &response, nil
}

func (s *registrationService) GetMyRegistrations(userID uint, requestID string) ([]RegistrationResponse, error) {
	utils.LogInfo("registrations", "get_my_registrations", "Fetching user registrations", requestID, userID, nil)

	registrations, err := s.repo.FindByUserID(userID)
	if err != nil {
		utils.LogError("registrations", "get_my_registrations", "Failed to fetch registrations: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	var responses []RegistrationResponse
	for _, r := range registrations {
		responses = append(responses, ToRegistrationResponse(r))
	}

	return responses, nil
}

func (s *registrationService) GetRegistrationByID(id uint, requestID string) (*RegistrationResponse, error) {
	utils.LogInfo("registrations", "get_by_id", "Fetching registration by ID", requestID, 0, map[string]any{
		"registration_id": id,
	})

	registration, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("registration not found")
		}
		return nil, err
	}

	response := ToRegistrationResponse(registration)
	return &response, nil
}

// ==========================================
// Admin Actions
// ==========================================

func (s *registrationService) GetPendingPayments(requestID string) ([]PendingPaymentResponse, error) {
	utils.LogInfo("registrations", "get_pending_payments", "Fetching pending payments", requestID, 0, nil)

	registrations, err := s.repo.FindPendingPayments()
	if err != nil {
		utils.LogError("registrations", "get_pending_payments", "Failed to fetch pending payments: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var responses []PendingPaymentResponse
	for _, r := range registrations {
		responses = append(responses, ToPendingPaymentResponse(r))
	}

	return responses, nil
}

func (s *registrationService) GetRegistrationsByTryOut(tryOutID uint, requestID string) ([]RegistrationResponse, error) {
	utils.LogInfo("registrations", "get_by_tryout", "Fetching registrations by try out", requestID, 0, map[string]any{
		"try_out_id": tryOutID,
	})

	registrations, err := s.repo.FindByTryOutID(tryOutID)
	if err != nil {
		utils.LogError("registrations", "get_by_tryout", "Failed to fetch registrations: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var responses []RegistrationResponse
	for _, r := range registrations {
		responses = append(responses, ToRegistrationResponse(r))
	}

	return responses, nil
}

func (s *registrationService) ApprovePayment(registrationID uint, adminUserID uint, requestID string) (*RegistrationResponse, error) {
	utils.LogInfo("registrations", "approve_payment", "Admin approving payment", requestID, adminUserID, map[string]any{
		"registration_id": registrationID,
	})

	registration, err := s.repo.FindByID(registrationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("registration not found")
		}
		return nil, err
	}

	if registration.PaymentStatus == entities.PaymentStatusApproved {
		return nil, errors.New("payment is already approved")
	}

	if registration.PaymentProofURL == "" {
		return nil, errors.New("no payment proof uploaded yet")
	}

	now := time.Now()
	registration.PaymentStatus = entities.PaymentStatusApproved
	registration.ApprovedByUserID = &adminUserID
	registration.ApprovedAt = &now
	registration.RejectionReason = ""

	if err := s.repo.Update(&registration); err != nil {
		utils.LogError("registrations", "approve_payment", "Failed to approve payment: "+err.Error(), requestID, adminUserID, nil)
		return nil, err
	}

	// Fetch with preload
	updatedReg, err := s.repo.FindByID(registration.ID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("registrations", "approve_payment", "Payment approved successfully", requestID, adminUserID, map[string]any{
		"registration_id": registrationID,
	})

	response := ToRegistrationResponse(updatedReg)
	return &response, nil
}

func (s *registrationService) RejectPayment(registrationID uint, input ApprovePaymentInput, adminUserID uint, requestID string) (*RegistrationResponse, error) {
	utils.LogInfo("registrations", "reject_payment", "Admin rejecting payment", requestID, adminUserID, map[string]any{
		"registration_id": registrationID,
	})

	registration, err := s.repo.FindByID(registrationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("registration not found")
		}
		return nil, err
	}

	if registration.PaymentStatus == entities.PaymentStatusApproved {
		return nil, errors.New("cannot reject an already approved payment")
	}

	registration.PaymentStatus = entities.PaymentStatusRejected
	registration.RejectionReason = input.RejectionReason
	registration.PaymentProofURL = "" // Clear proof so user can re-upload

	if err := s.repo.Update(&registration); err != nil {
		utils.LogError("registrations", "reject_payment", "Failed to reject payment: "+err.Error(), requestID, adminUserID, nil)
		return nil, err
	}

	// Fetch with preload
	updatedReg, err := s.repo.FindByID(registration.ID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("registrations", "reject_payment", "Payment rejected", requestID, adminUserID, map[string]any{
		"registration_id": registrationID,
		"reason":          input.RejectionReason,
	})

	response := ToRegistrationResponse(updatedReg)
	return &response, nil
}
