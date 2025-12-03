package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redukasquad/be-reduka/configs"
)

type Server struct {
	cfg    *configs.Config
	engine *gin.Engine
}

func NewServer(cfg *configs.Config) *Server {
	return &Server{
		cfg:    cfg,
		engine: gin.Default(),
	}
}

func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

func (s *Server) Run() {
	if err := s.engine.Run(":" + s.cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
