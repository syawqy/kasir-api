package handler_test

import (
	"kasir-api/internal/model"
)

type MockProductService struct {
	CreateFunc  func(product model.Product) (model.Product, error)
	GetAllFunc  func() ([]model.Product, error)
	GetByIDFunc func(id int) (model.Product, error)
	UpdateFunc  func(id int, product model.Product) (model.Product, error)
	DeleteFunc  func(id int) error
}

func (m *MockProductService) Create(product model.Product) (model.Product, error) {
	return m.CreateFunc(product)
}

func (m *MockProductService) GetAll() ([]model.Product, error) {
	return m.GetAllFunc()
}

func (m *MockProductService) GetByID(id int) (model.Product, error) {
	return m.GetByIDFunc(id)
}

func (m *MockProductService) Update(id int, product model.Product) (model.Product, error) {
	return m.UpdateFunc(id, product)
}

func (m *MockProductService) Delete(id int) error {
	return m.DeleteFunc(id)
}
