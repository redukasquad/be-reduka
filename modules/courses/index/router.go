package courses

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

// CourseIndexRouter registers all course index routes
func CourseIndexRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	courseRepo := NewRepository(migrations.GetDB())
	courseService := NewService(courseRepo)
	courseHandler := NewHandler(courseService)

	courses := router.Group("/courses")
	{
		// Public routes
		courses.GET("", courseHandler.GetAllCoursesHandler)
		courses.GET("/:id", courseHandler.GetCourseByIDHandler)
		courses.GET("/program/:programId", courseHandler.GetCoursesByProgramIDHandler)

		// Admin only routes
		courses.POST("", requireAuth, requireAdmin, courseHandler.CreateCourseHandler)
		courses.PUT("/:id", requireAuth, requireAdmin, courseHandler.UpdateCourseHandler)
		courses.DELETE("/:id", requireAuth, requireAdmin, courseHandler.DeleteCourseHandler)
	}
}
