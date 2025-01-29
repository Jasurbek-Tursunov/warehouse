package usecase

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
	"log/slog"
)

type ProductServiceImpl struct {
	logger      *slog.Logger
	productRepo repository.ProductRepository
}

func NewProductService(logger *slog.Logger, productRepo repository.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{logger: logger, productRepo: productRepo}
}

func (p *ProductServiceImpl) List(filters *dto.ProductQuery, paginate *dto.Paginator) ([]entity.Product, error) {
	return p.productRepo.List(filters, paginate)
}

func (p *ProductServiceImpl) Create(product *dto.CreateProduct) (*entity.Product, error) {
	// Validation
	var validationErrors []*entity.ValidationError
	if product.Name == "" {
		validationErrors = append(validationErrors, entity.NewValidationError(
			"name",
			"name cannot be empty",
		))
	}
	if product.Price <= 0 {
		validationErrors = append(validationErrors, entity.NewValidationError(
			"price",
			"price must be greater than 0",
		))
	}
	if product.Quantity < 0 {
		validationErrors = append(validationErrors, entity.NewValidationError(
			"quantity",
			"quantity cannot be negative",
		))
	}
	if len(validationErrors) > 0 {
		return nil, entity.WrapValidationError(validationErrors...)
	}
	return p.productRepo.Create(product)
}

func (p *ProductServiceImpl) Update(id int, product *dto.UpdateProduct) (*entity.Product, error) {
	// Validation
	var validationErrors []*entity.ValidationError
	if product.Name == "" {
		validationErrors = append(validationErrors, entity.NewValidationError(
			"name",
			"name cannot be empty",
		))
	}
	if product.Price <= 0 {
		validationErrors = append(validationErrors, entity.NewValidationError(
			"price",
			"price must be greater than 0",
		))
	}
	if product.Quantity < 0 {
		validationErrors = append(validationErrors, entity.NewValidationError(
			"quantity",
			"quantity cannot be negative",
		))
	}
	if len(validationErrors) > 0 {
		return nil, entity.WrapValidationError(validationErrors...)
	}

	return p.productRepo.Update(id, product)
}

func (p *ProductServiceImpl) Get(id int) (*entity.Product, error) {
	return p.productRepo.Get(id)
}

func (p *ProductServiceImpl) Delete(id int) error {
	return p.productRepo.Delete(id)
}
