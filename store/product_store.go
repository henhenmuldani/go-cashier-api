package store

import "go-cashier-api/model"

// In-memory product data
var products = []model.Product{
	{ID: 1, Name: "Indomie Godog", Price: 3500, Stock: 10},
	{ID: 2, Name: "Vit 1000ml", Price: 3000, Stock: 40},
	{ID: 3, Name: "Kecap", Price: 12000, Stock: 20},
}

// Query functions
// GetAllProducts returns all products
func GetAllProducts() []model.Product {
	return products
}

// GetProductByID returns a product by its ID
func GetProductByID(id int) (model.Product, bool) {
	// find product by ID
	// iterate over products
	for _, p := range products {
		// if ID matches, return product
		if p.ID == id {
			// found, return true
			return p, true
		}
	}
	// if not found, return false
	return model.Product{}, false
}

// Command functions
// CreateProduct adds a new product to the store
func CreateProduct(p model.Product) model.Product {
	// assign new ID
	p.ID = len(products) + 1
	// append to products slice
	products = append(products, p)
	// return created product
	return p
}

// UpdateProduct updates an existing product by its ID
func UpdateProduct(id int, updated model.Product) (model.Product, bool) {
	// find product by ID
	// iterate over products
	for i, p := range products {
		if p.ID == id {
			updated.ID = id
			products[i] = updated
			return updated, true
		}
	}
	return model.Product{}, false
}

func DeleteProduct(id int) bool {
	// find product by ID
	// iterate over products
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			return true
		}
	}
	return false
}
