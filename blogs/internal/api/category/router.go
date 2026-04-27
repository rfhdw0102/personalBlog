package category

import (
	"blogs/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/category")
	// 获取全部分类
	r.GET("/category/list", ctrl.GetCategories)
	router.Use(middleware.JWTAuthMiddleware(ctrl.redisRepo))
	{
		// 获取分类列表
		router.GET("/pageList", middleware.AdminAuthMiddleware(), ctrl.GetCategory)
		// 获取单个分类
		router.GET("/:id", middleware.AdminAuthMiddleware(), ctrl.GetCategoryOne)
		// 新建分类
		router.POST("", middleware.AdminAuthMiddleware(), ctrl.CreateCategory)
		// 删除分类
		router.DELETE("/:id", middleware.AdminAuthMiddleware(), ctrl.DeleteCategory)
		// 编辑分类
		router.PUT("/:id", middleware.AdminAuthMiddleware(), ctrl.UpdateCategory)

	}
}
