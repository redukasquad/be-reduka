package registrations

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func RegistrationRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	repo := NewRepository(migrations.GetDB())
	service := NewService(repo)
	handler := NewHandler(service)

	// User endpoints - require auth
	userRoutes := router.Group("/tryouts")
	userRoutes.Use(requireAuth)
	{
		// Register for a try out
		userRoutes.POST("/:id/register", handler.RegisterHandler)
	}

	// User registration management
	regRoutes := router.Group("/tryouts/registrations")
	regRoutes.Use(requireAuth)
	{
		regRoutes.POST("/:id/payment-proof", handler.UploadPaymentProofHandler)
	}

	// Get my registrations
	router.GET("/users/me/tryout-registrations", requireAuth, handler.GetMyRegistrationsHandler)

	// Admin endpoints
	adminRoutes := router.Group("/tryouts")
	adminRoutes.Use(requireAuth, requireAdmin)
	{
		// Get registrations for a try out
		adminRoutes.GET("/:id/registrations", handler.GetRegistrationsByTryOutHandler)
	}

	adminRegRoutes := router.Group("/tryouts/registrations")
	adminRegRoutes.Use(requireAuth, requireAdmin)
	{
		adminRegRoutes.GET("/pending", handler.GetPendingPaymentsHandler)
		adminRegRoutes.PUT("/:id/approve", handler.ApprovePaymentHandler)
		adminRegRoutes.PUT("/:id/reject", handler.RejectPaymentHandler)
	}
}
