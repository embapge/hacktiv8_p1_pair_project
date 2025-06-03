package entity

import "time"

type StatusOrder string

const (
	StatusProcessing		StatusOrder = "processing"
	StatusCompleted			StatusOrder = "completed"
	StatusCancel			StatusOrder = "cancel"
)

type Order struct{
	ID				int
	CustomerID		int
	NumberDisplay	string
	Date			time.Time
	Status			StatusOrder
	Total			float64
	CreatedAt		time.Time
	UpdatedAt		time.Time
	CreatedBy		int
	UpdatedBy		int
}