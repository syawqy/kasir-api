package handler

import (
	"encoding/json"
	"kasir-api/internal/model"
	"kasir-api/internal/service"
	"net/http"
	"strings"
	"time"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/report") {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getReport(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleCheckout godoc
// @Summary Process checkout/transaction
// @Description Create a new transaction with multiple items, calculate total, update stock
// @Tags transactions
// @Accept json
// @Produce json
// @Param items body model.CheckoutRequest true "Checkout items"
// @Success 201 {object} model.Transaction
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /checkout [post]
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/checkout" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodPost:
		h.checkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) checkout(w http.ResponseWriter, r *http.Request) {
	var req model.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		http.Error(w, "Items cannot be empty", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// getReport godoc
// @Summary Get sales report
// @Description Get sales report for a date range. Use /report/hari-ini for today's report.
// @Tags reports
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} repository.SalesReport
// @Failure 500 {object} map[string]string
// @Router /report [get]
func (h *TransactionHandler) getReport(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/report")

	var startDate, endDate string
	today := time.Now().Format("2006-01-02")

	if path == "/hari-ini" {
		// Today's report
		startDate = today
		endDate = today
	} else {
		// Date range from query params
		startDate = r.URL.Query().Get("start_date")
		endDate = r.URL.Query().Get("end_date")

		// Default to today if not provided
		if startDate == "" {
			startDate = today
		}
		if endDate == "" {
			endDate = today
		}
	}

	report, err := h.service.GetSalesReport(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
