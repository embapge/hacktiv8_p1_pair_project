package handler

import (
	"context"
	"database/sql"
	"fmt"
	"pairproject/entity"
	"pairproject/utils"
)

// ProductHandler adalah handler yang bertugas menangani proses terkait produk.
// Struct ini memiliki field DB untuk koneksi ke database.
type ProductHandler struct {
	DB *sql.DB
}

// CreateProduct digunakan untuk menambahkan produk baru ke dalam database.
// Hanya user dengan role "admin" yang diizinkan menjalankan fungsi ini.
func (h *ProductHandler) CreateProduct(ctx *context.Context, product *entity.Product) {
	// Mengambil user dari context untuk mengecek otorisasi
	user, ok := utils.GetUser(*ctx)
	if !ok || user.Role != "admin" {
		// Jika user tidak ditemukan atau bukan admin, tampilkan pesan tidak diizinkan
		fmt.Println("Unauthorized. Only admin can create products.")
		return
	}

	// Query SQL untuk menyisipkan data produk baru ke tabel 'products'
	query := `
		INSERT INTO products (name, stock, description, category_id, price, created_by, updated_by)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// Menjalankan query dengan parameter dari input produk
	_, err := h.DB.Exec(query, product.Name, product.Stock, product.Description, product.CategoryID, product.Price, user.ID, user.ID)
	if err != nil {
		// Jika terjadi error saat eksekusi query, tampilkan pesan kesalahan
		fmt.Printf("Failed to create product: %v\n", err)
		return
	}

	// Jika berhasil, tampilkan pesan sukses
	fmt.Println("Product created successfully.")
}