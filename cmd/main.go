package main

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/di"
	"log/slog"
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

	// TODO Setup Log

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

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
