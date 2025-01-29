package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Jasurbek-Tursunov/warehouse/internal/data/db/postgres"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
	"github.com/lib/pq"
)

type UserRepositoryImpl struct {
	store *postgres.Storage
}

func NewUserRepository(store *postgres.Storage) *UserRepositoryImpl {
	return &UserRepositoryImpl{store: store}
}

func (u *UserRepositoryImpl) Create(user *dto.CreateUser) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.store.Timeout)
	defer cancel()

	out := entity.User{
		Username: user.Username,
	}

	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	err := u.store.DB.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&out.ID)
	if err != nil {
		var errPQ *pq.Error
		switch {
		case errors.As(err, &errPQ) && errPQ.Code == "23505":
			return nil, entity.NewValidationError(
				"username",
				"User with username: "+user.Username+" already exist",
			)
		default:
			return nil, err
		}
	}

	return &out, nil
}

func (u *UserRepositoryImpl) Get(id int) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.store.Timeout)
	defer cancel()

	out := entity.User{ID: id}

	query := `SELECT id, username FROM users WHERE id = $1`
	row := u.store.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&out.ID,
		&out.Username,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewNotFoundError("user", id)
		}
		return nil, err
	}

	return &out, nil
}

func (u *UserRepositoryImpl) GetByUsername(username string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.store.Timeout)
	defer cancel()

	out := entity.User{
		Username: username,
	}

	query := `SELECT id, password FROM users WHERE username = $1`

	err := u.store.DB.QueryRowContext(ctx, query, username).Scan(&out.ID, &out.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewNotFoundError("user", username)
		}
		return nil, err
	}

	return &out, nil
}
