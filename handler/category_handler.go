package handler

import (
	"context"
	"database/sql"
	"fmt"
	"pairproject/entity"
	"pairproject/utils"
)

// CategoryHandler bertanggung jawab untuk menangani operasi terkait kategori.
// Struktur ini menyimpan koneksi ke database melalui field DB.
type CategoryHandler struct {
	DB *sql.DB // koneksi ke database
}

// CreateCategory digunakan untuk membuat kategori baru di sistem.
// Hanya pengguna dengan role "admin" yang diperbolehkan untuk menjalankan fungsi ini.
func (h *CategoryHandler) CreateCategory(ctx *context.Context, category *entity.Category) {
	// Ambil informasi user dari context (hasil login)
	user, ok := utils.GetUser(*ctx)

	// Validasi: hanya admin yang diizinkan membuat kategori
	if !ok || user.Role != "admin" {
		fmt.Println("Unauthorized. Only admin can create categories.")
		return
	}

	// Query SQL untuk menyisipkan data kategori baru ke dalam tabel categories.
	// Field created_by dan updated_by diisi dengan ID user yang sedang login (admin)
	query := `INSERT INTO categories (name, created_by, updated_by) VALUES (?, ?, ?)`

	// Eksekusi query dengan parameter yang diambil dari input dan user yang login
	_, err := h.DB.Exec(query, category.Name, user.ID, user.ID)
	if err != nil {
		// Jika terjadi error saat menyimpan ke database, tampilkan pesan error
		fmt.Printf("Failed to create category: %v\n", err)
		return
	}

	// Jika berhasil
	fmt.Println("Category created successfully.")
}