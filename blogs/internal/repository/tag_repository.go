package repository

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/entity"

	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) ListTags() ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}

func (r *tagRepository) GetTagList(req request.TagRequest) ([]entity.Tag, int64, error) {
	var tags []entity.Tag
	var total int64
	db := r.db.Model(&entity.Tag{})
	if req.Query != "" {
		db = db.Where("name LIKE ?", "%"+req.Query+"%")
	}
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&tags).Error
	return tags, total, err
}

func (r *tagRepository) CreateTag(tag entity.Tag) error {
	return r.db.Create(&tag).Error
}

func (r *tagRepository) DeleteTag(id int) error {
	tag := &entity.Tag{Base: entity.Base{ID: id}}
	return r.db.Delete(tag).Error
}

func (r *tagRepository) UpdateTag(tag entity.Tag) error {
	return r.db.Model(&entity.Tag{}).Where("id = ?", tag.ID).Updates(tag).Error
}

func (r *tagRepository) GetTagByID(id int) (entity.Tag, error) {
	var tag entity.Tag
	err := r.db.First(&tag, id).Error
	return tag, err
}
