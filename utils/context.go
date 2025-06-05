package utils

import (
	"context"
	"pairproject/entity"
)

// ContextKey adalah tipe kustom untuk key yang digunakan di context agar menghindari bentrok dengan key lain.
type ContextKey string

// userKey adalah key unik untuk menyimpan dan mengambil data user pada context.
const userKey ContextKey = "user"

// WithUser menyisipkan data user ke dalam context dengan key userKey.
// Mengembalikan context baru yang mengandung user.
func WithUser(ctx context.Context, user *entity.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetUser mengambil data user dari context menggunakan key userKey.
// Menggunakan type assertion agar memastikan tipe data benar *entity.User.
// Jika data tidak ada atau tipe salah, mengembalikan nil dan false.
func GetUser(ctx context.Context) (*entity.User, bool) {
	user, ok := ctx.Value(userKey).(*entity.User)
	return user, ok
}

// ClearUser menghapus data user dari context dengan menyisipkan nilai nil untuk key userKey.
// Mengembalikan context baru tanpa data user.
func ClearUser(ctx context.Context) context.Context {
	return context.WithValue(ctx, userKey, nil)
}