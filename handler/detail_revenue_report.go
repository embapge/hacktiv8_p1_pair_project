package handler

import (
	"context"
	"database/sql"
)

// RevenueDetail merepresentasikan data detail pendapatan yang akan ditampilkan pada laporan revenue.
type RevenueDetail struct {
	BillNumber    string  // Nomor tagihan (billing)
	PaymentDate   string  // Tanggal pembayaran
	Amount        float64 // Jumlah pembayaran
	Method        string  // Metode pembayaran (misal: cash, credit card)
	CustomerName  string  // Nama pelanggan
	OrderNumber   string  // Nomor pesanan/order
}

// RevenueHandler bertanggung jawab menangani operasi terkait laporan revenue.
type RevenueHandler struct {
	DB *sql.DB // Koneksi database
}

// GetRevenueDetails mengambil data detail pendapatan dari database.
// Data ini merupakan hasil join beberapa tabel seperti payments, billings, orders, customers, dll.
func (r *RevenueHandler) GetRevenueDetails(ctx *context.Context) ([]RevenueDetail, error) {
	// Query SQL mengambil data pembayaran dengan informasi terkait tagihan, pelanggan, dan pesanan.
	// Hasil diurutkan berdasarkan tanggal pembayaran secara menurun (terbaru di atas).
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

	// Eksekusi query dengan context agar dapat dibatalkan jika perlu (misal timeout)
	rows, err := r.DB.QueryContext(*ctx, query)
	if err != nil {
		// Jika error saat eksekusi query, kembalikan error
		return nil, err
	}
	defer rows.Close() // Pastikan rows ditutup saat fungsi selesai

	// Slice untuk menampung hasil query
	var results []RevenueDetail

	// Iterasi tiap baris hasil query
	for rows.Next() {
		var detail RevenueDetail
		// Scan data baris ke struct detail
		if err := rows.Scan(
			&detail.BillNumber,
			&detail.PaymentDate,
			&detail.Amount,
			&detail.Method,
			&detail.CustomerName,
			&detail.OrderNumber,
		); err != nil {
			// Jika error saat scanning data, kembalikan error
			return nil, err
		}
		// Tambahkan detail ke hasil
		results = append(results, detail)
	}

	// Cek error dari iterasi rows (misal error jaringan saat fetch data)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Kembalikan slice hasil dan nil error
	return results, nil
}