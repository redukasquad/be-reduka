package answers

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

// AnswerRouter registers all answer routes
func AnswerRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	answerRepo := NewRepository(migrations.GetDB())
	answerService := NewService(answerRepo)
	answerHandler := NewHandler(answerService)

	// Admin routes for viewing answers
	registrations := router.Group("/registrations")
	registrations.Use(requireAuth, requireAdmin)
	{
		registrations.GET("/:registrationId/answers", answerHandler.GetAnswersByRegistrationHandler)
	}
}
