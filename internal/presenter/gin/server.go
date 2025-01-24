package pgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Server struct {
	Port   int
	router *gin.Engine
	logger *slog.Logger
}

func NewServer() *Server {
	cfg := MustLoadEnvConfig()

	return &Server{
		Port:   cfg.Port,
		router: NewRouter(),
	}
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) Run() error {
	port := fmt.Sprintf(":%d", s.Port)

	if err := s.router.Run(port); err != nil {
		return err
	}
	return nil
}
