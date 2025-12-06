package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/packages/dto"
	"github.com/redukasquad/be-reduka/packages/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Register Failed", err.Error(), nil))
		return
	}

	user, err := h.service.Register(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Register Failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("User Registered", user))
}

func (h *Handler) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Login Failed", err.Error(), nil))
		return
	}

	token, err := h.service.Login(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Login Failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Login Success", gin.H{"token": token}))
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Verification Failed", "Code is required", nil))
		return
	}

	err := h.service.VerifyEmail(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Verification Failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Email Verified Successfully", nil))
}

func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Logout Success", nil))
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var input dto.ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Request Failed", err.Error(), nil))
		return
	}

	err := h.service.ForgotPassword(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Request Failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Reset Password Code Sent", nil))
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var input dto.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Reset Failed", err.Error(), nil))
		return
	}

	err := h.service.ResetPassword(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Reset Failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Password Reset Successfully", nil))
}
