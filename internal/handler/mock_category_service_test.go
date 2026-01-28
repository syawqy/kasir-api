package handler_test

import (
	"kasir-api/internal/model"
)

type MockCategoryService struct {
	CreateFunc  func(category model.Category) (model.Category, error)
	GetAllFunc  func() ([]model.Category, error)
	GetByIDFunc func(id int) (model.Category, error)
	UpdateFunc  func(id int, category model.Category) (model.Category, error)
	DeleteFunc  func(id int) error
}

func (m *MockCategoryService) Create(category model.Category) (model.Category, error) {
	return m.CreateFunc(category)
}

func (m *MockCategoryService) GetAll() ([]model.Category, error) {
	return m.GetAllFunc()
}

func (m *MockCategoryService) GetByID(id int) (model.Category, error) {
	return m.GetByIDFunc(id)
}

func (m *MockCategoryService) Update(id int, category model.Category) (model.Category, error) {
	return m.UpdateFunc(id, category)
}

func (m *MockCategoryService) Delete(id int) error {
	return m.DeleteFunc(id)
}
