package repository

import "github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
}
