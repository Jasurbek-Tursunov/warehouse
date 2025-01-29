package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Jasurbek-Tursunov/warehouse/internal/data/db/postgres"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
	"time"
)

type ProductRepositoryImpl struct {
	store *postgres.Storage
}

func NewProductRepository(store *postgres.Storage) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{store: store}
}

func (p *ProductRepositoryImpl) List(filters *dto.ProductQuery, paginate *dto.Paginator) ([]entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.store.Timeout)
	defer cancel()

	query := `SELECT id, name, price, quantity, created_at, updated_at FROM products`
	args := make([]any, 0)

	if filters.Name != "" {
		query += ` WHERE name ILIKE $1`
		args = append(args, "%"+filters.Name+"%")
	}

	switch filters.SortBy {
	case "last_create":
		query += ` ORDER BY created_at DESC`
	case "name":
		query += ` ORDER BY name ASC`
	default:
	}

	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, paginate.PageSize, max(paginate.Page-1, 0)*paginate.PageSize)

	rows, err := p.store.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return []entity.Product{}, err
	}

	products := make([]entity.Product, 0, paginate.PageSize)
	var product entity.Product
	for rows.Next() {
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			continue
		}

		products = append(products, product)
	}
	return products, nil
}

func (p *ProductRepositoryImpl) Create(product *dto.CreateProduct) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.store.Timeout)
	defer cancel()

	out := entity.Product{
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}

	query := `INSERT INTO products(name, price, quantity) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err := p.store.DB.QueryRowContext(ctx, query, product.Name, product.Price, product.Quantity).
		Scan(&out.ID, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (p *ProductRepositoryImpl) Update(id int, product *dto.UpdateProduct) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.store.Timeout)
	defer cancel()

	out := entity.Product{
		ID:       id,
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}

	query := `UPDATE products SET name = $1, price = $2, quantity = $3, updated_at = $4 
                WHERE id = $5 RETURNING created_at, updated_at`
	row := p.store.DB.QueryRowContext(ctx, query, product.Name, product.Price, product.Quantity, time.Now(), id)

	if err := row.Scan(&out.CreatedAt, &out.UpdatedAt); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, entity.NewNotFoundError("product", id)
		default:
			return nil, err
		}
	}

	return &out, nil
}

func (p *ProductRepositoryImpl) Get(id int) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.store.Timeout)
	defer cancel()

	query := `SELECT id, name, price, quantity, created_at, updated_at FROM products WHERE id = $1`
	row := p.store.DB.QueryRowContext(ctx, query, id)

	var out entity.Product
	err := row.Scan(
		&out.ID,
		&out.Name,
		&out.Price,
		&out.Quantity,
		&out.CreatedAt,
		&out.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, entity.NewNotFoundError("product", id)
		default:
			return nil, err
		}
	}

	return &out, nil
}

func (p *ProductRepositoryImpl) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.store.Timeout)
	defer cancel()

	query := `DELETE FROM products WHERE id = $1`
	result, err := p.store.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsCount == 0 {
		return entity.NewNotFoundError("product", id)
	}

	return nil
}
