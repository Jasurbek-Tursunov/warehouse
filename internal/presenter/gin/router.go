package gin

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/usecase"
	"github.com/Jasurbek-Tursunov/warehouse/internal/presenter/gin/hendler"
	swaggerfiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/Jasurbek-Tursunov/warehouse/docs"
)

func (s *Server) Register(authService usecase.AuthService, productService usecase.ProductService) {
	middleware := NewMiddleware(authService)

	auth := hendler.NewAuthHandler(authService)
	product := hendler.NewProductHandler(productService)

	s.router = gin.New()
	s.router.Use(middleware.SlogLogger(s.logger))
	authorized := s.router.Group("/").Use(middleware.Auth)

	s.router.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"pong": true}) })
	s.router.GET("/swagger/*any", swagger.WrapHandler(swaggerfiles.Handler))
	s.router.POST("/register", auth.Register)
	s.router.POST("/login", auth.Login)

	authorized.GET("/products", product.List)
	authorized.POST("/product/add", product.Create)
	authorized.GET("/product/:id", product.Get)
	authorized.PUT("/product/:id", product.Update)
	authorized.DELETE("/product/:id", product.Delete)
}
