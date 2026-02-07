package questions

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
	tryouts "github.com/redukasquad/be-reduka/modules/tryouts/index"
)

func QuestionRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	db := migrations.GetDB()
	repo := NewRepository(db)
	tryOutRepo := tryouts.NewRepository(db)
	service := NewService(repo, tryOutRepo)
	handler := NewHandler(service)

	// Admin/Tutor endpoints - require auth
	questionsAdmin := router.Group("/tryouts")
	questionsAdmin.Use(requireAuth)
	{
		// Get subtests with question count (for tutor dashboard)
		questionsAdmin.GET("/:id/subtests", handler.GetSubtestsHandler)

		// Get all questions for a try out
		questionsAdmin.GET("/:id/questions", handler.GetQuestionsByTryOutHandler)

		// Get questions by subtest
		questionsAdmin.GET("/:id/subtests/:subtestId/questions", handler.GetQuestionsBySubtestHandler)

		// Create question for a subtest
		questionsAdmin.POST("/:id/subtests/:subtestId/questions", handler.CreateQuestionHandler)
	}

	// Question-specific endpoints
	questionByID := router.Group("/tryouts/questions")
	questionByID.Use(requireAuth)
	{
		questionByID.PUT("/:questionId", handler.UpdateQuestionHandler)
		questionByID.DELETE("/:questionId", handler.DeleteQuestionHandler)
	}
}
