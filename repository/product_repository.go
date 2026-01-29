package repository

import (
	"database/sql"
	"errors"

	"go-cashier-api/model"
)

// implementation of repository pattern for product entity
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new instance of ProductRepository
// this called at main.go to initialize the repository with the database connection
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Query functions
// GetAllProducts returns all products
func (repo *ProductRepository) GetAll() ([]model.Product, error) {
	// query all products from database
	query := "SELECT id, name, price, stock FROM products"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over rows and scan into products slice
	products := make([]model.Product, 0)
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetProductByID returns a product by its ID
func (repo *ProductRepository) GetByID(id int) (*model.Product, error) {
	// query product by ID from database
	query := "SELECT id, name, price, stock FROM products WHERE id = $1"

	// scan result into product struct
	var p model.Product
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Command functions
// CreateProduct adds a new product to the store
func (repo *ProductRepository) Create(p *model.Product) error {
	// insert new product into database
	query := "INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id"
	err := repo.db.QueryRow(query, p.Name, p.Price, p.Stock).Scan(&p.ID)
	return err
}

// UpdateProduct updates an existing product by its ID
func (repo *ProductRepository) Update(product *model.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

// DeleteProduct removes a product by its ID
func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return err
}
