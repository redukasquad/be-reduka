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
	CreateAnswerHandler(c *gin.Context)
	DeleteAnswerHandler(c *gin.Context)
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

func (h *handler) CreateAnswerHandler(c *gin.Context) {
	requestID := getRequestID(c)

	var input CreateAnswerRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid request body", err.Error(), nil))
		return
	}

	answer, err := h.service.CreateAnswer(input, requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create answer", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Answer created successfully", answer))
}

func (h *handler) DeleteAnswerHandler(c *gin.Context) {
	requestID := getRequestID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Answer ID", "Answer ID must be a valid number", nil))
		return
	}

	err = h.service.DeleteAnswer(uint(id), requestID)
	if err != nil {
		if err.Error() == "answer not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Answer not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete answer", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Answer deleted successfully", nil))
}
