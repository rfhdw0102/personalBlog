package notification

import (
	"blogs/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/notification")
	router.Use(middleware.JWTAuthMiddleware(ctrl.redisRepo))
	{
		// 获取消息列表
		router.GET("/list", ctrl.List)
		// 已读某条消息
		router.PUT("/:id/read", ctrl.MarkAsRead)
		// 已读全部消息
		router.PUT("/read", ctrl.MarkAllAsRead)
		// 获取未读消息数量
		router.GET("/unread-count", ctrl.GetUnreadCount)
		// 删除某条消息
		router.DELETE("/:id", ctrl.Delete)
	}
}
