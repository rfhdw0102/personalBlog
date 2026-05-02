package repository

import "blogs/internal/model/entity"

type AboutRepository interface {
	Get() (*entity.About, error)
	Upsert(about *entity.About) error
}
