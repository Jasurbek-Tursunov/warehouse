package hendler

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service usecase.ProductService
}

func NewProductHandler(service usecase.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (p *ProductHandler) List(c *gin.Context)   {}
func (p *ProductHandler) Create(c *gin.Context) {}
func (p *ProductHandler) Update(c *gin.Context) {}
func (p *ProductHandler) Get(c *gin.Context)    {}
func (p *ProductHandler) Delete(c *gin.Context) {}
