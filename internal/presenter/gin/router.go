package gin

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/presenter/gin/hendler"
	libgin "github.com/gin-gonic/gin"
)

func NewRouter() *libgin.Engine {
	auth := hendler.AuthHandler{}
	product := hendler.ProductHandler{}

	router := libgin.Default()

	router.GET("/ping", func(c *libgin.Context) {
		c.JSON(200, libgin.H{"client": c.ClientIP(), "server": c.Request.Host})
	})

	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)

	router.GET("/products", product.List)
	router.POST("/products", product.Create)
	router.GET("/products/{}", product.Get)
	router.PUT("/products/{}", product.Update)
	router.DELETE("/products/{}", product.Delete)

	return router
}
