package lessons

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
	GetLessonsBySubjectHandler(c *gin.Context)
	GetLessonByIDHandler(c *gin.Context)
	CreateLessonHandler(c *gin.Context)
	UpdateLessonHandler(c *gin.Context)
	DeleteLessonHandler(c *gin.Context)
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

func (h *handler) GetLessonsBySubjectHandler(c *gin.Context) {
	requestID := getRequestID(c)
	subjectIDStr := c.Param("id")

	subjectID, err := strconv.ParseUint(subjectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Subject ID", "Subject ID must be a valid number", nil))
		return
	}

	lessons, err := h.service.GetBySubjectID(uint(subjectID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch lessons", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Lessons retrieved successfully", lessons))
}

func (h *handler) GetLessonByIDHandler(c *gin.Context) {
	requestID := getRequestID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	lesson, err := h.service.GetByID(uint(id), requestID)
	if err != nil {
		if err.Error() == "lesson not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Lesson not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch lesson", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Lesson retrieved successfully", lesson))
}

func (h *handler) CreateLessonHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	subjectIDStr := c.Param("id")

	subjectID, err := strconv.ParseUint(subjectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Subject ID", "Subject ID must be a valid number", nil))
		return
	}

	var input CreateLessonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	lesson, err := h.service.Create(uint(subjectID), input, requestID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create lesson", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Lesson created successfully", lesson))
}

func (h *handler) UpdateLessonHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	var input UpdateLessonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	lesson, err := h.service.Update(uint(id), input, requestID, userID)
	if err != nil {
		if err.Error() == "lesson not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Lesson not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to update lesson", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Lesson updated successfully", lesson))
}

func (h *handler) DeleteLessonHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	if err := h.service.Delete(uint(id), requestID, userID); err != nil {
		if err.Error() == "lesson not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Lesson not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete lesson", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Lesson deleted successfully", nil))
}
