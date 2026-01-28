package handler

import (
	"encoding/json"
	"kasir-api/internal/model"
	"kasir-api/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.URL.Path != "/categories" {
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

func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
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
// @Summary Get all categories
// @Description Get all categories
// @Tags categories
// @Produce json
// @Success 200 {array} model.Category
// @Failure 500 {object} map[string]string
// @Router /categories [get]
func (h *CategoryHandler) getAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

// create godoc
// @Summary Create a new category
// @Description Create a new category with the provided information
// @Tags categories
// @Accept json
// @Produce json
// @Param category body model.Category true "Category object"
// @Success 201 {object} model.Category
// @Failure 400 {object} map[string]string
// @Router /categories [post]
func (h *CategoryHandler) create(w http.ResponseWriter, r *http.Request) {
	var cat model.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	created, err := h.service.Create(cat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// getByID godoc
// @Summary Get category by ID
// @Description Get a single category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.Category
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [get]
func (h *CategoryHandler) getByID(w http.ResponseWriter, r *http.Request, id int) {
	cat, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(cat)
}

// update godoc
// @Summary Update a category
// @Description Update an existing category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body model.Category true "Category object"
// @Success 200 {object} model.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [put]
func (h *CategoryHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var cat model.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	updated, err := h.service.Update(id, cat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

// delete godoc
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Param id path int true "Category ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [delete]
func (h *CategoryHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
