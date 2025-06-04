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

type BillingHandler struct {
	DB *sql.DB
	Ctx *context.Context
}

func (b *BillingHandler) GenerateBill(o entity.Order) (entity.Billing, error) {
	var billing entity.Billing

	user, ok := utils.GetUser(*b.Ctx)
	if !ok {
		return billing, fmt.Errorf("Please Login.")
	}

	// Tax rate 10%
	 tax := o.Total * 0.10
	// Calculate total with tax
	total := o.Total + tax

	// Generate bill number
	numberDisplay := b.GenerateBillNumber()

	// Issue date dan due date (30 menit setelah issue)
	issueDate := time.Now()
	dueDate := issueDate.Add(30 * time.Minute)

	// Insert ke database
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

	billingID, err := res.LastInsertId()
	if err != nil {
		return billing, errors.New("Terjadi kesalahan saat mengambil id")
	}

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


func (b *BillingHandler) GenerateBillNumber() string {
	currentYearMonth := time.Now().Format("200601") // YYYYMMDD
	var lastNumber int

	// Query for the latest number_display for the current date from billings table
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
	err := b.DB.QueryRow(query, currentYearMonth).Scan(&lastNumber)
	if err != nil {
		lastNumber = 0
	}

	numberDisplay := fmt.Sprintf("BIL-%s-%03d", currentYearMonth, lastNumber+1)
	return numberDisplay
}

func (b *BillingHandler) GetBillByNumberDisplay(numberDisplay string) (entity.Billing, error) {
	var billing entity.Billing
	// user, ok := utils.GetUser(*b.Ctx)
	// if !ok {
	// 	return billing, fmt.Errorf("failed to get user from context")
	// }

	query := `
		SELECT billings.id, order_id, billings.number_display, issue_date, due_date, billings.status, tax, billings.total, billings.created_by
		FROM billings
		JOIN orders on orders.id = billings.order_id 
		WHERE billings.number_display = ?
		LIMIT 1
	`

	err := b.DB.QueryRow(query, numberDisplay).Scan(
	&billing.ID,             // 1: billings.id
	&billing.OrderID,        // 2: order_id
	&billing.NumberDisplay,  // 3: billings.number_display
	&billing.IssueDate,      // 4: issue_date ✅
	&billing.DueDate,        // 5: due_date
	&billing.Status,         // 6: billings.status
	&billing.Tax,            // 7: tax
	&billing.Total,          // 8: billings.total
	&billing.CreatedBy,      // 9: billings.created_by
)

	if err != nil {
		if err == sql.ErrNoRows {
			return billing, errors.New("Billing tidak ditemukan")
		}
		
		return billing, fmt.Errorf("failed to create bill: %s", err)
	}

	return billing, nil
}
