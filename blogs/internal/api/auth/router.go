package auth

import (
	"blogs/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/auth")
	{
		// 注册
		router.POST("/register", ctrl.Register)
		// 登录
		router.POST("/login", ctrl.Login)
		// 获取验证码
		router.POST("/code", ctrl.SendCode)
		// 忘记密码
		router.POST("/password", ctrl.ForgotPassword)
		// 获取RSA公钥
		router.GET("/key", ctrl.GetPubKeyHandler)
		// 退出登录
		router.POST("/logout", middleware.JWTAuthMiddleware(ctrl.redisRepo), ctrl.Logout)
		//图形验证码
		router.GET("/captcha", ctrl.GetCaptcha)
	}
}
