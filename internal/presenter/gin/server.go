package gin

import (
	"fmt"
	"github.com/Jasurbek-Tursunov/warehouse/pkg/config"
	libgin "github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

type Server struct {
	Port    int
	Timeout time.Duration
	logger  *slog.Logger
	router  *libgin.Engine
}

func NewServer(logger *slog.Logger) *Server {
	cfg := config.MustLoad[Config]()

	return &Server{
		Port:    cfg.Port,
		Timeout: cfg.Timeout,
		logger:  logger,
		//router:  router,
	}
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) GracefulStop() {

}

func (s *Server) Run() error {
	port := fmt.Sprintf(":%d", s.Port)

	if err := s.router.Run(port); err != nil {
		return err
	}
	return nil
}
