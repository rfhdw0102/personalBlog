package response

import (
	"time"
)

type LoginUserResponse struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role"`
	Status    int       `json:"status"`
	Token     string    `json:"token"`
}

type AuthorResponse struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Avatar       string `json:"avatar"`
	Introduction string `json:"introduction"`
	Qr           string `json:"qr"`
	ArticleCount int    `json:"article_count"`
	ViewCount    int    `json:"view_count"`
}

type AdminCard struct {
	UserCount     int `json:"user_count"`
	ArticleCount  int `json:"article_count"`
	CategoryCount int `json:"category_count"`
	TagCount      int `json:"tag_count"`
}
