package service

import (
	"go-cashier-api/model"
)

// ProductRepository defines the interface for product data operations
// This allows for decoupling the service layer from the data layer
type ProductRepository interface {
	GetAll() ([]model.Product, error)
	Create(data *model.Product) error
	GetByID(id int) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id int) error
}

// ProductService struct holds the repository dependency
type ProductService struct {
	repo ProductRepository
}

// NewProductService creates a new instance of ProductService
// this called at main.go to initialize the service with the repository
func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// GetAll retrieves all products using the repository
func (s *ProductService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll()
}

// Create adds a new product using the repository
func (s *ProductService) Create(data *model.Product) error {
	return s.repo.Create(data)
}

// GetByID retrieves a product by its ID using the repository
func (s *ProductService) GetByID(id int) (*model.Product, error) {
	return s.repo.GetByID(id)
}

// Update modifies an existing product using the repository
func (s *ProductService) Update(product *model.Product) error {
	return s.repo.Update(product)
}

// Delete removes a product by its ID using the repository
func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
