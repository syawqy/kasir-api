package service

import (
	"kasir-api/internal/model"
	"kasir-api/internal/repository"
)

type TransactionService interface {
	Checkout(items []model.CheckoutItem) (*model.Transaction, error)
	GetSalesReport(startDate, endDate string) (*repository.SalesReport, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) Checkout(items []model.CheckoutItem) (*model.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *transactionService) GetSalesReport(startDate, endDate string) (*repository.SalesReport, error) {
	return s.repo.GetSalesReport(startDate, endDate)
}
