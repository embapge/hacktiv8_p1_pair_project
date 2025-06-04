package entity

import "time"

type Product struct{
	ID			int
	Name		string
	Stock		int
	Description	string
	CategoryID	int
	Category	Category
	Price		float64
	CreatedAt	time.Time
	UpdatedAt	time.Time
	CreatedBy	int
	UpdatedBy	int
}