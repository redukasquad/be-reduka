package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func ResourceRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdminOrTutor gin.HandlerFunc) {
	resourceRepo := NewRepository(migrations.GetDB())
	resourceService := NewService(resourceRepo)
	resourceHandler := NewHandler(resourceService)

	lessonResources := router.Group("/lessons/:id")
	{
		lessonResources.GET("/resources", resourceHandler.GetResourcesByLessonHandler)
		lessonResources.POST("/resources", requireAuth, requireAdminOrTutor, resourceHandler.CreateResourceHandler)
	}

	resources := router.Group("/resources")
	resources.Use(requireAuth, requireAdminOrTutor)
	{
		resources.PUT("/:id", resourceHandler.UpdateResourceHandler)
		resources.DELETE("/:id", resourceHandler.DeleteResourceHandler)
	}
}
