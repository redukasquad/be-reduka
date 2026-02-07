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

type Handler interface {
	GetSubtestsHandler(c *gin.Context)
	GetQuestionsByTryOutHandler(c *gin.Context)
	GetQuestionsBySubtestHandler(c *gin.Context)
	CreateQuestionHandler(c *gin.Context)
	UpdateQuestionHandler(c *gin.Context)
	DeleteQuestionHandler(c *gin.Context)
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
// Handlers
// ==========================================

func (h *handler) GetSubtestsHandler(c *gin.Context) {
	requestID := getRequestID(c)
	tryOutIDStr := c.Param("id")

	tryOutID, err := strconv.ParseUint(tryOutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	subtests, err := h.service.GetSubtestsWithQuestionCount(uint(tryOutID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch subtests", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Subtests retrieved successfully", subtests))
}

func (h *handler) GetQuestionsByTryOutHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	tryOutIDStr := c.Param("id")

	tryOutID, err := strconv.ParseUint(tryOutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	// Check permission
	if err := h.service.CheckTutorPermission(uint(tryOutID), userID); err != nil {
		c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		return
	}

	questions, err := h.service.GetQuestionsByTryOut(uint(tryOutID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch questions", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Questions retrieved successfully", questions))
}

func (h *handler) GetQuestionsBySubtestHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	tryOutIDStr := c.Param("id")
	subtestIDStr := c.Param("subtestId")

	tryOutID, err := strconv.ParseUint(tryOutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Try Out ID", "ID must be a valid number", nil))
		return
	}

	subtestID, err := strconv.ParseUint(subtestIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Subtest ID", "ID must be a valid number", nil))
		return
	}

	// Check permission
	if err := h.service.CheckTutorPermission(uint(tryOutID), userID); err != nil {
		c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		return
	}

	questions, err := h.service.GetQuestionsBySubtest(uint(tryOutID), uint(subtestID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch questions", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Questions retrieved successfully", questions))
}

func (h *handler) CreateQuestionHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	tryOutIDStr := c.Param("id")
	subtestIDStr := c.Param("subtestId")

	tryOutID, err := strconv.ParseUint(tryOutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Try Out ID", "ID must be a valid number", nil))
		return
	}

	subtestID, err := strconv.ParseUint(subtestIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Subtest ID", "ID must be a valid number", nil))
		return
	}

	// Check permission
	if err := h.service.CheckTutorPermission(uint(tryOutID), userID); err != nil {
		c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		return
	}

	var input CreateQuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	question, err := h.service.CreateQuestion(uint(tryOutID), uint(subtestID), input, requestID, userID)
	if err != nil {
		if err.Error() == "try out not found" || err.Error() == "subtest not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Failed to create question", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Question created successfully", question))
}

func (h *handler) UpdateQuestionHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	questionIDStr := c.Param("questionId")

	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Question ID", "ID must be a valid number", nil))
		return
	}

	// Get question to check try out permission
	existingQuestion, err := h.service.GetQuestionByID(uint(questionID), requestID)
	if err != nil {
		if err.Error() == "question not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Question not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch question", err.Error(), nil))
		return
	}

	// Check permission
	if err := h.service.CheckTutorPermission(existingQuestion.TryOutID, userID); err != nil {
		c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		return
	}

	var input UpdateQuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	question, err := h.service.UpdateQuestion(uint(questionID), input, requestID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to update question", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Question updated successfully", question))
}

func (h *handler) DeleteQuestionHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	questionIDStr := c.Param("questionId")

	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Question ID", "ID must be a valid number", nil))
		return
	}

	// Get question to check try out permission
	existingQuestion, err := h.service.GetQuestionByID(uint(questionID), requestID)
	if err != nil {
		if err.Error() == "question not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Question not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch question", err.Error(), nil))
		return
	}

	// Check permission
	if err := h.service.CheckTutorPermission(existingQuestion.TryOutID, userID); err != nil {
		c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		return
	}

	if err := h.service.DeleteQuestion(uint(questionID), requestID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete question", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Question deleted successfully", nil))
}
