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
	User         int64     `json:"added_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	StockStatus  string    `json:"status"`
}
type ProductModel struct {
	DB *sql.DB
}

// product validation
func ValidateProduct(v *validator.Validator, product *Product) {

	v.Check(validator.NotBlank(product.PName), "product_name", "Name must be provided")
	v.Check(validator.NotZero(product.PQuantity), "product_quantity", "Quantity must be Provided")
	v.Check(validator.NotZeroF(product.PPrice), "product_price", "Price must be Provided")
	v.Check(validator.NotBlank(product.PDescription), "product_description", "Description must be Provided")

	v.Check(validator.MaxLengths(product.PName, 100), "product_name", "Name must be less than 100 characters")
	v.Check(validator.MaxLengths(product.PDescription, 255), "product_description", "Description must be less than 255 characters")

	v.Check(validator.NotPositive(product.PQuantity), "product_quantity", "Quantity must be greater than 0")
	v.Check(validator.NotPositiveF(product.PPrice), "product_price", "Price must be greater than 0")

}

// Insert record into the the products database
// func (m *ProductModel) Insert(product *Product) error {
// 	query := `
// 		INSERT INTO products (product_name, product_quantity, product_price, product_description, added_by, created_at, updated_at)
// 		VALUES ($1, $2, $3, $4, $5, now(), now())
// 		RETURNING product_id, created_at, updated_at
// 	`
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	//scan method
// 	return m.DB.QueryRowContext(
// 		ctx,
// 		query,
// 		product.PName,
// 		product.PQuantity,
// 		product.PPrice,
// 		product.PDescription,
// 		product.User,
// 	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

// }
// Insert a product associated with a user
func (m *ProductModel) Insert(product *Product) error {
	query := `
        INSERT INTO products (product_name, product_quantity, product_price, product_description, added_by, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, now(), now())
        RETURNING product_id, created_at, updated_at
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the user ID to link product to the user
	return m.DB.QueryRowContext(
		ctx,
		query,
		product.PName,
		product.PQuantity,
		product.PPrice,
		product.PDescription,
		product.User, // Pass the authenticated user's ID here
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

//Get all records from the products database
// func (m *ProductModel) GetAll() ([]*Product, error) {
// 	query := `
// 		SELECT product_id, product_name, product_quantity, product_price, product_description, added_by, created_at, updated_at
// 		FROM products
// 		WHERE product_id = $1 AND added_by = $2
// 	`
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	rows, err := m.DB.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	products := make([]*Product, 0)

// 	for rows.Next() {
// 		product := &Product{}
// 		if err := rows.Scan(
// 			&product.ID,
// 			&product.PName,
// 			&product.PQuantity,
// 			&product.PPrice,
// 			&product.PDescription,
// 			&product.CreatedAt,
// 			&product.UpdatedAt,
// 		); err != nil {
// 			return nil, err
// 		}
// 		products = append(products, product)
// 	}
// 	// Check for iteration error
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

//		return products, nil
//	}
//
// Get all products for a specific user
func (m *ProductModel) GetAllForUser(userID int64) ([]*Product, error) {
	query := `
        SELECT product_id, product_name, product_quantity, product_price, product_description, added_by, created_at, updated_at
        FROM products
        WHERE added_by = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
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
			&product.User,
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

// // search
// func (m *ProductModel) Search(query string) ([]*Product, error) {
// 	stmt := `
// 		SELECT product_id, product_name, product_description, product_quantity, product_price
// 		FROM products
// 		WHERE LOWER(product_name) LIKE LOWER($1) OR LOWER(product_description) LIKE LOWER($1)
// 		OR CAST(product_quantity AS TEXT) LIKE $1 OR CAST(product_price AS TEXT) LIKE $1
// 	`
// 	rows, err := m.DB.Query(stmt, "%"+query+"%")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var products []*Product
// 	for rows.Next() {
// 		var p Product
// 		err := rows.Scan(&p.ID, &p.PName, &p.PDescription, &p.PQuantity, &p.PPrice)
// 		if err != nil {
// 			return nil, err
// 		}
// 		products = append(products, &p)
// 	}

//		return products, nil
//	}
func (m *ProductModel) SearchByUser(query string, userID int64) ([]*Product, error) {
	stmt := `
		SELECT product_id, product_name, product_description, product_quantity, product_price, added_by
		FROM products
		WHERE added_by = $2 AND (
			LOWER(product_name) LIKE LOWER($1)
			OR LOWER(product_description) LIKE LOWER($1)
			OR CAST(product_quantity AS TEXT) LIKE $1
			OR CAST(product_price AS TEXT) LIKE $1
		)
	`

	rows, err := m.DB.Query(stmt, "%"+query+"%", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.PName, &p.PDescription, &p.PQuantity, &p.PPrice, &p.User)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}
