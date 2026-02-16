package uploads

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redukasquad/be-reduka/packages/utils"
)

type handler struct {
	service Service
}

type Handler interface {
	CreateImageHandler(c *gin.Context)
	DeleteImageHandler(c *gin.Context)
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

func (h *handler) CreateImageHandler(c *gin.Context) {
	requestID := getRequestID(c)

	var input CreateImageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	image, err := h.service.Create(input, requestID)
	if err != nil {
		if err.Error() == "image with this URL already exists" {
			c.JSON(http.StatusConflict, utils.BuildResponseFailed("Image already exists", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create image", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Image created successfully", image))
}

func (h *handler) DeleteImageHandler(c *gin.Context) {
	requestID := getRequestID(c)

	var input DeleteImageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid input", err.Error(), nil))
		return
	}

	fileId, err := h.service.DeleteByURL(input.URL, requestID)
	if err != nil {
		if err.Error() == "image not found" {
			c.JSON(http.StatusNotFound, utils.BuildResponseFailed("Image not found", err.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete image", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Image deleted successfully", gin.H{
		"fileId": fileId,
	}))
}
