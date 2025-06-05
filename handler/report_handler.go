package handler

import (
	"context"
	"database/sql"
	"fmt"
)

// ReportHandler bertugas untuk mengambil laporan terkait produk, tagihan, dan pendapatan dari database.
type ReportHandler struct {
	DB  *sql.DB
	Ctx *context.Context
}

// MostSoldProduct merepresentasikan produk dengan jumlah penjualan terbanyak.
type MostSoldProduct struct {
	ProductID int    // ID produk
	Name      string // Nama produk
	TotalSold int    // Total kuantitas produk terjual
}

// UnpaidBill merepresentasikan tagihan yang belum dibayar.
type UnpaidBill struct {
	ID           int     // ID tagihan
	BillNumber   string  // Nomor tagihan yang ditampilkan
	OrderNumber  string  // Nomor pesanan yang terkait dengan tagihan
	CustomerName string  // Nama pelanggan yang berhubungan dengan tagihan
	Tax          float64 // Pajak yang dikenakan pada tagihan
	Total        float64 // Total jumlah tagihan
	Status       string  // Status tagihan, harus 'unpaid'
	CreatedAt    string  // Waktu pembuatan tagihan
}

// RevenueDetail merepresentasikan detail pembayaran/tagihan untuk laporan pendapatan.
type RevenueDetail struct {
	BillNumber   string  // Nomor tagihan
	PaymentDate  string  // Tanggal pembayaran dilakukan
	Amount       float64 // Jumlah pembayaran
	Method       string  // Metode pembayaran (contoh: cash, credit card, transfer)
	CustomerName string  // Nama pelanggan yang membayar
	OrderNumber  string  // Nomor pesanan terkait pembayaran
}

// GetMostSoldProducts mengambil 5 produk dengan total penjualan terbanyak dari database.
func (r *ReportHandler) GetMostSoldProducts() ([]MostSoldProduct, error) {
	query := `
	SELECT
		p.id, p.name, SUM(od.qty) AS total_sold
	FROM products p
	JOIN order_details od ON p.id = od.product_id
	GROUP BY p.id, p.name
	ORDER BY total_sold DESC
	LIMIT 5
	`

	// Eksekusi query ke database untuk mendapatkan data produk terlaris
	rows, err := r.DB.Query(query)
	if err != nil {
		// Jika query gagal, langsung kembalikan error
		return nil, err
	}
	defer rows.Close() // Pastikan resource rows dibebaskan setelah selesai digunakan

	var results []MostSoldProduct

	// Iterasi tiap baris hasil query
	for rows.Next() {
		var p MostSoldProduct
		// Scan data kolom ke struct MostSoldProduct
		if err := rows.Scan(&p.ProductID, &p.Name, &p.TotalSold); err != nil {
			// Jika error saat scan, skip baris ini dan lanjut ke baris berikutnya tanpa menghentikan proses keseluruhan
			continue
		}
		results = append(results, p) // Tambahkan ke slice hasil
	}

	// Kembalikan slice produk dan nil error jika sukses
	return results, nil
}

// GetUnpaidBills mengambil daftar tagihan dengan status 'unpaid' beserta informasi terkait pesanan dan pelanggan.
func (h *ReportHandler) GetUnpaidBills() ([]UnpaidBill, error) {
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

	// Jalankan query ke database
	rows, err := h.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query unpaid bills: %w", err)
	}
	defer rows.Close()

	var bills []UnpaidBill

	// Iterasi hasil query
	for rows.Next() {
		var bill UnpaidBill
		// Scan hasil baris ke struct UnpaidBill
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
			// Jika scan gagal, return error untuk memudahkan debugging
			return nil, fmt.Errorf("failed to scan bill: %w", err)
		}
		bills = append(bills, bill) // Tambahkan ke slice hasil
	}

	// Kembalikan daftar tagihan unpaid dan nil error jika sukses
	return bills, nil
}

// GetRevenueDetails mengambil data detail pendapatan berupa pembayaran beserta info terkait pelanggan dan pesanan.
func (r *ReportHandler) GetRevenueDetails() ([]RevenueDetail, error) {
	query := `
	SELECT
		b.number_display AS bill_number,
		p.date AS payment_date,
		p.amount,
		p.method,
		c.name AS customer_name,
		o.number_display AS order_number
	FROM payments p
	JOIN billings b ON p.billing_id = b.id
	JOIN orders o ON b.order_id = o.id
	JOIN user_customers uc ON o.customer_id = uc.customer_id
	JOIN customers c ON uc.customer_id = c.id
	ORDER BY p.date DESC
	`

	// Eksekusi query untuk mengambil data pembayaran
	rows, err := r.DB.Query(query)
	if err != nil {
		// Jika query gagal, return error
		return nil, err
	}
	defer rows.Close()

	var results []RevenueDetail

	// Iterasi hasil query
	for rows.Next() {
		var detail RevenueDetail
		// Scan data baris ke struct RevenueDetail
		if err := rows.Scan(
			&detail.BillNumber,
			&detail.PaymentDate,
			&detail.Amount,
			&detail.Method,
			&detail.CustomerName,
			&detail.OrderNumber,
		); err != nil {
			// Jika scan gagal, return error
			return nil, err
		}
		results = append(results, detail) // Tambahkan ke slice hasil
	}

	// Cek apakah ada error dari iterasi rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Kembalikan data revenue dan nil error jika sukses
	return results, nil
}