package handler

import (
	"database/sql"
)

// ReportHandlerMostSold bertanggung jawab untuk menghasilkan laporan produk yang paling banyak terjual.
type ReportHandlerMostSold struct {
	DB *sql.DB // Koneksi ke database
}

// MostSoldProduct merepresentasikan data produk yang paling banyak terjual.
type MostSoldProduct struct {
	ProductID int    // ID produk
	Name      string // Nama produk
	TotalSold int    // Total jumlah produk terjual
}

// GetMostSoldProducts mengambil 5 produk dengan penjualan tertinggi dari database.
// Mengembalikan slice produk dan error jika terjadi kegagalan.
func (h *ReportHandlerMostSold) GetMostSoldProducts() ([]MostSoldProduct, error) {
	// Query SQL untuk menghitung total penjualan tiap produk
	// Menggabungkan tabel products dan order_details berdasarkan product_id
	// Mengelompokkan hasil berdasarkan id dan nama produk
	// Mengurutkan berdasarkan total penjualan dari yang terbesar
	// Membatasi hasil hanya 5 produk teratas
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
	rows, err := h.DB.Query(query)
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