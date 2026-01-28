package repository

import (
	"database/sql"
	"kasir-api/internal/model"
)

type CategoryRepository interface {
	Create(category model.Category) (model.Category, error)
	GetAll() ([]model.Category, error)
	GetByID(id int) (model.Category, error)
	Update(id int, category model.Category) (model.Category, error)
	Delete(id int) error
}

type postgresCategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &postgresCategoryRepository{db: db}
}

func (r *postgresCategoryRepository) Create(category model.Category) (model.Category, error) {
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	if err != nil {
		return model.Category{}, err
	}
	return category, nil
}

func (r *postgresCategoryRepository) GetAll() ([]model.Category, error) {
	rows, err := r.db.Query(`SELECT id, name, description FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *postgresCategoryRepository) GetByID(id int) (model.Category, error) {
	var c model.Category
	err := r.db.QueryRow(`SELECT id, name, description FROM categories WHERE id = $1`, id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return model.Category{}, err
	}
	return c, nil
}

func (r *postgresCategoryRepository) Update(id int, category model.Category) (model.Category, error) {
	query := `UPDATE categories SET name = $1, description = $2 WHERE id = $3 RETURNING id, name, description`
	var updated model.Category
	err := r.db.QueryRow(query, category.Name, category.Description, id).Scan(&updated.ID, &updated.Name, &updated.Description)
	if err != nil {
		return model.Category{}, err
	}
	return updated, nil
}

func (r *postgresCategoryRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
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
