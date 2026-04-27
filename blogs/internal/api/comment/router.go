package comment

import (
	"blogs/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/comment")
	{
		// 获取某个文章所有评论
		router.GET("/article/:id", ctrl.GetCommentsByArticle)
		// 获取某个文章所有隐藏评论
		router.GET("/article/:id/hide", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.GetCommentsByArticleHide)
		// 提交评论
		router.POST("", middleware.JWTAuthMiddleware(ctrl.redisRepo), ctrl.AddComment)
		// 隐藏评论
		router.PUT("/:id/hide", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.HideComment)
		// 取消隐藏评论
		router.PUT("/:id/unhide", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.UnHideComment)
		// 获取某个用户所收到的全部评论
		router.GET("/user", middleware.JWTAuthMiddleware(ctrl.redisRepo), ctrl.ListUserCommentsGet)
		// 获取某个用户给出的所有评论
		router.GET("/userTake", middleware.JWTAuthMiddleware(ctrl.redisRepo), ctrl.ListUserCommentsTake)
		// 获取某个作者隐藏的评论
		router.GET("/user/hide", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.ListUserCommentsHide)
		// 删除评论
		router.DELETE("/:id", middleware.JWTAuthMiddleware(ctrl.redisRepo), ctrl.DeleteComment)
	}
}
