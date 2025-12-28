package registrations

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

// RegistrationRouter registers all registration routes
func RegistrationRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	regRepo := NewRepository(migrations.GetDB())
	regService := NewService(regRepo)
	regHandler := NewHandler(regService)

	// User registration routes (requires auth)
	courses := router.Group("/courses")
	{
		courses.POST("/:courseId/register", requireAuth, regHandler.RegisterHandler)
		courses.GET("/registrations/me", requireAuth, regHandler.GetMyRegistrationsHandler)
	}

	// Admin routes for managing registrations
	adminCourses := router.Group("/courses")
	adminCourses.Use(requireAuth, requireAdmin)
	{
		adminCourses.GET("/:courseId/registrations", regHandler.GetRegistrationsByCourseHandler)
	}

	// Admin routes for approval/rejection
	registrations := router.Group("/registrations")
	registrations.Use(requireAuth, requireAdmin)
	{
		registrations.PUT("/:id/approve", regHandler.ApproveRegistrationHandler)
		registrations.PUT("/:id/reject", regHandler.RejectRegistrationHandler)
	}
}
