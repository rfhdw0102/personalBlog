package service

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/entity"
	"blogs/internal/repository"
	appErrors "blogs/pkg/errors"
)

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}

func (s *tagService) ListTags() ([]entity.Tag, error) {
	return s.tagRepo.ListTags()
}

func (s *tagService) GetTagList(req request.TagRequest) ([]entity.Tag, int64, error) {
	return s.tagRepo.GetTagList(req)
}

func (s *tagService) CreateTag(tag entity.Tag) error {
	if err := s.tagRepo.CreateTag(tag); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "创建标签失败", err)
	}
	return nil
}

func (s *tagService) DeleteTag(id int) error {
	if err := s.tagRepo.DeleteTag(id); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "删除标签失败", err)
	}
	return nil
}

func (s *tagService) UpdateTag(tag entity.Tag) error {
	if err := s.tagRepo.UpdateTag(tag); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "更新标签失败", err)
	}
	return nil
}

func (s *tagService) GetTagByID(id int) (entity.Tag, error) {
	return s.tagRepo.GetTagByID(id)
}
