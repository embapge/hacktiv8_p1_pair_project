package handler

import (
	"context"
	"database/sql"
	"fmt"
	"pairproject/entity"
	"pairproject/utils"
)

type CategoryHandler struct {
	DB *sql.DB
	Ctx *context.Context
}

func (c *CategoryHandler) CreateCategory(name string)(error){
	user, ok := utils.GetUser(*c.Ctx)
	if !ok {
		return fmt.Errorf("Please Login!")
	}

	query := `INSERT INTO categories (name, created_by) VALUES (?, ?)`
	_, err := c.DB.Exec(query, name, user.ID)
	if err != nil {
		// Jika terjadi error saat menyimpan ke database, tampilkan pesan error
		return fmt.Errorf("Terjadi kesalahan ketika membuat produk")
	}

	return nil
}

func (c *CategoryHandler) GetCategories() ([]entity.Category, error) {
	rows, err := c.DB.Query("SELECT id, name FROM categories")
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