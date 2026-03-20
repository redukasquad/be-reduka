package universities

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
	GetAllUniversitiesHandler(c *gin.Context)
	GetUniversityByIDHandler(c *gin.Context)
	CreateUniversityHandler(c *gin.Context)
	UpdateUniversityHandler(c *gin.Context)
	DeleteUniversityHandler(c *gin.Context)

	GetMajorsByUniversityHandler(c *gin.Context)
	CreateMajorHandler(c *gin.Context)
	UpdateMajorHandler(c *gin.Context)
	DeleteMajorHandler(c *gin.Context)

	GetMyTargetsHandler(c *gin.Context)
	AddTargetHandler(c *gin.Context)
	DeleteTargetHandler(c *gin.Context)

	GetUsersByUniversityHandler(c *gin.Context)
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func (h *handler) GetAllUniversitiesHandler(c *gin.Context) {
	search := c.Query("q")
	unis, err := h.service.GetAllUniversities(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch universities", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Universities retrieved", unis))
}

func (h *handler) GetUniversityByIDHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}
	uni, err := h.service.GetUniversityByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Not found", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("University retrieved", uni))
}

func (h *handler) CreateUniversityHandler(c *gin.Context) {
	var input CreateUniversityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}
	uni, err := h.service.CreateUniversity(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create", err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("University created", uni))
}

func (h *handler) UpdateUniversityHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}
	var input UpdateUniversityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}
	uni, err := h.service.UpdateUniversity(uint(id), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Update failed", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("University updated", uni))
}

func (h *handler) DeleteUniversityHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}
	if err := h.service.DeleteUniversity(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Delete failed", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("University deleted", nil))
}

func (h *handler) GetMajorsByUniversityHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}
	majors, err := h.service.GetMajorsByUniversity(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch majors", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Majors retrieved", majors))
}

func (h *handler) CreateMajorHandler(c *gin.Context) {
	uniID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || uniID == 0 {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid university ID", "University ID must be a valid non-zero number", nil))
		return
	}
	var input CreateMajorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}
	// Always use the URL param — ignore any body universityId
	input.UniversityID = uint(uniID)
	major, err := h.service.CreateMajor(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create major", err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Major created", major))
}

func (h *handler) UpdateMajorHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}
	var input UpdateMajorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}
	major, err := h.service.UpdateMajor(uint(id), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Update failed", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Major updated", major))
}

func (h *handler) DeleteMajorHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}
	if err := h.service.DeleteMajor(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Delete failed", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Major deleted", nil))
}

func (h *handler) GetMyTargetsHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")
	targets, err := h.service.GetMyTargets(uint(userID.(int)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch targets", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Targets retrieved", targets))
}

func (h *handler) AddTargetHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var input SetUserTargetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}
	target, err := h.service.AddTarget(uint(userID.(int)), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Failed to add target", err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Target added", target))
}

func (h *handler) DeleteTargetHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}
	if err := h.service.DeleteTarget(uint(id), uint(userID.(int))); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Delete failed", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Target removed", nil))
}

func (h *handler) GetUsersByUniversityHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}
	users, err := h.service.GetUsersByUniversity(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to fetch users", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Users retrieved", users))
}
