package universities

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func UniversityRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	repo := NewRepository(migrations.GetDB())
	service := NewService(repo)
	h := NewHandler(service)

	// Public: list & detail universities
	unis := router.Group("/universities")
	{
		unis.GET("", h.GetAllUniversitiesHandler)
		unis.GET("/:id", h.GetUniversityByIDHandler)
		unis.GET("/:id/majors", h.GetMajorsByUniversityHandler)
		unis.GET("/:id/users", requireAuth, requireAdmin, h.GetUsersByUniversityHandler)

		// Admin only
		unis.POST("", requireAuth, requireAdmin, h.CreateUniversityHandler)
		unis.PUT("/:id", requireAuth, requireAdmin, h.UpdateUniversityHandler)
		unis.DELETE("/:id", requireAuth, requireAdmin, h.DeleteUniversityHandler)
		unis.POST("/:id/majors", requireAuth, requireAdmin, h.CreateMajorHandler)
	}

	// Admin: manage majors
	majors := router.Group("/majors")
	majors.Use(requireAuth, requireAdmin)
	{
		majors.PUT("/:id", h.UpdateMajorHandler)
		majors.DELETE("/:id", h.DeleteMajorHandler)
	}

	// User: manage own targets
	targets := router.Group("/users/me/targets")
	targets.Use(requireAuth)
	{
		targets.GET("", h.GetMyTargetsHandler)
		targets.POST("", h.AddTargetHandler)
		targets.DELETE("/:id", h.DeleteTargetHandler)
	}
}
