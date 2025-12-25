package users

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func UserRouter(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	userRepo := NewRepository(migrations.GetDB())
	userService := NewService(userRepo)
	userHandler := NewHandler(userService)

	users := router.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("", userHandler.GetAllUsersHandler)
		users.GET("/:id", userHandler.GetUserByIDHandler)
		users.PUT("/:id", userHandler.UpdateUserHandler)
		users.DELETE("/:id", userHandler.DeleteUserHandler)
	}
}
