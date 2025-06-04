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

func (p *PaymentHandler) CreatePayment(billing entity.Billing, amount float64, paymentMethod entity.Method) error {
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
	return nil
}