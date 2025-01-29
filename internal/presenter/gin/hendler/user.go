package hendler

import (
	"errors"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	service usecase.AuthService
}

func NewAuthHandler(service usecase.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// @Summary Register
// @Tags auth
// @Description registrantion
// @Accept json
// @Produce json
// @Param input body dto.CreateUser  true  "User struct"
// @Success 200 {object} entity.User
// @Failure 400 {object} entity.Err
// @Failure 500 {object} entity.Err
// @Router /register [post]
func (a *AuthHandler) Register(c *gin.Context) {
	var data dto.CreateUser
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, entity.Err{
			Error:   http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		})
		return
	}

	user, err := a.service.Register(&data)
	if err != nil {
		var errsValidation *entity.ValidationErrors
		switch {
		case errors.As(err, &errsValidation):
			c.JSON(http.StatusBadRequest, entity.Err{
				Error:   http.StatusText(http.StatusBadRequest),
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, entity.Err{
				Error:   http.StatusText(http.StatusInternalServerError),
				Message: "Internal server error",
			})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Login
// @Tags auth
// @Description login (get auth token)
// @Accept json
// @Produce json
// @Param input body dto.Auth  true  "Auth struct"
// @Success 200 {object} entity.Token
// @Failure 400 {object} entity.Err
// @Failure 500 {object} entity.Err
// @Router /login [post]
func (a *AuthHandler) Login(c *gin.Context) {
	var data dto.Auth

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, entity.Err{
			Error:   http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		})
		return
	}

	token, err := a.service.Login(&data)
	if err != nil {
		var errsValidation *entity.ValidationErrors
		switch {
		case errors.As(err, &errsValidation):
			c.JSON(http.StatusBadRequest, entity.Err{
				Error:   http.StatusText(http.StatusBadRequest),
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, entity.Err{
				Error:   http.StatusText(http.StatusInternalServerError),
				Message: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, entity.Token{Token: token})
}
