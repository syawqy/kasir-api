package service

import (
	"errors"
	"kasir-api/internal/model"
	"kasir-api/internal/repository"
)

type CategoryService interface {
	Create(category model.Category) (model.Category, error)
	GetAll() ([]model.Category, error)
	GetByID(id int) (model.Category, error)
	Update(id int, category model.Category) (model.Category, error)
	Delete(id int) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(category model.Category) (model.Category, error) {
	if category.Name == "" {
		return model.Category{}, errors.New("name is required")
	}
	return s.repo.Create(category)
}

func (s *categoryService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()
}

func (s *categoryService) GetByID(id int) (model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) Update(id int, category model.Category) (model.Category, error) {
	if category.Name == "" {
		return model.Category{}, errors.New("name is required")
	}
	return s.repo.Update(id, category)
}

func (s *categoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
