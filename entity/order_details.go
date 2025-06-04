package entity

import "time"

type OrderDetail struct {
	ID			int
	OrderID 	int
	Order		Order
	ProductID	int
	Product		Product
	Qty			int
	Total			float64
	CreatedAt	time.Time
	UpdatedAt	time.Time
	CreatedBy	int
	UpdatedBy	int
}