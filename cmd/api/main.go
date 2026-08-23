package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"test-product-api/config"
	"test-product-api/internal/database"
	"test-product-api/internal/handler"
	"test-product-api/internal/middleware"
	"test-product-api/internal/repository"
	"test-product-api/internal/service"
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
	log.Printf("Сервер запущен на http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, wrapped))
}
