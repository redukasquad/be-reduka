package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/internal/domain"
	"github.com/redukasquad/be-reduka/pkg/response"
)

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
}

func NewAuthHandler(r *gin.Engine, us domain.AuthUsecase) {
	handler := &AuthHandler{
		AuthUsecase: us,
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.GET("/google/login", handler.GoogleLogin)
		auth.GET("/google/callback", handler.GoogleCallback)
		auth.GET("/verify-email", handler.VerifyEmail)
		auth.POST("/logout", handler.Logout)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := h.AuthUsecase.Register(c.Request.Context(), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to register", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully. Please check your email for verification.", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	res, err := h.AuthUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid credentials", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url, err := h.AuthUsecase.GoogleLogin(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate Google login URL", err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.Error(c, http.StatusBadRequest, "Code is required", nil)
		return
	}

	res, err := h.AuthUsecase.GoogleCallback(c.Request.Context(), code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to login with Google", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Google login successful", res)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.Error(c, http.StatusBadRequest, "Token is required", nil)
		return
	}

	if err := h.AuthUsecase.VerifyEmail(c.Request.Context(), token); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to verify email", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Email verified successfully", nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Since we are using JWT, logout is handled client-side by removing the token.
	// We can blacklist the token here if we implement a blacklist mechanism (e.g., Redis).
	// For now, just return success.
	response.Success(c, http.StatusOK, "Logout successful", nil)
}
