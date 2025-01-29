package usecase

import (
	"errors"
	"fmt"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
	"github.com/Jasurbek-Tursunov/warehouse/pkg/jwt"
	pass "github.com/Jasurbek-Tursunov/warehouse/pkg/password"
	"log/slog"
	"time"
)

type AuthServiceImpl struct {
	logger         *slog.Logger
	secretKey      string
	expireDuration time.Duration

	userRepos repository.UserRepository
}

func NewAuthService(logger *slog.Logger, userRepos repository.UserRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		logger:         logger,
		secretKey:      "secret",
		expireDuration: 24 * time.Hour,
		userRepos:      userRepos,
	}
}

func (a *AuthServiceImpl) Register(data *dto.CreateUser) (*entity.User, error) {
	// Validation
	var validationErrors []*entity.ValidationError
	if data.Username == "" {
		validationErrors = append(validationErrors, entity.NewValidationError(
			"username",
			"username cannot be empty",
		))
	}

	if data.Password == "" {
		validationErrors = append(validationErrors, entity.NewValidationError(
			"password",
			"password cannot be empty",
		))
	}

	if len(validationErrors) > 0 {
		return nil, entity.WrapValidationError(validationErrors...)
	}

	hashedPassword, err := pass.HashingPassword(data.Password)
	if err != nil {

		return nil, err
	}
	data.Password = hashedPassword

	user, err := a.userRepos.Create(data)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthServiceImpl) Login(data *dto.Auth) (string, error) {
	// Validation
	var validationErrors []*entity.ValidationError
	if data.Username == "" {
		validationErrors = append(validationErrors, entity.NewValidationError(
			"username",
			"username cannot be empty",
		))
	}

	if data.Password == "" {
		validationErrors = append(validationErrors, entity.NewValidationError(
			"password",
			"password cannot be empty",
		))
	}

	if len(validationErrors) > 0 {
		return "", entity.WrapValidationError(validationErrors...)
	}

	user, err := a.userRepos.GetByUsername(data.Username)
	if err != nil {
		return "", fmt.Errorf("user with this username not found: %w", err)
	}

	if !pass.AssertPassword(data.Password, user.Password) {
		return "", errors.New("invalid password")
	}

	return jwt.Encode(user.ID, a.expireDuration, a.secretKey)
}

func (a *AuthServiceImpl) Check(token string) error {
	id, err := jwt.Decode(token, a.secretKey)
	if err != nil {
		a.logger.Debug("failed decode token", "error", err.Error())
		return err
	}

	_, err = a.userRepos.Get(id)
	if err != nil {
		a.logger.Debug("failed to get user", "error", err.Error())
		return err
	}
	return nil
}
