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

	schema := `
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL, -- Changed from DATE to TEXT
    status TEXT NOT NULL CHECK (status IN ('processing', 'completed', 'cancel')),
    total NUMERIC NOT NULL DEFAULT 0 CHECK (total >= 0), -- Changed from DECIMAL to NUMERIC
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Handled by trigger below
	created_by INTEGER NOT NULL
);

INSERT INTO orders (id, date, status, created_by, total)
VALUES (1, '2025-06-05', 'processing', 1, 10000000.00);

CREATE TRIGGER update_orders_updated_at
AFTER UPDATE ON orders
FOR EACH ROW
BEGIN
    UPDATE orders SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

---

CREATE TABLE billings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	number_display TEXT NOT NULL,
    order_id INTEGER NOT NULL,
    due_date DATETIME,
	tax NUMERIC DEFAULT 0 CHECK (total >= 0) NOT NULL,
	status TEXT NOT NULL CHECK (status IN ('unpaid', 'lesspaid', 'paid', 'cancelled', 'refunded')) DEFAULT 'unpaid',
    total NUMERIC DEFAULT 0 CHECK (total >= 0) NOT NULL, -- Changed from DECIMAL to NUMERIC
    amount NUMERIC DEFAULT 0 CHECK (amount >= 0) NOT NULL, -- Changed from DECIMAL to NUMERIC
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Handled by trigger below
    created_by INTEGER NOT NULL,
	updated_by INTEGER NOT NULL
);

-- Contoh insert data (Added 'total' column and value to match table schema)
INSERT INTO billings (order_id, number_display, tax, total, status, created_by, updated_by)
VALUES 
(1, 'BIL-202506-001', 50000.00, 1700000.00, 'paid', 2, 2);

CREATE TABLE payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    billing_id INTEGER NOT NULL,
    date DATETIME DEFAULT CURRENT_TIMESTAMP,
    amount NUMERIC DEFAULT 0 CHECK (amount >= 0) NOT NULL, -- Changed from DECIMAL to NUMERIC
    method TEXT NOT NULL CHECK (method IN ('credit_card', 'va', 'transfer')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Added updated_at column
    FOREIGN KEY (billing_id) REFERENCES billings(id) ON DELETE CASCADE
);

-- Trigger validate_payment_amount (already compatible with SQLite)
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
func TestCreatePayment_Error(t *testing.T){
	db := SetupTestPaymentDB(t)
	ctx := utils.NewTestContextWithUser()
	handler := &PaymentHandler{DB: db, Ctx: &ctx}
	var method entity.Method = "va"
	billing := entity.Billing{ID: 1, DueDate: time.Now()}
	billingHandler := &BillingHandler{DB: db, Ctx: &ctx}
	err := handler.CreatePayment(billingHandler, billing, 30000000.00, method)

	require.Error(t, err) // Langsung gagal test kalau err == nil
	assert.Contains(t, err.Error(), "Total payment exceeds billing total")
	// assert.Contains(t, err.Error(), "Total payment exceeds billing total")
}