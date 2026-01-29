package users

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

// UserRouter sets up user routes. Pass RequireAuth and RequireAdmin middleware from caller.
func UserRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	userRepo := NewRepository(migrations.GetDB())
	userService := NewService(userRepo)
	userHandler := NewHandler(userService)

	users := router.Group("/users")
	users.Use(requireAuth)
	{
		// All authenticated users can access these
		users.GET("", userHandler.GetAllUsersHandler)
		users.GET("/:id", userHandler.GetUserByIDHandler)
		users.PUT("/:id", userHandler.UpdateUserHandler)

		// Admin only routes
		users.PATCH("/:id/role", requireAdmin, userHandler.SetRoleHandler)
		users.DELETE("/:id", requireAdmin, userHandler.DeleteUserHandler)
	}
}
