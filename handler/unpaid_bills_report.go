package handler

import (
	"context"
	"database/sql"
	"fmt"
)

// ReportHandlerUnpaidBill bertugas menangani laporan tagihan yang belum dibayar.
// Struct ini memiliki properti DB yang merupakan koneksi ke database.
type ReportHandlerUnpaidBill struct {
	DB *sql.DB
}

// UnpaidBill merepresentasikan struktur data dari tagihan yang belum dibayar,
// termasuk informasi bill, order, customer, dan detail nilai tagihan.
type UnpaidBill struct {
	ID           int     // ID tagihan
	BillNumber   string  // Nomor tagihan yang ditampilkan
	OrderNumber  string  // Nomor pesanan terkait
	CustomerName string  // Nama customer
	Tax          float64 // Pajak yang dikenakan
	Total        float64 // Total tagihan
	Status       string  // Status tagihan (harus 'unpaid')
	CreatedAt    string  // Tanggal pembuatan tagihan
}

// GetUnpaidBills mengambil daftar tagihan yang belum dibayar dari database.
// Mengembalikan slice dari UnpaidBill dan error jika terjadi kesalahan.
func (h *ReportHandlerUnpaidBill) GetUnpaidBills(ctx *context.Context) ([]UnpaidBill, error) {
	// Query SQL untuk mendapatkan data tagihan dengan status "unpaid"
	query := `
		SELECT 
			b.id,
			b.number_display,
			o.number_display AS order_number,
			c.name AS customer_name,
			b.tax,
			b.total,
			b.status,
			b.created_at
		FROM billings b
		JOIN orders o ON b.order_id = o.id
		JOIN customers c ON o.customer_id = c.id
		WHERE b.status = 'unpaid'
		ORDER BY b.created_at DESC
	`

	// Menjalankan query dengan konteks
	rows, err := h.DB.QueryContext(*ctx, query)
	if err != nil {
		// Jika gagal melakukan query, kembalikan error
		return nil, fmt.Errorf("failed to query unpaid bills: %w", err)
	}
	defer rows.Close()

	var bills []UnpaidBill

	// Melakukan iterasi terhadap setiap hasil query
	for rows.Next() {
		var bill UnpaidBill
		// Scan baris hasil query ke dalam struct UnpaidBill
		if err := rows.Scan(
			&bill.ID,
			&bill.BillNumber,
			&bill.OrderNumber,
			&bill.CustomerName,
			&bill.Tax,
			&bill.Total,
			&bill.Status,
			&bill.CreatedAt,
		); err != nil {
			// Jika gagal saat scanning, kembalikan error
			return nil, fmt.Errorf("failed to scan bill: %w", err)
		}
		// Menambahkan data ke slice hasil
		bills = append(bills, bill)
	}

	// Mengembalikan hasil berupa slice UnpaidBill dan nil error
	return bills, nil
}