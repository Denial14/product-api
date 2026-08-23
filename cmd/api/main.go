package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"test-product-api/config"
	"test-product-api/internal/database"
	"test-product-api/internal/handler"
	"test-product-api/internal/middleware"
	"test-product-api/internal/repository"
	"test-product-api/internal/service"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки конфига:", err)
	}

	db, err := database.New(cfg.Database.DSN(), cfg.Database.MaxConnections, cfg.Database.MaxIdleConnections, cfg.Database.ConnectionTimeout)
	if err != nil {
		log.Fatal("Ошибка подключение к бд:", err)
	}
	defer db.Close()

	log.Println("Подключено к БД:", cfg.Database.Name)

	productRepo := repository.NewProductRepository(db.Conn)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	r := mux.NewRouter()

	r.HandleFunc("/products", productHandler.Create).Methods("POST")
	r.HandleFunc("/products", productHandler.GetAll).Methods("GET")
	r.HandleFunc("/products/{id:[0-9]+}", productHandler.GetByID).Methods("GET")
	r.HandleFunc("/products/{id:[0-9]+}", productHandler.Update).Methods("PUT")
	r.HandleFunc("/products/{id:[0-9]+}", productHandler.Delete).Methods("DELETE")

	wrapped := middleware.Chain(r, middleware.Recovery, middleware.Logging)

	port := strconv.Itoa(cfg.App.Port)
	addr := ":" + port

	server := &http.Server{
		Addr:         addr,
		Handler:      wrapped,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Сервер запущен на http://localhost%s", addr)
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Ошибка сервера: %v", err)
		}
	}()

	quit := make(chan (os.Signal), 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Получен сигнал завершения...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = db.Close(); err != nil {
		log.Printf("Ошибка закрытия БД: %v", err)
	}

	if err = server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка остановки сервера: %v", err)
	}
	log.Println("Сервер остановлен корректно!")
}
