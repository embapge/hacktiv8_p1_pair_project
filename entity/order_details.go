package entity

import "time"

type OrderDetail struct {
	ID			int
	OrderID 	int
	ProductID	int
	Qty			int
	CreatedAt	time.Time
	UpdatedAt	time.Time
	CreatedBy	int
	UpdatedBy	int
}