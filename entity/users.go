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
	Customer	Customer
	Password	string
	CreateAt	time.Time
	UpdatedAt 	time.Time
}

type CustomerRegister struct{
	Name     string
	Address  string
	Email    string
	Phone    string
	Username string
	Password string
	Role     string // "admin" atau "customer"
}