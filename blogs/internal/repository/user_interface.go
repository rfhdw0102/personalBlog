package repository

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
)

type UserRepository interface {
	Create(user entity.User) error
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	GetByUsername(username string) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	UpdatePassword(email, password string) error
	GetByID(id int) (entity.User, error)
	GetAuthor() (response.AuthorResponse, error)
	Updates(id int, values map[string]interface{}) error
	Delete(id int) error
	List(req request.UserListRequest) ([]entity.User, int64, error)
	GetStats() (response.AdminCard, error)
}
