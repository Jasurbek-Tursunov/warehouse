package gin

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Middleware struct {
	service usecase.AuthService
}

func NewMiddleware(service usecase.AuthService) *Middleware {
	return &Middleware{service: service}
}

func (m *Middleware) Auth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "auth header require"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		return
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "auth type unsupported"})
		return
	}

	if err := m.service.Check(headerParts[1]); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
}
