package repository

import (
	"database/sql"
	"kasir-api/internal/model"
)

type ProductRepository interface {
	Create(product model.Product) (model.Product, error)
	GetAll() ([]model.Product, error)
	GetByID(id int) (model.Product, error)
	Update(id int, product model.Product) (model.Product, error)
	Delete(id int) error
}

type postgresProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Create(product model.Product) (model.Product, error) {
	query := `INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		return model.Product{}, err
	}
	return product, nil
}

func (r *postgresProductRepository) GetAll() ([]model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.stock, p.category_id,
			c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		var catID sql.NullInt64
		var catName, catDesc sql.NullString

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &catID, &catName, &catDesc); err != nil {
			return nil, err
		}

		// Populate category if it exists
		if catID.Valid {
			p.Category = &model.Category{
				ID:          int(catID.Int64),
				Name:        catName.String,
				Description: catDesc.String,
			}
		}

		products = append(products, p)
	}
	return products, nil
}

func (r *postgresProductRepository) GetByID(id int) (model.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.stock, p.category_id,
			c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`
	var p model.Product
	var catID sql.NullInt64
	var catName, catDesc sql.NullString

	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &catID, &catName, &catDesc)
	if err != nil {
		return model.Product{}, err
	}

	// Populate category if it exists
	if catID.Valid {
		p.Category = &model.Category{
			ID:          int(catID.Int64),
			Name:        catName.String,
			Description: catDesc.String,
		}
	}

	return p, nil
}

func (r *postgresProductRepository) Update(id int, product model.Product) (model.Product, error) {
	query := `UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5 RETURNING id, name, price, stock, category_id`
	var updated model.Product
	err := r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID, id).Scan(&updated.ID, &updated.Name, &updated.Price, &updated.Stock, &updated.CategoryID)
	if err != nil {
		return model.Product{}, err
	}
	return updated, nil
}

func (r *postgresProductRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
