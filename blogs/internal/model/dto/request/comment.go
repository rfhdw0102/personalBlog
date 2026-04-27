package request

type CommentSubmit struct {
	ParentID  int    `json:"parent_id"`
	ArticleID int    `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}
