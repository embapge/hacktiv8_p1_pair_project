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

// BillingHandler bertanggung jawab menangani proses terkait tagihan
type BillingHandler struct {
	DB  *sql.DB
	Ctx *context.Context
}

// BillingWithPaymentsSimple digunakan untuk mengambil informasi billing dan daftar pembayaran terkait
type BillingWithPaymentsSimple struct {
	BillingID     int
	OrderID       int
	NumberDisplay string
	Tax           float64
	Total         float64
	Status        string
	Payments      []struct {
		ID     int
		Amount float64
	}
}

// GenerateBill membuat tagihan berdasarkan informasi order
func (b *BillingHandler) GenerateBill(o entity.Order) (entity.Billing, error) {
	var billing entity.Billing

	// Ambil data user dari context, validasi login
	user, ok := utils.GetUser(*b.Ctx)
	if !ok {
		return billing, fmt.Errorf("Please Login.")
	}

	// Hitung pajak 10%
	tax := o.Total * 0.10
	// Hitung total setelah pajak
	total := o.Total + tax
	// Buat nomor tagihan
	numberDisplay := b.GenerateBillNumber()

	// Tanggal issue dan due (30 menit ke depan)
	issueDate := time.Now()
	dueDate := issueDate.Add(30 * time.Minute)

	// Insert tagihan ke DB
	insertQuery := `
		INSERT INTO billings (order_id, tax, total, number_display, issue_date, due_date, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	res, err := b.DB.Exec(
		insertQuery,
		o.ID,
		tax,
		total,
		numberDisplay,
		issueDate,
		dueDate,
		user.ID,
	)
	if err != nil {
		return billing, errors.New("Kesalahan membuat tagihan")
	}

	// Ambil ID dari billing yang baru dibuat
	billingID, err := res.LastInsertId()
	if err != nil {
		return billing, errors.New("Terjadi kesalahan saat mengambil id")
	}

	// Kembalikan data billing
	billing = entity.Billing{
		ID:            int(billingID),
		OrderID:       o.ID,
		Tax:           tax,
		Total:         total,
		NumberDisplay: numberDisplay,
		IssueDate:     issueDate,
		DueDate:       dueDate,
		CreatedBy:     user.ID,
	}

	return billing, nil
}

// GenerateBillNumber membuat nomor tagihan otomatis berdasarkan bulan dan urutan terakhir
func (b *BillingHandler) GenerateBillNumber() string {
	currentYearMonth := time.Now().Format("200601") // Contoh: "202506"
	var lastNumber int

	query := `
		SELECT 
			COALESCE(
				CAST(SUBSTR(number_display, 14, 3) AS UNSIGNED),
				0
			) AS last_number
		FROM billings
		WHERE SUBSTR(number_display, 5, 6) = ?
		ORDER BY last_number DESC
		LIMIT 1
	`
	// Ambil nomor urut terakhir dari bulan yang sama
	err := b.DB.QueryRow(query, currentYearMonth).Scan(&lastNumber)
	if err != nil {
		lastNumber = 0
	}

	// Format: BIL-YYYYMM-XXX
	numberDisplay := fmt.Sprintf("BIL-%s-%03d", currentYearMonth, lastNumber+1)
	return numberDisplay
}

// GetBillByNumberDisplay mencari tagihan berdasarkan nomor tagihan untuk customer yang sedang login
func (b *BillingHandler) GetBillByNumberDisplay(numberDisplay string) (entity.Billing, error) {
	var billing entity.Billing
	user, ok := utils.GetUser(*b.Ctx)
	if !ok {
		return billing, fmt.Errorf("Please Login.")
	}

	query := `
		SELECT billings.id, order_id, billings.number_display, issue_date, due_date, billings.status, tax, billings.total, billings.created_by
		FROM billings
		JOIN orders on orders.id = billings.order_id
		WHERE billings.number_display = ? AND orders.customer_id = ? AND (billings.status = 'unpaid' OR billings.status = 'lesspaid')
		LIMIT 1
	`

	// Cari billing untuk customer saat ini dan status unpaid/lesspaid
	err := b.DB.QueryRow(query, numberDisplay, user.Customer.ID).Scan(
		&billing.ID,
		&billing.OrderID,
		&billing.NumberDisplay,
		&billing.IssueDate,
		&billing.DueDate,
		&billing.Status,
		&billing.Tax,
		&billing.Total,
		&billing.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return billing, errors.New("Billing tidak ditemukan")
		}
		return billing, fmt.Errorf("Terjadi kesalahan: %s", err)
	}

	return billing, nil
}

// GetBillingWithSimplePayments mengambil data billing dan semua pembayaran terkait (jika ada)
func (b *BillingHandler) GetBillingWithSimplePayments(billingID int) (BillingWithPaymentsSimple, error) {
	var result BillingWithPaymentsSimple

	query := `
		SELECT 
			billings.id, orders.id, billings.number_display, billings.tax, billings.total, billings.status, 
			payments.id, payments.amount
		FROM billings
		JOIN orders on orders.id = billings.order_id
		LEFT JOIN payments ON payments.billing_id = billings.id
		WHERE billings.id = ?
	`

	rows, err := b.DB.Query(query, billingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, err
		}
		return result, err
	}
	defer rows.Close()

	// Iterasi hasil dan isi struct result
	for rows.Next() {
		var (
			billingID     int
			orderID       int
			numberDisplay string
			tax           float64
			total         float64
			status        string
			paymentID     sql.NullInt64
			paymentAmount sql.NullFloat64
		)

		err := rows.Scan(&billingID, &orderID, &numberDisplay, &tax, &total, &status, &paymentID, &paymentAmount)
		if err != nil {
			return result, err
		}

		result.BillingID = billingID
		result.OrderID = orderID
		result.NumberDisplay = numberDisplay
		result.Tax = tax
		result.Total = total
		result.Status = status

		// Tambahkan payment jika tersedia
		result.Payments = append(result.Payments, struct {
			ID     int
			Amount float64
		}{
			ID:     int(paymentID.Int64),
			Amount: paymentAmount.Float64,
		})
	}

	return result, nil
}

// UpdateOrderAndBillingStatus memperbarui status billing dan order berdasarkan total pembayaran
func (b *BillingHandler) UpdateOrderAndBillingStatus(billingID int) error {
	// Mulai transaksi
	tx, err := b.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ambil billing dan semua payment
	billPayments, err := b.GetBillingWithSimplePayments(billingID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Hitung total pembayaran
	var total float64
	for _, payment := range billPayments.Payments {
		total += payment.Amount
	}

	// Update status billing dan order sesuai pembayaran
	if total >= billPayments.Total {
		_, err = tx.Exec("UPDATE billings SET status = 'paid' WHERE id = ?", billingID)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec("UPDATE orders SET status = 'completed' WHERE id = ?", billPayments.OrderID)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		_, err = tx.Exec("UPDATE billings SET status = 'lesspaid' WHERE id = ?", billingID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaksi
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Terjadi kesalahan saat commit transaksi: %v", err)
	}

	return nil
}