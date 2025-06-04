package utils

import (
	"context"
	"pairproject/entity"
)

type ContextKey string // Penambahan tipe string contextKey menghindari penggunan key yang sama pada context 

const userKey ContextKey = "user"

func WithUser(ctx context.Context, user *entity.User) context.Context {
    return context.WithValue(ctx, userKey, user)
}

func GetUser(ctx context.Context) (*entity.User, bool) {
    user, ok := ctx.Value(userKey).(*entity.User) // Menggunakan type assertion memastikan bahwa type datanya merupakan entity.User. Jika berbeda maka akan return false
    return user, ok
}

func ClearUser(ctx context.Context) context.Context {
    return context.WithValue(ctx, userKey, nil)
}