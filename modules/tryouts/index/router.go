package tryouts

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func TryOutIndexRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	repo := NewRepository(migrations.GetDB())
	service := NewService(repo)
	handler := NewHandler(service)

	// Public endpoints (anyone can view published try outs)
	tryouts := router.Group("/tryouts")
	{
		tryouts.GET("", handler.GetAllTryOutsHandler)
		tryouts.GET("/:id", handler.GetTryOutByIDHandler)
	}

	// Admin only endpoints
	tryoutsAdmin := router.Group("/tryouts")
	tryoutsAdmin.Use(requireAuth, requireAdmin)
	{
		tryoutsAdmin.POST("", handler.CreateTryOutHandler)
		tryoutsAdmin.PUT("/:id", handler.UpdateTryOutHandler)
		tryoutsAdmin.DELETE("/:id", handler.DeleteTryOutHandler)

		// Tutor permissions management
		tryoutsAdmin.GET("/:id/tutors", handler.GetTutorPermissionsHandler)
		tryoutsAdmin.POST("/:id/tutors", handler.GrantTutorPermissionHandler)
		tryoutsAdmin.DELETE("/:id/tutors/:userId", handler.RevokeTutorPermissionHandler)
	}
}
