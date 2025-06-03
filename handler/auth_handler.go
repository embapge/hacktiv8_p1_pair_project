package handler

import (
	"database/sql"
	"fmt"
	"log"
	"pairproject/entity"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

// // Fungsi untuk menambah user baru (contoh, bisa diubah sesuai kebutuhan)
func (h *AuthHandler) Register(username, email, password, role string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Gagal hash password:", err)
		return
	}
	_, err = h.DB.Exec("INSERT INTO users(username, email, password, role) VALUES (?,?,?,?)",
		strings.TrimSpace(username), strings.TrimSpace(email), hashedPassword, strings.TrimSpace(role))
	if err != nil {
		log.Println("Gagal tambah user:", err)
	} else {
		fmt.Println("User berhasil ditambahkan!")
	}
}

// Fungsi untuk login user
func (h *AuthHandler) LoginUser(username, password string) (*entity.User, error) {
	var user entity.User

	err := h.DB.QueryRow(
		"SELECT id, username, email, role, password FROM users WHERE username = ? AND password = ?", 
		strings.TrimSpace(username), strings.TrimSpace(password),
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Password)

	fmt.Println("Ketik User:", username)
	fmt.Println("Hashed User:", user)
	// fmt.Println("Hashed from DB:", hashedPassword)
	// fmt.Println("Input password:", password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Password atau username salah")
		}
		
		return nil, err
	}

	return &user, nil
}