package courses

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redukasquad/be-reduka/modules/dto"
	"github.com/redukasquad/be-reduka/packages/utils"
)

type handler struct {
	service Service
}

type Handler interface {
	GetAllCoursesHandler(c *gin.Context)
	GetCourseByIDHandler(c *gin.Context)
	GetCoursesByProgramIDHandler(c *gin.Context)
	CreateCourseHandler(c *gin.Context)
	UpdateCourseHandler(c *gin.Context)
	DeleteCourseHandler(c *gin.Context)
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
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

func (h *handler) GetAllCoursesHandler(c *gin.Context) {
	requestID := getRequestID(c)

	var params dto.ListQueryParams
	err := c.ShouldBindQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid query params", err.Error(), nil))
		return
	}

	courses, err := h.service.GetAll(params, requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch courses", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Courses retrieved successfully", courses))
}

func (h *handler) GetCourseByIDHandler(c *gin.Context) {
	requestID := getRequestID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	course, err := h.service.GetByID(uint(id), requestID)
	if err != nil {
		if err.Error() == "course not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Course not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch course", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Course retrieved successfully", course))
}

func (h *handler) GetCoursesByProgramIDHandler(c *gin.Context) {
	requestID := getRequestID(c)
	programIDStr := c.Param("id")

	programID, err := strconv.ParseUint(programIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Program ID", "Program ID must be a valid number", nil))
		return
	}

	courses, err := h.service.GetByProgramID(uint(programID), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch courses", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Courses retrieved successfully", courses))
}

func (h *handler) CreateCourseHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)

	var input CreateCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	course, err := h.service.Create(input, requestID, userID)
	if err != nil {
		if err.Error() == "course with this name already exists" {
			c.JSON(http.StatusConflict, utils.BuildResponseFailed("Course already exists", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create course", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Course created successfully", course))
}

func (h *handler) UpdateCourseHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	var input UpdateCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	course, err := h.service.Update(uint(id), input, requestID, userID)
	if err != nil {
		if err.Error() == "course not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Course not found", err.Error(), nil))
			return
		}
		if err.Error() == "course with this name already exists" {
			c.JSON(http.StatusConflict, utils.BuildResponseFailed("Course name conflict", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to update course", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Course updated successfully", course))
}

func (h *handler) DeleteCourseHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	if err := h.service.Delete(uint(id), requestID, userID); err != nil {
		if err.Error() == "course not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Course not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete course", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Course deleted successfully", nil))
}
