package entity

import "time"

type Category struct{
	ID			int
	Name		string
	CreatedAt	time.Time
	UpdatedAt	time.Time
	CreatedBy	time.Time
	UpdatedBy	time.Time	
}