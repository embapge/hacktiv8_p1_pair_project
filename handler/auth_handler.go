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

// Register menambahkan user baru sekaligus data customer terkait.
// Menerima input struct CustomerRegister dan mengembalikan error jika ada.
func (h *AuthHandler) Register(cust *entity.CustomerRegister) error {
	// Insert data user ke tabel users dengan role 'customer'
	resultUser, err := h.DB.Exec(
		"INSERT INTO users(username, email, password, role) VALUES (?, ?, ?, 'customer')",
		strings.TrimSpace(cust.Username), strings.TrimSpace(cust.Email), cust.Password,
	)

	if err != nil {
		// Jika error dari MySQL, cek apakah error duplikat entry (1062)
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				// Username atau email sudah terdaftar
				return fmt.Errorf("duplicate entry. Silahkan gunakan data lainnya")
			}
		}
		// Error lain dikembalikan apa adanya
		return err
	}

	// Ambil ID user yang baru saja dibuat
	userID, err := resultUser.LastInsertId()
	if err != nil {
		return err
	}

	// Insert data customer ke tabel customers dengan relasi created_by ke userID
	result, err := h.DB.Exec(
		"INSERT INTO customers(name, address, email, phone_number, created_by) VALUES (?, ?, ?, ?, ?)",
		strings.TrimSpace(cust.Name),
		strings.TrimSpace(cust.Address),
		strings.TrimSpace(cust.Email),
		strings.TrimSpace(cust.Phone),
		userID,
	)
	
	if err != nil {
		// Cek error duplikat entry pada customer juga
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return fmt.Errorf("duplicate entry. Silahkan gunakan data lainnya")
			}
		}
		return err
	}

	// Ambil ID customer yang baru dibuat
	customerID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Hubungkan user dan customer pada tabel user_customers
	_, err = h.DB.Exec(
		"INSERT INTO user_customers(customer_id, user_id) VALUES (?, ?)",
		customerID, userID,
	)

	if err != nil {
		return errors.New("Terjadi kesalahan dalam mendaftarkan akun.")
	}

	return nil
}

// LoginUser melakukan validasi login dengan username dan password.
// Mengembalikan objek User dan error jika login gagal.
func (h *AuthHandler) LoginUser(username, password string) (*entity.User, error) {
	var user entity.User
	var customerID sql.NullInt64

	// Query user berdasarkan username dan password yang diberikan
	// Menggunakan LEFT JOIN untuk mengambil customer_id jika ada
	err := h.DB.QueryRow(
		"SELECT id, username, email, role, password, user_customers.customer_id FROM users LEFT JOIN user_customers on user_customers.user_id = users.id WHERE username = ? AND password = ?", 
		strings.TrimSpace(username), strings.TrimSpace(password),
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Password, &customerID)

	if err != nil {
		if err == sql.ErrNoRows {
			// Jika tidak ada data, berarti username/password salah
			return nil, fmt.Errorf("Password atau username salah")
		}
		
		return nil, err
	}

	// Jika customerID valid, set properti Customer pada user
	if customerID.Valid {
		user.Customer = entity.Customer{ID: int(customerID.Int64)}
	} else {
		user.Customer = entity.Customer{}
	}	

	return &user, nil
}