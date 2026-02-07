package attempts

import (
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
)

func AttemptRouter(router *gin.RouterGroup, requireAuth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	repo := NewRepository(migrations.GetDB())
	service := NewService(repo)
	handler := NewHandler(service)

	// Start attempt from registration
	router.POST("/tryouts/registrations/:regId/start", requireAuth, handler.StartAttemptHandler)

	// Attempt operations
	attemptRoutes := router.Group("/tryouts/attempts")
	attemptRoutes.Use(requireAuth)
	{
		attemptRoutes.GET("/:attemptId/current", handler.GetCurrentStateHandler)
		attemptRoutes.GET("/:attemptId/subtests/:subtestId/start", handler.StartSubtestHandler)
		attemptRoutes.POST("/:attemptId/subtests/:subtestId/submit", handler.SubmitSubtestHandler)
		attemptRoutes.POST("/:attemptId/finish", handler.FinishAttemptHandler)
		attemptRoutes.GET("/:attemptId/results", handler.GetResultsHandler)
	}

	// Public leaderboard
	router.GET("/tryouts/:id/leaderboard", handler.GetLeaderboardHandler)
}
