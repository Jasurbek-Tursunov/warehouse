package usecase

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
)

type ProductService interface {
	List(filters dto.ProductQuery, paginate dto.Paginator) ([]entity.Product, error)
	Create(product dto.CreateProduct) (entity.Product, error)
	Update(id int, product dto.UpdateProduct) (entity.Product, error)
	Get(id int) (entity.Product, error)
	Delete(id string) error
}
