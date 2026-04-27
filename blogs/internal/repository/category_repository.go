package repository

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/entity"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) ListCategories() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetCategoryList(req request.CategoryRequest) ([]entity.Category, int64, error) {
	var categories []entity.Category
	var total int64
	db := r.db.Model(&entity.Category{})
	if req.Query != "" {
		db = db.Where("name LIKE ?", "%"+req.Query+"%")
	}
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&categories).Error
	return categories, total, err
}

func (r *categoryRepository) CreateCategory(category *entity.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) DeleteCategory(id int) error {
	category := &entity.Category{Base: entity.Base{ID: id}}
	return r.db.Delete(category).Error
}

func (r *categoryRepository) UpdateCategory(category *entity.Category) error {
	return r.db.Model(&entity.Category{}).Where("id = ?", category.ID).Updates(category).Error
}

func (r *categoryRepository) GetCategoryByID(id int) (*entity.Category, error) {
	var category entity.Category
	err := r.db.First(&category, id).Error
	return &category, err
}
