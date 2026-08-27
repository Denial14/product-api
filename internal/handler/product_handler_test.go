package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"test-product-api/internal/models"

	"github.com/gorilla/mux"
)

type mockService struct {
	products []models.Product
	nextID   int
	err      error
}

func (m *mockService) Create(p models.Product) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	if p.Model == "" || p.Company == "" || p.Price < 0 {
		return 0, fmt.Errorf("invalid product data")
	}
	m.nextID++
	p.ID = m.nextID
	m.products = append(m.products, p)
	return p.ID, nil
}

func (m *mockService) GetByID(id int) (*models.Product, error) {
	for _, p := range m.products {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("product not found")
}

func (m *mockService) GetAll() ([]models.Product, error) {
	return m.products, nil
}

func (m *mockService) Update(id int, req models.UpdateProductRequest) (*models.Product, error) {
	for i, p := range m.products {
		if p.ID == id {
			if req.Company != "" {
				m.products[i].Company = req.Company
			}
			if req.Model != "" {
				m.products[i].Model = req.Model
			}
			if req.Price >= 0 {
				m.products[i].Price = req.Price
			}
			return &m.products[i], nil
		}
	}
	return nil, fmt.Errorf("product not found")
}

func (m *mockService) Delete(id int) error {
	for i, p := range m.products {
		if p.ID == id {
			m.products = append(m.products[:i], m.products[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("product not found")
}

func setupTestHandler() *ProductHandler {
	svc := &mockService{products: []models.Product{}, nextID: 0}
	return NewProductHandler(svc)
}

func TestProductHandler_Create(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name       string
		payload    interface{}
		wantStatus int
	}{
		{
			name:       "great creating",
			payload:    models.CreateProductRequest{Model: "iPhone", Company: "Apple", Price: 1200},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "empty model",
			payload:    models.CreateProductRequest{Model: "", Company: "Apple", Price: 990},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid JSON",
			payload:    "invalid json",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler.Create(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestProductHandler_GetAll(t *testing.T) {
	handler := setupTestHandler()

	handler.service.Create(models.Product{Model: "A", Company: "CoA", Price: 100})
	handler.service.Create(models.Product{Model: "B", Company: "CoB", Price: 200})

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var products []models.Product
	json.NewDecoder(w.Body).Decode(&products)

	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}
}

func TestProductHandler_GetByID(t *testing.T) {
	handler := setupTestHandler()

	handler.service.Create(models.Product{Model: "Test", Company: "TestCo", Price: 100})

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{"existing ID", "1", http.StatusOK},
		{"non-existent ID", "999", http.StatusNotFound},
		{"invalid ID", "abc", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			w := httptest.NewRecorder()

			handler.GetByID(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestProductHandler_Update(t *testing.T) {
	handler := setupTestHandler()

	handler.service.Create(models.Product{Model: "Old", Company: "OldCo", Price: 100})

	tests := []struct {
		name       string
		id         string
		payload    interface{}
		wantStatus int
	}{
		{
			name:       "successful update",
			id:         "1",
			payload:    models.UpdateProductRequest{Price: 999},
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existent ID",
			id:         "999",
			payload:    models.UpdateProductRequest{Price: 999},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid ID",
			id:         "abc",
			payload:    models.UpdateProductRequest{Price: 999},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPut, "/products/"+tt.id, bytes.NewReader(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			w := httptest.NewRecorder()

			handler.Update(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestProductHandler_Delete(t *testing.T) {
	handler := setupTestHandler()

	handler.service.Create(models.Product{Model: "ToDelete", Company: "DelCo", Price: 100})

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{"successful delete", "1", http.StatusNoContent},
		{"non-existent ID", "999", http.StatusNotFound},
		{"invalid ID", "abc", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/products/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			w := httptest.NewRecorder()

			handler.Delete(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
