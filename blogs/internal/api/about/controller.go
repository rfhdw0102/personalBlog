package about

import (
	"blogs/internal/model/entity"
	"blogs/internal/repository"
	"blogs/internal/service"
	appErrors "blogs/pkg/errors"
	"blogs/pkg/responses"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	redisRepo    repository.RedisRepository
	aboutService service.AboutService
}

func NewController(redisRepo repository.RedisRepository, aboutService service.AboutService) *Controller {
	return &Controller{redisRepo: redisRepo, aboutService: aboutService}
}

func (ctrl *Controller) Get(c *gin.Context) {
	about, err := ctrl.aboutService.Get()
	if err != nil {
		responses.FromError(c, err, appErrors.CodeNotFound, "获取关于页面信息失败")
		return
	}
	responses.Success(c, about)
}

func (ctrl *Controller) Update(c *gin.Context) {
	var about entity.About
	if err := c.ShouldBindJSON(&about); err != nil {
		responses.BadRequest(c, "参数错误")
		return
	}
	if err := ctrl.aboutService.Update(&about); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "更新关于页面信息失败")
		return
	}
	responses.Success(c, about)
}
