package repository

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
)

type UserRepository interface {
	Create(user *dto.CreateUser) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
}
