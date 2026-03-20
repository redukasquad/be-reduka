package registrations

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func RegistrationRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	repo := NewRepository(migrations.GetDB())
	service := NewService(repo)
	handler := NewHandler(service)

	// User: register for a tryout
	userRoutes := router.Group("/tryouts")
	userRoutes.Use(requireAuth)
	{
		userRoutes.POST("/:id/register", handler.RegisterHandler)
	}

	// Get my registrations
	router.GET("/users/me/tryout-registrations", requireAuth, handler.GetMyRegistrationsHandler)

	// Admin: get registrations for a tryout
	adminRoutes := router.Group("/tryouts")
	adminRoutes.Use(requireAuth, requireAdmin)
	{
		adminRoutes.GET("/:id/registrations", handler.GetRegistrationsByTryOutHandler)
	}

	// User: upload payment proof — path /payment-proof/:id to avoid wildcard conflict with DELETE /:id
	proofRoutes := router.Group("/tryouts/registrations")
	proofRoutes.Use(requireAuth)
	{
		proofRoutes.POST("/payment-proof/:id", handler.UploadPaymentProofHandler)
	}

	// Admin: manage registrations — separate group, no wildcard conflict
	adminRegRoutes := router.Group("/tryouts/registrations")
	adminRegRoutes.Use(requireAuth, requireAdmin)
	{
		adminRegRoutes.GET("/pending", handler.GetPendingPaymentsHandler)
		adminRegRoutes.PUT("/:id/approve", handler.ApprovePaymentHandler)
		adminRegRoutes.PUT("/:id/reject", handler.RejectPaymentHandler)
		adminRegRoutes.DELETE("/:id", handler.DeleteRegistrationHandler)
	}
}
