package repository

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
)

type ProductRepository interface {
	List() ([]entity.Product, error)
	Create(product dto.CreateProduct) (entity.Product, error)
	Update(product dto.UpdateProduct) (entity.Product, error)
	Get(id int) (entity.Product, error)
	Delete(id string) error
}
