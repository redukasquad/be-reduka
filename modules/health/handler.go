package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthRouter registers health check endpoint
func HealthRouter(r *gin.RouterGroup) {
	r.GET("/health", HealthCheck)
}

// HealthCheck returns the health status of the API
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "API is running",
	})
}
