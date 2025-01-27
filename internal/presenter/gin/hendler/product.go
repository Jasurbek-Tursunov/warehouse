package hendler

import (
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/repository/dto"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	service usecase.ProductService
}

func NewProductHandler(service usecase.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (p *ProductHandler) List(c *gin.Context) {
	filters := dto.ProductQuery{}
	paginate := dto.Paginator{}

	products, err := p.service.List(&filters, &paginate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, products)
}
func (p *ProductHandler) Create(c *gin.Context) {
	var data dto.CreateProduct
	if err := c.BindJSON(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := p.service.Create(&data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, product)
}

func (p *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var data dto.UpdateProduct

	product, err := p.service.Update(id, &data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, product)
}

func (p *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	product, err := p.service.Get(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, product)
}
func (p *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err = p.service.Delete(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
