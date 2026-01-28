package services

import (
	"kasir-api-irfan/models"
	"kasir-api-irfan/repositories"
)

type CategoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *CategoryService) Create(category *models.Category) error {
	return s.categoryRepo.Create(category)
}

func (s *CategoryService) Update(category *models.Category) error {
	return s.categoryRepo.Update(category)
}

func (s *CategoryService) Delete(id int) error {
	return s.categoryRepo.Delete(id)
}
