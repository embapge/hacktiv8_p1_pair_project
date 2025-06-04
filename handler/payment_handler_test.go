package handler

import (
	"database/sql"
	"pairproject/entity"
	"pairproject/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupTestPaymentDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed open db: %v", err)
	}

	// Buat schema tabel minimal untuk test
	schema := `
	-- Tabel billings (SQLite3)
CREATE TABLE billings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    billing_id INTEGER NOT NULL,
    due_date DATETIME NOT NULL,
    amount NUMERIC DEFAULT 0 CHECK (amount >= 0) NOT NULL,
    method TEXT NOT NULL CHECK (method IN ('credit_card', 'va', 'transfer')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL
);

-- Contoh insert data
INSERT INTO billings (billing_id, due_date, amount, method, created_by)
VALUES (1001, '2025-06-05 20:21:13', 10000000.00, 'transfer', 1);

-- Tabel payments (SQLite3)
CREATE TABLE payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    billing_id INTEGER NOT NULL,
    date DATETIME DEFAULT CURRENT_TIMESTAMP,
    amount NUMERIC DEFAULT 0 CHECK (amount >= 0) NOT NULL,
    method TEXT NOT NULL CHECK (method IN ('credit_card', 'va', 'transfer')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    FOREIGN KEY (billing_id) REFERENCES billings(id) ON DELETE CASCADE
);

CREATE TRIGGER update_billings_updated_at
AFTER UPDATE ON billings
FOR EACH ROW
BEGIN
    UPDATE billings SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE TRIGGER validate_payment_amount
BEFORE INSERT ON payments
FOR EACH ROW
BEGIN
    -- Jika total pembayaran melebihi jumlah tagihan, hentikan proses
    SELECT 
        CASE 
            WHEN (
                (SELECT IFNULL(SUM(amount), 0) FROM payments WHERE billing_id = NEW.billing_id)
                + NEW.amount
            ) > (SELECT amount FROM billings WHERE id = NEW.billing_id)
            THEN RAISE(ABORT, 'Total payment exceeds billing total')
        END;
END;
	`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("failed create schema: %v", err)
	}
	return db
}
func TestCreatePayment_Success(t *testing.T){
	db := SetupTestPaymentDB(t)
	ctx := utils.NewTestContextWithUser()
	handler := &PaymentHandler{DB: db, Ctx: &ctx}
	var method entity.Method = "va"
	billing := entity.Billing{ID: 1, DueDate: time.Now().Add(1 * time.Minute)}
	err := handler.CreatePayment(billing, 1000000.0, method)

	if err != nil{
		t.Fatalf("failed create payment: %v", err)
	}

	assert.Equal(t, true, true, "Payment berhasil dibuat")

	billing = entity.Billing{ID: 1, DueDate: time.Now().Truncate(1 * time.Minute)}
	err = handler.CreatePayment(billing, 1000000.0, method)

	assert.Contains(t, err.Error(), "cannot create payment: order is past due date")
}
func TestCreatePayment_Error(t *testing.T){
	db := SetupTestPaymentDB(t)
	ctx := utils.NewTestContextWithUser()
	handler := &PaymentHandler{DB: db, Ctx: &ctx}
	var method entity.Method = "va"
	billing := entity.Billing{ID: 1, DueDate: time.Now()}
	err := handler.CreatePayment(billing, 30000000.00, method)

	require.Error(t, err) // Langsung gagal test kalau err == nil
	assert.Contains(t, err.Error(), "Total payment exceeds billing total")
	// assert.Contains(t, err.Error(), "Total payment exceeds billing total")
}