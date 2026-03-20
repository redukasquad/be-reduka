package subjects

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func SubjectRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdminOrTutor gin.HandlerFunc) {
	subjectRepo := NewRepository(migrations.GetDB())
	subjectService := NewService(subjectRepo)
	subjectHandler := NewHandler(subjectService)

	// NOTE: /courses/:id/classes is registered in CourseIndexRouter to avoid Gin wildcard conflict

	classes := router.Group("/classes")
	{
		classes.GET("/:id", subjectHandler.GetSubjectByIDHandler)
		classes.PUT("/:id", requireAuth, requireAdminOrTutor, subjectHandler.UpdateSubjectHandler)
		classes.DELETE("/:id", requireAuth, requireAdminOrTutor, subjectHandler.DeleteSubjectHandler)
	}
}

// RegisterCourseSubRoutes registers /courses/:id/classes under an existing courseByID group
// to avoid Gin wildcard conflicts.
func RegisterCourseSubRoutes(courseByID *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdminOrTutor gin.HandlerFunc, h Handler) {
	courseByID.GET("/classes", h.GetSubjectsByCourseHandler)
	courseByID.POST("/classes", requireAuth, requireAdminOrTutor, h.CreateSubjectHandler)
}
