package repository

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
)

type ArticleRepository interface {
	List(req request.ArticleListQuery) ([]response.ArticleInfo, int64, error)
	Create(article *entity.Article, tagIDs []int, tagNames []string) error
	Update(article *entity.Article, tagIDs []int, tagNames []string) error
	Delete(id int) error
	GetByID(id int) (*response.ArticleInfo, error)
	LikeArticle(articleID int, userID int) error
	UnlikeArticle(articleID int, userID int) error
	UpdateViewAndLike(articleID int, view int, like int) error
	IsLiked(articleID int, userID int) (bool, error)
}
