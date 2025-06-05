package handler

import (
	"context"
	"database/sql"
	"fmt"
	"pairproject/entity"
	"pairproject/utils"
)

type OrderDetailHandler struct {
	DB *sql.DB
	Ctx *context.Context
}

func (od *OrderDetailHandler) UpdateDetail(id int, qty int) (entity.OrderDetail, error){
	var orderDetail entity.OrderDetail
	user, ok := utils.GetUser(*od.Ctx)
	if !ok {
		return orderDetail, fmt.Errorf("failed to get user from context")
	}

	res, err := od.DB.Exec("UPDATE order_details SET qty = ?, updated_by = ? WHERE id = ?", qty, user.ID, id)

	if err != nil{
		return orderDetail, fmt.Errorf("Terjadi kesalahan update data: %s", err)
	}

	_, err = res.LastInsertId()
	if err != nil {
		return orderDetail, fmt.Errorf("Terjadi kesalahan mengambil order detail id: %s", err)
	}

	return entity.OrderDetail{Qty: qty}, nil
}