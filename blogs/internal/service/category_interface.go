package service

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/entity"
)

type CategoryService interface {
	ListCategories() ([]entity.Category, error)
	GetCategoryList(req request.CategoryRequest) ([]entity.Category, int64, error)
	CreateCategory(category *entity.Category) error
	DeleteCategory(id int) error
	UpdateCategory(category *entity.Category) error
	GetCategoryByID(id int) (*entity.Category, error)
}
