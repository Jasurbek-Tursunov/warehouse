package main

import (
	conf "github.com/Jasurbek-Tursunov/warehouse/internal/config"
	"github.com/Jasurbek-Tursunov/warehouse/internal/di"
	"github.com/Jasurbek-Tursunov/warehouse/pkg/config"
	"github.com/Jasurbek-Tursunov/warehouse/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

// @title Warehouse API
// @version 1.0

// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

func main() {
	// TODO Load Config
	cfg := config.MustLoad[conf.Config]()

	// TODO Setup Log
	logger := log.SetupLogger(cfg.Env)

	// TODO Run server
	container := di.NewContainer(logger)
	defer container.Close()

	container.InitStore()
	container.InitUserRepo()
	container.InitAuthService()
	container.InitProductRepo()
	container.InitProductService()
	container.InitServer()

	go func() {
		container.Server.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	// TODO Stop server
	container.Server.GracefulStop()
}
