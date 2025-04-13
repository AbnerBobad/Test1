package data

import (
	"context"
	"database/sql"
	"time"
)

type Product struct {
	ID           int64     `json:"product_id"`
	PName        string    `json:"product_name"`
	PQuantity    int64     `json:"product_quantity"`
	PPrice       float64   `json:"product_price"`
	PDescription string    `json:"product_description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	StockStatus  string    `json:"status"`
}
type ProductModel struct {
	DB *sql.DB
}

// Insert record into the the products database
func (m *ProductModel) Insert(product *Product) error {
	query := `
		INSERT INTO products (product_name, product_quantity, product_price, product_description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, now(), now())
		RETURNING product_id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//scan method
	return m.DB.QueryRowContext(
		ctx,
		query,
		product.PName,
		product.PQuantity,
		product.PPrice,
		product.PDescription,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

}

// // Get all records from the products database
func (m *ProductModel) GetAll() ([]*Product, error) {
	query := `
		SELECT product_id, product_name, product_quantity, product_price, product_description, created_at, updated_at
		FROM products
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*Product, 0)

	for rows.Next() {
		product := &Product{}
		if err := rows.Scan(
			&product.ID,
			&product.PName,
			&product.PQuantity,
			&product.PPrice,
			&product.PDescription,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	// Check for iteration error
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// // Update record in the products database
// func (m *ProductModel) Update(p *Product) error {
// 	query := `
// 		UPDATE products
// 		SET name = $1, description = $2, quantity = $3, price = $4, updated_at = now()
// 		WHERE id = $5
// 	`
// 	args := []any{p.Name, p.Description, p.Quantity, p.Price, p.ID}
// 	_, err := m.DB.Exec(query, args...)
// 	return err
// }

// // Delete record from the products database
// func (m *ProductModel) Delete(id int64) error {
// 	query := `
// 		DELETE FROM products
// 		WHERE id = $1
// 	`
// 	_, err := m.DB.Exec(query, id)
// 	return err
// }

// // Get a single record from the products database
// func (m *ProductModel) Get(id int64) (*Product, error) {
// 	query := `
// 		SELECT id, name, description, quantity, price, created_at, updated_at
// 		FROM products
// 		WHERE id = $1
// 	`
// 	p := &Product{}
// 	err := m.DB.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Quantity, &p.Price, &p.CreatedAt, &p.UpdatedAt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return p, nil
// }
