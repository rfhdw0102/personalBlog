package response

type CommentInfo struct {
	ID             int            `json:"id"`
	ArticleID      int            `json:"article_id"`
	UserID         int            `json:"user_id"`
	Username       string         `json:"username"`
	Avatar         string         `json:"avatar"`
	ArticleTitle   string         `json:"article_title"`
	Content        string         `json:"content"`
	ParentID       int            `json:"parent_id"`
	ParentUsername string         `json:"parent_username"`
	CreatedAt      string         `json:"created_at"`
	Children       []*CommentInfo `json:"children,omitempty" gorm:"-"`
}
