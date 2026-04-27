package request

type ArticleListQuery struct {
	Query      string `form:"query"`
	Status     string `form:"status"`
	CategoryID int    `form:"category_id"`
	TagID      int    `form:"tag_id"`
	Sort       int    `form:"sort"`
	Page       int    `form:"page,required,default=1"`
	PageSize   int    `form:"pageSize,default=10"`
}

type ArticleSubmit struct {
	Title      string   `json:"title" binding:"required"`
	Summary    string   `json:"summary"`
	Content    string   `json:"content" binding:"required"`
	CoverImage string   `json:"cover_image"`
	CategoryID int      `json:"category_id"`
	TagIDs     []int    `json:"tag_ids"`
	TagNames   []string `json:"tag_names"`
	Status     string   `json:"status" binding:"required,oneof=draft published"`
}

type ArticleEdit struct {
	ID         int      `json:"id" binding:"required"`
	Title      string   `json:"title" binding:"required"`
	Summary    string   `json:"summary"`
	Content    string   `json:"content" binding:"required"`
	CoverImage string   `json:"cover_image"`
	CategoryID int      `json:"category_id"`
	TagIDs     []int    `json:"tag_ids"`
	TagNames   []string `json:"tag_names"`
	Status     string   `json:"status" binding:"required,oneof=draft published"`
}
