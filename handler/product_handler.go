package handler

import (
	"encoding/json" //Encode/decode JSON  API response
	"net/http"      //HTTP server & request handling
	"strconv"       //Convert string to number (for ID from URL)
	"strings"       //String manipulation (trim, split, etc)

	"go-cashier-api/model" // Import model package
	"go-cashier-api/store" // Import store package
)

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Extract ID from URL path if present
	path := strings.TrimPrefix(r.URL.Path, "/api/products")
	path = strings.Trim(path, "/")

	// Handle different HTTP methods
	switch r.Method {
	case http.MethodGet:
		// GET /api/products
		// If no ID is provided, return all products
		if path == "" {
			json.NewEncoder(w).Encode(store.GetAllProducts())
			return
		}

		// GET /api/products/{id}
		// If ID is provided, return specific product
		// Convert ID from string to integer
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		// Fetch product by ID
		product, found := store.GetProductByID(id)
		if !found {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// Return the product as JSON
		json.NewEncoder(w).Encode(product)
	case http.MethodPost:
		// POST /api/products
		// If path is not empty, return 404
		if path != "" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Decode request body into Product struct
		var newProduct model.Product
		if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Create new product
		created := store.CreateProduct(newProduct)
		// Return created product with 201 status
		w.WriteHeader(http.StatusCreated)
		// Return the created product as JSON
		json.NewEncoder(w).Encode(created)
	case http.MethodPut:
		// PUT /api/products/{id}
		// Convert ID from string to integer
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		// Decode request body into Product struct
		var product model.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}

		// Update product
		updated, ok := store.UpdateProduct(id, product)
		if !ok {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// Return the updated product as JSON
		json.NewEncoder(w).Encode(updated)
	case http.MethodDelete:
		// DELETE /api/products/{id}
		// Convert ID from string to integer
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		// Delete product
		if !store.DeleteProduct(id) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// Return success message
		json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
	default:
		// If method is not supported, return 405
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
