package questions

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

// QuestionRouter registers all question routes
func QuestionRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	questionRepo := NewRepository(migrations.GetDB())
	questionService := NewService(questionRepo)
	questionHandler := NewHandler(questionService)

	// Public routes - anyone can view questions for a course
	courses := router.Group("/courses")
	{
		courses.GET("/:courseId/questions", questionHandler.GetQuestionsByCourseHandler)
	}

	// Admin routes for managing questions
	adminCourses := router.Group("/courses")
	adminCourses.Use(requireAuth, requireAdmin)
	{
		adminCourses.POST("/:courseId/questions", questionHandler.CreateQuestionHandler)
	}

	// Admin routes for individual question management
	questions := router.Group("/questions")
	questions.Use(requireAuth, requireAdmin)
	{
		questions.PUT("/:id", questionHandler.UpdateQuestionHandler)
		questions.DELETE("/:id", questionHandler.DeleteQuestionHandler)
	}
}
