package handler

import (
	"context"
	"database/sql"
	"fmt"
	"pairproject/entity"
	"pairproject/utils"
)

// OrderDetailHandler bertanggung jawab menangani logika terkait tabel order_details.
type OrderDetailHandler struct {
	DB  *sql.DB           // Koneksi ke database
	Ctx *context.Context  // Context untuk mengambil informasi user yang login
}

// UpdateDetail mengubah jumlah (qty) dari sebuah item order_detail berdasarkan ID.
// Parameter:
// - id: ID dari order_detail yang ingin di-update
// - qty: jumlah baru yang akan diset
// Return:
// - entity.OrderDetail berisi qty yang diupdate
// - error jika terjadi kesalahan
func (od *OrderDetailHandler) UpdateDetail(id int, qty int) (entity.OrderDetail, error) {
	var orderDetail entity.OrderDetail

	// Ambil user dari context, untuk menyimpan informasi siapa yang mengupdate
	user, ok := utils.GetUser(*od.Ctx)
	if !ok {
		// Jika user tidak ditemukan di context, kembalikan error
		return orderDetail, fmt.Errorf("failed to get user from context")
	}

	// Eksekusi perintah update untuk mengubah qty dan updated_by
	res, err := od.DB.Exec(
		"UPDATE order_details SET qty = ?, updated_by = ? WHERE id = ?",
		qty, user.ID, id,
	)
	if err != nil {
		// Jika terjadi error saat eksekusi query update
		return orderDetail, fmt.Errorf("Terjadi kesalahan update data: %s", err)
	}

	// ⚠️ Catatan: `LastInsertId()` hanya digunakan untuk operasi INSERT.
	// Dalam konteks UPDATE, pemanggilan `LastInsertId()` ini tidak valid dan sebaiknya dihapus.
	_, err = res.LastInsertId()
	if err != nil {
		return orderDetail, fmt.Errorf("Terjadi kesalahan mengambil order detail id: %s", err)
	}

	// Kembalikan struct OrderDetail dengan qty baru
	return entity.OrderDetail{Qty: qty}, nil
}