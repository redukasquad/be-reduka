package auth

import (
	"context" // Add missing import
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/packages/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOAuth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type handler struct {
	service Service
}

type Handler interface {
	RegisterHandler(c *gin.Context)
	LoginHandler(c *gin.Context)
	VerifyEmailHandler(c *gin.Context)
	ResendVerificationHandler(c *gin.Context)
	LogoutHandler(c *gin.Context)
	ForgotPasswordHandler(c *gin.Context)
	ResetPasswordHandler(c *gin.Context)
	GoogleLoginHandler(c *gin.Context)
	GoogleCallbackHandler(c *gin.Context)
	MeHandler(c *gin.Context)
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func (h *handler) RegisterHandler(c *gin.Context) {
	var input RegisterInput
	var validationErrors []utils.ValidationError

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "unknown field") {
			fieldName := strings.TrimPrefix(errStr, "json: unknown field ")
			fieldName = strings.Trim(fieldName, "\"")
			validationErrors = append(validationErrors, utils.ValidationError{
				Field: fieldName,
				Error: "field is not allowed in registration",
			})
			c.JSON(http.StatusBadRequest, utils.BuildValidationErrorResponse("Invalid request body", validationErrors))
			return
		}
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid JSON format", err.Error(), nil))
		return
	}

	if input.Username == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "username",
			Error: "username cannot be empty",
		})
	}
	if input.Email == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "email",
			Error: "email cannot be empty",
		})
	}
	if input.Password == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "password",
			Error: "password cannot be empty",
		})
	}

	if input.Email != "" && !strings.Contains(input.Email, "@") {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "email",
			Error: "email format is invalid",
		})
	}

	if input.Password != "" && len(input.Password) < 6 {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "password",
			Error: "password must be at least 6 characters",
		})
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, utils.BuildValidationErrorResponse("Validation failed", validationErrors))
		return
	}

	user, err := h.service.Register(input)
	if err != nil {
		if strings.Contains(err.Error(), "email already registered") {
			validationErrors = append(validationErrors, utils.ValidationError{
				Field: "email",
				Error: "email is already registered",
			})
			c.JSON(http.StatusConflict, utils.BuildValidationErrorResponse("Registration failed", validationErrors))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Registration failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("User registered successfully", UserResponseJSON(user)))
}

func (h *handler) LoginHandler(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := []utils.ValidationError{}
		if strings.Contains(err.Error(), "Email") {
			validationErrors = append(validationErrors, utils.ValidationError{Field: "email", Error: "email is required and must be valid"})
		}
		if strings.Contains(err.Error(), "Password") {
			validationErrors = append(validationErrors, utils.ValidationError{Field: "password", Error: "password is required"})
		}
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, utils.BuildValidationErrorResponse("Validation failed", validationErrors))
			return
		}
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Login failed", err.Error(), nil))
		return
	}

	token, err := h.service.Login(input)
	if err != nil {
		if strings.Contains(err.Error(), "email not verified") {
			c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Email not verified", "Please verify your email before logging in. Use POST /api/v1/auth/resend-verification to get a new code.", nil))
			return
		}
		if strings.Contains(err.Error(), "this account uses Google login") {
			c.JSON(http.StatusForbidden, utils.BuildResponseFailed("Wrong login method", "This account was registered with Google. Please use Google login instead.", nil))
			return
		}
		if strings.Contains(err.Error(), "invalid email or password") {
			validationErrors := []utils.ValidationError{
				{Field: "email", Error: "invalid email or password"},
			}
			c.JSON(http.StatusUnauthorized, utils.BuildValidationErrorResponse("Login failed", validationErrors))
			return
		}
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Login failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Login successful", gin.H{"token": token}))
}

func (h *handler) VerifyEmailHandler(c *gin.Context) {
	var input VerifyEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := []utils.ValidationError{}
		if strings.Contains(err.Error(), "Email") {
			validationErrors = append(validationErrors, utils.ValidationError{Field: "email", Error: "email is required and must be valid"})
		}
		if strings.Contains(err.Error(), "Code") {
			validationErrors = append(validationErrors, utils.ValidationError{Field: "code", Error: "verification code is required"})
		}
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, utils.BuildValidationErrorResponse("Validation failed", validationErrors))
			return
		}
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Verification failed", err.Error(), nil))
		return
	}

	err := h.service.VerifyEmail(input.Email, input.Code)
	if err != nil {
		if strings.Contains(err.Error(), "email not found") {
			validationErrors := []utils.ValidationError{
				{Field: "email", Error: "email is not registered"},
			}
			c.JSON(http.StatusNotFound, utils.BuildValidationErrorResponse("Email not found", validationErrors))
			return
		}
		if strings.Contains(err.Error(), "email already verified") {
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Already verified", "email is already verified", nil))
			return
		}
		if strings.Contains(err.Error(), "invalid verification code") {
			validationErrors := []utils.ValidationError{
				{Field: "code", Error: "invalid verification code"},
			}
			c.JSON(http.StatusBadRequest, utils.BuildValidationErrorResponse("Verification failed", validationErrors))
			return
		}
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Verification failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Email verified successfully", nil))
}

func (h *handler) ResendVerificationHandler(c *gin.Context) {
	var input ResendVerificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := []utils.ValidationError{
			{Field: "email", Error: "email is required and must be valid"},
		}
		c.JSON(http.StatusBadRequest, utils.BuildValidationErrorResponse("Validation failed", validationErrors))
		return
	}

	err := h.service.ResendVerificationCode(input.Email)
	if err != nil {
		if strings.Contains(err.Error(), "email not found") {
			validationErrors := []utils.ValidationError{
				{Field: "email", Error: "email is not registered"},
			}
			c.JSON(http.StatusNotFound, utils.BuildValidationErrorResponse("Email not found", validationErrors))
			return
		}
		if strings.Contains(err.Error(), "email already verified") {
			c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Already verified", "email is already verified", nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to send verification code", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Verification code sent to your email", nil))
}

func (h *handler) LogoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Logout Success", nil))
}

func (h *handler) ForgotPasswordHandler(c *gin.Context) {
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

func (h *handler) ResetPasswordHandler(c *gin.Context) {
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

func (h *handler) GoogleLoginHandler(c *gin.Context) {
	googleConfig := getGoogleOauthConfig()
	url := googleConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *handler) GoogleCallbackHandler(c *gin.Context) {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	if errParam := c.Query("error"); errParam != "" {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/auth/google/error?error="+errParam)
		return
	}

	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/auth/google/error?error=code_not_found")
		return
	}

	googleConfig := getGoogleOauthConfig()
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/auth/google/error?error=token_exchange_failed")
		return
	}

	client := googleConfig.Client(context.Background(), token)
	oauth2Service, err := googleOAuth2.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/auth/google/error?error=oauth_service_failed")
		return
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/auth/google/error?error=user_info_failed")
		return
	}

	_, jwtToken, err := h.service.LoginOrRegisterWithGoogle(userInfo)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/auth/google/error?error=login_failed")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/auth/google/success?token="+jwtToken)
}

func (h *handler) MeHandler(c *gin.Context) {
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

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("User Profile", UserResponseJSON(*user)))
}
