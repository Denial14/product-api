package repository

import (
	"database/sql"
	"fmt"
	"test-product-api/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user models.User) error {
	query := `INSERT INTO users (email, password_hash, full_name) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, user.Email, user.PasswordHash, user.FullName)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password_hash, full_name, created_at FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, full_name, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.FullName, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
