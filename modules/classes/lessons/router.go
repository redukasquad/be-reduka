package lessons

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func LessonRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdminOrTutor gin.HandlerFunc) {
	lessonRepo := NewRepository(migrations.GetDB())
	lessonService := NewService(lessonRepo)
	lessonHandler := NewHandler(lessonService)

	subjectLessons := router.Group("/subjects/:id")
	{
		subjectLessons.GET("/lessons", lessonHandler.GetLessonsBySubjectHandler)
		subjectLessons.POST("/lessons", requireAuth, requireAdminOrTutor, lessonHandler.CreateLessonHandler)
	}

	lessons := router.Group("/lessons")
	{
		lessons.GET("/:id", lessonHandler.GetLessonByIDHandler)
		lessons.PUT("/:id", requireAuth, requireAdminOrTutor, lessonHandler.UpdateLessonHandler)
		lessons.DELETE("/:id", requireAuth, requireAdminOrTutor, lessonHandler.DeleteLessonHandler)
	}
}
