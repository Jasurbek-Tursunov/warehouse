package presenter

import "github.com/Jasurbek-Tursunov/warehouse/internal/domain/usecase"

type Server interface {
	Register(usecase.AuthService, usecase.ProductService)
	MustRun()
	GracefulStop()
}
