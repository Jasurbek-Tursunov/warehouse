package dto

type CreateProduct struct {
	Name     string  `json:"name" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	Quantity int     `json:"quantity" validate:"required"`
}

type UpdateProduct struct {
	Name     string  `json:"name" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	Quantity int     `json:"quantity" validate:"required"`
}

type ProductQuery struct {
}
