package handler

import (
	"database/sql"
	"pairproject/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite" // SQLite driver CGO-less
)

func SetupTestAuthDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err, "failed to open in-memory SQLite DB")

	_, err = db.Exec(`
		CREATE TABLE users ( 
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			username VARCHAR(100) NOT NULL UNIQUE, 
			email VARCHAR(100) NOT NULL UNIQUE, 
			password VARCHAR(255) NOT NULL, 
			role TEXT NOT NULL CHECK (role IN ('admin', 'staff', 'customer')), 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		); 

		INSERT INTO users (username, email, password, role) VALUES
		('admin01', 'admin01@example.com', '123456', 'admin'),
		('staff01', 'staff01@example.com', '123456', 'staff'),
		('staff02', 'staff02@example.com', '123456', 'staff'),
		('custuser1', 'cust1@example.com', '123456', 'customer'),
		('custuser2', 'cust2@example.com', '123456', 'customer');

		CREATE TABLE customers ( 
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			name VARCHAR(100) NOT NULL, 
			address TEXT NOT NULL, 
			email VARCHAR(100) NOT NULL UNIQUE, 
			phone_number VARCHAR(20) NOT NULL UNIQUE, 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER, 
			updated_by INTEGER, 
			FOREIGN KEY (created_by) REFERENCES users(id), 
			FOREIGN KEY (updated_by) REFERENCES users(id)
		); 

		CREATE TABLE user_customers ( 
			user_id INTEGER NOT NULL, 
			customer_id INTEGER NOT NULL PRIMARY KEY, 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id), 
			FOREIGN KEY (customer_id) REFERENCES customers(id) 
		); 

		INSERT INTO customers (name, address, email, phone_number, created_by, updated_by)
		VALUES 
		('John Doe', 'Jl. Merdeka No. 123, Jakarta', 'john@example.com', '081234567890', 2, 2),
		('Jane Smith', 'Jl. Sudirman No. 10, Bandung', 'jane@example.com', '089876543210', 3, 3);

		INSERT INTO user_customers (user_id, customer_id)
		VALUES 
		(4, 1),
		(5, 2);
	`)
	require.NoError(t, err, "failed to initialize schema and data")

	return db
}

func TestRegister_SUCCESS(t *testing.T){
	db := SetupTestAuthDB(t)
	handler := &AuthHandler{DB: db}
	customer := entity.CustomerRegister{Name: "Mohammad Barata", Address: "KP", Email: "mohammadbarata.mb@gmail.com", Phone: "08922387737", Username: "embapge", Password: "1234567"}

	err := handler.Register(&customer)

	assert.NoError(t, err)
}

func TestRegister_FAILED(t *testing.T){
	db := SetupTestAuthDB(t)
	handler := &AuthHandler{DB: db}

	// Name     string
	// Address  string
	// Email    string
	// Phone    string
	// Username string
	// Password string
	// Role     string // "admin" atau "customer"

	customer := entity.CustomerRegister{Name: "Didit", Address: "KP", Email: "mohammadbarata.mb@gmail.com", Phone: "08922387737", Username: "custuser1", Password: "1234567"}

	err := handler.Register(&customer)

	assert.Contains(t, err.Error(),"UNIQUE", "Terjadi error karena username telah terdaftar")
}

func TestLoginUser_SUCCESS(t *testing.T){
	db := SetupTestAuthDB(t)
	handler := &AuthHandler{DB: db}
	user, err := handler.LoginUser("admin01", "123456")
	if err != nil{
		t.Fatal(err)
	}
	assert.NotEqual(t, user.ID, 0, "Berhasil login")
	user, err = handler.LoginUser("staff01", "123456")
	if err != nil{
		t.Fatal(err)
	}

	assert.NotEqual(t, user.ID, 0, "Berhasil login")
	user, err = handler.LoginUser("custuser2", "123456")
	if err != nil{
		t.Fatal(err)
	}

	assert.NotEqual(t, user.ID, 0, "Berhasil login")
}

func TestLoginUser_FAILED(t *testing.T){
	db := SetupTestAuthDB(t)
	handler := &AuthHandler{DB: db}
	_, err := handler.LoginUser("admin01", "1234567")
	
	assert.Contains(t, err.Error(), "Password atau username salah", "Berhasil login")
}
