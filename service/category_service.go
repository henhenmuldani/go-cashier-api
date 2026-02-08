package service

import (
	"errors"
	"strings"

	"go-cashier-api/model"
	"go-cashier-api/repository"
)

type CategoryService interface {
	GetAll() ([]model.Category, error)
	GetByID(id int) (*model.Category, error)
	Create(category *model.Category) error
	Update(id int, category *model.Category) error
	Delete(id int) error
}

type CategoryServiceImpl struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{repo: repo}
}

func (s *CategoryServiceImpl) GetAll() ([]model.Category, error) {
	categories, err := s.repo.GetAll()

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CategoryServiceImpl) Create(category *model.Category) error {
	// Business validation
	if strings.TrimSpace(category.Name) == "" {
		return errors.New("category name is required")
	}

	// Check for duplicate name (business rule)
	// ...

	return s.repo.Create(category)
}

func (s *CategoryServiceImpl) GetByID(id int) (*model.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("category not found")
	}

	return category, nil
}

func (s *CategoryServiceImpl) Update(id int, category *model.Category) error {
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

func (s *CategoryServiceImpl) Delete(id int) error {
	// 1. Consider implementing a soft delete pattern
	//    (add DeletedAt field to your model)

	// 2. Check business constraints BEFORE attempting deletion
	// In a real application, you would check if category has products
	// before deleting
	// hasProducts, err := s.productRepo.HasProductsInCategory(id)
	// if err != nil {
	//     return fmt.Errorf("failed to check category products: %w", err)
	// }
	// if hasProducts {
	//     return errors.New("cannot delete category with existing products")
	// }

	// 3. Single attempt with proper error handling
	rowsAffected, err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// 4. Consistent error handling
	if rowsAffected == 0 {
		// You might want to distinguish between "not found"
		// and "already deleted"
		return errors.New("category not found")
	}

	// 5. Optional: Clear cache or trigger events

	return nil
}
