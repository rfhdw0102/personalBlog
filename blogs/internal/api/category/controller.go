package category

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
	redisRepo       repository.RedisRepository
	categoryService service.CategoryService
}

func NewController(redisRepo repository.RedisRepository, categoryService service.CategoryService) *Controller {
	return &Controller{redisRepo: redisRepo, categoryService: categoryService}
}

func (ctrl *Controller) GetCategories(c *gin.Context) {
	categories, err := ctrl.categoryService.ListCategories()
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取分类失败")
		return
	}
	responses.Success(c, categories)
}

func (ctrl *Controller) GetCategory(c *gin.Context) {
	var req request.CategoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responses.BadRequest(c, "请求参数错误")
		return
	}
	categories, total, err := ctrl.categoryService.GetCategoryList(req)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取分类列表失败")
		return
	}
	responses.Success(c, responses.NewPageResponse(categories, total, req.Page, req.PageSize))
}

func (ctrl *Controller) GetCategoryOne(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	category, err := ctrl.categoryService.GetCategoryByID(id)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取分类失败")
		return
	}
	responses.Success(c, category)
}

func (ctrl *Controller) CreateCategory(c *gin.Context) {
	var category entity.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		responses.BadRequest(c, "请求参数错误")
		return
	}
	if category.Name == "" {
		responses.BadRequest(c, "分类名称不能为空")
		return
	}
	if err := ctrl.categoryService.CreateCategory(&category); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "创建分类失败")
		return
	}
	responses.Success(c, category)
}

func (ctrl *Controller) DeleteCategory(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	if err := ctrl.categoryService.DeleteCategory(id); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "删除分类失败")
		return
	}
	responses.Success(c, nil)
}

func (ctrl *Controller) UpdateCategory(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	var category entity.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		responses.BadRequest(c, "请求参数错误")
		return
	}
	category.ID = id
	if err := ctrl.categoryService.UpdateCategory(&category); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "更新分类失败")
		return
	}
	responses.Success(c, category)
}
