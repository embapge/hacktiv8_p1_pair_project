package handler

import (
	"context"
	"database/sql"
	"fmt"
)

type ReportHandler struct {
	DB *sql.DB
	Ctx *context.Context
}

type MostSoldProduct struct {
	ProductID int    // ID produk
	Name      string // Nama produk
	TotalSold int    // Total jumlah produk terjual
}

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

type RevenueDetail struct {
	BillNumber    string  // Nomor tagihan (billing)
	PaymentDate   string  // Tanggal pembayaran
	Amount        float64 // Jumlah pembayaran
	Method        string  // Metode pembayaran (misal: cash, credit card)
	CustomerName  string  // Nama pelanggan
	OrderNumber   string  // Nomor pesanan/order
}

func (r *ReportHandler) GetMostSoldProducts() ([]MostSoldProduct, error) {
query := `
	SELECT
			p.id, p.name, SUM(od.qty) AS total_sold
	FROM
			products p
	JOIN
			order_details od ON p.id = od.product_id
	GROUP BY 
			p.id, p.name
	ORDER BY
			total_sold DESC
		LIMIT 5
	`

       // Eksekusi query ke database
       rows, err := r.DB.Query(query)
       if err != nil {
               // Jika query gagal, kembalikan error
               return nil, err
       }
       defer rows.Close() // Pastikan rows ditutup setelah selesai

       // Slice untuk menampung hasil produk terlaris
       var results []MostSoldProduct

       // Iterasi setiap baris hasil query
       for rows.Next() {
               var p MostSoldProduct
               // Scan data dari baris ke struct MostSoldProduct
               if err := rows.Scan(&p.ProductID, &p.Name, &p.TotalSold); err != nil {
                       // Jika error saat scan data, lanjutkan ke baris berikutnya tanpa menghentikan seluruh proses
                       continue
               }
               // Tambahkan produk hasil scan ke slice results
               results = append(results, p)
       }

       // Kembalikan slice hasil dan nil error
       return results, nil
}

func (h *ReportHandler) GetUnpaidBills() ([]UnpaidBill, error) {
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
	rows, err := h.DB.Query(query)
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
		); 
		
		err != nil {
				// Jika gagal saat scanning, kembalikan error
				return nil, fmt.Errorf("failed to scan bill: %w", err)
		}

		// Menambahkan data ke slice hasil
		bills = append(bills, bill)
	}

       // Mengembalikan hasil berupa slice UnpaidBill dan nil error
	return bills, nil
}

func (r *ReportHandler) GetRevenueDetails() ([]RevenueDetail, error) {
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
	rows, err := r.DB.Query(query)
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