package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redukasquad/be-reduka/database/migrations"
	"github.com/redukasquad/be-reduka/modules/auth"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	migrations.ConnectDatabase()

	r := gin.Default()

	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		auth.AuthRouter(v1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
