package request

type PageRequest struct {
	Page     int `form:"page" json:"page" binding:"required,min=1"`                 // 页码，从1开始
	PageSize int `form:"pageSize" json:"pageSize" binding:"required,min=1,max=100"` // 每页大小，最大100
}
