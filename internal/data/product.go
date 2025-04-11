package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/AbnerBobad/final_project/internal/validator"
)

type Product struct {
	ID           int64     `json:"product_id"`
	PName        string    `json:"product_name"`
	PQuantity    int64     `json:"product_quantity"`
	PPrice       float64   `json:"product_price"`
	PDescription string    `json:"product_description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type ProductModel struct {
	DB *sql.DB
}

// validation checker
func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.PName != "", "product_name", "Must be provided")
	v.Check(product.PQuantity != " ", "product_quantity", "Must be provided")
	v.Check(product.PPrice != "", "product_price", "Must be provided")

	v.Check(len(product.PName) <= 100, "product_name", "Must not exceed 100 characters")
	v.Check(product.PQuantity > 0, "product_quantity", "Must be greater than 0")
	v.Check(product.PPrice > 0.0, "product_price", "Must be greater than 0")
	v.Check(len(product.PDescription) <= 255, "product_description", "Must not exceed 255 characters")
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
