package repository

import (
	"database/sql"
	"fmt"
	"test-product-api/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(p models.Product) (int, error) {
	var id int
	query := `INSERT INTO products (model, company, price) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, p.Model, p.Company, p.Price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Ошибка создания продукта: %w", err)
	}
	return id, nil
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	var p models.Product

	query := `SELECT id, model, company, price FROM products WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Model, &p.Company, &p.Price)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Товар с id = %d не найден", id)
	}

	if err != nil {
		return nil, fmt.Errorf("Ошибка получения продукта: %w", err)
	}

	return &p, nil
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	query := `SELECT id, model, company, price FROM products ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения списка продуктов: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err = rows.Scan(&p.ID, &p.Model, &p.Company, &p.Price); err != nil {
			return nil, fmt.Errorf("Ошибка сканирования: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) Update(p models.Product) error {
	query := `UPDATE products SET model = $1, company = $2, price = $3 WHERE id = $4`
	result, err := r.db.Exec(query, p.Model, p.Company, p.Price, p.ID)

	if err != nil {
		return fmt.Errorf("Ошибка обновления продукта: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("Товар с id = %d не найден", p.ID)
	}
	return nil
}

func (r *ProductRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Ошибка удаления продукта: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("Товар с id = %d не найден", id)
	}
	return nil
}
