package auth

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"pairproject/handler"
	"pairproject/utils"
	"strings"
)

// Login menangani proses login pengguna dengan input username dan password dari terminal.
// Fungsi ini akan mengembalikan context yang sudah ditambahkan informasi user jika login berhasil.
func Login(db *sql.DB, ctx *context.Context) (context.Context, error) {
	// Membuat reader untuk membaca input dari terminal
	reader := bufio.NewReader(os.Stdin)

	// Minta pengguna untuk memasukkan username
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')

	// Minta pengguna untuk memasukkan password
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')

	// Menghapus spasi dan newline dari input
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	// Membuat instance dari handler untuk autentikasi
	userHandler := handler.AuthHandler{DB: db}

	// Memanggil fungsi LoginUser untuk memverifikasi username dan password
	user, err := userHandler.LoginUser(username, password)
	if err != nil {
		// Jika gagal login, kembalikan context lama dan error
		return *ctx, err
	}

	// Jika login berhasil, tambahkan data user ke dalam context dan kembalikan
	return utils.WithUser(*ctx, user), nil
}