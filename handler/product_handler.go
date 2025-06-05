package handler

import (
	"context"
	"database/sql"
	"fmt"
	"pairproject/entity"
	"pairproject/utils"
)

// ProductHandler bertanggung jawab untuk menangani operasi terkait produk,
// seperti mengambil daftar produk dan menambahkan produk baru.
// DB digunakan untuk akses database, dan Ctx menyimpan context dengan informasi user.
type ProductHandler struct {
	DB  *sql.DB
	Ctx *context.Context
}

// GetProducts mengambil semua produk dari database beserta kategori mereka.
// Produk dikembalikan dalam bentuk slice []entity.Product.
func (p *ProductHandler) GetProducts() ([]entity.Product, error) {
	// Query SQL untuk mengambil data produk dengan join ke tabel kategori untuk mendapatkan nama kategori
	rows, err := p.DB.Query(`
		SELECT p.id, p.name, p.stock, p.description, c.name as category_name, p.price
		FROM products p
		JOIN categories c ON p.category_id = c.id
	`)
	if err != nil {
		// Jika terjadi error saat query, langsung return error tersebut
		return nil, err
	}
	defer rows.Close() // Pastikan rows ditutup saat fungsi selesai agar tidak terjadi memory leak

	var products []entity.Product // Slice untuk menampung data produk hasil query

	// Iterasi setiap baris hasil query
	for rows.Next() {
		var p entity.Product    // Variabel untuk menampung data produk per baris
		var c entity.Category   // Variabel untuk menampung nama kategori sementara

		// Scan data dari baris ke struct produk dan category nama
		// Field category_name dari query akan kita simpan ke c.Name sementara (tidak dipakai dalam products sekarang)
		if err := rows.Scan(&p.ID, &p.Name, &p.Stock, &p.Description, &c.Name, &p.Price); err != nil {
			return nil, err // Jika error saat scanning data, kembalikan error
		}

		// Tambahkan produk yang sudah di-scan ke dalam slice products
		products = append(products, p)
	}

	// Kembalikan slice produk dan nil error (berhasil)
	return products, nil
}

// CreateProduct menambahkan produk baru ke database.
// Memerlukan data produk dari parameter dan user harus sudah login (diperiksa lewat context).
func (p *ProductHandler) CreateProduct(product entity.Product) error {
	// Ambil data user dari context, user harus sudah login untuk dapat membuat produk
	user, ok := utils.GetUser(*p.Ctx)
	if !ok {
		return fmt.Errorf("Please Login!") // Jika user tidak ditemukan di context, return error
	}

	// Query SQL untuk memasukkan data produk baru ke tabel products
	query := `
		INSERT INTO products (name, stock, description, category_id, price, created_by)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	// Jalankan query dengan parameter dari input product dan ID user sebagai created_by
	_, err := p.DB.Exec(query, product.Name, product.Stock, product.Description, product.CategoryID, product.Price, user.ID)
	if err != nil {
		// Jika terjadi error saat eksekusi query insert, return error dengan pesan generik
		return fmt.Errorf("Terjadi kesalahan ketika membuat produk")
	}

	// Jika tidak ada error, kembalikan nil (berhasil)
	return nil
}