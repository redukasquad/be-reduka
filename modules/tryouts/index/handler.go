package tryouts

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
	// Try Out
	GetAllTryOutsHandler(c *gin.Context)
	GetTryOutByIDHandler(c *gin.Context)
	CreateTryOutHandler(c *gin.Context)
	UpdateTryOutHandler(c *gin.Context)
	DeleteTryOutHandler(c *gin.Context)

	// Tutor Permissions
	GetTutorPermissionsHandler(c *gin.Context)
	GrantTutorPermissionHandler(c *gin.Context)
	RevokeTutorPermissionHandler(c *gin.Context)
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

func isAdmin(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}
	roleStr, ok := role.(string)
	if !ok {
		return false
	}
	return roleStr == "ADMIN"
}

// ==========================================
// Try Out Handlers
// ==========================================

func (h *handler) GetAllTryOutsHandler(c *gin.Context) {
	requestID := getRequestID(c)

	var params dto.ListQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid query params", err.Error(), nil))
		return
	}

	tryOuts, err := h.service.GetAll(params, requestID, isAdmin(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch try outs", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Try outs retrieved successfully", tryOuts))
}

func (h *handler) GetTryOutByIDHandler(c *gin.Context) {
	requestID := getRequestID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	tryOut, err := h.service.GetByID(uint(id), requestID)
	if err != nil {
		if err.Error() == "try out not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Try out not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch try out", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Try out retrieved successfully", tryOut))
}

func (h *handler) CreateTryOutHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)

	var input CreateTryOutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	tryOut, err := h.service.Create(input, requestID, userID)
	if err != nil {
		if err.Error() == "try out with this name already exists" {
			c.JSON(http.StatusConflict, utils.BuildResponseFailed("Try out already exists", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create try out", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Try out created successfully", tryOut))
}

func (h *handler) UpdateTryOutHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	var input UpdateTryOutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	tryOut, err := h.service.Update(uint(id), input, requestID, userID)
	if err != nil {
		if err.Error() == "try out not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Try out not found", err.Error(), nil))
			return
		}
		if err.Error() == "try out with this name already exists" {
			c.JSON(http.StatusConflict, utils.BuildResponseFailed("Try out name conflict", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to update try out", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Try out updated successfully", tryOut))
}

func (h *handler) DeleteTryOutHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	if err := h.service.Delete(uint(id), requestID, userID); err != nil {
		if err.Error() == "try out not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Try out not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete try out", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Try out deleted successfully", nil))
}

// ==========================================
// Tutor Permission Handlers
// ==========================================

func (h *handler) GetTutorPermissionsHandler(c *gin.Context) {
	requestID := getRequestID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	permissions, err := h.service.GetTutorPermissions(uint(id), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch tutor permissions", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Tutor permissions retrieved successfully", permissions))
}

func (h *handler) GrantTutorPermissionHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a valid number", nil))
		return
	}

	var input GrantTutorPermissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	permission, err := h.service.GrantTutorPermission(uint(id), input, requestID, userID)
	if err != nil {
		if err.Error() == "try out not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Try out not found", err.Error(), nil))
			return
		}
		if err.Error() == "tutor already has permission for this try out" {
			c.JSON(http.StatusConflict, utils.BuildResponseFailed("Permission already exists", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to grant permission", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Tutor permission granted successfully", permission))
}

func (h *handler) RevokeTutorPermissionHandler(c *gin.Context) {
	requestID := getRequestID(c)
	userID := getUserID(c)
	idStr := c.Param("id")
	tutorIDStr := c.Param("userId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid Try Out ID", "ID must be a valid number", nil))
		return
	}

	tutorID, err := strconv.ParseUint(tutorIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid User ID", "User ID must be a valid number", nil))
		return
	}

	if err := h.service.RevokeTutorPermission(uint(id), uint(tutorID), requestID, userID); err != nil {
		if err.Error() == "tutor permission not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Permission not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to revoke permission", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Tutor permission revoked successfully", nil))
}
