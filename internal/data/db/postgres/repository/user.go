package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Jasurbek-Tursunov/warehouse/internal/data/db/postgres"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
	"github.com/lib/pq"
	"log/slog"
)

type UserRepositoryImpl struct {
	logger *slog.Logger
	store  *postgres.Storage
}

func NewUserRepository(logger *slog.Logger, store *postgres.Storage) *UserRepositoryImpl {
	return &UserRepositoryImpl{logger, store}
}

func (u *UserRepositoryImpl) Create(user *dto.CreateUser) (*entity.User, error) {
	const op = "postgres.repository.user.Create"

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
			err = entity.WrapValidationError(entity.NewValidationError(
				"username",
				"User with username: "+user.Username+" already exist",
			))
			u.logger.Debug(op+"failed validation", "error", err.Error())
		default:
			u.logger.Warn(op+"failed executing", "error", err.Error())
		}
		return nil, err
	}

	return &out, nil
}

func (u *UserRepositoryImpl) Get(id int) (*entity.User, error) {
	const op = "postgres.repository.user.Get"

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
			err = entity.NewNotFoundError("user", id)
			u.logger.Debug(op+" failed not found", "error", err.Error())
			return nil, err
		}
		u.logger.Warn(op+" failed executing", "error", err.Error())
		return nil, err
	}

	return &out, nil
}

func (u *UserRepositoryImpl) GetByUsername(username string) (*entity.User, error) {
	const op = "postgres.repository.user.GetByUsername"

	ctx, cancel := context.WithTimeout(context.Background(), u.store.Timeout)
	defer cancel()

	out := entity.User{
		Username: username,
	}

	query := `SELECT id, password FROM users WHERE username = $1`

	err := u.store.DB.QueryRowContext(ctx, query, username).Scan(&out.ID, &out.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = entity.NewNotFoundError("user", username)
			u.logger.Debug(op+" failed not found", "error", err.Error())
			return nil, err
		}
		u.logger.Warn(op+" failed executing", "error", err.Error())
		return nil, err
	}

	return &out, nil
}
