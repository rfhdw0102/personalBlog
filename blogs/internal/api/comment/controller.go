package comment

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/entity"
	"blogs/internal/repository"
	"blogs/internal/service"
	appErrors "blogs/pkg/errors"
	"blogs/pkg/responses"
	"blogs/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	redisRepo      repository.RedisRepository
	commentService service.CommentService
}

func NewController(redisRepo repository.RedisRepository, commentService service.CommentService) *Controller {
	return &Controller{redisRepo: redisRepo, commentService: commentService}
}

func (ctrl *Controller) GetCommentsByArticle(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	comments, err := ctrl.commentService.GetComments(id, 1)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取评论失败")
		return
	}
	responses.Success(c, comments)
}

func (ctrl *Controller) GetCommentsByArticleHide(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	comments, err := ctrl.commentService.GetComments(id, 0)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取评论失败")
		return
	}
	responses.Success(c, comments)
}

func (ctrl *Controller) AddComment(c *gin.Context) {
	var req request.CommentSubmit
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "参数错误")
		return
	}

	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}

	comment := &entity.Comment{
		ArticleID: req.ArticleID,
		UserID:    uid,
		Content:   req.Content,
		ParentID:  req.ParentID,
		Status:    1,
	}

	if err := ctrl.commentService.AddComment(comment); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "发表评论失败")
		return
	}

	responses.Success(c, comment)
}

func (ctrl *Controller) HideComment(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}
	if err := ctrl.commentService.HideComment(id, uid, 0); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "隐藏评论失败")
		return
	}
	responses.Success(c, "成功")
}

func (ctrl *Controller) UnHideComment(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}
	if err := ctrl.commentService.HideComment(id, uid, 1); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "隐藏评论失败")
		return
	}
	responses.Success(c, "成功")
}

func (ctrl *Controller) DeleteComment(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}
	if err := ctrl.commentService.DeleteComment(id, uid); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "删除评论失败")
		return
	}
	responses.Success(c, "成功")
}

func (ctrl *Controller) ListUserCommentsGet(c *gin.Context) {
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}

	var req request.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responses.BadRequest(c, "参数错误")
		return
	}

	comments, total, err := ctrl.commentService.ListUserCommentsGet(uid, req.Page, req.PageSize, 1)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取评论列表失败")
		return
	}

	responses.Success(c, responses.NewPageResponse(comments, total, req.Page, req.PageSize))
}

func (ctrl *Controller) ListUserCommentsTake(c *gin.Context) {
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}

	var req request.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responses.BadRequest(c, "参数错误")
		return
	}

	comments, total, err := ctrl.commentService.ListUserCommentsTake(uid, req.Page, req.PageSize)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取评论列表失败")
		return
	}

	responses.Success(c, responses.NewPageResponse(comments, total, req.Page, req.PageSize))
}

func (ctrl *Controller) ListUserCommentsHide(c *gin.Context) {
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}

	var req request.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responses.BadRequest(c, "参数错误")
		return
	}

	comments, total, err := ctrl.commentService.ListUserCommentsGet(uid, req.Page, req.PageSize, 0)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取评论列表失败")
		return
	}

	responses.Success(c, responses.NewPageResponse(comments, total, req.Page, req.PageSize))
}
