package main

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/di"
	gin "github.com/Jasurbek-Tursunov/warehouse/internal/presenter/gin"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO Load Config

	// TODO Setup Log

	server := gin.NewServer()
	server.MustRun()

	// TODO Run server
	container := di.NewContainer()
	defer container.Close()

	go func() {
		container.Server.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	// TODO Stop server
	container.Server.GracefulStop()
}
