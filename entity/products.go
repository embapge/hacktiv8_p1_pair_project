package entity

import "time"

type Product struct{
	ID			int
	Name		int
	Stock		int
	Description	string
	CategoryID	int
	Price		float64
	CreatedAt	time.Time
	UpdatedAt	time.Time
	CreatedBy	int
	UpdatedBy	int
}