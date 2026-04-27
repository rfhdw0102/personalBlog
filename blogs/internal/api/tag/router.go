package tag

import (
	"blogs/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/tag")
	router.Use()
	{
		// 获得所有标签
		router.GET("/list", ctrl.GetTags)
		// 获取分页标签
		router.GET("/pageList", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.GetTag)
		// 获取单个标签
		router.GET("/:id", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.GetTagOne)
		// 新建标签
		router.POST("", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.CreateTag)
		// 删除标签
		router.DELETE("/:id", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.DeleteTag)
		// 编辑标签
		router.PUT("/:id", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.UpdateTag)
	}
}
