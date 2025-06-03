package entity

import "time"

type UserCustomer struct{
	UserID		int
	CustomerID	int
	CreatedAt	time.Time
	UpdatedAt	time.Time
}