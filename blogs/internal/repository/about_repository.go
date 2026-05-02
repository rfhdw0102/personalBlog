package repository

import (
	"blogs/internal/model/entity"
	"errors"

	"gorm.io/gorm"
)

type aboutRepository struct {
	db *gorm.DB
}

func NewAboutRepository(db *gorm.DB) AboutRepository {
	return &aboutRepository{db: db}
}

func (r *aboutRepository) Get() (*entity.About, error) {
	var about entity.About
	err := r.db.First(&about).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.About{}, nil
	}
	if err != nil {
		return nil, err
	}
	return &about, nil
}

func (r *aboutRepository) Upsert(about *entity.About) error {
	var existing entity.About
	err := r.db.First(&existing).Error
	if err == nil {
		about.ID = existing.ID
		return r.db.Save(about).Error
	}
	return r.db.Create(about).Error
}
