package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Category model
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Global in-memory storage
var (
	categories  = []Category{}
	nextID      = 1
	categoryMux sync.Mutex // Mutex for thread-safe access
)

func main() {
	// Register handlers
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/categories/", categoryByIDHandler)

	fmt.Println("Server starting on port 8080...")
	// Start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

// categoriesHandler handles GET /categories and POST /categories
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Exact match check for /categories to avoid unintended matches if any
	if r.URL.Path != "/categories" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getAllCategories(w, r)
	case http.MethodPost:
		createCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// categoryByIDHandler handles GET, PUT, DELETE for /categories/{id}
func categoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	// Clean potentially trailing slash
	idStr = strings.TrimSuffix(idStr, "/") 
	
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getCategoryByID(w, r, id)
	case http.MethodPut:
		updateCategory(w, r, id)
	case http.MethodDelete:
		deleteCategory(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllCategories(w http.ResponseWriter, r *http.Request) {
	categoryMux.Lock()
	defer categoryMux.Unlock()

	json.NewEncoder(w).Encode(categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCat Category
	if err := json.NewDecoder(r.Body).Decode(&newCat); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	categoryMux.Lock()
	defer categoryMux.Unlock()

	newCat.ID = nextID
	nextID++
	categories = append(categories, newCat)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCat)
}

func getCategoryByID(w http.ResponseWriter, r *http.Request, id int) {
	categoryMux.Lock()
	defer categoryMux.Unlock()

	for _, cat := range categories {
		if cat.ID == id {
			json.NewEncoder(w).Encode(cat)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func updateCategory(w http.ResponseWriter, r *http.Request, id int) {
	var updatedData Category
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	categoryMux.Lock()
	defer categoryMux.Unlock()

	for i, cat := range categories {
		if cat.ID == id {
			// Update fields - keep ID same
			categories[i].Name = updatedData.Name
			categories[i].Description = updatedData.Description
			
			// Return the updated object
			updatedCat := categories[i]
			json.NewEncoder(w).Encode(updatedCat)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request, id int) {
	categoryMux.Lock()
	defer categoryMux.Unlock()

	for i, cat := range categories {
		if cat.ID == id {
			// Remove from slice
			categories = append(categories[:i], categories[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}
