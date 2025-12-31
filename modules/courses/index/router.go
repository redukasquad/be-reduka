package courses

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func CourseIndexRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	courseRepo := NewRepository(migrations.GetDB())
	courseService := NewService(courseRepo)
	courseHandler := NewHandler(courseService)

	courses := router.Group("/courses")
	{
		courses.GET("", courseHandler.GetAllCoursesHandler)
		courses.POST("", requireAuth, requireAdmin, courseHandler.CreateCourseHandler)
	}

	courseByID := router.Group("/courses")
	{
		courseByID.GET("/:id", courseHandler.GetCourseByIDHandler)
		courseByID.PUT("/:id", requireAuth, requireAdmin, courseHandler.UpdateCourseHandler)
		courseByID.DELETE("/:id", requireAuth, requireAdmin, courseHandler.DeleteCourseHandler)
	}

	router.GET("/programs/:id/courses", courseHandler.GetCoursesByProgramIDHandler)
}
