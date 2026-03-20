package courses

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
	"github.com/redukasquad/be-reduka/modules/classes/subjects"
)

func CourseIndexRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc, requireAdminOrTutor gin.HandlerFunc) {
	courseRepo := NewRepository(migrations.GetDB())
	courseService := NewService(courseRepo)
	courseHandler := NewHandler(courseService)

	courses := router.Group("/courses")
	{
		courses.GET("", courseHandler.GetAllCoursesHandler)
		courses.POST("", requireAuth, requireAdminOrTutor, courseHandler.CreateCourseHandler)
	}

	// Single group for /courses/:id to avoid Gin wildcard conflicts
	courseByID := router.Group("/courses/:id")
	{
		courseByID.GET("", courseHandler.GetCourseByIDHandler)
		courseByID.PUT("", requireAuth, requireAdmin, courseHandler.UpdateCourseHandler)
		courseByID.DELETE("", requireAuth, requireAdmin, courseHandler.DeleteCourseHandler)

		// Register /courses/:id/classes here to share the same wildcard group
		subjectRepo := subjects.NewRepository(migrations.GetDB())
		subjectSvc := subjects.NewService(subjectRepo)
		subjectHandler := subjects.NewHandler(subjectSvc)
		subjects.RegisterCourseSubRoutes(courseByID, requireAuth, requireAdminOrTutor, subjectHandler)
	}

	// Tutor: lihat courses milik sendiri
	router.GET("/tutor/my-courses", requireAuth, requireAdminOrTutor, courseHandler.GetMyCoursesTutorHandler)

	router.GET("/programs/:id/courses", courseHandler.GetCoursesByProgramIDHandler)
}
