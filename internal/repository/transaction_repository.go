package repository

import (
	"database/sql"
	"fmt"
	"kasir-api/internal/model"
	"strings"

	"github.com/lib/pq"
)

type SalesReport struct {
	TotalRevenue   int            `json:"total_revenue"`
	TotalTransaksi int            `json:"total_transaksi"`
	ProdukTerlaris ProdukTerlaris `json:"produk_terlaris"`
}

type ProdukTerlaris struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

type TransactionRepository interface {
	CreateTransaction(items []model.CheckoutItem) (*model.Transaction, error)
	GetSalesReport(startDate, endDate string) (*SalesReport, error)
}

type postgresTransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &postgresTransactionRepository{db: db}
}

func (r *postgresTransactionRepository) GetSalesReport(startDate, endDate string) (*SalesReport, error) {
	// Get total revenue and total transactions
	var totalRevenue, totalTransaksi int
	query := `
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) BETWEEN $1 AND $2
	`
	err := r.db.QueryRow(query, startDate, endDate).Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return nil, err
	}

	// Get best-selling product
	var produkNama string
	var qtyTerjual int
	query = `
		SELECT p.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) BETWEEN $1 AND $2
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`
	err = r.db.QueryRow(query, startDate, endDate).Scan(&produkNama, &qtyTerjual)
	if err == sql.ErrNoRows {
		produkNama = ""
		qtyTerjual = 0
	} else if err != nil {
		return nil, err
	}

	return &SalesReport{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: ProdukTerlaris{
			Nama:       produkNama,
			QtyTerjual: qtyTerjual,
		},
	}, nil
}

func (r *postgresTransactionRepository) CreateTransaction(items []model.CheckoutItem) (*model.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]model.TransactionDetail, 0, len(items))

	if len(items) == 0 {
		var transactionID int
		err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
		if err != nil {
			return nil, err
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}

		return &model.Transaction{
			ID:          transactionID,
			TotalAmount: totalAmount,
			Details:     details,
		}, nil
	}

	type productRow struct {
		ID    int
		Name  string
		Price int
		Stock int
	}

	qtyByID := make(map[int]int, len(items))
	uniqueIDs := make([]int, 0, len(items))
	for _, item := range items {
		if _, ok := qtyByID[item.ProductID]; !ok {
			uniqueIDs = append(uniqueIDs, item.ProductID)
		}
		qtyByID[item.ProductID] += item.Quantity
	}

	rows, err := tx.Query("SELECT id, name, price, stock FROM products WHERE id = ANY($1::int[]) FOR UPDATE", pq.Array(uniqueIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make(map[int]productRow, len(uniqueIDs))
	for rows.Next() {
		var p productRow
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products[p.ID] = p
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, item := range items {
		if _, ok := products[item.ProductID]; !ok {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
	}

	for productID, qty := range qtyByID {
		p := products[productID]
		if p.Stock < qty {
			return nil, fmt.Errorf("insufficient stock for product %s (available: %d, requested: %d)", p.Name, p.Stock, qty)
		}
	}

	updateArgs := make([]interface{}, 0, len(uniqueIDs)*2)
	var updateQuery strings.Builder
	updateQuery.WriteString("UPDATE products SET stock = stock - v.qty FROM (VALUES ")
	argPos := 1
	for i, productID := range uniqueIDs {
		if i > 0 {
			updateQuery.WriteString(",")
		}
		updateQuery.WriteString(fmt.Sprintf("($%d::int, $%d::int)", argPos, argPos+1))
		updateArgs = append(updateArgs, productID, qtyByID[productID])
		argPos += 2
	}
	updateQuery.WriteString(") AS v(id, qty) WHERE products.id = v.id")

	_, err = tx.Exec(updateQuery.String(), updateArgs...)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		p := products[item.ProductID]
		subtotal := p.Price * item.Quantity
		totalAmount += subtotal
		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: p.Name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	if len(details) > 0 {
		insertArgs := make([]interface{}, 0, len(details)*4)
		var insertQuery strings.Builder
		insertQuery.WriteString("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ")
		argPos = 1
		for i := range details {
			if i > 0 {
				insertQuery.WriteString(",")
			}
			insertQuery.WriteString(fmt.Sprintf("($%d::int, $%d::int, $%d::int, $%d::int)", argPos, argPos+1, argPos+2, argPos+3))
			insertArgs = append(insertArgs, transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
			argPos += 4
		}

		_, err = tx.Exec(insertQuery.String(), insertArgs...)
		if err != nil {
			return nil, err
		}
	}

	for i := range details {
		details[i].TransactionID = transactionID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
