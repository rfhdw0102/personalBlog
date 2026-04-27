package notification

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/repository"
	"blogs/internal/service"
	appErrors "blogs/pkg/errors"
	"blogs/pkg/responses"
	"blogs/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	notificationService service.NotificationService
	redisRepo           repository.RedisRepository
}

func NewController(notificationService service.NotificationService, redisRepo repository.RedisRepository) *Controller {
	return &Controller{
		notificationService: notificationService,
		redisRepo:           redisRepo,
	}
}

func (ctrl *Controller) List(c *gin.Context) {
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}
	var req request.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responses.BadRequest(c, "参数错误")
		return
	}

	notifications, total, err := ctrl.notificationService.ListNotifications(uid, req.Page, req.PageSize)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取通知消息列表失败")
		return
	}
	responses.Success(c, responses.NewPageResponse(notifications, total, req.Page, req.PageSize))
}

func (ctrl *Controller) MarkAsRead(c *gin.Context) {
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	if err := ctrl.notificationService.MarkAsRead(id, uid); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "标记消息已读失败")
		return
	}

	responses.Success(c, nil)
}

func (ctrl *Controller) MarkAllAsRead(c *gin.Context) {
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}

	if err := ctrl.notificationService.MarkAllAsRead(uid); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "全部标记已读失败")
		return
	}

	responses.Success(c, nil)
}

func (ctrl *Controller) GetUnreadCount(c *gin.Context) {
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}

	count, err := ctrl.notificationService.GetUnreadCount(uid)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取未读数失败")
		return
	}

	responses.Success(c, gin.H{"count": count})
}

func (ctrl *Controller) Delete(c *gin.Context) {
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}

	if err := ctrl.notificationService.DeleteNotification(id, uid); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "删除通知失败")
		return
	}

	responses.Success(c, nil)
}
