package service

import (
	"errors"
	"strings"

	"go-cashier-api/model"
	"go-cashier-api/repository"
)

// ProductService interface defines the methods for product service
type ProductService interface {
	GetAll(name string) ([]model.Product, error)
	GetByID(id int) (*model.Product, error)
	Create(product *model.Product) error
	Update(id int, product *model.Product) error
	Delete(id int) error
}

type ProductServiceImpl struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductService creates a new instance of ProductService
// this called at main.go to initialize the service with the repository
func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &ProductServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// GetAll retrieves all products using the repository
func (s *ProductServiceImpl) GetAll(name string) ([]model.Product, error) {
	products, err := s.productRepo.GetAll(name)

	if err != nil {
		return nil, err
	}

	return products, nil
}

// Create adds a new product using the repository
func (s *ProductServiceImpl) Create(product *model.Product) error {
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
	_, err := s.categoryRepo.GetByID(product.CategoryID)
	if err != nil {
		return errors.New("category does not exist")
	}

	return s.productRepo.Create(product)

}

// GetByID retrieves a product by its ID using the repository
func (s *ProductServiceImpl) GetByID(id int) (*model.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("product not found")
	}

	return product, nil

}

// Update modifies an existing product using the repository
func (s *ProductServiceImpl) Update(id int, product *model.Product) error {
	existing, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("product not found")
	}

	// 2. Apply partial updates
	updated := false

	if strings.TrimSpace(product.Name) != "" && product.Name != existing.Name {
		existing.Name = product.Name
		updated = true
	}

	if product.Price != existing.Price {
		existing.Price = product.Price
		updated = true
	}

	if product.Stock != existing.Stock {
		existing.Stock = product.Stock
		updated = true
	}

	if product.CategoryID > 0 {
		// Check if new category exists
		_, err := s.categoryRepo.GetByID(product.CategoryID)
		if err != nil {
			return errors.New("category does not exist")
		}
		existing.CategoryID = product.CategoryID
	}

	// 3. Save if changes were made
	if !updated {
		return nil // No changes needed
	}

	rowsAffected, err := s.productRepo.Update(existing)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("failed to update category")
	}

	return nil
}

// Delete removes a product by its ID using the repository
func (s *ProductServiceImpl) Delete(id int) error {
	rowsAffected, err := s.productRepo.Delete(id)
	if err != nil {
		return err
	}

	// 4. Consistent error handling
	if rowsAffected == 0 {
		// You might want to distinguish between "not found"
		// and "already deleted"
		return errors.New("product not found")
	}

	return nil
}
