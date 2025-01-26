package main

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/di"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO Load Config

	// TODO Setup Log

	// TODO Run server
	container := di.NewContainer()
	defer container.Close()

	container.InitStore()
	container.InitUserRepo()
	container.InitAuthService()
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
