package entity

import (
	"time"
)

type StatusOrder string

const (
	StatusProcessing		StatusOrder = "processing"
	StatusCompleted			StatusOrder = "completed"
	StatusCancel			StatusOrder = "cancel"
)

type Order struct{
	ID				int
	CustomerID		int
	Customer		Customer
	NumberDisplay	string
	Date			time.Time
	Status			StatusOrder
	Total			float64
	Details			[]OrderDetail
	CreatedAt		time.Time
	UpdatedAt		time.Time
	CreatedBy		int
	UpdatedBy		int
}

type OrderProduct struct{
	ProductId, Qty int
}