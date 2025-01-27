package usecase

import (
	"errors"
	"fmt"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strconv"
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
		return "", fmt.Errorf("user with this username not found: %w", err)
	}

	hashedPassword, err := hashingPassword(password)
	if err != nil {
		return "", err
	}

	if assertPassword(user.Password, hashedPassword) {
		return "", errors.New("invalid password")
	}

	return encode(user.ID, a.expireDuration, a.secretKey)
}

func (a *AuthServiceImpl) Check(token string) error {
	id, err := decode(token, a.secretKey)
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

func hashingPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func assertPassword(password, hashPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)) == nil
}

func encode(id int, expireDuration time.Duration, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ID:        strconv.Itoa(id),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	return token.SignedString([]byte(secretKey))
}

func decode(accessToken, secretKey string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return strconv.Atoi(claims.ID)
	}

	return 0, errors.New("invalid token")
}
