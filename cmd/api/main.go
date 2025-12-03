package main

import (
	"log"
	"time"

	"github.com/redukasquad/be-reduka/configs"
	_http "github.com/redukasquad/be-reduka/internal/delivery/http"
	"github.com/redukasquad/be-reduka/internal/repository"
	"github.com/redukasquad/be-reduka/internal/server"
	"github.com/redukasquad/be-reduka/internal/usecase"
	"github.com/redukasquad/be-reduka/pkg/database"
)

func main() {
	// Load Config
	cfg := configs.LoadConfig()

	// Init Database
	db := database.ConnectDB(cfg)

	// Auto Migrate
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Init Server
	s := server.NewServer(cfg)

	// Init Layers
	timeoutContext := 2 * time.Second

	// Repositories
	userRepo := repository.NewUserRepository(db)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, timeoutContext, cfg)

	// Init Handlers
	_http.NewAuthHandler(s.GetEngine(), authUsecase)

	// Run Server
	s.Run()
}
