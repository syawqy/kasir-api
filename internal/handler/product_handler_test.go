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

func TestCreateProduct(t *testing.T) {
	mockService := &MockProductService{
		CreateFunc: func(product model.Product) (model.Product, error) {
			product.ID = 100
			return product, nil
		},
	}
	h := handler.NewProductHandler(mockService)

	payload := []byte(`{"name":"Smartphone", "price":5000000, "stock":10, "category_id":1}`)
	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(h.HandleProducts)
	handlerFunc.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var created model.Product
	err = json.Unmarshal(rr.Body.Bytes(), &created)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.ID != 100 {
		t.Errorf("expected ID 100, got %v", created.ID)
	}
}
