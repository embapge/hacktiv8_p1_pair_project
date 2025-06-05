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

type OrderHandler struct {
	DB *sql.DB
	Ctx *context.Context
}

func (o *OrderHandler) CreateOrder(oProducts []entity.OrderProduct) (entity.Order, error) {
	var order entity.Order

	user, ok := utils.GetUser(*o.Ctx)
	if !ok {
		return order, fmt.Errorf("Please Login!")
	}

	tx, err := o.DB.Begin()
	if err != nil {
		tx.Rollback()
		return order, err
	}

	numberDisplay := o.GenerateOrderNumber(tx)

	// Insert into orders
	createdDate := time.Now().Format("2006-01-02")
	res, err := tx.Exec("INSERT INTO orders (number_display, customer_id, date, created_by) VALUES (?, ?, ?, ?)", numberDisplay, user.Customer.ID, createdDate, user.ID)
	if err != nil {
		tx.Rollback()
		return order, errors.New("Terjadi kesalahan membuat order")
	}

	orderID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return order, err
	}

	stmt, err := tx.Prepare("INSERT INTO order_details (order_id, product_id, qty, created_by) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return order, errors.New("Terjadi kesalahan membuat detail order")
	}
	defer stmt.Close()

	var orderDetails []entity.OrderDetail
	for _, op := range oProducts {
		_, err := stmt.Exec(orderID, op.ProductId, op.Qty, user.ID)
		if err != nil {
			tx.Rollback()
			return order, errors.New("Terjadi kesalahan membuat detail order")
		}

		orderDetails = append(orderDetails, entity.OrderDetail{
			OrderID:   int(orderID),
			ProductID: op.ProductId,
			Qty:       op.Qty,
			CreatedBy: user.ID,
		})
	}

	order = entity.Order{
		ID:            int(orderID),
		NumberDisplay: numberDisplay,
		CreatedBy:        user.ID,
		Details:  orderDetails,
	}

	err = tx.Commit()
	if err != nil {
		return order, fmt.Errorf("Terjadi kesalahan saat commit transaksi: %v", err)
	}

	return order, nil
}

func (o *OrderHandler) GenerateOrderNumber(tx *sql.Tx) (string) {
	currentYearMonth := time.Now().Format("200601") // YYYYMMDD
	var lastNumber int

	// Query for the latest number_display for the current date
	query := `
		SELECT 
			COALESCE(
				CAST(SUBSTR(number_display, 14, 3) AS UNSIGNED),
				0
			) AS last_number
		FROM orders
		WHERE SUBSTR(number_display, 5, 6) = ?
		ORDER BY last_number DESC
		LIMIT 1
	`
	err := tx.QueryRow(query, currentYearMonth).Scan(&lastNumber)
	if err != nil{
		lastNumber = 0
	}

	numberDisplay := fmt.Sprintf("ORD-%s-%03d", currentYearMonth, lastNumber+1)
	return numberDisplay
}

func (o *OrderHandler) GetOrders() ([]entity.Order, error) {
	var orders []entity.Order

	user, ok := utils.GetUser(*o.Ctx)
	if !ok {
		return orders, fmt.Errorf("failed to get user from context")
	}

	query := `
		SELECT 
			ord.id, ord.number_display, ord.date, ord.status, ord.total, ord.created_by, od.id, od.product_id, od.qty, od.qty * p.price as od_total ,od.created_by, p.price, p.name
		FROM orders ord
		JOIN order_details od ON ord.id = od.order_id
		JOIN products p ON od.product_id = p.id
		WHERE ord.customer_id = ? AND status = "processing"
		ORDER BY ord.id DESC
	`

	rows, err := o.DB.Query(query, user.Customer.ID)
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	orderMap := make(map[int]*entity.Order)

	for rows.Next() {
		var (
			orderID        int
			numberDisplay  string
			date           string
			status         entity.StatusOrder
			total          float64
			createdBy      int
			orderDetailID  int
			productID      int
			qty            int
			subtotal          float64
			detailCreatedBy int
			price          float64
			productName		string
		)

		err := rows.Scan(&orderID, &numberDisplay, &date, &status, &total, &createdBy, &orderDetailID, &productID, &qty, &subtotal, &detailCreatedBy, &price, &productName)
		if err != nil {
			return orders, err
		}

		if _, exists := orderMap[orderID]; !exists {
			parsedDate, _ := time.Parse("2006-01-02", date)
			orderMap[orderID] = &entity.Order{
				ID:            orderID,
				NumberDisplay: numberDisplay,
				Date:          parsedDate,
				Status:        status,
				Total:        total,
				CreatedBy:     createdBy,
				Details:       []entity.OrderDetail{},
			}
		}

		orderMap[orderID].Details = append(orderMap[orderID].Details, entity.OrderDetail{
			ID:        orderDetailID,
			OrderID:   orderID,
			ProductID: productID,
			Qty:       qty,
			Total: 	   subtotal,
			CreatedBy: detailCreatedBy,
			Product:   entity.Product{ID: productID, Name: productName, Price: price},
		})
	}

	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	return orders, nil
}

func (o *OrderHandler) GetOrderByNumberDisplay(numberDisplay string) (entity.Order, error) {
	var order entity.Order
	user, ok := utils.GetUser(*o.Ctx)
	if !ok {
		return order, fmt.Errorf("failed to get user from context")
	}

	query := `
		SELECT id, number_display, date, status, total, created_by
		FROM orders
		WHERE number_display = ? AND customer_id = ?
		LIMIT 1
	`

	err := o.DB.QueryRow(query, numberDisplay, user.Customer.ID).Scan(
		&order.ID,
		&order.NumberDisplay,
		&order.Date,
		&order.Status,
		&order.Total,
		&order.CreatedBy,
		)

	if err != nil {
		return order, err
	}

	return order, nil
}