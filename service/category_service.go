package service

import (
	"errors"
	"go-cashier-api/model"
	"go-cashier-api/repository"
	"strings"
)

type CategoryService interface {
	GetAll() ([]model.Category, error)
	GetByID(id int) (*model.Category, error)
	Create(category *model.Category) error
	Update(id int, category *model.Category) error
	Delete(id int) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()

}

func (s *categoryService) Create(category *model.Category) error {
	// Business validation
	if strings.TrimSpace(category.Name) == "" {
		return errors.New("category name is required")
	}

	// Check for duplicate name (business rule)
	// ...

	return s.repo.Create(category)
}

func (s *categoryService) GetByID(id int) (*model.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("category not found")
	}

	return category, nil
}

func (s *categoryService) Update(id int, category *model.Category) error {
	// 1. Get existing
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("category not found")
	}

	// 2. Apply partial updates
	updated := false

	if strings.TrimSpace(category.Name) != "" && category.Name != existing.Name {
		existing.Name = category.Name
		updated = true
	}

	if category.Description != existing.Description {
		existing.Description = category.Description
		updated = true
	}

	// 3. Save if changes were made
	if !updated {
		return nil // No changes needed
	}

	rowsAffected, err := s.repo.Update(existing)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("failed to update category")
	}

	return nil
}

func (s *categoryService) Delete(id int) error {
	// Check if category exists
	rowsAffected, err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("category not found")
	}

	// In a real application, you would check if category has products
	// before deleting

	return nil
}
