package usecase

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
)

type AuthService interface {
	Register(*dto.CreateUser) (*entity.User, error)
	Login(*dto.Auth) (string, error)
	Check(token string) error
}
