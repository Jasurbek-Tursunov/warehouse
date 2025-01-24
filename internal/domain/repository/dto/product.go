package dto

type CreateProduct struct {
	Name     string `json:"name" validate:"required"`
	Price    uint   `json:"price" validate:"required"`
	Quantity uint   `json:"quantity" validate:"required"`
}

type UpdateProduct struct {
	Name     string `json:"name" validate:"required"`
	Price    uint   `json:"price" validate:"required"`
	Quantity uint   `json:"quantity" validate:"required"`
}

type ProductQuery struct {
}
