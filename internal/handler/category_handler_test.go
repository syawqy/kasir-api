package handler_test

import (
	"bytes"
	"encoding/json"
	"kasir-api/internal/handler"
	"kasir-api/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCategory(t *testing.T) {
	mockService := &MockCategoryService{
		CreateFunc: func(category model.Category) (model.Category, error) {
			category.ID = 1
			return category, nil
		},
	}
	h := handler.NewCategoryHandler(mockService)

	payload := []byte(`{"name":"Electronics", "description":"Gadgets"}`)
	req, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(h.HandleCategories)
	handlerFunc.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var created model.Category
	err = json.Unmarshal(rr.Body.Bytes(), &created)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.ID != 1 {
		t.Errorf("expected ID 1, got %v", created.ID)
	}
}

func TestGetAllCategories(t *testing.T) {
	mockService := &MockCategoryService{
		GetAllFunc: func() ([]model.Category, error) {
			return []model.Category{{ID: 1, Name: "A", Description: "B"}}, nil
		},
	}
	h := handler.NewCategoryHandler(mockService)

	req, err := http.NewRequest("GET", "/categories", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(h.HandleCategories)
	handlerFunc.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
