package service

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/entity"
	"blogs/internal/repository"
	appErrors "blogs/pkg/errors"
)

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) ListCategories() ([]entity.Category, error) {
	return s.categoryRepo.ListCategories()
}

func (s *categoryService) GetCategoryList(req request.CategoryRequest) ([]entity.Category, int64, error) {
	return s.categoryRepo.GetCategoryList(req)
}

func (s *categoryService) CreateCategory(category *entity.Category) error {
	if err := s.categoryRepo.CreateCategory(category); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "创建分类失败", err)
	}
	return nil
}

func (s *categoryService) DeleteCategory(id int) error {
	if err := s.categoryRepo.DeleteCategory(id); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "删除分类失败", err)
	}
	return nil
}

func (s *categoryService) UpdateCategory(category *entity.Category) error {
	if err := s.categoryRepo.UpdateCategory(category); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "更新分类失败", err)
	}
	return nil
}

func (s *categoryService) GetCategoryByID(id int) (*entity.Category, error) {
	return s.categoryRepo.GetCategoryByID(id)
}
