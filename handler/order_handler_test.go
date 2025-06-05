package handler

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"pairproject/entity"
	"pairproject/utils"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite" // Driver SQLite untuk testing
)

// SetupTestOrderDB membuat database in-memory SQLite untuk keperluan unit testing.
// Tabel yang dibuat: orders, order_details, products, dan trigger untuk menghitung total otomatis.
func SetupTestOrderDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed open db: %v", err)
	}

	currentYearMonth := time.Now().Format("200601") // Format YYYYMM

	// Buat tabel, trigger, dan isi data awal
	initialDB := fmt.Sprintf(`
		CREATE TABLE orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			number_display TEXT,
			customer_id INTEGER,
			date DATETIME DEFAULT CURRENT_TIMESTAMP,
			status TEXT DEFAULT 'processing' CHECK (status IN ('processing', 'completed', 'cancel')),
			created_by INTEGER,
			total REAL DEFAULT 0
		);

		INSERT INTO orders (number_display, customer_id, created_by)
		VALUES ('ORD-%s-001', 1, 1);

		CREATE TABLE order_details (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER,
			product_id INTEGER,
			qty INTEGER,
			created_by INTEGER
		);

		CREATE TABLE products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			stock INTEGER DEFAULT 0 NOT NULL,
			description TEXT,
			category_id INTEGER NOT NULL,
			price REAL DEFAULT 0 NOT NULL
		);

		INSERT INTO products (name, stock, description, category_id, price)
		VALUES 
		('Adjustable Dumbbell 20kg', 10, 'Customizable weight dumbbell for home workouts', 1, 100000.00),
		('Treadmill Compact X100', 4, 'Foldable treadmill with digital display', 1, 50000.00),
		('Tent 2-Person Waterproof', 8, 'Outdoor tent ideal for camping', 2, 875000.00),
		('Camping Stove Mini', 15, 'Portable gas stove for outdoor use', 2, 320000.00),
		('Whey Protein 1kg', 20, 'Chocolate flavor protein supplement', 3, 450000.00),
		('Electrolyte Drink Pack (12x)', 30, 'Hydration booster during exercise', 3, 180000.00);

		CREATE TRIGGER trg_order_details_after_insert
		AFTER INSERT ON order_details
		BEGIN
			UPDATE orders
			SET total = (
				SELECT IFNULL(SUM(od.qty * p.price), 0)
				FROM order_details od
				JOIN products p ON od.product_id = p.id
				WHERE od.order_id = NEW.order_id
			)
			WHERE id = NEW.order_id;
		END;
	`, currentYearMonth)

	_, err = db.Exec(initialDB)
	if err != nil {
		t.Fatalf("failed create schema: %v", err)
	}

	return db
}

// TestCreateOrder menguji proses pembuatan order dan menghitung total dari order_details.
func TestCreateOrder(t *testing.T) {
	db := SetupTestOrderDB(t)
	defer db.Close()

	ctx := utils.NewTestContextWithUser()
	handler := &OrderHandler{DB: db, Ctx: &ctx}

	// Produk yang ingin dipesan
	orderProducts := []entity.OrderProduct{
		{ProductId: 1, Qty: 2}, // 2 x 100000 = 200000
		{ProductId: 2, Qty: 3}, // 3 x 50000  = 150000
	}

	order, err := handler.CreateOrder(orderProducts)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	assert.NotEqual(t, order.ID, 0, "Order ID harus berhasil dibuat")
	assert.Equal(t, len(order.Details), len(orderProducts), "Jumlah produk harus sesuai")

	// Validasi order berdasarkan nomor
	currentYearMonth := time.Now().Format("200601")
	order, err = handler.GetOrderByNumberDisplay(fmt.Sprintf("ORD-%s-002", currentYearMonth))
	if err != nil {
		t.Fatalf("GetOrderByNumberDisplay failed: %v", err)
	}

	assert.Equal(t, float64(350000), order.Total, "Total order harus sesuai (200000 + 150000)")
}

// TestGenerateOrderNumber menguji penomoran otomatis order berdasarkan nomor terakhir di database.
func TestGenerateOrderNumber(t *testing.T) {
	db := SetupTestOrderDB(t)
	defer db.Close()

	ctx := context.Background()
	handler := &OrderHandler{DB: db, Ctx: &ctx}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed begin tx: %v", err)
	}

	num := handler.GenerateOrderNumber(tx)
	if num == "" {
		t.Errorf("GenerateOrderNumber returned empty string")
	}

	currentYearMonth := time.Now().Format("200601")
	expected := fmt.Sprintf("ORD-%s-002", currentYearMonth)

	assert.Equal(t, expected, num, "Nomor order harus naik jadi ORD-YYYYMM-002")

	if err := tx.Commit(); err != nil {
		t.Fatalf("failed commit tx: %v", err)
	}
}