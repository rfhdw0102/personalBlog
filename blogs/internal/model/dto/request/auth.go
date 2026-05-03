package request

type AuthRequest struct {
	KeyID           string `json:"key_id"`
	Username        string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Code            string `json:"code" binding:"required,len=6"`
}

type LoginRequest struct {
	KeyID     string `json:"key_id"`
	Username  string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Password  string `json:"password" binding:"required"`
	CaptchaID string `json:"captcha_id" binding:"required"`
	Captcha   string `json:"captcha" binding:"required,len=4"`
}

type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordRequest struct {
	KeyID           string `json:"key_id"`
	Email           string `json:"email" binding:"required,email"`
	Code            string `json:"code" binding:"required,len=6"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
