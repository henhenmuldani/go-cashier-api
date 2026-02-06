package service

import (
	"errors"
	"go-cashier-api/model"
	"go-cashier-api/repository"
	"log"
	"strings"
)

// ProductService interface defines the methods for product service
type ProductService interface {
	GetAll() ([]model.Product, error)
	GetByID(id int) (*model.Product, error)
	Create(product *model.Product) error
	Update(id int, product *model.Product) error
	Delete(id int) error
}

type productService struct {
	productRepo repository.ProductRepository
	// categoryRepo repository.CategoryRepository
}

// NewProductService creates a new instance of ProductService
// this called at main.go to initialize the service with the repository
func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

// GetAll retrieves all products using the repository
func (s *productService) GetAll() ([]model.Product, error) {
	products, err := s.productRepo.GetAll()

	if err != nil {
		return nil, err
	}

	return products, nil
}

// Create adds a new product using the repository
func (s *productService) Create(product *model.Product) error {
	// Validate input
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name is required")
	}

	if product.Price <= 0 {
		return errors.New("product price must be positive")
	}

	if product.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	// Check if category exists
	// _, err := s.categoryRepo.GetByID(product.CategoryID)
	// if err != nil {
	// 	return errors.New("category does not exist")
	// }

	return s.productRepo.Create(product)

}

// GetByID retrieves a product by its ID using the repository
func (s *productService) GetByID(id int) (*model.Product, error) {
	return s.productRepo.GetByID(id)
}

// Update modifies an existing product using the repository
func (s *productService) Update(id int, product *model.Product) error {
	existing, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}
	log.Printf("Existing product: %+v", existing)
	log.Printf("Update data: %+v", product)
	// Update fields
	if product.Name != "" { // Also consider updating Name if needed
		existing.Name = product.Name
	}
	if product.Price > 0 {
		existing.Price = product.Price
	}

	if product.Stock >= 0 {
		existing.Stock = product.Stock
	}

	if product.CategoryID > 0 { // Add this check
		existing.CategoryID = product.CategoryID
	}
	// if product.CategoryID > 0 {
	// 	// Check if new category exists
	// 	_, err := s.categoryRepo.GetByID(product.CategoryID)
	// 	if err != nil {
	// 		return errors.New("category does not exist")
	// 	}
	// 	existing.CategoryID = product.CategoryID
	// }

	return s.productRepo.Update(existing)
}

// Delete removes a product by its ID using the repository
func (s *productService) Delete(id int) error {
	// Check if product exists
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.productRepo.Delete(id)
}
