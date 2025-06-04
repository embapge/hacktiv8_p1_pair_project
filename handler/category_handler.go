package handler

import (
	"database/sql"
	"pairproject/entity"
)

type CategoryHandler struct {
	DB *sql.DB
}

func (h *CategoryHandler) GetCategories() ([]entity.Category, error) {
	rows, err := h.DB.Query("SELECT id, name FROM categories")
	var emptyCategory []entity.Category
	if err != nil {
		return emptyCategory, err
	}
	defer rows.Close()

	var categories []entity.Category

	for rows.Next() {
		var c entity.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return emptyCategory, err
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return emptyCategory, err
	}
	return categories, nil
}