package service

import (
	"errors"
	"kasir-api/internal/model"
	"kasir-api/internal/repository"
)

type ProductService interface {
	Create(product model.Product) (model.Product, error)
	GetAll() ([]model.Product, error)
	GetByID(id int) (model.Product, error)
	Update(id int, product model.Product) (model.Product, error)
	Delete(id int) error
}

type productService struct {
	repo    repository.ProductRepository
	catRepo repository.CategoryRepository
}

func NewProductService(repo repository.ProductRepository, catRepo repository.CategoryRepository) ProductService {
	return &productService{repo: repo, catRepo: catRepo}
}

func (s *productService) Create(product model.Product) (model.Product, error) {
	if product.Name == "" {
		return model.Product{}, errors.New("name is required")
	}
	if product.Price < 0 {
		return model.Product{}, errors.New("price cannot be negative")
	}

	// Validate category exists
	_, err := s.catRepo.GetByID(product.CategoryID)
	if err != nil {
		return model.Product{}, errors.New("category not found")
	}

	return s.repo.Create(product)
}

func (s *productService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) GetByID(id int) (model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) Update(id int, product model.Product) (model.Product, error) {
	if product.Name == "" {
		return model.Product{}, errors.New("name is required")
	}
	if product.Price < 0 {
		return model.Product{}, errors.New("price cannot be negative")
	}

	// Validate category exists if category_id is being updated/set
	if product.CategoryID != 0 {
		_, err := s.catRepo.GetByID(product.CategoryID)
		if err != nil {
			return model.Product{}, errors.New("category not found")
		}
	}

	return s.repo.Update(id, product)
}

func (s *productService) Delete(id int) error {
	return s.repo.Delete(id)
}
