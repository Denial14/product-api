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

	if err = db.Ping(); err != nil {
		t.Fatalf("Error pinging test database: %v", err)
	}

	return db
}

func SetupTestDBProducts(t *testing.T) (*sql.DB, func()) {
	db := SetupTestDB(t)

	createReq := `CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    model VARCHAR(100) NOT NULL,
    company VARCHAR(100) NOT NULL,
    price INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(createReq)
	if err != nil {
		t.Fatalf("Error creating products table: %v", err)
	}

	_, err = db.Exec(`TRUNCATE TABLE products RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("Error cleaning products table: %v", err)
	}

	cleanup := func() {
		db.Exec(`TRUNCATE TABLE products RESTART IDENTITY CASCADE`)
		db.Close()
	}

	return db, cleanup
}

func SetupTestDBAll(t *testing.T) (*sql.DB, func()) {
	db := SetupTestDB(t)

	createProducts := `CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    model VARCHAR(100) NOT NULL,
    company VARCHAR(100) NOT NULL,
    price INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(createProducts)
	if err != nil {
		t.Fatalf("Error creating products table: %v", err)
	}

	createUsers := `CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err = db.Exec(createUsers)
	if err != nil {
		t.Fatalf("Error creating users table: %v", err)
	}

	db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)

	db.Exec(`TRUNCATE TABLE products, users RESTART IDENTITY CASCADE`)

	cleanup := func() {
		db.Exec(`TRUNCATE TABLE products, users RESTART IDENTITY CASCADE`)
		db.Close()
	}

	return db, cleanup
}

func SetupTestDBUsers(t *testing.T) (*sql.DB, func()) {
	db := SetupTestDB(t)

	createReq := `CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(createReq)
	if err != nil {
		t.Fatalf("Error creating users table: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
	if err != nil {
		t.Fatalf("Error creating index: %v", err)
	}

	_, err = db.Exec(`TRUNCATE TABLE users RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("Error cleaning users table: %v", err)
	}

	cleanup := func() {
		db.Exec(`TRUNCATE TABLE users RESTART IDENTITY CASCADE`)
		db.Close()
	}

	return db, cleanup
}

func CleanTestDB(db *sql.DB, tables ...string) {
	if len(tables) == 0 {
		db.Exec(`TRUNCATE TABLE products, users RESTART IDENTITY CASCADE`)
		return
	}

	for _, table := range tables {
		db.Exec(`TRUNCATE TABLE ` + table + ` RESTART IDENTITY CASCADE`)
	}
}
