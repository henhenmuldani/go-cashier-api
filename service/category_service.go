package service

import (
	"go-cashier-api/model"
)

type CategoryRepository interface {
	GetAll() ([]model.Category, error)
	Create(data *model.Category) error
	GetByID(id int) (*model.Category, error)
	Update(category *model.Category) error
	Delete(id int) error
}

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) Create(data *model.Category) error {
	return s.repo.Create(data)
}

func (s *CategoryService) GetByID(id int) (*model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Update(category *model.Category) error {
	return s.repo.Update(category)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
