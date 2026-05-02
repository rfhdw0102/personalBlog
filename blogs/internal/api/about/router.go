package about

import (
	"blogs/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/about")
	{
		router.GET("", ctrl.Get)
		router.PUT("", middleware.JWTAuthMiddleware(ctrl.redisRepo), middleware.AdminAuthMiddleware(), ctrl.Update)
	}
}
