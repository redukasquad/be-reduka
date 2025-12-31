package questions

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func QuestionRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	questionRepo := NewRepository(migrations.GetDB())
	questionService := NewService(questionRepo)
	questionHandler := NewHandler(questionService)

	courseQuestions := router.Group("/courses/:id")
	{
		courseQuestions.GET("/questions", questionHandler.GetQuestionsByCourseHandler)
		courseQuestions.POST("/questions", requireAuth, requireAdmin, questionHandler.CreateQuestionHandler)
	}

	questions := router.Group("/questions")
	questions.Use(requireAuth, requireAdmin)
	{
		questions.PUT("/:id", questionHandler.UpdateQuestionHandler)
		questions.DELETE("/:id", questionHandler.DeleteQuestionHandler)
	}
}
