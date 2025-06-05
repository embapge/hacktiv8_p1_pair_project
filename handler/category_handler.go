package handler

import (
	"context"
	"database/sql"
	"fmt"
	"pairproject/entity"
	"pairproject/utils"
)

// CategoryHandler adalah struct yang menyimpan dependensi database dan context.
// Struct ini digunakan untuk menangani semua logika bisnis terkait kategori produk.
type CategoryHandler struct {
	DB  *sql.DB           // koneksi database
	Ctx *context.Context  // context untuk mengambil user yang sedang login
}

// CreateCategory adalah fungsi untuk menambahkan kategori baru ke dalam database.
// Parameter: name (string) - nama kategori yang ingin dibuat.
func (c *CategoryHandler) CreateCategory(name string) error {
	// Ambil informasi user yang sedang login dari context
	user, ok := utils.GetUser(*c.Ctx)
	if !ok {
		// Jika user tidak ditemukan (belum login), kembalikan error
		return fmt.Errorf("Please Login!")
	}

	// Query SQL untuk menambahkan kategori baru
	query := `INSERT INTO categories (name, created_by) VALUES (?, ?)`
	
	// Eksekusi query dengan parameter nama kategori dan ID user yang membuat
	_, err := c.DB.Exec(query, name, user.ID)
	if err != nil {
		// Jika gagal menyimpan ke database, kembalikan error
		return fmt.Errorf("Terjadi kesalahan ketika membuat kategori")
	}

	// Jika berhasil, kembalikan nil (tidak ada error)
	return nil
}

// GetCategories mengambil semua kategori dari database dan mengembalikannya dalam bentuk slice.
func (c *CategoryHandler) GetCategories() ([]entity.Category, error) {
	// Jalankan query SQL untuk mengambil semua kategori
	rows, err := c.DB.Query("SELECT id, name FROM categories")
	var emptyCategory []entity.Category
	if err != nil {
		// Jika gagal query, kembalikan slice kosong dan error
		return emptyCategory, err
	}
	// Pastikan rows ditutup setelah selesai digunakan
	defer rows.Close()

	// Slice untuk menyimpan hasil kategori
	var categories []entity.Category

	// Looping setiap baris hasil query
	for rows.Next() {
		var c entity.Category
		// Scan data hasil query ke dalam struct Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			// Jika gagal scan, kembalikan slice kosong dan error
			return emptyCategory, err
		}
		// Tambahkan kategori ke dalam slice
		categories = append(categories, c)
	}

	// Cek jika terjadi error saat membaca rows
	if err := rows.Err(); err != nil {
		return emptyCategory, err
	}

	// Kembalikan data kategori
	return categories, nil
}