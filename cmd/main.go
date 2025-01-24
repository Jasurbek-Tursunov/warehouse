package main

import gin "github.com/Jasurbek-Tursunov/warehouse/internal/presenter/gin"

func main() {
	// TODO Load Config

	// TODO Setup Log

	// TODO Run server
	server := gin.NewServer()
	server.MustRun()

	// TODO Stop server
}
