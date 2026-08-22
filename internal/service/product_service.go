package service

import (
	"fmt"
	"test-product-api/internal/models"
	"test-product-api/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(p models.Product) (int, error) {
	if p.Model == "" {
		return 0, fmt.Errorf("Модель не может быть пустой")
	}
	if p.Company == "" {
		return 0, fmt.Errorf("Компания не может быть пустой")
	}
	if p.Price < 0 {
		return 0, fmt.Errorf("Цена не может быть пустой")
	}
	return s.repo.Create(p)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID должен быть положительным")
	}
	return s.repo.GetByID(id)
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) Update(id int, req models.UpdateProductRequest) (*models.Product, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID должен быть положительным")
	}

	// Проверяем, существует ли товар
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Сохраняем старые значения
	updated := *existing

	// Обновляем только те поля, которые пришли
	if req.Model != "" {
		updated.Model = req.Model
	}
	if req.Company != "" {
		updated.Company = req.Company
	}
	if req.Price >= 0 {
		updated.Price = req.Price
	}

	if err = s.repo.Update(updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *ProductService) Delete(id int) error {
	if id <= 0 {
		return fmt.Errorf("ID должен быть положительным")
	}
	return s.repo.Delete(id)
}
