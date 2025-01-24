package entity

import "time"

type Product struct {
	ID        uint
	name      string
	price     float64
	quantity  int
	createdAt time.Time
	updatedAt time.Time
}
