package handler

import (
	"encoding/json" //Encode/decode JSON  API response
	"net/http"      //HTTP server & request handling
	"strconv"       //Convert string to number (for ID from URL)
	"strings"       //String manipulation (trim, split, etc)

	"go-cashier-api/model"   // Import model package
	"go-cashier-api/service" // Import service package
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	service *service.ProductService
}

// NewProductHandler creates a new ProductHandler with the given ProductService
func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

// HandleProducts handles HTTP requests for products based on the method and URL path
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
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
			products, err := h.service.GetAll()
			if err != nil {
				http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(products)
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
		product, err := h.service.GetByID(id)
		if err != nil {
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
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Disallow unknown fields
		if err := decoder.Decode(&newProduct); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Create new product
		err := h.service.Create(&newProduct)
		if err != nil {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return
		}
		// Return created product with 201 status
		w.WriteHeader(http.StatusCreated)
		// Return the created product as JSON
		json.NewEncoder(w).Encode(newProduct)
	case http.MethodPut:
		// PUT /api/products/{id}
		if path == "" {
			http.Error(w, "Product ID required", http.StatusBadRequest)
			return
		}

		// Convert ID from string to integer
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		// Decode request body into Product struct
		var product model.Product
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Disallow unknown fields
		if err := decoder.Decode(&product); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// Update product
		product.ID = id // Ensure the ID is set from the URL
		err = h.service.Update(&product)
		if err != nil {
			http.Error(w, "Failed to update product", http.StatusInternalServerError)
			return
		}
		// Return the updated product as JSON
		json.NewEncoder(w).Encode(product)
	case http.MethodDelete:
		// DELETE /api/products/{id}
		// Convert ID from string to integer
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		// Delete product
		err = h.service.Delete(id)
		if err != nil {
			http.Error(w, "Failed to delete product", http.StatusInternalServerError)
			return
		}

		// Return success message
		json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
	default:
		// If method is not supported, return 405
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
