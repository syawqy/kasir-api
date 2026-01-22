package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// resetStore resets the global state for testing
func resetStore() {
	categoryMux.Lock()
	categories = []Category{}
	nextID = 1
	categoryMux.Unlock()
}

func TestCreateCategory(t *testing.T) {
	resetStore()

	payload := []byte(`{"name":"Electronics", "description":"Gadgets"}`)
	req, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(categoriesHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var created Category
	err = json.Unmarshal(rr.Body.Bytes(), &created)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.ID != 1 {
		t.Errorf("expected ID 1, got %v", created.ID)
	}
	if created.Name != "Electronics" {
		t.Errorf("expected Name Electronics, got %v", created.Name)
	}
}

func TestGetAllCategories(t *testing.T) {
	resetStore()

	// Seed one item
	categoryMux.Lock()
	categories = append(categories, Category{ID: 1, Name: "A", Description: "B"})
	categoryMux.Unlock()

	req, err := http.NewRequest("GET", "/categories", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(categoriesHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var list []Category
	err = json.Unmarshal(rr.Body.Bytes(), &list)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(list) != 1 {
		t.Errorf("expected 1 item, got %v", len(list))
	}
}

func TestGetCategoryByID(t *testing.T) {
	resetStore()
	categoryMux.Lock()
	categories = append(categories, Category{ID: 10, Name: "Test", Description: "Desc"})
	categoryMux.Unlock()

	req, err := http.NewRequest("GET", "/categories/10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(categoryByIDHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var cat Category
	err = json.Unmarshal(rr.Body.Bytes(), &cat)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if cat.ID != 10 {
		t.Errorf("expected ID 10, got %v", cat.ID)
	}
}

func TestUpdateCategory(t *testing.T) {
	resetStore()
	categoryMux.Lock()
	categories = append(categories, Category{ID: 2, Name: "Old", Description: "Old Desc"})
	categoryMux.Unlock()

	payload := []byte(`{"name":"New", "description":"New Desc"}`)
	req, err := http.NewRequest("PUT", "/categories/2", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(categoryByIDHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verify update in memory
	categoryMux.Lock()
	updated := categories[0]
	categoryMux.Unlock()

	if updated.Name != "New" {
		t.Errorf("expected Name New, got %v", updated.Name)
	}
}

func TestDeleteCategory(t *testing.T) {
	resetStore()
	categoryMux.Lock()
	categories = append(categories, Category{ID: 5, Name: "Del", Description: "To Delete"})
	categoryMux.Unlock()

	req, err := http.NewRequest("DELETE", "/categories/5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(categoryByIDHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	categoryMux.Lock()
	count := len(categories)
	categoryMux.Unlock()

	if count != 0 {
		t.Errorf("expected 0 items, got %v", count)
	}
}
