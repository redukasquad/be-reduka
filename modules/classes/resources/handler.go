package resources

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
	GetResourcesByLessonHandler(c *gin.Context)
	CreateResourceHandler(c *gin.Context)
	UpdateResourceHandler(c *gin.Context)
	DeleteResourceHandler(c *gin.Context)
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

func (h *handler) GetResourcesByLessonHandler(c *gin.Context) {
	requestID := getRequestID(c)
	lessonIDStr := c.Param("id")

	lessonID, err := strconv.ParseUint(lessonIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Lesson ID", "Lesson ID must be a valid number", nil))
		return
	}

	resources, err := h.service.GetByLessonID(uint(lessonID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch resources", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Resources retrieved successfully", resources))
}

func (h *handler) CreateResourceHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	lessonIDStr := c.Param("id")

	lessonID, err := strconv.ParseUint(lessonIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Lesson ID", "Lesson ID must be a valid number", nil))
		return
	}

	var input CreateResourceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	resource, err := h.service.Create(uint(lessonID), input, requestID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create resource", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Resource created successfully", resource))
}

func (h *handler) UpdateResourceHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	var input UpdateResourceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	resource, err := h.service.Update(uint(id), input, requestID, userID)
	if err != nil {
		if err.Error() == "resource not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Resource not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to update resource", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Resource updated successfully", resource))
}

func (h *handler) DeleteResourceHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	if err := h.service.Delete(uint(id), requestID, userID); err != nil {
		if err.Error() == "resource not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Resource not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete resource", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Resource deleted successfully", nil))
}
