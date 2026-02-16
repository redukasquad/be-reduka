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
	"github.com/redukasquad/be-reduka/modules/classes"
	"github.com/redukasquad/be-reduka/modules/courses"
	"github.com/redukasquad/be-reduka/modules/health"
	"github.com/redukasquad/be-reduka/modules/programs"
	"github.com/redukasquad/be-reduka/modules/tryouts"
	"github.com/redukasquad/be-reduka/modules/uploads"
	"github.com/redukasquad/be-reduka/modules/users"
	"github.com/redukasquad/be-reduka/packages/utils"
)

func getAllowedOrigins() []string {
	origins := []string{"http://localhost:3000", "http://localhost:5173"}

	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
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
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}

	migrations.ConnectDatabase()

	utils.InitLogger()
	r := gin.Default()

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
		health.HealthRouter(v1)

		auth.AuthRouter(v1)
		users.UserRouter(v1, middleware.RequireAuth(), middleware.RequireAdmin())
		programs.ProgramRouter(v1, middleware.RequireAuth(), middleware.RequireAdmin())
		courses.CoursesRouter(v1, middleware.RequireAuth(), middleware.RequireAdmin())
		classes.ClassesRouter(v1, middleware.RequireAuth(), middleware.RequireAdminOrTutor())
		tryouts.TryOutsRouter(v1, middleware.RequireAuth(), middleware.RequireAdmin())
		uploads.UploadRouter(v1, middleware.RequireAuth(), middleware.RequireAdmin())
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
