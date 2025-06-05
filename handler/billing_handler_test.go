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
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite" // SQLite driver tanpa CGO (cocok untuk testing)
)

// SetupBillingAndOrdersDB membuat database SQLite in-memory untuk pengujian billing dan orders
func SetupBillingAndOrdersDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err, "Gagal membuka database in-memory")

	// Ambil format tahun-bulan saat ini untuk penomoran tagihan
	currentYearMonth := time.Now().Format("200601") // contoh: "202506"

	// Buat schema dan seed data untuk orders dan billings
	schema := fmt.Sprintf(`
		CREATE TABLE orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			number_display TEXT,
			customer_id INTEGER,
			date DATETIME DEFAULT CURRENT_TIMESTAMP,
			total NUMERIC DEFAULT 0 CHECK (total >= 0) NOT NULL, 
			created_by INTEGER
		);

		-- Tambah 1 data order sebagai dasar tagihan
		INSERT INTO orders (number_display, customer_id, created_by, total)
		VALUES ("ORD-%s-001", 1, 1, 1000000);

		CREATE TABLE billings ( 
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			order_id INTEGER NOT NULL, 
			number_display TEXT, 
			issue_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			due_date TIMESTAMP,
			tax NUMERIC DEFAULT 0 CHECK (tax >= 0) NOT NULL, 
			total NUMERIC DEFAULT 0 CHECK (total >= 0) NOT NULL, 
			status TEXT DEFAULT 'unpaid' CHECK (status IN ('unpaid', 'paid', 'cancelled', 'refunded')), 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER NOT NULL, 
			updated_by INTEGER
		);

		-- Tambahkan 2 tagihan awal agar test GenerateBillNumber menghasilkan BIL-XXX-003
		INSERT INTO billings (order_id, number_display, tax, total, status, created_by, updated_by)
		VALUES 
			(1, 'BIL-%s-001', 50000.00, 1700000.00, 'paid', 2, 2),
			(1, 'BIL-%s-002', 30000.00, 630000.00, 'unpaid', 3, 3);
	`, currentYearMonth, currentYearMonth, currentYearMonth)

	_, err = db.Exec(schema)
	require.NoError(t, err, "Gagal membuat schema dan seeding data")

	return db
}

// TestGenerateBill menguji fungsi GenerateBill() untuk menghitung tagihan dari order
func TestGenerateBill(t *testing.T) {
	db := SetupBillingAndOrdersDB(t)

	// Buat context user palsu (dari utils) agar billing punya created_by yang valid
	ctx := utils.NewTestContextWithUser()
	handler := &BillingHandler{DB: db, Ctx: &ctx}

	// Buat order baru dengan total 350.000
	order := entity.Order{ID: 1, Total: 350000}

	// Jalankan GenerateBill()
	billing, err := handler.GenerateBill(order)
	require.NoError(t, err)

	// Pastikan perhitungan pajak dan total tagihan sesuai (10% PPN)
	assert.Equal(t, 35000.0, billing.Tax)
	assert.Equal(t, 385000.0, billing.Total)
	assert.NotZero(t, billing.ID) // Billing berhasil dibuat
}

// TestGenerateBillNumber menguji nomor tagihan yang di-generate otomatis
func TestGenerateBillNumber(t *testing.T) {
	db := SetupBillingAndOrdersDB(t)

	// Tidak butuh user context hanya untuk generate nomor
	ctx := context.Background()
	handler := &BillingHandler{DB: db, Ctx: &ctx}

	// Jalankan generate nomor tagihan
	num := handler.GenerateBillNumber()

	// Seharusnya menjadi BIL-<YYYYMM>-003 karena sudah ada 2 tagihan sebelumnya
	expected := fmt.Sprintf("BIL-%s-003", time.Now().Format("200601"))
	assert.Equal(t, expected, num)
}