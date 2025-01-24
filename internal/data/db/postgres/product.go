package postgres

import (
	"context"
	"errors"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
)

type ProductRepositoryImpl struct {
	store *Storage
}

func NewProductRepository(store *Storage) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{store: store}
}

func (p *ProductRepositoryImpl) Create(product *dto.CreateProduct) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.store.timeout)
	defer cancel()

	out := entity.Product{
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}

	query := `INSERT INTO products(name, price, quantity) VALUES ($1, $2, $3)`

	err := p.store.conn.QueryRowContext(ctx, query, product.Name, product.Price, product.Quantity).
		Scan(&out.ID, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (p *ProductRepositoryImpl) Update(id int, product *dto.UpdateProduct) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.store.timeout)
	defer cancel()

	out := entity.Product{
		ID:       id,
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}

	query := `UPDATE products SET name = $1, price = $2, quantity = $3 WHERE id = $4`
	result, err := p.store.conn.ExecContext(ctx, query, product.Name, product.Price, product.Quantity, id)
	if err != nil {
		return nil, err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowCount == 0 {
		return nil, errors.New("product not found") // TODO create template errors
	}

	return &out, nil
}
