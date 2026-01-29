package programs

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

// Handler interface defines the HTTP handlers for programs
type Handler interface {
	GetAllProgramsHandler(c *gin.Context)
	GetProgramByIDHandler(c *gin.Context)
	CreateProgramHandler(c *gin.Context)
	UpdateProgramHandler(c *gin.Context)
	DeleteProgramHandler(c *gin.Context)
}

// NewHandler creates a new program handler
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

// GetAllProgramsHandler handles GET /programs
func (h *handler) GetAllProgramsHandler(c *gin.Context) {
	requestID := getRequestID(c)

	programs, err := h.service.GetAll(requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch programs", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Programs retrieved successfully", programs))
}

// GetProgramByIDHandler handles GET /programs/:id
func (h *handler) GetProgramByIDHandler(c *gin.Context) {
	requestID := getRequestID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	program, err := h.service.GetByID(uint(id), requestID)
	if err != nil {
		if err.Error() == "program not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Program not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch program", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Program retrieved successfully", program))
}

// CreateProgramHandler handles POST /programs
func (h *handler) CreateProgramHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)

	var input CreateProgramInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	program, err := h.service.Create(input, requestID, userID)
	if err != nil {
		if err.Error() == "program with this name already exists" {
			c.JSON(http.StatusConflict, utils.BuildResponseFailed("Program already exists", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create program", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Program created successfully", program))
}

// UpdateProgramHandler handles PUT /programs/:id
func (h *handler) UpdateProgramHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	var input UpdateProgramInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	program, err := h.service.Update(uint(id), input, requestID, userID)
	if err != nil {
		if err.Error() == "program not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Program not found", err.Error(), nil))
			return
		}
		if err.Error() == "program with this name already exists" {
			c.JSON(http.StatusConflict, utils.BuildResponseFailed("Program name conflict", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to update program", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Program updated successfully", program))
}

// DeleteProgramHandler handles DELETE /programs/:id
func (h *handler) DeleteProgramHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	if err := h.service.Delete(uint(id), requestID, userID); err != nil {
		if err.Error() == "program not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Program not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete program", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Program deleted successfully", nil))
}
