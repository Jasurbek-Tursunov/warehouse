package postgres

import (
	"context"
	"database/sql"
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

func (p *ProductRepositoryImpl) List(filters *dto.ProductQuery, paginate *dto.Paginator) ([]entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.store.timeout)
	defer cancel()

	args := []any{
		paginate.PageSize,
		max(paginate.Page-1, 0) * paginate.PageSize,
	}

	query := `SELECT id, name, price, quantity, created_at, updated_at FROM products LIMIT $1 OFFSET $2`
	rows, err := p.store.conn.QueryContext(ctx, query, args...)
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

func (p *ProductRepositoryImpl) Get(id int) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.store.timeout)
	defer cancel()

	query := `SELECT id, name, price, quantity, created_at, updated_at FROM products WHERE id = $1`
	row := p.store.conn.QueryRowContext(ctx, query, id)

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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NotFoundError
		}
		return nil, err
	}

	return &out, nil
}

func (p *ProductRepositoryImpl) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.store.timeout)
	defer cancel()

	query := `DELETE FROM products WHERE id = $1`
	result, err := p.store.conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsCount == 0 {
		return entity.NotFoundError
	}

	return nil
}
