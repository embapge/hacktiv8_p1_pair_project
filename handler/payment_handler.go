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

type PaymentHandler struct {
	DB  *sql.DB
	Ctx *context.Context
}

func (p *PaymentHandler) CreatePayment(billingHandler *BillingHandler, billing entity.Billing, amount float64, paymentMethod entity.Method) error {
	user, ok := utils.GetUser(*p.Ctx)
	if !ok {
		return fmt.Errorf("failed to get user from context")
	}
	
	if time.Now().After(billing.DueDate) {
		return errors.New("cannot create payment: order is past due date")
	}

	insertQuery := "INSERT INTO payments (billing_id, amount, created_by, method) VALUES (?, ?, ?, ?)"
	_, err := p.DB.Exec(insertQuery, billing.ID, amount, user.ID, string(paymentMethod))
	if err != nil {
		return fmt.Errorf("Gagal membuat payment: %s", err)
	}
	
	err = billingHandler.UpdateOrderAndBillingStatus(billing.ID)
	if err != nil{
		return fmt.Errorf("Gagal mengupdate order dan billing: %s", err)
	}
		
	return nil
}

func (p *PaymentHandler) GetPaymentsByBillingID(billingID int) ([]entity.Payment, error) {
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

	var payments []entity.Payment
	for rows.Next() {
		var pmt entity.Payment
		var method string

		err := rows.Scan(
			&pmt.ID,
			&pmt.BillingID,
			&pmt.Date,
			&pmt.Amount,
			&method,
			&pmt.CreatedAt,
			&pmt.UpdatedAt,
			&pmt.CreatedBy,
			&pmt.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment row: %w", err)
		}

		pmt.Method = entity.Method(method)
		payments = append(payments, pmt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return payments, nil
}