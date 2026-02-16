package uploads

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func UploadRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	uploadRepo := NewRepository(migrations.GetDB())
	uploadService := NewService(uploadRepo)
	uploadHandler := NewHandler(uploadService)

	uploads := router.Group("/uploads")
	{
		uploads.POST("", requireAuth, requireAdmin, uploadHandler.CreateImageHandler)
		uploads.DELETE("", requireAuth, requireAdmin, uploadHandler.DeleteImageHandler)
	}
}
