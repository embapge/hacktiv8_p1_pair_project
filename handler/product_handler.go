package handler

import (
	"context"
	"database/sql"
	"fmt"
	"pairproject/entity"
	"pairproject/utils"
)

type ProductHandler struct {
	DB *sql.DB
	Ctx *context.Context
}


func (p *ProductHandler) GetProducts() ([]entity.Product, error) {
	rows, err := p.DB.Query(`
		SELECT p.id, p.name, p.stock, p.description, c.name as category_name, p.price
		FROM products p
		JOIN categories c ON p.category_id = c.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var p entity.Product
		var c entity.Category
		if err := rows.Scan(&p.ID, &p.Name, &p.Stock, &p.Description, &c.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (p *ProductHandler) CreateProduct(product entity.Product)(error){
	user, ok := utils.GetUser(*p.Ctx)
	if !ok {
		return fmt.Errorf("Please Login!")
	}

	// Query SQL untuk menyisipkan data produk baru ke tabel 'products'
	query := `
			INSERT INTO products (name, stock, description, category_id, price, created_by)
			VALUES (?, ?, ?, ?, ?, ?)
	`

	// Menjalankan query dengan parameter dari input produk
	_, err := p.DB.Exec(query, product.Name, product.Stock, product.Description, product.CategoryID, product.Price, user.ID)
	if err != nil {
		return fmt.Errorf("Terjadi kesalahan ketika membuat produk")
	}

	return nil
}
