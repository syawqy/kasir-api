package handler

import (
	"encoding/json"
	"kasir-api/internal/model"
	"kasir-api/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.URL.Path != "/products" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAll godoc
// @Summary Get all products
// @Description Get all products with their category information
// @Tags products
// @Produce json
// @Success 200 {array} model.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *ProductHandler) getAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}

// create godoc
// @Summary Create a new product
// @Description Create a new product with the provided information
// @Tags products
// @Accept json
// @Produce json
// @Param product body model.Product true "Product object"
// @Success 201 {object} model.Product
// @Failure 400 {object} map[string]string
// @Router /products [post]
func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	created, err := h.service.Create(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// getByID godoc
// @Summary Get product by ID
// @Description Get a single product by ID with category information
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} model.Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request, id int) {
	p, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}

// update godoc
// @Summary Update a product
// @Description Update an existing product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body model.Product true "Product object"
// @Success 200 {object} model.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [put]
func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	updated, err := h.service.Update(id, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

// delete godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
