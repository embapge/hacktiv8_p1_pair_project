package entity

import "time"

type Customer struct {
	ID				int
	Name			string
	Address			string
	Email			string
	PhoneNumber		string
	CreatedAt		time.Time
	UpdatedAt		time.Time
}