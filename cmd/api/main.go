package api

import (
	"log"
	"os"
	"strings"
	"time"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/database/migrations"
	"github.com/redukasquad/be-reduka/middleware"
	"github.com/redukasquad/be-reduka/modules/auth"
	"github.com/redukasquad/be-reduka/modules/classes"
	"github.com/redukasquad/be-reduka/modules/courses"
	"github.com/redukasquad/be-reduka/modules/health"
	"github.com/redukasquad/be-reduka/modules/programs"
	"github.com/redukasquad/be-reduka/modules/tryouts"
	"github.com/redukasquad/be-reduka/modules/universities"
	"github.com/redukasquad/be-reduka/modules/uploads"
	"github.com/redukasquad/be-reduka/modules/users"
	"github.com/redukasquad/be-reduka/packages/utils"
)

var router *gin.Engine

func getAllowedOrigins() []string {
	origins := []string{"http://localhost:3000", "http://localhost:5173"}

	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		for _, url := range strings.Split(frontendURL, ",") {
			url = strings.TrimSpace(url)
			if url != "" {
				origins = append(origins, url)
			}
		}
	}

	return origins
}

func init() {
	utils.InitLogger()
	migrations.ConnectDatabase()

	router = gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     getAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	v1 := api.Group("/v1")
	{
		health.HealthRouter(v1)

		auth.AuthRouter(v1)
		users.UserRouter(v1, middleware.RequireAuth(), middleware.RequireAdmin())
		programs.ProgramRouter(v1, middleware.RequireAuth(), middleware.RequireAdmin())
		courses.CoursesRouter(v1, middleware.RequireAuth(), middleware.RequireAdmin(), middleware.RequireAdminOrTutor())
		classes.ClassesRouter(v1, middleware.RequireAuth(), middleware.RequireAdminOrTutor())
		tryouts.TryOutsRouter(v1, middleware.RequireAuth(), middleware.RequireAdminOrTutor())
		universities.UniversityRouter(v1, middleware.RequireAuth(), middleware.RequireAdmin())
		uploads.UploadRouter(v1, middleware.RequireAuth(), middleware.RequireAdminOrTutorOrUser())
	}

	log.Println("Router initialized successfully")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}