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
	_ "modernc.org/sqlite" // contoh driver SQLite untuk test
)

func SetupTestOrderDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed open db: %v", err)
	}

	currentYearMonth := time.Now().Format("200601")

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

	// Buat schema tabel minimal untuk test
	schema := initialDB
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("failed create schema: %v", err)
	}
	return db
}
func TestCreateOrder(t *testing.T) {
	db := SetupTestOrderDB(t)
	defer db.Close()
	ctx := utils.NewTestContextWithUser()
	handler := &OrderHandler{DB: db, Ctx: &ctx}

	orderProducts := []entity.OrderProduct{
		{ProductId: 1, Qty: 2}, // 200000
		{ProductId: 2, Qty: 3}, // 150000
	}

	order, err := handler.CreateOrder(orderProducts)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	assert.NotEqual(t, order.ID, 0, "Order Id berhasil dibuat")
	assert.Equal(t, len(order.Details), len(orderProducts), "Jumlah product terbuat sesuai yaitu %d", len(orderProducts))
	
	currentYearMonth := time.Now().Format("200601")
	order, err = handler.GetOrderByNumberDisplay(fmt.Sprintf("ORD-%s-002", currentYearMonth))
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	
	assert.Equal(t, order.Total, float64(350000))
}

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

	assert.NotEmpty(t, num)
	assert.Equal(t, fmt.Sprintf("ORD-%s-002", string(currentYearMonth)), num, "Penomoran order sudai sesuai")

	err = tx.Commit()
	if err != nil {
		t.Fatalf("failed commit tx: %v", err)
	}
}
