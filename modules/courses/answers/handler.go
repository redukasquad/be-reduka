package answers

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
	GetAnswersByRegistrationHandler(c *gin.Context)
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

func (h *handler) GetAnswersByRegistrationHandler(c *gin.Context) {
	requestID := getRequestID(c)
	registrationIDStr := c.Param("registrationId")

	registrationID, err := strconv.ParseUint(registrationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Registration ID", "Registration ID must be a valid number", nil))
		return
	}

	answers, err := h.service.GetByRegistrationID(uint(registrationID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch answers", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Answers retrieved successfully", answers))
}
