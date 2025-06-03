package entity

import "time"

type Method string

const (
	MethodCredit			Method = "credit_card"
	MethodVA				Method = "va"
	MethodTransfer			Method = "transfer"
)

type Payment struct {
	ID			int
	BillingID	int
	Date		time.Time
	Amount		float64
	Method		Method
	CreatedAt	time.Time
	UpdatedAt	time.Time
	CreatedBy	int
	UpdatedBy	int
}