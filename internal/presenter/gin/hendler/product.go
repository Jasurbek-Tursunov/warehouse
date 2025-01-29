package hendler

import (
	"errors"
	"github.com/Jasurbek-Tursunov/warehouse/internal/domain/entity"
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

// @Security Bearer
// @Summary List
// @Tags product
// @Description list products data
// @Accept json
// @Produce json
// @Param name query string  false  "Search by product name"
// @Param sort_by query string  false  "Sort field"
// @Param limit query int false  "Limit for paginate"
// @Param page query int false  "Page for paginate"
// @Success 200 {object} []entity.Product
// @Failure 500 {object} entity.Err
// @Router /products [get]
func (p *ProductHandler) List(c *gin.Context) {
	filters := dto.ProductQuery{
		Name:   c.Query("name"),
		SortBy: c.Query("sort_by"),
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit < 1 {
		limit = 10
	}

	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}

	paginate := dto.Paginator{
		PageSize: limit,
		Page:     page,
	}

	products, err := p.service.List(&filters, &paginate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Err{
			Error:   http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, products)
}

// @Security Bearer
// @Summary Create
// @Tags product
// @Description create product
// @Accept json
// @Produce json
// @Param input body dto.CreateProduct  true  "Product struct"
// @Success 200 {object} entity.Product
// @Failure 400 {object} entity.Err
// @Failure 500 {object} entity.Err
// @Router /product/add [post]
func (p *ProductHandler) Create(c *gin.Context) {
	var data dto.CreateProduct
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, entity.Err{
			Error:   http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		})
		return
	}

	product, err := p.service.Create(&data)
	if err != nil {
		var errValidation *entity.ValidationError
		switch {
		case errors.As(err, &errValidation):
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

	c.JSON(http.StatusOK, product)
}

// @Security Bearer
// @Summary Get
// @Tags product
// @Description get product data
// @Accept json
// @Produce json
// @Param id path int  true  "Product ID"
// @Success 200 {object} entity.Product
// @Failure 404 {object} entity.Err
// @Failure 500 {object} entity.Err
// @Router /product/{id} [get]
func (p *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Err{
			Error:   http.StatusText(http.StatusNotFound),
			Message: "object with this id not found",
		})
		return
	}

	product, err := p.service.Get(id)
	if err != nil {
		var errNotFound *entity.NotFoundError
		switch {
		case errors.As(err, &errNotFound):
			c.JSON(http.StatusNotFound, entity.Err{
				Error:   http.StatusText(http.StatusNotFound),
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

	c.JSON(http.StatusOK, product)
}

// @Security Bearer
// @Summary Update
// @Tags product
// @Description update product data
// @Accept json
// @Produce json
// @Param id path int  true  "Product ID"
// @Param input body dto.UpdateProduct  true  "Product struct"
// @Success 200 {object} entity.Product
// @Failure 400 {object} entity.Err
// @Failure 404 {object} entity.Err
// @Failure 500 {object} entity.Err
// @Router /product/{id} [put]
func (p *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Err{
			Error:   http.StatusText(http.StatusNotFound),
			Message: entity.NewNotFoundError("product", id).Error(),
		})
		return
	}

	var data dto.UpdateProduct
	if err = c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, entity.Err{
			Error:   http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		})
		return
	}

	product, err := p.service.Update(id, &data)
	if err != nil {
		var errNotFound *entity.NotFoundError
		var errValidation *entity.ValidationError
		switch {
		case errors.As(err, &errNotFound):
			c.JSON(http.StatusNotFound, entity.Err{
				Error:   http.StatusText(http.StatusNotFound),
				Message: err.Error(),
			})
		case errors.As(err, &errValidation):
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

	c.JSON(http.StatusOK, product)
}

// @Security Bearer
// @Summary Delete
// @Tags product
// @Description delete product data
// @Accept json
// @Produce json
// @Param id path int  true  "Product ID"
// @Success 204
// @Failure 404 {object} entity.Err
// @Failure 500 {object} entity.Err
// @Router /product/{id} [delete]
func (p *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Err{
			Error:   http.StatusText(http.StatusNotFound),
			Message: entity.NewNotFoundError("product", id).Error(),
		})
		return
	}

	err = p.service.Delete(id)
	if err != nil {
		var errNotFound *entity.NotFoundError
		switch {
		case errors.As(err, &errNotFound):
			c.JSON(http.StatusNotFound, entity.Err{
				Error:   http.StatusText(http.StatusNotFound),
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

	c.Status(http.StatusNoContent)
}
