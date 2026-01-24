package handler

import (
	"encoding/json" //Encode/decode JSON  API response
	"net/http"      //HTTP server & request handling
	"strconv"       //Convert string to number (for ID from URL)
	"strings"       //String manipulation (trim, split, etc)

	"go-cashier-api/model" // Import model package
	"go-cashier-api/store" // Import store package
)

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Extract ID from URL path if present
	path := strings.TrimPrefix(r.URL.Path, "/api/categories")
	path = strings.Trim(path, "/")

	// Handle different HTTP methods
	switch r.Method {
	case http.MethodGet:
		// GET /api/categories
		// If no ID is provided, return all categories
		if path == "" {
			json.NewEncoder(w).Encode(store.GetAllCategories())
			return
		}

		// GET /api/categories/{id}
		// If ID is provided, return specific category
		// Convert ID from string to integer
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		// Fetch category by ID
		category, found := store.GetCategoryByID(id)
		if !found {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		// Return the category as JSON
		json.NewEncoder(w).Encode(category)
	case http.MethodPost:
		// POST /api/categories
		// If path is not empty, return 404
		if path != "" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Decode request body into Category struct
		var newCategory model.Category
		if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Create new category
		created := store.CreateCategory(newCategory)
		// Return created category with 201 status
		w.WriteHeader(http.StatusCreated)
		// Return the created category as JSON
		json.NewEncoder(w).Encode(created)
	case http.MethodPut:
		// PUT /api/categories/{id}
		// Convert ID from string to integer
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		// Decode request body into Category struct
		var category model.Category
		if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}

		// Update category
		updated, ok := store.UpdateCategory(id, category)
		if !ok {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		// Return the updated category as JSON
		json.NewEncoder(w).Encode(updated)
	case http.MethodDelete:
		// DELETE /api/categories/{id}
		// Convert ID from string to integer
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		// Delete category
		if !store.DeleteCategory(id) {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		// Return success message
		json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
	default:
		// If method is not supported, return 405
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
