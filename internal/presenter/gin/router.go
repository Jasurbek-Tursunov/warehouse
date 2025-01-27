package gin

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/usecase"
	"github.com/Jasurbek-Tursunov/warehouse/internal/presenter/gin/hendler"
	libgin "github.com/gin-gonic/gin"
)

func (s *Server) Register(authService usecase.AuthService, productService usecase.ProductService) {
	middleware := NewMiddleware(authService)

	auth := hendler.NewAuthHandler(authService)
	product := hendler.NewProductHandler(productService)

	s.router = libgin.Default()
	authorized := s.router.Group("/").Use(middleware.Auth)

	s.router.GET("/ping", func(c *libgin.Context) {
		c.JSON(200, libgin.H{"client": c.ClientIP(), "server": c.Request.Host})
	})

	s.router.POST("/register", auth.Register)
	s.router.POST("/login", auth.Login)

	authorized.GET("/products", product.List)
	authorized.POST("/products", product.Create)
	authorized.GET("/products/:id", product.Get)
	authorized.PUT("/products/:id", product.Update)
	authorized.DELETE("/products/:id", product.Delete)
}
