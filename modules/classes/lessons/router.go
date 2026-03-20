package lessons

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func LessonRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdminOrTutor gin.HandlerFunc) {
	lessonRepo := NewRepository(migrations.GetDB())
	lessonService := NewService(lessonRepo)
	lessonHandler := NewHandler(lessonService)

	classLessons := router.Group("/classes/:id")
	{
		classLessons.GET("/lessons", lessonHandler.GetLessonsByClassHandler)
		classLessons.POST("/lessons", requireAuth, requireAdminOrTutor, lessonHandler.CreateLessonHandler)
	}

	lessons := router.Group("/lessons")
	{
		lessons.GET("/:id", lessonHandler.GetLessonByIDHandler)
		lessons.PUT("/:id", requireAuth, requireAdminOrTutor, lessonHandler.UpdateLessonHandler)
		lessons.DELETE("/:id", requireAuth, requireAdminOrTutor, lessonHandler.DeleteLessonHandler)
	}
}
