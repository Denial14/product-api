package test

import (
	"database/sql"
	_ "github.com/lib/pq"
	"testing"
)

func SetupTestDB(t *testing.T) *sql.DB {
	dsn := "postgres://postgres:postgres@localhost:5432/productdb_test?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Error connection to database: %v", err)
	}

	createReq := `CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    model VARCHAR(100) NOT NULL,
    company VARCHAR(100) NOT NULL,
    price INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err = db.Exec(createReq)
	if err != nil {
		t.Fatalf("Error creating database: %v", err)
	}

	_, err = db.Exec(`TRUNCATE TABLE products RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("Error cleaning table: %v", err)
	}
	return db
}

func CleanTestDB(db *sql.DB) {
	db.Exec(`TRUNCATE TABLE products RESTART IDENTITY CASCADE`)
}
