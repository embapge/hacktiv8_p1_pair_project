package entity

import "time"

type Users struct {
	ID			int
	Username	string
	Email		string
	Password	string
	Create_at	time.Time
	Updated_at 	time.Time
}