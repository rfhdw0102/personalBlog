package service

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/entity"
)

type TagService interface {
	ListTags() ([]entity.Tag, error)
	GetTagList(req request.TagRequest) ([]entity.Tag, int64, error)
	CreateTag(tag entity.Tag) error
	DeleteTag(id int) error
	UpdateTag(tag entity.Tag) error
	GetTagByID(id int) (entity.Tag, error)
}
