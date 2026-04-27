package request

type TagRequest struct {
	Query    string `form:"query"`
	Page     int    `form:"page,required,default=1"`
	PageSize int    `form:"pageSize,default=10"`
}
