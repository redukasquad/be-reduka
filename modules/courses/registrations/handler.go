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
	RegisterHandler(c *gin.Context)
	GetMyRegistrationsHandler(c *gin.Context)
	GetRegistrationsByCourseHandler(c *gin.Context)
	ApproveRegistrationHandler(c *gin.Context)
	RejectRegistrationHandler(c *gin.Context)
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
	if id, ok := userID.(int); ok {
		return uint(id)
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

func (h *handler) RegisterHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	courseIDStr := c.Param("id")

	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Course ID", "Course ID must be a valid number", nil))
		return
	}

	var input RegisterCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		input = RegisterCourseInput{Answers: []AnswerInput{}}
	}

	registration, err := h.service.Register(uint(courseID), userID, input, requestID)
	if err != nil {
		if err.Error() == "you have already registered for this course" {
			c.JSON(http.StatusConflict, utils.BuildResponseFailed("Already registered", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to register", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Registration submitted successfully. Please wait for admin approval.", registration))
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

func (h *handler) GetRegistrationsByCourseHandler(c *gin.Context) {
	requestID := getRequestID(c)
	courseIDStr := c.Param("id")

	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Course ID", "Course ID must be a valid number", nil))
		return
	}

	registrations, err := h.service.GetRegistrationsByCourse(uint(courseID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch registrations", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Registrations retrieved successfully", registrations))
}

func (h *handler) ApproveRegistrationHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	registration, err := h.service.ApproveRegistration(uint(id), requestID, userID)
	if err != nil {
		if err.Error() == "registration not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Registration not found", err.Error(), nil))
			return
		}
		if err.Error() == "registration is not in pending status" {
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Cannot approve", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to approve registration", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Registration approved successfully", registration))
}

func (h *handler) RejectRegistrationHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	registration, err := h.service.RejectRegistration(uint(id), requestID, userID)
	if err != nil {
		if err.Error() == "registration not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Registration not found", err.Error(), nil))
			return
		}
		if err.Error() == "registration is not in pending status" {
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Cannot reject", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to reject registration", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Registration rejected successfully", registration))
}
