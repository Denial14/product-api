package service

import (
	"fmt"
	"test-product-api/internal/models"
	"testing"
)

type mockRepo struct {
	products []models.Product
	nextID   int
	err      error
}

func (m *mockRepo) Create(p models.Product) (int, error) {
	if m.err != nil {
		return 0, m.err
	}

	m.nextID++
	p.ID = m.nextID
	m.products = append(m.products, p)
	return p.ID, nil
}

func (m *mockRepo) GetByID(id int) (*models.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, p := range m.products {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("Не найден")
}

func (m *mockRepo) GetAll() ([]models.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.products, nil
}

func (m *mockRepo) Update(p models.Product) error {
	if m.err != nil {
		return m.err
	}
	for i, prod := range m.products {
		if prod.ID == p.ID {
			m.products[i] = p
			return nil
		}
	}
	return fmt.Errorf("Не найден")
}

func (m *mockRepo) Delete(id int) error {
	if m.err != nil {
		return m.err
	}
	for _, p := range m.products {
		if p.ID == id {
			m.products = append(m.products[:id], m.products[id+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Не найден")
}

func TestProductService_Create(t *testing.T) {
	repo := &mockRepo{products: []models.Product{}, nextID: 0}
	svc := NewProductService(repo)

	tests := []struct {
		name    string
		product models.Product
		wantErr bool
	}{
		{
			name:    "Успешное создание",
			product: models.Product{Model: "Test", Company: "TestCo", Price: 999},
			wantErr: false,
		},
		{
			name:    "Missing company",
			product: models.Product{Model: "Test", Company: "", Price: 100},
			wantErr: true,
		},
		{
			name:    "Missing model",
			product: models.Product{Model: "", Company: "Test", Price: 10},
			wantErr: true,
		},
		{
			name:    "Negative price",
			product: models.Product{Model: "Test", Company: "TestCo", Price: -777},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(tt.product)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error, but receive nothing")
			} else if !tt.wantErr && err != nil {
				t.Errorf("No expected error, but got: %v", err)
			}
		})
	}
}

func TestProductService_GetByID(t *testing.T) {
	repo := &mockRepo{products: []models.Product{}, nextID: 0}
	svc := NewProductService(repo)

	product := models.Product{Model: "Test", Company: "TestCo", Price: 999}
	id, _ := svc.Create(product)

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "Right ID",
			id:      id,
			wantErr: false,
		},
		{
			name:    "ID not exist",
			id:      999,
			wantErr: true,
		},
		{
			name:    "Negative ID",
			id:      -3,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GetByID(tt.id)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error, but receive nothing")
			} else if !tt.wantErr && err != nil {
				t.Errorf("No expected error, but got: %v", err)
			}
		})
	}
}

func TestProductService_Update(t *testing.T) {
	repo := &mockRepo{products: []models.Product{}, nextID: 0}
	svc := NewProductService(repo)

	product := models.Product{Model: "Test", Company: "TestCo", Price: 100}
	id, _ := svc.Create(product)

	tests := []struct {
		name    string
		id      int
		req     models.UpdateProductRequest
		wantErr bool
	}{
		{
			name:    "успешное обновление",
			id:      id,
			req:     models.UpdateProductRequest{Price: 999},
			wantErr: false,
		},
		{
			name:    "несуществующий ID",
			id:      999,
			req:     models.UpdateProductRequest{Price: 999},
			wantErr: true,
		},
		{
			name:    "отрицательный ID",
			id:      -1,
			req:     models.UpdateProductRequest{Price: 999},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Update(tt.id, tt.req)
			if tt.wantErr && err == nil {
				t.Errorf("ожидалась ошибка, но её нет")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("не ожидалась ошибка, но получили: %v", err)
			}
		})
	}
}

func TeTestProductService_Delete(t *testing.T) {
	repo := &mockRepo{products: []models.Product{}, nextID: 0}
	svc := NewProductService(repo)

	prod := models.Product{Model: "Test", Company: "TestCo", Price: 333}
	id, _ := svc.Create(prod)

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "успешное удаление",
			id:      id,
			wantErr: false,
		},
		{
			name:    "несуществующий ID",
			id:      999,
			wantErr: true,
		},
		{
			name:    "отрицательный ID",
			id:      -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Delete(tt.id)
			if tt.wantErr && err == nil {
				t.Errorf("ожидалась ошибка, но её нет")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("не ожидалась ошибка, но получили: %v", err)
			}
		})
	}
}
