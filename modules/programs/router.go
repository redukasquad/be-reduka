package programs

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func ProgramRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	programRepo := NewRepository(migrations.GetDB())
	programService := NewService(programRepo)
	programHandler := NewHandler(programService)

	programs := router.Group("/programs")
	{
		programs.GET("", programHandler.GetAllProgramsHandler)
		programs.GET("/:id", programHandler.GetProgramByIDHandler)

		programs.POST("", requireAuth, requireAdmin, programHandler.CreateProgramHandler)
		programs.PUT("/:id", requireAuth, requireAdmin, programHandler.UpdateProgramHandler)
		programs.DELETE("/:id", requireAuth, requireAdmin, programHandler.DeleteProgramHandler)
	}
}