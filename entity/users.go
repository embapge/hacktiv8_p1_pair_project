package entity

import "time"

type Role string

const (
	RoleAdmin		Role = "admin"
	RoleCustomer	Role = "customer"
)

type User struct {
	ID			int
	Username	string
	Email		string
	Role		Role
	Password	string
	CreateAt	time.Time
	UpdatedAt 	time.Time
}