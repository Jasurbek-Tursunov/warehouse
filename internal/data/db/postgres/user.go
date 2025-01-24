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

	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	err := u.store.conn.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&out.ID)
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
	query := `SELECT * FROM users WHERE username = $1`

	err := u.store.conn.QueryRowContext(ctx, query, username).Scan(&out.ID)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
