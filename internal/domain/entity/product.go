package entity

import "time"

type Product struct {
	ID        int
	Name      string
	Price     float64
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
