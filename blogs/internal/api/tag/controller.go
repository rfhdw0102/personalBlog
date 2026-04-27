package tag

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/entity"
	"blogs/internal/repository"
	"blogs/internal/service"
	appErrors "blogs/pkg/errors"
	"blogs/pkg/logger"
	"blogs/pkg/responses"
	"blogs/pkg/utils"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	redisRepo  repository.RedisRepository
	tagService service.TagService
}

func NewController(redisRepo repository.RedisRepository, tagService service.TagService) *Controller {
	return &Controller{redisRepo: redisRepo, tagService: tagService}
}

func (ctrl *Controller) GetTags(c *gin.Context) {
	tags, err := ctrl.tagService.ListTags()
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取标签失败")
		return
	}
	responses.Success(c, tags)
}

func (ctrl *Controller) GetTag(c *gin.Context) {
	var req request.TagRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responses.BadRequest(c, "请求参数错误")
		return
	}
	tags, total, err := ctrl.tagService.GetTagList(req)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取标签列表失败")
		return
	}
	responses.Success(c, responses.NewPageResponse(tags, total, req.Page, req.PageSize))
}

func (ctrl *Controller) CreateTag(c *gin.Context) {
	var tag entity.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		responses.BadRequest(c, "请求参数错误")
		return
	}
	if tag.Name == "" {
		responses.BadRequest(c, "标签名称不能为空")
		return
	}
	if err := ctrl.tagService.CreateTag(tag); err != nil {
		logger.Error("=======创建标签失败=======", zap.Error(err))
		responses.FromError(c, err, appErrors.CodeInternalError, "创建标签失败")
		return
	}
	responses.Success(c, tag)
}

func (ctrl *Controller) GetTagOne(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	tag, err := ctrl.tagService.GetTagByID(id)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取标签失败")
		return
	}
	responses.Success(c, tag)
}

func (ctrl *Controller) DeleteTag(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	if err := ctrl.tagService.DeleteTag(id); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "删除标签失败")
		return
	}
	responses.Success(c, nil)
}

func (ctrl *Controller) UpdateTag(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	var tag entity.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		responses.BadRequest(c, "请求参数错误")
		return
	}
	tag.ID = id
	if err := ctrl.tagService.UpdateTag(tag); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "更新标签失败")
		return
	}
	responses.Success(c, tag)
}
