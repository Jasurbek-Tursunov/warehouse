package hendler

import (
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

func (a *AuthHandler) Register(c *gin.Context) {
	var data dto.CreateUser
	if err := c.BindJSON(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.service.Register(data.Username, data.Password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func (a *AuthHandler) Login(c *gin.Context) {

}
