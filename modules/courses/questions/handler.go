package questions

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

// Handler interface defines the HTTP handlers for questions
type Handler interface {
	GetQuestionsByCourseHandler(c *gin.Context)
	CreateQuestionHandler(c *gin.Context)
	UpdateQuestionHandler(c *gin.Context)
	DeleteQuestionHandler(c *gin.Context)
}

// NewHandler creates a new question handler
func NewHandler(service Service) Handler {
	return &handler{service: service}
}

// getRequestID gets or generates a request ID from context
func getRequestID(c *gin.Context) string {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
}

// getUserID gets the user ID from context (set by auth middleware)
func getUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

// GetQuestionsByCourseHandler handles GET /courses/:id/questions
func (h *handler) GetQuestionsByCourseHandler(c *gin.Context) {
	requestID := getRequestID(c)
	courseIDStr := c.Param("id")

	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Course ID", "Course ID must be a valid number", nil))
		return
	}

	questions, err := h.service.GetByCourseID(uint(courseID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch questions", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Questions retrieved successfully", questions))
}

// CreateQuestionHandler handles POST /courses/:id/questions
func (h *handler) CreateQuestionHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	courseIDStr := c.Param("id")

	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Course ID", "Course ID must be a valid number", nil))
		return
	}

	var input CreateQuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	question, err := h.service.Create(uint(courseID), input, requestID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create question", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Question created successfully", question))
}

// UpdateQuestionHandler handles PUT /questions/:id
func (h *handler) UpdateQuestionHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	var input UpdateQuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	question, err := h.service.Update(uint(id), input, requestID, userID)
	if err != nil {
		if err.Error() == "question not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Question not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to update question", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Question updated successfully", question))
}

// DeleteQuestionHandler handles DELETE /questions/:id
func (h *handler) DeleteQuestionHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	if err := h.service.Delete(uint(id), requestID, userID); err != nil {
		if err.Error() == "question not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Question not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete question", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Question deleted successfully", nil))
}
