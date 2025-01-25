package di

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/data"
	"github.com/Jasurbek-Tursunov/warehouse/internal/data/db/postgres"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository"
	dusecase "github.com/Jasurbek-Tursunov/warehouse/internal/domain/usecase"
	"github.com/Jasurbek-Tursunov/warehouse/internal/presenter"
	"github.com/Jasurbek-Tursunov/warehouse/internal/usecase"
)

type StorageType int

type Container struct {
	store data.Storage

	userRepo    repository.UserRepository
	productRepo repository.ProductRepository

	authService    dusecase.AuthService
	productService dusecase.ProductService

	Server presenter.Server
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) InitStore() {
	c.store = postgres.NewStorage()
	c.store.MustConnect()
}

func (c *Container) InitUserRepo() {
	if s, ok := c.store.(*postgres.Storage); ok {
		c.userRepo = postgres.NewUserRepository(s)
	}
}

func (c *Container) InitProductRepo() {
	if s, ok := c.store.(*postgres.Storage); ok {
		c.productRepo = postgres.NewProductRepository(s)
	}
}

func (c *Container) InitAuthService() {
	c.authService = usecase.NewAuthService(c.userRepo)
}

func (c *Container) InitProductService() {
	c.productService = usecase.NewProductService(c.productRepo)
}

func (c *Container) InitGinServer() {}

func (c *Container) Close() {
	if c.store != nil {
		c.store.Close()
	}
}
