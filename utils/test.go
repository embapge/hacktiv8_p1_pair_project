package utils

import (
	"context"
	"pairproject/entity"
)

const userTestKey ContextKey = "user"

func NewTestContextWithUser() context.Context {
	user := &entity.User{
		ID: 1,
		Customer: entity.Customer{
			ID: 1,
		},
	}
	ctx := context.Background()
	return context.WithValue(ctx, userTestKey, user)
}