package request

type UpdateUserRequest struct {
	KeyID           string `json:"key_id"`
	Username        string `json:"username" binding:"omitempty,min=3,max=20"`
	Password        string `json:"password" binding:"omitempty"`
	ConfirmPassword string `json:"confirm_password" binding:"omitempty"`
	Role            string `json:"role" binding:"omitempty,oneof=admin user"`
	Email           string `json:"email" binding:"omitempty,email"`
	Status          int    `json:"status" binding:"omitempty"`
	Introduction    string `json:"introduction" binding:"omitempty"`
	Qr              string `json:"qr" binding:"omitempty"`
}

type CreateUserRequest struct {
	KeyID           string `json:"key_id"`
	Username        string `json:"username" binding:"required,min=3,max=20"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Role            string `json:"role" binding:"required,oneof=admin user"`
}

type UserListRequest struct {
	Query    string `form:"query"`
	Status   string `form:"status"`
	Page     int    `form:"page,required,default=1"`
	PageSize int    `form:"pageSize,default=10"`
}
