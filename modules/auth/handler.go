package auth

import (
	"context" // Add missing import
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOAuth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"github.com/redukasquad/be-reduka/packages/utils"
)


type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterHandler(c *gin.Context) {
	var input RegisterInput
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

func (h *Handler) LoginHandler(c *gin.Context) {
	var input LoginInput
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

func (h *Handler) VerifyEmailHandler(c *gin.Context) {
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

func (h *Handler) LogoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Logout Success", nil))
}

func (h *Handler) ForgotPasswordHandler(c *gin.Context) {
	var input ForgotPasswordInput
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

func (h *Handler) ResetPasswordHandler(c *gin.Context) {
	var input ResetPasswordInput
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

func getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			googleOAuth2.UserinfoEmailScope,
			googleOAuth2.UserinfoProfileScope,
		},
		Endpoint: google.Endpoint,
	}
}

func (h *Handler) GoogleLoginHandler(c *gin.Context) {
	googleConfig := getGoogleOauthConfig()
	url := googleConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleCallbackHandler(c *gin.Context) {
	if errParam := c.Query("error"); errParam != "" {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Google Login Failed", "Google returned error: "+errParam, nil))
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Google Login Failed", "Code not found. Make sure you are not visiting this URL directly.", nil))
		return
	}

	googleConfig := getGoogleOauthConfig()
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Google Login Failed", "Failed to exchange token: "+err.Error(), nil))
		return
	}

	client := googleConfig.Client(context.Background(), token)
	oauth2Service, err := googleOAuth2.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Google Login Failed", "Failed to create oauth service: "+err.Error(), nil))
		return
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Google Login Failed", "Failed to get user info: "+err.Error(), nil))
		return
	}

	user, jwtToken, err := h.service.LoginOrRegisterWithGoogle(userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Google Login Failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Google Login Success", gin.H{
		"token": jwtToken,
		"user":  user,
	}))
}

func (h *Handler) MeHandler(c *gin.Context) {
	// Get user_id from context (set by middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.BuildResponseFailed("Unauthorized", "User ID not found in context", nil))
		return
	}

	user, err := h.service.Me(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to get user profile", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("User Profile", user))
}
