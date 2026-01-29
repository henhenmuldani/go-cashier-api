package handler

import (
	"encoding/json" //Encode/decode JSON  API response
	"net/http"      //HTTP server & request handling
	"strconv"       //Convert string to number (for ID from URL)
	"strings"       //String manipulation (trim, split, etc)

	"go-cashier-api/model"   // Import model package
	"go-cashier-api/service" // Import service package
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
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
			categories, err := h.service.GetAll()
			if err != nil {
				http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(categories)
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
		category, err := h.service.GetByID(id)
		if err != nil {
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
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Disallow unknown fields
		if err := decoder.Decode(&newCategory); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Create new category
		err := h.service.Create(&newCategory)
		if err != nil {
			http.Error(w, "Failed to create category", http.StatusInternalServerError)
			return
		}
		// Return created category with 201 status
		w.WriteHeader(http.StatusCreated)
		// Return the created category as JSON
		json.NewEncoder(w).Encode(newCategory)
	case http.MethodPut:
		// PUT /api/categories/{id}
		if path == "" {
			http.Error(w, "Category ID required", http.StatusBadRequest)
			return
		}

		// Convert ID from string to integer
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		// Decode request body into Category struct
		var category model.Category
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Disallow unknown fields
		if err := decoder.Decode(&category); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// Update category
		category.ID = id // Ensure the ID is set from the URL
		err = h.service.Update(&category)
		if err != nil {
			http.Error(w, "Failed to update category", http.StatusInternalServerError)
			return
		}

		// Return the updated category as JSON
		json.NewEncoder(w).Encode(category)
	case http.MethodDelete:
		// DELETE /api/categories/{id}
		// Convert ID from string to integer
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		// Delete category
		err = h.service.Delete(id)
		if err != nil {
			http.Error(w, "Failed to delete category", http.StatusInternalServerError)
			return
		}

		// Return success message
		json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
	default:
		// If method is not supported, return 405
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
