package handler

import (
	"encoding/json" //Encode/decode JSON  API response
	"net/http"      //HTTP server & request handling
	"strconv"       //Convert string to number (for ID from URL)
	"strings"       //String manipulation (trim, split, etc)

	"go-cashier-api/model"        // Import model package
	"go-cashier-api/pkg/response" //
	"go-cashier-api/service"      // Import service package
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w)
	case http.MethodPost:
		h.create(w, r)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r)
	case http.MethodPut:
		h.update(w, r)
	case http.MethodDelete:
		h.delete(w, r)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// getAll godoc
// @Summary Get all categories
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {array} model.Category
// @Router /api/categories [get]
func (h *CategoryHandler) getAll(w http.ResponseWriter) {
	categories, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}
	// Optional: Return empty array if no categories found
	if categories == nil {
		categories = []model.Category{}
	}

	response.JSON(w, http.StatusOK, categories)
}

// create godoc
// @Summary Create category
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body model.CreateCategoryRequestSwagger true "Create category payload"
// @Router /api/categories [post]
func (h *CategoryHandler) create(w http.ResponseWriter, r *http.Request) {
	var newCategory model.Category
	// Validate content type
	if r.Header.Get("Content-Type") != "application/json" {
		response.Error(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Disallow unknown fields
	if err := decoder.Decode(&newCategory); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the category data (could also be in service)
	if strings.TrimSpace(newCategory.Name) == "" {
		response.Error(w, http.StatusBadRequest, "Category name is required")
		return
	}

	if err := h.service.Create(&newCategory); err != nil {
		// Map service errors to appropriate HTTP status codes
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "exists") {
			statusCode = http.StatusConflict
		}
		response.Error(w, statusCode, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, newCategory)
}

// ISSUE: This assumes path is exactly "/api/categories/{id}"
// Better to use URL parameters or a router that extracts ID
// getByID godoc
// @Summary Get category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) getByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL - THIS IS FRAGILE
	// Better approach: Use a router like gorilla/mux, chi, or httprouter
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	// Fetch category by ID
	category, err := h.service.GetByID(id)
	if err != nil {
		// Differentiate between not found and internal errors
		if strings.Contains(err.Error(), "not found") {
			response.Error(w, http.StatusNotFound, "Category not found")
		} else {
			response.Error(w, http.StatusInternalServerError, "Failed to fetch category")
		}
		return
	}

	// Return the entire category object
	response.JSON(w, http.StatusOK, category)
}

// update godoc
// @Summary Update category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body model.CreateCategoryRequestSwagger true "Update category payload"
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	// Decode request body
	if r.Header.Get("Content-Type") != "application/json" {
		response.Error(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	// Decode request body into Category struct
	var category model.Category
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Disallow unknown fields
	if err := decoder.Decode(&category); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Update category
	err = h.service.Update(id, &category)
	if err != nil {
		// Map errors to appropriate status codes
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "duplicate") {
			statusCode = http.StatusConflict
		}
		response.Error(w, statusCode, err.Error())
		return
	}
	// Return success response - optionally fetch updated category
	updatedCategory, _ := h.service.GetByID(id)
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Category updated successfully",
		"data":    updatedCategory,
	})
}

// delete godoc
// @Summary Delete category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	// Delete category
	if err := h.service.Delete(id); err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "in use") || strings.Contains(err.Error(), "constraint") {
			statusCode = http.StatusConflict
		}
		response.Error(w, statusCode, err.Error())
		return
	}

	// Return success message
	response.JSON(w, http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
