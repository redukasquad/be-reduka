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

	// User routes - use /registrations path to avoid conflict with /courses/:id
	registrations := router.Group("/registrations")
	{
		// User routes
		registrations.GET("/me", requireAuth, regHandler.GetMyRegistrationsHandler)

		// Admin routes for approval/rejection
		registrations.PUT("/:id/approve", requireAuth, requireAdmin, regHandler.ApproveRegistrationHandler)
		registrations.PUT("/:id/reject", requireAuth, requireAdmin, regHandler.RejectRegistrationHandler)
	}

	// Course-specific registration routes - nested under courses with specific path
	// Using /courses/:id/register and /courses/:id/registrations
	courseRegistrations := router.Group("/courses/:id")
	{
		courseRegistrations.POST("/register", requireAuth, regHandler.RegisterHandler)
		courseRegistrations.GET("/registrations", requireAuth, requireAdmin, regHandler.GetRegistrationsByCourseHandler)
	}
}
