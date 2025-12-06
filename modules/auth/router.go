package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
	"github.com/redukasquad/be-reduka/modules/users"
)

func AuthRouter(router *gin.RouterGroup) {
	userRepo := users.NewRepository(migrations.GetDB())
	authService := NewService(userRepo)
	authHandler := NewHandler(authService)

	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.GET("/verify-email", authHandler.VerifyEmail)
	router.POST("/logout", authHandler.Logout)
}
