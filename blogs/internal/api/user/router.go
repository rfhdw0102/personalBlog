package user

import (
	"blogs/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/user")
	{
		// 获取个人信息
		router.GET("/info", middleware.JWTAuthMiddleware(ctrl.redisRepo), ctrl.GetInfo)
		// 获取博客作者信息
		router.GET("/author", ctrl.GetAuthor)
		// 编辑用户信息
		router.PUT("", middleware.JWTAuthMiddleware(ctrl.redisRepo), ctrl.Update)
		// 上传头像
		router.POST("/avatar", middleware.JWTAuthMiddleware(ctrl.redisRepo), ctrl.UploadAvatar)
	}
	// 管理员权限路由
	admin := r.Group("/admin")
	admin.Use(middleware.JWTAuthMiddleware(ctrl.redisRepo))
	admin.Use(middleware.AdminAuthMiddleware())
	{
		// 上传二维码
		admin.POST("/qr", ctrl.PostQr)
		// 获取全部用户信息
		admin.GET("/list", ctrl.List)
		// 创建用户
		admin.POST("", ctrl.Create)
		// 删除用户
		admin.DELETE("/:id", ctrl.Delete)
		// 更新用户
		admin.PUT("/:id", ctrl.AdminUpdate)
		// 获取卡片数据
		admin.GET("/stats", ctrl.GetStats)
	}
}
