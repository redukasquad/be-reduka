package registrations

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redukasquad/be-reduka/packages/utils"
)

type handler struct {
	service Service
}

type Handler interface {
	// User endpoints
	RegisterHandler(c *gin.Context)
	UploadPaymentProofHandler(c *gin.Context)
	GetMyRegistrationsHandler(c *gin.Context)

	// Admin endpoints
	GetPendingPaymentsHandler(c *gin.Context)
	GetRegistrationsByTryOutHandler(c *gin.Context)
	ApprovePaymentHandler(c *gin.Context)
	RejectPaymentHandler(c *gin.Context)
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func getRequestID(c *gin.Context) string {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
}

func getUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	switch id := userID.(type) {
	case int:
		return uint(id)
	case uint:
		return id
	case float64:
		return uint(id)
	}
	return 0
}

// ==========================================
// User Handlers
// ==========================================

func (h *handler) RegisterHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	tryOutIDStr := c.Param("id")

	tryOutID, err := strconv.ParseUint(tryOutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	registration, err := h.service.Register(uint(tryOutID), userID, requestID)
	if err != nil {
		switch err.Error() {
		case "try out not found":
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Try out not found", err.Error(), nil))
		case "registration has not started yet", "registration period has ended":
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Registration period invalid", err.Error(), nil))
		case "you are already registered for this try out":
			c.JSON(http.StatusConflict, utils.BuildResponseFailed("Already registered", err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to register", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Registration successful", registration))
}

func (h *handler) UploadPaymentProofHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	registrationIDStr := c.Param("id")

	registrationID, err := strconv.ParseUint(registrationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	var input UploadPaymentProofInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	registration, err := h.service.UploadPaymentProof(uint(registrationID), input, userID, requestID)
	if err != nil {
		switch err.Error() {
		case "registration not found":
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Registration not found", err.Error(), nil))
		case "you can only upload payment proof for your own registration":
			c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		case "payment is already approved", "no payment required for free try out":
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid operation", err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to upload payment proof", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Payment proof uploaded successfully", registration))
}

func (h *handler) GetMyRegistrationsHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)

	registrations, err := h.service.GetMyRegistrations(userID, requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch registrations", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Registrations retrieved successfully", registrations))
}

// ==========================================
// Admin Handlers
// ==========================================

func (h *handler) GetPendingPaymentsHandler(c *gin.Context) {
	requestID := getRequestID(c)

	payments, err := h.service.GetPendingPayments(requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch pending payments", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Pending payments retrieved successfully", payments))
}

func (h *handler) GetRegistrationsByTryOutHandler(c *gin.Context) {
	requestID := getRequestID(c)
	tryOutIDStr := c.Param("id")

	tryOutID, err := strconv.ParseUint(tryOutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	registrations, err := h.service.GetRegistrationsByTryOut(uint(tryOutID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch registrations", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Registrations retrieved successfully", registrations))
}

func (h *handler) ApprovePaymentHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	registrationIDStr := c.Param("id")

	registrationID, err := strconv.ParseUint(registrationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	registration, err := h.service.ApprovePayment(uint(registrationID), userID, requestID)
	if err != nil {
		switch err.Error() {
		case "registration not found":
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Registration not found", err.Error(), nil))
		case "payment is already approved", "no payment proof uploaded yet":
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid operation", err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to approve payment", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Payment approved successfully", registration))
}

func (h *handler) RejectPaymentHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	registrationIDStr := c.Param("id")

	registrationID, err := strconv.ParseUint(registrationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	var input ApprovePaymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	registration, err := h.service.RejectPayment(uint(registrationID), input, userID, requestID)
	if err != nil {
		switch err.Error() {
		case "registration not found":
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Registration not found", err.Error(), nil))
		case "cannot reject an already approved payment":
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid operation", err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to reject payment", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Payment rejected", registration))
}
