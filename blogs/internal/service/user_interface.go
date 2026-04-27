package service

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
)

type UserService interface {
	GetByID(id int) (entity.User, error)
	GetAuthor() (response.AuthorResponse, error)
	Update(id int, req request.UpdateUserRequest, key string) error
	Delete(id int) error
	List(req request.UserListRequest) ([]entity.User, int64, error)
	UpdateAvatar(id int, avatarPath string) error
	Create(req request.CreateUserRequest, key string) error
	ConsumePrivateKey(keyID string) (string, error)
	GetStats() (response.AdminCard, error)
}
