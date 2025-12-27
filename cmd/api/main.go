package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redukasquad/be-reduka/database/migrations"
	"github.com/redukasquad/be-reduka/middleware"
	"github.com/redukasquad/be-reduka/modules/auth"
	"github.com/redukasquad/be-reduka/modules/health"
	"github.com/redukasquad/be-reduka/modules/programs"
	"github.com/redukasquad/be-reduka/modules/users"
	"github.com/redukasquad/be-reduka/packages/utils"
)

// getAllowedOrigins returns CORS allowed origins from environment
func getAllowedOrigins() []string {
	// Default development origins
	origins := []string{"http://localhost:3000", "http://localhost:5173"}

	// Add production frontend URL if set
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		// Support multiple URLs separated by comma
		for _, url := range strings.Split(frontendURL, ",") {
			url = strings.TrimSpace(url)
			if url != "" && url != "http://localhost:3000" && url != "http://localhost:5173" {
				origins = append(origins, url)
			}
		}
	}

	return origins
}

func main() {
	// Load .env file if exists (for local development)
	// In production, environment variables are set directly
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}

	migrations.ConnectDatabase()

	utils.InitLogger()
	r := gin.Default()

	// CORS configuration with dynamic origins
	r.Use(cors.New(cors.Config{
		AllowOrigins:     getAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		// Health check endpoint (no auth required)
		health.HealthRouter(v1)

		auth.AuthRouter(v1)
		users.UserRouter(v1, middleware.RequireAuth(), middleware.RequireAdmin())
		programs.ProgramRouter(v1, middleware.RequireAuth(), middleware.RequireAdmin())
	}

	port := os.Getenv("GOLANG_PORT")
	if port == "" {
		port = "8888"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
