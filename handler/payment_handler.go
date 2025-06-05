package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pairproject/entity"
	"pairproject/utils"
	"time"
)

// PaymentHandler adalah struct yang bertugas menangani logika terkait pembayaran.
// DB digunakan untuk koneksi database dan Ctx menyimpan informasi pengguna yang sedang login.
type PaymentHandler struct {
	DB  *sql.DB
	Ctx *context.Context
}

// CreatePayment membuat entri pembayaran baru untuk suatu tagihan (billing).
// Jika pembayaran melebihi batas yang diizinkan oleh trigger DB, maka akan gagal.
// Setelah pembayaran berhasil, status order dan billing akan diperbarui melalui billingHandler.
func (p *PaymentHandler) CreatePayment(billingHandler *BillingHandler, billing entity.Billing, amount float64, paymentMethod entity.Method) error {
	// Ambil informasi user dari context
	user, ok := utils.GetUser(*p.Ctx)
	if !ok {
		return fmt.Errorf("failed to get user from context")
	}
	
	// Cek apakah sudah melewati batas waktu pembayaran
	if time.Now().After(billing.DueDate) {
		return errors.New("cannot create payment: order is past due date")
	}

	// Query untuk menyimpan data payment baru
	insertQuery := "INSERT INTO payments (billing_id, amount, created_by, method) VALUES (?, ?, ?, ?)"
	_, err := p.DB.Exec(insertQuery, billing.ID, amount, user.ID, string(paymentMethod))
	if err != nil {
		return fmt.Errorf("Gagal membuat payment: %s", err)
	}
	
	// Setelah payment berhasil dibuat, update status billing dan order terkait
	err = billingHandler.UpdateOrderAndBillingStatus(billing.ID)
	if err != nil {
		return fmt.Errorf("Gagal mengupdate order dan billing: %s", err)
	}
		
	return nil // Berhasil
}

// GetPaymentsByBillingID mengambil semua riwayat pembayaran berdasarkan ID tagihan (billing_id)
func (p *PaymentHandler) GetPaymentsByBillingID(billingID int) ([]entity.Payment, error) {
	// Query untuk mengambil semua payment berdasarkan billing_id
	query := `
		SELECT 
			id, billing_id, date, amount, method, created_at, updated_at, created_by, IFNULL(updated_by, 0)
		FROM payments
		WHERE billing_id = ?
		ORDER BY date ASC
	`

	rows, err := p.DB.Query(query, billingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query payments: %w", err)
	}
	defer rows.Close()

	var payments []entity.Payment // Slice untuk menampung hasil pembayaran

	// Iterasi setiap baris hasil query
	for rows.Next() {
		var pmt entity.Payment
		var method string // Sementara tampung method dalam bentuk string

		// Scan data ke dalam struct Payment
		err := rows.Scan(
			&pmt.ID,
			&pmt.BillingID,
			&pmt.Date,
			&pmt.Amount,
			&method, // method disimpan sementara sebagai string
			&pmt.CreatedAt,
			&pmt.UpdatedAt,
			&pmt.CreatedBy,
			&pmt.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment row: %w", err)
		}

		pmt.Method = entity.Method(method) // Konversi method ke enum bertipe Method
		payments = append(payments, pmt)   // Tambahkan ke slice hasil
	}

	// Cek error saat iterasi (bukan error dari query)
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return payments, nil // Kembalikan list pembayaran
}