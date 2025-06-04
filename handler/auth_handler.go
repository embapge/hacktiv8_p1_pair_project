package handler

import (
	"database/sql"
	"fmt"
	"pairproject/entity"
	"strings"
)

type AuthHandler struct {
	DB *sql.DB
}

// // Fungsi untuk menambah user baru (contoh, bisa diubah sesuai kebutuhan)
func (h *AuthHandler) Register(cust *entity.CustomerRegister) (string, error) {
	// Insert user and return the inserted user ID
	resultUser, err := h.DB.Exec(
		"INSERT INTO users(username, email, password, role) VALUES (?, ?, ?, 'customer')",
		strings.TrimSpace(cust.Username), strings.TrimSpace(cust.Email), cust.Password,
	)

	if err != nil {
		return "", err
	}

	userID, err := resultUser.LastInsertId()
	if err != nil {
		return "", err
	}

	result, err := h.DB.Exec(
		"INSERT INTO customers(name, address, email, phone_number, created_by) VALUES (?, ?, ?, ?, ?)",
		strings.TrimSpace(cust.Name),
		strings.TrimSpace(cust.Address),
		strings.TrimSpace(cust.Email),
		strings.TrimSpace(cust.Phone),
		userID,
	)
	
	if err != nil {
		return "Gagal masuk data customer", err
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		return "Gagal mendapatkan ID customer", err
	}

	_, err = h.DB.Exec(
		"INSERT INTO user_customers(customer_id, user_id) VALUES (?, ?)",
		customerID, userID,
	)

	if err != nil {
		return "", err
	}

	return "Data berhasil masuk", nil
}

// Fungsi untuk login user
func (h *AuthHandler) LoginUser(username, password string) (*entity.User, error) {
	var user entity.User

	err := h.DB.QueryRow(
		"SELECT id, username, email, role, password, user_customers.customer_id FROM users JOIN user_customers on user_customers.user_id = users.id WHERE username = ? AND password = ?", 
		strings.TrimSpace(username), strings.TrimSpace(password),
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Password, &user.Customer.ID)

	fmt.Println("Ketik User:", username)
	fmt.Println("Hashed User:", user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Password atau username salah")
		}
		
		return nil, err
	}

	return &user, nil
}