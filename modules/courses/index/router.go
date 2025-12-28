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

		// Admin only routes
		courses.POST("", requireAuth, requireAdmin, courseHandler.CreateCourseHandler)
	}

	// Routes with :id parameter - separate group to avoid conflicts
	courseByID := router.Group("/courses")
	{
		courseByID.GET("/:id", courseHandler.GetCourseByIDHandler)
		courseByID.PUT("/:id", requireAuth, requireAdmin, courseHandler.UpdateCourseHandler)
		courseByID.DELETE("/:id", requireAuth, requireAdmin, courseHandler.DeleteCourseHandler)
	}

	// Routes by program - use different path to avoid conflict
	router.GET("/programs/:programId/courses", courseHandler.GetCoursesByProgramIDHandler)
}
