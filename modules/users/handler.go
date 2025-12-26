package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/packages/utils"
)

type handler struct {
	service Service
}

type Handler interface {
	GetAllUsersHandler(c *gin.Context)
	GetUserByIDHandler(c *gin.Context)
	UpdateUserHandler(c *gin.Context)
	SetRoleHandler(c *gin.Context)
	DeleteUserHandler(c *gin.Context)
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func (h *handler) GetAllUsersHandler(c *gin.Context) {
	users, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to get users", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Users retrieved", users))
}

func (h *handler) GetUserByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildResponseFailed("User not found", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("User retrieved", user))
}

func (h *handler) UpdateUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	user, err := h.service.Update(id, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Update failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("User updated", user))
}

// SetRoleHandler allows admin to set a user's role (ADMIN only)
func (h *handler) SetRoleHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}

	var input SetRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := []utils.ValidationError{
			{Field: "role", Error: "role is required and must be STUDENT, TUTOR, or ADMIN"},
		}
		c.JSON(http.StatusBadRequest, utils.BuildValidationErrorResponse("Validation failed", validationErrors))
		return
	}

	user, err := h.service.SetRole(id, input.Role)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("User not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to set role", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("User role updated", user))
}

func (h *handler) DeleteUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Delete failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("User deleted", nil))
}
