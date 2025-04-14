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

// Update modifies an existing product in the database.
func (m *ProductModel) Update(product *Product) error {
	query := `
		UPDATE products
		SET product_name = $1,
			product_quantity = $2,
			product_price = $3,
			product_description = $4,
			updated_at = now()
		WHERE product_id = $5
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query,
		product.PName,
		product.PQuantity,
		product.PPrice,
		product.PDescription,
		product.ID,
	)

	return err
}

// Delete removes a product from the database.
func (m *ProductModel) Delete(id int64) error {
	query := `DELETE FROM products WHERE product_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// GetByID fetches a product by its ID
func (m *ProductModel) GetByID(id int64) (*Product, error) {
	query := `
		SELECT product_id, product_name, product_quantity, product_price, product_description, created_at, updated_at
		FROM products
		WHERE product_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	product := &Product{}
	err := row.Scan(
		&product.ID,
		&product.PName,
		&product.PQuantity,
		&product.PPrice,
		&product.PDescription,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}
