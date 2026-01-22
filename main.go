package main

import (
	"encoding/json" //Encode/decode JSON  API response)
	"fmt"           //Print to console (fmt.Println)
	"net/http"      //HTTP server & handling
	"strconv"       //Convert string to number (for ID from URL)
	"strings"       //String manipulation (trim, split, etc)
)

// Product struct to represent product data, including JSON tags for serialization
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

// In-memory product data
var products = []Product{
	{ID: 1, Name: "Indomie Godog", Price: 3500, Stock: 10},
	{ID: 2, Name: "Vit 1000ml", Price: 3000, Stock: 40},
	{ID: 3, Name: "kecap", Price: 12000, Stock: 20},
}

// Handler to get product by ID
func getProductByID(w http.ResponseWriter, r *http.Request) {
	// parsing ID from URL
	// url : /api/product/123 => ID = 123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	// change string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	// find product by ID
	for _, p := range products {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	// if product not found
	http.Error(w, "Product not found", http.StatusNotFound)
}

// Handler to update product by ID
func updateProductByID(w http.ResponseWriter, r *http.Request) {
	// parsing ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	// change string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	// read data from request
	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// find and update product by ID
	for i, p := range products {
		if p.ID == id {
			products[i].Name = updatedProduct.Name
			products[i].Price = updatedProduct.Price
			products[i].Stock = updatedProduct.Stock
			// can like this too:
			// products[i].ID = id // ID must be set again
			// products[i] = updatedProduct

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products[i])
			return
		}
	}

	// if product not found
	http.Error(w, "Product not found", http.StatusNotFound)
}

// Handler to delete product by ID
func deleteProductByID(w http.ResponseWriter, r *http.Request) {
	// parsing ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	// change string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	// find and delete product by ID
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Product deleted successfully",
			})
			return
		}
	}

	// if product not found
	http.Error(w, "Product not found", http.StatusNotFound)
}

// Main function to start the server and define routes
func main() {
	// GET localhost:8080/api/products/{id}
	// PUT localhost:8080/api/products/{id}
	// DELETE localhost:8080/api/products/{id}
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProductByID(w, r)
		case http.MethodPut:
			updateProductByID(w, r)
		case http.MethodDelete:
			deleteProductByID(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET localhost:8080/api/products
	// POST localhost:8080/api/products
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// return all product data
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products)

		case http.MethodPost:
			// read data from request
			var newProduct Product
			err := json.NewDecoder(r.Body).Decode(&newProduct)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// insert data to Array Products
			newProduct.ID = len(products) + 1
			products = append(products, newProduct)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(newProduct)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// 	Start the server
	fmt.Println("Starting server on localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
