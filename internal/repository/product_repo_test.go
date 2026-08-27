package repository

import (
	"test-product-api/internal/models"
	"test-product-api/internal/test"
	"testing"
)

func TestProductRepository_Create(t *testing.T) {
	db := test.SetupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	product := models.Product{Model: "Test model",
		Company: "Test company",
		Price:   100}

	id, err := repo.Create(product)
	if err != nil {
		t.Fatalf("Error repository.Create: %v", err)
	}

	if id <= 0 {
		t.Errorf("Invalid ID: %d", id)
	}

	saved, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("Error getting: %v", err)
	}

	if saved.Model != product.Model {
		t.Errorf("Model: expected %s, but got %s", product.Model, saved.Model)
	}

	if saved.Company != product.Company {
		t.Errorf("Company: expected %s, but got %s", product.Company, saved.Company)
	}

	if saved.Price != product.Price {
		t.Errorf("Price: expected %d, but got %d", product.Price, saved.Price)
	}
}

func TestProductRepository_GetByID(t *testing.T) {
	db := test.SetupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	product := models.Product{Model: "Test", Company: "Company", Price: 50}

	id, err := repo.Create(product)
	if err != nil {
		t.Fatalf("Error creating: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name: "Valid ID", id: id, wantErr: false,
		},
		{
			name: "Invalid ID", id: 999, wantErr: true,
		},
		{
			name: "Negative ID", id: -333, wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = repo.GetByID(tt.id)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error, but no got")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}

func TestProductRepository_GetAll(t *testing.T) {
	db := test.SetupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	products := []models.Product{
		{Model: "A", Company: "Co1", Price: 100},
		{Model: "B", Company: "Co2", Price: 200},
		{Model: "C", Company: "Co3", Price: 300},
	}

	for _, p := range products {
		_, err := repo.Create(p)
		if err != nil {
			t.Fatalf("Error creating: %v", err)
		}
	}

	result, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Error getting list of products: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 products, but got: %d", len(products))
	}
}

func TestProductRepository_Update(t *testing.T) {
	db := test.SetupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	product := models.Product{Model: "Old", Company: "OldCo", Price: 10}

	id, err := repo.Create(product)
	if err != nil {
		t.Fatalf("Error creating: %v", err)
	}

	product.Model = "New"
	product.Company = "NewCo"
	product.Price = 200
	product.ID = id

	err = repo.Update(product)
	if err != nil {
		t.Fatalf("Error update product: %v", err)
	}

	saved, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("Error getting product: %v", err)
	}

	if saved.Model != "New" {
		t.Errorf("Модель: ожидали New, получили %s", saved.Model)
	}
	if saved.Price != 200 {
		t.Errorf("Цена: ожидали 200, получили %d", saved.Price)
	}
}

func TestProductRepository_Delete(t *testing.T) {
	db := test.SetupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	product := models.Product{Model: "ToDelete", Company: "DelCo", Price: 100}
	id, err := repo.Create(product)
	if err != nil {
		t.Fatalf("Error creating: %v", err)
	}

	err = repo.Delete(id)
	if err != nil {
		t.Fatalf("Error delete: %v", err)
	}

	_, err = repo.GetByID(id)
	if err == nil {
		t.Errorf("Products has been deleted, but he is exist: %v", err)
	}
}
