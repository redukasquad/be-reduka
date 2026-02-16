package answers

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func AnswerRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	answerRepo := NewRepository(migrations.GetDB())
	answerService := NewService(answerRepo)
	answerHandler := NewHandler(answerService)

	registrations := router.Group("/registrations")
	registrations.Use(requireAuth, requireAdmin)
	{
		registrations.GET("/:registrationId/answers", answerHandler.GetAnswersByRegistrationHandler)
		registrations.POST("/:registrationId/answers", answerHandler.CreateAnswerHandler)
		registrations.DELETE("/answers/:id", answerHandler.DeleteAnswerHandler)
	}
}
