package handler

import (
	"encoding/json" //Encode/decode JSON  API response
	"net/http"      //HTTP server & request handling
	"strconv"
	"strings"

	"go-cashier-api/model"        // Import model package
	"go-cashier-api/pkg/response" // Import response
	"go-cashier-api/service"      // Import service package
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler creates a new ProductHandler with the given ProductService
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// HandleProducts - GET /api/produk
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w)
	case http.MethodPost:
		h.create(w, r)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleProductByID - GET/PUT/DELETE /api/produk/{id}
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
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
// @Summary Get all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {array} model.ProductResponseSwagger
// @Router /api/products [get]
func (h *ProductHandler) getAll(w http.ResponseWriter) {
	// GET /api/products
	products, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}
	response.JSON(w, http.StatusOK, products)
	// response.JSON(w, http.StatusOK, map[string]interface{}{
	// 	"data": products,
	// })
}

// create godoc
// @Summary Create product
// @Tags Products
// @Accept json
// @Produce json
// @Param product body model.CreateProductRequestSwagger true "Create product payload"
// @Router /api/products [post]
func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	// Decode request body into Product struct
	var newProduct model.Product
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Disallow unknown fields
	if err := decoder.Decode(&newProduct); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create new product
	err := h.service.Create(&newProduct)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	// Return the created product as JSON
	response.JSON(w, http.StatusCreated, newProduct)

}

// getByID godoc
// @Summary Get product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Router /api/products/{id} [get]
func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	// Fetch product by ID
	product, err := h.service.GetByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Product not found")
		return
	}

	// Return the product as JSON
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"name":          product.Name,
		"price":         product.Price,
		"stock":         product.Stock,
		"category_id":   product.CategoryID,
		"category_name": product.Category.Name,
	})
}

// update godoc
// @Summary Update product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body model.CreateProductRequestSwagger true "Update product payload"
// @Router /api/products/{id} [put]
func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// Decode request body into Product struct
	var product model.Product
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Disallow unknown fields
	if err := decoder.Decode(&product); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Update product
	product.ID = id // Ensure the ID is set from the URL
	err = h.service.Update(id, &product)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	// Return the updated product as JSON
	response.JSON(w, http.StatusOK, map[string]interface{}{"message": "Product updated successfully", "data": product})
}

// delete godoc
// @Summary Delete product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Router /api/products/{id} [delete]
func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// Delete product
	if err := h.service.Delete(id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return success message
	response.JSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
