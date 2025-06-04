package entity

import "time"


type StatusBilling string

const (
	StatusUnpaid		StatusBilling = "unpaid"
	StatusPaid			StatusBilling = "paid"
	StatusCancelled		StatusBilling = "cancelled"
	StatusRefunded		StatusBilling = "refunded"
)

type Billing struct {
	ID				int
	OrderID 		int
	IssueDate 		time.Time
	DueDate 		time.Time
	NumberDisplay	string
	Tax				float64
	Total			float64
	Status			StatusBilling
	CreatedAt		time.Time
	UpdatedAt		time.Time
	CreatedBy		int
	UpdatedBy		int
}