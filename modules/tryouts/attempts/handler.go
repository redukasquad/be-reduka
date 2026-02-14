package attempts

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
	StartAttemptHandler(c *gin.Context)
	GetCurrentStateHandler(c *gin.Context)
	StartSubtestHandler(c *gin.Context)
	SubmitSubtestHandler(c *gin.Context)
	FinishAttemptHandler(c *gin.Context)
	GetResultsHandler(c *gin.Context)
	GetLeaderboardHandler(c *gin.Context)
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

func (h *handler) StartAttemptHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	regIDStr := c.Param("id")

	regID, err := strconv.ParseUint(regIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Registration ID", "ID must be a valid number", nil))
		return
	}

	attempt, err := h.service.StartAttempt(uint(regID), userID, requestID)
	if err != nil {
		switch err.Error() {
		case "registration not found":
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Registration not found", err.Error(), nil))
		case "you can only start your own registration":
			c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		case "payment must be approved before starting":
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Payment required", err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to start attempt", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Attempt started successfully", attempt))
}

func (h *handler) GetCurrentStateHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	attemptIDStr := c.Param("attemptId")

	attemptID, err := strconv.ParseUint(attemptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Attempt ID", "ID must be a valid number", nil))
		return
	}

	state, err := h.service.GetCurrentState(uint(attemptID), userID, requestID)
	if err != nil {
		switch err.Error() {
		case "attempt not found":
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Attempt not found", err.Error(), nil))
		case "you can only view your own attempt":
			c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to get current state", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Current state retrieved", state))
}

func (h *handler) StartSubtestHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	attemptIDStr := c.Param("attemptId")
	subtestIDStr := c.Param("subtestId")

	attemptID, err := strconv.ParseUint(attemptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Attempt ID", "ID must be a valid number", nil))
		return
	}

	subtestID, err := strconv.ParseUint(subtestIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Subtest ID", "ID must be a valid number", nil))
		return
	}

	questions, err := h.service.StartSubtest(uint(attemptID), uint(subtestID), userID, requestID)
	if err != nil {
		switch err.Error() {
		case "attempt not found", "subtest not found":
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Not found", err.Error(), nil))
		case "you can only access your own attempt":
			c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		case "attempt is not in progress":
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid state", err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to start subtest", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Subtest started", questions))
}

func (h *handler) SubmitSubtestHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	attemptIDStr := c.Param("attemptId")
	subtestIDStr := c.Param("subtestId")

	attemptID, err := strconv.ParseUint(attemptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Attempt ID", "ID must be a valid number", nil))
		return
	}

	subtestID, err := strconv.ParseUint(subtestIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Subtest ID", "ID must be a valid number", nil))
		return
	}

	var input SubmitSubtestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	result, err := h.service.SubmitSubtest(uint(attemptID), uint(subtestID), input, userID, requestID)
	if err != nil {
		switch err.Error() {
		case "attempt not found":
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Attempt not found", err.Error(), nil))
		case "you can only submit your own answers":
			c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		case "attempt is not in progress", "subtest not started", "subtest already submitted":
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid state", err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to submit subtest", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Subtest submitted successfully", result))
}

func (h *handler) FinishAttemptHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	attemptIDStr := c.Param("attemptId")

	attemptID, err := strconv.ParseUint(attemptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Attempt ID", "ID must be a valid number", nil))
		return
	}

	result, err := h.service.FinishAttempt(uint(attemptID), userID, requestID)
	if err != nil {
		switch err.Error() {
		case "attempt not found":
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Attempt not found", err.Error(), nil))
		case "you can only finish your own attempt":
			c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		case "attempt is already completed":
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid state", err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to finish attempt", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Attempt finished successfully", result))
}

func (h *handler) GetResultsHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	attemptIDStr := c.Param("attemptId")

	attemptID, err := strconv.ParseUint(attemptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Attempt ID", "ID must be a valid number", nil))
		return
	}

	result, err := h.service.GetAttemptResults(uint(attemptID), userID, requestID)
	if err != nil {
		switch err.Error() {
		case "attempt not found":
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Attempt not found", err.Error(), nil))
		case "you can only view your own results":
			c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Permission denied", err.Error(), nil))
		case "attempt is not completed yet":
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid state", err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to get results", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Results retrieved successfully", result))
}

func (h *handler) GetLeaderboardHandler(c *gin.Context) {
	requestID := getRequestID(c)
	tryOutIDStr := c.Param("id")

	tryOutID, err := strconv.ParseUint(tryOutIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Try Out ID", "ID must be a valid number", nil))
		return
	}

	leaderboard, err := h.service.GetLeaderboard(uint(tryOutID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to get leaderboard", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Leaderboard retrieved successfully", leaderboard))
}
