package handler

import (
	"context"
	"database/sql"
	"testing"

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

	// Buat schema tabel minimal untuk test
	schema := `
	CREATE TABLE orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		number_display TEXT,
		customer_id INTEGER,
		date TEXT,
		created_by INTEGER
	);
	CREATE TABLE order_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id INTEGER,
		product_id INTEGER,
		qty INTEGER,
		created_by INTEGER
	);
	`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("failed create schema: %v", err)
	}
	return db
}
func TestCreateOrder(t *testing.T) {
	db := SetupTestOrderDB(t)
	ctx := utils.NewTestContextWithUser()
	handler := &OrderHandler{DB: db, Ctx: &ctx}

	orderProducts := []entity.OrderProduct{
		{ProductId: 1, Qty: 2},
		{ProductId: 2, Qty: 3},
	}

	order, err := handler.CreateOrder(orderProducts)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	assert.NotEqual(t, order.ID, 0, "Order Id berhasil dibuat")
	assert.Equal(t, len(order.Details), len(orderProducts), "Jumlah product terbuat sesuai yaitu %d", len(orderProducts))
}

func TestGenerateOrderNumber(t *testing.T) {
	db := SetupTestOrderDB(t)
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

	assert.NotEmpty(t, num)
	assert.Equal(t, "ORD-202506-001", num, "Hasil perjumlahan harus 5")

	err = tx.Commit()
	if err != nil {
		t.Fatalf("failed commit tx: %v", err)
	}
}
