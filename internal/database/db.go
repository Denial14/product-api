package database

import (
	"database/sql"
	"fmt"
	"time"
)

type DB struct {
	Conn *sql.DB
}

func New(dsn string, maxConnections, maxIdleConnections int, connectionTimeout string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения бд: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Ошибка Ping: %w", err)
	}

	db.SetMaxOpenConns(maxConnections)
	db.SetMaxIdleConns(maxIdleConnections)

	timeout, err := time.ParseDuration(connectionTimeout)
	if err == nil {
		db.SetConnMaxLifetime(timeout)
	}

	return &DB{Conn: db}, nil
}
