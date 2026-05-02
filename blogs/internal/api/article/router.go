package article

import (
	"blogs/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/article")
	{
		// 获取文章列表
		router.GET("/list", ctrl.List)
		// 获取单个文章
		router.GET("/:id", ctrl.Get)
		// 增加文章浏览量
		router.POST("/:id/view", ctrl.IncrementView)
		// 获得是否已点赞
		router.GET("/:id/like-status", ctrl.CheckLikeStatus)
		// 获取上一篇文章和下一篇文章
		router.GET("/:id/adjacent", ctrl.GetAdjacent)
		// 创建文章
		router.POST("", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.Create)
		// 编辑文章
		router.PUT("", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.Update)
		// 删除文章
		router.DELETE("/:id", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.Delete)
		// 点赞
		router.POST("/like/:id", middleware.JWTAuthMiddleware(ctrl.redisRepo), ctrl.LikeArticle)
		// 取消点赞
		router.POST("/unlike/:id", middleware.JWTAuthMiddleware(ctrl.redisRepo), ctrl.UnlikeArticle)
		// 上传封面
		router.POST("/coverImage", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.UploadCoverImage)
		// 上传内容图片
		router.POST("/contentImage", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.UploadContentImage)
	}
}
