package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
	"github.com/redukasquad/be-reduka/modules/users"

	"github.com/redukasquad/be-reduka/middleware"
)

func AuthRouter(router *gin.RouterGroup) {
	userRepo := users.NewRepository(migrations.GetDB())
	authService := NewService(userRepo)
	authHandler := NewHandler(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.RegisterHandler)
		auth.POST("/login", authHandler.LoginHandler)
		auth.POST("/verify-email", authHandler.VerifyEmailHandler)
		auth.POST("/resend-verification", authHandler.ResendVerificationHandler)
		auth.POST("/logout", authHandler.LogoutHandler)
		auth.POST("/forgot-password", authHandler.ForgotPasswordHandler)
		auth.POST("/reset-password", authHandler.ResetPasswordHandler)
		auth.GET("/google/login", authHandler.GoogleLoginHandler)
		auth.GET("/google/callback", authHandler.GoogleCallbackHandler)

		auth.GET("/me", middleware.RequireAuth(), authHandler.MeHandler)
	}
}
