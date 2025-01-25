package usecase

import "github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"

type AuthService interface {
	Register(username, password string) (*entity.User, error)
	Login(username, password string) (string, error)
}
