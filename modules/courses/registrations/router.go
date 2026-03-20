package registrations

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func RegistrationRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc, requireAdminOrTutor gin.HandlerFunc) {
	regRepo := NewRepository(migrations.GetDB())
	regService := NewService(regRepo)
	regHandler := NewHandler(regService)

	registrations := router.Group("/registrations")
	{
		registrations.GET("/me", requireAuth, regHandler.GetMyRegistrationsHandler)
		registrations.PUT("/:id/approve", requireAuth, requireAdminOrTutor, regHandler.ApproveRegistrationHandler)
		registrations.PUT("/:id/reject", requireAuth, requireAdminOrTutor, regHandler.RejectRegistrationHandler)
	}

	courseRegistrations := router.Group("/courses/:id")
	{
		courseRegistrations.POST("/register", requireAuth, regHandler.RegisterHandler)
		courseRegistrations.GET("/registrations", requireAuth, requireAdminOrTutor, regHandler.GetRegistrationsByCourseHandler)
	}
}
