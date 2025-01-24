package usecase

import (
	"errors"
	"fmt"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type AuthServiceImpl struct {
	userRepos      repository.UserRepository
	secretKey      string
	expireDuration time.Duration
}

func NewAuthService(userRepos repository.UserRepository) *AuthServiceImpl {
	return &AuthServiceImpl{userRepos: userRepos}
}

func (a *AuthServiceImpl) Register(username, password string) (*entity.User, error) {
	hashedPassword, err := hashingPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := a.userRepos.Create(
		&dto.CreateUser{
			Username: username,
			Password: hashedPassword,
		},
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthServiceImpl) Login(username, password string) (string, error) {
	user, err := a.userRepos.GetByUsername(username)
	if err != nil {
		return "", errors.New("user with this username not found")
	}

	hashedPassword, err := hashingPassword(password)
	if err != nil {
		return "", err
	}

	if user.Password != hashedPassword {
		return "", errors.New("invalid password")
	}

	return a.encode(int(user.ID))
}

func hashingPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (a *AuthServiceImpl) encode(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        strconv.Itoa(id),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.expireDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	return token.SignedString([]byte(a.secretKey))
}

func (a *AuthServiceImpl) decode(accessToken string) (int, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.RegisteredClaims); ok && token.Valid {
		return strconv.Atoi(claims.ID)
	}

	return 0, errors.New("invalid token")
}
