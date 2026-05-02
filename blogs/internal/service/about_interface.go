package service

import "blogs/internal/model/entity"

type AboutService interface {
	Get() (*entity.About, error)
	Update(about *entity.About) error
}
