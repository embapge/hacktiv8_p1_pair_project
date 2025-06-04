package handler

import (
	"database/sql"
	"pairproject/entity"
)

type ProductHandler struct {
	DB *sql.DB
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
