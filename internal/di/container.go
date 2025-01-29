package di

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/data"
	"github.com/Jasurbek-Tursunov/warehouse/internal/data/db/postgres"
	repoPostgres "github.com/Jasurbek-Tursunov/warehouse/internal/data/db/postgres/repository"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository"
	dusecase "github.com/Jasurbek-Tursunov/warehouse/internal/domain/usecase"
	"github.com/Jasurbek-Tursunov/warehouse/internal/presenter"
	"github.com/Jasurbek-Tursunov/warehouse/internal/presenter/gin"
	"github.com/Jasurbek-Tursunov/warehouse/internal/usecase"
	"log/slog"
)

type Container struct {
	logger *slog.Logger

	store data.Storage

	userRepo    repository.UserRepository
	productRepo repository.ProductRepository

	authService    dusecase.AuthService
	productService dusecase.ProductService

	Server presenter.Server
}

func NewContainer(logger *slog.Logger) *Container {
	return &Container{logger: logger}
}

func (c *Container) InitStore(storeType StorageType) {
	switch storeType {
	case StoragePostgres:
		c.store = postgres.NewStorage(c.logger)
		c.store.MustConnect()
	case StorageMock:
		// TODO
	case StorageMySQL:
		// TODO
	default:
		// nil
	}
}

func (c *Container) InitUserRepo() {
	if s, ok := c.store.(*postgres.Storage); ok {
		c.userRepo = repoPostgres.NewUserRepository(c.logger, s)
	}
}

func (c *Container) InitProductRepo() {
	if s, ok := c.store.(*postgres.Storage); ok {
		c.productRepo = repoPostgres.NewProductRepository(c.logger, s)
	}
}

func (c *Container) InitAuthService() {
	c.authService = usecase.NewAuthService(c.logger, c.userRepo)
}

func (c *Container) InitProductService() {
	c.productService = usecase.NewProductService(c.logger, c.productRepo)
}

func (c *Container) InitServer() {
	c.Server = gin.NewServer(c.logger)
	c.Server.Register(c.authService, c.productService)
}

func (c *Container) Close() {
	if c.store != nil {
		c.store.Close()
	}
}
