package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"pairproject/entity"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type AuthHandler struct {
	DB *sql.DB
}

// // Fungsi untuk menambah user baru (contoh, bisa diubah sesuai kebutuhan)
func (h *AuthHandler) Register(cust *entity.CustomerRegister) error {
	// Insert user and return the inserted user ID
	resultUser, err := h.DB.Exec(
		"INSERT INTO users(username, email, password, role) VALUES (?, ?, ?, 'customer')",
		strings.TrimSpace(cust.Username), strings.TrimSpace(cust.Email), cust.Password,
	)

	if err != nil {
		// Cek apakah error dari driver MySQL
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				// Duplicate entry
				return fmt.Errorf("duplicate entry. Silahkan gunakan data lainnya")
			}
		}
		// Error lain
		return err
	}

	userID, err := resultUser.LastInsertId()
	if err != nil {
		return err
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
		// Cek apakah error dari driver MySQL
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				// Duplicate entry
				return fmt.Errorf("duplicate entry. Silahkan gunakan data lainnya")
			}
		}
		// Error lain
		return err
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	_, err = h.DB.Exec(
		"INSERT INTO user_customers(customer_id, user_id) VALUES (?, ?)",
		customerID, userID,
	)

	if err != nil {
		return errors.New("Terjadi kesalahan dalam mendaftarkan akun.")
	}

	return nil
}

// Fungsi untuk login user
func (h *AuthHandler) LoginUser(username, password string) (*entity.User, error) {
	var user entity.User
	var customerID sql.NullInt64

	err := h.DB.QueryRow(
		"SELECT id, username, email, role, password, user_customers.customer_id FROM users LEFT JOIN user_customers on user_customers.user_id = users.id WHERE username = ? AND password = ?", 
		strings.TrimSpace(username), strings.TrimSpace(password),
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Password, &customerID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Password atau username salah")
		}
		
		return nil, err
	}

	if customerID.Valid {
		user.Customer = entity.Customer{ID: int(customerID.Int64)}
	} else {
		user.Customer = entity.Customer{}
	}	

	return &user, nil
}