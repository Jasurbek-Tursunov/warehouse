package postgres

import (
	"context"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
)

type UserRepositoryImpl struct {
	store *Storage
}

func NewUserRepository(store *Storage) *UserRepositoryImpl {
	return &UserRepositoryImpl{store: store}
}

func (u *UserRepositoryImpl) Create(user *dto.CreateUser) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.store.timeout)
	defer cancel()

	out := entity.User{
		Username: user.Username,
	}

	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	err := u.store.conn.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&out.ID)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (u *UserRepositoryImpl) Get(id int) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.store.timeout)
	defer cancel()

	out := entity.User{ID: id}

	query := `SELECT id, username FROM users WHERE id = $1`
	row := u.store.conn.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&out.ID,
		&out.Username,
	)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (u *UserRepositoryImpl) GetByUsername(username string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.store.timeout)
	defer cancel()

	out := entity.User{
		Username: username,
	}

	query := `SELECT id FROM users WHERE username = $1`

	err := u.store.conn.QueryRowContext(ctx, query, username).Scan(&out.ID)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
