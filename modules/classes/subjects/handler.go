package subjects

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
	GetSubjectsByCourseHandler(c *gin.Context)
	GetSubjectByIDHandler(c *gin.Context)
	CreateSubjectHandler(c *gin.Context)
	UpdateSubjectHandler(c *gin.Context)
	DeleteSubjectHandler(c *gin.Context)
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

func (h *handler) GetSubjectsByCourseHandler(c *gin.Context) {
	requestID := getRequestID(c)
	courseIDStr := c.Param("id")

	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Course ID", "Course ID must be a valid number", nil))
		return
	}

	subjects, err := h.service.GetByCourseID(uint(courseID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch subjects", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Subjects retrieved successfully", subjects))
}

func (h *handler) GetSubjectByIDHandler(c *gin.Context) {
	requestID := getRequestID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	subject, err := h.service.GetByID(uint(id), requestID)
	if err != nil {
		if err.Error() == "subject not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Subject not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch subject", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Subject retrieved successfully", subject))
}

func (h *handler) CreateSubjectHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	courseIDStr := c.Param("id")

	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Course ID", "Course ID must be a valid number", nil))
		return
	}

	var input CreateSubjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	subject, err := h.service.Create(uint(courseID), input, requestID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create subject", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Subject created successfully", subject))
}

func (h *handler) UpdateSubjectHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	var input UpdateSubjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	subject, err := h.service.Update(uint(id), input, requestID, userID)
	if err != nil {
		if err.Error() == "subject not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Subject not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to update subject", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Subject updated successfully", subject))
}

func (h *handler) DeleteSubjectHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	if err := h.service.Delete(uint(id), requestID, userID); err != nil {
		if err.Error() == "subject not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Subject not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete subject", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Subject deleted successfully", nil))
}
