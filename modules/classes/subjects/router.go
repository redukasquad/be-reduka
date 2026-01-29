package subjects

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func SubjectRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdminOrTutor gin.HandlerFunc) {
	subjectRepo := NewRepository(migrations.GetDB())
	subjectService := NewService(subjectRepo)
	subjectHandler := NewHandler(subjectService)

	courseSubjects := router.Group("/courses/:id")
	{
		courseSubjects.GET("/subjects", subjectHandler.GetSubjectsByCourseHandler)
		courseSubjects.POST("/subjects", requireAuth, requireAdminOrTutor, subjectHandler.CreateSubjectHandler)
	}

	subjects := router.Group("/subjects")
	{
		subjects.GET("/:id", subjectHandler.GetSubjectByIDHandler)
		subjects.PUT("/:id", requireAuth, requireAdminOrTutor, subjectHandler.UpdateSubjectHandler)
		subjects.DELETE("/:id", requireAuth, requireAdminOrTutor, subjectHandler.DeleteSubjectHandler)
	}
}
