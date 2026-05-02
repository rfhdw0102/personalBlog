package response

import "blogs/internal/model/entity"

type AdjacentArticle struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type AdjacentArticles struct {
	Prev *AdjacentArticle `json:"prev"`
	Next *AdjacentArticle `json:"next"`
}

type ArticleInfo struct {
	ID           int          `json:"id"`
	Title        string       `json:"title"`
	Content      string       `json:"content"`
	Summary      string       `json:"summary"`
	UserID       int          `json:"user_id"`
	Username     string       `json:"username"`
	Avatar       string       `json:"avatar"`
	CoverImage   string       `json:"cover_image"`
	CategoryID   int          `json:"category_id"`
	CategoryName string       `json:"category_name"`
	Tags         []entity.Tag `json:"tags"`
	Status       string       `json:"status"`
	ViewCount    int          `json:"view_count"`
	LikeCount    int          `json:"like_count"`
	UpdatedAt    string       `json:"updated_at"`
}
